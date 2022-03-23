package setup

import (
	"context"
	"database/sql"
	_ "embed"
)

var (
	//go:embed 01_sql/adminapi.sql
	createAdminViews string
	//go:embed 01_sql/auth.sql
	createAuthViews string
	//go:embed 01_sql/authz.sql
	createAuthzViews string
	//go:embed 01_sql/notification.sql
	createNotificationViews string
	//go:embed 01_sql/projections.sql
	createProjections string
)

type ProjectionTable struct {
	dbClient *sql.DB
}

func (mig *ProjectionTable) Execute(ctx context.Context) error {
	stmt := createAdminViews + createAuthViews + createAuthzViews + createNotificationViews + createProjections
	_, err := mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *ProjectionTable) String() string {
	return "01_tables"
}
