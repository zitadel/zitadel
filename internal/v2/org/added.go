package org

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var uniqueOrgName = "org_name"

var (
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	Added AddedEvent
)

type addedPayload struct {
	Name string `json:"name"`
}

type AddedEvent addedEvent
type addedEvent = eventstore.Event[addedPayload]

func AddedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*AddedEvent, error) {
	event, err := eventstore.EventFromStorage[addedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*AddedEvent)(event), nil
}

func (e AddedEvent) IsType(typ string) bool {
	return typ == "org.added"
}

var _ eventstore.Command = (*AddedCommand)(nil)

type AddedCommand struct {
	creator string
	addedPayload
}

func NewAddedCommand(ctx context.Context, name string) (*AddedCommand, error) {
	if name = strings.TrimSpace(name); name == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-mruNY", "Errors.Invalid.Argument")
	}
	return &AddedCommand{
		creator: authz.GetCtxData(ctx).UserID,
		addedPayload: addedPayload{
			Name: name,
		},
	}, nil
}

// Creator implements eventstore.Command.
func (a *AddedCommand) Creator() string {
	return a.creator
}

// Payload implements eventstore.Command.
func (a *AddedCommand) Payload() any {
	return a.addedPayload
}

// Revision implements [eventstore.action].
func (*AddedCommand) Revision() uint16 {
	return 1
}

// UniqueConstraints implements [eventstore.Command].
func (e *AddedCommand) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		eventstore.NewAddEventUniqueConstraint(uniqueOrgName, e.Name, "Errors.Org.AlreadyExists"),
	}
}

// Type implements [eventstore.action].
func (*AddedCommand) Type() string {
	return "org.added"
}
