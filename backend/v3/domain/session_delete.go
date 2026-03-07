package domain

import (
	"context"
	"strings"
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
	if cmd.DeletedAt.IsZero() {
		return nil, nil
	}
	return []eventstore.Command{
		session.NewTerminateEvent(ctx, &session.NewAggregate(cmd.ID, authz.GetInstance(ctx).InstanceID()).Aggregate),
	}, nil
}

func sessionDeletePermissionCheckCondition(ctx context.Context, sessionRepo SessionRepository, id, token string, decryptor SessionTokenDecryptor) (database.Condition, error) {
	if token != "" {
		sessionID, tokenID, err := decryptor(ctx, token)
		if err != nil || sessionID != id {
			return nil, zerrors.ThrowInvalidArgumentf(err, "SESS-S3gq1", "Errors.Session.TokenInvalid")
		}
		return database.Or(database.Exists("sessions", sessionRepo.TokenIDCondition(tokenID)),
			database.Permission(domain.PermissionSessionDelete, true),
		), nil
	}
	return database.Or(
		database.Exists("sessions", sessionRepo.UserIDCondition(authz.GetCtxData(ctx).UserID)),
		database.Permission(domain.PermissionSessionDelete, true), // TODO: implement check
	), nil
}

// Execute implements [Commander].
func (cmd *DeleteSessionCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	sessionRepo := opts.sessionRepo
	instance := authz.GetInstance(ctx)

	permCheck, err := sessionDeletePermissionCheckCondition(ctx, sessionRepo, cmd.ID, cmd.Token, opts.sessionTokenDecryptor)
	if err != nil {
		return err
	}

	deletedRows, err := sessionRepo.Delete(ctx, opts.DB(),
		sessionRepo.PrimaryKeyCondition(instance.InstanceID(), cmd.ID),
		permCheck,
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
	if cmd.ID = strings.TrimSpace(cmd.ID); cmd.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "SESS-3n9fs", "Errors.IDMissing")
	}
	return nil
}

var (
	_ Commander     = (*DeleteOrgCommand)(nil)
	_ Transactional = (*DeleteOrgCommand)(nil)
)
