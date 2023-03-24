package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Session = (*Session)(nil)

var (
	ErrNoTokens     = errors.New("no tokens provided")
	ErrInvalidToken = errors.New("invalid tokens provided")
)

// Session is the [idp.Session] implementation for the JWT provider
type Session struct {
	*Provider
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
	s.Tokens.IDTokenClaims, err = s.validateToken(ctx, s.Tokens.IDToken)
	if err != nil {
		return nil, err
	}
	return &User{s.Tokens.IDTokenClaims}, nil
}

func (s *Session) validateToken(ctx context.Context, token string) (oidc.IDTokenClaims, error) {
	logging.Debug("begin token validation")
	// TODO: be able to specify them in the template: https://github.com/zitadel/zitadel/issues/5322
	offset := 3 * time.Second
	maxAge := time.Hour
	claims := oidc.EmptyIDTokenClaims()
	payload, err := oidc.ParseToken(token, claims)
	if err != nil {
		return nil, fmt.Errorf("%w: malformed jwt payload: %v", ErrInvalidToken, err)
	}

	if err = oidc.CheckIssuer(claims, s.Provider.issuer); err != nil {
		return nil, fmt.Errorf("%w: invalid issuer: %v", ErrInvalidToken, err)
	}

	logging.Debug("begin signature validation")
	keySet := rp.NewRemoteKeySet(http.DefaultClient, s.Provider.keysEndpoint)
	if err = oidc.CheckSignature(ctx, token, payload, claims, nil, keySet); err != nil {
		return nil, fmt.Errorf("%w: invalid signature: %v", ErrInvalidToken, err)
	}

	if !claims.GetExpiration().IsZero() {
		if err = oidc.CheckExpiration(claims, offset); err != nil {
			return nil, fmt.Errorf("%w: expired: %v", ErrInvalidToken, err)
		}
	}

	if !claims.GetIssuedAt().IsZero() {
		if err = oidc.CheckIssuedAt(claims, maxAge, offset); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
		}
	}
	return claims, nil
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

func (u *User) GetPhone() domain.PhoneNumber {
	return domain.PhoneNumber(u.IDTokenClaims.GetPhoneNumber())
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

func (u *User) GetEmail() domain.EmailAddress {
	return domain.EmailAddress(u.IDTokenClaims.GetEmail())
}
