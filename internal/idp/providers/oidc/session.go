package oidc

import (
	"context"

	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Session = (*Session)(nil)

// Session is the idp.Session implementation for the OIDC provider
type Session struct {
	Provider *Provider
	AuthURL  string
	Code     string
	Tokens   *oidc.Tokens
}

func (s *Session) GetAuthURL() string {
	return s.AuthURL
}

func (s *Session) FetchUser(ctx context.Context) (user idp.User, err error) {
	if s.Tokens == nil {
		if err = s.authorize(ctx); err != nil {
			return idp.User{}, err
		}
	}
	info, err := rp.Userinfo(
		s.Tokens.AccessToken,
		s.Tokens.TokenType,
		s.Tokens.IDTokenClaims.GetSubject(),
		s.Provider.RelyingParty,
	)
	if err != nil {
		return idp.User{}, err
	}
	userFromClaims(info, &user)
	return user, nil
}

func (s *Session) authorize(ctx context.Context) error {
	if s.Code == "" {
		return ErrCodeMissing
	}
	tokens, err := rp.CodeExchange(ctx, s.Code, s.Provider.RelyingParty)
	if err != nil {
		return err
	}
	s.Tokens = tokens
	return nil
}

// maps the oidc.UserInfo to an idp.User using the default OIDC claims
func userFromClaims(info oidc.UserInfo, user *idp.User) {
	user.ID = info.GetSubject()
	user.FirstName = info.GetGivenName()
	user.LastName = info.GetFamilyName()
	user.DisplayName = info.GetName()
	user.NickName = info.GetNickname()
	user.PreferredUsername = info.GetPreferredUsername()
	user.Email = info.GetEmail()
	user.IsEmailVerified = info.IsEmailVerified()
	user.Phone = info.GetPhoneNumber()
	user.IsPhoneVerified = info.IsPhoneNumberVerified()
	user.PreferredLanguage = info.GetLocale()
	user.AvatarURL = info.GetPicture()
	user.Profile = info.GetProfile()
}
