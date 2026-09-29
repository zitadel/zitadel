package eventstore

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type recordingTx struct {
	stmts []string
}

func (r *recordingTx) ExecContext(_ context.Context, query string, _ ...any) (sql.Result, error) {
	r.stmts = append(r.stmts, query)
	return noopResult{}, nil
}

func (r *recordingTx) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	return nil, errors.New("unexpected query")
}

func (r *recordingTx) Commit() error { return nil }

func (r *recordingTx) Rollback() error { return nil }

type noopResult struct{}

func (noopResult) LastInsertId() (int64, error) { return 0, nil }
func (noopResult) RowsAffected() (int64, error) { return 1, nil }

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
