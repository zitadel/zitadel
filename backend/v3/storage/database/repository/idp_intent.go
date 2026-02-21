package repository

import (
	"context"
	"net/url"
	"time"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type idpIntentRepository struct{}

func (i idpIntentRepository) qualifiedTableName() string {
	return "zitadel.identity_provider_intents"
}

func (i idpIntentRepository) unqualifiedTableName() string {
	return "identity_provider_intents"
}

func IDPIntentRepository() domain.IDPIntentRepository {
	return new(idpIntentRepository)
}

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

const queryIDPIntentStmt = `
SELECT
	identity_provider_intents.instance_id
	, identity_provider_intents.id
	, identity_provider_intents.state
	, identity_provider_intents.success_url
	, identity_provider_intents.failure_url
	, identity_provider_intents.created_at
	, identity_provider_intents.updated_at
	, identity_provider_intents.idp_id
	, identity_provider_intents.idp_arguments
	, identity_provider_intents.idp_user
	, identity_provider_intents.idp_user_id
	, identity_provider_intents.idp_username
	, identity_provider_intents.user_id
	, identity_provider_intents.idp_access_token
	, identity_provider_intents.idp_id_token
	, identity_provider_intents.idp_entry_attributes
	, identity_provider_intents.request_id
	, identity_provider_intents.assertion
	, identity_provider_intents.succeeded_at
	, identity_provider_intents.failed_at
	, identity_provider_intents.fail_reason
	, identity_provider_intents.expires_at
FROM
	zitadel.identity_provider_intents
`

// Create implements [domain.IDPIntentRepository].
func (i idpIntentRepository) Create(ctx context.Context, client database.QueryExecutor, intent *domain.IDPIntent) error {
	var createdAt any = database.DefaultInstruction
	if !intent.CreatedAt.IsZero() {
		createdAt = intent.CreatedAt
	}
	var updatedAt = createdAt
	if !intent.UpdatedAt.IsZero() {
		updatedAt = intent.UpdatedAt
	}

	builder := new(database.StatementBuilder)
	builder.WriteString(`INSERT INTO ` + i.qualifiedTableName() + ` (instance_id, id, success_url, failure_url, idp_id, idp_arguments, created_at, updated_at) VALUES ( `)
	builder.WriteArgs(intent.InstanceID, intent.ID, intent.SuccessURL.String(), intent.FailureURL.String(), intent.IDPID, intent.IDPArguments, createdAt, updatedAt)
	builder.WriteString(` ) RETURNING created_at, updated_at`)
	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&intent.CreatedAt, &intent.UpdatedAt)
}

// Delete implements [domain.IDPIntentRepository].
func (i idpIntentRepository) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return deleteOne(ctx, client, i, condition)
}

