package oidc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/dop251/goja"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/actions/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (s *Server) userInfo(ctx context.Context, userID, projectID string, scope, roleAudience []string) (_ *oidc.UserInfo, err error) {
	roleAudience, requestedRoles := prepareRoles(ctx, projectID, scope, roleAudience)
	qu, err := s.query.GetOIDCUserInfo(ctx, userID, roleAudience)
	if err != nil {
		return nil, err
	}

	userInfo := userInfoToOIDC(projectID, qu, scope, roleAudience, requestedRoles, s.assetAPIPrefix(ctx))
	return userInfo, s.userinfoFlows(ctx, qu, userInfo)
}

// prepareRoles scans the requested scopes, appends to roleAudiendce and returns the requestedRoles.
//
// When [ScopeProjectsRoles] is present and roleAudience was empty,
// project IDs with the [domain.ProjectIDScope] prefix are added to the roleAudience.
//
// Scopes with [ScopeProjectRolePrefix] are added to requestedRoles.
//
// If the resulting requestedRoles or roleAudience are not not empty,
// the current projectID will always be parts or roleAudience.
// Else nil, nil is returned.
func prepareRoles(ctx context.Context, projectID string, scope, roleAudience []string) (ra, requestedRoles []string) {
	// if all roles are requested take the audience for those from the scopes
	if slices.Contains(scope, ScopeProjectsRoles) && len(roleAudience) == 0 {
		roleAudience = domain.AddAudScopeToAudience(ctx, roleAudience, scope)
	}
	requestedRoles = make([]string, 0, len(scope))
	for _, s := range scope {
		if role, ok := strings.CutPrefix(s, ScopeProjectRolePrefix); ok {
			requestedRoles = append(requestedRoles, role)
		}
	}
	if len(requestedRoles) == 0 && len(roleAudience) == 0 {
		return nil, nil
	}

	if projectID != "" && !slices.Contains(roleAudience, projectID) {
		roleAudience = append(roleAudience, projectID)
	}
	return roleAudience, requestedRoles
}

func userInfoToOIDC(projectID string, user *query.OIDCUserInfo, scope, roleAudience, requestedRoles []string, assetPrefix string) *oidc.UserInfo {
	out := new(oidc.UserInfo)
	for _, s := range scope {
		switch s {
		case oidc.ScopeOpenID:
			out.Subject = user.User.ID
		case oidc.ScopeEmail:
			out.UserInfoEmail = userInfoEmailToOIDC(user.User)
		case oidc.ScopeProfile:
			out.UserInfoProfile = userInfoProfileToOidc(user.User, assetPrefix)
		case oidc.ScopePhone:
			out.UserInfoPhone = userInfoPhoneToOIDC(user.User)
		case oidc.ScopeAddress:
			//TODO: handle address for human users as soon as implemented
		case ScopeUserMetaData:
			setUserInfoMetadata(user.Metadata, out)
		case ScopeResourceOwner:
			setUserInfoOrgClaims(user, out)
		default:
			if claim, ok := strings.CutPrefix(s, domain.OrgDomainPrimaryScope); ok {
				out.AppendClaims(domain.OrgDomainPrimaryClaim, claim)
			}
			if claim, ok := strings.CutPrefix(s, domain.OrgIDScope); ok {
				out.AppendClaims(domain.OrgIDClaim, claim)
				setUserInfoOrgClaims(user, out)
			}
		}
	}

	// prevent returning obtained grants if none where requested
	if (projectID != "" && len(requestedRoles) > 0) || len(roleAudience) > 0 {
		setUserInfoRoleClaims(out, newProjectRoles(projectID, user.UserGrants, requestedRoles))
	}
	return out
}

func userInfoEmailToOIDC(user *query.User) oidc.UserInfoEmail {
	if human := user.Human; human != nil {
		return oidc.UserInfoEmail{
			Email:         string(human.Email),
			EmailVerified: oidc.Bool(human.IsEmailVerified),
		}
	}
	return oidc.UserInfoEmail{}
}

func userInfoProfileToOidc(user *query.User, assetPrefix string) oidc.UserInfoProfile {
	if human := user.Human; human != nil {
		return oidc.UserInfoProfile{
			Name:              human.DisplayName,
			GivenName:         human.FirstName,
			FamilyName:        human.LastName,
			Nickname:          human.NickName,
			Picture:           domain.AvatarURL(assetPrefix, user.ResourceOwner, user.Human.AvatarKey),
			Gender:            getGender(human.Gender),
			Locale:            oidc.NewLocale(human.PreferredLanguage),
			UpdatedAt:         oidc.FromTime(user.ChangeDate),
			PreferredUsername: user.PreferredLoginName,
		}
	}
	if machine := user.Machine; machine != nil {
		return oidc.UserInfoProfile{
			Name:              machine.Name,
			UpdatedAt:         oidc.FromTime(user.ChangeDate),
			PreferredUsername: user.PreferredLoginName,
		}
	}
	return oidc.UserInfoProfile{}
}

