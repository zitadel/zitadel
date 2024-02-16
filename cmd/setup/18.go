package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 18.sql
	addLowerFieldsToLoginNames string
)

type AddLowerFieldsToLoginNames struct {
	dbClient *database.DB
}

func (mig *AddLowerFieldsToLoginNames) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addLowerFieldsToLoginNames)
	return err
}

func (mig *AddLowerFieldsToLoginNames) String() string {
	return "18_add_lower_fields_to_login_names"
}
