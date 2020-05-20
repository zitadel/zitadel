package oidc

import (
	"context"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/op"

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
	client, err := o.repo.ApplicationByClientID(ctx, id)
	if err != nil {
		return nil, err
	}
	return ClientFromBusiness(client, o.defaultLoginURL, o.tokenLifetime)
}

func (o *OPStorage) AuthorizeClientIDSecret(ctx context.Context, id string, secret string) error {
	return o.repo.AuthorizeOIDCApplication(ctx, id, secret)
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
	userInfo := new(oidc.Userinfo)
	for _, scope := range scopes {
		switch scope {
		case scopeOpenID:
			userInfo.Subject = user.AggregateID
		case scopeEmail:
			if user.Email == nil {
				continue
			}
			userInfo.Email = user.EmailAddress
			userInfo.EmailVerified = user.IsEmailVerified
		case scopeProfile:
			if user.Profile == nil {
				continue
			}
			userInfo.Name = user.FirstName + " " + user.LastName
			userInfo.FamilyName = user.LastName
			userInfo.GivenName = user.FirstName
			userInfo.Nickname = user.NickName
			userInfo.PreferredUsername = user.UserName
			userInfo.UpdatedAt = user.ChangeDate
			userInfo.Gender = oidc.Gender(getGender(user.Gender))
		case scopePhone:
			if user.Phone == nil {
				continue
			}
			userInfo.PhoneNumber = user.PhoneNumber
			userInfo.PhoneNumberVerified = user.IsPhoneVerified
		case scopeAddress:
			if user.Address == nil {
				continue
			}
			userInfo.Address.StreetAddress = user.StreetAddress
			userInfo.Address.Locality = user.Locality
			userInfo.Address.Region = user.Region
			userInfo.Address.PostalCode = user.PostalCode
			userInfo.Address.Country = user.Country
		}
	}
	return userInfo, nil
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
