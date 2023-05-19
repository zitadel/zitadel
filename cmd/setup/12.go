package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 12.sql
	changeEvents string
)

type ChangeEvents struct {
	dbClient *database.DB
}

func (mig *ChangeEvents) Execute(ctx context.Context) error {
	_, err := mig.dbClient.ExecContext(ctx, changeEvents)
	return err
}

func (mig *ChangeEvents) String() string {
	return "12_events_push"
}
