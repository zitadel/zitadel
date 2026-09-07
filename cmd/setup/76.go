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

type Users14LoginEqualityIndexes struct {
	dbClient *database.DB
}

func (mig *Users14LoginEqualityIndexes) Execute(ctx context.Context, _ eventstore.Event) error {
	var exists bool
	err := mig.dbClient.QueryRowContext(ctx, func(r *sql.Row) error {
		return r.Scan(&exists)
	}, "SELECT exists(SELECT 1 FROM information_schema.tables WHERE table_schema = 'projections' AND table_name = 'users14')")
	if err != nil || !exists {
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

func (mig *Users14LoginEqualityIndexes) String() string {
	return "76_users14_login_equality_indexes"
}
