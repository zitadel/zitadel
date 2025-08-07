package setup

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

var (
	//go:embed 39.sql
	deleteStaleOrgFields string
)

type DeleteStaleOrgFields struct {
	eventstore *eventstore.Eventstore
}

func (mig *DeleteStaleOrgFields) Execute(ctx context.Context, _ eventstore.Event) error {
	instances, err := mig.eventstore.InstanceIDs(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsInstanceIDs).
			OrderDesc().
			AddQuery().
			AggregateTypes("instance").
			EventTypes(instance.InstanceAddedEventType).
			Builder(),
	)
	if err != nil {
		return err
	}
	for i, instance := range instances {
		logging.WithFields("instance_id", instance, "migration", mig.String(), "progress", fmt.Sprintf("%d/%d", i+1, len(instances))).Info("execute delete query")
		if _, err := mig.eventstore.Client().ExecContext(ctx, deleteStaleOrgFields, instance); err != nil {
			return err
		}
	}
	return nil
}

func (*DeleteStaleOrgFields) Check(map[string]any) bool {
	return true
}

func (*DeleteStaleOrgFields) String() string {
	return "repeatable_delete_stale_org_fields"
}
