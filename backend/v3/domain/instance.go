package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/cache"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Instance struct {
	ID              string    `json:"id,omitempty" db:"id"`
	Name            string    `json:"name,omitempty" db:"name"`
	DefaultOrgID    string    `json:"defaultOrgId,omitempty" db:"default_org_id"`
	IAMProjectID    string    `json:"iamProjectId,omitempty" db:"iam_project_id"`
	ConsoleClientID string    `json:"consoleClientId,omitempty" db:"console_client_id"`
	ConsoleAppID    string    `json:"consoleAppId,omitempty" db:"console_app_id"`
	DefaultLanguage string    `json:"defaultLanguage,omitempty" db:"default_language"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`

	Domains []*InstanceDomain `json:"domains,omitempty" db:"domains"`
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
	// `qualified` indicates if the column should be qualified with the table name.
	IDColumn(qualified bool) database.Column
	// NameColumn returns the column for the name field.
	// `qualified` indicates if the column should be qualified with the table name.
	NameColumn(qualified bool) database.Column
	// DefaultOrgIDColumn returns the column for the default org id field
	// `qualified` indicates if the column should be qualified with the table name.
	DefaultOrgIDColumn(qualified bool) database.Column
	// IAMProjectIDColumn returns the column for the default IAM org id field
	// `qualified` indicates if the column should be qualified with the table name.
	IAMProjectIDColumn(qualified bool) database.Column
	// ConsoleClientIDColumn returns the column for the default IAM org id field
	// `qualified` indicates if the column should be qualified with the table name.
	ConsoleClientIDColumn(qualified bool) database.Column
	// ConsoleAppIDColumn returns the column for the console client id field
	// `qualified` indicates if the column should be qualified with the table name.
	ConsoleAppIDColumn(qualified bool) database.Column
	// DefaultLanguageColumn returns the column for the default language field
	// `qualified` indicates if the column should be qualified with the table name.
	DefaultLanguageColumn(qualified bool) database.Column
	// CreatedAtColumn returns the column for the created at field.
	// `qualified` indicates if the column should be qualified with the table name.
	CreatedAtColumn(qualified bool) database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	// `qualified` indicates if the column should be qualified with the table name.
	UpdatedAtColumn(qualified bool) database.Column
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

	Get(ctx context.Context, opts ...database.QueryOption) (*Instance, error)
	List(ctx context.Context, opts ...database.QueryOption) ([]*Instance, error)

	Create(ctx context.Context, instance *Instance) error
	Update(ctx context.Context, id string, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)

	// Domains returns the domain sub repository for the instance.
	// If shouldLoad is true, the domains will be loaded from the database and written to the [Instance].Domains field.
	// If shouldLoad is set to true once, the Domains field will be set event if shouldLoad is false in the future.
	Domains(shouldLoad bool) InstanceDomainRepository
}

type CreateInstance struct {
	Name string `json:"name"`
}
