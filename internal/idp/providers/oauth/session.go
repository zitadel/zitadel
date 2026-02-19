package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
	"golang.org/x/text/language"
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
func (s *Session) GetAuth(ctx context.Context) (idp.Auth, error) {
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
	// Check if the provider is Facebook
	if s.isFacebookProvider() {
		user, err := s.fetchFacebookUser(ctx)
		return user, err
	}

	if s.Tokens == nil {
		if err = s.authorize(ctx); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest("GET", s.Provider.userEndpoint, nil)
	if err != nil {
		return nil, err
	}

	authHeader := s.Tokens.TokenType + " " + s.Tokens.AccessToken
	req.Header.Set("authorization", authHeader)

	user := s.Provider.User()
	if err := httphelper.HttpRequest(s.Provider.RelyingParty.HttpClient(), req, &user); err != nil {
		return nil, err
	}

	return user, nil
}

// isFacebookProvider checks if the current provider is Facebook
func (s *Session) isFacebookProvider() bool {
	// Check by provider name
	if strings.Contains(strings.ToLower(s.Provider.Name()), "facebook") {
		return true
	}

	// Check by user endpoint
	if strings.Contains(s.Provider.userEndpoint, "graph.facebook.com") {
		return true
	}

	// Check by issuer if available
	if s.Provider.RelyingParty != nil {
		issuer := s.Provider.RelyingParty.Issuer()
		if strings.Contains(issuer, "facebook.com") {
			return true
		}
	}

	return false
}

// fetchFacebookUser handles Facebook-specific user fetching
func (s *Session) fetchFacebookUser(ctx context.Context) (user idp.User, err error) {
	if s.Tokens == nil {
		if err = s.authorize(ctx); err != nil {
			return nil, err
		}
	}

	// Use the access token to get user info from Facebook's /me endpoint
	endpoint := "https://graph.facebook.com/me?fields=id,name,email,first_name,last_name,picture&access_token=" + s.Tokens.AccessToken

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Picture   struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	// Create a user object that implements the idp.User interface
	u := &FacebookUser{
		ID:        userInfo.ID,
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
		Picture:   userInfo.Picture.Data.URL,
	}

	return u, nil
}

// FacebookUser implements the idp.User interface for Facebook users
type FacebookUser struct {
	ID        string
	Email     string
	Name      string
	FirstName string
	LastName  string
	Picture   string
}

func (u *FacebookUser) GetID() string {
	return u.ID
}

func (u *FacebookUser) GetFirstName() string {
	return u.FirstName
}

func (u *FacebookUser) GetLastName() string {
	return u.LastName
}

func (u *FacebookUser) GetDisplayName() string {
	return u.Name
}

func (u *FacebookUser) GetNickname() string {
	return ""
}

func (u *FacebookUser) GetPreferredUsername() string {
	return ""
}

func (u *FacebookUser) GetEmail() domain.EmailAddress {
	return domain.EmailAddress(u.Email)
}

func (u *FacebookUser) IsEmailVerified() bool {
	return u.Email != ""
}

func (u *FacebookUser) GetPhone() domain.PhoneNumber {
	return ""
}

func (u *FacebookUser) IsPhoneVerified() bool {
	return false
}

func (u *FacebookUser) GetPreferredLanguage() language.Tag {
	return language.Und
}

func (u *FacebookUser) GetAvatarURL() string {
	return u.Picture
}

func (u *FacebookUser) GetProfile() string {
	return ""
}

func (s *Session) ExpiresAt() time.Time {
	if s.Tokens == nil {
		return time.Time{}
	}
	return s.Tokens.Expiry
}

func (s *Session) authorize(ctx context.Context) (err error) {
	if s.Code == "" {
		return ErrCodeMissing
	}

	var opts []rp.CodeExchangeOpt
	if s.CodeVerifier != "" {
		opts = append(opts, rp.WithCodeVerifier(s.CodeVerifier))
	}

	// Check if the provider is Facebook and set custom token URL
	if s.isFacebookProvider() {
		s.Provider.RelyingParty.OAuthConfig().Endpoint.TokenURL = "https://graph.facebook.com/v23.0/oauth/access_token"
		s.Provider.RelyingParty.OAuthConfig().Scopes = []string{"email", "public_profile"}
	}

	s.Tokens, err = rp.CodeExchange[*oidc.IDTokenClaims](ctx, s.Code, s.Provider.RelyingParty, opts...)
	return err
}
