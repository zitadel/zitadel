package setup

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 58/*.sql
	replaceLoginNames3View embed.FS
)

type ReplaceLoginNames3View struct {
	dbClient *database.DB
}

func (mig *ReplaceLoginNames3View) Execute(ctx context.Context, _ eventstore.Event) error {
	var exists bool
	err := mig.dbClient.QueryRowContext(ctx, func(r *sql.Row) error {
		return r.Scan(&exists)
	}, "SELECT exists(SELECT 1 from information_schema.views WHERE table_schema = 'projections' AND table_name = 'login_names3')")

	if err != nil || !exists {
		return err
	}

	statements, err := readStatements(replaceLoginNames3View, "58")
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.WithFields("file", stmt.file, "migration", mig.String()).Info("execute statement")
		if _, err := mig.dbClient.ExecContext(ctx, stmt.query); err != nil {
			return fmt.Errorf("%s %s: %w", mig.String(), stmt.file, err)
		}
	}
	return nil
}

func (mig *ReplaceLoginNames3View) String() string {
	return "58_replace_login_names3_view"
}
