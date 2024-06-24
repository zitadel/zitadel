package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (es *Eventstore) Search(ctx context.Context, conditions ...map[eventstore.SearchFieldType]any) (result []*eventstore.SearchResult, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer span.EndWithError(err)

	var builder strings.Builder
	args := buildSearchStatement(ctx, &builder, conditions...)

	err = es.client.QueryContext(
		ctx,
		func(rows *sql.Rows) error {
			for rows.Next() {
				var (
					res         eventstore.SearchResult
					textValue   sql.Null[string]
					numberValue pgtype.Numeric
				)

				err = rows.Scan(
					&res.Aggregate.InstanceID,
					&res.Aggregate.ResourceOwner,
					&res.Aggregate.Type,
					&res.Aggregate.ID,
					&res.Object.Type,
					&res.Object.ID,
					&res.Object.Revision,
					&res.FieldName,
					&textValue,
					&numberValue,
				)
				if err != nil {
					return err
				}

				if numberValue.Valid {
					value, err := numberValue.MarshalJSON()
					if err != nil {
						return err
					}
					res.Value = value
					res.ValueType = eventstore.SearchValueTypeNumeric
				} else if textValue.Valid {
					res.Value = []byte(textValue.V)
					res.ValueType = eventstore.SearchValueTypeString
				}
				result = append(result, &res)
			}
			return nil
		},
		builder.String(),
		args...,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

const searchQueryPrefix = `SELECT instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, text_value, number_value FROM eventstore.search WHERE instance_id = $1`

func buildSearchStatement(ctx context.Context, builder *strings.Builder, conditions ...map[eventstore.SearchFieldType]any) []any {
	args := make([]any, 0, len(conditions)*4+1)
	args = append(args, authz.GetInstance(ctx).InstanceID())

	builder.WriteString(searchQueryPrefix)

	if len(conditions) == 0 {
		return args
	}

	builder.WriteString(" AND (")

	for i, condition := range conditions {
		if i > 0 {
			builder.WriteString(" OR ")
		}
		builder.WriteRune('(')
		args = append(args, buildSearchCondition(builder, len(args)+1, condition)...)
		builder.WriteRune(')')
	}
	builder.WriteRune(')')

	return args
}

func buildSearchCondition(builder *strings.Builder, index int, conditions map[eventstore.SearchFieldType]any) []any {
	args := make([]any, 0, len(conditions))

	for field, value := range conditions {
		if len(args) > 0 {
			builder.WriteString(" AND ")
		}
		builder.WriteString(searchFieldNameByType(field))
		builder.WriteString(" = $")
		builder.WriteString(strconv.Itoa(index + len(args)))
		args = append(args, value)
	}

	return args
}

func handleSearchCommands(ctx context.Context, tx *sql.Tx, commands []eventstore.Command) error {
	for _, command := range commands {
		if len(command.SearchOperations()) > 0 {
			if err := handleSearchOperations(ctx, tx, command.SearchOperations()); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleSearchOperations(ctx context.Context, tx *sql.Tx, operations []*eventstore.SearchOperation) error {
	for _, operation := range operations {
		if operation.Set != nil {
			if err := handleSearchSet(ctx, tx, operation.Set); err != nil {
				return err
			}
			continue
		}
		if operation.Remove != nil {
			if err := handleSearchDelete(ctx, tx, operation.Remove); err != nil {
				return err
			}
		}
	}

	return nil
}

func handleSearchSet(ctx context.Context, tx *sql.Tx, field *eventstore.SearchField) error {
	if len(field.UpsertConflictFields) == 0 {
		return handleSearchInsert(ctx, tx, field)
	}
	return handleSearchUpsert(ctx, tx, field)
}

const (
	insertSearchNumeric = `INSERT INTO eventstore.search (instance_id, resource_owner, aggregate_type, aggregate_id, field_name, object_type, object_id, object_revision, number_value) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	insertSearchText    = `INSERT INTO eventstore.search (instance_id, resource_owner, aggregate_type, aggregate_id, field_name, object_type, object_id, object_revision, text_value) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
)

func handleSearchInsert(ctx context.Context, tx *sql.Tx, field *eventstore.SearchField) error {
	stmt := insertSearchText
	if field.Value.IsNumeric() {
		stmt = insertSearchNumeric
	}

	_, err := tx.ExecContext(
		ctx,
		stmt,

		field.Aggregate.InstanceID,
		field.Aggregate.ResourceOwner,
		field.Aggregate.Type,
		field.Aggregate.ID,
		field.FieldName,
		field.Object.Type,
		field.Object.ID,
		field.Object.Revision,
		field.Value.Value,
	)
	return err
}

const (
	searchUpsertNumericPrefix = `WITH upsert AS (UPDATE eventstore.search SET (instance_id, resource_owner, aggregate_type, aggregate_id, field_name, object_type, object_id, object_revision, number_value) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE `
	searchUpsertTextPrefix    = `WITH upsert AS (UPDATE eventstore.search SET (instance_id, resource_owner, aggregate_type, aggregate_id, field_name, object_type, object_id, object_revision, text_value) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE `

	searchUpsertNumericSuffix = ` RETURNING * ) INSERT INTO eventstore.search (instance_id, resource_owner, aggregate_type, aggregate_id, field_name, object_type, object_id, object_revision, number_value) SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9 WHERE NOT EXISTS (SELECT 1 FROM upsert)`
	searchUpsertTextSuffix    = ` RETURNING * ) INSERT INTO eventstore.search (instance_id, resource_owner, aggregate_type, aggregate_id, field_name, object_type, object_id, object_revision, text_value) SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9 WHERE NOT EXISTS (SELECT 1 FROM upsert)`
)

func handleSearchUpsert(ctx context.Context, tx *sql.Tx, field *eventstore.SearchField) error {
	var builder strings.Builder
	if field.Value.IsNumeric() {
		builder.WriteString(searchUpsertNumericPrefix)
	} else {
		builder.WriteString(searchUpsertTextPrefix)
	}

	for i, field := range field.UpsertConflictFields {
		if i > 0 {
			builder.WriteString(" AND ")
		}
		name, index := searchFieldNameAndIndexByTypeForPush(field)

		builder.WriteString(name)
		builder.WriteString(" = ")
		builder.WriteString(index)
	}

	if field.Value.IsNumeric() {
		builder.WriteString(searchUpsertNumericSuffix)
	} else {
		builder.WriteString(searchUpsertTextSuffix)
	}

	_, err := tx.ExecContext(
		ctx,
		builder.String(),

		field.Aggregate.InstanceID,
		field.Aggregate.ResourceOwner,
		field.Aggregate.Type,
		field.Aggregate.ID,
		field.FieldName,
		field.Object.Type,
		field.Object.ID,
		field.Object.Revision,
		field.Value.Value,
	)
	return err
}

const removeSearch = `DELETE FROM eventstore.search WHERE `

func handleSearchDelete(ctx context.Context, tx *sql.Tx, clauses map[eventstore.SearchFieldType]any) error {
	var (
		builder strings.Builder
		args    = make([]any, 0, len(clauses))
	)
	builder.WriteString(removeSearch)

	for fieldName, value := range clauses {
		if len(args) > 0 {
			builder.WriteString(" AND ")
		}
		builder.WriteString(searchFieldNameByType(fieldName))

		builder.WriteString(" = $")
		builder.WriteString(strconv.Itoa(len(args) + 1))

		args = append(args, value)
	}
	_, err := tx.ExecContext(ctx, builder.String(), args...)
	return err
}

func searchFieldNameByType(typ eventstore.SearchFieldType) string {
	switch typ {
	case eventstore.SearchFieldTypeAggregateID:
		return "aggregate_id"
	case eventstore.SearchFieldTypeAggregateType:
		return "aggregate_type"
	case eventstore.SearchFieldTypeInstanceID:
		return "instance_id"
	case eventstore.SearchFieldTypeResourceOwner:
		return "resource_owner"
	case eventstore.SearchFieldTypeFieldName:
		return "field_name"
	case eventstore.SearchFieldTypeObjectType:
		return "object_type"
	case eventstore.SearchFieldTypeObjectID:
		return "object_id"
	case eventstore.SearchFieldTypeObjectRevision:
		return "object_revision"
	case eventstore.SearchFieldTypeTextValue:
		return "text_value"
	case eventstore.SearchFieldTypeNumericValue:
		return "number_value"
	}
	return ""
}

func searchFieldNameAndIndexByTypeForPush(typ eventstore.SearchFieldType) (string, string) {
	switch typ {
	case eventstore.SearchFieldTypeAggregateID:
		return "aggregate_id", "$4"
	case eventstore.SearchFieldTypeAggregateType:
		return "aggregate_type", "$3"
	case eventstore.SearchFieldTypeInstanceID:
		return "instance_id", "$1"
	case eventstore.SearchFieldTypeResourceOwner:
		return "resource_owner", "$2"
	case eventstore.SearchFieldTypeFieldName:
		return "field_name", "$5"
	case eventstore.SearchFieldTypeObjectType:
		return "object_type", "$6"
	case eventstore.SearchFieldTypeObjectID:
		return "object_id", "$7"
	case eventstore.SearchFieldTypeObjectRevision:
		return "object_revision", "$8"
	case eventstore.SearchFieldTypeTextValue:
		return "text_value", "$9"
	case eventstore.SearchFieldTypeNumericValue:
		return "number_value", "$9"
	}
	return "", ""
}
