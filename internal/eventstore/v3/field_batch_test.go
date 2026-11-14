package eventstore

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	sqldb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/internal/database/postgres"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type countingFieldTx struct {
	database.Transaction
	statements int
}

func (tx *countingFieldTx) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	tx.statements++
	return tx.Transaction.Exec(ctx, query, args...)
}

func TestFieldBatch(t *testing.T) {
	config, stop := postgres.StartEmbedded()
	defer stop()
	db, err := sql.Open("pgx", config.GetConnectionURL())
	require.NoError(t, err)
	defer db.Close()
	ctx := t.Context()
	_, err = db.ExecContext(ctx, "CREATE SCHEMA eventstore")
	require.NoError(t, err)
	ddl, err := os.ReadFile("../../../cmd/setup/28.sql")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, string(ddl))
	require.NoError(t, err)

	aggregate := &eventstore.Aggregate{InstanceID: "instance", ResourceOwner: "org", Type: "group", ID: "group"}
	field := func(id string, value any) *eventstore.FieldOperation {
		return eventstore.SetField(aggregate, eventstore.Object{Type: "group", ID: id, Revision: 1}, "member",
			&eventstore.Value{Value: value, ShouldIndex: true},
			eventstore.FieldTypeInstanceID, eventstore.FieldTypeResourceOwner, eventstore.FieldTypeAggregateType,
			eventstore.FieldTypeAggregateID, eventstore.FieldTypeObjectType, eventstore.FieldTypeObjectID,
			eventstore.FieldTypeObjectRevision, eventstore.FieldTypeFieldName)
	}
	remove := func(id string) *eventstore.FieldOperation {
		return eventstore.RemoveSearchFieldsByAggregateAndObjectAndField(aggregate, eventstore.Object{Type: "group", ID: id, Revision: 1}, "member")
	}
	run := func(operations []*eventstore.FieldOperation, statements int, wantErr bool) {
		t.Helper()
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		counted := &countingFieldTx{Transaction: sqldb.SQLTx(tx)}
		err = handleFieldOperations(ctx, counted, operations)
		if wantErr {
			require.Error(t, err)
			require.NoError(t, tx.Rollback())
		} else {
			require.NoError(t, err)
			require.NoError(t, tx.Commit())
		}
		require.Equal(t, statements, counted.statements)
	}
	count := func(want int) {
		t.Helper()
		var got int
		require.NoError(t, db.QueryRowContext(ctx, "SELECT count(*) FROM eventstore.fields").Scan(&got))
		require.Equal(t, want, got)
	}
	value := func(id, want string) {
		t.Helper()
		var got string
		require.NoError(t, db.QueryRowContext(ctx, "SELECT text_value FROM eventstore.fields WHERE instance_id = 'instance' AND object_id = $1", id).Scan(&got))
		require.Equal(t, want, got)
	}

	// Large enough to exceed the parameter limit of a VALUES list with eleven columns.
	const size = 7000
	adds, removes := make([]*eventstore.FieldOperation, size), make([]*eventstore.FieldOperation, size)
	for i := range adds {
		id := fmt.Sprintf("member-%d", i)
		adds[i], removes[i] = field(id, id), remove(id)
	}
	run(adds, 1, false)
	count(size)
	var originalID string
	require.NoError(t, db.QueryRowContext(ctx, "SELECT id FROM eventstore.fields WHERE object_id = 'member-0'").Scan(&originalID))
	adds[0].Set.Value.Value = `changed,"{}\雪`
	run(adds, 1, false)
	count(size)
	value("member-0", `changed,"{}\雪`)
	var updatedID string
	require.NoError(t, db.QueryRowContext(ctx, "SELECT id FROM eventstore.fields WHERE object_id = 'member-0'").Scan(&updatedID))
	require.Equal(t, originalID, updatedID)
	run(removes, 1, false)
	count(0)

	// Same keys in different instances must remain distinct, even in one batch.
	other := field("a", "other")
	otherAggregate := *aggregate
	otherAggregate.InstanceID = "other"
	other.Set.Aggregate = &otherAggregate
	run([]*eventstore.FieldOperation{field("a", "first"), other, field("a", "last"), field("b", "b")}, 2, false)
	count(3)
	value("a", "last")
	run([]*eventstore.FieldOperation{remove("a"), remove("b"), field("a", "after remove"), field("b", "b")}, 2, false)
	value("a", "after remove")
	run([]*eventstore.FieldOperation{field("a", "before broad delete"), field("b", "b"),
		eventstore.RemoveSearchFieldsByAggregate(aggregate), field("a", "after broad delete"), field("b", "b")}, 3, false)
	count(3)
	value("a", "after broad delete")

	// Removes by aggregate batch together, but not with removes by field.
	otherGroup := *aggregate
	otherGroup.ID = "other-group"
	run([]*eventstore.FieldOperation{field("a", "a"), field("b", "b")}, 1, false)
	otherField := field("x", "x")
	otherField.Set.Aggregate = &otherGroup
	run([]*eventstore.FieldOperation{otherField}, 1, false)
	count(4)
	run([]*eventstore.FieldOperation{eventstore.RemoveSearchFieldsByAggregate(aggregate), eventstore.RemoveSearchFieldsByAggregate(&otherGroup),
		remove("never")}, 2, false)
	count(1)
	run([]*eventstore.FieldOperation{field("a", "after broad delete"), field("b", "b")}, 1, false)

	// A failed later statement must roll back the preceding batch as well.
	run([]*eventstore.FieldOperation{remove("a"), remove("b"), field("bad", []string{"not primitive"}), field("c", "c")}, 2, true)
	count(3)
	value("a", "after broad delete")
	// Do not coalesce repeated keys and accidentally skip an invalid intermediate write.
	run([]*eventstore.FieldOperation{field("a", []string{"invalid"}), field("a", "valid")}, 1, true)
	value("a", "after broad delete")
	badJSON := field("bad", make(chan int))
	run([]*eventstore.FieldOperation{field("c", "c"), badJSON}, 0, true)
	count(3)
}

func TestBatchFieldKeyFallback(t *testing.T) {
	for _, operation := range []*eventstore.FieldOperation{
		{Set: &eventstore.Field{Value: eventstore.Value{MustBeUnique: true}}},
		{Set: &eventstore.Field{UpsertConflictFields: []eventstore.FieldType{eventstore.FieldTypeObjectID}}},
		{Set: &eventstore.Field{UpsertConflictFields: make([]eventstore.FieldType, 8)}},
		{Remove: map[eventstore.FieldType]any{eventstore.FieldTypeInstanceID: "i", eventstore.FieldTypeObjectType: "o",
			eventstore.FieldTypeObjectID: "id", eventstore.FieldTypeFieldName: "f"}},
	} {
		_, kind := batchFieldKey(operation)
		require.Equal(t, fieldBatchKindNone, kind)
	}
	_, kind := batchFieldKey(eventstore.RemoveSearchFieldsByAggregate(&eventstore.Aggregate{}))
	require.Equal(t, fieldBatchKindRemoveByAggregate, kind)
}
