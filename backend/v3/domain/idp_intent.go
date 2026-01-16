package domain

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate enumer -type  IDPIntentState -transform lower -trimprefix IDPIntentState -sql
type IDPIntentState uint8

const (
	IDPIntentStateStarted IDPIntentState = iota
	IDPIntentStateSucceeded
	IDPIntentStateFailed
)

type IDPIntent struct {
	ID              string             `json:"id,omitempty" db:"id"`
	InstanceID      string             `json:"instance_id,omitempty" db:"instance_id"`
	State           IDPIntentState     `json:"state" db:"state"`
	SuccessURL      *url.URL           `json:"success_url,omitempty" db:"success_url"`
	FailureURL      *url.URL           `json:"failure_url,omitempty" db:"failure_url"`
	CreatedAt       time.Time          `json:"created_at,omitzero" db:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at,omitzero" db:"updated_at"`
	IDPID           string             `json:"idp_id,omitempty" db:"idp_id"`
	IDPArguments    map[string]any     `json:"idp_arguments,omitempty" db:"idp_arguments"`
	IDPUser         []byte             `json:"idp_user,omitempty" db:"idp_user"`
	IDPUserID       string             `json:"idp_user_id,omitempty" db:"idp_user_id"`
	IDPUsername     string             `json:"idp_username,omitempty" db:"idp_username"`
	UserID          string             `json:"user_id,omitempty" db:"user_id"`
	IDPAccessToken  []byte             `json:"idp_access_token,omitempty" db:"idp_access_token"`
	IDPIDToken      string             `json:"idp_id_token,omitempty" db:"idp_id_token"`
	EntryAttributes IDPEntryAttributes `json:"idp_entry_attributes,omitempty" db:"idp_entry_attributes"`
	RequestID       string             `json:"request_id,omitempty" db:"request_id"`
	Assertion       []byte             `json:"assertion,omitempty" db:"assertion"`
	SucceededAt     *time.Time         `json:"succeeded_at,omitzero" db:"succeeded_at"`
	FailReason      string             `json:"fail_reason,omitempty" db:"fail_reason"`
	ExpiresAt       *time.Time         `json:"expires_at,omitzero" db:"expires_at"`
}

type IDPEntryAttributes map[string][]string

func NewIDPEntryAttributes(attrs map[string][]string) *IDPEntryAttributes {
	cast := IDPEntryAttributes(attrs)
	return &cast
}

func (iea *IDPEntryAttributes) String() ([]byte, error) {
	return json.Marshal(iea)
}

// idpIntentColumns defines all the columns of the identity_provider_intents table.
type idpIntentColumns interface {
	// IDColumn returns the column for the field id
	IDColumn() database.Column
	// InstanceIDColumn returns the column for the field instance_id
	InstanceIDColumn() database.Column
	// StateColumn returns the column for the field state
	StateColumn() database.Column
	// SuccessURLColumn returns the column for the field success_url
	SuccessURLColumn() database.Column
	// FailureURLColumn returns the column for the field failure_url
	FailureURLColumn() database.Column
	// CreatedAtColumn returns the column for the field created_at
	CreatedAtColumn() database.Column
	// UpdatedAtColumn returns the column for the field updated_at
	UpdatedAtColumn() database.Column
	// IDPIDColumn returns the column for the field idp_id
	IDPIDColumn() database.Column
	// IDPArgumentsColumn returns the column for the field idp_arguments
	IDPArgumentsColumn() database.Column
	// IDPUserColumn returns the column for the field idp_user
	IDPUserColumn() database.Column
	// IDPUserIDColumn returns the column for the field idp_user_id
	IDPUserIDColumn() database.Column
	// IDPUsernameColumn returns the column for the field idp_username
	IDPUsernameColumn() database.Column
	// UserIDColumn returns the column for the field user_id
	UserIDColumn() database.Column
	// IDPAccessTokenColumn returns the column for the field idp_access_token
	IDPAccessTokenColumn() database.Column
	// IDPIDTokenColumn returns the column for the field idp_id_token
	IDPIDTokenColumn() database.Column
	// IDPEntryAttributesColumn returns the column for the field idp_entry_attributes
	IDPEntryAttributesColumn() database.Column
	// RequestIDColumn returns the column for the field request_id
	RequestIDColumn() database.Column
	// AssertionColumn returns the column for the field assertion
	AssertionColumn() database.Column
	// SucceededAtColumn returns the column for the field succeeded_at
	SucceededAtColumn() database.Column
	// FailReasonColumn returns the column for the field fail_reason
	FailReasonColumn() database.Column
	// ExpiresAtColumn returns the column for the field expires_at
	ExpiresAtColumn() database.Column
}

