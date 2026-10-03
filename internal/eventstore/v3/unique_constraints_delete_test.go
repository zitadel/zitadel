package eventstore

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var _ database.Transaction = (*recordingTx)(nil)

type recordingTx struct {
	stmts []string
}

func (r *recordingTx) Exec(_ context.Context, stmt string, _ ...any) (int64, error) {
	r.stmts = append(r.stmts, stmt)
	return 1, nil
}

func (r *recordingTx) Query(context.Context, string, ...any) (database.Rows, error) {
	return nil, nil
}

func (r *recordingTx) QueryRow(context.Context, string, ...any) database.Row {
	return nil
}

func (r *recordingTx) Begin(context.Context) (database.Transaction, error) {
	return r, nil
}

func (r *recordingTx) Commit(context.Context) error { return nil }

func (r *recordingTx) Rollback(context.Context) error { return nil }

func (r *recordingTx) End(_ context.Context, err error) error { return err }

func TestHandleUniqueConstraintsSplitsKeyedAndOwnerDeletes(t *testing.T) {
	tx := &recordingTx{}
	cmd := &mockCommand{
		aggregate: mockAggregate("project-1"),
		constraints: []*eventstore.UniqueConstraint{
			eventstore.NewRemoveUniqueConstraint("project_names", "docsro"),
			eventstore.NewRemoveUniqueConstraintsByOwner(eventstore.UniqueConstraintOwnerProject, "project-1"),
		},
	}

	err := handleUniqueConstraints(t.Context(), tx, []eventstore.Command{cmd})
	require.NoError(t, err)
	require.Len(t, tx.stmts, 2)

	assert.Contains(t, tx.stmts[0], "unique_type")
	assert.NotContains(t, tx.stmts[0], "owners @>")
	assert.False(t, strings.Contains(tx.stmts[0], " OR "))

	assert.Contains(t, tx.stmts[1], "owners @>")
	assert.NotContains(t, tx.stmts[1], "unique_type")
	assert.False(t, strings.Contains(tx.stmts[1], " OR "))
}
