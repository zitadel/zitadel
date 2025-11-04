package repository

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.UserMetadataRepository = (*userMetadata)(nil)

type userMetadata struct{}

func UserMetadataRepository() domain.UserMetadataRepository {
	return new(userMetadata)
}

func (o userMetadata) qualifiedTableName() string {
	return "zitadel.user_metadata"
}

func (o userMetadata) unqualifiedTableName() string {
	return "user_metadata"
}

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

const queryUserMetadataStmt = `SELECT instance_id, user_id, key, value, created_at, updated_at ` +
	`FROM zitadel.user_metadata`

// Get implements [domain.UserMetadataRepository].
func (o userMetadata) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.UserMetadata, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if !options.Condition.IsRestrictingColumn(o.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(o.InstanceIDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(queryUserMetadataStmt)
	options.Write(&builder)

	return scanUserMetadata(ctx, client, &builder)
}

// List implements [domain.UserMetadataRepository].
func (o userMetadata) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.UserMetadata, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if !options.Condition.IsRestrictingColumn(o.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(o.InstanceIDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(queryUserMetadataStmt)
	options.Write(&builder)

	return scanUserMetadataList(ctx, client, &builder)
}

// Set implements [domain.UserMetadataRepository].
func (o userMetadata) Set(ctx context.Context, client database.QueryExecutor, metadata ...*domain.UserMetadata) error {
	if len(metadata) == 0 {
		return database.ErrNoChanges
	}

	var builder database.StatementBuilder
	builder.WriteString(`INSERT INTO zitadel.user_metadata (instance_id, user_id, key, value, created_at, updated_at) VALUES `)
	for i, m := range metadata {
		var createdAt, updatedAt any = database.DefaultInstruction, database.NullInstruction
		if !m.CreatedAt.IsZero() {
			createdAt = m.CreatedAt
		}
		if !m.UpdatedAt.IsZero() {
			updatedAt = m.UpdatedAt
		}
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteRune('(')
		builder.WriteArgs(m.InstanceID, m.UserID, m.Key, m.Value, createdAt, updatedAt)
		builder.WriteRune(')')
	}
	builder.WriteString(` ON CONFLICT (instance_id, user_id, key) DO UPDATE SET value = EXCLUDED.value, updated_at = EXCLUDED.updated_at RETURNING created_at, updated_at`)

	res, err := client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return err
	}
	if res.Err() != nil {
		return res.Err()
	}
	dates := make([]changedDates, 0, len(metadata))
	err = res.(database.CollectableRows).Collect(&dates)
	if err != nil {
		return err
	}
	if len(dates) != len(metadata) {
		return errors.New("repository: returned changed dates count does not match list count")
	}
	for i, d := range dates {
		metadata[i].CreatedAt = d.CreatedAt
		metadata[i].UpdatedAt = d.UpdatedAt
	}
	return nil
}

// Remove implements [domain.UserMetadataRepository].
func (o userMetadata) Remove(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if err := checkPKCondition(o, condition); err != nil {
		return 0, err
	}

	var builder database.StatementBuilder
	builder.WriteString(`DELETE FROM zitadel.user_metadata `)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// PrimaryKeyCondition implements [domain.UserMetadataRepository].
func (o userMetadata) PrimaryKeyCondition(instanceID string, userID string, key string) database.Condition {
	return database.And(
		o.InstanceIDCondition(instanceID),
		o.UserIDCondition(userID),
		o.KeyCondition(database.TextOperationEqual, key),
	)
}

// InstanceIDCondition implements [domain.UserMetadataRepository].
func (o userMetadata) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(o.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// KeyCondition implements [domain.UserMetadataRepository].
func (o userMetadata) KeyCondition(op database.TextOperation, key string) database.Condition {
	return database.NewTextCondition(o.KeyColumn(), op, key)
}

// UserIDCondition implements [domain.UserMetadataRepository].
func (o userMetadata) UserIDCondition(userID string) database.Condition {
	return database.NewTextCondition(o.UserIDColumn(), database.TextOperationEqual, userID)
}

// ValueCondition implements [domain.UserMetadataRepository].
func (o userMetadata) ValueCondition(op database.BytesOperation, value []byte) database.Condition {
	return database.NewBytesCondition[[]byte](database.SHA256Column(o.ValueColumn()), op, database.SHA256Value(value))
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// PrimaryKeyColumns implements [domain.UserMetadataRepository].
func (o userMetadata) PrimaryKeyColumns() []database.Column {
	return []database.Column{o.InstanceIDColumn(), o.UserIDColumn(), o.KeyColumn()}
}

// CreatedAtColumn implements [domain.UserMetadataRepository].
func (o userMetadata) CreatedAtColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "created_at")
}

// InstanceIDColumn implements [domain.UserMetadataRepository].
func (o userMetadata) InstanceIDColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "instance_id")
}

// KeyColumn implements [domain.UserMetadataRepository].
func (o userMetadata) KeyColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "key")
}

// UserIDColumn implements [domain.UserMetadataRepository].
func (o userMetadata) UserIDColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "user_id")
}

// UpdatedAtColumn implements [domain.UserMetadataRepository].
func (o userMetadata) UpdatedAtColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "updated_at")
}

// ValueColumn implements [domain.UserMetadataRepository].
func (o userMetadata) ValueColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "value")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

func scanUserMetadata(ctx context.Context, client database.Querier, builder *database.StatementBuilder) (*domain.UserMetadata, error) {
	rows, err := client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	metadata := new(domain.UserMetadata)
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

func scanUserMetadataList(ctx context.Context, client database.Querier, builder *database.StatementBuilder) ([]*domain.UserMetadata, error) {
	rows, err := client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var metadataList []*domain.UserMetadata
	if err := rows.(database.CollectableRows).Collect(&metadataList); err != nil {
		return nil, err
	}
	return metadataList, nil
}
