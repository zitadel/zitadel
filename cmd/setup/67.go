package setup

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 67.sql
	syncMemberRoleFields string
)

type SyncMemberRoleFields struct {
	dbClient *database.DB
}

func (mig *SyncMemberRoleFields) Execute(ctx context.Context, _ eventstore.Event) error {
	var exists bool
	err := mig.dbClient.QueryRowContext(
		ctx,
		func(row *sql.Row) error {
			return row.Scan(&exists)
		},
		"SELECT EXISTS(SELECT FROM pg_catalog.pg_tables WHERE schemaname = 'projections' and tablename = 'instance_members4')")
	if err != nil || !exists {
		return err
	}
	_, err = mig.dbClient.ExecContext(ctx, syncMemberRoleFields)
	return err
}

func (mig *SyncMemberRoleFields) String() string {
	return "67_sync_member_role_fields"
}
