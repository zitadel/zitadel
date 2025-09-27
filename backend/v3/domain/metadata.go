package domain

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Metadata struct {
	InstanceID string `json:"instanceId,omitempty" db:"instance_id"`
	Key        string `json:"key,omitempty" db:"key"`
	// Value is a byte slice that might be json encoded.
	// the API does not require json encoding so we keep it as a byte slice here.
	Value []byte `json:"value,omitempty" db:"value"`

	CreatedAt time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt,omitzero" db:"updated_at"`
}

type MetadataColumns interface {
	// InstanceIDColumn returns the column for the instance id field.
	InstanceIDColumn() database.Column
	// KeyColumn returns the column for the key field.
	KeyColumn() database.Column
	// ValueColumn returns the column for the value field.
	ValueColumn() database.Column
	// CreatedAtColumn returns the column for the created at field.
	CreatedAtColumn() database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	UpdatedAtColumn() database.Column
}

type MetadataConditions interface {
	// InstanceIDCondition returns a filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
	// KeyCondition returns a filter on the key field.
	KeyCondition(op database.TextOperation, key string) database.Condition
	// ValueCondition returns a filter on the value field.
	ValueCondition(op database.BytesOperation, value []byte) database.Condition
}