func userInfoPhoneToOIDC(user *query.User) oidc.UserInfoPhone {
	if human := user.Human; human != nil {
		return oidc.UserInfoPhone{
			PhoneNumber:         string(human.Phone),
			PhoneNumberVerified: human.IsPhoneVerified,
		}
	}
	return oidc.UserInfoPhone{}
}

func setUserInfoMetadata(metadata []query.UserMetadata, out *oidc.UserInfo) {
	if len(metadata) == 0 {
		return
	}
	mdmap := make(map[string]string, len(metadata))
	for _, md := range metadata {
		mdmap[md.Key] = base64.RawURLEncoding.EncodeToString(md.Value)
	}
	out.AppendClaims(ClaimUserMetaData, mdmap)
}

func setUserInfoOrgClaims(user *query.OIDCUserInfo, out *oidc.UserInfo) {
	if org := user.Org; org != nil {
		out.AppendClaims(ClaimResourceOwner+"id", org.ID)
		out.AppendClaims(ClaimResourceOwner+"name", org.Name)
		out.AppendClaims(ClaimResourceOwner+"primary_domain", org.PrimaryDomain)
	}
}

func setUserInfoRoleClaims(userInfo *oidc.UserInfo, roles *projectsRoles) {
	if roles != nil && len(roles.projects) > 0 {
		if roles, ok := roles.projects[roles.requestProjectID]; ok {
			userInfo.AppendClaims(ClaimProjectRoles, roles)
		}
		for projectID, roles := range roles.projects {
			userInfo.AppendClaims(fmt.Sprintf(ClaimProjectRolesFormat, projectID), roles)
		}
	}
}

func (s *Server) userinfoFlows(ctx context.Context, qu *query.OIDCUserInfo, userInfo *oidc.UserInfo) error {
	queriedActions, err := s.query.GetActiveActionsByFlowAndTriggerType(ctx, domain.FlowTypeCustomiseToken, domain.TriggerTypePreUserinfoCreation, qu.User.ResourceOwner)
	if err != nil {
		return err
	}

	ctxFields := actions.SetContextFields(
		actions.SetFields("v1",
			actions.SetFields("claims", userinfoClaims(userInfo)),
			actions.SetFields("getUser", func(c *actions.FieldConfig) interface{} {
				return func(call goja.FunctionCall) goja.Value {
					return object.UserFromQuery(c, qu.User)
				}
			}),
			actions.SetFields("user",
				actions.SetFields("getMetadata", func(c *actions.FieldConfig) interface{} {
					return func(goja.FunctionCall) goja.Value {
						return object.UserMetadataListFromSlice(c, qu.Metadata)
					}
				}),
				actions.SetFields("grants", func(c *actions.FieldConfig) interface{} {
					return object.UserGrantsFromSlice(c, qu.UserGrants)
				}),
			),
		),
	)

	for _, action := range queriedActions {
		actionCtx, cancel := context.WithTimeout(ctx, action.Timeout())
		claimLogs := []string{}

		apiFields := actions.WithAPIFields(
			actions.SetFields("v1",
				actions.SetFields("userinfo",
					actions.SetFields("setClaim", func(key string, value interface{}) {
						if userInfo.Claims[key] == nil {
							userInfo.AppendClaims(key, value)
							return
						}
						claimLogs = append(claimLogs, fmt.Sprintf("key %q already exists", key))
					}),
					actions.SetFields("appendLogIntoClaims", func(entry string) {
						claimLogs = append(claimLogs, entry)
					}),
				),
				actions.SetFields("claims",
					actions.SetFields("setClaim", func(key string, value interface{}) {
						if userInfo.Claims[key] == nil {
							userInfo.AppendClaims(key, value)
							return
						}
						claimLogs = append(claimLogs, fmt.Sprintf("key %q already exists", key))
					}),
					actions.SetFields("appendLogIntoClaims", func(entry string) {
						claimLogs = append(claimLogs, entry)
					}),
				),
				actions.SetFields("user",
					actions.SetFields("setMetadata", func(call goja.FunctionCall) goja.Value {
						if len(call.Arguments) != 2 {
							panic("exactly 2 (key, value) arguments expected")
						}
						key := call.Arguments[0].Export().(string)
						val := call.Arguments[1].Export()

						value, err := json.Marshal(val)
						if err != nil {
							logging.WithError(err).Debug("unable to marshal")
							panic(err)
						}

						metadata := &domain.Metadata{
							Key:   key,
							Value: value,
						}
						if _, err = s.command.SetUserMetadata(ctx, metadata, userInfo.Subject, qu.User.ResourceOwner); err != nil {
							logging.WithError(err).Info("unable to set md in action")
							panic(err)
						}
						return nil
					}),
				),
			),
		)

		err = actions.Run(
			actionCtx,
			ctxFields,
			apiFields,
			action.Script,
			action.Name,
			append(actions.ActionToOptions(action), actions.WithHTTP(actionCtx), actions.WithUUID(actionCtx))...,
		)
		cancel()
		if err != nil {
			return err
		}
		if len(claimLogs) > 0 {
			userInfo.AppendClaims(fmt.Sprintf(ClaimActionLogFormat, action.Name), claimLogs)
		}
	}

	return nil
}
