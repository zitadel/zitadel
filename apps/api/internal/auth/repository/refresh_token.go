package repository

import (
	"context"

	"github.com/zitadel/zitadel/internal/user/model"
)

type RefreshTokenRepository interface {
	RefreshTokenByID(ctx context.Context, tokenID, userID string) (*model.RefreshTokenView, error)
	RefreshTokenByToken(ctx context.Context, refreshToken string) (*model.RefreshTokenView, error)
	SearchMyRefreshTokens(ctx context.Context, userID string, request *model.RefreshTokenSearchRequest) (*model.RefreshTokenSearchResponse, error)
}
