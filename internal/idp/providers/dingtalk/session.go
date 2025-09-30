package dingtalk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
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
	// Step 1: Get basic user info using personal access token
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

	// Step 2: Try to get corporate email and mobile using app access token
	// This is optional and may fail if the app doesn't have the required permissions
	if corpEmail, corpMobile, err := s.getCorpEmailAndMobile(ctx, user.UnionID); err == nil {
		// Use corporate email if available and valid
		if corpEmail != "" {
			user.Email = domain.EmailAddress(corpEmail)
		}
		// Use corporate mobile if available and valid
		if corpMobile != "" {
			user.Mobile = corpMobile
		}
	}

	return &user, nil
}

// getCorpEmailAndMobile tries to get corporate email and mobile using app access token
func (s *Session) getCorpEmailAndMobile(ctx context.Context, unionID string) (string, string, error) {
	// Step 1: Get app access token
	appAccessToken, err := s.getAppAccessToken(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to get app access token: %w", err)
	}

	// Step 2: Get user ID using union ID
	userID, err := s.getUserIDByUnionID(ctx, unionID, appAccessToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user ID: %w", err)
	}

	// Step 3: Get corporate user details
	corpEmail, corpMobile, err := s.getUserCorpDetails(ctx, userID, appAccessToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to get corp details: %w", err)
	}

	return corpEmail, corpMobile, nil
}

// getAppAccessToken gets application access token for internal API calls
func (s *Session) getAppAccessToken(ctx context.Context) (string, error) {
	tokenRequest := map[string]string{
		"appKey":    s.clientID,
		"appSecret": s.clientSecret,
	}

	requestBody, err := json.Marshal(tokenRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal app token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.dingtalk.com/v1.0/oauth2/accessToken", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create app token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute app token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("app token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse struct {
		AccessToken string `json:"accessToken"`
		ExpireIn    int    `json:"expireIn"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("failed to decode app token response: %w", err)
	}

	return tokenResponse.AccessToken, nil
}

// getUserIDByUnionID gets internal user ID using union ID
func (s *Session) getUserIDByUnionID(ctx context.Context, unionID, appAccessToken string) (string, error) {
	requestBody := map[string]string{
		"unionid": unionID,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user ID request: %w", err)
	}

	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/user/getbyunionid?access_token=%s", appAccessToken)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create user ID request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute user ID request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("user ID request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		ErrCode    int    `json:"errcode"`
		ErrMessage string `json:"errmsg"`
		Result     struct {
			UserID string `json:"userid"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode user ID response: %w", err)
	}

	if response.ErrCode == 60121 {
		return "", fmt.Errorf("该应用只允许本企业内部用户登录，您不属于该企业，无法登录")
	} else if response.ErrCode != 0 {
		return "", fmt.Errorf("DingTalk API error: %s", response.ErrMessage)
	}

	return response.Result.UserID, nil
}

// getUserCorpDetails gets corporate email and mobile using user ID
func (s *Session) getUserCorpDetails(ctx context.Context, userID, appAccessToken string) (string, string, error) {
	requestBody := map[string]string{
		"userid": userID,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal corp details request: %w", err)
	}

	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/v2/user/get?access_token=%s", appAccessToken)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", "", fmt.Errorf("failed to create corp details request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to execute corp details request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("corp details request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		ErrCode    int    `json:"errcode"`
		ErrMessage string `json:"errmsg"`
		Result     struct {
			Mobile    string `json:"mobile"`
			Email     string `json:"email"`
			JobNumber string `json:"job_number"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", "", fmt.Errorf("failed to decode corp details response: %w", err)
	}

	if response.ErrCode != 0 {
		return "", "", fmt.Errorf("DingTalk API error: %s", response.ErrMessage)
	}

	return response.Result.Email, response.Result.Mobile, nil
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
