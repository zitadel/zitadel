package policy

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

type LabelPolicyReadModel struct {
	eventstore.ReadModel

	PrimaryColor   string
	SecondaryColor string
}

func (rm *LabelPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *LabelPolicyAddedEvent:
			rm.PrimaryColor = e.PrimaryColor
			rm.SecondaryColor = e.SecondaryColor
		case *LabelPolicyChangedEvent:
			rm.PrimaryColor = e.PrimaryColor
			rm.SecondaryColor = e.SecondaryColor
		}
	}
	return rm.ReadModel.Reduce()
}

type LabelPolicyWriteModel struct {
	eventstore.WriteModel

	PrimaryColor   string
	SecondaryColor string
}

func (wm *LabelPolicyWriteModel) Reduce() error {
	return errors.ThrowUnimplemented(nil, "POLIC-xJjvN", "reduce unimpelemnted")
}

type LabelPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PrimaryColor   string `json:"primaryColor,omitempty"`
	SecondaryColor string `json:"secondaryColor,omitempty"`
}

func (e *LabelPolicyAddedEvent) CheckPrevious() bool {
	return true
}

func (e *LabelPolicyAddedEvent) Data() interface{} {
	return e
}

func NewLabelPolicyAddedEvent(
	base *eventstore.BaseEvent,
	primaryColor,
	secondaryColor string,
) *LabelPolicyAddedEvent {

	return &LabelPolicyAddedEvent{
		BaseEvent:      *base,
		PrimaryColor:   primaryColor,
		SecondaryColor: secondaryColor,
	}
}

func LabelPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &LabelPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-puqv4", "unable to unmarshal label policy")
	}

	return e, nil
}

type LabelPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PrimaryColor   string `json:"primaryColor,omitempty"`
	SecondaryColor string `json:"secondaryColor,omitempty"`
}

func (e *LabelPolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *LabelPolicyChangedEvent) Data() interface{} {
	return e
}

func NewLabelPolicyChangedEvent(
	base *eventstore.BaseEvent,
	current *LabelPolicyWriteModel,
	primaryColor,
	secondaryColor string,
) *LabelPolicyChangedEvent {

	e := &LabelPolicyChangedEvent{
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

func LabelPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &LabelPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-qhfFb", "unable to unmarshal label policy")
	}

	return e, nil
}

type LabelPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *LabelPolicyRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *LabelPolicyRemovedEvent) Data() interface{} {
	return nil
}

func NewLabelPolicyRemovedEvent(base *eventstore.BaseEvent) *LabelPolicyRemovedEvent {
	return &LabelPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func LabelPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &LabelPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
