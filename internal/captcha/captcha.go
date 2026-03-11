package captcha

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

// CaptchaVerifier verifies captcha challenge tokens from the client.
type CaptchaVerifier interface {
	// Verify checks a captcha response token. Returns true if the challenge
	// was solved successfully.
	Verify(ctx context.Context, token string, remoteIP string) (bool, error)
	// SiteKey returns the public site key for the client-side widget.
	SiteKey() string
	// Provider returns the provider name (e.g. "turnstile", "hcaptcha", "recaptcha").
	Provider() string
}

// CaptchaConfig holds captcha provider configuration.
type CaptchaConfig struct {
	Enabled   bool          `yaml:"enabled"`
	Provider  string        `yaml:"provider"` // turnstile, hcaptcha, recaptcha
	SiteKey   string        `yaml:"siteKey"`
	SecretKey string        `yaml:"secretKey"`
	VerifyURL string        `yaml:"verifyURL"`
	Timeout   time.Duration `yaml:"timeout"`
}

func (c CaptchaConfig) IsEnabled() bool {
	return c.Enabled && c.SecretKey != ""
}

// NewCaptchaVerifier creates a verifier from config.
func NewCaptchaVerifier(cfg CaptchaConfig, httpClient *http.Client) CaptchaVerifier {
	if !cfg.IsEnabled() {
		return nil
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: cfg.timeout()}
	}
	switch strings.ToLower(cfg.Provider) {
	case "turnstile":
		return &TurnstileVerifier{
			siteKey:   cfg.SiteKey,
			secretKey: cfg.SecretKey,
			verifyURL: cfg.verifyURL(),
			client:    httpClient,
		}
	case "hcaptcha":
		return &TurnstileVerifier{
			siteKey:   cfg.SiteKey,
			secretKey: cfg.SecretKey,
			verifyURL: defaultIfEmpty(cfg.VerifyURL, "https://hcaptcha.com/siteverify"),
			client:    httpClient,
			provider:  "hcaptcha",
		}
	case "recaptcha":
		return &TurnstileVerifier{
			siteKey:   cfg.SiteKey,
			secretKey: cfg.SecretKey,
			verifyURL: defaultIfEmpty(cfg.VerifyURL, "https://www.google.com/recaptcha/api/siteverify"),
			client:    httpClient,
			provider:  "recaptcha",
		}
	default:
		return &TurnstileVerifier{
			siteKey:   cfg.SiteKey,
			secretKey: cfg.SecretKey,
			verifyURL: cfg.verifyURL(),
			client:    httpClient,
		}
	}
}

func (c CaptchaConfig) timeout() time.Duration {
	if c.Timeout > 0 {
		return c.Timeout
	}
	return 5 * time.Second
}

func (c CaptchaConfig) verifyURL() string {
	return defaultIfEmpty(c.VerifyURL, "https://challenges.cloudflare.com/turnstile/v0/siteverify")
}

func defaultIfEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// TurnstileVerifier implements CaptchaVerifier for Cloudflare Turnstile.
// It also works for hCaptcha and reCAPTCHA since they share the same
// POST form + JSON response protocol.
type TurnstileVerifier struct {
	siteKey   string
	secretKey string
	verifyURL string
	client    *http.Client
	provider  string // override for Provider() — defaults to "turnstile"
}

type turnstileResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

func (v *TurnstileVerifier) SiteKey() string { return v.siteKey }

func (v *TurnstileVerifier) Provider() string {
	if v.provider != "" {
		return v.provider
	}
	return "turnstile"
}

func (v *TurnstileVerifier) Verify(ctx context.Context, token string, remoteIP string) (bool, error) {
	if token == "" {
		return false, nil
	}

	form := url.Values{
		"secret":   {v.secretKey},
		"response": {token},
	}
	if remoteIP != "" {
		form.Set("remoteip", remoteIP)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.verifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return false, fmt.Errorf("captcha request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("captcha verify: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return false, fmt.Errorf("captcha read response: %w", err)
	}

	var result turnstileResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("captcha parse response: %w", err)
	}

	if !result.Success {
		logging.Info(ctx, "risk.captcha.verification_failed",
			slog.Any("captcha_error_codes", result.ErrorCodes),
			slog.String("captcha_provider", v.Provider()),
		)
	}

	return result.Success, nil
}
