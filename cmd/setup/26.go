package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 26.sql
	addTokenProjectID string
)

type AddProjectIDToAuthTokens struct {
	dbClient *database.DB
}

func (mig *AddProjectIDToAuthTokens) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addTokenProjectID)
	return err
}

func (mig *AddProjectIDToAuthTokens) String() string {
	return "26_add_project_id_col_to_auth_tokens"
}
