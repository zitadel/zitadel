package eventstore_test

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore/v3"
)

var _ eventstore.Cmd = (*UserAddedCommand)(nil)

type UserAddedCommand struct {
	eventstore.Command `json:"-"`

	Username string `json:"username"`
}

func NewUserAddedCommand(ctx context.Context) *UserAddedCommand {
	return &UserAddedCommand{
		Command: *eventstore.NewCommand(
			"user.added",
			1,
			eventstore.NewEditorFromCtx(ctx),
			eventstore.NewAggregate(ctx, "user-1", "user"),
		),
	}
}

func (cmd *UserAddedCommand) Payload() any {
	return cmd
}
