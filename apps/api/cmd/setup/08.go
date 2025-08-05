package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 08/08.sql
	tokenIndexes08 string
)

type AuthTokenIndexes struct {
	dbClient *database.DB
}

func (mig *AuthTokenIndexes) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, tokenIndexes08)
	return err
}

func (mig *AuthTokenIndexes) String() string {
	return "08_auth_token_indexes"
}
