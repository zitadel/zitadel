package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 54.sql
	instancePositionIndex string
)

type InstancePositionIndex struct {
	dbClient *database.DB
}

func (mig *InstancePositionIndex) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, instancePositionIndex)
	return err
}

func (mig *InstancePositionIndex) String() string {
	return "54_instance_position_index_remove_again"
}
