package command

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/v2/storage/eventstore"
)

var (
	_ eventstore.EventCommander = (*setEmail)(nil)
)

type setEmail struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

func SetEmail(userID, email string) *setEmail {
	return &setEmail{
		UserID: userID,
		Email:  email,
	}
}

// Event implements [eventstore.EventCommander].
func (c *setEmail) Event() *eventstore.Event {
	panic("unimplemented")
}

// Name implements [pattern.Command].
func (c *setEmail) Name() string {
	return "user.v2.set_email"
}

// Execute implements [pattern.Command].
func (c *setEmail) Execute(ctx context.Context) error {
	// Implementation of the command execution
	return nil
}
