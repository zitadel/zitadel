package policy

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/asset"
)

const (
	LabelPolicyAddedEventType     = "policy.label.added"
	LabelPolicyChangedEventType   = "policy.label.changed"
	LabelPolicyActivatedEventType = "policy.label.activated"

	LabelPolicyLogoAddedEventType   = "policy.label.logo.added"
	LabelPolicyLogoRemovedEventType = "policy.label.logo.removed"
	LabelPolicyIconAddedEventType   = "policy.label.icon.added"
	LabelPolicyIconRemovedEventType = "policy.label.icon.removed"

	LabelPolicyLogoDarkAddedEventType   = "policy.label.logo.dark.added"
	LabelPolicyLogoDarkRemovedEventType = "policy.label.logo.dark.removed"
	LabelPolicyIconDarkAddedEventType   = "policy.label.icon.dark.added"
	LabelPolicyIconDarkRemovedEventType = "policy.label.icon.dark.removed"

	LabelPolicyFontAddedEventType   = "policy.label.font.added"
	LabelPolicyFontRemovedEventType = "policy.label.font.removed"

	LabelPolicyRemovedEventType = "policy.label.removed"
)

type LabelPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PrimaryColor        string `json:"primaryColor,omitempty"`
	BackgroundColor     string `json:"backgroundColor,omitempty"`
	WarnColor           string `json:"warnColor,omitempty"`
	FontColor           string `json:"fontColor,omitempty"`
	PrimaryColorDark    string `json:"primaryColorDark,omitempty"`
	BackgroundColorDark string `json:"backgroundColorDark,omitempty"`
	WarnColorDark       string `json:"warnColorDark,omitempty"`
	FontColorDark       string `json:"fontColorDark,omitempty"`
	HideLoginNameSuffix bool   `json:"hideLoginNameSuffix,omitempty"`
	ErrorMsgPopup       bool   `json:"errorMsgPopup,omitempty"`
	DisableWatermark    bool   `json:"disableMsgPopup,omitempty"`
}

func (e *LabelPolicyAddedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyAddedEvent(
	base *eventstore.BaseEvent,
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
) *LabelPolicyAddedEvent {

	return &LabelPolicyAddedEvent{
		BaseEvent:           *base,
		PrimaryColor:        primaryColor,
		BackgroundColor:     backgroundColor,
		WarnColor:           warnColor,
		FontColor:           fontColor,
		PrimaryColorDark:    primaryColorDark,
		BackgroundColorDark: backgroundColorDark,
		WarnColorDark:       warnColorDark,
		FontColorDark:       fontColorDark,
		HideLoginNameSuffix: hideLoginNameSuffix,
		ErrorMsgPopup:       errorMsgPopup,
		DisableWatermark:    disableWatermark,
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

	PrimaryColor        *string `json:"primaryColor,omitempty"`
	BackgroundColor     *string `json:"backgroundColor,omitempty"`
	WarnColor           *string `json:"warnColor,omitempty"`
	FontColor           *string `json:"fontColor,omitempty"`
	PrimaryColorDark    *string `json:"primaryColorDark,omitempty"`
	BackgroundColorDark *string `json:"backgroundColorDark,omitempty"`
	WarnColorDark       *string `json:"warnColorDark,omitempty"`
	FontColorDark       *string `json:"fontColorDark,omitempty"`
	HideLoginNameSuffix *bool   `json:"hideLoginNameSuffix,omitempty"`
	ErrorMsgPopup       *bool   `json:"errorMsgPopup,omitempty"`
	DisableWatermark    *bool   `json:"disableMsgPopup,omitempty"`
}

func (e *LabelPolicyChangedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []LabelPolicyChanges,
) (*LabelPolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "POLICY-Asfd3", "Errors.NoChangesFound")
	}
	changeEvent := &LabelPolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type LabelPolicyChanges func(*LabelPolicyChangedEvent)

func ChangePrimaryColor(primaryColor string) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.PrimaryColor = &primaryColor
	}
}

func ChangeBackgroundColor(background string) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.BackgroundColor = &background
	}
}

func ChangeWarnColor(warnColor string) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.WarnColor = &warnColor
	}
}

func ChangeFontColor(fontColor string) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.FontColor = &fontColor
	}
}

func ChangePrimaryColorDark(primaryColorDark string) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.PrimaryColorDark = &primaryColorDark
	}
}

func ChangeBackgroundColorDark(backgroundColorDark string) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.BackgroundColorDark = &backgroundColorDark
	}
}

func ChangeWarnColorDark(warnColorDark string) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.WarnColorDark = &warnColorDark
	}
}

func ChangeFontColorDark(fontColorDark string) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.FontColorDark = &fontColorDark
	}
}

func ChangeHideLoginNameSuffix(hideLoginNameSuffix bool) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.HideLoginNameSuffix = &hideLoginNameSuffix
	}
}

func ChangeErrorMsgPopup(errMsgPopup bool) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.ErrorMsgPopup = &errMsgPopup
	}
}

func ChangeDisableWatermark(disableWatermark bool) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.DisableWatermark = &disableWatermark
	}
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

type LabelPolicyActivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *LabelPolicyActivatedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyActivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyActivatedEvent(base *eventstore.BaseEvent) *LabelPolicyActivatedEvent {
	return &LabelPolicyActivatedEvent{
		BaseEvent: *base,
	}
}

func LabelPolicyActivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &LabelPolicyActivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type LabelPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *LabelPolicyRemovedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

