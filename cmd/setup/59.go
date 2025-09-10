package setup

import (
	"context"
	"fmt"

	"github.com/muhlemmer/gu"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type SetupWebkeys struct {
	eventstore *eventstore.Eventstore
	commands   *command.Commands
}

func (mig *SetupWebkeys) Execute(ctx context.Context, _ eventstore.Event) error {
	instances, err := mig.eventstore.InstanceIDs(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsInstanceIDs).
			OrderDesc().
			AddQuery().
			AggregateTypes(instance.AggregateType).
			EventTypes(instance.InstanceAddedEventType).
			Builder().ExcludeAggregateIDs().
			AggregateTypes(instance.AggregateType).
			EventTypes(instance.InstanceRemovedEventType).
			Builder(),
	)
	if err != nil {
		return fmt.Errorf("%s get instance IDs: %w", mig, err)
	}

	for _, instance := range instances {
		ctx := authz.WithInstanceID(ctx, instance)
		logging.Info("prepare initial webkeys for instance", "instance_id", instance, "migration", mig)
		_, err := mig.commands.SetInstanceFeatures(ctx, &command.InstanceFeatures{
			WebKey: gu.Ptr(true),
		})
		if err != nil {
			return fmt.Errorf("%s set webkey instance feature: %w", mig, err)
		}
	}
	return nil
}

func (mig *SetupWebkeys) String() string {
	return "59_setup_webkeys_2"
}
