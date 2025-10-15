package domain

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type UpdateInstanceCommand struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Events implements Commander.
func (u *UpdateInstanceCommand) Events(ctx context.Context, _ *InvokeOpts) ([]eventstore.Command, error) {
	return []eventstore.Command{
		instance.NewInstanceChangedEvent(ctx, &instance.NewAggregate(u.ID).Aggregate, u.Name),
	}, nil
}

var _ Commander = (*UpdateInstanceCommand)(nil)

func NewUpdateInstanceCommand(id, name string) *UpdateInstanceCommand {
	return &UpdateInstanceCommand{
		ID:   id,
		Name: name,
	}
}

// Execute implements [Commander]
func (u *UpdateInstanceCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()

	instanceRepo := opts.instanceRepo

	updateCount, err := instanceRepo.Update(
		ctx,
		pool,
		u.ID,
		database.NewChange(instanceRepo.NameColumn(), u.Name),
	)
	if err != nil {
		return err
	}

	if updateCount == 0 {
		err = zerrors.ThrowNotFound(nil, "DOM-ghfov1", "Errors.Instance.NotFound")
		return err
	}
	if updateCount > 1 {
		err = zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, updateCount), "DOM-HlrNmD", "unexpected number of rows updated")
		return err
	}

	return err
}

// String implements [Commander]
func (u *UpdateInstanceCommand) String() string {
	return "UpdateInstanceCommand"
}

// Validate implements [Commander]
func (u *UpdateInstanceCommand) Validate(ctx context.Context, opts *InvokeOpts) error {
	if u.ID = strings.TrimSpace(u.ID); u.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-wSs6kG", "invalid instance ID")
	}
	if u.Name = strings.TrimSpace(u.Name); u.Name == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-FPJcLC", "invalid instance name")
	}

	if authZErr := opts.Permissions.CheckInstancePermission(ctx, InstanceWritePermission); authZErr != nil {
		return zerrors.ThrowPermissionDenied(authZErr, "DOM-M5ObLP", "permission denied")
	}

	instanceRepo := opts.instanceRepo

	instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.IDCondition(u.ID)))
	if err != nil {
		return err
	}

	if instance.Name == u.Name {
		err = zerrors.ThrowPreconditionFailed(nil, "DOM-5MrT21", "Errors.Instance.NotChanged")
		return err
	}
	return nil
}
