package captcha

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// RecaptchaV2Config stores site and secret keys
type RecaptchaV2Config struct {
	SiteKey   string
	SecretKey string
}

// RecaptchaV2 handles reCAPTCHA v2 validation
type RecaptchaV2 struct {
	config RecaptchaV2Config
	client *http.Client
}

// Initialize reCAPTCHA v2 with site key and secret key
func (r *RecaptchaV2) Initialize(config any) error {
	cfg, ok := config.(RecaptchaV2Config)
	if !ok {
		return fmt.Errorf("invalid config type for RecaptchaV2: %T", config)
	}
	if cfg.SiteKey == "" || cfg.SecretKey == "" {
		return errors.New("site key and secret key are required")
	}
	if r.client == nil {
		r.client = &http.Client{Timeout: 5 * time.Second}
	}
	r.config = cfg
	return nil
}

// GetToken retrieves the token from the request
func (r *RecaptchaV2) GetToken(req *http.Request) (string, error) {
	token := req.FormValue("g-recaptcha-response")
	if token == "" {
		return "", errors.New("missing captcha token")
	}

	return token, nil
}

// ValidateToken verifies the token using Google's API
func (r *RecaptchaV2) ValidateToken(token string) (bool, error) {
	if token == "" {
		return false, errors.New("token cannot be empty")
	}

	// Google reCAPTCHA verification URL
	apiURL := "https://www.google.com/recaptcha/api/siteverify"

	// Prepare form data
	formData := url.Values{}
	formData.Set("secret", r.config.SecretKey) // Secret key for backend validation
	formData.Set("response", token)            // Token received from client

	// Make HTTP request
	resp, err := r.client.PostForm(apiURL, formData)
	if err != nil {
		return false, fmt.Errorf("failed to contact reCAPTCHA server: %w", err)
	}
	defer resp.Body.Close()

	// Parse JSON response
	var result struct {
		Success     bool     `json:"success"`
		ChallengeTs string   `json:"challenge_ts"`
		Hostname    string   `json:"hostname"`
		ErrorCodes  []string `json:"error-codes,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to parse reCAPTCHA response: %w", err)
	}

	// Check if verification was successful
	if !result.Success {
		return false, fmt.Errorf("reCAPTCHA verification failed: %v", result.ErrorCodes)
	}

	return true, nil
}
