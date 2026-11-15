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

//go:embed field_batch_remove_aggregate.sql
var fieldBatchRemoveAggregate string

type fieldBatchKind int

const (
	fieldBatchKindNone fieldBatchKind = iota
	fieldBatchKindSet
	fieldBatchKindRemoveField
	// Removes by the four aggregate columns, like [eventstore.RemoveSearchFieldsByAggregate].
	fieldBatchKindRemoveByAggregate
)

const aggregateKeyLen = 4

// Only full-identity operations without value uniqueness, and removes by full
// identity or by aggregate, can be reordered within a statement.
func batchFieldKey(operation *eventstore.FieldOperation) (key [8]string, kind fieldBatchKind) {
	if field := operation.Set; field != nil {
		if field.Value.MustBeUnique || len(field.UpsertConflictFields) != len(key) {
			return key, fieldBatchKindNone
		}
		var seen [8]bool
		for _, typ := range field.UpsertConflictFields {
			if typ < 0 || int(typ) >= len(seen) || seen[typ] {
				return key, fieldBatchKindNone
			}
			seen[typ] = true
		}
		return [8]string{string(field.Aggregate.Type), field.Aggregate.ID, field.Aggregate.InstanceID,
			field.Aggregate.ResourceOwner, field.Object.Type, field.Object.ID,
			strconv.Itoa(int(field.Object.Revision)), field.FieldName}, fieldBatchKindSet
	}
	kind = fieldBatchKindRemoveField
	if len(operation.Remove) == aggregateKeyLen {
		kind = fieldBatchKindRemoveByAggregate
	} else if len(operation.Remove) != len(key) {
		return key, fieldBatchKindNone
	}
	for i := range len(operation.Remove) {
		switch value := operation.Remove[eventstore.FieldType(i)].(type) {
		case string:
			key[i] = value
		case eventstore.AggregateType:
			key[i] = string(value)
		case uint8:
			key[i] = strconv.Itoa(int(value))
		default:
			return key, fieldBatchKindNone
		}
	}
	return key, kind
}

func handleFieldBatch(ctx context.Context, tx database.Transaction, kind fieldBatchKind, operations []*eventstore.FieldOperation) error {
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
	switch kind {
	case fieldBatchKindSet:
		stmt = fieldBatchSet
		args = append(args, values, indexed)
	case fieldBatchKindRemoveByAggregate:
		stmt = fieldBatchRemoveAggregate
		args = args[:aggregateKeyLen]
	case fieldBatchKindNone, fieldBatchKindRemoveField:
	}
	_, err := tx.Exec(ctx, stmt, args...)
	return err
}
