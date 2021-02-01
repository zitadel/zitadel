package repository

import (
	"context"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"time"
)

type TokenRepository interface {
	CreateToken(ctx context.Context, agentID, clientID, subject string, audience, scopes []string, lifetime time.Duration) (string, time.Time, error)
	IsTokenValid(ctx context.Context, userID, tokenID string) (bool, error)
	TokenByID(ctx context.Context, userID, tokenID string) (*usr_model.TokenView, error)
}
