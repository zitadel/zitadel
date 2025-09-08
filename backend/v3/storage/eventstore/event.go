package eventstore

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	legacy_db "github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type Event struct {
	AggregateType string `json:"aggregateType"`
	AggregateID   string `json:"aggregateId"`
	Type          string `json:"type"`
	Payload       any    `json:"payload,omitempty"`
}

func Publish(ctx context.Context, events []*Event, db database.Executor) error {
	for _, event := range events {
		_, err := db.Exec(ctx, `INSERT INTO events (aggregate_type, aggregate_id) VALUES ($1, $2)`, event.AggregateType, event.AggregateID)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteLegacyEvents(ctx context.Context, es eventstore.Pusher, client database.QueryExecutor, commands ...eventstore.Command) error {
	_, err := es.Push(ctx, LegacyContextQueryExecutorAdapter{client}, commands...)
	return err
}

type LegacyContextQueryExecutorAdapter struct{ database.QueryExecutor }

// ExecContext implements database.ContextQueryExecuter.
func (l LegacyContextQueryExecutorAdapter) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	affected, err := l.QueryExecutor.Exec(ctx, query, args...)
	return &sqlResult{rowsAffected: affected}, err
}

// QueryContext implements database.ContextQueryExecuter.
func (l LegacyContextQueryExecutorAdapter) QueryContext(ctx context.Context, query string, args ...any) (legacy_db.Rows, error) {
	rows, err := l.QueryExecutor.Query(ctx, query, args...)
	return rows, err
}

var _ legacy_db.ContextQueryExecuter = (*LegacyContextQueryExecutorAdapter)(nil)

type sqlResult struct {
	rowsAffected int64
}

// LastInsertId implements [sql.Result].
// Its never used in Zitadel so it always returns 0, nil.
func (s *sqlResult) LastInsertId() (int64, error) {
	return 0, nil
}

// RowsAffected implements [sql.Result].
func (s *sqlResult) RowsAffected() (int64, error) {
	return s.rowsAffected, nil
}

var _ sql.Result = (*sqlResult)(nil)