type idpIntentConditions interface {
	// PrimaryKeyCondition returns a filter on the primary key fields.
	PrimaryKeyCondition(instanceID, id string) database.Condition
	// IDCondition returns an equal filter on the ID
	IDCondition(id string) database.Condition
	// InstanceIDCondition returns an equal filter on the InstanceID
	InstanceIDCondition(instanceID string) database.Condition
	// StateCondition returns an equal filter on the State
	StateCondition(state IDPIntentState) database.Condition
	// IDPIDCondition returns an equal filter on the IDPID
	IDPIDCondition(idpID string) database.Condition
	// IDPUserIDCondition returns an equal filter on the IDPUserID
	IDPUserIDCondition(idpUserID string) database.Condition
	// IDPUsernameCondition returns an equal filter on the IDPUsername
	IDPUsernameCondition(idpUsername string) database.Condition
	// UserIDCondition returns an equal filter on the UserID
	UserIDCondition(userID string) database.Condition
	// RequestIDCondition returns an equal filter on the RequestID
	RequestIDCondition(requestID string) database.Condition
	// CreatedAtCondition returns a filter on the created at field.
	CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition
	// UpdatedAtCondition returns a filter on the updated at field.
	UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition
}

type idpIntentChanges interface {
	// SetState sets the state of the IDP intent
	SetState(state IDPIntentState) database.Change
	// SetSuccessURL sets the successURL of the IDP intent
	SetSuccessURL(successURL url.URL) database.Change
	// SetFailureURL sets the failureURL of the IDP intent
	SetFailureURL(failureURL url.URL) database.Change
	// SetIDPID sets the idpID of the IDP intent
	SetIDPID(idpID string) database.Change
	// SetIDPArguments sets the idpArguments of the IDP intent
	SetIDPArguments(idpArguments []byte) database.Change
	// SetIDPUser sets the idpUser of the IDP intent
	SetIDPUser(idpUser []byte) database.Change
	// SetIDPUserID sets the idpUserID of the IDP intent
	SetIDPUserID(idpUserID string) database.Change
	// SetIDPUsername sets the idpUsername of the IDP intent
	SetIDPUsername(idpUsername string) database.Change
	// SetUserID sets the userID of the IDP intent
	SetUserID(userID string) database.Change
	// SetIDPAccessToken sets the idpAccessToken of the IDP intent
	SetIDPAccessToken(idpAccessToken []byte) database.Change
	// SetIDPIDToken sets the idpIDToken of the IDP intent
	SetIDPIDToken(idpIDToken string) database.Change
	// SetIDPEntryAttributes sets the idpEntryAttributes of the IDP intent
	SetIDPEntryAttributes(idpEntryAttributes []byte) database.Change
	// SetRequestID sets the requestID of the IDP intent
	SetRequestID(requestID string) database.Change
	// SetAssertion sets the assertion of the IDP intent
	SetAssertion(assertion []byte) database.Change
	// SetFailReason sets the fail reason of the IDP intent
	SetFailReason(reason string) database.Change
	// SetExpiresAt sets the expiration of the IDP intent
	SetExpiresAt(expiration time.Time) database.Change
	// SetSucceededAt sets the succededAt time of the IDP intent
	SetSucceededAt(succeededAt time.Time) database.Change
}

//go:generate mockgen -typed -package domainmock -destination ./mock/idp_intent.mock.go . IDPIntentRepository
type IDPIntentRepository interface {
	Repository

	idpIntentColumns
	idpIntentConditions
	idpIntentChanges

	// Get returns an idp intent based on the given condition.
	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*IDPIntent, error)
	// Create creates a new idp intent.
	Create(ctx context.Context, client database.QueryExecutor, intent *IDPIntent) error
	// Update one or more existing idp intents.
	// The condition must include at least the instanceID of the idp intent to update.
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	// Delete removes idp intents based on the given condition.
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}
