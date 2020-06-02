package repository

import (
	"context"
)

type TokenVerifierRepository interface {
	VerifyAccessToken(ctx context.Context, appName string) (string, string, string, error)
}
