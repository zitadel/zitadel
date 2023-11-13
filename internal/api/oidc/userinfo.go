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
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (s *Server) getUserInfoWithRoles(ctx context.Context, userID, projectID string, scope, roleAudience []string) (_ *oidc.UserInfo, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	userInfoChan := make(chan *userInfoResult)
	go s.getUserInfo(ctx, userID, userInfoChan)

	rolesChan := make(chan *assertRolesResult)
	go s.assertRoles(ctx, userID, projectID, scope, roleAudience, rolesChan)

	var (
		userInfoResult    *userInfoResult
		assertRolesResult *assertRolesResult
	)

	// make sure both channels are always read,
	// and cancel the context on first error
	for i := 0; i < 2; i++ {
		var resErr error

		select {
		case userInfoResult = <-userInfoChan:
			resErr = userInfoResult.err
		case assertRolesResult = <-rolesChan:
			resErr = assertRolesResult.err
		}

		if resErr == nil {
			continue
		}
		cancel()

		// we only care for the first error that occured,
		// as the next error is most probably a context error.
		if err == nil {
			err = resErr
		}
	}

	userInfo := userInfoToOIDC(userInfoResult.userInfo, scope, s.assetAPIPrefix(ctx))
	setUserInfoRoleClaims(userInfo, assertRolesResult.projectsRoles)

	return userInfo, s.userinfoFlows(ctx, userInfoResult.userInfo, assertRolesResult.userGrants, userInfo)
}

type userInfoResult struct {
	userInfo *query.OIDCUserInfo
	err      error
}

func (s *Server) getUserInfo(ctx context.Context, userID string, rc chan<- *userInfoResult) {
	userInfo, err := s.query.GetOIDCUserInfo(ctx, userID)
	rc <- &userInfoResult{
		userInfo: userInfo,
		err:      err,
	}
}

type assertRolesResult struct {
	userGrants    *query.UserGrants
	projectsRoles *projectsRoles
	err           error
}

func (s *Server) assertRoles(ctx context.Context, userID, projectID string, scope, roleAudience []string, rc chan<- *assertRolesResult) {
	userGrands, projectsRoles, err := func() (*query.UserGrants, *projectsRoles, error) {
		// if all roles are requested take the audience for those from the scopes
		if slices.Contains(scope, ScopeProjectsRoles) {
			roleAudience = domain.AddAudScopeToAudience(ctx, roleAudience, scope)
		}

		requestedRoles := make([]string, 0, len(scope))
		for _, s := range scope {
			if role, ok := strings.CutPrefix(s, ScopeProjectRolePrefix); ok {
				requestedRoles = append(requestedRoles, role)
			}
		}

		if len(requestedRoles) == 0 && len(roleAudience) == 0 {
			return nil, nil, nil
		}

		// ensure the projectID of the requesting is part of the roleAudience
		if projectID != "" {
			roleAudience = append(roleAudience, projectID)
		}
		queries := make([]query.SearchQuery, 0, 2)
		projectQuery, err := query.NewUserGrantProjectIDsSearchQuery(roleAudience)
		if err != nil {
			return nil, nil, err
		}
		queries = append(queries, projectQuery)
		userIDQuery, err := query.NewUserGrantUserIDSearchQuery(userID)
		if err != nil {
			return nil, nil, err
		}
		queries = append(queries, userIDQuery)
		grants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{
			Queries: queries,
		}, false, false) // triggers disabled
		if err != nil {
			return nil, nil, err
		}
		roles := new(projectsRoles)
		// if specific roles where requested, check if they are granted and append them in the roles list
		if len(requestedRoles) > 0 {
			for _, requestedRole := range requestedRoles {
				for _, grant := range grants.UserGrants {
					checkGrantedRoles(roles, grant, requestedRole, grant.ProjectID == projectID)
				}
			}
			return grants, roles, nil
		}
		// no specific roles were requested, so convert any grants into roles
		for _, grant := range grants.UserGrants {
			for _, role := range grant.Roles {
				roles.Add(grant.ProjectID, role, grant.ResourceOwner, grant.OrgPrimaryDomain, grant.ProjectID == projectID)
			}
		}
		return grants, roles, nil
	}()

	rc <- &assertRolesResult{
		userGrants:    userGrands,
		projectsRoles: projectsRoles,
		err:           err,
	}
}

func userInfoToOIDC(user *query.OIDCUserInfo, scope []string, assetPrefix string) *oidc.UserInfo {
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
			if strings.HasPrefix(s, domain.OrgDomainPrimaryScope) {
				out.AppendClaims(domain.OrgDomainPrimaryClaim, strings.TrimPrefix(s, domain.OrgDomainPrimaryScope))
			}
			if strings.HasPrefix(s, domain.OrgIDScope) {
				out.AppendClaims(domain.OrgIDClaim, strings.TrimPrefix(s, domain.OrgIDScope))
				setUserInfoOrgClaims(user, out)
			}
		}
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
	return oidc.UserInfoProfile{
		UpdatedAt:         oidc.FromTime(user.ChangeDate),
		PreferredUsername: user.PreferredLoginName,
	}
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

func (s *Server) userinfoFlows(ctx context.Context, user *query.OIDCUserInfo, userGrants *query.UserGrants, userInfo *oidc.UserInfo) error {
	queriedActions, err := s.query.GetActiveActionsByFlowAndTriggerType(ctx, domain.FlowTypeCustomiseToken, domain.TriggerTypePreUserinfoCreation, user.User.ResourceOwner, false)
	if err != nil {
		return err
	}

	ctxFields := actions.SetContextFields(
		actions.SetFields("v1",
			actions.SetFields("claims", userinfoClaims(userInfo)),
			actions.SetFields("getUser", func(c *actions.FieldConfig) interface{} {
				return func(call goja.FunctionCall) goja.Value {
					return object.UserFromQuery(c, user.User)
				}
			}),
			actions.SetFields("user",
				actions.SetFields("getMetadata", func(c *actions.FieldConfig) interface{} {
					return func(goja.FunctionCall) goja.Value {
						return object.UserMetadataListFromSlice(c, user.Metadata)
					}
				}),
				actions.SetFields("grants", func(c *actions.FieldConfig) interface{} {
					return object.UserGrantsFromQuery(c, userGrants)
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
						if _, err = s.command.SetUserMetadata(ctx, metadata, userInfo.Subject, user.User.ResourceOwner); err != nil {
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
