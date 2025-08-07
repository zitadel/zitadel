package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 26.sql
	authUsers3 string
)

type AuthUsers3 struct {
	dbClient *database.DB
}

func (mig *AuthUsers3) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, authUsers3)
	return err
}

func (mig *AuthUsers3) String() string {
	return "26_auth_users3"
}
