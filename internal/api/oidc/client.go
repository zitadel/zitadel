package oidc

import (
	"context"
	"strings"

	"golang.org/x/text/language"
	"gopkg.in/square/go-jose.v2"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/op"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	proj_model "github.com/caos/zitadel/internal/project/model"
	user_model "github.com/caos/zitadel/internal/user/model"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
)

const (
	scopeOpenID  = "openid"
	scopeProfile = "profile"
	scopeEmail   = "email"
	scopePhone   = "phone"
	scopeAddress = "address"

	ScopeProjectRolePrefix = "urn:zitadel:iam:org:project:role:"
	ClaimProjectRoles      = "urn:zitadel:iam:org:project:roles"

	oidcCtx = "oidc"
)

func (o *OPStorage) GetClientByClientID(ctx context.Context, id string) (op.Client, error) {
	client, err := o.repo.ApplicationByClientID(ctx, id)
	if err != nil {
		return nil, err
	}
	if client.State != proj_model.AppStateActive {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-sdaGg", "client is not active")
	}
	projectRoles, err := o.repo.ProjectRolesByProjectID(client.ProjectID)
	if err != nil {
		return nil, err
	}
	allowedScopes := make([]string, len(projectRoles))
	for i, role := range projectRoles {
		allowedScopes[i] = ScopeProjectRolePrefix + role.Key
	}
	return ClientFromBusiness(client, o.defaultLoginURL, o.defaultAccessTokenLifetime, o.defaultIdTokenLifetime, allowedScopes)
}

func (o *OPStorage) GetKeyByIDAndUserID(ctx context.Context, keyID, userID string) (*jose.JSONWebKey, error) {
	key, err := o.repo.MachineKeyByID(ctx, keyID)
	if err != nil {
		return nil, err
	}
	if key.UserID != userID {
		return nil, errors.ThrowPermissionDenied(nil, "OIDC-24jm3", "key from different user")
	}
	publicKey, err := crypto.BytesToPublicKey(key.PublicKey)
	if err != nil {
		return nil, err
	}
	return &jose.JSONWebKey{
		KeyID: key.ID,
		Use:   "sig",
		Key:   publicKey,
	}, nil
}

func (o *OPStorage) AuthorizeClientIDSecret(ctx context.Context, id string, secret string) error {
	ctx = authz.SetCtxData(ctx, authz.CtxData{
		UserID: oidcCtx,
		OrgID:  oidcCtx,
	})
	return o.repo.AuthorizeOIDCApplication(ctx, id, secret)
}

func (o *OPStorage) GetUserinfoFromToken(ctx context.Context, tokenID, origin string) (oidc.UserInfo, error) {
	token, err := o.repo.ValidTokenByID(ctx, tokenID)
	if err != nil {
		return nil, errors.ThrowPermissionDenied(nil, "OIDC-Dsfb2", "token is not valid or has expired")
	}
	if token.ApplicationID != "" {
		app, err := o.repo.ApplicationByClientID(ctx, token.ApplicationID)
		if err != nil {
			return nil, err
		}
		if origin != "" && !http.IsOriginAllowed(app.OriginAllowList, origin) {
			return nil, errors.ThrowPermissionDenied(nil, "OIDC-da1f3", "origin is not allowed")
		}
	}
	return o.GetUserinfoFromScopes(ctx, token.UserID, token.ApplicationID, token.Scopes)
}

