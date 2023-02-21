package setup

import (
	"context"
	"database/sql"
	_ "embed"
)

var (
	//go:embed 08.sql
	tokenIndexes07 string
)

type AuthTokenIndexes struct {
	dbClient *sql.DB
}

func (mig *AuthTokenIndexes) Execute(ctx context.Context) error {
	_, err := mig.dbClient.ExecContext(ctx, tokenIndexes07)
	return err
}

func (mig *AuthTokenIndexes) String() string {
	return "08_auth_token_indexes"
}
