package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type UserIdentityProviderLink struct {
	InstanceID         string `json:"instanceID" db:"instance_id"`
	UserID             string `json:"userID" db:"user_id"`
	IdentityProviderID string `json:"identityProviderID" db:"identity_provider_id"`
	ProvidedID         string `json:"providedID" db:"provided_id"`
	ProvidedUsername   string `json:"providedUsername" db:"provided_username"`

	CreatedAt time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

//go:generate mockgen -typed -package domainmock -destination ./mock/user_identity_provider_link.mock.go . UserIdentityProviderLinkRepository
type UserIdentityProviderLinkRepository interface {
	Repository
	userIdentityProviderLinkColumns
	userIdentityProviderLinkConditions
	userIdentityProviderLinkChanges

	Get(ctx context.Context, client database.QueryExecutor, condition database.Condition, opts ...database.QueryOption) (*UserIdentityProviderLink, error)
	List(ctx context.Context, client database.QueryExecutor, condition database.Condition, opts ...database.QueryOption) ([]*UserIdentityProviderLink, error)

	Create(ctx context.Context, client database.QueryExecutor, link *UserIdentityProviderLink) error
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}

type userIdentityProviderLinkColumns interface {
	InstanceIDColumn() database.Column
	UserIDColumn() database.Column
	IdentityProviderIDColumn() database.Column
	ProvidedIDColumn() database.Column
	ProvidedUsernameColumn() database.Column

	CreatedAtColumn() database.Column
	UpdatedAtColumn() database.Column
}

type userIdentityProviderLinkConditions interface {
	PrimaryKeyCondition(instanceID, idpID, providedID string) database.Condition
	InstanceIDCondition(op database.TextOperation, id string) database.Condition
	UserIDCondition(op database.TextOperation, id string) database.Condition
	IdentityProviderIDCondition(op database.TextOperation, id string) database.Condition
	ProvidedIDCondition(op database.TextOperation, id string) database.Condition
	ProvidedUsernameCondition(op database.TextOperation, username string) database.Condition
}

type userIdentityProviderLinkChanges interface {
	SetProvidedUsername(username string) database.Change
	SetUpdatedAt(updatedAt time.Time) database.Change
	SetProvidedUserID(providedID string) database.Change
}
