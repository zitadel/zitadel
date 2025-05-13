package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
)

type CreateUserCommand struct {
	user  *User
	email *SetEmailCommand
}

var (
	_ Commander = (*CreateUserCommand)(nil)
	_ eventer   = (*CreateUserCommand)(nil)
)

func NewCreateHumanCommand(username string, opts ...CreateHumanOpt) *CreateUserCommand {
	cmd := &CreateUserCommand{
		user: &User{
			Username: username,
			Traits:   &Human{},
		},
	}

	for _, opt := range opts {
		opt.applyOnCreateHuman(cmd)
	}
	return cmd
}

// Events implements [eventer].
func (c *CreateUserCommand) Events() []*eventstore.Event {
	panic("unimplemented")
}

// Execute implements [Commander].
func (c *CreateUserCommand) Execute(ctx context.Context, opts *CommandOpts) error {
	if err := c.ensureUserID(); err != nil {
		return err
	}
	c.email.UserID = c.user.ID
	if err := opts.Invoke(ctx, c.email); err != nil {
		return err
	}
	return nil
}

type CreateHumanOpt interface {
	applyOnCreateHuman(*CreateUserCommand)
}

type createHumanIDOpt string

// applyOnCreateHuman implements [CreateHumanOpt].
func (c createHumanIDOpt) applyOnCreateHuman(cmd *CreateUserCommand) {
	cmd.user.ID = string(c)
}

var _ CreateHumanOpt = (*createHumanIDOpt)(nil)

func CreateHumanWithID(id string) CreateHumanOpt {
	return createHumanIDOpt(id)
}

func (c *CreateUserCommand) ensureUserID() (err error) {
	if c.user.ID != "" {
		return nil
	}
	c.user.ID, err = generateID()
	return err
}
