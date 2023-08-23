package setup

import (
	"context"
	"embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 13/cockroach/13.sql
	//go:embed 13/postgres/13.sql
	changeEvents embed.FS
)

type ChangeEvents struct {
	dbClient *database.DB
}

func (mig *ChangeEvents) Execute(ctx context.Context) error {
	stmt, err := readStmt(changeEvents, "13", mig.dbClient.Type(), "13.sql")
	if err != nil {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *ChangeEvents) String() string {
	return "13_events_push"
}
