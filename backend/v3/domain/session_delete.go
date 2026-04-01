package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
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
	// deletedRows is set after successful execution, used to determine if the session was deleted or not.
	deletedRows int64
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
	if cmd.deletedRows == 0 {
		return nil, nil
	}
	return []eventstore.Command{
		session.NewTerminateEvent(ctx, &session.NewAggregate(cmd.ID, authz.GetInstance(ctx).InstanceID()).Aggregate),
	}, nil
}

var (
	SlugSessionTokenInvalid   = NewSlug("session", "token_invalid")
	SlugAuthMissingPermission = NewSlug("auth", "missing_permission")
	SlugRequestInvalid        = NewSlug("request", "invalid")
	SlugInternalError         = NewSlug("zitadel", "internal_error")
)

var (
	ErrSessionTokenInvalid = func(err error) error {
		return zerrors.ThrowPermissionDeniedSlug(err, SlugSessionTokenInvalid, "the provided session token is invalid, either the token is malformed, expired or does not match the session")
	}
	ErrAuthMissingPermission = func(err error, message string) error {
		return zerrors.ThrowPermissionDeniedSlug(err, SlugAuthMissingPermission, message)
	}
	ErrIDMissing         = zerrors.ThrowInvalidArgumentSlug(nil, SlugRequestInvalid, "validation failed: id is required", json.RawMessage{"id": "required"})
	ErrInstanceIDMissing = zerrors.ThrowInvalidArgumentSlug(nil, SlugRequestInvalid, "validation failed: instance_id is required", map[string]interface{}{"instance_id": "required"})
	ErrInternal          = func(err error, message string) error {
		return zerrors.ThrowInternalSlug(err, SlugInternalError, message, nil)
	}
	ErrMoreThanOneRowAffected = func(message string, rows int64) error {
		return zerrors.ThrowInternalSlug(nil, SlugInternalError, message, map[string]interface{}{"rows": rows})
	}
)

func (cmd *DeleteSessionCommand) sessionDeletePermissionCheckCondition(ctx context.Context, sessionRepo SessionRepository, decryptor SessionTokenDecryptor) (database.Condition, error) {
	if !cmd.MustCheckPermission {
		return nil, nil
	}
	if cmd.Token != "" {
		sessionID, tokenID, err := decryptor(ctx, cmd.Token)
		if err != nil || sessionID != cmd.ID {
			return nil, ErrSessionTokenInvalid(err)
		}
		return database.Or(
			sessionRepo.TokenIDCondition(tokenID),
			database.PermissionCheck(SessionDeletePermission, true),
		), nil
	}
	return database.Or(
		sessionRepo.UserIDCondition(authz.GetCtxData(ctx).UserID),
		database.PermissionCheck(SessionDeletePermission, true), // TODO: implement check once permissions are implemented
	), nil
}

// Execute implements [Commander].
func (cmd *DeleteSessionCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	sessionRepo := opts.sessionRepo
	instance := authz.GetInstance(ctx)

	permCheck, err := cmd.sessionDeletePermissionCheckCondition(ctx, sessionRepo, opts.sessionTokenDecryptor)
	if err != nil {
		return err
	}

	deletedRows, deletedAt, err := sessionRepo.Delete(ctx, opts.DB(),
		sessionRepo.PrimaryKeyCondition(instance.InstanceID(), cmd.ID),
		permCheck,
	)
	if err != nil {
		if errors.Is(err, new(database.PermissionError)) {
			return ErrAuthMissingPermission(err, "insufficient permissions to delete session, require `session.delete` permission, ownership of the session or current session token")
		}
		return ErrInternal(err, "an unexpected error occurred while deleting the session")
	}

	if deletedRows > 1 {
		return ErrMoreThanOneRowAffected(fmt.Sprintf("expected 1 session to be deleted, got %d", deletedRows), deletedRows)
	}

	cmd.DeletedAt = deletedAt
	cmd.deletedRows = deletedRows
	return nil
}

// String implements [Commander].
func (*DeleteSessionCommand) String() string {
	return "DeleteSessionCommand"
}

// Validate implements [Commander].
func (cmd *DeleteSessionCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if cmd.ID = strings.TrimSpace(cmd.ID); cmd.ID == "" {
		return ErrIDMissing
	}
	return nil
}

var (
	_ Commander     = (*DeleteSessionCommand)(nil)
	_ Transactional = (*DeleteSessionCommand)(nil)
)
