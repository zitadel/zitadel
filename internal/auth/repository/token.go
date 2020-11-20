package repository

import (
	"context"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"time"
)

type TokenRepository interface {
	CreateToken(ctx context.Context, agentID, clientID, userID string, audience, scopes []string, lifetime time.Duration) (*usr_model.Token, error)
	IsTokenValid(ctx context.Context, userID, tokenID string) (bool, error)
	TokenByID(ctx context.Context, userID, tokenID string) (*usr_model.TokenView, error)
}