// Get implements [domain.IDPIntentRepository].
func (i idpIntentRepository) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPIntent, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if options.Condition == nil || !options.Condition.IsRestrictingColumn(i.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(i.InstanceIDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(queryIDPIntentStmt)
	options.Write(&builder)

	return scanIDPIntent(ctx, client, &builder)
}

type rawIDPIntent struct {
	*domain.IDPIntent
	SuccessURL  string  `json:"success_url,omitempty" db:"success_url"`
	FailureURL  string  `json:"failure_url,omitempty" db:"failure_url"`
	IDPUserID   *string `json:"idp_user_id,omitempty" db:"idp_user_id"`
	IDPUsername *string `json:"idp_username,omitempty" db:"idp_username"`
	UserID      *string `json:"user_id,omitempty" db:"user_id"`
	IDPIDToken  *string `json:"idp_id_token,omitempty" db:"idp_id_token"`
	RequestID   *string `json:"request_id,omitempty" db:"request_id"`
	FailReason  *string `json:"fail_reason,omitempty" db:"fail_reason"`
}

func scanIDPIntent(ctx context.Context, querier database.QueryExecutor, builder *database.StatementBuilder) (*domain.IDPIntent, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	raw := new(rawIDPIntent)
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(raw); err != nil {
		return nil, err
	}
	return rawIDPIntentToDomain(raw)
}

func rawIDPIntentToDomain(raw *rawIDPIntent) (*domain.IDPIntent, error) {
	successURL, _ := url.Parse(raw.SuccessURL)
	failureURL, _ := url.Parse(raw.FailureURL)
	raw.IDPIntent.SuccessURL = successURL
	raw.IDPIntent.FailureURL = failureURL
	raw.IDPIntent.IDPUserID = gu.Value(raw.IDPUserID)
	raw.IDPIntent.IDPUsername = gu.Value(raw.IDPUsername)
	raw.IDPIntent.UserID = gu.Value(raw.UserID)
	raw.IDPIntent.IDPIDToken = gu.Value(raw.IDPIDToken)
	raw.IDPIntent.RequestID = gu.Value(raw.RequestID)
	raw.IDPIntent.FailReason = gu.Value(raw.FailReason)

	return raw.IDPIntent, nil
}

// Update implements [domain.IDPIntentRepository].
func (i idpIntentRepository) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if changes == nil {
		return 0, database.ErrNoChanges
	}
	if condition == nil || !condition.IsRestrictingColumn(i.InstanceIDColumn()) {
		return 0, database.NewMissingConditionError(i.InstanceIDColumn())
	}
	if !database.Changes(changes).IsOnColumn(i.UpdatedAtColumn()) {
		changes = append(changes, database.NewChange(i.UpdatedAtColumn(), database.DefaultInstruction))
	}

	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.identity_provider_intents SET `)

	err := database.Changes(changes).Write(&builder)
	if err != nil {
		return 0, err
	}
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetAssertion implements [domain.idpIntentChanges].
func (i idpIntentRepository) SetAssertion(assertion []byte) database.Change {
	return database.NewChange(i.AssertionColumn(), assertion)
}

// SetFailureURL implements [domain.idpIntentChanges].
func (i idpIntentRepository) SetFailureURL(failureURL url.URL) database.Change {
	return database.NewChange(i.FailureURLColumn(), failureURL.String())
}

// SetIDPAccessToken implements [domain.idpIntentChanges].
func (i idpIntentRepository) SetIDPAccessToken(idpAccessToken []byte) database.Change {
	return database.NewChange(i.IDPAccessTokenColumn(), idpAccessToken)
}

// SetIDPIDToken implements [domain.IDPIntentRepository].
func (i *idpIntentRepository) SetIDPIDToken(idpIDToken string) database.Change {
	return database.NewChange(i.IDPIDTokenColumn(), idpIDToken)
}

// SetIDPArguments implements [domain.idpIntentChanges].
func (i idpIntentRepository) SetIDPArguments(idpArguments []byte) database.Change {
	return database.NewChange(i.IDPArgumentsColumn(), idpArguments)
}

// SetIDPEntryAttributes implements [domain.idpIntentChanges].
func (i idpIntentRepository) SetIDPEntryAttributes(idpEntryAttributes []byte) database.Change {
	return database.NewChange(i.IDPEntryAttributesColumn(), idpEntryAttributes)
}

// SetIDPID implements [domain.idpIntentChanges].
func (i idpIntentRepository) SetIDPID(idpID string) database.Change {
	return database.NewChange(i.IDPIDColumn(), idpID)
}

// SetIDPUser implements [domain.idpIntentChanges].
func (i idpIntentRepository) SetIDPUser(idpUser []byte) database.Change {
	return database.NewChange(i.IDPUserColumn(), idpUser)
}

// SetIDPUserID implements [domain.idpIntentChanges].
func (i idpIntentRepository) SetIDPUserID(idpUserID string) database.Change {
	return database.NewChange(i.IDPUserIDColumn(), idpUserID)
}

// SetIDPUsername implements [domain.idpIntentChanges].
func (i idpIntentRepository) SetIDPUsername(idpUsername string) database.Change {
	return database.NewChange(i.IDPUsernameColumn(), idpUsername)
}

// SetRequestID implements [domain.idpIntentChanges].
func (i idpIntentRepository) SetRequestID(requestID string) database.Change {
	return database.NewChange(i.RequestIDColumn(), requestID)
}

// SetState implements [domain.idpIntentChanges].
func (i idpIntentRepository) SetState(state domain.IDPIntentState) database.Change {
	return database.NewChange(i.StateColumn(), state)
}

// SetSuccessURL implements [domain.idpIntentChanges].
func (i idpIntentRepository) SetSuccessURL(successURL url.URL) database.Change {
	return database.NewChange(i.SuccessURLColumn(), successURL.String())
}

// SetUserID implements [domain.idpIntentChanges].
func (i idpIntentRepository) SetUserID(userID string) database.Change {
	return database.NewChange(i.UserIDColumn(), userID)
}

// SetExpiresAt implements [domain.IDPIntentRepository].
func (i *idpIntentRepository) SetExpiresAt(expiration time.Time) database.Change {
	return database.NewChange(i.ExpiresAtColumn(), expiration)
}

// SetSucceededAt implements [domain.IDPIntentRepository].
func (i *idpIntentRepository) SetSucceededAt(succeededAt time.Time) database.Change {
	return database.NewChange(i.SucceededAtColumn(), succeededAt)
}

// SetFailedAt implements [domain.IDPIntentRepository].
func (i *idpIntentRepository) SetFailedAt(failedAt time.Time) database.Change {
	return database.NewChange(i.FailedAtColumn(), failedAt)
}

// SetFailReason implements [domain.IDPIntentRepository].
func (i *idpIntentRepository) SetFailReason(reason string) database.Change {
	return database.NewChange(i.FailReasonColumn(), reason)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// CreatedAtCondition implements [domain.idpIntentConditions].
func (i idpIntentRepository) CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition {
	return database.NewNumberCondition(i.CreatedAtColumn(), op, createdAt)
}

// IDCondition implements [domain.idpIntentConditions].
func (i idpIntentRepository) IDCondition(id string) database.Condition {
	return database.NewTextCondition(i.IDColumn(), database.TextOperationEqual, id)
}

// IDPIDCondition implements [domain.idpIntentConditions].
func (i idpIntentRepository) IDPIDCondition(idpID string) database.Condition {
	return database.NewTextCondition(i.IDPIDColumn(), database.TextOperationEqual, idpID)
}

// IDPUserIDCondition implements [domain.idpIntentConditions].
func (i idpIntentRepository) IDPUserIDCondition(idpUserID string) database.Condition {
	return database.NewTextCondition(i.IDPUserIDColumn(), database.TextOperationEqual, idpUserID)
}

// IDPUsernameCondition implements [domain.idpIntentConditions].
func (i idpIntentRepository) IDPUsernameCondition(idpUsername string) database.Condition {
	return database.NewTextCondition(i.IDPUsernameColumn(), database.TextOperationEqual, idpUsername)
}

// InstanceIDCondition implements [domain.idpIntentConditions].
func (i idpIntentRepository) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(i.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// PrimaryKeyCondition implements [domain.idpIntentConditions].
func (i idpIntentRepository) PrimaryKeyCondition(instanceID, id string) database.Condition {
	return database.And(
		i.InstanceIDCondition(instanceID),
		i.IDCondition(id),
	)
}

// RequestIDCondition implements [domain.idpIntentConditions].
func (i idpIntentRepository) RequestIDCondition(requestID string) database.Condition {
	return database.NewTextCondition(i.RequestIDColumn(), database.TextOperationEqual, requestID)
}

// StateCondition implements [domain.idpIntentConditions].
func (i idpIntentRepository) StateCondition(state domain.IDPIntentState) database.Condition {
	return database.NewTextCondition(i.StateColumn(), database.TextOperationEqual, state.String())
}

// UpdatedAtCondition implements [domain.idpIntentConditions].
func (i idpIntentRepository) UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition {
	return database.NewNumberCondition(i.UpdatedAtColumn(), op, updatedAt)
}

// UserIDCondition implements [domain.idpIntentConditions].
func (i idpIntentRepository) UserIDCondition(userID string) database.Condition {
	return database.NewTextCondition(i.UserIDColumn(), database.TextOperationEqual, userID)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// AssertionColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) AssertionColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "assertion")
}

// CreatedAtColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) CreatedAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "created_at")
}

// ExpiresAtColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) ExpiresAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "expires_at")
}

// FailureURLColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) FailureURLColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "failure_url")
}

// IDColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) IDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "id")
}

// IDPAccessTokenColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) IDPAccessTokenColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_access_token")
}

// IDPIDTokenColumn implements [domain.IDPIntentRepository].
func (i *idpIntentRepository) IDPIDTokenColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_id_token")
}

// IDPArgumentsColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) IDPArgumentsColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_arguments")
}

// IDPEntryAttributesColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) IDPEntryAttributesColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_entry_attributes")
}

// IDPIDColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) IDPIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_id")
}

// IDPUserColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) IDPUserColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_user")
}

// IDPUserIDColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) IDPUserIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_user_id")
}

// IDPUsernameColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) IDPUsernameColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "idp_username")
}

// InstanceIDColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) InstanceIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "instance_id")
}

// PrimaryKeyColumns implements [domain.IDPIntentRepository].
func (i idpIntentRepository) PrimaryKeyColumns() []database.Column {
	return []database.Column{i.IDColumn(), i.InstanceIDColumn()}
}

// RequestIDColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) RequestIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "request_id")
}

// StateColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) StateColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "state")
}

// SucceededAtColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) SucceededAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "succeeded_at")
}

// SuccessURLColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) SuccessURLColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "success_url")
}

// FailReasonColumn implements [domain.IDPIntentRepository].
func (i *idpIntentRepository) FailReasonColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "fail_reason")
}

// FailedAtColumn implements [domain.IDPIntentRepository].
func (i *idpIntentRepository) FailedAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "failed_at")
}

// UpdatedAtColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) UpdatedAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "updated_at")
}

// UserIDColumn implements [domain.idpIntentColumns].
func (i idpIntentRepository) UserIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "user_id")
}

var _ domain.IDPIntentRepository = (*idpIntentRepository)(nil)
