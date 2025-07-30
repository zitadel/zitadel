package domain

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type DomainValidationType string

const (
	DomainValidationTypeDNS  DomainValidationType = "dns"
	DomainValidationTypeHTTP DomainValidationType = "http"
)

type DomainType string

const (
	DomainTypeCustom  DomainType = "custom"
	DomainTypeTrusted DomainType = "trusted"
)

type domainColumns interface {
	// InstanceIDColumn returns the column for the instance id field.
	// `qualified` indicates if the column should be qualified with the table name.
	InstanceIDColumn(qualified bool) database.Column
	// DomainColumn returns the column for the domain field.
	// `qualified` indicates if the column should be qualified with the table name.
	DomainColumn(qualified bool) database.Column
	// IsPrimaryColumn returns the column for the is primary field.
	// `qualified` indicates if the column should be qualified with the table name.
	IsPrimaryColumn(qualified bool) database.Column
	// CreatedAtColumn returns the column for the created at field.
	// `qualified` indicates if the column should be qualified with the table name.
	CreatedAtColumn(qualified bool) database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	// `qualified` indicates if the column should be qualified with the table name.
	UpdatedAtColumn(qualified bool) database.Column
}

type domainConditions interface {
	// InstanceIDCondition returns a filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
	// DomainCondition returns a filter on the domain field.
	DomainCondition(op database.TextOperation, domain string) database.Condition
	// IsPrimaryCondition returns a filter on the is primary field.
	IsPrimaryCondition(isPrimary bool) database.Condition
}

type domainChanges interface {
	// SetPrimary sets a domain as primary based on the condition.
	// All other domains will be set to non-primary.
	//
	// An error is returned if:
	// - The condition identifies multiple domains.
	// - The condition does not identify any domain.
	//
	// This is a no-op if:
	// - The domain is already primary.
	// - No domain matches the condition.
	SetPrimary() database.Change
	// SetUpdatedAt sets the updated at column.
	// This is used for reducing events.
	SetUpdatedAt(t time.Time) database.Change
}
