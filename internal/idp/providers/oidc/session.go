package oidc

import (
	"context"
	"errors"

	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
)

var ErrCodeMissing = errors.New("no auth code provided")

var _ idp.Session = (*Session)(nil)

// Session is the [idp.Session] implementation for the OIDC provider.
type Session struct {
	Provider *Provider
	AuthURL  string
	Code     string
	Tokens   *oidc.Tokens[*oidc.IDTokenClaims]
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
			return nil, err
		}
	}
	info, err := rp.Userinfo(
		s.Tokens.AccessToken,
		s.Tokens.TokenType,
		s.Tokens.IDTokenClaims.GetSubject(),
		s.Provider.RelyingParty,
	)
	if err != nil {
		return nil, err
	}
	if s.Provider.useIDToken {
		info = s.Tokens.IDTokenClaims.GetUserInfo()
	}
	u := s.Provider.userInfoMapper(info)
	return u, nil
}

func (s *Session) authorize(ctx context.Context) (err error) {
	if s.Code == "" {
		return ErrCodeMissing
	}
	s.Tokens, err = rp.CodeExchange[*oidc.IDTokenClaims](ctx, s.Code, s.Provider.RelyingParty)
	return err
}

func NewUser(info *oidc.UserInfo) *User {
	return &User{UserInfo: info}
}

type User struct {
	*oidc.UserInfo
}

func (u *User) GetID() string {
	return u.Subject
}

func (u *User) GetFirstName() string {
	return u.GivenName
}

func (u *User) GetLastName() string {
	return u.FamilyName
}

func (u *User) GetDisplayName() string {
	return u.Name
}

func (u *User) GetNickname() string {
	return u.Nickname
}

func (u *User) GetPreferredUsername() string {
	return u.PreferredUsername
}

func (u *User) GetEmail() domain.EmailAddress {
	return domain.EmailAddress(u.UserInfo.Email)
}

func (u *User) IsEmailVerified() bool {
	return bool(u.EmailVerified)
}

func (u *User) GetPhone() domain.PhoneNumber {
	return domain.PhoneNumber(u.PhoneNumber)
}

func (u *User) IsPhoneVerified() bool {
	return u.PhoneNumberVerified
}

func (u *User) GetPreferredLanguage() language.Tag {
	return u.Locale.Tag()
}

func (u *User) GetAvatarURL() string {
	return u.Picture
}

func (u *User) GetProfile() string {
	return u.Profile
}
