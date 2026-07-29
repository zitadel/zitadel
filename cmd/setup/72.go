package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 72.sql
	addColumnsToLoginNamesView string
)

type AddColumnsToLoginNamesView struct {
	dbClient *database.DB
}

func (mig *AddColumnsToLoginNamesView) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addColumnsToLoginNamesView)
	return err
}

func (mig *AddColumnsToLoginNamesView) String() string {
	return "72_add_columns_to_login_names_view"
}
