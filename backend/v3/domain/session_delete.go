package domain

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type DeleteSessionCommand struct {
	ID                   string  `json:"id"`
	Token                *string `json:"token"`
	SessionTokenVerifier func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)
}

func NewDeleteSessionCommand(
	id string,
	token *string,
	sessionTokenVerifier func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error),
) *DeleteSessionCommand {
	return &DeleteSessionCommand{
		ID:                   id,
		Token:                token,
		SessionTokenVerifier: sessionTokenVerifier,
	}
}

// RequiresTransaction implements [Transactional].
func (cmd *DeleteSessionCommand) RequiresTransaction() {}

// Events implements [Commander].
func (cmd *DeleteSessionCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	return []eventstore.Command{
		session.NewTerminateEvent(ctx, &session.NewAggregate(cmd.ID, authz.GetInstance(ctx).InstanceID()).Aggregate),
	}, nil
}

// Execute implements [Commander].
func (cmd *DeleteSessionCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	sessionRepo := opts.sessionRepo

	deletedRows, err := sessionRepo.Delete(ctx, opts.DB(),
		sessionRepo.PrimaryKeyCondition(authz.GetInstance(ctx).InstanceID(), cmd.ID),
	)

	if deletedRows > 1 {
		err = zerrors.ThrowInternalf(nil, "DOM-wv33rsKpRw", "expecting 1 row deleted, got %d", deletedRows)
		return err
	}

	if deletedRows < 1 {
		err = zerrors.ThrowNotFound(nil, "DOM-g1lDb1qs1f", "session not found")
	}
	return err
}

// String implements [Commander].
func (DeleteSessionCommand) String() string {
	return "DeleteSessionCommand"
}

// Validate implements [Commander].
func (cmd *DeleteSessionCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	sessionRepo := opts.sessionRepo

	sessionToDelete, errGetSession := sessionRepo.Get(ctx, opts.DB(), database.WithCondition(
		sessionRepo.PrimaryKeyCondition(authz.GetInstance(ctx).InstanceID(), cmd.ID),
	))
	if errGetSession != nil {
		if errors.Is(errGetSession, &database.NoRowFoundError{}) {
			return zerrors.ThrowNotFound(errGetSession, "DOM-8KYOH3", "Errors.Session.NotFound")
		}
	}
	if cmd.Token != nil {
		if err := cmd.SessionTokenVerifier(ctx, *cmd.Token, cmd.ID, sessionToDelete.TokenID); err != nil {
			return err
		}
	}
	return err
}

var (
	_ Commander     = (*DeleteOrgCommand)(nil)
	_ Transactional = (*DeleteOrgCommand)(nil)
)
