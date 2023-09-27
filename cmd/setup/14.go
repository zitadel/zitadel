package setup

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 14/cockroach/*.sql
	//go:embed 14/postgres/*.sql
	newEventsTable embed.FS
)

type NewEventsTable struct {
	dbClient *database.DB
}

func (mig *NewEventsTable) Execute(ctx context.Context) error {
	migrations, err := newEventsTable.ReadDir("14/" + mig.dbClient.Type())
	if err != nil {
		return err
	}
	// check if events2 table was already created from init
	if mig.dbClient.Type() == "cockroach" {
		var count int
		err := mig.dbClient.QueryRowContext(
			ctx,
			func(r *sql.Row) error {
				return r.Scan(&count)
			},
			"SELECT count(*) FROM [show tables from eventstore] WHERE table_name = 'events2'",
		)
		if err != nil {
			return err
		}
		if count == 1 {
			return nil
		}
	} else if mig.dbClient.Type() == "postgres" {
		// TODO: check for events2 table
	}
	for _, migration := range migrations {
		stmt, err := readStmt(newEventsTable, "14", mig.dbClient.Type(), migration.Name())
		if err != nil {
			return err
		}
		stmt = strings.ReplaceAll(stmt, "{{.username}}", mig.dbClient.Username())

		logging.WithFields("migration", mig.String(), "file", migration.Name()).Debug("execute statement")

		_, err = mig.dbClient.ExecContext(ctx, stmt)
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
