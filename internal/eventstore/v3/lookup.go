package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"strconv"
	"strings"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func handleLookupCommands(ctx context.Context, tx *sql.Tx, commands []eventstore.Command) error {
	for _, command := range commands {
		if len(command.LookupOperations()) > 0 {
			if err := handleLookupOperations(ctx, tx, command.LookupOperations()); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleLookupOperations(ctx context.Context, tx *sql.Tx, operations []*eventstore.LookupOperation) error {
	var (
		builder strings.Builder
		args    = make([]any, 0, len(operations)*6)
	)

	for _, operation := range operations {
		if operation.Add != nil {
			args = append(args, handleLookupUpsert(&builder, len(args), operation.Add)...)
			continue
		}
		args = append(args, handleLookupDelete(&builder, len(args), operation.Remove)...)
	}

	_, err := tx.ExecContext(ctx, builder.String(), args...)
	return err
}

const (
	insertLookupNumeric = `INSERT INTO eventstore.lookup_fields (instance_id, resource_owner, aggregate_type, aggregate_id, field_name, number_value) VALUES (`
	insertLookupText    = `INSERT INTO eventstore.lookup_fields (instance_id, resource_owner, aggregate_type, aggregate_id, field_name, text_value) VALUES (`
	onConflictLookup    = ` ON CONFLICT (instance_id, resource_owner, aggregate_type, aggregate_id, field_name) DO UPDATE SET number_value = EXCLUDED.number_value, text_value = EXCLUDED.text_value`
)

func handleLookupUpsert(builder *strings.Builder, index int, operation *eventstore.LookupField) []any {
	if operation.Value.IsNumeric() {
		builder.WriteString(insertLookupNumeric)
	} else {
		builder.WriteString(insertLookupText)
	}
	upsertParameters(builder, index)
	builder.WriteRune(')')
	if operation.IsUpsert {
		builder.WriteString(onConflictLookup)
	}
	builder.WriteRune(';')

	return []any{
		operation.Aggregate.InstanceID,
		operation.Aggregate.ResourceOwner,
		operation.Aggregate.Type,
		operation.Aggregate.ID,
		operation.FieldName,
		operation.Value.Value,
	}
}

func upsertParameters(builder *strings.Builder, index int) {
	builder.WriteRune('$')
	builder.WriteString(strconv.Itoa(index + 1))
	builder.WriteString(", $")
	builder.WriteString(strconv.Itoa(index + 2))
	builder.WriteString(", $")
	builder.WriteString(strconv.Itoa(index + 3))
	builder.WriteString(", $")
	builder.WriteString(strconv.Itoa(index + 4))
	builder.WriteString(", $")
	builder.WriteString(strconv.Itoa(index + 5))
	builder.WriteString(", $")
	builder.WriteString(strconv.Itoa(index + 6))
}

const removeLookup = `DELETE FROM eventstore.lookup_fields WHERE `

func handleLookupDelete(builder *strings.Builder, index int, operation map[eventstore.LookupFieldType]*eventstore.LookupValue) []any {
	args := make([]any, 0, len(operation))
	builder.WriteString(removeLookup)

	for fieldName, value := range operation {
		if len(args) > 0 {
			builder.WriteString(" AND ")
		}
		builder.WriteString(fieldNameByType(fieldName, value))
		builder.WriteString(" = $")
		builder.WriteString(strconv.Itoa(index + len(args) + 1))
		args = append(args, value.Value)
	}
	builder.WriteRune(';')
	return args
}

func fieldNameByType(typ eventstore.LookupFieldType, value *eventstore.LookupValue) string {
	switch typ {
	case eventstore.LookupFieldTypeAggregateID:
		return "aggregate_id"
	case eventstore.LookupFieldTypeAggregateType:
		return "aggregate_type"
	case eventstore.LookupFieldTypeInstanceID:
		return "instance_id"
	case eventstore.LookupFieldTypeResourceOwner:
		return "resource_owner"
	case eventstore.LookupFieldTypeFieldName:
		return "field_name"
	case eventstore.LookupFieldTypeValue:
		if value.IsNumeric() {
			return "number_value"
		} else {
			return "text_value"
		}
	}
	return ""
}
