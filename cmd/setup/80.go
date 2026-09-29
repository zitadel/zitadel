package setup

import (
	"context"
	"embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 80/*.sql
	users14InstanceResourceOwnerIndex embed.FS
)

const users14InstanceResourceOwnerIdx = "users14_instance_resource_owner_idx"

type Users14InstanceResourceOwnerIndex struct {
	dbClient *database.DB
}

func (mig *Users14InstanceResourceOwnerIndex) Execute(ctx context.Context, _ eventstore.Event) error {
	return executeConcurrentIndexSQL(ctx, mig.dbClient, mig.String(), users14TableExistsQuery, []string{users14InstanceResourceOwnerIdx}, users14InstanceResourceOwnerIndex, "80")
}

func (mig *Users14InstanceResourceOwnerIndex) String() string {
	return "80_users14_instance_resource_owner_idx"
}
