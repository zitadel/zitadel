package setup

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 76/*.sql
	users14LoginEqualityIndexes embed.FS
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
	users14UsernameLowerIdx      = "users14_username_lower_idx"
	users14HumansPhoneLowerIdx   = "users14_humans_phone_lower_idx"
	dropInvalidIndexConcurrently = "DROP INDEX CONCURRENTLY IF EXISTS projections."
)

var users14LoginEqualityIndexNames = []string{
	users14UsernameLowerIdx,
	users14HumansPhoneLowerIdx,
}

type Users14LoginEqualityIndexes struct {
	dbClient *database.DB
}

func (mig *Users14LoginEqualityIndexes) Execute(ctx context.Context, _ eventstore.Event) error {
	var exists bool
	err := mig.dbClient.QueryRowContext(ctx, func(r *sql.Row) error {
		return r.Scan(&exists)
	}, users14TableExistsQuery)
	if err != nil || !exists {
		return err
	}

	if err := mig.dropInvalidIndexes(ctx); err != nil {
		return err
	}

	statements, err := readStatements(users14LoginEqualityIndexes, "76")
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.Info(ctx, "execute statement", "file", stmt.file, "migration", mig.String())
		if _, err := mig.dbClient.ExecContext(ctx, stmt.query); err != nil {
			return fmt.Errorf("%s %s: %w", mig.String(), stmt.file, err)
		}
	}
	return nil
}

func (mig *Users14LoginEqualityIndexes) dropInvalidIndexes(ctx context.Context) error {
	for _, name := range users14LoginEqualityIndexNames {
		var invalid bool
		err := mig.dbClient.QueryRowContext(ctx, func(r *sql.Row) error {
			return r.Scan(&invalid)
		}, users14InvalidIndexQuery, name)
		if err != nil {
			return fmt.Errorf("%s check invalid index %s: %w", mig.String(), name, err)
		}
		if !invalid {
			continue
		}
		drop := dropInvalidIndexConcurrently + name
		logging.Info(ctx, "drop invalid leftover index", "index", name, "migration", mig.String())
		if _, err := mig.dbClient.ExecContext(ctx, drop); err != nil {
			return fmt.Errorf("%s drop invalid index %s: %w", mig.String(), name, err)
		}
	}
	return nil
}

func (mig *Users14LoginEqualityIndexes) String() string {
	return "76_users14_login_equality_indexes"
}
