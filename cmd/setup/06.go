package setup

import (
	"context"
	"database/sql"
	_ "embed"
)

var (
	//go:embed 06/adminapi.sql
	createAdminViews06 string
	//go:embed 06/auth.sql
	createAuthViews06 string
)

type OwnerRemoveColumns struct {
	dbClient *sql.DB
}

func (mig *OwnerRemoveColumns) Execute(ctx context.Context) error {
	stmt := createAdminViews06 + createAuthViews06
	_, err := mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *OwnerRemoveColumns) String() string {
	return "06_resource_owner_columns"
}

func (mig *OwnerRemoveColumns) ShouldSkip() bool {
	return false
}
