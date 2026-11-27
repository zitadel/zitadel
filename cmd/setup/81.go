package setup

import (
	"context"
	"embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 81/*.sql
	userGrants5InstanceUserIDIndex embed.FS
)

const (
	userGrants5InstanceUserIDIdx = "user_grants5_instance_user_id_idx"
	userGrants5TableExistsQuery  = "SELECT exists(SELECT 1 FROM information_schema.tables WHERE table_schema = 'projections' AND table_name = 'user_grants5')"
)

type UserGrants5InstanceUserIDIndex struct {
	dbClient *database.DB
}

func (mig *UserGrants5InstanceUserIDIndex) Execute(ctx context.Context, _ eventstore.Event) error {
	return executeConcurrentIndexSQL(ctx, mig.dbClient, mig.String(), userGrants5TableExistsQuery, []string{userGrants5InstanceUserIDIdx}, userGrants5InstanceUserIDIndex, "81")
}

func (mig *UserGrants5InstanceUserIDIndex) String() string {
	return "81_user_grants5_instance_user_id_idx"
}
