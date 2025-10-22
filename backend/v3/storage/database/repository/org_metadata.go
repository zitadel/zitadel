package repository

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.OrganizationMetadataRepository = (*orgMetadata)(nil)

type orgMetadata struct{}

func OrganizationMetadataRepository() domain.OrganizationMetadataRepository {
	return new(orgMetadata)
}

func (o orgMetadata) qualifiedTableName() string {
	return "zitadel.organization_metadata"
}

func (o orgMetadata) unqualifiedTableName() string {
	return "organization_metadata"
}

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

const queryOrganizationMetadataStmt = `SELECT instance_id, organization_id, key, value, created_at, updated_at ` +
	`FROM zitadel.organization_metadata`

// Get implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.OrganizationMetadata, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if !options.Condition.IsRestrictingColumn(o.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(o.InstanceIDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(queryOrganizationMetadataStmt)
	options.Write(&builder)

	return scanOrganizationMetadata(ctx, client, &builder)
}

// List implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.OrganizationMetadata, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if !options.Condition.IsRestrictingColumn(o.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(o.InstanceIDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(queryOrganizationMetadataStmt)
	options.Write(&builder)

	return scanOrganizationMetadataList(ctx, client, &builder)
}

// Set implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) Set(ctx context.Context, client database.QueryExecutor, metadata ...*domain.OrganizationMetadata) error {
	if len(metadata) == 0 {
		return database.ErrNoChanges
	}

	var builder database.StatementBuilder
	builder.WriteString(`INSERT INTO zitadel.organization_metadata (instance_id, organization_id, key, value, created_at, updated_at) VALUES `)
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
		builder.WriteArgs(m.InstanceID, m.OrganizationID, m.Key, m.Value, createdAt, updatedAt)
		builder.WriteRune(')')
	}
	builder.WriteString(` ON CONFLICT (instance_id, organization_id, key) DO UPDATE SET value = EXCLUDED.value, updated_at = EXCLUDED.updated_at RETURNING created_at, updated_at`)

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

// Remove implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) Remove(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if err := checkPKCondition(o, condition); err != nil {
		return 0, err
	}

	var builder database.StatementBuilder
	builder.WriteString(`DELETE FROM zitadel.organization_metadata `)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// PrimaryKeyCondition implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) PrimaryKeyCondition(instanceID string, orgID string, key string) database.Condition {
	return database.And(
		o.InstanceIDCondition(instanceID),
		o.OrganizationIDCondition(orgID),
		o.KeyCondition(database.TextOperationEqual, key),
	)
}

// InstanceIDCondition implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(o.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// KeyCondition implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) KeyCondition(op database.TextOperation, key string) database.Condition {
	return database.NewTextCondition(o.KeyColumn(), op, key)
}

// OrganizationIDCondition implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) OrganizationIDCondition(orgID string) database.Condition {
	return database.NewTextCondition(o.OrganizationIDColumn(), database.TextOperationEqual, orgID)
}

// ValueCondition implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) ValueCondition(op database.BytesOperation, value []byte) database.Condition {
	return database.NewBytesCondition[[]byte](database.SHA256Column(o.ValueColumn()), op, database.SHA256Value(value))
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// PrimaryKeyColumns implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) PrimaryKeyColumns() []database.Column {
	return []database.Column{o.InstanceIDColumn(), o.OrganizationIDColumn(), o.KeyColumn()}
}

// CreatedAtColumn implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) CreatedAtColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "created_at")
}

// InstanceIDColumn implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) InstanceIDColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "instance_id")
}

// KeyColumn implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) KeyColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "key")
}

// OrganizationIDColumn implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) OrganizationIDColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "organization_id")
}

// UpdatedAtColumn implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) UpdatedAtColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "updated_at")
}

// ValueColumn implements [domain.OrganizationMetadataRepository].
func (o orgMetadata) ValueColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "value")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

func scanOrganizationMetadata(ctx context.Context, client database.Querier, builder *database.StatementBuilder) (*domain.OrganizationMetadata, error) {
	rows, err := client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	metadata := new(domain.OrganizationMetadata)
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

func scanOrganizationMetadataList(ctx context.Context, client database.Querier, builder *database.StatementBuilder) ([]*domain.OrganizationMetadata, error) {
	rows, err := client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var metadataList []*domain.OrganizationMetadata
	if err := rows.(database.CollectableRows).Collect(&metadataList); err != nil {
		return nil, err
	}
	return metadataList, nil
}
