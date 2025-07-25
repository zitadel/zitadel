package domain

// import (
// 	"context"

// 	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
// )

// // SetEmailCommand sets the email address of a user.
// // If allows verification as a sub command.
// // The verification command is executed after the email address is set.
// // The verification command is executed in the same transaction as the email address update.
// type SetEmailCommand struct {
// 	UserID       string `json:"userId"`
// 	Email        string `json:"email"`
// 	verification Commander
// }

// var (
// 	_ Commander      = (*SetEmailCommand)(nil)
// 	_ eventer        = (*SetEmailCommand)(nil)
// 	_ CreateHumanOpt = (*SetEmailCommand)(nil)
// )

// type SetEmailOpt interface {
// 	applyOnSetEmail(*SetEmailCommand)
// }

// func NewSetEmailCommand(userID, email string, verificationType SetEmailOpt) *SetEmailCommand {
// 	cmd := &SetEmailCommand{
// 		UserID: userID,
// 		Email:  email,
// 	}
// 	verificationType.applyOnSetEmail(cmd)
// 	return cmd
// }

// // String implements [Commander].
// func (cmd *SetEmailCommand) String() string {
// 	return "SetEmailCommand"
// }

// func (cmd *SetEmailCommand) Execute(ctx context.Context, opts *CommandOpts) error {
// 	close, err := opts.EnsureTx(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() { err = close(ctx, err) }()
// 	// userStatement(opts.DB).Human().ByID(cmd.UserID).SetEmail(ctx, cmd.Email)
// 	repo := userRepo(opts.DB).Human()
// 	err = repo.Update(ctx, repo.IDCondition(cmd.UserID), repo.SetEmailAddress(cmd.Email))
// 	if err != nil {
// 		return err
// 	}

// 	return opts.Invoke(ctx, cmd.verification)
// }

// // Events implements [eventer].
// func (cmd *SetEmailCommand) Events() []*eventstore.Event {
// 	return []*eventstore.Event{
// 		{
// 			AggregateType: "user",
// 			AggregateID:   cmd.UserID,
// 			Type:          "user.email.set",
// 			Payload:       cmd,
// 		},
// 	}
// }

// // applyOnCreateHuman implements [CreateHumanOpt].
// func (cmd *SetEmailCommand) applyOnCreateHuman(createUserCmd *CreateUserCommand) {
// 	createUserCmd.email = cmd
// }
