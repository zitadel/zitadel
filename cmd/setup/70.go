package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 70.sql
	addWrittenByV3Column string
)

type AddWrittenByV3Column struct {
	dbClient *database.DB
}

func (mig *AddWrittenByV3Column) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addWrittenByV3Column)
	return err
}

func (mig *AddWrittenByV3Column) String() string {
	return "70_add_written_by_v3_column"
}
