package authz

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/errors"
)

const (
	BearerPrefix = "Bearer "
)

func verifyAccessToken(ctx context.Context, token string, t TokenVerifier) (string, string, string, error) {
	parts := strings.Split(token, BearerPrefix)
	if len(parts) != 2 {
		return "", "", "", errors.ThrowUnauthenticated(nil, "AUTH-7fs1e", "invalid auth header")
	}
	return t.VerifyAccessToken(ctx, parts[1])
}
