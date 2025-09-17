package repository

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------
var _ domain.OrganizationMetadataRepository = (*orgMetadata)(nil)

type orgMetadata struct {
	repository
	org *org
}

const queryOrganizationMetadataStmt = `SELECT instance_id, org_id, key, value, created_at, updated_at ` +
	`FROM zitadel.org_metadata`

// Get implements [domain.OrganizationMetadataRepository].
func (o *orgMetadata) Get(ctx context.Context, opts ...database.QueryOption) (*domain.OrganizationMetadata, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if !options.Condition.ContainsColumn(o.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(o.InstanceIDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(queryOrganizationMetadataStmt)
	options.Write(&builder)

	return scanOrganizationMetadata(ctx, o.client, &builder)
}

// List implements [domain.OrganizationMetadataRepository].
func (o *orgMetadata) List(ctx context.Context, opts ...database.QueryOption) ([]*domain.OrganizationMetadata, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if !options.Condition.ContainsColumn(o.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(o.InstanceIDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(queryOrganizationMetadataStmt)
	options.Write(&builder)

	return scanOrganizationMetadataList(ctx, o.client, &builder)
}

// Set implements [domain.OrganizationMetadataRepository].
func (o *orgMetadata) Set(ctx context.Context, metadata ...*domain.OrganizationMetadata) error {
	if len(metadata) == 0 {
		return database.ErrNoChanges
	}

	var builder database.StatementBuilder
	builder.WriteString(`INSERT INTO zitadel.org_metadata (instance_id, org_id, key, value, created_at, updated_at) VALUES `)
	for i, m := range metadata {
		var createdAt, updatedAt any = database.DefaultInstruction, database.DefaultInstruction
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
		builder.WriteArgs(m.InstanceID, m.OrgID, m.Key, m.Value, createdAt, updatedAt)
		builder.WriteRune(')')
	}
	builder.WriteString(` ON CONFLICT (instance_id, org_id, key) DO UPDATE SET value = EXCLUDED.value, updated_at = EXCLUDED.updated_at RETURNING created_at, updated_at`)

	res, err := o.client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil || res.Err() != nil {
		return errors.Join(err, res.Err())
	}
	return res.(database.CollectableRows).Collect(&metadata)
}

// Remove implements [domain.OrganizationMetadataRepository].
func (o *orgMetadata) Remove(ctx context.Context, condition database.Condition) (int64, error) {
	var builder database.StatementBuilder
	if !condition.ContainsColumn(o.InstanceIDColumn()) {
		return 0, database.NewMissingConditionError(o.InstanceIDColumn())
	}
	if !condition.ContainsColumn(o.OrgIDColumn()) {
		return 0, database.NewMissingConditionError(o.OrgIDColumn())
	}

	builder.WriteString(`DELETE FROM zitadel.org_metadata `)
	writeCondition(&builder, condition)

	return o.client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// InstanceIDCondition implements [domain.OrganizationMetadataRepository].
func (o *orgMetadata) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(o.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// KeyCondition implements [domain.OrganizationMetadataRepository].
func (o *orgMetadata) KeyCondition(op database.TextOperation, key string) database.Condition {
	return database.NewTextCondition(o.KeyColumn(), op, key)
}

// OrgIDCondition implements [domain.OrganizationMetadataRepository].
func (o *orgMetadata) OrgIDCondition(orgID string) database.Condition {
	return database.NewTextCondition(o.OrgIDColumn(), database.TextOperationEqual, orgID)
}

// ValueCondition implements [domain.OrganizationMetadataRepository].
func (o *orgMetadata) ValueCondition(op database.BytesOperation, value []byte) database.Condition {
	// return database.NewValueCondition(o.ValueColumn(), value)
	return database.NewBytesCondition(o.ValueColumn(), op, value)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// CreatedAtColumn implements [domain.OrganizationMetadataRepository].
func (o *orgMetadata) CreatedAtColumn() *database.Column {
	return database.NewColumn("org_metadata", "created_at")
}

// InstanceIDColumn implements [domain.OrganizationMetadataRepository].
func (o *orgMetadata) InstanceIDColumn() *database.Column {
	return database.NewColumn("org_metadata", "instance_id")
}

// KeyColumn implements [domain.OrganizationMetadataRepository].
func (o *orgMetadata) KeyColumn() *database.Column {
	return database.NewColumn("org_metadata", "key")
}

// OrgIDColumn implements [domain.OrganizationMetadataRepository].
func (o *orgMetadata) OrgIDColumn() *database.Column {
	return database.NewColumn("org_metadata", "org_id")
}

// UpdatedAtColumn implements [domain.OrganizationMetadataRepository].
func (o *orgMetadata) UpdatedAtColumn() *database.Column {
	return database.NewColumn("org_metadata", "updated_at")
}

// ValueColumn implements [domain.OrganizationMetadataRepository].
func (o *orgMetadata) ValueColumn() *database.Column {
	return database.NewColumn("org_metadata", "value")
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
