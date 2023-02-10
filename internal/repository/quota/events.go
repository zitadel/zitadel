package quota

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

type Unit uint

const (
	UniqueQuotaNameType           = "quota_units"
	UniqueQuotaNotificationIDType = "quota_notification"
	eventTypePrefix               = eventstore.EventType("quota.")
	AddedEventType                = eventTypePrefix + "added"
	NotifiedEventType             = eventTypePrefix + "notified"
	RemovedEventType              = eventTypePrefix + "removed"
)

const (
	Unimplemented Unit = iota
	RequestsAllAuthenticated
	ActionsAllRunsSeconds
)

func NewAddQuotaUnitUniqueConstraint(unit Unit) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueQuotaNameType,
		strconv.FormatUint(uint64(unit), 10),
		"Errors.Quota.AlreadyExists",
	)
}

func NewRemoveQuotaNameUniqueConstraint(unit Unit) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		UniqueQuotaNameType,
		strconv.FormatUint(uint64(unit), 10),
	)
}

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Unit          Unit                      `json:"unit"`
	From          time.Time                 `json:"from"`
	ResetInterval time.Duration             `json:"interval,omitempty"`
	Amount        uint64                    `json:"amount"`
	Limit         bool                      `json:"limit"`
	Notifications []*AddedEventNotification `json:"notifications,omitempty"`
}

type AddedEventNotification struct {
	ID      string `json:"id"`
	Percent uint16 `json:"percent"`
	Repeat  bool   `json:"repeat,omitempty"`
	CallURL string `json:"callUrl"`
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddQuotaUnitUniqueConstraint(e.Unit)}
}

func NewAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	unit Unit,
	from time.Time,
	resetInterval time.Duration,
	amount uint64,
	limit bool,
	notifications []*AddedEventNotification,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedEventType,
		),
		Unit:          unit,
		From:          from,
		ResetInterval: resetInterval,
		Amount:        amount,
		Limit:         limit,
		Notifications: notifications,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ACTION-4n8vs", "unable to unmarshal quota added")
	}

	return e, nil
}

type NotifiedEvent struct {
	eventstore.BaseEvent `json:"-"`
	Unit                 Unit      `json:"unit"`
	ID                   string    `json:"id"`
	CallURL              string    `json:"callURL"`
	PeriodStart          time.Time `json:"periodStart"`
	Threshold            uint16    `json:"threshold"`
	Usage                uint64    `json:"usage"`
}

func (e *NotifiedEvent) Data() interface{} {
	return e
}

func (e *NotifiedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewNotifiedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	unit Unit,
	id string,
	callURL string,
	periodStart time.Time,
	threshold uint16,
	usage uint64,
) *NotifiedEvent {
	return &NotifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			NotifiedEventType,
		),
		Unit:        unit,
		ID:          id,
		CallURL:     callURL,
		PeriodStart: periodStart,
		Threshold:   threshold,
		Usage:       usage,
	}
}

func NotifiedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &NotifiedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ACTION-4n8vs", "unable to unmarshal quota notified")
	}

	return e, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	Unit                 Unit `json:"unit"`
}

func (e *RemovedEvent) Data() interface{} {
	return e
}

func (e *RemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveQuotaNameUniqueConstraint(e.Unit)}
}

func NewRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	unit Unit,
) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			RemovedEventType,
		),
		Unit: unit,
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ACTION-4bReE", "unable to unmarshal quota removed")
	}

	return e, nil
}
