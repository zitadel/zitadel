package repository

import (
	"context"
	usr_model "github.com/caos/zitadel/internal/user/model"
)

type RefreshTokenRepository interface {
	//IsTokenValid(ctx context.Context, userID, tokenID string) (bool, error)
	RefreshTokenByID(ctx context.Context, refreshToken string) (*usr_model.RefreshTokenView, error)
}
