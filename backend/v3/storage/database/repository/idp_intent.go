package repository

import (
	"context"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

type idpIntent struct{}

func (i idpIntent) qualifiedTableName() string {
	return "zitadel.identity_provider_intents"
}

func (i idpIntent) unqualifiedTableName() string {
	return "identity_provider_intents"
}

func IDPIntentRepository() domain.IDPIntentRepository {
	return new(idpIntent)
}

// Create implements [domain.IDPIntentRepository].
func (i *idpIntent) Create(ctx context.Context, client database.QueryExecutor, intent *domain.IDPIntent) error {
	panic("unimplemented")
}

// Delete implements [domain.IDPIntentRepository].
func (i *idpIntent) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	panic("unimplemented")
}

// Get implements [domain.IDPIntentRepository].
func (i *idpIntent) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPIntent, error) {
	panic("unimplemented")
}

// Update implements [domain.IDPIntentRepository].
func (i *idpIntent) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	panic("unimplemented")
}

// LoadIdentityProvider implements [domain.IDPIntentRepository].
func (i *idpIntent) LoadIdentityProvider() domain.IDPIntentRepository {
	panic("unimplemented")
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetAssertion implements [domain.IDPIntentRepository].
func (i *idpIntent) SetAssertion(assertion string) database.Change {
	panic("unimplemented")
}

// SetFailureURL implements [domain.IDPIntentRepository].
func (i *idpIntent) SetFailureURL(failureURL url.URL) database.Change {
	panic("unimplemented")
}

// SetIDPAccessToken implements [domain.IDPIntentRepository].
func (i *idpIntent) SetIDPAccessToken(idpAccessToken string) database.Change {
	panic("unimplemented")
}

// SetIDPArguments implements [domain.IDPIntentRepository].
func (i *idpIntent) SetIDPArguments(idpArguments map[string]any) database.Change {
	panic("unimplemented")
}

// SetIDPEntryAttributes implements [domain.IDPIntentRepository].
func (i *idpIntent) SetIDPEntryAttributes(idpEntryAttributes map[string][]string) database.Change {
	panic("unimplemented")
}

// SetIDPID implements [domain.IDPIntentRepository].
func (i *idpIntent) SetIDPID(idpID string) database.Change {
	panic("unimplemented")
}

// SetIDPUser implements [domain.IDPIntentRepository].
func (i *idpIntent) SetIDPUser(idpUser []byte) database.Change {
	panic("unimplemented")
}

// SetIDPUserID implements [domain.IDPIntentRepository].
func (i *idpIntent) SetIDPUserID(idpUserID string) database.Change {
	panic("unimplemented")
}

// SetIDPUsername implements [domain.IDPIntentRepository].
func (i *idpIntent) SetIDPUsername(idpUsername string) database.Change {
	panic("unimplemented")
}

// SetRequestID implements [domain.IDPIntentRepository].
func (i *idpIntent) SetRequestID(requestID string) database.Change {
	panic("unimplemented")
}

// SetState implements [domain.IDPIntentRepository].
func (i *idpIntent) SetState(state domain.IDPIntentState) database.Change {
	panic("unimplemented")
}

// SetSuccessURL implements [domain.IDPIntentRepository].
func (i *idpIntent) SetSuccessURL(successURL url.URL) database.Change {
	panic("unimplemented")
}

// SetUserID implements [domain.IDPIntentRepository].
func (i *idpIntent) SetUserID(userID string) database.Change {
	panic("unimplemented")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// CreatedAtCondition implements [domain.IDPIntentRepository].
func (i *idpIntent) CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition {
	panic("unimplemented")
}

// IDCondition implements [domain.IDPIntentRepository].
func (i *idpIntent) IDCondition(id string) database.Condition {
	panic("unimplemented")
}

// IDPIDCondition implements [domain.IDPIntentRepository].
func (i *idpIntent) IDPIDCondition(idpID string) database.Condition {
	panic("unimplemented")
}

// IDPUserIDCondition implements [domain.IDPIntentRepository].
func (i *idpIntent) IDPUserIDCondition(idpUserID string) database.Condition {
	panic("unimplemented")
}

// IDPUsernameCondition implements [domain.IDPIntentRepository].
func (i *idpIntent) IDPUsernameCondition(idpUsername string) database.Condition {
	panic("unimplemented")
}

// InstanceIDCondition implements [domain.IDPIntentRepository].
func (i *idpIntent) InstanceIDCondition(instanceID string) database.Condition {
	panic("unimplemented")
}

// PrimaryKeyCondition implements [domain.IDPIntentRepository].
func (i *idpIntent) PrimaryKeyCondition(instanceID string, id string) database.Condition {
	panic("unimplemented")
}

// RequestIDCondition implements [domain.IDPIntentRepository].
func (i *idpIntent) RequestIDCondition(requestID string) database.Condition {
	panic("unimplemented")
}

// StateCondition implements [domain.IDPIntentRepository].
func (i *idpIntent) StateCondition(state domain.IDPIntentState) database.Condition {
	panic("unimplemented")
}

// UpdatedAtCondition implements [domain.IDPIntentRepository].
func (i *idpIntent) UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition {
	panic("unimplemented")
}

// UserIDCondition implements [domain.IDPIntentRepository].
func (i *idpIntent) UserIDCondition(userID string) database.Condition {
	panic("unimplemented")
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// AssertionColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) AssertionColumn() database.Column {
	panic("unimplemented")
}

// CreatedAtColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) CreatedAtColumn() database.Column {
	panic("unimplemented")
}

// ExpiresAtColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) ExpiresAtColumn() database.Column {
	panic("unimplemented")
}

// FailureURLColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) FailureURLColumn() database.Column {
	panic("unimplemented")
}

// IDColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) IDColumn() database.Column {
	panic("unimplemented")
}

// IDPAccessTokenColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) IDPAccessTokenColumn() database.Column {
	panic("unimplemented")
}

// IDPArgumentsColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) IDPArgumentsColumn() database.Column {
	panic("unimplemented")
}

// IDPEntryAttributesColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) IDPEntryAttributesColumn() database.Column {
	panic("unimplemented")
}

// IDPIDColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) IDPIDColumn() database.Column {
	panic("unimplemented")
}

// IDPUserColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) IDPUserColumn() database.Column {
	panic("unimplemented")
}

// IDPUserIDColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) IDPUserIDColumn() database.Column {
	panic("unimplemented")
}

// IDPUsernameColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) IDPUsernameColumn() database.Column {
	panic("unimplemented")
}

// InstanceIDColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) InstanceIDColumn() database.Column {
	panic("unimplemented")
}

// MaxIDPIntentLifetimeColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) MaxIDPIntentLifetimeColumn() database.Column {
	panic("unimplemented")
}

// PrimaryKeyColumns implements [domain.IDPIntentRepository].
func (i *idpIntent) PrimaryKeyColumns() []database.Column {
	panic("unimplemented")
}

// RequestIDColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) RequestIDColumn() database.Column {
	panic("unimplemented")
}

// StateColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) StateColumn() database.Column {
	panic("unimplemented")
}

// SucceededAtColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) SucceededAtColumn() database.Column {
	panic("unimplemented")
}

// SuccessURLColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) SuccessURLColumn() database.Column {
	panic("unimplemented")
}

// UpdatedAtColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) UpdatedAtColumn() database.Column {
	panic("unimplemented")
}

// UserIDColumn implements [domain.IDPIntentRepository].
func (i *idpIntent) UserIDColumn() database.Column {
	panic("unimplemented")
}

var _ domain.IDPIntentRepository = (*idpIntent)(nil)
