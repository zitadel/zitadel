package quota

import (
	"encoding/json"
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
		string(unit),
		"Errors.Quota.AlreadyExists",
	)
}

func NewRemoveQuotaNameUniqueConstraint(unit Unit) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		UniqueQuotaNameType,
		string(unit))
}

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Unit          Unit                      `json:"unit"`
	From          time.Time                 `json:"from"`
	Interval      time.Duration             `json:"interval,omitempty"`
	Amount        uint64                    `json:"amount"`
	Limit         bool                      `json:"limit"`
	Notifications []*AddedEventNotification `json:"notifications,omitempty"`
}

type AddedEventNotification struct {
	ID      string `json:"id"`
	Percent uint64 `json:"percent"`
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
	base *eventstore.BaseEvent,
	unit Unit,
	from time.Time,
	interval time.Duration,
	amount uint64,
	limit bool,
	notifications []*AddedEventNotification, // todo: redefine struct to receive here and convert to AddedEventNotification slice?
) *AddedEvent {
	return &AddedEvent{
		BaseEvent:     *base,
		Unit:          unit,
		From:          from,
		Interval:      interval,
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
	ID                   string `json:"id"`
	Unit                 Unit   `json:"unit"`
	Threshold            uint64 `json:"threshold"`
	Usage                uint64 `json:"usage"`
}

func (e *NotifiedEvent) Data() interface{} {
	return e
}

func (e *NotifiedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewNotifiedEvent(
	base *eventstore.BaseEvent,
	unit Unit,
	id string,
	threshold uint64,
	usage uint64,
) *NotifiedEvent {
	return &NotifiedEvent{
		BaseEvent: *base,
		Unit:      unit,
		ID:        id,
		Threshold: threshold,
		Usage:     usage,
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

	unit Unit
}

func (e *RemovedEvent) Data() interface{} {
	return e
}

func (e *RemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveQuotaNameUniqueConstraint(e.unit)}
}

func NewRemovedEvent(
	base *eventstore.BaseEvent,
	unit Unit,
) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *base,
		unit:      unit,
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	return &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
