package target

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	eventTypePrefix          eventstore.EventType = "execution."
	SetRequestEventType                           = eventTypePrefix + "request.set"
	SetResponseEventType                          = eventTypePrefix + "response.set"
	SetFunctionEventType                          = eventTypePrefix + "function.set"
	SetEventEventType                             = eventTypePrefix + "event.set"
	RemovedRequestEventType                       = eventTypePrefix + "request.removed"
	RemovedResponseEventType                      = eventTypePrefix + "response.removed"
	RemovedFunctionEventType                      = eventTypePrefix + "function.removed"
	RemovedEventEventType                         = eventTypePrefix + "event.removed"
)

type setEvent struct {
	Target  string `json:"target"`
	Include string `json:"include"`
}

type SetRequestEvent struct {
	*eventstore.BaseEvent `json:"-"`

	Method  string `json:"method"`
	Service string `json:"service"`
	All     bool   `json:"all"`

	setEvent
}

func (e *SetRequestEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *SetRequestEvent) Payload() any {
	return e
}

func (e *SetRequestEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.Method != "" {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionRequestMethod(e.Method))}
	}
	if e.Service != "" {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionRequestService(e.Service))}
	}
	if e.All {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionRequestAll())}
	}
	return nil
}

func NewSetRequestEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	method string,
	service string,
	all bool,
	target string,
	include string,
) *SetRequestEvent {
	return &SetRequestEvent{
		eventstore.NewBaseEventForPush(
			ctx, aggregate, SetRequestEventType,
		),
		method, service, all,
		setEvent{target, include},
	}
}

func SetRequestEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &SetRequestEvent{
		BaseEvent: eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "TARGET-fx8f8yfbn1", "unable to unmarshal execution request set")
	}

	return added, nil
}

type RemovedRequestEvent struct {
	*eventstore.BaseEvent `json:"-"`

	method  string
	service string
	all     bool
}

func (e *RemovedRequestEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *RemovedRequestEvent) Payload() any {
	return e
}

func (e *RemovedRequestEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.method != "" {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionRequestMethod(e.method))}
	}
	if e.service != "" {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionRequestService(e.service))}
	}
	if e.all {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionRequestAll())}
	}
	return nil
}

func NewRemovedRequestEvent(ctx context.Context, aggregate *eventstore.Aggregate, method, service string, all bool) *RemovedRequestEvent {
	return &RemovedRequestEvent{eventstore.NewBaseEventForPush(ctx, aggregate, RemovedRequestEventType), method, service, all}
}

func RemovedRequestEventMapper(event eventstore.Event) (eventstore.Event, error) {
	removed := &RemovedRequestEvent{
		BaseEvent: eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(removed)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "TARGET-0kuc12c7bc", "unable to unmarshal execution removed")
	}

	return removed, nil
}

type SetResponseEvent struct {
	*eventstore.BaseEvent `json:"-"`

	Method  string `json:"method"`
	Service string `json:"service"`
	All     bool   `json:"all"`

	setEvent
}

func (e *SetResponseEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *SetResponseEvent) Payload() any {
	return e
}

func (e *SetResponseEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.Method != "" {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionResponseMethod(e.Method))}
	}
	if e.Service != "" {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionResponseService(e.Service))}
	}
	if e.All {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionResponseAll())}
	}
	return nil
}

func NewSetResponseEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	method string,
	service string,
	all bool,
	target string,
	include string,
) *SetResponseEvent {
	return &SetResponseEvent{
		eventstore.NewBaseEventForPush(
			ctx, aggregate, SetResponseEventType,
		),
		method, service, all,
		setEvent{target, include}}
}

func SetResponseEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &SetResponseEvent{
		BaseEvent: eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "TARGET-fx8f8yfbn1", "unable to unmarshal execution response set")
	}

	return added, nil
}

type RemovedResponseEvent struct {
	*eventstore.BaseEvent `json:"-"`

	method  string
	service string
	all     bool
}

func (e *RemovedResponseEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *RemovedResponseEvent) Payload() any {
	return e
}

func (e *RemovedResponseEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.method != "" {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionResponseMethod(e.method))}
	}
	if e.service != "" {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionResponseService(e.service))}
	}
	if e.all {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionResponseAll())}
	}
	return nil
}

