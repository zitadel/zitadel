package domain

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddedEvent struct {
	Name string `json:"domain"`

	creator string
}

func NewAddedEvent(ctx context.Context, name string) (*AddedEvent, error) {
	if name = strings.TrimSpace(name); name == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "DOMAI-HaT0m", "Errors.Invalid.Argument")
	}
	return &AddedEvent{
		Name:    name,
		creator: authz.GetCtxData(ctx).UserID,
	}, nil
}

// Creator implements eventstore.Command.
func (e *AddedEvent) Creator() string {
	return e.creator
}

// Payload implements eventstore.Command.
func (e *AddedEvent) Payload() any {
	return e
}

// Revision implements eventstore.Command.
func (*AddedEvent) Revision() uint16 {
	return 1
}

func (*AddedEvent) Type() string {
	return "domain.added"
}
