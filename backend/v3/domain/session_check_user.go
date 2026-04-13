package domain

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CheckUserCommand struct {
	parent CheckUserParent

	userID    *string
	loginName *string

	user   *User
	factor *SessionFactorUser
}

// Result implements [Querier].
func (cmd *CheckUserCommand) Result() *User {
	return cmd.user
}

func NewCheckUserCommand(parent CheckUserParent, userID, loginName *string) *CheckUserCommand {
	cmd := &CheckUserCommand{
		parent:    parent,
		userID:    userID,
		loginName: loginName,
	}
	cmd.parent.SetUserConditionProvider(cmd.userCondition)
	return cmd
}

// Events implements [Commander].
func (cmd *CheckUserCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	var preferredLanguage *language.Tag
	if cmd.user.Human != nil && !cmd.user.Human.PreferredLanguage.IsRoot() {
		preferredLanguage = &cmd.user.Human.PreferredLanguage
	}
	fetchedSession, err := cmd.parent.FetchSession(ctx, opts)
	if err != nil {
		return nil, err
	}
	return []eventstore.Command{
		session.NewUserCheckedEvent(
			ctx,
			&session.NewAggregate(fetchedSession.ID, fetchedSession.InstanceID).Aggregate,
			cmd.user.ID,
			cmd.user.OrganizationID,
			cmd.factor.LastVerifiedAt,
			preferredLanguage,
		),
	}, nil
}

// Execute implements [Commander].
func (cmd *CheckUserCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	close, err := opts.ensureIsolated(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = close(err)
	}()
	user, err := cmd.parent.FetchUser(ctx, opts)
	if err != nil && zerrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if user.State != UserStateActive {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-vgDIu9", "Errors.User.NotActive")
	}
	session, err := cmd.parent.FetchSession(ctx, opts)
	if err != nil {
		return err
	}
	if session.UserID != "" && user.ID != "" && session.UserID != user.ID {
		return zerrors.ThrowInvalidArgument(nil, "DOM-78g1TV", "user change not possible")
	}
	cmd.factor = &SessionFactorUser{
		UserID:         user.ID,
		LastVerifiedAt: time.Now(), // TODO(adlerhurst): use a consistent time source
	}
	cmd.user = user

	return nil
}

// String implements [Commander].
func (cmd *CheckUserCommand) String() string {
	return "CheckUserCommand"
}

// Validate implements [Commander].
func (cmd *CheckUserCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if cmd.userID == nil && cmd.loginName == nil {
		return zerrors.ThrowInvalidArgument(nil, "DOMAI-D0UTe", "neither login name nor id provided")
	}
	return nil
}

func (cmd *CheckUserCommand) userCondition(ctx context.Context, opts *InvokeOpts) (condition database.Condition) {
	if cmd.userID != nil {
		return opts.userRepo.IDCondition(*cmd.userID)
	}
	return opts.userRepo.LoginNameCondition(database.TextOperationEqualIgnoreCase, *cmd.loginName)
}

// checkResult implements [sessionCheckSubCommand].
func (cmd *CheckUserCommand) checkResult() SessionFactor {
	return cmd.factor
}

var (
	_ Commander              = (*CheckUserCommand)(nil)
	_ Querier[*User]         = (*CheckUserCommand)(nil)
	_ sessionCheckSubCommand = (*CheckUserCommand)(nil)
)

//go:generate mockgen -typed -package domainmock -destination ./mock/session_check_user_parent.mock.go . CheckUserParent
type CheckUserParent interface {
	// SetUserConditionProvider is used to set the user condition provider for the command.
	SetUserConditionProvider(provider UserConditionProvider)

	// FetchSession is used to fetch the session.
	FetchSession(ctx context.Context, opts *InvokeOpts) (session *Session, err error)
	// FetchUser is used to fetch the user based on the condition set by setUserConditionProvider.
	// It might get called multiple times, so it should be implemented with caching in mind.
	FetchUser(ctx context.Context, opts *InvokeOpts) (user *User, err error)
}

type UserConditionProvider func(ctx context.Context, opts *InvokeOpts) (condition database.Condition)
