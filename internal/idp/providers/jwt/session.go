package jwt

import (
	"context"

	"github.com/zitadel/oidc/v2/pkg/oidc"

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
		return idp.User{}, ErrNoTokens
	}
	mapTokenToUser(s.Tokens.IDTokenClaims, &user)
	return user, nil
}

func mapTokenToUser(claims oidc.IDTokenClaims, user *idp.User) {
	user.ID = claims.GetSubject()
	user.AvatarURL = claims.GetPicture()
	user.DisplayName = claims.GetName()
	user.Email = claims.GetEmail()
	user.FirstName = claims.GetGivenName()
	user.LastName = claims.GetFamilyName()
	user.NickName = claims.GetNickname()
}
