package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/zitadel/zitadel/apps/cli/internal/config"
	"golang.org/x/oauth2"
)

// TokenSource returns the correct oauth2.TokenSource for the given context.
// Resolution order:
//  1. ZITADEL_TOKEN env var → PAT
//  2. context.AuthMethod == "pat" → PAT from config
//  3. context.Token (cached interactive token) with refresh support
func TokenSource(_ context.Context, cfg *config.Context) (oauth2.TokenSource, error) {
	// 1. Env var override
	if token := os.Getenv("ZITADEL_TOKEN"); token != "" {
		return PATTokenSource(token), nil
	}

	// 2. PAT from config
	if cfg.AuthMethod == "pat" && cfg.PAT != "" {
		return PATTokenSource(cfg.PAT), nil
	}

	// 3. Cached interactive token with auto-refresh
	if cfg.Token != "" {
		return newInteractiveTokenSource(cfg), nil
	}

	return nil, fmt.Errorf("no token available; run 'zitadel login' or set ZITADEL_TOKEN")
}

// interactiveTokenSource wraps a cached access token with refresh logic.
type interactiveTokenSource struct {
	cfg *config.Context
}

func newInteractiveTokenSource(cfg *config.Context) *interactiveTokenSource {
	return &interactiveTokenSource{cfg: cfg}
}

func (s *interactiveTokenSource) Token() (*oauth2.Token, error) {
	tok := &oauth2.Token{
		AccessToken:  s.cfg.Token,
		RefreshToken: s.cfg.RefreshToken,
		TokenType:    "Bearer",
	}

	if s.cfg.TokenExpiry != "" {
		if expiry, err := time.Parse(time.RFC3339, s.cfg.TokenExpiry); err == nil {
			tok.Expiry = expiry
		}
	}

	// If the token has no expiry or is still valid, return it as-is.
	if tok.Expiry.IsZero() || time.Until(tok.Expiry) > 30*time.Second {
		return tok, nil
	}

	// Token expired or about to expire — try refresh.
	if s.cfg.RefreshToken == "" {
		return nil, fmt.Errorf("access token expired and no refresh token available; run 'zitadel login' again")
	}

	refreshed, err := s.refresh()
	if err != nil {
		return nil, fmt.Errorf("token refresh failed (run 'zitadel login' again): %w", err)
	}
	return refreshed, nil
}

func (s *interactiveTokenSource) refresh() (*oauth2.Token, error) {
	issuer := s.cfg.Instance
	if !strings.HasPrefix(issuer, "http://") && !strings.HasPrefix(issuer, "https://") {
		issuer = "https://" + issuer
	}

	oauthCfg := &oauth2.Config{
		ClientID: s.cfg.ClientID,
		Endpoint: oauth2.Endpoint{
			TokenURL: issuer + "/oauth/v2/token",
		},
	}

	tok := &oauth2.Token{RefreshToken: s.cfg.RefreshToken}
	newTok, err := oauthCfg.TokenSource(context.Background(), tok).Token()
	if err != nil {
		return nil, err
	}

	// Persist the refreshed tokens back to config.
	s.cfg.Token = newTok.AccessToken
	if newTok.RefreshToken != "" {
		s.cfg.RefreshToken = newTok.RefreshToken
	}
	if !newTok.Expiry.IsZero() {
		s.cfg.TokenExpiry = newTok.Expiry.Format(time.RFC3339)
	}

	if err := s.persistConfig(); err != nil {
		log.Printf("warning: could not persist refreshed token: %v", err)
	}

	return newTok, nil
}

func (s *interactiveTokenSource) persistConfig() error {
	fullCfg, err := config.Load()
	if err != nil {
		return err
	}
	name := fullCfg.ActiveContext
	if name == "" {
		return nil
	}
	ctx, ok := fullCfg.Contexts[name]
	if !ok {
		return nil
	}
	ctx.Token = s.cfg.Token
	ctx.RefreshToken = s.cfg.RefreshToken
	ctx.TokenExpiry = s.cfg.TokenExpiry
	fullCfg.Contexts[name] = ctx
	return config.Save(fullCfg)
}
