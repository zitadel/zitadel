package setup

import (
	"context"
	"embed"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 08/cockroach/08.sql
	//go:embed 08/postgres/08.sql
	tokenIndexes08 embed.FS
)

type AuthTokenIndexes struct {
	dbClient *database.DB
}

func (mig *AuthTokenIndexes) Execute(ctx context.Context) error {
	stmt, err := readStmt(tokenIndexes08, "08", mig.dbClient.Type(), "08.sql")
	if err != nil {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *AuthTokenIndexes) String() string {
	return "08_auth_token_indexes"
}
