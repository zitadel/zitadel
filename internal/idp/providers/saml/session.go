package saml

import (
	"context"
	"net/http"

	"github.com/crewjam/saml/samlsp"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Session = (*Session)(nil)

// Session is the [idp.Session] implementation for the SAML provider.
type Session struct {
	Provider *Provider
	AuthURL  string

	RequestID string
	Request   *http.Request
}

// GetAuthURL implements the [idp.Session] interface.
func (s *Session) GetAuthURL() string {
	return s.AuthURL
}

// FetchUser implements the [idp.Session] interface.
func (s *Session) FetchUser(ctx context.Context) (user idp.User, err error) {
	sp, err := samlsp.New(*s.Provider.spOptions)
	if err != nil {
		return nil, err
	}

	assertion, err := sp.ServiceProvider.ParseResponse(s.Request, []string{s.RequestID})
	if err != nil {
		return nil, err
	}

	userMapper := &UserMapper{}
	userMapper.SetID(assertion.Subject.NameID)
	return userMapper, nil
}
