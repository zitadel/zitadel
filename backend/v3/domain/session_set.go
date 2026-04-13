package domain

import (
	"bytes"
	"context"
	"slices"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/activity"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SetSessionCommand struct {
	instanceID string
	sessionID  string

	checks     []sessionCheckSubCommand
	challenges []sessionChallengeSubCommand

	lifetime *time.Duration
	metadata []*SessionMetadata

	user    *lazyGetter[*User]
	session lazyGetter[*Session]
}

// Result implements [Querier].
func (cmd *SetSessionCommand) Result() *Session {
	return nil
	// return cmd.session
}

type SetSessionOption interface {
	ApplyOnSetSessionCommand(cmd *SetSessionCommand)
}

func NewSetSessionCommand(instanceID, sessionID string, opts ...SetSessionOption) *SetSessionCommand {
	cmd := &SetSessionCommand{
		instanceID: instanceID,
		sessionID:  sessionID,
		session: lazyGetter[*Session]{
			get: func(ctx context.Context, opts *InvokeOpts) (*Session, error) {
				return opts.sessionRepo.Get(ctx, opts.DB(),
					database.WithCondition(opts.sessionRepo.PrimaryKeyCondition(instanceID, sessionID)),
				)
			},
		},
	}

	for _, opt := range opts {
		opt.ApplyOnSetSessionCommand(cmd)
	}

	return cmd
}

// Events implements [Commander].
func (cmd *SetSessionCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	fetchedSession, err := cmd.FetchSession(ctx, opts)
	if err != nil {
		return nil, err
	}
	aggregate := &session.NewAggregate(fetchedSession.ID, fetchedSession.InstanceID).Aggregate
	commands := []eventstore.Command{
		session.NewTokenSetEvent(ctx, aggregate, fetchedSession.TokenID),
	}
	if cmd.lifetime != nil {
		commands = append(commands, session.NewLifetimeSetEvent(ctx, aggregate, *cmd.lifetime))
	}

	activity.TriggerWithoutOrg(ctx, fetchedSession.UserID, activity.SessionAPI)

	if len(cmd.metadata) == 0 {
		return commands, nil
	}
	metadata := make(map[string][]byte, len(cmd.metadata))
	for _, md := range cmd.metadata {
		metadata[md.Key] = md.Value
	}
	// the following logic is copied from the `ChangeMetadata`-method in https://github.com/zitadel/zitadel/blob/main/internal/command/session.go
	var changed bool
	for key, value := range metadata {
		idx := slices.IndexFunc(fetchedSession.Metadata, func(m *SessionMetadata) bool {
			return m.Key == key
		})

		if len(value) != 0 {
			// if a value is provided, and it's not equal, change it
			if !bytes.Equal(fetchedSession.Metadata[idx].Value, value) {
				fetchedSession.Metadata[idx].Value = value
				changed = true
			}
		} else {
			// if there's no / an empty value, we only need to remove it on existing entries
			if idx != -1 {
				delete(metadata, key)
				changed = true
			}
		}
	}
	if changed {
		commands = append(commands, session.NewMetadataSetEvent(ctx, aggregate, metadata))
	}

	return commands, nil
}

// Execute implements [Commander].
func (cmd *SetSessionCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	changes := make(database.Changes, 0, 3)
	if cmd.lifetime != nil {
		changes = append(changes, opts.sessionRepo.SetLifetime(*cmd.lifetime))
	}
	if len(cmd.metadata) > 0 {
		changes = append(changes, opts.sessionRepo.SetMetadata(cmd.metadata))
	}

	for _, check := range cmd.checks {
		changes = append(changes, opts.sessionRepo.SetFactor(check.checkResult()))
	}
	for _, challenge := range cmd.challenges {
		changes = append(changes, opts.sessionRepo.SetChallenge(challenge.challengeResult()))
	}
	_, err = opts.sessionRepo.Update(ctx, opts.DB(), opts.sessionRepo.PrimaryKeyCondition(cmd.instanceID, cmd.sessionID), changes...)
	return err
}

// String implements [Commander].
func (cmd *SetSessionCommand) String() string {
	return "SetSessionCommand"
}

// Validate implements [Commander].
func (cmd *SetSessionCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	_, err = cmd.FetchSession(ctx, opts)
	return err
}

// SetUserConditionProvider implements [CheckUserParent].
func (cmd *SetSessionCommand) SetUserConditionProvider(provider UserConditionProvider) {
	cmd.user = &lazyGetter[*User]{
		get: func(ctx context.Context, opts *InvokeOpts) (*User, error) {
			return opts.userRepo.Get(ctx, opts.DB(), database.WithCondition(database.And(
				opts.userRepo.InstanceIDCondition(cmd.instanceID),
				provider(ctx, opts),
			)))
		},
	}
}

// FetchSession implements [CheckPasswordParent] and [CheckUserParent].
func (cmd *SetSessionCommand) FetchSession(ctx context.Context, opts *InvokeOpts) (session *Session, err error) {
	return cmd.session.fetch(ctx, opts)
}

// FetchUser implements [CheckUserParent].
func (cmd *SetSessionCommand) FetchUser(ctx context.Context, opts *InvokeOpts) (user *User, err error) {
	fetchedSession, err := cmd.FetchSession(ctx, opts)
	if err != nil {
		return nil, err
	}
	if cmd.user == nil && fetchedSession.UserID == "" {
		return nil, zerrors.ThrowNotFound(nil, "DOMAI-mY4Y2", "no user information provided")
	}
	if cmd.user == nil {
		cmd.user = &lazyGetter[*User]{
			get: func(ctx context.Context, opts *InvokeOpts) (*User, error) {
				return opts.userRepo.Get(ctx, opts.DB(), database.WithCondition(
					opts.userRepo.PrimaryKeyCondition(cmd.instanceID, fetchedSession.UserID),
				))
			},
		}
	}
	fetchedUser, err := cmd.user.fetch(ctx, opts)
	if err != nil {
		return nil, err
	}
	if fetchedSession.UserID != "" && fetchedUser.ID != fetchedSession.UserID {
		return nil, zerrors.ThrowInternal(nil, "DOMAI-2g4f2", "user information does not match session")
	}
	return fetchedUser, nil
}

// reloadUser implements [CheckUserParent].
func (cmd *SetSessionCommand) reloadUser(ctx context.Context, opts *InvokeOpts) (user *User, err error) {
	return cmd.user.reload(ctx, opts)
}

var (
	_ Commander           = (*SetSessionCommand)(nil)
	_ CheckUserParent     = (*SetSessionCommand)(nil)
	_ CheckPasswordParent = (*SetSessionCommand)(nil)
	_ Querier[*Session]   = (*SetSessionCommand)(nil)
)
