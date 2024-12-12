package setup

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type FillFieldsForInstanceDomains struct {
	eventstore *eventstore.Eventstore
}

func (mig *FillFieldsForInstanceDomains) Execute(ctx context.Context, _ eventstore.Event) error {
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
	for _, instance := range instances {
		ctx := authz.WithInstanceID(ctx, instance)
		if err := projection.InstanceDomainFields.Trigger(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (mig *FillFieldsForInstanceDomains) String() string {
	return "41_fill_fields_for_instance_domains"
}
