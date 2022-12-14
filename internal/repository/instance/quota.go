package instance

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

var (
	QuotaAddedEventType    = instanceEventTypePrefix + quota.AddedEventType
	QuotaNotifiedEventType = instanceEventTypePrefix + quota.NotifiedEventType
	QuotaRemovedEventType  = instanceEventTypePrefix + quota.RemovedEventType
)

type QuotaAddedEvent struct {
	quota.AddedEvent `json:",inline"`
}

func NewQuotaAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	unit quota.Unit,
	from time.Time,
	interval time.Duration,
	amount uint64,
	limit bool,
	notifications []*quota.AddedEventNotification,
) *QuotaAddedEvent {
	return &QuotaAddedEvent{
		AddedEvent: *quota.NewAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				QuotaAddedEventType,
			),
			unit,
			from,
			interval,
			amount,
			limit,
			notifications,
		),
	}
}

func QuotaAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := quota.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &QuotaAddedEvent{AddedEvent: *e.(*quota.AddedEvent)}, nil
}

type QuotaNotifiedEvent struct {
	quota.NotifiedEvent `json:",inline"`
}

func NewQuotaNotifiedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	unit quota.Unit,
	id string,
	threshold uint64,
	usage uint64,
) *QuotaNotifiedEvent {
	return &QuotaNotifiedEvent{
		NotifiedEvent: *quota.NewNotifiedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				QuotaNotifiedEventType,
			),
			unit,
			id,
			threshold,
			usage,
		),
	}
}

func QuotaNotifiedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := quota.NotifiedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &QuotaNotifiedEvent{NotifiedEvent: *e.(*quota.NotifiedEvent)}, nil
}

type QuotaRemovedEvent struct {
	quota.RemovedEvent `json:",inline"`
}

func NewQuotaRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	unit quota.Unit,
) *QuotaRemovedEvent {
	return &QuotaRemovedEvent{
		RemovedEvent: *quota.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				QuotaRemovedEventType,
			),
			unit,
		),
	}
}

func QuotaRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := quota.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &QuotaRemovedEvent{RemovedEvent: *e.(*quota.RemovedEvent)}, nil
}
