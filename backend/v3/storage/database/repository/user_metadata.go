package repository

import (
	"strings"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userMetadata struct{}

func (userMetadata) qualifiedTableName() string {
	return "zitadel.user_metadata"
}

func (userMetadata) unqualifiedTableName() string {
	return "user_metadata"
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetMetadata implements [domain.UserRepository].
func (u userMetadata) SetMetadata(metadata ...*domain.Metadata) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("INSERT INTO zitadel.user_metadata(instance_id, user_id, key, value, created_at, updated_at)")
			builder.WriteString(" SELECT existing_user.instance_id, existing_user.id, md.key::TEXT, md.value::BYTEA, md.created_at::TIMESTAMPTZ, md.updated_at::TIMESTAMPTZ")
			builder.WriteString(" FROM existing_user CROSS JOIN (VALUES ")
			for i, md := range metadata {
				if i > 0 {
					builder.WriteString(", ")
				}
				builder.WriteRune('(')
				// The complex placeholder handling is needed to append the type casts
				// if we use WriteArgs directly, pgx resolves all args as TEXT instead of the correct type
				placeholders := make([]string, 4)
				placeholders[0] = builder.AppendArg(md.Key)
				placeholders[1] = builder.AppendArg(md.Value) + "::BYTEA"
				var createdAt, updatedAt any = database.NowInstruction, database.NullInstruction
				if !md.CreatedAt.IsZero() {
					createdAt = md.CreatedAt
				}
				if !md.UpdatedAt.IsZero() {
					updatedAt = md.UpdatedAt
				}
				placeholders[2] = builder.AppendArg(createdAt) + "::TIMESTAMPTZ"
				placeholders[3] = builder.AppendArg(updatedAt) + "::TIMESTAMPTZ"
				builder.WriteString(strings.Join(placeholders, ", "))
				builder.WriteRune(')')
			}
			builder.WriteString(") AS md(key, value, created_at, updated_at) ON CONFLICT (")
			database.Columns{
				u.instanceIDColumn(),
				u.userIDColumn(),
				u.keyColumn(),
			}.WriteUnqualified(builder)
			builder.WriteString(") DO UPDATE SET value = EXCLUDED.value, updated_at = EXCLUDED.updated_at WHERE ")
			u.valueColumn().WriteQualified(builder)
			builder.WriteString(" IS DISTINCT FROM EXCLUDED.value")
		},
		nil,
	)
}

// RemoveMetadata implements [domain.UserRepository].
func (u userMetadata) RemoveMetadata(condition database.Condition) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("DELETE FROM zitadel.user_metadata USING ")
			builder.WriteString(existingUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(existingUser.InstanceIDColumn(), u.instanceIDColumn()),
				database.NewColumnCondition(existingUser.IDColumn(), u.userIDColumn()),
				condition,
			))
		},
		nil,
	)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (u userMetadata) MetadataConditions() domain.UserMetadataConditions {
	return u
}

// MetadataKeyCondition implements [domain.UserMetadataConditions].
func (u userMetadata) MetadataKeyCondition(op database.TextOperation, key string) database.Condition {
	return database.NewTextCondition(u.keyColumn(), op, key)
}

// MetadataValueCondition implements [domain.UserMetadataConditions].
func (u userMetadata) MetadataValueCondition(op database.BytesOperation, value []byte) database.Condition {
	return database.NewBytesCondition[[]byte](database.SHA256Column(u.valueColumn()), op, database.SHA256Value(value))
}

var _ domain.UserMetadataConditions = (*userMetadata)(nil)

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userMetadata) instanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "instance_id")
}

func (u userMetadata) userIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "user_id")
}

func (u userMetadata) keyColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "key")
}

func (u userMetadata) valueColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "value")
}
