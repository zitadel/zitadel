package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/zitadel/zitadel/apps/cli/internal/config"
	"golang.org/x/oauth2"
)

// TokenSource returns the correct oauth2.TokenSource for the given context.
// Resolution order:
//  1. ZITADEL_TOKEN env var → PAT
//  2. context.AuthMethod == "pat" → PAT from config
//  3. context.Token (cached interactive token)
func TokenSource(_ context.Context, cfg *config.Context) (oauth2.TokenSource, error) {
	// 1. Env var override
	if token := os.Getenv("ZITADEL_TOKEN"); token != "" {
		return PATTokenSource(token), nil
	}

	// 2. PAT from config
	if cfg.AuthMethod == "pat" && cfg.PAT != "" {
		return PATTokenSource(cfg.PAT), nil
	}

	// 3. Cached interactive token
	if cfg.Token != "" {
		return PATTokenSource(cfg.Token), nil
	}

	return nil, fmt.Errorf("no token available; run 'zitadel login' or set ZITADEL_TOKEN")
}
