package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

var (
	LabelPolicyAddedEventType     = instanceEventTypePrefix + policy.LabelPolicyAddedEventType
	LabelPolicyChangedEventType   = instanceEventTypePrefix + policy.LabelPolicyChangedEventType
	LabelPolicyActivatedEventType = instanceEventTypePrefix + policy.LabelPolicyActivatedEventType

	LabelPolicyLogoAddedEventType       = instanceEventTypePrefix + policy.LabelPolicyLogoAddedEventType
	LabelPolicyLogoRemovedEventType     = instanceEventTypePrefix + policy.LabelPolicyLogoRemovedEventType
	LabelPolicyIconAddedEventType       = instanceEventTypePrefix + policy.LabelPolicyIconAddedEventType
	LabelPolicyIconRemovedEventType     = instanceEventTypePrefix + policy.LabelPolicyIconRemovedEventType
	LabelPolicyLogoDarkAddedEventType   = instanceEventTypePrefix + policy.LabelPolicyLogoDarkAddedEventType
	LabelPolicyLogoDarkRemovedEventType = instanceEventTypePrefix + policy.LabelPolicyLogoDarkRemovedEventType
	LabelPolicyIconDarkAddedEventType   = instanceEventTypePrefix + policy.LabelPolicyIconDarkAddedEventType
	LabelPolicyIconDarkRemovedEventType = instanceEventTypePrefix + policy.LabelPolicyIconDarkRemovedEventType

	LabelPolicyFontAddedEventType   = instanceEventTypePrefix + policy.LabelPolicyFontAddedEventType
	LabelPolicyFontRemovedEventType = instanceEventTypePrefix + policy.LabelPolicyFontRemovedEventType
)

type LabelPolicyAddedEvent struct {
	policy.LabelPolicyAddedEvent
}

func NewLabelPolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	primaryColor,
	backgroundColor,
	warnColor,
	fontColor,
	primaryColorDark,
	backgroundColorDark,
	warnColorDark,
	fontColorDark string,
	hideLoginNameSuffix,
	errorMsgPopup,
	disableWatermark bool,
	themeMode domain.LabelPolicyThemeMode,
) *LabelPolicyAddedEvent {
	return &LabelPolicyAddedEvent{
		LabelPolicyAddedEvent: *policy.NewLabelPolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LabelPolicyAddedEventType),
			primaryColor,
			backgroundColor,
			warnColor,
			fontColor,
			primaryColorDark,
			backgroundColorDark,
			warnColorDark,
			fontColorDark,
			hideLoginNameSuffix,
			errorMsgPopup,
			disableWatermark,
			themeMode),
	}
}

func LabelPolicyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LabelPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyAddedEvent{LabelPolicyAddedEvent: *e.(*policy.LabelPolicyAddedEvent)}, nil
}

type LabelPolicyChangedEvent struct {
	policy.LabelPolicyChangedEvent
}

func NewLabelPolicyChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []policy.LabelPolicyChanges,
) (*LabelPolicyChangedEvent, error) {
	changedEvent, err := policy.NewLabelPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			LabelPolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &LabelPolicyChangedEvent{LabelPolicyChangedEvent: *changedEvent}, nil
}

func LabelPolicyChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LabelPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyChangedEvent{LabelPolicyChangedEvent: *e.(*policy.LabelPolicyChangedEvent)}, nil
}

type LabelPolicyActivatedEvent struct {
	policy.LabelPolicyActivatedEvent
}

func NewLabelPolicyActivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *LabelPolicyActivatedEvent {
	return &LabelPolicyActivatedEvent{
		LabelPolicyActivatedEvent: *policy.NewLabelPolicyActivatedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LabelPolicyActivatedEventType),
		),
	}
}

func LabelPolicyActivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LabelPolicyActivatedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyActivatedEvent{LabelPolicyActivatedEvent: *e.(*policy.LabelPolicyActivatedEvent)}, nil
}

type LabelPolicyLogoAddedEvent struct {
	policy.LabelPolicyLogoAddedEvent
}

func NewLabelPolicyLogoAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	storageKey string,
) *LabelPolicyLogoAddedEvent {
	return &LabelPolicyLogoAddedEvent{
		LabelPolicyLogoAddedEvent: *policy.NewLabelPolicyLogoAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LabelPolicyLogoAddedEventType),
			storageKey,
		),
	}
}

func LabelPolicyLogoAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LabelPolicyLogoAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyLogoAddedEvent{LabelPolicyLogoAddedEvent: *e.(*policy.LabelPolicyLogoAddedEvent)}, nil
}

type LabelPolicyLogoRemovedEvent struct {
	policy.LabelPolicyLogoRemovedEvent
}

func NewLabelPolicyLogoRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	storageKey string,
) *LabelPolicyLogoRemovedEvent {
	return &LabelPolicyLogoRemovedEvent{
		LabelPolicyLogoRemovedEvent: *policy.NewLabelPolicyLogoRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LabelPolicyLogoRemovedEventType),
			storageKey,
		),
	}
}

func LabelPolicyLogoRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LabelPolicyLogoRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyLogoRemovedEvent{LabelPolicyLogoRemovedEvent: *e.(*policy.LabelPolicyLogoRemovedEvent)}, nil
}

type LabelPolicyIconAddedEvent struct {
	policy.LabelPolicyIconAddedEvent
}

func NewLabelPolicyIconAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	storageKey string,
) *LabelPolicyIconAddedEvent {
	return &LabelPolicyIconAddedEvent{
		LabelPolicyIconAddedEvent: *policy.NewLabelPolicyIconAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LabelPolicyIconAddedEventType),
			storageKey,
		),
	}
}

func LabelPolicyIconAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LabelPolicyIconAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyIconAddedEvent{LabelPolicyIconAddedEvent: *e.(*policy.LabelPolicyIconAddedEvent)}, nil
}

type LabelPolicyIconRemovedEvent struct {
	policy.LabelPolicyIconRemovedEvent
}

func NewLabelPolicyIconRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	storageKey string,
) *LabelPolicyIconRemovedEvent {
	return &LabelPolicyIconRemovedEvent{
		LabelPolicyIconRemovedEvent: *policy.NewLabelPolicyIconRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LabelPolicyIconRemovedEventType),
			storageKey,
		),
	}
}

func LabelPolicyIconRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LabelPolicyIconRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyIconRemovedEvent{LabelPolicyIconRemovedEvent: *e.(*policy.LabelPolicyIconRemovedEvent)}, nil
}

type LabelPolicyLogoDarkAddedEvent struct {
	policy.LabelPolicyLogoDarkAddedEvent
}

func NewLabelPolicyLogoDarkAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	storageKey string,
) *LabelPolicyLogoDarkAddedEvent {
	return &LabelPolicyLogoDarkAddedEvent{
		LabelPolicyLogoDarkAddedEvent: *policy.NewLabelPolicyLogoDarkAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LabelPolicyLogoDarkAddedEventType),
			storageKey,
		),
	}
}

func LabelPolicyLogoDarkAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LabelPolicyLogoDarkAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyLogoDarkAddedEvent{LabelPolicyLogoDarkAddedEvent: *e.(*policy.LabelPolicyLogoDarkAddedEvent)}, nil
}

type LabelPolicyLogoDarkRemovedEvent struct {
	policy.LabelPolicyLogoDarkRemovedEvent
}

func NewLabelPolicyLogoDarkRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	storageKey string,
) *LabelPolicyLogoDarkRemovedEvent {
	return &LabelPolicyLogoDarkRemovedEvent{
		LabelPolicyLogoDarkRemovedEvent: *policy.NewLabelPolicyLogoDarkRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LabelPolicyLogoDarkRemovedEventType),
			storageKey,
		),
	}
}

func LabelPolicyLogoDarkRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LabelPolicyLogoDarkRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyLogoDarkRemovedEvent{LabelPolicyLogoDarkRemovedEvent: *e.(*policy.LabelPolicyLogoDarkRemovedEvent)}, nil
}

type LabelPolicyIconDarkAddedEvent struct {
	policy.LabelPolicyIconDarkAddedEvent
}

func NewLabelPolicyIconDarkAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	storageKey string,
) *LabelPolicyIconDarkAddedEvent {
	return &LabelPolicyIconDarkAddedEvent{
		LabelPolicyIconDarkAddedEvent: *policy.NewLabelPolicyIconDarkAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LabelPolicyIconDarkAddedEventType),
			storageKey,
		),
	}
}

func LabelPolicyIconDarkAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LabelPolicyIconDarkAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyIconDarkAddedEvent{LabelPolicyIconDarkAddedEvent: *e.(*policy.LabelPolicyIconDarkAddedEvent)}, nil
}

type LabelPolicyIconDarkRemovedEvent struct {
	policy.LabelPolicyIconDarkRemovedEvent
}

func NewLabelPolicyIconDarkRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	storageKey string,
) *LabelPolicyIconDarkRemovedEvent {
	return &LabelPolicyIconDarkRemovedEvent{
		LabelPolicyIconDarkRemovedEvent: *policy.NewLabelPolicyIconDarkRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LabelPolicyIconDarkRemovedEventType),
			storageKey,
		),
	}
}

func LabelPolicyIconDarkRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LabelPolicyIconDarkRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyIconDarkRemovedEvent{LabelPolicyIconDarkRemovedEvent: *e.(*policy.LabelPolicyIconDarkRemovedEvent)}, nil
}

type LabelPolicyFontAddedEvent struct {
	policy.LabelPolicyFontAddedEvent
}

func NewLabelPolicyFontAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	storageKey string,
) *LabelPolicyFontAddedEvent {
	return &LabelPolicyFontAddedEvent{
		LabelPolicyFontAddedEvent: *policy.NewLabelPolicyFontAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LabelPolicyFontAddedEventType),
			storageKey,
		),
	}
}

func LabelPolicyFontAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LabelPolicyFontAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyFontAddedEvent{LabelPolicyFontAddedEvent: *e.(*policy.LabelPolicyFontAddedEvent)}, nil
}

type LabelPolicyFontRemovedEvent struct {
	policy.LabelPolicyFontRemovedEvent
}

func NewLabelPolicyFontRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	storageKey string,
) *LabelPolicyFontRemovedEvent {
	return &LabelPolicyFontRemovedEvent{
		LabelPolicyFontRemovedEvent: *policy.NewLabelPolicyFontRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LabelPolicyFontRemovedEventType),
			storageKey,
		),
	}
}

func LabelPolicyFontRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LabelPolicyFontRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyFontRemovedEvent{LabelPolicyFontRemovedEvent: *e.(*policy.LabelPolicyFontRemovedEvent)}, nil
}

type LabelPolicyAssetsRemovedEvent struct {
	policy.LabelPolicyAssetsRemovedEvent
}

func (e *LabelPolicyAssetsRemovedEvent) Payload() interface{} {
	return nil
}

func (e *LabelPolicyAssetsRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func LabelPolicyAssetsRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LabelPolicyAssetsRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyAssetsRemovedEvent{LabelPolicyAssetsRemovedEvent: *e.(*policy.LabelPolicyAssetsRemovedEvent)}, nil
}
