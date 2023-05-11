package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 11.sql
	addEventCreatedAt string
)

type AddEventCreatedAt struct {
	dbClient *database.DB
}

func (mig *AddEventCreatedAt) Execute(ctx context.Context) error {
	_, err := mig.dbClient.ExecContext(ctx, addEventCreatedAt)
	return err
}

func (mig *AddEventCreatedAt) String() string {
	return "11_event_created_at"
}
