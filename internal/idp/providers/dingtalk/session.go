package dingtalk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Session = (*Session)(nil)

// Session implements a DingTalk-specific session that handles DingTalk's non-standard OAuth flow
type Session struct {
	AuthURL      string
	provider     *Provider
	code         string
	clientID     string
	clientSecret string
	accessToken  string
	expiresAt    time.Time
}

// NewSession creates a new DingTalk session
func NewSession(provider *Provider, code string, clientID, clientSecret string) *Session {
	return &Session{
		provider:     provider,
		code:         code,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// FetchUser implements the idp.Session interface and handles DingTalk's custom OAuth flow
func (s *Session) FetchUser(ctx context.Context) (idp.User, error) {
	// Step 1: Get user access token using DingTalk's custom flow
	accessToken, err := s.getUserAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user access token: %w", err)
	}

	// Step 2: Get user info using the access token
	user, err := s.getUserInfo(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return user, nil
}

// getUserAccessToken exchanges the authorization code for a user access token
// This follows DingTalk's custom OAuth flow instead of standard OAuth
func (s *Session) getUserAccessToken(ctx context.Context) (string, error) {
	// Return cached token if still valid
	if s.accessToken != "" && time.Now().Before(s.expiresAt) {
		return s.accessToken, nil
	}

	// Prepare request body for DingTalk's token endpoint
	tokenRequest := map[string]interface{}{
		"clientId":     s.clientID,
		"clientSecret": s.clientSecret,
		"code":         s.code,
		"grantType":    "authorization_code",
	}

	requestBody, err := json.Marshal(tokenRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token request: %w", err)
	}

	// Make request to DingTalk's token endpoint
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResponse struct {
		AccessToken string `json:"accessToken"`
		ExpireIn    int    `json:"expireIn"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResponse.AccessToken == "" {
		return "", fmt.Errorf("no access token in response")
	}

	// Cache the token and set expiration time
	s.accessToken = tokenResponse.AccessToken
	if tokenResponse.ExpireIn > 0 {
		s.expiresAt = time.Now().Add(time.Duration(tokenResponse.ExpireIn) * time.Second)
	} else {
		// Default expiration of 1 hour if not provided
		s.expiresAt = time.Now().Add(time.Hour)
	}

	return s.accessToken, nil
}

// getUserInfo fetches user information using DingTalk's API
func (s *Session) getUserInfo(ctx context.Context, accessToken string) (*User, error) {
	// Make request to DingTalk's user info endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", profileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}

	// DingTalk uses a custom header for access token
	req.Header.Set("x-acs-dingtalk-access-token", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute user info request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user info request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse user info response
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user info response: %w", err)
	}

	return &user, nil
}

// GetAuth implements the [idp.Session] interface.
func (s *Session) GetAuth(ctx context.Context) (idp.Auth, error) {
	return &idp.RedirectAuth{RedirectURL: s.AuthURL}, nil
}

// PersistentParameters implements the [idp.Session] interface.
func (s *Session) PersistentParameters() map[string]any {
	return nil // DingTalk doesn't need persistent parameters like PKCE
}

// ExpiresAt implements the [idp.Session] interface.
func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}
