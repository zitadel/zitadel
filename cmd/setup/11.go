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
	step10   *CorrectCreationDate
	dbClient *database.DB
}

func (mig *AddEventCreatedAt) Execute(ctx context.Context) error {
	// execute step 10 again because events created after the first execution of step 10
	// could still have the wrong ordering of sequences and creation date
	if err := mig.step10.Execute(ctx); err != nil {
		return err
	}
	_, err := mig.dbClient.ExecContext(ctx, addEventCreatedAt)
	return err
}

func (mig *AddEventCreatedAt) String() string {
	return "11_event_created_at"
}
