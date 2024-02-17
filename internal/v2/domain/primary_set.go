package domain

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type PrimarySetEvent struct {
	Name string `json:"domain"`

	creator string
}

func NewSetPrimaryEvent(ctx context.Context, name string) (*PrimarySetEvent, error) {
	if name = strings.TrimSpace(name); name == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "DOMAI-6ofTB", "Errors.Invalid.Argument")
	}
	return &PrimarySetEvent{
		Name:    name,
		creator: authz.GetCtxData(ctx).UserID,
	}, nil
}

// Creator implements [eventstore.action].
func (a *PrimarySetEvent) Creator() string {
	return a.creator
}

// Payload implements [eventstore.Command].
func (a *PrimarySetEvent) Payload() any {
	return a
}

// Revision implements [eventstore.action].
func (*PrimarySetEvent) Revision() uint16 {
	return 1
}

func (*PrimarySetEvent) Type() string {
	return "domain.primary.set"
}
