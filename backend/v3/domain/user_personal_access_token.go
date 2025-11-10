package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type PersonalAccessToken struct {
	InstanceID string `json:"instanceId,omitempty" db:"instance_id"`
	UserID     string `json:"userId,omitempty" db:"user_id"`
	ID         string `json:"tokenId,omitempty" db:"token_id"`

	CreatedAt time.Time `json:"createdAt,omitzero" db:"created_at"`
	ExpiresAt time.Time `json:"expiresAt,omitzero" db:"expires_at"`
	Scopes    []string  `json:"scopes,omitempty" db:"scopes"`
}

type personalAccessTokenConditions interface {
	PrimaryKeyCondition(instanceID, tokenID string) database.Condition
	UserIDCondition(userID string) database.Condition
	InstanceIDCondition(instanceID string) database.Condition
}

type PersonalAccessTokenRepository interface {
	personalAccessTokenConditions

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOpts) (*PersonalAccessToken, error)
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOpts) ([]*PersonalAccessToken, error)
	Create(ctx context.Context, client database.QueryExecutor, pat *PersonalAccessToken) error
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}
