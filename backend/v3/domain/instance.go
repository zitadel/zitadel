package domain

import (
	"context"
	"time"

	"golang.org/x/text/language"

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

	Domains []*InstanceDomain `json:"domains,omitempty" db:"-"`
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
	// DefaultOrgIDColumn returns the column for the default org id field
	DefaultOrgIDColumn() database.Column
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
	// SetUpdatedAt sets the updated at column.
	SetUpdatedAt(time time.Time) database.Change
	// SetIAMProject sets the iam project column.
	SetIAMProject(id string) database.Change
	// SetDefaultOrg sets the default org column.
	SetDefaultOrg(id string) database.Change
	// SetDefaultLanguage sets the default language column.
	SetDefaultLanguage(language language.Tag) database.Change
	// SetConsoleClientID sets the console client id column.
	SetConsoleClientID(id string) database.Change
	// SetConsoleAppID sets the console app id column.
	SetConsoleAppID(id string) database.Change
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
	// If shouldLoad is set to true once, the Domains field will be set even if shouldLoad is false in the future.
	Domains(shouldLoad bool) InstanceDomainRepository
}

type CreateInstance struct {
	Name string `json:"name"`
}