type LabelPolicyLogoAddedEvent struct {
	asset.AddedEvent
}

func (e *LabelPolicyLogoAddedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyLogoAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyLogoAddedEvent(base *eventstore.BaseEvent, storageKey string) *LabelPolicyLogoAddedEvent {
	return &LabelPolicyLogoAddedEvent{
		*asset.NewAddedEvent(base, storageKey),
	}
}

func LabelPolicyLogoAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := asset.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyLogoAddedEvent{*e.(*asset.AddedEvent)}, nil
}

type LabelPolicyLogoRemovedEvent struct {
	asset.RemovedEvent
}

func (e *LabelPolicyLogoRemovedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyLogoRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyLogoRemovedEvent(base *eventstore.BaseEvent, storageKey string) *LabelPolicyLogoRemovedEvent {
	return &LabelPolicyLogoRemovedEvent{
		*asset.NewRemovedEvent(base, storageKey),
	}
}

func LabelPolicyLogoRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := asset.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyLogoRemovedEvent{*e.(*asset.RemovedEvent)}, nil
}

type LabelPolicyIconAddedEvent struct {
	asset.AddedEvent
}

func (e *LabelPolicyIconAddedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyIconAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyIconAddedEvent(base *eventstore.BaseEvent, storageKey string) *LabelPolicyIconAddedEvent {
	return &LabelPolicyIconAddedEvent{
		*asset.NewAddedEvent(base, storageKey),
	}
}

func LabelPolicyIconAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := asset.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyIconAddedEvent{*e.(*asset.AddedEvent)}, nil
}

type LabelPolicyIconRemovedEvent struct {
	asset.RemovedEvent
}

func (e *LabelPolicyIconRemovedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyIconRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyIconRemovedEvent(base *eventstore.BaseEvent, storageKey string) *LabelPolicyIconRemovedEvent {
	return &LabelPolicyIconRemovedEvent{
		*asset.NewRemovedEvent(base, storageKey),
	}
}

func LabelPolicyIconRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := asset.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyIconRemovedEvent{*e.(*asset.RemovedEvent)}, nil
}

type LabelPolicyLogoDarkAddedEvent struct {
	asset.AddedEvent
}

func (e *LabelPolicyLogoDarkAddedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyLogoDarkAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyLogoDarkAddedEvent(base *eventstore.BaseEvent, storageKey string) *LabelPolicyLogoDarkAddedEvent {
	return &LabelPolicyLogoDarkAddedEvent{
		*asset.NewAddedEvent(base, storageKey),
	}
}

func LabelPolicyLogoDarkAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := asset.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyLogoDarkAddedEvent{*e.(*asset.AddedEvent)}, nil
}

type LabelPolicyLogoDarkRemovedEvent struct {
	asset.RemovedEvent
}

func (e *LabelPolicyLogoDarkRemovedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyLogoDarkRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyLogoDarkRemovedEvent(base *eventstore.BaseEvent, storageKey string) *LabelPolicyLogoDarkRemovedEvent {
	return &LabelPolicyLogoDarkRemovedEvent{
		*asset.NewRemovedEvent(base, storageKey),
	}
}

func LabelPolicyLogoDarkRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := asset.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyLogoDarkRemovedEvent{*e.(*asset.RemovedEvent)}, nil
}

type LabelPolicyIconDarkAddedEvent struct {
	asset.AddedEvent
}

func (e *LabelPolicyIconDarkAddedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyIconDarkAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyIconDarkAddedEvent(base *eventstore.BaseEvent, storageKey string) *LabelPolicyIconDarkAddedEvent {
	return &LabelPolicyIconDarkAddedEvent{
		*asset.NewAddedEvent(base, storageKey),
	}
}

func LabelPolicyIconDarkAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := asset.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyIconDarkAddedEvent{*e.(*asset.AddedEvent)}, nil
}

type LabelPolicyIconDarkRemovedEvent struct {
	asset.RemovedEvent
}

func (e *LabelPolicyIconDarkRemovedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyIconDarkRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyIconDarkRemovedEvent(base *eventstore.BaseEvent, storageKey string) *LabelPolicyIconDarkRemovedEvent {
	return &LabelPolicyIconDarkRemovedEvent{
		*asset.NewRemovedEvent(base, storageKey),
	}
}

func LabelPolicyIconDarkRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := asset.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyIconDarkRemovedEvent{*e.(*asset.RemovedEvent)}, nil
}

type LabelPolicyFontAddedEvent struct {
	asset.AddedEvent
}

func (e *LabelPolicyFontAddedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyFontAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyFontAddedEvent(base *eventstore.BaseEvent, storageKey string) *LabelPolicyFontAddedEvent {
	return &LabelPolicyFontAddedEvent{
		*asset.NewAddedEvent(base, storageKey),
	}
}

func LabelPolicyFontAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := asset.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyFontAddedEvent{*e.(*asset.AddedEvent)}, nil
}

type LabelPolicyFontRemovedEvent struct {
	asset.RemovedEvent
}

func (e *LabelPolicyFontRemovedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyFontRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyFontRemovedEvent(base *eventstore.BaseEvent, storageKey string) *LabelPolicyFontRemovedEvent {
	return &LabelPolicyFontRemovedEvent{
		*asset.NewRemovedEvent(base, storageKey),
	}
}

func LabelPolicyFontRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := asset.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyFontRemovedEvent{*e.(*asset.RemovedEvent)}, nil
}
