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
	SetEventType             = eventTypePrefix + "set"
	NotifiedEventType        = eventTypePrefix + "notified"
	NotificationDueEventType = eventTypePrefix + "notificationdue"
	RemovedEventType         = eventTypePrefix + "removed"
)

const (
	Unimplemented Unit = iota
	RequestsAllAuthenticated
	ActionsAllRunsSeconds
)

func NewRemoveQuotaNameUniqueConstraint(unit Unit) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueQuotaNameType,
		strconv.FormatUint(uint64(unit), 10),
	)
}

// SetEvent describes that a quota is added or modified and contains only changed properties
type SetEvent struct {
	eventstore.BaseEvent `json:"-"`
	Unit                 Unit                     `json:"unit"`
	From                 *time.Time               `json:"from,omitempty"`
	ResetInterval        *time.Duration           `json:"interval,omitempty"`
	Amount               *uint64                  `json:"amount,omitempty"`
	Limit                *bool                    `json:"limit,omitempty"`
	Notifications        *[]*SetEventNotification `json:"notifications,omitempty"`
}

type SetEventNotification struct {
	ID      string `json:"id"`
	Percent uint16 `json:"percent"`
	Repeat  bool   `json:"repeat"`
	CallURL string `json:"callUrl"`
}

func (e *SetEvent) Payload() interface{} {
	return e
}

func (e *SetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewSetEvent(
	base *eventstore.BaseEvent,
	unit Unit,
	changes ...QuotaChange,
) *SetEvent {
	changedEvent := &SetEvent{
		BaseEvent: *base,
		Unit:      unit,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent
}

type QuotaChange func(*SetEvent)

func ChangeAmount(amount uint64) QuotaChange {
	return func(e *SetEvent) {
		e.Amount = &amount
	}
}

func ChangeLimit(limit bool) QuotaChange {
	return func(e *SetEvent) {
		e.Limit = &limit
	}
}

func ChangeFrom(from time.Time) QuotaChange {
	return func(event *SetEvent) {
		event.From = &from
	}
}

func ChangeResetInterval(interval time.Duration) QuotaChange {
	return func(event *SetEvent) {
		event.ResetInterval = &interval
	}
}

func ChangeNotifications(notifications []*SetEventNotification) QuotaChange {
	return func(event *SetEvent) {
		event.Notifications = &notifications
	}
}

func SetEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &SetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUOTA-kmIpI", "unable to unmarshal quota set")
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
