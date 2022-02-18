package project

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	SAMLConfigAddedType   = applicationEventTypePrefix + "config.saml.added"
	SAMLConfigChangedType = applicationEventTypePrefix + "config.saml.changed"
)

type SAMLConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID       string `json:"appId"`
	EntityID    string `json:"entityId"`
	Metadata    string `json:"metadata,omitempty"`
	MetadataURL string `json:"metadata_url,omitempty"`
}

func (e *SAMLConfigAddedEvent) Data() interface{} {
	return e
}

func (e *SAMLConfigAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewSAMLConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
	entityID string,
	metadata string,
	metadataURL string,
) *SAMLConfigAddedEvent {
	return &SAMLConfigAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SAMLConfigAddedType,
		),
		AppID:       appID,
		EntityID:    entityID,
		Metadata:    metadata,
		MetadataURL: metadataURL,
	}
}

func SAMLConfigAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &SAMLConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SAML-BDd15", "unable to unmarshal saml config")
	}

	return e, nil
}

type SAMLConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID       string  `json:"appId"`
	EntityID    string  `json:"entityId"`
	Metadata    *string `json:"metadata,omitempty"`
	MetadataURL *string `json:"metadata_url,omitempty"`
}

func (e *SAMLConfigChangedEvent) Data() interface{} {
	return e
}

func (e *SAMLConfigChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewSAMLConfigChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
	entityID string,
	changes []SAMLConfigChanges,
) (*SAMLConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "SAML-i8idç", "Errors.NoChangesFound")
	}

	changeEvent := &SAMLConfigChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SAMLConfigChangedType,
		),
		AppID: appID,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type SAMLConfigChanges func(event *SAMLConfigChangedEvent)

func ChangeMetadata(metadata string) func(event *SAMLConfigChangedEvent) {
	return func(e *SAMLConfigChangedEvent) {
		e.Metadata = &metadata
	}
}

func ChangeMetadataURL(metadataURL string) func(event *SAMLConfigChangedEvent) {
	return func(e *SAMLConfigChangedEvent) {
		e.MetadataURL = &metadataURL
	}
}

func SAMLConfigChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &SAMLConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SAML-BFd15", "unable to unmarshal saml config")
	}

	return e, nil
}
