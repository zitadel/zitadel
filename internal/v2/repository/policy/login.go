package policy

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/idp/provider"
)

const (
	LoginPolicyAddedEventType              = "policy.login.added"
	LoginPolicyChangedEventType            = "policy.login.changed"
	LoginPolicyRemovedEventType            = "policy.login.removed"
	LoginPolicyIDPProviderAddedEventType   = "policy.login." + provider.AddedEventType
	LoginPolicyIDPProviderRemovedEventType = "policy.login." + provider.RemovedEventType
)

type LoginPolicyReadModel struct {
	eventstore.ReadModel

	AllowUserNamePassword bool
	AllowRegister         bool
	AllowExternalIDP      bool
}

func (rm *LoginPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *LoginPolicyAddedEvent:
			rm.AllowUserNamePassword = e.AllowUserNamePassword
			rm.AllowExternalIDP = e.AllowExternalIDP
			rm.AllowRegister = e.AllowRegister
		case *LoginPolicyChangedEvent:
			rm.AllowUserNamePassword = e.AllowUserNamePassword
			rm.AllowExternalIDP = e.AllowExternalIDP
			rm.AllowRegister = e.AllowRegister
		}
	}
	return rm.ReadModel.Reduce()
}

type LoginPolicyWriteModel struct {
	eventstore.WriteModel

	AllowUserNamePassword bool
	AllowRegister         bool
	AllowExternalIDP      bool
}

func (wm *LoginPolicyWriteModel) Reduce() error {
	return errors.ThrowUnimplemented(nil, "POLIC-xJjvN", "reduce unimpelemnted")
}

type LoginPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AllowUserNamePassword bool `json:"allowUsernamePassword"`
	AllowRegister         bool `json:"allowRegister"`
	AllowExternalIDP      bool `json:"allowExternalIdp"`
}

func (e *LoginPolicyAddedEvent) CheckPrevious() bool {
	return true
}

func (e *LoginPolicyAddedEvent) Data() interface{} {
	return e
}

func NewLoginPolicyAddedEvent(
	base *eventstore.BaseEvent,
	allowUserNamePassword,
	allowRegister,
	allowExternalIDP bool,
) *LoginPolicyAddedEvent {

	return &LoginPolicyAddedEvent{
		BaseEvent:             *base,
		AllowExternalIDP:      allowExternalIDP,
		AllowRegister:         allowRegister,
		AllowUserNamePassword: allowUserNamePassword,
	}
}

func LoginPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &LoginPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-nWndT", "unable to unmarshal policy")
	}

	return e, nil
}

type LoginPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AllowUserNamePassword bool `json:"allowUsernamePassword,omitempty"`
	AllowRegister         bool `json:"allowRegister"`
	AllowExternalIDP      bool `json:"allowExternalIdp"`
}

func (e *LoginPolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *LoginPolicyChangedEvent) Data() interface{} {
	return e
}

func NewLoginPolicyChangedEvent(
	base *eventstore.BaseEvent,
	current *LoginPolicyWriteModel,
	allowUserNamePassword,
	allowRegister,
	allowExternalIDP bool,
) *LoginPolicyChangedEvent {

	e := &LoginPolicyChangedEvent{
		BaseEvent: *base,
	}

	if current.AllowUserNamePassword != allowUserNamePassword {
		e.AllowUserNamePassword = allowUserNamePassword
	}
	if current.AllowRegister != allowRegister {
		e.AllowRegister = allowRegister
	}
	if current.AllowExternalIDP != allowExternalIDP {
		e.AllowExternalIDP = allowExternalIDP
	}

	return e
}

func LoginPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &LoginPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-ehssl", "unable to unmarshal policy")
	}

	return e, nil
}

type LoginPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *LoginPolicyRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *LoginPolicyRemovedEvent) Data() interface{} {
	return nil
}

func NewLoginPolicyRemovedEvent(base *eventstore.BaseEvent) *LoginPolicyRemovedEvent {
	return &LoginPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func LoginPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &LoginPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type IDPProviderWriteModel struct {
	provider.WriteModel
}

func (wm *IDPProviderWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *IDPProviderAddedEvent:
			wm.WriteModel.AppendEvents(&e.AddedEvent)
		}
	}
}

type IDPProviderAddedEvent struct {
	provider.AddedEvent
}

func NewIDPProviderAddedEvent(
	base *eventstore.BaseEvent,
	idpConfigID string,
	idpProviderType provider.Type,
) *IDPProviderAddedEvent {

	return &IDPProviderAddedEvent{
		AddedEvent: *provider.NewAddedEvent(
			base,
			idpConfigID,
			idpProviderType),
	}
}

func IDPProviderAddedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := provider.AddedEventEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPProviderAddedEvent{
		AddedEvent: *e.(*provider.AddedEvent),
	}, nil
}

type IDPProviderRemovedEvent struct {
	provider.RemovedEvent
}

func NewIDPProviderRemovedEvent(
	base *eventstore.BaseEvent,
	idpConfigID string,
) *IDPProviderRemovedEvent {

	return &IDPProviderRemovedEvent{
		RemovedEvent: *provider.NewRemovedEvent(base, idpConfigID),
	}
}

func IDPProviderRemovedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := provider.RemovedEventEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPProviderRemovedEvent{
		RemovedEvent: *e.(*provider.RemovedEvent),
	}, nil
}
