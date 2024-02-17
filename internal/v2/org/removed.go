package org

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

var (
	_ eventstore.Command = (*RemovedEvent)(nil)
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	Removed *RemovedEvent
)

type RemovedEvent struct {
	creator string
}

func NewRemovedEvent(ctx context.Context) *RemovedEvent {
	return &RemovedEvent{
		creator: authz.GetCtxData(ctx).UserID,
	}
}

// Creator implements [eventstore.action].
func (e *RemovedEvent) Creator() string {
	return e.creator
}

// Payload implements [eventstore.Command].
func (*RemovedEvent) Payload() any {
	return nil
}

// Revision implements [eventstore.action].
func (*RemovedEvent) Revision() uint16 {
	return 1
}

// UniqueConstraints implements [eventstore.Command].
func (e *RemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		// TODO: soon as filter works
		eventstore.NewRemoveUniqueConstraint(uniqueOrgName, "test"),
		eventstore.NewRemoveUniqueConstraint(uniqueOrgDomain, "test.localhost"),
	}
}

// Type implements [eventstore.action].
func (*RemovedEvent) Type() string {
	return "org.removed"
}
