package setup

import (
	"context"
	"embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 12/cockroach/12.sql
	//go:embed 12/postgres/12.sql
	changeEvents embed.FS
)

type ChangeEvents struct {
	dbClient *database.DB
}

func (mig *ChangeEvents) Execute(ctx context.Context) error {
	stmt, err := readStmt(changeEvents, "12", mig.dbClient.Type(), "12.sql")
	if err != nil {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *ChangeEvents) String() string {
	return "12_events_push"
}
