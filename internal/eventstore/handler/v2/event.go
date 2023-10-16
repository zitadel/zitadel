package handler

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	schedulerSucceeded = eventstore.EventType("system.projections.scheduler.succeeded")
	aggregateType      = eventstore.AggregateType("system")
	aggregateID        = "SYSTEM"
)

func (h *Handler) didProjectionInitialize(ctx context.Context) bool {
	events, err := h.es.Filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID("").
		AddQuery().
		AggregateTypes(aggregateType).
		AggregateIDs(aggregateID).
		EventTypes(schedulerSucceeded).
		EventData(map[string]interface{}{
			"name": h.projection.Name(),
		}).
		Builder(),
	)
	return len(events) > 0 && err == nil
}

func (h *Handler) setSucceededOnce(ctx context.Context) error {
	_, err := h.es.Push(ctx, &ProjectionSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(ctx,
			eventstore.NewAggregate(ctx, aggregateID, aggregateType, "v1"),
			schedulerSucceeded,
		),
		Name: h.projection.Name(),
	})
	return err
}

type ProjectionSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
	Name                 string `json:"name"`
}

func (p *ProjectionSucceededEvent) Payload() interface{} {
	return p
}

func (p *ProjectionSucceededEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