func (o *OPStorage) GetUserinfoFromScopes(ctx context.Context, userID, applicationID string, scopes []string) (oidc.UserInfo, error) {
	user, err := o.repo.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	userInfo := oidc.NewUserInfo()
	roles := make([]string, 0)
	for _, scope := range scopes {
		switch scope {
		case oidc.ScopeOpenID:
			userInfo.SetSubject(user.ID)
		case oidc.ScopeEmail:
			if user.HumanView == nil {
				continue
			}
			userInfo.SetEmail(user.Email, user.IsEmailVerified)
		case oidc.ScopeProfile:
			userInfo.SetPreferredUsername(user.PreferredLoginName)
			userInfo.SetUpdatedAt(user.ChangeDate)
			if user.HumanView != nil {
				userInfo.SetName(user.DisplayName)
				userInfo.SetFamilyName(user.LastName)
				userInfo.SetGivenName(user.FirstName)
				userInfo.SetNickname(user.NickName)
				userInfo.SetGender(oidc.Gender(getGender(user.Gender)))
				locale, _ := language.Parse(user.PreferredLanguage)
				userInfo.SetLocale(locale)
			} else {
				userInfo.SetName(user.MachineView.Name)
			}
		case oidc.ScopePhone:
			if user.HumanView == nil {
				continue
			}
			userInfo.SetPhone(user.Phone, user.IsPhoneVerified)
		case oidc.ScopeAddress:
			if user.HumanView == nil {
				continue
			}
			if user.StreetAddress == "" && user.Locality == "" && user.Region == "" && user.PostalCode == "" && user.Country == "" {
				continue
			}
			userInfo.SetAddress(oidc.NewUserInfoAddress(user.StreetAddress, user.Locality, user.Region, user.PostalCode, user.Country, ""))
		default:
			if strings.HasPrefix(scope, ScopeProjectRolePrefix) {
				roles = append(roles, strings.TrimPrefix(scope, ScopeProjectRolePrefix))
			}
		}
	}

	if len(roles) == 0 || applicationID == "" {
		return userInfo, nil
	}
	projectRoles, err := o.assertRoles(ctx, userID, applicationID, roles)
	if err != nil {
		return nil, err
	}
	if len(projectRoles) > 0 {
		userInfo.AppendClaims(ClaimProjectRoles, projectRoles)
	}

	return userInfo, nil
}

func (o *OPStorage) GetPrivateClaimsFromScopes(ctx context.Context, userID, applicationID string, scopes []string) (claims map[string]interface{}, err error) {
	roles := make([]string, 0)
	for _, scope := range scopes {
		if strings.HasPrefix(scope, ScopeProjectRolePrefix) {
			roles = append(roles, strings.TrimPrefix(scope, ScopeProjectRolePrefix))
		}
	}
	if len(roles) == 0 || applicationID == "" {
		return nil, nil
	}
	projectRoles, err := o.assertRoles(ctx, userID, applicationID, roles)
	if err != nil {
		return nil, err
	}
	if len(projectRoles) > 0 {
		claims = map[string]interface{}{ClaimProjectRoles: projectRoles}
	}
	return claims, err
}

func (o *OPStorage) assertRoles(ctx context.Context, userID, applicationID string, requestedRoles []string) (map[string]map[string]string, error) {
	app, err := o.repo.ApplicationByClientID(ctx, applicationID)
	if err != nil {
		return nil, err
	}
	grants, err := o.repo.UserGrantsByProjectAndUserID(app.ProjectID, userID)
	if err != nil {
		return nil, err
	}
	projectRoles := make(map[string]map[string]string)
	for _, requestedRole := range requestedRoles {
		for _, grant := range grants {
			checkGrantedRoles(projectRoles, grant, requestedRole)
		}
	}
	return projectRoles, nil
}

func checkGrantedRoles(roles map[string]map[string]string, grant *grant_model.UserGrantView, requestedRole string) {
	for _, grantedRole := range grant.RoleKeys {
		if requestedRole == grantedRole {
			appendRole(roles, grantedRole, grant.ResourceOwner, grant.OrgPrimaryDomain)
		}
	}
}

func appendRole(roles map[string]map[string]string, role, orgID, orgPrimaryDomain string) {
	if roles[role] == nil {
		roles[role] = make(map[string]string, 0)
	}
	roles[role][orgID] = orgPrimaryDomain
}

func getGender(gender user_model.Gender) string {
	switch gender {
	case user_model.GenderFemale:
		return "female"
	case user_model.GenderMale:
		return "male"
	case user_model.GenderDiverse:
		return "diverse"
	}
	return ""
}
