package domain

import (
	"context"
	"errors"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type DeleteSessionCommand struct {
	ID                  string `json:"id"`
	Token               string `json:"token"`
	MustCheckPermission bool   `json:"mustCheckPermission,omitempty"`

	// DeletedAt is set after successful execution.
	DeletedAt time.Time `json:"deletedAt,omitempty"`
}

func NewDeleteSessionCommand(
	id string,
	token string,
	mustCheckPermission bool,
) *DeleteSessionCommand {
	return &DeleteSessionCommand{
		ID:                  id,
		Token:               token,
		MustCheckPermission: mustCheckPermission,
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
	instance := authz.GetInstance(ctx)

	deletedRows, err := sessionRepo.Delete(ctx, opts.DB(),
		sessionRepo.PrimaryKeyCondition(instance.InstanceID(), cmd.ID),
	)
	if err != nil {
		return err
	}

	if deletedRows > 1 {
		return zerrors.ThrowInternalf(nil, "DOM-wv33rsKpRw", "expecting 1 row deleted, got %d", deletedRows)
	}

	if deletedRows == 1 {
		// TODO(LS): Change this with the real update date when SessionRepo.Delete()
		// returns the timestamp. See https://github.com/zitadel/zitadel/issues/10881
		cmd.DeletedAt = time.Now()
	}
	return nil
}

func (cmd *DeleteSessionCommand) checkPermission(ctx context.Context, session *Session, opts *InvokeOpts) error {
	var id, tokenID, userID, userResourceOwner string
	if session != nil {
		id = session.ID
		tokenID = session.TokenID
		userID = session.UserID
	}
	if cmd.Token != "" {
		return opts.sessionTokenVerifier(ctx, cmd.Token, id, tokenID)
	}
	if userID != "" && userID == authz.GetCtxData(ctx).UserID {
		return nil
	}
	return opts.Permissions.CheckOrganizationPermission(ctx, domain.PermissionSessionDelete, userResourceOwner)
}

// String implements [Commander].
func (*DeleteSessionCommand) String() string {
	return "DeleteSessionCommand"
}

// Validate implements [Commander].
func (cmd *DeleteSessionCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if cmd.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "SESS-3n9fs", "Errors.IDMissing")
	}
	if !cmd.MustCheckPermission {
		return nil
	}
	sessionRepo := opts.sessionRepo
	instance := authz.GetInstance(ctx)

	sessionToDelete, err := sessionRepo.Get(ctx, opts.DB(),
		database.WithCondition(
			sessionRepo.PrimaryKeyCondition(instance.InstanceID(), cmd.ID),
		),
	)
	if err != nil {
		if !errors.Is(err, &database.NoRowFoundError{}) {
			return err
		}
	}
	return cmd.checkPermission(ctx, sessionToDelete, opts)
}

var (
	_ Commander     = (*DeleteOrgCommand)(nil)
	_ Transactional = (*DeleteOrgCommand)(nil)
)
