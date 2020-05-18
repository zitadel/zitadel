package repository

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/token/model"
)

type TokenRepository interface {
	CreateToken(ctx context.Context, agentID, applicationID, userID string, lifetime time.Duration) (*model.Token, error)
	IsTokenValid(ctx context.Context, tokenID string) (bool, error)
}
