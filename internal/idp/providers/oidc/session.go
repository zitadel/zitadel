package oidc

import (
	"context"
	"errors"

	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
)

var ErrCodeMissing = errors.New("no auth code provided")

var _ idp.Session = (*Session)(nil)

// Session is the [idp.Session] implementation for the OIDC provider.
type Session struct {
	Provider *Provider
	AuthURL  string
	Code     string
	Tokens   *oidc.Tokens
}

// GetAuthURL implements the [idp.Session] interface.
func (s *Session) GetAuthURL() string {
	return s.AuthURL
}

// FetchUser implements the [idp.Session] interface.
// It will execute an OIDC code exchange if needed to retrieve the tokens,
// call the userinfo endpoint and map the received information into an [idp.User].
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

func (s *Session) authorize(ctx context.Context) (err error) {
	if s.Code == "" {
		return ErrCodeMissing
	}
	s.Tokens, err = rp.CodeExchange(ctx, s.Code, s.Provider.RelyingParty)
	return err
}

// maps the oidc.UserInfo to an [idp.User] using the default OIDC claims
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
