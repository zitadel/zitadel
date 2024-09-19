package repository

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/user/model"
)

type RefreshTokenRepository interface {
	RefreshTokenByID(ctx context.Context, tokenID, userID string) (*model.RefreshTokenView, error)
	RefreshTokenByToken(ctx context.Context, refreshToken string) (*model.RefreshTokenView, error)
	SearchMyRefreshTokens(ctx context.Context, userID string, request *model.RefreshTokenSearchRequest) (*model.RefreshTokenSearchResponse, error)
}
