package oidc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/actions/object"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/execution"
	"github.com/zitadel/zitadel/internal/query"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (s *Server) UserInfo(ctx context.Context, r *op.Request[oidc.UserInfoRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()

	features := authz.GetFeatures(ctx)
	if features.LegacyIntrospection {
		return s.LegacyServer.UserInfo(ctx, r)
	}
	if features.TriggerIntrospectionProjections {
		query.TriggerOIDCUserInfoProjections(ctx)
	}

	token, err := s.verifyAccessToken(ctx, r.Data.AccessToken)
	if err != nil {
		return nil, op.NewStatusError(oidc.ErrAccessDenied().WithDescription("access token invalid").WithParent(err).WithReturnParentToClient(authz.GetFeatures(ctx).DebugOIDCParentError), http.StatusUnauthorized)
	}

	var (
		projectID string
		assertion bool
	)
	if token.clientID != "" {
		projectID, assertion, err = s.query.GetOIDCUserinfoClientByID(ctx, token.clientID)
		// token.clientID might contain a username (e.g. client credentials) -> ignore the not found
		if err != nil && !zerrors.IsNotFound(err) {
			return nil, err
		}
	}

	userInfo, err := s.userInfo(
		token.userID,
		token.scope,
		projectID,
		assertion,
		true,
		false,
	)(ctx, true, domain.TriggerTypePreUserinfoCreation)
	if err != nil {
		if !zerrors.IsNotFound(err) {
			return nil, err
		}
		return nil, op.NewStatusError(oidc.ErrAccessDenied().WithDescription("no active user").WithParent(err).WithReturnParentToClient(authz.GetFeatures(ctx).DebugOIDCParentError), http.StatusUnauthorized)
	}
	return op.NewResponse(userInfo), nil
}

// userInfo gets the user's data based on the scope.
// The returned UserInfo contains standard and reserved claims, documented
// here: https://zitadel.com/docs/apis/openidoauth/claims.
//
// User information is only retrieved once from the database.
// However, each time, role claims are asserted and also action flows will trigger.
//
// projectID is an optional parameter which defines the default audience when there are any (or all) role claims requested.
// projectRoleAssertion sets the default of returning all project roles, only if no specific roles were requested in the scope.
// roleAssertion decides whether the roles will be returned (in the token or response)
// userInfoAssertion decides whether the user information (profile data like name, email, ...) are returned
//
// currentProjectOnly can be set to use the current project ID only and ignore the audience from the scope.
// It should be set in cases where the client doesn't need to know roles outside its own project,
// for example an introspection client.
func (s *Server) userInfo(
	userID string,
	scope []string,
	projectID string,
	projectRoleAssertion, userInfoAssertion, currentProjectOnly bool,
) func(ctx context.Context, roleAssertion bool, triggerType domain.TriggerType) (_ *oidc.UserInfo, err error) {
	var (
		once                         sync.Once
		rawUserInfo                  *oidc.UserInfo
		qu                           *query.OIDCUserInfo
		roleAudience, requestedRoles []string
	)
	return func(ctx context.Context, roleAssertion bool, triggerType domain.TriggerType) (_ *oidc.UserInfo, err error) {
		once.Do(func() {
			ctx, span := tracing.NewSpan(ctx)
			defer func() { span.EndWithError(err) }()

			roleAudience, requestedRoles = prepareRoles(ctx, scope, projectID, projectRoleAssertion, currentProjectOnly)
			roleOrgIDs := domain.RoleOrgIDsFromScope(scope)
			qu, err = s.query.GetOIDCUserInfo(ctx, userID, roleAudience, roleOrgIDs...)
			if err != nil {
				return
			}
			rawUserInfo = userInfoToOIDC(qu, userInfoAssertion, scope, s.assetAPIPrefix(ctx))
		})
		if err != nil {
			return nil, err
		}
		// copy the userinfo to make sure the assert roles and actions use their own copy (e.g. map)
		userInfo := &oidc.UserInfo{
			Subject:         rawUserInfo.Subject,
			UserInfoProfile: rawUserInfo.UserInfoProfile,
			UserInfoEmail:   rawUserInfo.UserInfoEmail,
			UserInfoPhone:   rawUserInfo.UserInfoPhone,
			Address:         rawUserInfo.Address,
			Claims:          maps.Clone(rawUserInfo.Claims),
		}
		assertRoles(projectID, qu, roleAudience, requestedRoles, roleAssertion, userInfo)
		return userInfo, s.userinfoFlows(ctx, qu, userInfo, triggerType)
	}
}

// prepareRoles scans the requested scopes and builds the requested roles
// and the audience for which roles need to be asserted.
//
// Scopes with [ScopeProjectRolePrefix] are added to requestedRoles.
// When [ScopeProjectsRoles] is present project IDs with the [domain.ProjectIDScope]
// prefix are added to the returned audience.
//
// If projectRoleAssertion is true and there were no specific roles requested,
// the current projectID will always be parts of the returned audience.
func prepareRoles(ctx context.Context, scope []string, projectID string, projectRoleAssertion, currentProjectOnly bool) (roleAudience, requestedRoles []string) {
	for _, s := range scope {
		if role, ok := strings.CutPrefix(s, ScopeProjectRolePrefix); ok {
			requestedRoles = append(requestedRoles, role)
		}
	}

	// If roles are requested take the audience for those from the scopes,
	// when currentProjectOnly is not set.
	if !currentProjectOnly && (len(requestedRoles) > 0 || slices.Contains(scope, ScopeProjectsRoles)) {
		roleAudience = domain.AddAudScopeToAudience(ctx, roleAudience, scope)
	}

	// When either:
	// - Project role assertion is set;
	// - Roles for the current project (only) are requested;
	// - There is already a roleAudience requested through scope;
	// - There are requested roles through the scope;
	// and the projectID is not empty, projectID must be part of the roleAudience.
	if (projectRoleAssertion || currentProjectOnly || len(roleAudience) > 0 || len(requestedRoles) > 0) && projectID != "" && !slices.Contains(roleAudience, projectID) {
		roleAudience = append(roleAudience, projectID)
	}
	return roleAudience, requestedRoles
}

func userInfoToOIDC(user *query.OIDCUserInfo, userInfoAssertion bool, scope []string, assetPrefix string) *oidc.UserInfo {
	out := &oidc.UserInfo{
		Subject: user.User.ID,
	}
	for _, s := range scope {
		switch s {
		case oidc.ScopeEmail:
			if !userInfoAssertion {
				continue
			}
			out.UserInfoEmail = userInfoEmailToOIDC(user.User)
		case oidc.ScopeProfile:
			if !userInfoAssertion {
				continue
			}
			out.UserInfoProfile = userInfoProfileToOidc(user.User, assetPrefix)
		case oidc.ScopePhone:
			if !userInfoAssertion {
				continue
			}
			out.UserInfoPhone = userInfoPhoneToOIDC(user.User)
		case oidc.ScopeAddress:
			if !userInfoAssertion {
				continue
			}
			// TODO: handle address for human users as soon as implemented
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
	return out
}

func assertRoles(projectID string, user *query.OIDCUserInfo, roleAudience, requestedRoles []string, assertion bool, info *oidc.UserInfo) {
	if !assertion {
		return
	}
	// prevent returning obtained grants if none where requested
	if (projectID != "" && len(requestedRoles) > 0) || len(roleAudience) > 0 {
		setUserInfoRoleClaims(info, newProjectRoles(projectID, user.UserGrants, requestedRoles))
	}
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
		out.AppendClaims(ClaimResourceOwnerID, org.ID)
		out.AppendClaims(ClaimResourceOwnerName, org.Name)
		out.AppendClaims(ClaimResourceOwnerPrimaryDomain, org.PrimaryDomain)
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

func (s *Server) userinfoFlows(ctx context.Context, qu *query.OIDCUserInfo, userInfo *oidc.UserInfo, triggerType domain.TriggerType) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	queriedActions, err := s.query.GetActiveActionsByFlowAndTriggerType(ctx, domain.FlowTypeCustomiseToken, triggerType, qu.User.ResourceOwner)
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
					return object.UserGrantsFromSlice(ctx, s.query, c, qu.UserGrants)
				}),
			),
			actions.SetFields("org",
				actions.SetFields("getMetadata", func(c *actions.FieldConfig) interface{} {
					return func(goja.FunctionCall) goja.Value {
						return object.GetOrganizationMetadata(ctx, s.query, c, qu.User.ResourceOwner)
					}
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
						if strings.HasPrefix(key, ClaimPrefix) {
							return
						}
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
						if strings.HasPrefix(key, ClaimPrefix) {
							return
						}
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

	var function string
	switch triggerType {
	case domain.TriggerTypePreUserinfoCreation:
		function = exec_repo.ID(domain.ExecutionTypeFunction, domain.ActionFunctionPreUserinfo.LocalizationKey())
	case domain.TriggerTypePreAccessTokenCreation:
		function = exec_repo.ID(domain.ExecutionTypeFunction, domain.ActionFunctionPreAccessToken.LocalizationKey())
	case domain.TriggerTypeUnspecified, domain.TriggerTypePostAuthentication, domain.TriggerTypePreCreation, domain.TriggerTypePostCreation, domain.TriggerTypePreSAMLResponseCreation:
		// added for linting, there should never be any trigger type be used here besides PreUserinfo and PreAccessToken
		return err
	}

	if function == "" {
		return nil
	}
	executionTargets, err := execution.QueryExecutionTargetsForFunction(ctx, s.query, function)
	if err != nil {
		return err
	}
	info := &ContextInfo{
		Function:     function,
		UserInfo:     userInfo,
		User:         qu.User,
		UserMetadata: qu.Metadata,
		Org:          qu.Org,
		UserGrants:   qu.UserGrants,
	}

	resp, err := execution.CallTargets(ctx, executionTargets, info)
	if err != nil {
		return err
	}
	contextInfoResponse, ok := resp.(*ContextInfoResponse)
	if !ok || contextInfoResponse == nil {
		return nil
	}
	claimLogs := make([]string, 0)
	for _, metadata := range contextInfoResponse.SetUserMetadata {
		if _, err = s.command.SetUserMetadata(ctx, metadata, userInfo.Subject, qu.User.ResourceOwner); err != nil {
			claimLogs = append(claimLogs, fmt.Sprintf("failed to set user metadata key %q", metadata.Key))
		}
	}
	for _, claim := range contextInfoResponse.AppendClaims {
		if strings.HasPrefix(claim.Key, ClaimPrefix) {
			continue
		}
		if userInfo.Claims[claim.Key] == nil {
			userInfo.AppendClaims(claim.Key, claim.Value)
			continue
		}
		claimLogs = append(claimLogs, fmt.Sprintf("key %q already exists", claim.Key))
	}
	claimLogs = append(claimLogs, contextInfoResponse.AppendLogClaims...)
	if len(claimLogs) > 0 {
		userInfo.AppendClaims(fmt.Sprintf(ClaimActionLogFormat, function), claimLogs)
	}

	return nil
}

type ContextInfo struct {
	Function     string               `json:"function,omitempty"`
	UserInfo     *oidc.UserInfo       `json:"userinfo,omitempty"`
	User         *query.User          `json:"user,omitempty"`
	UserMetadata []query.UserMetadata `json:"user_metadata,omitempty"`
	Org          *query.UserInfoOrg   `json:"org,omitempty"`
	UserGrants   []query.UserGrant    `json:"user_grants,omitempty"`
	Response     *ContextInfoResponse `json:"response,omitempty"`
}

type ContextInfoResponse struct {
	SetUserMetadata []*domain.Metadata `json:"set_user_metadata,omitempty"`
	AppendClaims    []*AppendClaim     `json:"append_claims,omitempty"`
	AppendLogClaims []string           `json:"append_log_claims,omitempty"`
}

type AppendClaim struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func (c *ContextInfo) GetHTTPRequestBody() []byte {
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return data
}

func (c *ContextInfo) SetHTTPResponseBody(resp []byte) error {
	if !json.Valid(resp) {
		return zerrors.ThrowPreconditionFailed(nil, "ACTION-4m9s2", "Errors.Execution.ResponseIsNotValidJSON")
	}
	if c.Response == nil {
		c.Response = &ContextInfoResponse{}
	}
	return json.Unmarshal(resp, c.Response)
}

func (c *ContextInfo) GetContent() any {
	return c.Response
}
