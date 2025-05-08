package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/cache"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Instance struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

type instanceCacheIndex uint8

const (
	instanceCacheIndexUndefined instanceCacheIndex = iota
	instanceCacheIndexID
)

// Keys implements the [cache.Entry].
func (i *Instance) Keys(index instanceCacheIndex) (key []string) {
	if index == instanceCacheIndexID {
		return []string{i.ID}
	}
	return nil
}

var _ cache.Entry[instanceCacheIndex, string] = (*Instance)(nil)

type instanceColumns interface {
	// IDColumn returns the column for the id field.
	IDColumn() database.Column
	// NameColumn returns the column for the name field.
	NameColumn() database.Column
	// CreatedAtColumn returns the column for the created at field.
	CreatedAtColumn() database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	UpdatedAtColumn() database.Column
	// DeletedAtColumn returns the column for the deleted at field.
	DeletedAtColumn() database.Column
}

type instanceConditions interface {
	// IDCondition returns an equal filter on the id field.
	IDCondition(instanceID string) database.Condition
	// NameCondition returns a filter on the name field.
	NameCondition(op database.TextOperation, name string) database.Condition
}

type instanceChanges interface {
	// SetName sets the name column.
	SetName(name string) database.Change
}

type InstanceRepository interface {
	instanceColumns
	instanceConditions
	instanceChanges

	Member() MemberRepository

	Get(ctx context.Context, opts ...database.QueryOption) (*Instance, error)

	Create(ctx context.Context, instance *Instance) error
	Update(ctx context.Context, condition database.Condition, changes ...database.Change) error
	Delete(ctx context.Context, condition database.Condition) error
}

type CreateInstance struct {
	Name string `json:"name"`
}
