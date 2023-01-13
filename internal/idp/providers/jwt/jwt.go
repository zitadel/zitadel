package jwt

import (
	"errors"
	"net/url"

	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
)

const queryAuthRequestID = "authRequestID"

var _ idp.Provider = (*Provider)(nil)

var ErrNoTokens = errors.New("no tokens")

type Provider struct {
	name         string
	issuer       string
	jwtEndpoint  string
	keysEndpoint string
	headerName   string
}

func New(issuer, jwtEndpoint, keysEndpoint, headerName string) (*Provider, error) {
	provider := &Provider{
		issuer:       issuer,
		jwtEndpoint:  jwtEndpoint,
		keysEndpoint: keysEndpoint,
		headerName:   headerName,
	}

	return provider, nil
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) BeginAuth(state string) (idp.Session, error) {
	redirect, err := url.Parse(p.jwtEndpoint)
	if err != nil {
		return nil, err
	}
	q := redirect.Query()
	q.Set(queryAuthRequestID, state)
	//TODO: userAgentID
	redirect.RawQuery = q.Encode()
	return &Session{AuthURL: redirect.String()}, nil
}

func (p *Provider) FetchUser(session idp.Session) (user idp.User, err error) {
	jwtSession := session.(*Session)
	if jwtSession.Tokens == nil {
		return idp.User{}, ErrNoTokens
	}
	err = mapTokenToUser(jwtSession.Tokens.IDTokenClaims, &user)
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
