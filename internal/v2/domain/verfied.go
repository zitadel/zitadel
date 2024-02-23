package domain

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type VerifiedEvent struct {
	Name string `json:"domain"`

	creator string
}

func NewVerifiedEvent(ctx context.Context, name string) (*VerifiedEvent, error) {
	if name = strings.TrimSpace(name); name == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "DOMAI-2zkf1", "Errors.Invalid.Argument")
	}
	return &VerifiedEvent{
		Name:    name,
		creator: authz.GetCtxData(ctx).UserID,
	}, nil
}

// Creator implements [eventstore.action].
func (a *VerifiedEvent) Creator() string {
	return a.creator
}

// Payload implements [eventstore.Command].
func (a *VerifiedEvent) Payload() any {
	return a
}

// Revision implements [eventstore.action].
func (*VerifiedEvent) Revision() uint16 {
	return 1
}

func (*VerifiedEvent) Type() string {
	return "domain.verified"
}
