package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 28.sql
	addFieldTable string
)

type AddFieldsTable struct {
	dbClient *database.DB
}

func (mig *AddFieldsTable) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addFieldTable)
	return err
}

func (mig *AddFieldsTable) String() string {
	return "28_add_search_table"
}
