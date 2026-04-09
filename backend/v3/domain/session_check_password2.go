package domain

import (
	"context"
	"errors"
	"time"

	"github.com/zitadel/passwap"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CheckPasswordCommand struct {
	parent   CheckPasswordParent
	password string
	tarpit   tarpitFn
	verify   verifierFn

	validatedUser        *User
	updatedPassword      string
	passwordVerification VerificationType

	factor          *SessionFactorPassword
	verificationErr error
}

type PasswordCheckOption func(*CheckPasswordCommand)

func WithTarpitFunc(tarpitFunc tarpitFn) PasswordCheckOption {
	return func(cmd *CheckPasswordCommand) {
		cmd.tarpit = tarpitFunc
	}
}

func WithVerifierFn(verifierFn verifierFn) PasswordCheckOption {
	return func(cmd *CheckPasswordCommand) {
		cmd.verify = verifierFn
	}
}

func NewCheckPasswordCommand(parent CheckPasswordParent, password string, options ...PasswordCheckOption) *CheckPasswordCommand {
	cmd := &CheckPasswordCommand{
		parent:   parent,
		password: password,
		tarpit:   sysConfig.Tarpit.Tarpit(),
		verify:   passwordHasher.Verify,
	}
	for _, option := range options {
		option(cmd)
	}
	return cmd
}

// Events implements [Commander].
func (cmd *CheckPasswordCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if cmd.passwordVerification == nil {
		return nil, nil
	}

	fetchedSession, err := cmd.parent.FetchSession(ctx, opts)
	if err != nil {
		return nil, err
	}

	commands := make([]eventstore.Command, 0, 3)
	userAgg := &user.NewAggregate(cmd.validatedUser.ID, cmd.validatedUser.OrganizationID).Aggregate

	switch cmd.passwordVerification.(type) {
	case *VerificationTypeSucceeded:
		commands = append(commands,
			user.NewHumanPasswordCheckSucceededEvent(ctx, userAgg, nil),
			session.NewPasswordCheckedEvent(ctx,
				&session.NewAggregate(fetchedSession.ID, fetchedSession.InstanceID).Aggregate,
				cmd.factor.LastVerifiedAt,
			),
		)
		if cmd.updatedPassword != "" {
			commands = append(commands, user.NewHumanPasswordHashUpdatedEvent(ctx, userAgg, cmd.updatedPassword))
		}
	case *VerificationTypeFailed:
		commands = append(commands, user.NewHumanPasswordCheckFailedEvent(ctx, userAgg, nil))
		if cmd.validatedUser.State == UserStateLocked {
			commands = append(commands, user.NewUserLockedEvent(ctx, userAgg))
		}
	}

	return commands, nil
}

// Execute implements [Commander].
func (cmd *CheckPasswordCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if cmd.verificationErr != nil {
		cmd.tarpit(uint64(cmd.validatedUser.Human.Password.FailedAttempts + 1))
		return cmd.writeUserChanges(ctx, opts)
	}

	close, err := opts.ensureIsolated(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = close(err)
	}()

	latestUser, err := cmd.parent.reloadUser(ctx, opts)
	if err != nil {
		return err
	}
	if latestUser.Human == nil {
		return zerrors.ThrowInvalidArgument(nil, "DOM-9n8sdf", "user is not human")
	}

	if latestUser.Human.Password.Hash != cmd.validatedUser.Human.Password.Hash {
		return zerrors.ThrowInvalidArgument(nil, "DOM-9n8sdf", "password has changed since last check")
	}

	err = cmd.writeUserChanges(ctx, opts)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *CheckPasswordCommand) writeUserChanges(ctx context.Context, opts *InvokeOpts) error {
	changes := make([]database.Change, 0, 2)
	switch cmd.passwordVerification.(type) {
	case *VerificationTypeSucceeded:
		changes = append(changes, opts.userRepo.Human().SetLastSuccessfulPasswordCheck(time.Time{}))
		if cmd.updatedPassword != "" {
			changes = append(changes, opts.userRepo.Human().SetPassword(cmd.updatedPassword))
		}
	case *VerificationTypeFailed:
		changes = append(changes, opts.userRepo.Human().IncrementPasswordFailedAttempts())
		// TODO(adlerhurst): lock user if failed attempts exceed threshold
	}
	_, err := opts.userRepo.Human().Update(ctx, opts.DB(),
		opts.userRepo.PrimaryKeyCondition(cmd.validatedUser.InstanceID, cmd.validatedUser.ID),
		changes...,
	)
	return err
}

// String implements [Commander].
func (cmd *CheckPasswordCommand) String() string {
	return "CheckPasswordCommand"
}

// Validate implements [Commander].
func (cmd *CheckPasswordCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if cmd.parent == nil {
		return zerrors.ThrowInvalidArgument(nil, "DOM-9n8sdf", "parent must not be nil")
	}
	if cmd.password == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-9n8sdf", "password must not be empty")
	}
	return nil
}

// PreValidate implements [PreValidator].
func (cmd *CheckPasswordCommand) PreValidate(ctx context.Context, opts *InvokeOpts) (err error) {
	close, err := opts.ensureIsolated(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = close(err)
	}()

	cmd.validatedUser, err = cmd.parent.FetchUser(ctx, opts)
	if err != nil {
		return err
	}

	if cmd.validatedUser.Human == nil {
		return zerrors.ThrowInvalidArgument(nil, "DOM-9n8sdf", "user is not human")
	}

	cmd.updatedPassword, cmd.verificationErr = cmd.verify(cmd.validatedUser.Human.Password.Hash, cmd.password)
	if cmd.verificationErr == nil {
		cmd.passwordVerification = new(VerificationTypeSucceeded)
		cmd.factor = &SessionFactorPassword{
			LastVerifiedAt: time.Now(), // TODO(adlerhurst): use a consistent time source
		}
		return nil
	}

	cmd.factor = &SessionFactorPassword{
		LastFailedAt: time.Now(), // TODO(adlerhurst): use a consistent time source
	}
	if errors.Is(cmd.verificationErr, passwap.ErrPasswordMismatch) {
		cmd.verificationErr = zerrors.ThrowInvalidArgument(
			NewPasswordVerificationError(cmd.validatedUser.Human.Password.FailedAttempts+1),
			"DOM-3gcfDV",
			"Errors.User.Password.Invalid",
		)
	}
	cmd.passwordVerification = new(VerificationTypeFailed)
	return cmd.verificationErr
}

// checkResult implements [sessionCheckSubCommand].
func (cmd *CheckPasswordCommand) checkResult() SessionFactor {
	return cmd.factor
}

type CheckPasswordParent interface {
	FetchSession(ctx context.Context, opts *InvokeOpts) (session *Session, err error)
	FetchUser(ctx context.Context, opts *InvokeOpts) (user *User, err error)
	reloadUser(ctx context.Context, opts *InvokeOpts) (user *User, err error)
}

var (
	_ Commander              = (*CheckPasswordCommand)(nil)
	_ PreValidator           = (*CheckPasswordCommand)(nil)
	_ sessionCheckSubCommand = (*CheckPasswordCommand)(nil)
)
