package repository

import (
	"context"
)

type TokenVerifierRepository interface {
	VerifyAccessToken(ctx context.Context, appName string) (string, string, string, error)
	ProjectIDByClientID(ctx context.Context, clientID string) (string, error)
}
