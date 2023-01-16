package jwt

import (
	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Session = (*Session)(nil)

// Session is the idp.Session implementation for the JWT provider
type Session struct {
	AuthURL string
	Tokens  *oidc.Tokens
}

func (s *Session) GetAuthURL() string {
	return s.AuthURL
}

func (s *Session) FetchUser() (user idp.User, err error) {
	if s.Tokens == nil {
		return idp.User{}, ErrNoTokens
	}
	err = mapTokenToUser(s.Tokens.IDTokenClaims, &user)
	return user, err
}

func mapTokenToUser(claims oidc.IDTokenClaims, user *idp.User) error {
	user.ID = claims.GetSubject()
	user.AvatarURL = claims.GetPicture()
	user.DisplayName = claims.GetName()
	user.Email = claims.GetEmail()
	user.FirstName = claims.GetGivenName()
	user.LastName = claims.GetFamilyName()
	user.NickName = claims.GetNickname()
	return nil
}
