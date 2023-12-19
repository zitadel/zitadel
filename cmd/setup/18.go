package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 18.sql
	addLowerFieldsToLoginNames string
)

type AddLowerFieldsToLoginNames struct {
	dbClient *database.DB
}

func (mig *AddLowerFieldsToLoginNames) Execute(ctx context.Context) error {
	_, err := mig.dbClient.ExecContext(ctx, addLowerFieldsToLoginNames)
	return err
}

func (mig *AddLowerFieldsToLoginNames) String() string {
	return "18_add_lower_fields_to_login_names"
}
