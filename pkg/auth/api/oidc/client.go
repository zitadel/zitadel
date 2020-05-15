package oidc

import (
	"context"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/op"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/user/model"
)

const (
	scopeOpenID  = "openid"
	scopeProfile = "profile"
	scopeEmail   = "email"
	scopePhone   = "phone"
	scopeAddress = "address"
)

func (o *OPStorage) GetClientByClientID(ctx context.Context, id string) (op.Client, error) {
	apps, err := o.processor.SearchApplications(ctx, &model.ApplicationSearchRequest{
		Limit: 1,
		Queries: []*model.ApplicationSearchQuery{
			&model.ApplicationSearchQuery{Key: model.OIDCClientID, Method: model.Contains, Value: id},
		},
	})
	if err != nil || apps.TotalResult != 1 {
		return nil, caos_errs.ThrowNotFound(err, "OIDC-idixXj", "client not found")
	}
	if apps.Result[0].AppState != model.AppStateACTIVE {
		return nil, caos_errs.ThrowNotFound(err, "OIDC-obFvv", "inactive application")
	}
	return ClientFromBusiness(apps.Result[0], o.defaultLoginURL, o.tokenLifetime)
}

func (o *OPStorage) AuthorizeClientIDSecret(ctx context.Context, id string, secret string) error {
	_, err := o.processor.AuthorizeApplication(ctx, &model.AuthorizeApplication{
		AppType: model.AppTypeOIDC, OIDCAuth: &model.OIDCAuthorization{ClientID: id, ClientSecret: secret},
	})
	return err
}

func (o *OPStorage) GetUserinfoFromToken(ctx context.Context, tokenID string) (*oidc.Userinfo, error) {
	token, err := o.repo.TokenByID(ctx, tokenID)
	if err != nil {
		return nil, err
	}
	return o.GetUserinfoFromScopes(ctx, token.UserID, token.Scopes)
}

func (o *OPStorage) GetUserinfoFromScopes(ctx context.Context, userID string, scopes []string) (*oidc.Userinfo, error) {
	user, err := o.repo.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	userinfo := new(oidc.Userinfo)
	for _, scope := range scopes {
		switch scope {
		case scopeOpenID:
			userinfo.Subject = user.AggregateID
		case scopeEmail:
			if user.Email == nil {
				continue
			}
			userinfo.Email = user.EmailAddress
			userinfo.EmailVerified = user.IsEmailVerified
		case scopeProfile:
			if user.Profile == nil {
				continue
			}
			userinfo.Name = user.FirstName + " " + user.LastName
			userinfo.FamilyName = user.LastName
			userinfo.GivenName = user.FirstName
			userinfo.Nickname = user.NickName
			userinfo.PreferredUsername = user.UserName
			userinfo.UpdatedAt = user.ChangeDate
			userinfo.Gender = oidc.Gender(getGender(user.Gender))
		case scopePhone:
			if user.Phone == nil {
				continue
			}
			userinfo.PhoneNumber = user.PhoneNumber
			userinfo.PhoneNumberVerified = user.IsPhoneVerified
		case scopeAddress:
			if user.Address == nil {
				continue
			}
			userinfo.Address.StreetAddress = user.StreetAddress
			userinfo.Address.Locality = user.Locality
			userinfo.Address.Region = user.Region
			userinfo.Address.PostalCode = user.PostalCode
			userinfo.Address.Country = user.Country
		}
	}
	return userinfo, nil
}

func getGender(gender model.Gender) string {
	switch gender {
	case model.GENDER_FEMALE:
		return "female"
	case model.GENDER_MALE:
		return "male"
	case model.GENDER_DIVERSE:
		return "diverse"
	}
	return ""
}
