package org

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const AddedType = eventTypePrefix + "added"

type addedPayload struct {
	Name string `json:"name"`
}

type AddedEvent eventstore.Event[addedPayload]

var _ eventstore.TypeChecker = (*AddedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *AddedEvent) ActionType() string {
	return AddedType
}

func AddedEventFromStorage(event *eventstore.StorageEvent) (e *AddedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Nf3tr", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[addedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &AddedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}

const uniqueOrgName = "org_name"

func NewAddedCommand(ctx context.Context, name string) (*eventstore.Command, error) {
	if name = strings.TrimSpace(name); name == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-mruNY", "Errors.Invalid.Argument")
	}
	return &eventstore.Command{
		Action: eventstore.Action[any]{
			Creator:  authz.GetCtxData(ctx).UserID,
			Type:     AddedType,
			Revision: 1,
			Payload: addedPayload{
				Name: name,
			},
		},
		UniqueConstraints: []*eventstore.UniqueConstraint{
			eventstore.NewAddEventUniqueConstraint(uniqueOrgName, name, "Errors.Org.AlreadyExists"),
		},
	}, nil
}
