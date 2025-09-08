package domain

// import (
// 	"context"
// 	"time"
// )

// // EmailVerifiedCommand verifies an email address for a user.
// type EmailVerifiedCommand struct {
// 	UserID string `json:"userId"`
// 	Email  *Email `json:"email"`
// }

// func NewEmailVerifiedCommand(userID string, isVerified bool) *EmailVerifiedCommand {
// 	return &EmailVerifiedCommand{
// 		UserID: userID,
// 		Email: &Email{
// 			VerifiedAt: time.Time{},
// 		},
// 	}
// }

// // String implements [Commander].
// func (cmd *EmailVerifiedCommand) String() string {
// 	return "EmailVerifiedCommand"
// }

// var (
// 	_ Commander   = (*EmailVerifiedCommand)(nil)
// 	_ SetEmailOpt = (*EmailVerifiedCommand)(nil)
// )

// // Execute implements [Commander]
// func (cmd *EmailVerifiedCommand) Execute(ctx context.Context, opts *CommandOpts) error {
// 	repo := userRepo(opts.DB).Human()
// 	return repo.Update(ctx, repo.IDCondition(cmd.UserID), repo.SetEmailVerifiedAt(time.Time{}))
// }

// // applyOnSetEmail implements [SetEmailOpt]
// func (cmd *EmailVerifiedCommand) applyOnSetEmail(setEmailCmd *SetEmailCommand) {
// 	cmd.UserID = setEmailCmd.UserID
// 	cmd.Email.Address = setEmailCmd.Email
// 	setEmailCmd.verification = cmd
// }

// // SendCodeCommand sends a verification code to the user's email address.
// // If the URLTemplate is not set it will use the default of the organization / instance.
// type SendCodeCommand struct {
// 	UserID      string  `json:"userId"`
// 	Email       string  `json:"email"`
// 	URLTemplate *string `json:"urlTemplate"`
// 	generator   *generateCodeCommand
// }

// var (
// 	_ Commander   = (*SendCodeCommand)(nil)
// 	_ SetEmailOpt = (*SendCodeCommand)(nil)
// )

// func NewSendCodeCommand(userID string, urlTemplate *string) *SendCodeCommand {
// 	return &SendCodeCommand{
// 		UserID:      userID,
// 		generator:   &generateCodeCommand{},
// 		URLTemplate: urlTemplate,
// 	}
// }

// // String implements [Commander].
// func (cmd *SendCodeCommand) String() string {
// 	return "SendCodeCommand"
// }

// // Execute implements [Commander]
// func (cmd *SendCodeCommand) Execute(ctx context.Context, opts *CommandOpts) error {
// 	if err := cmd.ensureEmail(ctx, opts); err != nil {
// 		return err
// 	}
// 	if err := cmd.ensureURL(ctx, opts); err != nil {
// 		return err
// 	}

// 	if err := opts.Invoker.Invoke(ctx, cmd.generator, opts); err != nil {
// 		return err
// 	}
// 	// TODO: queue notification

// 	return nil
// }

// func (cmd *SendCodeCommand) ensureEmail(ctx context.Context, opts *CommandOpts) error {
// 	if cmd.Email != "" {
// 		return nil
// 	}
// 	repo := userRepo(opts.DB).Human()
// 	email, err := repo.GetEmail(ctx, repo.IDCondition(cmd.UserID))
// 	if err != nil || !email.VerifiedAt.IsZero() {
// 		return err
// 	}
// 	cmd.Email = email.Address
// 	return nil
// }

// func (cmd *SendCodeCommand) ensureURL(ctx context.Context, opts *CommandOpts) error {
// 	if cmd.URLTemplate != nil && *cmd.URLTemplate != "" {
// 		return nil
// 	}
// 	_, _ = ctx, opts
// 	// TODO: load default template
// 	return nil
// }

// // applyOnSetEmail implements [SetEmailOpt]
// func (cmd *SendCodeCommand) applyOnSetEmail(setEmailCmd *SetEmailCommand) {
// 	cmd.UserID = setEmailCmd.UserID
// 	cmd.Email = setEmailCmd.Email
// 	setEmailCmd.verification = cmd
// }

// // ReturnCodeCommand creates the code and returns it to the caller.
// // The caller gets the code by calling the Code field after the command got executed.
// type ReturnCodeCommand struct {
// 	UserID    string `json:"userId"`
// 	Email     string `json:"email"`
// 	Code      string `json:"code"`
// 	generator *generateCodeCommand
// }

// var (
// 	_ Commander   = (*ReturnCodeCommand)(nil)
// 	_ SetEmailOpt = (*ReturnCodeCommand)(nil)
// )

// func NewReturnCodeCommand(userID string) *ReturnCodeCommand {
// 	return &ReturnCodeCommand{
// 		UserID:    userID,
// 		generator: &generateCodeCommand{},
// 	}
// }

// // String implements [Commander].
// func (cmd *ReturnCodeCommand) String() string {
// 	return "ReturnCodeCommand"
// }

// // Execute implements [Commander]
// func (cmd *ReturnCodeCommand) Execute(ctx context.Context, opts *CommandOpts) error {
// 	if err := cmd.ensureEmail(ctx, opts); err != nil {
// 		return err
// 	}
// 	if err := opts.Invoker.Invoke(ctx, cmd.generator, opts); err != nil {
// 		return err
// 	}
// 	cmd.Code = cmd.generator.code
// 	return nil
// }

// func (cmd *ReturnCodeCommand) ensureEmail(ctx context.Context, opts *CommandOpts) error {
// 	if cmd.Email != "" {
// 		return nil
// 	}
// 	repo := userRepo(opts.DB).Human()
// 	email, err := repo.GetEmail(ctx, repo.IDCondition(cmd.UserID))
// 	if err != nil || !email.VerifiedAt.IsZero() {
// 		return err
// 	}
// 	cmd.Email = email.Address
// 	return nil
// }

// // applyOnSetEmail implements [SetEmailOpt]
// func (cmd *ReturnCodeCommand) applyOnSetEmail(setEmailCmd *SetEmailCommand) {
// 	cmd.UserID = setEmailCmd.UserID
// 	cmd.Email = setEmailCmd.Email
// 	setEmailCmd.verification = cmd
// }
