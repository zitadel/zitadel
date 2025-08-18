package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueUserIDPLinkType  = "external_idps"
	UserIDPLinkEventPrefix = humanEventPrefix + "externalidp."
	idpLoginEventPrefix    = humanEventPrefix + "externallogin."

	UserIDPLinkAddedType               = UserIDPLinkEventPrefix + "added"
	UserIDPLinkRemovedType             = UserIDPLinkEventPrefix + "removed"
	UserIDPLinkCascadeRemovedType      = UserIDPLinkEventPrefix + "cascade.removed"
	UserIDPExternalIDMigratedType      = UserIDPLinkEventPrefix + "id.migrated"
	UserIDPExternalUsernameChangedType = UserIDPLinkEventPrefix + "username.changed"

	UserIDPLoginCheckSucceededType = idpLoginEventPrefix + "check.succeeded"
)

func NewAddUserIDPLinkUniqueConstraint(idpConfigID, externalUserID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueUserIDPLinkType,
		idpConfigID+externalUserID,
		"Errors.User.ExternalIDP.AlreadyExists")
}

func NewRemoveUserIDPLinkUniqueConstraint(idpConfigID, externalUserID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueUserIDPLinkType,
		idpConfigID+externalUserID)
}

type UserIDPLinkAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID    string `json:"idpConfigId,omitempty"`
	ExternalUserID string `json:"userId,omitempty"`
	DisplayName    string `json:"displayName,omitempty"`
}

func (e *UserIDPLinkAddedEvent) Payload() any {
	return e
}

func (e *UserIDPLinkAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddUserIDPLinkUniqueConstraint(e.IDPConfigID, e.ExternalUserID)}
}

func NewUserIDPLinkAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID,
	displayName,
	externalUserID string,
) *UserIDPLinkAddedEvent {
	return &UserIDPLinkAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserIDPLinkAddedType,
		),
		IDPConfigID:    idpConfigID,
		DisplayName:    displayName,
		ExternalUserID: externalUserID,
	}
}

func UserIDPLinkAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &UserIDPLinkAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-6M9sd", "unable to unmarshal user external idp added")
	}

	return e, nil
}

type UserIDPLinkRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID    string `json:"idpConfigId"`
	ExternalUserID string `json:"userId,omitempty"`
}

func (e *UserIDPLinkRemovedEvent) Payload() any {
	return e
}

func (e *UserIDPLinkRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveUserIDPLinkUniqueConstraint(e.IDPConfigID, e.ExternalUserID)}
}

func NewUserIDPLinkRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID,
	externalUserID string,
) *UserIDPLinkRemovedEvent {
	return &UserIDPLinkRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserIDPLinkRemovedType,
		),
		IDPConfigID:    idpConfigID,
		ExternalUserID: externalUserID,
	}
}

func UserIDPLinkRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &UserIDPLinkRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-eAWoT", "unable to unmarshal user external idp removed")
	}

	return e, nil
}

type UserIDPLinkCascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID    string `json:"idpConfigId"`
	ExternalUserID string `json:"userId,omitempty"`
}

func (e *UserIDPLinkCascadeRemovedEvent) Payload() any {
	return e
}

func (e *UserIDPLinkCascadeRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveUserIDPLinkUniqueConstraint(e.IDPConfigID, e.ExternalUserID)}
}

func NewUserIDPLinkCascadeRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID,
	externalUserID string,
) *UserIDPLinkCascadeRemovedEvent {
	return &UserIDPLinkCascadeRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserIDPLinkCascadeRemovedType,
		),
		IDPConfigID:    idpConfigID,
		ExternalUserID: externalUserID,
	}
}

func UserIDPLinkCascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &UserIDPLinkCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-dKGqO", "unable to unmarshal user external idp cascade removed")
	}

	return e, nil
}

type UserIDPCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *UserIDPCheckSucceededEvent) Payload() any {
	return e
}

func (e *UserIDPCheckSucceededEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewUserIDPCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo) *UserIDPCheckSucceededEvent {
	return &UserIDPCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserIDPLoginCheckSucceededType,
		),
		AuthRequestInfo: info,
	}
}

func UserIDPCheckSucceededEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &UserIDPCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-oikSS", "unable to unmarshal user external idp check succeeded")
	}

	return e, nil
}

type UserIDPExternalIDMigratedEvent struct {
	eventstore.BaseEvent `json:"-"`
	IDPConfigID          string `json:"idpConfigId"`
	PreviousID           string `json:"previousId"`
	NewID                string `json:"newId"`
}

func (e *UserIDPExternalIDMigratedEvent) Payload() any {
	return e
}

func (e *UserIDPExternalIDMigratedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *UserIDPExternalIDMigratedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewUserIDPExternalIDMigratedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID,
	previousID,
	newID string,
) *UserIDPExternalIDMigratedEvent {
	return &UserIDPExternalIDMigratedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserIDPExternalIDMigratedType,
		),
		IDPConfigID: idpConfigID,
		PreviousID:  previousID,
		NewID:       newID,
	}
}

type UserIDPExternalUsernameEvent struct {
	eventstore.BaseEvent `json:"-"`
	IDPConfigID          string `json:"idpConfigId"`
	ExternalUserID       string `json:"userId"`
	ExternalUsername     string `json:"username"`
}

func (e *UserIDPExternalUsernameEvent) Payload() any {
	return e
}

func (e *UserIDPExternalUsernameEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *UserIDPExternalUsernameEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewUserIDPExternalUsernameEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID,
	externalUserID,
	externalUsername string,
) *UserIDPExternalUsernameEvent {
	return &UserIDPExternalUsernameEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserIDPExternalUsernameChangedType,
		),
		IDPConfigID:      idpConfigID,
		ExternalUserID:   externalUserID,
		ExternalUsername: externalUsername,
	}
}
