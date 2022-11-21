package setup

import (
	"context"
	"database/sql"
	_ "embed"
)

var (
	//go:embed 06/adminapi.sql
	createAdminViews05 string
	//go:embed 06/auth.sql
	createAuthViews05 string
)

type OwnerRemoveColumns struct {
	dbClient *sql.DB
}

func (mig *OwnerRemoveColumns) Execute(ctx context.Context) error {
	stmt := createAdminViews05 + createAuthViews05
	_, err := mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *OwnerRemoveColumns) String() string {
	return "06_resource_owner_columns"
}
