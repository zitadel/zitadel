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

// SetAssertion implements [domain.idpIntentChanges].
func (i *idpIntent) SetAssertion(assertion string) database.Change {
	panic("unimplemented")
}

// SetFailureURL implements [domain.idpIntentChanges].
func (i *idpIntent) SetFailureURL(failureURL url.URL) database.Change {
	panic("unimplemented")
}

// SetIDPAccessToken implements [domain.idpIntentChanges].
func (i *idpIntent) SetIDPAccessToken(idpAccessToken string) database.Change {
	panic("unimplemented")
}

// SetIDPArguments implements [domain.idpIntentChanges].
func (i *idpIntent) SetIDPArguments(idpArguments map[string]any) database.Change {
	panic("unimplemented")
}

// SetIDPEntryAttributes implements [domain.idpIntentChanges].
func (i *idpIntent) SetIDPEntryAttributes(idpEntryAttributes map[string][]string) database.Change {
	panic("unimplemented")
}

// SetIDPID implements [domain.idpIntentChanges].
func (i *idpIntent) SetIDPID(idpID string) database.Change {
	panic("unimplemented")
}

// SetIDPUser implements [domain.idpIntentChanges].
func (i *idpIntent) SetIDPUser(idpUser []byte) database.Change {
	panic("unimplemented")
}

// SetIDPUserID implements [domain.idpIntentChanges].
func (i *idpIntent) SetIDPUserID(idpUserID string) database.Change {
	panic("unimplemented")
}

// SetIDPUsername implements [domain.idpIntentChanges].
func (i *idpIntent) SetIDPUsername(idpUsername string) database.Change {
	panic("unimplemented")
}

// SetRequestID implements [domain.idpIntentChanges].
func (i *idpIntent) SetRequestID(requestID string) database.Change {
	panic("unimplemented")
}

// SetState implements [domain.idpIntentChanges].
func (i *idpIntent) SetState(state domain.IDPIntentState) database.Change {
	panic("unimplemented")
}

// SetSuccessURL implements [domain.idpIntentChanges].
func (i *idpIntent) SetSuccessURL(successURL url.URL) database.Change {
	panic("unimplemented")
}

// SetUserID implements [domain.idpIntentChanges].
func (i *idpIntent) SetUserID(userID string) database.Change {
	panic("unimplemented")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// CreatedAtCondition implements [domain.idpIntentConditions].
func (i *idpIntent) CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition {
	return database.NewNumberCondition(i.CreatedAtColumn(), op, createdAt)
}

// IDCondition implements [domain.idpIntentConditions].
func (i *idpIntent) IDCondition(id string) database.Condition {
	return database.NewTextCondition(i.IDColumn(), database.TextOperationEqual, id)
}

// IDPIDCondition implements [domain.idpIntentConditions].
func (i *idpIntent) IDPIDCondition(idpID string) database.Condition {
	return database.NewTextCondition(i.IDPIDColumn(), database.TextOperationEqual, idpID)
}

// IDPUserIDCondition implements [domain.idpIntentConditions].
func (i *idpIntent) IDPUserIDCondition(idpUserID string) database.Condition {
	return database.NewTextCondition(i.IDPUserIDColumn(), database.TextOperationEqual, idpUserID)
}

// IDPUsernameCondition implements [domain.idpIntentConditions].
func (i *idpIntent) IDPUsernameCondition(idpUsername string) database.Condition {
	return database.NewTextCondition(i.IDPUsernameColumn(), database.TextOperationEqual, idpUsername)
}

// InstanceIDCondition implements [domain.idpIntentConditions].
func (i *idpIntent) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(i.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// PrimaryKeyCondition implements [domain.idpIntentConditions].
func (i *idpIntent) PrimaryKeyCondition(instanceID, id string) database.Condition {
	return database.And(
		i.InstanceIDCondition(instanceID),
		i.IDCondition(id),
	)
}

// RequestIDCondition implements [domain.idpIntentConditions].
func (i *idpIntent) RequestIDCondition(requestID string) database.Condition {
	return database.NewTextCondition(i.RequestIDColumn(), database.TextOperationEqual, requestID)
}

// StateCondition implements [domain.idpIntentConditions].
func (i *idpIntent) StateCondition(state domain.IDPIntentState) database.Condition {
	return database.NewTextCondition(i.StateColumn(), database.TextOperationEqual, state.String())
}

// UpdatedAtCondition implements [domain.idpIntentConditions].
func (i *idpIntent) UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition {
	return database.NewNumberCondition(i.UpdatedAtColumn(), op, updatedAt)
}

// UserIDCondition implements [domain.idpIntentConditions].
func (i *idpIntent) UserIDCondition(userID string) database.Condition {
	return database.NewTextCondition(i.UserIDColumn(), database.TextOperationEqual, userID)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// AssertionColumn implements [domain.idpIntentColumns].
func (i *idpIntent) AssertionColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "assertion")
}

// CreatedAtColumn implements [domain.idpIntentColumns].
func (i *idpIntent) CreatedAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "created_at")
}

// ExpiresAtColumn implements [domain.idpIntentColumns].
func (i *idpIntent) ExpiresAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "expires_at")
}

// FailureURLColumn implements [domain.idpIntentColumns].
func (i *idpIntent) FailureURLColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "failure_url")
}

// IDColumn implements [domain.idpIntentColumns].
func (i *idpIntent) IDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "id")
}

// IDPAccessTokenColumn implements [domain.idpIntentColumns].
func (i *idpIntent) IDPAccessTokenColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_access_token")
}

// IDPArgumentsColumn implements [domain.idpIntentColumns].
func (i *idpIntent) IDPArgumentsColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_arguments")
}

// IDPEntryAttributesColumn implements [domain.idpIntentColumns].
func (i *idpIntent) IDPEntryAttributesColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_entry_attributes")
}

// IDPIDColumn implements [domain.idpIntentColumns].
func (i *idpIntent) IDPIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_id")
}

// IDPUserColumn implements [domain.idpIntentColumns].
func (i *idpIntent) IDPUserColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_user")
}

// IDPUserIDColumn implements [domain.idpIntentColumns].
func (i *idpIntent) IDPUserIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_user_id")
}

// IDPUsernameColumn implements [domain.idpIntentColumns].
func (i *idpIntent) IDPUsernameColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_username")
}

// InstanceIDColumn implements [domain.idpIntentColumns].
func (i *idpIntent) InstanceIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "instance_id")
}

// MaxIDPIntentLifetimeColumn implements [domain.idpIntentColumns].
func (i *idpIntent) MaxIDPIntentLifetimeColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "max_idp_intent_lifetime")
}

// PrimaryKeyColumns implements [domain.IDPIntentRepository].
func (i *idpIntent) PrimaryKeyColumns() []database.Column {
	return []database.Column{i.IDColumn(), i.InstanceIDColumn()}
}

// RequestIDColumn implements [domain.idpIntentColumns].
func (i *idpIntent) RequestIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "request_id")
}

// StateColumn implements [domain.idpIntentColumns].
func (i *idpIntent) StateColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "state")
}

// SucceededAtColumn implements [domain.idpIntentColumns].
func (i *idpIntent) SucceededAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "succeeded_at")
}

// SuccessURLColumn implements [domain.idpIntentColumns].
func (i *idpIntent) SuccessURLColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "success_url")
}

// UpdatedAtColumn implements [domain.idpIntentColumns].
func (i *idpIntent) UpdatedAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "updated_at")
}

// UserIDColumn implements [domain.idpIntentColumns].
func (i *idpIntent) UserIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "user_id")
}

var _ domain.IDPIntentRepository = (*idpIntent)(nil)
