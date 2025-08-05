package oidc

import (
	"context"
	"errors"
	"time"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

var ErrCodeMissing = errors.New("no auth code provided")

var _ idp.Session = (*Session)(nil)

// Session is the [idp.Session] implementation for the OIDC provider.
type Session struct {
	Provider     *Provider
	AuthURL      string
	CodeVerifier string
	Code         string
	Tokens       *oidc.Tokens[*oidc.IDTokenClaims]
}

func NewSession(provider *Provider, code string, idpArguments map[string]any) *Session {
	verifier, _ := idpArguments[oauth.CodeVerifier].(string)
	return &Session{Provider: provider, Code: code, CodeVerifier: verifier}
}

// GetAuth implements the [idp.Session] interface.
func (s *Session) GetAuth(ctx context.Context) (idp.Auth, error) {
	return idp.Redirect(s.AuthURL)
}

// PersistentParameters implements the [idp.Session] interface.
func (s *Session) PersistentParameters() map[string]any {
	if s.CodeVerifier == "" {
		return nil
	}
	return map[string]any{oauth.CodeVerifier: s.CodeVerifier}
}

// FetchUser implements the [idp.Session] interface.
// It will execute an OIDC code exchange if needed to retrieve the tokens,
// call the userinfo endpoint and map the received information into an [idp.User].
func (s *Session) FetchUser(ctx context.Context) (user idp.User, err error) {
	if s.Tokens == nil {
		if err = s.Authorize(ctx); err != nil {
			return nil, err
		}
	}

	var info *oidc.UserInfo
	if s.Provider.useIDToken {
		info = s.Tokens.IDTokenClaims.GetUserInfo()
	} else {
		info, err = rp.Userinfo[*oidc.UserInfo](ctx,
			s.Tokens.AccessToken,
			s.Tokens.TokenType,
			s.Tokens.IDTokenClaims.GetSubject(),
			s.Provider.RelyingParty,
		)
		if err != nil {
			return nil, err
		}
	}
	u := s.Provider.userInfoMapper(info)
	return u, nil
}

func (s *Session) ExpiresAt() time.Time {
	if s.Tokens == nil {
		return time.Time{}
	}
	return s.Tokens.Expiry
}

func (s *Session) Authorize(ctx context.Context) (err error) {
	if s.Code == "" {
		return ErrCodeMissing
	}
	var opts []rp.CodeExchangeOpt
	if s.CodeVerifier != "" {
		opts = append(opts, rp.WithCodeVerifier(s.CodeVerifier))
	}
	s.Tokens, err = rp.CodeExchange[*oidc.IDTokenClaims](ctx, s.Code, s.Provider.RelyingParty, opts...)
	return err
}

func NewUser(info *oidc.UserInfo) *User {
	return &User{UserInfo: info}
}

func InitUser() *User {
	return &User{UserInfo: &oidc.UserInfo{}}
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
