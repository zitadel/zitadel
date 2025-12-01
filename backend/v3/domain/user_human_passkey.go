package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type PasskeyType uint8

const (
	PasskeyTypePasswordless PasskeyType = iota
	PasskeyTypeU2F
)

type PasskeyState uint8

const (
	PasskeyStateUnverified PasskeyState = iota
	PasskeyStateVerified
)

type Passkey struct {
	InstanceID string `json:"instanceId,omitempty" db:"instance_id"`
	UserID     string `json:"userId,omitempty" db:"user_id"`
	TokenID    string `json:"tokenId,omitempty" db:"token_id"`

	CreatedAt time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" db:"updated_at"`

	Type           PasskeyType `json:"type,omitempty" db:"type"`
	Name           string      `json:"name,omitempty" db:"name"`
	SignCount      uint32      `json:"signCount,omitempty" db:"sign_count"`
	Challenge      []byte      `json:"challenge,omitempty" db:"challenge"`
	RelyingPartyID string      `json:"relyingPartyId,omitempty" db:"relying_party_id"`
}

type passkeyConditions interface {
	PrimaryKeyCondition(instanceID, tokenID string) database.Condition
	UserIDCondition(userID string) database.Condition
	InstanceIDCondition(instanceID string) database.Condition
}

type passkeyChanges interface {
	// SetSignCount sets the sign count column.
	SetSignCount(signCount uint32) database.Change
	// SetUpdatedAt sets the updated at column.
	SetUpdatedAt(updatedAt time.Time) database.Change
	// SetIsVerified sets the is verified column.
	SetState(state PasskeyState) database.Change
}

//go:generate mockgen -typed -package domainmock -destination ./mock/user_human_passkey.mock.go . PasskeyRepository
type PasskeyRepository interface {
	passkeyChanges
	passkeyConditions

	Get(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*Passkey, error)
	List(ctx context.Context, client database.QueryExecutor, condition database.Condition, limit, offset int) ([]*Passkey, error)

	Add(ctx context.Context, client database.QueryExecutor, passkey *Passkey) error
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)

	SetInitializationVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification VerificationType) (int64, error)
	GetInitializationVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*Verification, error)
	IncrementFailedInitializationAttempts(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}
