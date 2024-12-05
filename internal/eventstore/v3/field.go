package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type fieldValue struct {
	value []byte
}

func (value *fieldValue) Unmarshal(ptr any) error {
	return json.Unmarshal(value.value, ptr)
}

func (es *Eventstore) FillFields(ctx context.Context, events ...eventstore.FillFieldsEvent) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer span.End()

	tx, err := es.client.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	return handleFieldFillEvents(ctx, tx, events)
}

// Search implements the [eventstore.Search] method
func (es *Eventstore) Search(ctx context.Context, conditions ...map[eventstore.FieldType]any) (result []*eventstore.SearchResult, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer span.EndWithError(err)

	var builder strings.Builder
	args := buildSearchStatement(ctx, &builder, conditions...)

	err = es.client.QueryContext(
		ctx,
		func(rows *sql.Rows) error {
			for rows.Next() {
				var (
					res   eventstore.SearchResult
					value fieldValue
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
					&value.value,
				)
				if err != nil {
					return err
				}
				res.Value = &value

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

const searchQueryPrefix = `SELECT instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value FROM eventstore.fields WHERE instance_id = $1`

func buildSearchStatement(ctx context.Context, builder *strings.Builder, conditions ...map[eventstore.FieldType]any) []any {
	args := make([]any, 0, len(conditions)*4+1)
	args = append(args, authz.GetInstance(ctx).InstanceID())

	builder.WriteString(searchQueryPrefix)

	builder.WriteString(" AND ")
	if len(conditions) > 1 {
		builder.WriteRune('(')
	}
	for i, condition := range conditions {
		if i > 0 {
			builder.WriteString(" OR ")
		}
		if len(condition) > 1 {
			builder.WriteRune('(')
		}
		args = append(args, buildSearchCondition(builder, len(args)+1, condition)...)
		if len(condition) > 1 {
			builder.WriteRune(')')
		}
	}
	if len(conditions) > 1 {
		builder.WriteRune(')')
	}

	return args
}

func buildSearchCondition(builder *strings.Builder, index int, conditions map[eventstore.FieldType]any) []any {
	args := make([]any, 0, len(conditions))

	orderedCondition := make([]eventstore.FieldType, 0, len(conditions))
	for field := range conditions {
		orderedCondition = append(orderedCondition, field)
	}
	slices.Sort(orderedCondition)

	for _, field := range orderedCondition {
		if len(args) > 0 {
			builder.WriteString(" AND ")
		}
		builder.WriteString(fieldNameByType(field, conditions[field]))
		builder.WriteString(" = $")
		builder.WriteString(strconv.Itoa(index + len(args)))
		args = append(args, conditions[field])
	}

	return args
}

func (es *Eventstore) handleFieldCommands(ctx context.Context, tx database.Tx, commands []eventstore.Command) error {
	for _, command := range commands {
		if len(command.Fields()) > 0 {
			if err := handleFieldOperations(ctx, tx, command.Fields()); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleFieldFillEvents(ctx context.Context, tx database.Tx, events []eventstore.FillFieldsEvent) error {
	for _, event := range events {
		if len(event.Fields()) > 0 {
			if err := handleFieldOperations(ctx, tx, event.Fields()); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleFieldOperations(ctx context.Context, tx database.Tx, operations []*eventstore.FieldOperation) error {
	for _, operation := range operations {
		if operation.Set != nil {
			if err := handleFieldSet(ctx, tx, operation.Set); err != nil {
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

func handleFieldSet(ctx context.Context, tx database.Tx, field *eventstore.Field) error {
	if len(field.UpsertConflictFields) == 0 {
		return handleSearchInsert(ctx, tx, field)
	}
	return handleSearchUpsert(ctx, tx, field)
}

const (
	insertField = `INSERT INTO eventstore.fields (instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value, value_must_be_unique, should_index) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
)

func handleSearchInsert(ctx context.Context, tx database.Tx, field *eventstore.Field) error {
	value, err := json.Marshal(field.Value.Value)
	if err != nil {
		return zerrors.ThrowInvalidArgument(err, "V3-fcrW1", "unable to marshal field value")
	}
	_, err = tx.ExecContext(
		ctx,
		insertField,

		field.Aggregate.InstanceID,
		field.Aggregate.ResourceOwner,
		field.Aggregate.Type,
		field.Aggregate.ID,
		field.Object.Type,
		field.Object.ID,
		field.Object.Revision,
		field.FieldName,
		value,
		field.Value.MustBeUnique,
		field.Value.ShouldIndex,
	)
	return err
}

const (
	fieldsUpsertPrefix = `WITH upsert AS (UPDATE eventstore.fields SET (instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value, value_must_be_unique, should_index) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) WHERE `
	fieldsUpsertSuffix = ` RETURNING * ) INSERT INTO eventstore.fields (instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value, value_must_be_unique, should_index) SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11 WHERE NOT EXISTS (SELECT 1 FROM upsert)`
)

func handleSearchUpsert(ctx context.Context, tx database.Tx, field *eventstore.Field) error {
	value, err := json.Marshal(field.Value.Value)
	if err != nil {
		return zerrors.ThrowInvalidArgument(err, "V3-fcrW1", "unable to marshal field value")
	}

	_, err = tx.ExecContext(
		ctx,
		writeUpsertField(field.UpsertConflictFields),

		field.Aggregate.InstanceID,
		field.Aggregate.ResourceOwner,
		field.Aggregate.Type,
		field.Aggregate.ID,
		field.Object.Type,
		field.Object.ID,
		field.Object.Revision,
		field.FieldName,
		value,
		field.Value.MustBeUnique,
		field.Value.ShouldIndex,
	)
	return err
}

func writeUpsertField(fields []eventstore.FieldType) string {
	var builder strings.Builder

	builder.WriteString(fieldsUpsertPrefix)
	for i, fieldName := range fields {
		if i > 0 {
			builder.WriteString(" AND ")
		}
		name, index := searchFieldNameAndIndexByTypeForPush(fieldName)

		builder.WriteString(name)
		builder.WriteString(" = ")
		builder.WriteString(index)
	}
	builder.WriteString(fieldsUpsertSuffix)

	return builder.String()
}

const removeSearch = `DELETE FROM eventstore.fields WHERE `

func handleSearchDelete(ctx context.Context, tx database.Tx, clauses map[eventstore.FieldType]any) error {
	if len(clauses) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "V3-oqlBZ", "no conditions")
	}
	stmt, args := writeDeleteField(clauses)
	_, err := tx.ExecContext(ctx, stmt, args...)
	return err
}

func writeDeleteField(clauses map[eventstore.FieldType]any) (string, []any) {
	var (
		builder strings.Builder
		args    = make([]any, 0, len(clauses))
	)
	builder.WriteString(removeSearch)

	orderedCondition := make([]eventstore.FieldType, 0, len(clauses))
	for field := range clauses {
		orderedCondition = append(orderedCondition, field)
	}
	slices.Sort(orderedCondition)

	for _, fieldName := range orderedCondition {
		if len(args) > 0 {
			builder.WriteString(" AND ")
		}
		builder.WriteString(fieldNameByType(fieldName, clauses[fieldName]))

		builder.WriteString(" = $")
		builder.WriteString(strconv.Itoa(len(args) + 1))

		args = append(args, clauses[fieldName])
	}

	return builder.String(), args
}

func fieldNameByType(typ eventstore.FieldType, value any) string {
	switch typ {
	case eventstore.FieldTypeAggregateID:
		return "aggregate_id"
	case eventstore.FieldTypeAggregateType:
		return "aggregate_type"
	case eventstore.FieldTypeInstanceID:
		return "instance_id"
	case eventstore.FieldTypeResourceOwner:
		return "resource_owner"
	case eventstore.FieldTypeFieldName:
		return "field_name"
	case eventstore.FieldTypeObjectType:
		return "object_type"
	case eventstore.FieldTypeObjectID:
		return "object_id"
	case eventstore.FieldTypeObjectRevision:
		return "object_revision"
	case eventstore.FieldTypeValue:
		return valueColumn(value)
	}
	return ""
}

func searchFieldNameAndIndexByTypeForPush(typ eventstore.FieldType) (string, string) {
	switch typ {
	case eventstore.FieldTypeInstanceID:
		return "instance_id", "$1"
	case eventstore.FieldTypeResourceOwner:
		return "resource_owner", "$2"
	case eventstore.FieldTypeAggregateType:
		return "aggregate_type", "$3"
	case eventstore.FieldTypeAggregateID:
		return "aggregate_id", "$4"
	case eventstore.FieldTypeObjectType:
		return "object_type", "$5"
	case eventstore.FieldTypeObjectID:
		return "object_id", "$6"
	case eventstore.FieldTypeObjectRevision:
		return "object_revision", "$7"
	case eventstore.FieldTypeFieldName:
		return "field_name", "$8"
	case eventstore.FieldTypeValue:
		return "value", "$9"
	}
	return "", ""
}

func valueColumn(value any) string {
	//nolint: exhaustive
	switch reflect.TypeOf(value).Kind() {
	case reflect.Bool:
		return "bool_value"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return "number_value"
	case reflect.String:
		return "text_value"
	}
	return ""
}
