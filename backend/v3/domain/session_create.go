package domain

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/activity"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CreateSessionCommand struct {
	session *Session

	checks     []sessionCheckSubCommand
	challenges []sessionChallengeSubCommand

	user lazyGetter[*User]
}

// Result implements [Querier].
func (cmd *CreateSessionCommand) Result() *Session {
	return cmd.session
}

type CreateSessionOption interface {
	ApplyOnCreateSessionCommand(cmd *CreateSessionCommand)
}

func NewCreateSessionCommand(instanceID, creatorID string, userAgent *SessionUserAgent, opts ...CreateSessionOption) *CreateSessionCommand {
	cmd := &CreateSessionCommand{
		session: &Session{
			InstanceID: instanceID,
			CreatorID:  creatorID,
			UserAgent:  userAgent,
		},
	}

	for _, opt := range opts {
		opt.ApplyOnCreateSessionCommand(cmd)
	}

	return cmd
}

// Events implements [Commander].
func (cmd *CreateSessionCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	var userAgent *domain.UserAgent
	if cmd.session.UserAgent != nil {
		userAgent = &domain.UserAgent{
			FingerprintID: cmd.session.UserAgent.FingerprintID,
			IP:            cmd.session.UserAgent.IP,
			Description:   cmd.session.UserAgent.Description,
			Header:        cmd.session.UserAgent.Header,
		}
	}
	aggregate := &session.NewAggregate(cmd.session.ID, cmd.session.InstanceID).Aggregate
	commands := []eventstore.Command{
		session.NewAddedEvent(ctx,
			aggregate,
			userAgent,
		),
		session.NewTokenSetEvent(ctx, aggregate, cmd.session.TokenID),
	}

	activity.TriggerWithoutOrg(ctx, cmd.session.UserID, activity.SessionAPI)
	if len(cmd.session.Metadata) == 0 {
		return commands, nil
	}
	metadata := make(map[string][]byte, len(cmd.session.Metadata))
	for _, md := range cmd.session.Metadata {
		metadata[md.Key] = md.Value
	}
	return append(commands, session.NewMetadataSetEvent(ctx, aggregate, metadata)), nil
}

// Execute implements [Commander].
func (cmd *CreateSessionCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if cmd.session.ID == "" {
		cmd.session.ID = opts.MustNewID()
	}
	for _, check := range cmd.checks {
		cmd.session.Factors = append(cmd.session.Factors, check.checkResult())
	}
	for _, challenge := range cmd.challenges {
		cmd.session.Challenges = append(cmd.session.Challenges, challenge.challengeResult())
	}
	return opts.sessionRepo.Create(ctx, opts.DB(), cmd.session)
}

// String implements [Commander].
func (cmd *CreateSessionCommand) String() string {
	return "CreateSessionCommand"
}

// Validate implements [Commander].
func (cmd *CreateSessionCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if cmd.session.InstanceID = strings.TrimSpace(cmd.session.InstanceID); cmd.session.InstanceID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-9n8sdf", "invalid instance ID")
	}
	if cmd.session.CreatorID = strings.TrimSpace(cmd.session.CreatorID); cmd.session.CreatorID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-9n8sdf", "invalid creator ID")
	}
	return nil
}

// SetUserConditionProvider implements [CheckUserParent].
func (cmd *CreateSessionCommand) SetUserConditionProvider(provider UserConditionProvider) {
	cmd.user = lazyGetter[*User]{
		get: func(ctx context.Context, opts *InvokeOpts) (*User, error) {
			return opts.userRepo.Get(ctx, opts.DB(), database.WithCondition(database.And(
				opts.userRepo.InstanceIDCondition(cmd.session.InstanceID),
				provider(ctx, opts),
			)))
		},
	}
}

// FetchSession implements [CheckPasswordParent] and [CheckUserParent].
func (cmd *CreateSessionCommand) FetchSession(ctx context.Context, opts *InvokeOpts) (session *Session, err error) {
	if cmd.session.ID == "" {
		cmd.session.ID = opts.MustNewID()
	}
	return cmd.session, nil
}

// FetchUser implements [CheckUserParent].
func (cmd *CreateSessionCommand) FetchUser(ctx context.Context, opts *InvokeOpts) (user *User, err error) {
	return cmd.user.fetch(ctx, opts)
}

// reloadUser implements [CheckUserParent].
func (cmd *CreateSessionCommand) reloadUser(ctx context.Context, opts *InvokeOpts) (user *User, err error) {
	return cmd.user.reload(ctx, opts)
}

var (
	_ Commander           = (*CreateSessionCommand)(nil)
	_ CheckUserParent     = (*CreateSessionCommand)(nil)
	_ CheckPasswordParent = (*CreateSessionCommand)(nil)
	_ Querier[*Session]   = (*CreateSessionCommand)(nil)
)
