package setup

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
)

const (
	users14TableExistsQuery  = "SELECT exists(SELECT 1 FROM information_schema.tables WHERE table_schema = 'projections' AND table_name = 'users14')"
	users14InvalidIndexQuery = `SELECT EXISTS (
	SELECT 1
	FROM pg_index i
	JOIN pg_class c ON c.oid = i.indexrelid
	JOIN pg_namespace n ON n.oid = c.relnamespace
	WHERE n.nspname = 'projections'
		AND c.relname = $1
		AND NOT i.indisvalid
)`
	dropInvalidIndexConcurrently = "DROP INDEX CONCURRENTLY IF EXISTS projections."
)

// executeConcurrentIndexSQL skips when the table is missing, drops invalid leftovers
// of indexNames (failed CREATE INDEX CONCURRENTLY), then runs folder statements.
func executeConcurrentIndexSQL(ctx context.Context, db *database.DB, migrationName, tableExistsQuery string, indexNames []string, fs embed.FS, folder string) error {
	var exists bool
	err := db.QueryRowContext(ctx, func(r *sql.Row) error {
		return r.Scan(&exists)
	}, tableExistsQuery)
	if err != nil || !exists {
		return err
	}

	if err := dropInvalidProjectionIndexes(ctx, db, migrationName, indexNames); err != nil {
		return err
	}

	statements, err := readStatements(fs, folder)
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.Info(ctx, "execute statement", "file", stmt.file, "migration", migrationName)
		if _, err := db.ExecContext(ctx, stmt.query); err != nil {
			return fmt.Errorf("%s %s: %w", migrationName, stmt.file, err)
		}
	}
	return nil
}

func dropInvalidProjectionIndexes(ctx context.Context, db *database.DB, migrationName string, names []string) error {
	for _, name := range names {
		var invalid bool
		err := db.QueryRowContext(ctx, func(r *sql.Row) error {
			return r.Scan(&invalid)
		}, users14InvalidIndexQuery, name)
		if err != nil {
			return fmt.Errorf("%s check invalid index %s: %w", migrationName, name, err)
		}
		if !invalid {
			continue
		}
		drop := dropInvalidIndexConcurrently + name
		logging.Info(ctx, "drop invalid leftover index", "index", name, "migration", migrationName)
		if _, err := db.ExecContext(ctx, drop); err != nil {
			return fmt.Errorf("%s drop invalid index %s: %w", migrationName, name, err)
		}
	}
	return nil
}
