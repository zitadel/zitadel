package eventstore

import (
	"context"
	_ "embed"
	"encoding/json"
	"strconv"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

//go:embed field_batch_set.sql
var fieldBatchSet string

//go:embed field_batch_remove.sql
var fieldBatchRemove string

// Only full-identity operations without value uniqueness can be reordered within a statement.
func batchFieldKey(operation *eventstore.FieldOperation) (key [8]string, eligible bool) {
	if field := operation.Set; field != nil {
		if field.Value.MustBeUnique || len(field.UpsertConflictFields) != len(key) {
			return key, false
		}
		var seen [8]bool
		for _, typ := range field.UpsertConflictFields {
			if typ < 0 || int(typ) >= len(seen) || seen[typ] {
				return key, false
			}
			seen[typ] = true
		}
		return [8]string{string(field.Aggregate.Type), field.Aggregate.ID, field.Aggregate.InstanceID,
			field.Aggregate.ResourceOwner, field.Object.Type, field.Object.ID,
			strconv.Itoa(int(field.Object.Revision)), field.FieldName}, true
	}
	if len(operation.Remove) != len(key) {
		return key, false
	}
	for i := range key {
		switch value := operation.Remove[eventstore.FieldType(i)].(type) {
		case string:
			key[i] = value
		case eventstore.AggregateType:
			key[i] = string(value)
		case uint8:
			key[i] = strconv.Itoa(int(value))
		default:
			return key, false
		}
	}
	return key, true
}

func handleFieldBatch(ctx context.Context, tx database.Transaction, operations []*eventstore.FieldOperation) error {
	var columns [8][]string
	for i := range columns {
		columns[i] = make([]string, len(operations))
	}
	values := make([]string, len(operations))
	indexed := make([]bool, len(operations))
	for i, operation := range operations {
		key, _ := batchFieldKey(operation)
		for j, value := range key {
			columns[j][i] = value
		}
		if field := operation.Set; field != nil {
			value, err := json.Marshal(field.Value.Value)
			if err != nil {
				return zerrors.ThrowInvalidArgument(err, "V3-fcrW1", "unable to marshal field value")
			}
			values[i] = string(value)
			indexed[i] = field.Value.ShouldIndex
		}
	}
	args := make([]any, 0, 10)
	for _, column := range columns {
		args = append(args, column)
	}
	stmt := fieldBatchRemove
	if operations[0].Set != nil {
		stmt = fieldBatchSet
		args = append(args, values, indexed)
	}
	_, err := tx.Exec(ctx, stmt, args...)
	return err
}
