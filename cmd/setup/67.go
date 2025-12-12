package setup

import (
	"context"
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
	_, err := mig.dbClient.ExecContext(ctx, syncMemberRoleFields)
	return err
}

func (mig *SyncMemberRoleFields) String() string {
	return "67_sync_member_role_fields"
}
