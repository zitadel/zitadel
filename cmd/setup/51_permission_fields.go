package setup

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type FillPermissionFields struct {
	eventstore *eventstore.Eventstore
}

func (mig *FillPermissionFields) Execute(ctx context.Context, _ eventstore.Event) error {
	instances, err := mig.eventstore.InstanceIDs(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsInstanceIDs).
			OrderDesc().
			AddQuery().
			AggregateTypes("instance").
			EventTypes(instance.InstanceAddedEventType).
			Builder().ExcludeAggregateIDs().
			AggregateTypes("instance").
			EventTypes(instance.InstanceRemovedEventType).
			Builder(),
	)
	if err != nil {
		return err
	}
	for _, instance := range instances {
		ctx := authz.WithInstanceID(ctx, instance)
		if err := projection.PermissionFields.Trigger(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (mig *FillPermissionFields) String() string {
	return "51_fill_permission_fields"
}
