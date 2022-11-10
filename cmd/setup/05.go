package setup

import (
	"context"
	"database/sql"
	"embed"
)

var (
	//go:embed 05/cockroach/columns.sql
	//go:embed 05/postgres/columns.sql
	readOwnerRemovedColumnsStmts embed.FS
)

type OwnerRemovedColumns struct {
	dbClient *sql.DB
	dbType   string
}

func (mig *OwnerRemovedColumns) Execute(ctx context.Context) error {
	readOwnerRemovedColumnsStmts, err := readOwnerRemovedColumnsStmt(mig.dbType)
	if err != nil {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, readOwnerRemovedColumnsStmts)
	return err
}

func (mig *OwnerRemovedColumns) String() string {
	return "05_owner_removed_columns"
}

func readOwnerRemovedColumnsStmt(typ string) (string, error) {
	readOwnerRemovedColumnsStmts, err := readOwnerRemovedColumnsStmts.ReadFile("05/" + typ + "/columns.sql")
	return string(readOwnerRemovedColumnsStmts), err
}
