package jwt

import (
	"context"

	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Session = (*Session)(nil)

// Session is the [idp.Session] implementation for the JWT provider
type Session struct {
	AuthURL string
	Tokens  *oidc.Tokens
}

// GetAuthURL implements the [idp.Session] interface
func (s *Session) GetAuthURL() string {
	return s.AuthURL
}

// FetchUser implements the [idp.Session] interface.
// It will map the received idToken into an [idp.User].
func (s *Session) FetchUser(ctx context.Context) (user idp.User, err error) {
	if s.Tokens == nil {
		return nil, ErrNoTokens
	}
	return &User{s.Tokens.IDTokenClaims}, nil
}

type User struct {
	oidc.IDTokenClaims
}

func (u *User) GetID() string {
	return u.IDTokenClaims.GetSubject()
}

func (u *User) GetFirstName() string {
	return u.IDTokenClaims.GetGivenName()
}

func (u *User) GetLastName() string {
	return u.IDTokenClaims.GetFamilyName()
}

func (u *User) GetDisplayName() string {
	return u.IDTokenClaims.GetName()
}

func (u *User) GetNickname() string {
	return u.IDTokenClaims.GetNickname()
}

func (u *User) GetPhone() string {
	return u.IDTokenClaims.GetPhoneNumber()
}

func (u *User) IsPhoneVerified() bool {
	return u.IDTokenClaims.IsPhoneNumberVerified()
}

func (u *User) GetPreferredLanguage() language.Tag {
	return u.IDTokenClaims.GetLocale()
}

func (u *User) GetAvatarURL() string {
	return u.IDTokenClaims.GetPicture()
}
