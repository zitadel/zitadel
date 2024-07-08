package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 31.sql
	addAggregateIndexToFields string
)

type AddAggregateIndexToFields struct {
	dbClient *database.DB
}

func (mig *AddAggregateIndexToFields) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addAggregateIndexToFields)
	return err
}

func (mig *AddAggregateIndexToFields) String() string {
	return "31_add_aggregate_index_to_fields"
}
