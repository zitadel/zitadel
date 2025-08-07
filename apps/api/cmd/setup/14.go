package setup

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 14/*.sql
	newEventsTable embed.FS
)

type NewEventsTable struct {
	dbClient *database.DB
}

func (mig *NewEventsTable) Execute(ctx context.Context, _ eventstore.Event) error {
	// if events already exists events2 is created during a setup job
	var count int
	err := mig.dbClient.QueryRowContext(ctx,
		func(row *sql.Row) error {
			if err := row.Scan(&count); err != nil {
				return err
			}
			return row.Err()
		},
		"SELECT count(*) FROM information_schema.tables WHERE table_schema = 'eventstore' AND table_name like 'events2'",
	)
	if err != nil || count == 1 {
		return err
	}

	statements, err := readStatements(newEventsTable, "14")
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		stmt.query = strings.ReplaceAll(stmt.query, "{{.username}}", mig.dbClient.Username())
		logging.WithFields("file", stmt.file, "migration", mig.String()).Info("execute statement")
		_, err = mig.dbClient.ExecContext(ctx, stmt.query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mig *NewEventsTable) String() string {
	return "14_events_push"
}

func (mig *NewEventsTable) ContinueOnErr(err error) bool {
	pgErr := new(pgconn.PgError)
	if errors.As(err, &pgErr) {
		return pgErr.Code == "42P01"
	}
	return false
}
