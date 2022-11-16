package setup

import (
	"context"
	"database/sql"
	_ "embed"
)

var (
	//go:embed 05/adminapi.sql
	createAdminViews05 string
	//go:embed 05/auth.sql
	createAuthViews05 string
)

type ProjectionTable05 struct {
	dbClient *sql.DB
}

func (mig *ProjectionTable05) Execute(ctx context.Context) error {
	stmt := createAdminViews05 + createAuthViews05
	_, err := mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *ProjectionTable05) String() string {
	return "05_tables"
}
