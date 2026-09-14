package setup

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 79/*.sql
	backfillUniqueConstraintOwnersFS embed.FS
)

type BackfillUniqueConstraintOwners struct {
	dbClient *database.DB
}

func (mig *BackfillUniqueConstraintOwners) Execute(ctx context.Context, _ eventstore.Event) error {
	statements, err := readStatements(backfillUniqueConstraintOwnersFS, "79")
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.Info(ctx, "backfill unique constraint owners", "file", stmt.file, "migration", mig.String())
		if _, err := mig.dbClient.ExecContext(ctx, stmt.query); err != nil {
			if isUndefinedTable(err) {
				logging.Info(ctx, "skip unique constraint owners backfill, relation missing", "file", stmt.file, "migration", mig.String())
				continue
			}
			return err
		}
	}

	var unmatched int64
	err = mig.dbClient.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&unmatched)
	}, `SELECT COUNT(*) FROM eventstore.unique_constraints WHERE owners = '{}'`)
	if err != nil {
		return err
	}
	logging.Info(ctx, "unique constraint owners backfill complete", "unmatched", unmatched, "migration", mig.String())
	return nil
}

func (mig *BackfillUniqueConstraintOwners) String() string {
	return "79_backfill_unique_constraint_owners"
}

func isUndefinedTable(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "42P01"
}
