package setup

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 58.sql
	replaceLoginNames3View string
)

type ReplaceLoginNames3View struct {
	dbClient *database.DB
}

func (mig *ReplaceLoginNames3View) Execute(ctx context.Context, _ eventstore.Event) error {
	var exists bool
	err := mig.dbClient.QueryRow(func(r *sql.Row) error {
		return r.Scan(&exists)
	}, "SELECT exists(SELECT 1 from information_schema.views WHERE table_schema = 'projections' AND table_name = 'login_names3')")

	if err != nil || !exists {
		return err
	}

	_, err = mig.dbClient.ExecContext(ctx, replaceLoginNames3View)
	return err
}

func (mig *ReplaceLoginNames3View) String() string {
	return "58_replace_login_names3_view"
}
