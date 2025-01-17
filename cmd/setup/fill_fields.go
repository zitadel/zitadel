package setup

import (
	"context"
	"fmt"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type RepeatableFillFields struct {
	eventstore *eventstore.Eventstore
	handlers   []*handler.FieldHandler
}

func (mig *RepeatableFillFields) Execute(ctx context.Context, _ eventstore.Event) error {
	instances, err := mig.eventstore.InstanceIDs(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsInstanceIDs).
			OrderDesc().
			AddQuery().
			AggregateTypes(instance.AggregateType).
			EventTypes(instance.InstanceAddedEventType).
			Builder(),
	)
	if err != nil {
		return err
	}
	for _, instance := range instances {
		ctx := authz.WithInstanceID(ctx, instance)
		for _, handler := range mig.handlers {
			logging.WithFields("migration", mig.String(), "instance_id", instance, "handler", handler.String()).Info("run fields trigger")
			if err := handler.Trigger(ctx); err != nil {
				return fmt.Errorf("%s: %s: %w", mig.String(), handler.String(), err)
			}
		}
	}
	return nil
}

func (mig *RepeatableFillFields) String() string {
	return "repeatable_fill_fields"
}

func (f *RepeatableFillFields) Check(lastRun map[string]interface{}) bool {
	return true
}
