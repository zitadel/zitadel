package oidc

import (
	"context"

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
)

const (
	scopeOpenID  = "openid"
	scopeProfile = "profile"
	scopeEmail   = "email"
	scopePhone   = "phone"
	scopeAddress = "address"

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
	return ClientFromBusiness(client, o.defaultLoginURL, o.defaultAccessTokenLifetime, o.defaultIdTokenLifetime)
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

func (o *OPStorage) GetUserinfoFromToken(ctx context.Context, tokenID, origin string) (*oidc.Userinfo, error) {
	token, err := o.repo.TokenByID(ctx, tokenID)
	if err != nil {
		return nil, err
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
	return o.GetUserinfoFromScopes(ctx, token.UserID, token.Scopes)
}

func (o *OPStorage) GetUserinfoFromScopes(ctx context.Context, userID string, scopes []string) (*oidc.Userinfo, error) {
	user, err := o.repo.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	userInfo := new(oidc.Userinfo)
	for _, scope := range scopes {
		switch scope {
		case scopeOpenID:
			userInfo.Subject = user.ID
		case scopeEmail:
			if user.HumanView == nil {
				continue
			}
			userInfo.Email = user.Email
			userInfo.EmailVerified = user.IsEmailVerified
		case scopeProfile:
			userInfo.PreferredUsername = user.PreferredLoginName
			userInfo.UpdatedAt = user.ChangeDate
			if user.HumanView != nil {
				userInfo.Name = user.DisplayName
				userInfo.FamilyName = user.LastName
				userInfo.GivenName = user.FirstName
				userInfo.Nickname = user.NickName
				userInfo.Gender = oidc.Gender(getGender(user.Gender))
				userInfo.Locale, err = language.Parse(user.PreferredLanguage)
			} else {
				userInfo.Name = user.MachineView.Name
			}
		case scopePhone:
			if user.HumanView == nil {
				continue
			}
			userInfo.PhoneNumber = user.Phone
			userInfo.PhoneNumberVerified = user.IsPhoneVerified
		case scopeAddress:
			if user.HumanView == nil {
				continue
			}
			userInfo.Address.StreetAddress = user.StreetAddress
			userInfo.Address.Locality = user.Locality
			userInfo.Address.Region = user.Region
			userInfo.Address.PostalCode = user.PostalCode
			userInfo.Address.Country = user.Country
		default:
			userInfo.Authorizations = append(userInfo.Authorizations, scope)
		}
	}
	return userInfo, nil
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
