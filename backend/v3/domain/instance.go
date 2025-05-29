package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/cache"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Instance struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	DefaultOrgID    string     `json:"default_org_id"`
	IAMProjectID    string     `json:"iam_project_id"`
	ConsoleClientId string     `json:"console_client_id"`
	ConsoleAppID    string     `json:"console_app_id"`
	DefaultLanguage string     `json:"default_language"`
	CreatedAt       time.Time  `json:"-"`
	UpdatedAt       time.Time  `json:"-"`
	DeletedAt       *time.Time `json:"-"`
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

// instanceColumns define all the columns of the instance table.
type instanceColumns interface {
	// IDColumn returns the column for the id field.
	IDColumn() database.Column
	// NameColumn returns the column for the name field.
	NameColumn() database.Column
	// DefaultOrgIdColumn returns the column for the default org id field
	DefaultOrgIdColumn() database.Column
	// IAMProjectIDColumn returns the column for the default IAM org id field
	IAMProjectIDColumn() database.Column
	// ConsoleClientIDColumn returns the column for the default IAM org id field
	ConsoleClientIDColumn() database.Column
	// ConsoleAppIDColumn returns the column for the console client id field
	ConsoleAppIDColumn() database.Column
	// DefaultLanguageColumn returns the column for the default language field
	DefaultLanguageColumn() database.Column
	// CreatedAtColumn returns the column for the created at field.
	CreatedAtColumn() database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	UpdatedAtColumn() database.Column
	// DeletedAtColumn returns the column for the deleted at field.
	DeletedAtColumn() database.Column
}

// instanceConditions define all the conditions for the instance table.
type instanceConditions interface {
	// IDCondition returns an equal filter on the id field.
	IDCondition(instanceID string) database.Condition
	// NameCondition returns a filter on the name field.
	NameCondition(op database.TextOperation, name string) database.Condition
}

// instanceChanges define all the changes for the instance table.
type instanceChanges interface {
	// SetName sets the name column.
	SetName(name string) database.Change
}

// InstanceRepository is the interface for the instance repository.
type InstanceRepository interface {
	instanceColumns
	instanceConditions
	instanceChanges

	// TODO
	// Member returns the member repository which is a sub repository of the instance repository.
	// Member() MemberRepository

	Get(ctx context.Context, opts ...database.Condition) (*Instance, error)

	Create(ctx context.Context, instance *Instance) error
	Update(ctx context.Context, condition database.Condition, changes ...database.Change) error
	Delete(ctx context.Context, condition database.Condition) error
}

type CreateInstance struct {
	Name string `json:"name"`
}
