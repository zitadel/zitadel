package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type RefreshToken struct {
	TokenID string `json:"tokenId,omitempty"`
}

type refreshTokenConditions interface {
	PrimaryKeyCondition(instanceID, tokenID string) database.Condition
	UserIDCondition(userID string) database.Condition
	InstanceIDCondition(instanceID string) database.Condition
}

type refreshTokenChanges interface{}

type RefreshTokenRepository interface {
	refreshTokenChanges
	refreshTokenConditions

	Get(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*RefreshToken, error)
	List(ctx context.Context, client database.QueryExecutor, condition database.Condition) ([]*RefreshToken, error)

	Create(ctx context.Context, client database.QueryExecutor, refreshToken *RefreshToken) error
	Update(ctx context.Context, client database.QueryExecutor, refreshToken *RefreshToken, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}
