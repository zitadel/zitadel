package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 70.sql
	signalsSchema string
)

type SignalsSchema struct {
	dbClient *database.DB
}

func (mig *SignalsSchema) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, signalsSchema)
	return err
}

func (mig *SignalsSchema) String() string {
	return "70_signals_schema"
}
