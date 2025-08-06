package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 61.sql
	addScopesToSession string
)

type Sessions8AddScopes struct {
	dbClient *database.DB
}

func (mig *Sessions8AddScopes) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addScopesToSession)
	return err
}

func (mig *Sessions8AddScopes) String() string {
	return "61_sessions8_add_scopes"
}
