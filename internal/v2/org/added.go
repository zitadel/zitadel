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
	_ eventstore.Command = (*AddedEvent)(nil)
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	Added *AddedEvent
)

type AddedEvent struct {
	Name string `json:"name"`

	creator string
}

func NewAddedEvent(ctx context.Context, name string) (*AddedEvent, error) {
	if name = strings.TrimSpace(name); name == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-mruNY", "Errors.Invalid.Argument")
	}
	return &AddedEvent{
		Name:    name,
		creator: authz.GetCtxData(ctx).UserID,
	}, nil
}

// Creator implements [eventstore.action].
func (a *AddedEvent) Creator() string {
	return a.creator
}

// Payload implements [eventstore.Command].
func (a *AddedEvent) Payload() any {
	return a
}

// Revision implements [eventstore.action].
func (*AddedEvent) Revision() uint16 {
	return 1
}

// UniqueConstraints implements [eventstore.Command].
func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		eventstore.NewAddEventUniqueConstraint(uniqueOrgName, e.Name, "Errors.Org.AlreadyExists"),
	}
}

// Type implements [eventstore.action].
func (*AddedEvent) Type() string {
	return "org.added"
}
