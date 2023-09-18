package quota

import (
	"context"
	"strconv"
	"time"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type Unit uint

const (
	UniqueQuotaNameType      = "quota_units"
	eventTypePrefix          = eventstore.EventType("quota.")
	AddedEventType           = eventTypePrefix + "added"
	NotifiedEventType        = eventTypePrefix + "notified"
	NotificationDueEventType = eventTypePrefix + "notificationdue"
	RemovedEventType         = eventTypePrefix + "removed"
)

const (
	Unimplemented Unit = iota
	RequestsAllAuthenticated
	ActionsAllRunsSeconds
)

func NewAddQuotaUnitUniqueConstraint(unit Unit) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueQuotaNameType,
		strconv.FormatUint(uint64(unit), 10),
		"Errors.Quota.AlreadyExists",
	)
}

func NewRemoveQuotaNameUniqueConstraint(unit Unit) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
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

func (e *AddedEvent) Payload() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddQuotaUnitUniqueConstraint(e.Unit)}
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

func AddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUOTA-4n8vs", "unable to unmarshal quota added")
	}

	return e, nil
}

type NotificationDueEvent struct {
	eventstore.BaseEvent `json:"-"`
	Unit                 Unit      `json:"unit"`
	ID                   string    `json:"id"`
	CallURL              string    `json:"callURL"`
	PeriodStart          time.Time `json:"periodStart"`
	Threshold            uint16    `json:"threshold"`
	Usage                uint64    `json:"usage"`
}

func (n *NotificationDueEvent) Payload() interface{} {
	return n
}

func (n *NotificationDueEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewNotificationDueEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	unit Unit,
	id string,
	callURL string,
	periodStart time.Time,
	threshold uint16,
	usage uint64,
) *NotificationDueEvent {
	return &NotificationDueEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			NotificationDueEventType,
		),
		Unit:        unit,
		ID:          id,
		CallURL:     callURL,
		PeriodStart: periodStart,
		Threshold:   threshold,
		Usage:       usage,
	}
}

func NotificationDueEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &NotificationDueEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUOTA-k56rT", "unable to unmarshal notification due")
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
	DueEventID           string    `json:"dueEventID"`
}

func (e *NotifiedEvent) Payload() interface{} {
	return e
}

func (e *NotifiedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewNotifiedEvent(
	ctx context.Context,
	id string,
	dueEvent *NotificationDueEvent,
) *NotifiedEvent {
	aggregate := dueEvent.Aggregate()
	return &NotifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			NotifiedEventType,
		),
		ID:         id,
		DueEventID: dueEvent.ID,
		// Deprecated: dereference the NotificationDueEvent
		Unit: dueEvent.Unit,
		// Deprecated: dereference the NotificationDueEvent
		CallURL: dueEvent.CallURL,
		// Deprecated: dereference the NotificationDueEvent
		PeriodStart: dueEvent.PeriodStart,
		// Deprecated: dereference the NotificationDueEvent
		Threshold: dueEvent.Threshold,
		// Deprecated: dereference the NotificationDueEvent
		Usage: dueEvent.Usage,
	}
}

func NotifiedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &NotifiedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUOTA-4n8vs", "unable to unmarshal quota notified")
	}

	return e, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	Unit                 Unit `json:"unit"`
}

func (e *RemovedEvent) Payload() interface{} {
	return e
}

func (e *RemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveQuotaNameUniqueConstraint(e.Unit)}
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

func RemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUOTA-4bReE", "unable to unmarshal quota removed")
	}

	return e, nil
}
