package project

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueEntityIDType    = "entity_ids"
	SAMLConfigAddedType   = applicationEventTypePrefix + "config.saml.added"
	SAMLConfigChangedType = applicationEventTypePrefix + "config.saml.changed"
)

type SAMLConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID        string              `json:"appId"`
	EntityID     string              `json:"entityId"`
	Metadata     []byte              `json:"metadata,omitempty"`
	MetadataURL  string              `json:"metadata_url,omitempty"`
	LoginVersion domain.LoginVersion `json:"loginVersion,omitempty"`
	LoginBaseURI string              `json:"loginBaseURI,omitempty"`
}

func (e *SAMLConfigAddedEvent) Payload() interface{} {
	return e
}

func NewAddSAMLConfigEntityIDUniqueConstraint(entityID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueEntityIDType,
		entityID,
		"Errors.Project.App.SAMLEntityIDAlreadyExists")
}

func NewRemoveSAMLConfigEntityIDUniqueConstraint(entityID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueEntityIDType,
		entityID)
}

func (e *SAMLConfigAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddSAMLConfigEntityIDUniqueConstraint(e.EntityID)}
}

func NewSAMLConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
	entityID string,
	metadata []byte,
	metadataURL string,
	loginVersion domain.LoginVersion,
	loginBaseURI string,
) *SAMLConfigAddedEvent {
	return &SAMLConfigAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SAMLConfigAddedType,
		),
		AppID:        appID,
		EntityID:     entityID,
		Metadata:     metadata,
		MetadataURL:  metadataURL,
		LoginVersion: loginVersion,
		LoginBaseURI: loginBaseURI,
	}
}

func SAMLConfigAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &SAMLConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "SAML-BDd15", "unable to unmarshal saml config")
	}

	return e, nil
}

type SAMLConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID        string               `json:"appId"`
	EntityID     string               `json:"entityId"`
	Metadata     []byte               `json:"metadata,omitempty"`
	MetadataURL  *string              `json:"metadata_url,omitempty"`
	LoginVersion *domain.LoginVersion `json:"loginVersion,omitempty"`
	LoginBaseURI *string              `json:"loginBaseURI,omitempty"`
	oldEntityID  string
}

func (e *SAMLConfigChangedEvent) Payload() interface{} {
	return e
}

func (e *SAMLConfigChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.EntityID != "" {
		return []*eventstore.UniqueConstraint{
			NewRemoveSAMLConfigEntityIDUniqueConstraint(e.oldEntityID),
			NewAddSAMLConfigEntityIDUniqueConstraint(e.EntityID),
		}
	}
	return nil
}

func NewSAMLConfigChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
	oldEntityID string,
	changes []SAMLConfigChanges,
) (*SAMLConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "SAML-i8idç", "Errors.NoChangesFound")
	}

	changeEvent := &SAMLConfigChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SAMLConfigChangedType,
		),
		AppID:       appID,
		oldEntityID: oldEntityID,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type SAMLConfigChanges func(event *SAMLConfigChangedEvent)

func ChangeMetadata(metadata []byte) func(event *SAMLConfigChangedEvent) {
	return func(e *SAMLConfigChangedEvent) {
		e.Metadata = metadata
	}
}

func ChangeMetadataURL(metadataURL string) func(event *SAMLConfigChangedEvent) {
	return func(e *SAMLConfigChangedEvent) {
		e.MetadataURL = &metadataURL
	}
}

func ChangeEntityID(entityID string) func(event *SAMLConfigChangedEvent) {
	return func(e *SAMLConfigChangedEvent) {
		e.EntityID = entityID
	}
}

func ChangeSAMLLoginVersion(loginVersion domain.LoginVersion) func(event *SAMLConfigChangedEvent) {
	return func(e *SAMLConfigChangedEvent) {
		e.LoginVersion = &loginVersion
	}
}
func ChangeSAMLLoginBaseURI(loginBaseURI string) func(event *SAMLConfigChangedEvent) {
	return func(e *SAMLConfigChangedEvent) {
		e.LoginBaseURI = &loginBaseURI
	}
}

func SAMLConfigChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &SAMLConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "SAML-BFd15", "unable to unmarshal saml config")
	}

	return e, nil
}
