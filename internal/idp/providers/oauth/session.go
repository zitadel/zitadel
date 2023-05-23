package oauth

import (
	"context"
	"errors"
	"net/http"

	"github.com/zitadel/oidc/v2/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v2/pkg/http"
	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
)

var ErrCodeMissing = errors.New("no auth code provided")

var _ idp.Session = (*Session)(nil)

// Session is the [idp.Session] implementation for the OAuth2.0 provider.
type Session struct {
	AuthURL string
	Code    string
	Tokens  *oidc.Tokens[*oidc.IDTokenClaims]

	Provider *Provider
}

// GetAuthURL implements the [idp.Session] interface.
func (s *Session) GetAuthURL() string {
	return s.AuthURL
}

// FetchUser implements the [idp.Session] interface.
// It will execute an OAuth 2.0 code exchange if needed to retrieve the access token,
// call the specified userEndpoint and map the received information into an [idp.User].
func (s *Session) FetchUser(ctx context.Context) (user idp.User, err error) {
	if s.Tokens == nil {
		if err = s.authorize(ctx); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest("GET", s.Provider.userEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("authorization", s.Tokens.TokenType+" "+s.Tokens.AccessToken)
	mapper := s.Provider.userMapper()
	if err := httphelper.HttpRequest(s.Provider.RelyingParty.HttpClient(), req, &mapper); err != nil {
		return nil, err
	}
	return mapper, nil
}

func (s *Session) authorize(ctx context.Context) (err error) {
	if s.Code == "" {
		return ErrCodeMissing
	}
	s.Tokens, err = rp.CodeExchange[*oidc.IDTokenClaims](ctx, s.Code, s.Provider.RelyingParty)
	if err != nil {
		return err
	}
	return nil
}