func NewRemovedResponseEvent(ctx context.Context, aggregate *eventstore.Aggregate, method, service string, all bool) *RemovedResponseEvent {
	return &RemovedResponseEvent{eventstore.NewBaseEventForPush(ctx, aggregate, RemovedResponseEventType), method, service, all}
}

func RemovedResponseEventMapper(event eventstore.Event) (eventstore.Event, error) {
	removed := &RemovedResponseEvent{
		BaseEvent: eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(removed)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "TARGET-0kuc12c7bc", "unable to unmarshal execution removed")
	}

	return removed, nil
}

type SetFunctionEvent struct {
	*eventstore.BaseEvent `json:"-"`

	Name string `json:"name"`

	setEvent
}

func (e *SetFunctionEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *SetFunctionEvent) Payload() any {
	return e
}

func (e *SetFunctionEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionFunction(e.Name))}
}

func NewSetFunctionEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
	target string,
	include string,
) *SetFunctionEvent {
	return &SetFunctionEvent{
		eventstore.NewBaseEventForPush(
			ctx, aggregate, SetFunctionEventType,
		),
		name,
		setEvent{target, include}}
}

func SetFunctionEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &SetFunctionEvent{
		BaseEvent: eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "TARGET-fx8f8yfbn1", "unable to unmarshal execution function set")
	}

	return added, nil
}

type RemovedFunctionEvent struct {
	*eventstore.BaseEvent `json:"-"`

	name string
}

func (e *RemovedFunctionEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *RemovedFunctionEvent) Payload() any {
	return e
}

func (e *RemovedFunctionEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionFunction(e.name))}
}

func NewRemovedFunctionEvent(ctx context.Context, aggregate *eventstore.Aggregate, name string) *RemovedFunctionEvent {
	return &RemovedFunctionEvent{eventstore.NewBaseEventForPush(ctx, aggregate, RemovedFunctionEventType), name}
}

func RemovedFunctionEventMapper(event eventstore.Event) (eventstore.Event, error) {
	removed := &RemovedFunctionEvent{
		BaseEvent: eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(removed)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "TARGET-0kuc12c7bc", "unable to unmarshal execution removed")
	}

	return removed, nil
}

type SetEventEvent struct {
	*eventstore.BaseEvent `json:"-"`

	Name  string `json:"name"`
	Group string `json:"group"`
	All   bool   `json:"all"`

	setEvent
}

func (e *SetEventEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *SetEventEvent) Payload() any {
	return e
}

func (e *SetEventEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.Name != "" {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionEvent(e.Name))}
	}
	if e.Group != "" {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionEventGroup(e.Group))}
	}
	if e.All {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionEventAll())}
	}
	return nil
}

func NewSetEventEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
	group string,
	all bool,
	target string,
	include string,
) *SetEventEvent {
	return &SetEventEvent{
		eventstore.NewBaseEventForPush(
			ctx, aggregate, SetEventEventType,
		),
		name, group, all,
		setEvent{target, include}}
}

func SetEventEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &SetEventEvent{
		BaseEvent: eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "TARGET-fx8f8yfbn1", "unable to unmarshal execution event set")
	}

	return added, nil
}

type RemovedEventEvent struct {
	*eventstore.BaseEvent `json:"-"`

	name  string
	group string
	all   bool
}

func (e *RemovedEventEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *RemovedEventEvent) Payload() any {
	return e
}

func (e *RemovedEventEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.name != "" {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionEvent(e.name))}
	}
	if e.group != "" {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionEventGroup(e.group))}
	}
	if e.all {
		return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(ConditionEventAll())}
	}
	return nil
}

func NewRemovedEventEvent(ctx context.Context, aggregate *eventstore.Aggregate, method, service string, all bool) *RemovedEventEvent {
	return &RemovedEventEvent{eventstore.NewBaseEventForPush(ctx, aggregate, RemovedEventEventType), method, service, all}
}

func RemovedEventEventMapper(event eventstore.Event) (eventstore.Event, error) {
	removed := &RemovedEventEvent{
		BaseEvent: eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(removed)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "TARGET-0kuc12c7bc", "unable to unmarshal execution removed")
	}

	return removed, nil
}
