package repository

import (
	"context"
	usr_model "github.com/caos/zitadel/internal/user/model"
)

type TokenRepository interface {
	IsTokenValid(ctx context.Context, userID, tokenID string) (bool, error)
	TokenByID(ctx context.Context, userID, tokenID string) (*usr_model.TokenView, error)
}
