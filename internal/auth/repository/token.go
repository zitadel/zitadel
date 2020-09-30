package repository

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/token/model"
)

type TokenRepository interface {
	CreateToken(ctx context.Context, agentID, applicationID, userID string, audience, scopes []string, lifetime time.Duration) (*model.Token, error)
	ValidTokenByID(ctx context.Context, tokenID string) (*model.Token, error)
}
