package github

import (
	"context"
	"net/http"
	"strings"
	"time"

	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

var _ idp.Session = (*Session)(nil)

// Session extends the [oauth.Session] to be able to handle private email addresses.
type Session struct {
	*Provider
	Code         string
	IDPArguments map[string]any

	OAuthSession *oauth.Session
}

func NewSession(provider *Provider, code string, idpArguments map[string]any) *Session {
	return &Session{Provider: provider, Code: code, IDPArguments: idpArguments}
}

// GetAuth implements the [idp.Provider] interface by calling the wrapped [oauth.Session].
func (s *Session) GetAuth(ctx context.Context) (idp.Auth, error) {
	return s.oauth().GetAuth(ctx)
}

// PersistentParameters implements the [idp.Session] interface by calling the wrapped [oauth.Session].
func (s *Session) PersistentParameters() map[string]any {
	return s.oauth().PersistentParameters()
}

// FetchUser implements the [idp.Session] interface.
// It will execute an OAuth 2.0 code exchange if needed to retrieve the access token,
// call the specified userEndpoint and map the received information into an [idp.User].
// In case of a specific TenantID as [TenantType] it will additionally extract the id_token and validate it.
func (s *Session) FetchUser(ctx context.Context) (user idp.User, err error) {
	user, err = s.oauth().FetchUser(ctx)
	if err != nil {
		return nil, err
	}
	if user.GetEmail() == "" {
		userInternal := user.(*User)
		for _, scope := range s.Provider.Provider.OAuthConfig().Scopes {
			if strings.TrimSpace(scope) == "user:email" {
				email, err := s.getPrivateMail()
				if err != nil {
					return user, nil
				}
				userInternal.Email = domain.EmailAddress(email)
				return userInternal, nil
			}
		}
	}
	return user, nil
}

func (s *Session) ExpiresAt() time.Time {
	if s.OAuthSession == nil {
		return time.Time{}
	}
	return s.OAuthSession.ExpiresAt()
}

// Tokens returns the [oidc.Tokens] of the underlying [oauth.Session].
func (s *Session) Tokens() *oidc.Tokens[*oidc.IDTokenClaims] {
	return s.oauth().Tokens
}

func (s *Session) oauth() *oauth.Session {
	if s.OAuthSession != nil {
		return s.OAuthSession
	}
	s.OAuthSession = oauth.NewSession(s.Provider.Provider, s.Code, s.IDPArguments)
	return s.OAuthSession
}

type Email struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func (s *Session) getPrivateMail() (email string, err error) {
	req, err := http.NewRequest("GET", s.Provider.emailURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("authorization", s.oauth().Tokens.TokenType+" "+s.oauth().Tokens.AccessToken)

	emailList := make([]Email, 0)
	if err := httphelper.HttpRequest(s.Provider.HttpClient(), req, &emailList); err != nil {
		return "", err
	}
	for _, v := range emailList {
		if v.Primary && v.Verified {
			return v.Email, nil
		}
	}
	return email, nil
}
