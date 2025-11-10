package domain

import (
	"net"
	"net/http"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type SessionUserAgent struct {
	InstanceID    string      `json:"instanceId,omitempty" db:"instance_id"`
	FingerprintID *string     `json:"fingerprintId,omitempty" db:"fingerprint_id"`
	Description   *string     `json:"description,omitempty" db:"description"`
	IP            net.IP      `json:"ip,omitempty" db:"ip"`
	Header        http.Header `json:"header,omitempty" db:"headers"`
}

type sessionUserAgentColumns interface {
	// InstanceIDColumn returns the column for the instance id field.
	InstanceIDColumn() database.Column
	// FingerprintIDColumn returns the column for the fingerprint id field.
	FingerprintIDColumn() database.Column
	// IPColumn returns the column for the ip field.
	IPColumn() database.Column
	// DescriptionColumn returns the column for the description field.
	DescriptionColumn() database.Column
	// HeadersColumn returns the column for the headers field.
	HeadersColumn() database.Column
}

type sessionUserAgentConditions interface {
	// PrimaryKeyCondition returns a filter on the primary key fields.
	PrimaryKeyCondition(instanceID, fingerprintID string) database.Condition
	// InstanceIDCondition returns an equal filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
	// FingerprintIDCondition returns an equal filter on the fingerprint id field.
	FingerprintIDCondition(fingerprintID string) database.Condition
}
