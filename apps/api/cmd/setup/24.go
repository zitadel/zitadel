package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 24.sql
	addTokenActor string
)

type AddActorToAuthTokens struct {
	dbClient *database.DB
}

func (mig *AddActorToAuthTokens) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addTokenActor)
	return err
}

func (mig *AddActorToAuthTokens) String() string {
	return "24_add_actor_col_to_auth_tokens"
}
