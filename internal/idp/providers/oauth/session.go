package oauth

import (
	"context"
	"errors"
	"net/http"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
)

var ErrCodeMissing = errors.New("no auth code provided")

const (
	CodeVerifier = "codeVerifier"
)

var _ idp.Session = (*Session)(nil)

// Session is the [idp.Session] implementation for the OAuth2.0 provider.
type Session struct {
	AuthURL      string
	CodeVerifier string
	Code         string
	Tokens       *oidc.Tokens[*oidc.IDTokenClaims]

	Provider *Provider
}

func NewSession(provider *Provider, code string, idpArguments map[string]any) *Session {
	verifier, _ := idpArguments[CodeVerifier].(string)
	return &Session{Provider: provider, Code: code, CodeVerifier: verifier}
}

// GetAuth implements the [idp.Session] interface.
func (s *Session) GetAuth(ctx context.Context) (string, bool) {
	return idp.Redirect(s.AuthURL)
}

// PersistentParameters implements the [idp.Session] interface.
func (s *Session) PersistentParameters() map[string]any {
	if s.CodeVerifier == "" {
		return nil
	}
	return map[string]any{CodeVerifier: s.CodeVerifier}
}

// FetchUser implements the [idp.Session] interface.
// It will execute an OAuth 2.0 code exchange if needed to retrieve the access token,
// call the specified userEndpoint and map the received information into an [idp.User].
func (s *Session) FetchUser(ctx context.Context) (_ idp.User, err error) {
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
	user := s.Provider.User()
	if err := httphelper.HttpRequest(s.Provider.RelyingParty.HttpClient(), req, &user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Session) authorize(ctx context.Context) (err error) {
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
