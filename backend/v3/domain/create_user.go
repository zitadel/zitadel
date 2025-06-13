package domain

// import (
// 	"context"

// 	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
// )

// // CreateUserCommand adds a new user including the email verification for humans.
// // In the future it might make sense to separate the command into two commands:
// // - CreateHumanCommand: creates a new human user
// // - CreateMachineCommand: creates a new machine user
// type CreateUserCommand struct {
// 	user  *User
// 	email *SetEmailCommand
// }

// var (
// 	_ Commander = (*CreateUserCommand)(nil)
// 	_ eventer   = (*CreateUserCommand)(nil)
// )

// // opts heavily reduces the complexity for email verification because each type of verification is a simple option which implements the [Commander] interface.
// func NewCreateHumanCommand(username string, opts ...CreateHumanOpt) *CreateUserCommand {
// 	cmd := &CreateUserCommand{
// 		user: &User{
// 			Username: username,
// 			Traits:   &Human{},
// 		},
// 	}

// 	for _, opt := range opts {
// 		opt.applyOnCreateHuman(cmd)
// 	}
// 	return cmd
// }

// // String implements [Commander].
// func (cmd *CreateUserCommand) String() string {
// 	return "CreateUserCommand"
// }

// // Events implements [eventer].
// func (c *CreateUserCommand) Events() []*eventstore.Event {
// 	return []*eventstore.Event{
// 		{
// 			AggregateType: "user",
// 			AggregateID:   c.user.ID,
// 			Type:          "user.added",
// 			Payload:       c.user,
// 		},
// 	}
// }

// // Execute implements [Commander].
// func (c *CreateUserCommand) Execute(ctx context.Context, opts *CommandOpts) error {
// 	if err := c.ensureUserID(); err != nil {
// 		return err
// 	}
// 	c.email.UserID = c.user.ID
// 	if err := opts.Invoke(ctx, c.email); err != nil {
// 		return err
// 	}
// 	return nil
// }

// type CreateHumanOpt interface {
// 	applyOnCreateHuman(*CreateUserCommand)
// }

// type createHumanIDOpt string

// // applyOnCreateHuman implements [CreateHumanOpt].
// func (c createHumanIDOpt) applyOnCreateHuman(cmd *CreateUserCommand) {
// 	cmd.user.ID = string(c)
// }

// var _ CreateHumanOpt = (*createHumanIDOpt)(nil)

// func CreateHumanWithID(id string) CreateHumanOpt {
// 	return createHumanIDOpt(id)
// }

// func (c *CreateUserCommand) ensureUserID() (err error) {
// 	if c.user.ID != "" {
// 		return nil
// 	}
// 	c.user.ID, err = generateID()
// 	return err
// }
