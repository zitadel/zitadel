package label

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	LabelPolicyAddedEventType   = "policy.label.added"
	LabelPolicyChangedEventType = "policy.label.changed"
	LabelPolicyRemovedEventType = "policy.label.removed"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PrimaryColor   string `json:"primaryColor,omitempty"`
	SecondaryColor string `json:"secondaryColor,omitempty"`
}

func (e *AddedEvent) CheckPrevious() bool {
	return true
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func NewAddedEvent(
	base *eventstore.BaseEvent,
	primaryColor,
	secondaryColor string,
) *AddedEvent {

	return &AddedEvent{
		BaseEvent:      *base,
		PrimaryColor:   primaryColor,
		SecondaryColor: secondaryColor,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-puqv4", "unable to unmarshal label policy")
	}

	return e, nil
}

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PrimaryColor   string `json:"primaryColor,omitempty"`
	SecondaryColor string `json:"secondaryColor,omitempty"`
}

func (e *ChangedEvent) CheckPrevious() bool {
	return true
}

func (e *ChangedEvent) Data() interface{} {
	return e
}

func NewChangedEvent(
	base *eventstore.BaseEvent,
	current *WriteModel,
	primaryColor,
	secondaryColor string,
) *ChangedEvent {

	e := &ChangedEvent{
		BaseEvent: *base,
	}
	if primaryColor != "" && current.PrimaryColor != primaryColor {
		e.PrimaryColor = primaryColor
	}
	if secondaryColor != "" && current.SecondaryColor != secondaryColor {
		e.SecondaryColor = secondaryColor
	}

	return e
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-qhfFb", "unable to unmarshal label policy")
	}

	return e, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *RemovedEvent) CheckPrevious() bool {
	return true
}

func (e *RemovedEvent) Data() interface{} {
	return nil
}

func NewRemovedEvent(base *eventstore.BaseEvent) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *base,
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
