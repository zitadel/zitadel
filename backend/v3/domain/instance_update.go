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
	ID               string `json:"id"`
	Name             string `json:"name"`
	ShouldSkipUpdate bool   `json:"should_skip_update"`
}

// RequiresTransaction implements [Transactional].
func (u *UpdateInstanceCommand) RequiresTransaction() {}

// Events implements Commander.
func (u *UpdateInstanceCommand) Events(ctx context.Context, _ *InvokeOpts) ([]eventstore.Command, error) {
	if u.ShouldSkipUpdate {
		return nil, nil
	}

	return []eventstore.Command{
		instance.NewInstanceChangedEvent(ctx, &instance.NewAggregate(u.ID).Aggregate, u.Name),
	}, nil
}

var (
	_ Commander     = (*UpdateInstanceCommand)(nil)
	_ Transactional = (*UpdateInstanceCommand)(nil)
)

func NewUpdateInstanceCommand(id, name string) *UpdateInstanceCommand {
	return &UpdateInstanceCommand{
		ID:   id,
		Name: name,
	}
}

// Execute implements [Commander]
func (u *UpdateInstanceCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if u.ShouldSkipUpdate {
		return
	}

	instanceRepo := opts.instanceRepo

	updateCount, err := instanceRepo.Update(
		ctx,
		opts.DB(),
		u.ID,
		database.NewChange(instanceRepo.NameColumn(), u.Name),
	)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-PkVMNR", "Errors.Instance.Update")
	}

	if updateCount == 0 {
		return zerrors.ThrowNotFound(nil, "DOM-ghfov1", "Errors.Instance.NotFound")
	}
	if updateCount > 1 {
		return zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, updateCount), "DOM-HlrNmD", "Errors.Instance.UpdateMismatch")
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
		return zerrors.ThrowInvalidArgument(nil, "DOM-wSs6kG", "Errors.Instance.ID")
	}
	if u.Name = strings.TrimSpace(u.Name); u.Name == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-FPJcLC", "Errors.Instance.Name")
	}

	if authZErr := opts.Permissions.CheckInstancePermission(ctx, InstanceWritePermission); authZErr != nil {
		return zerrors.ThrowPermissionDenied(authZErr, "DOM-M5ObLP", "Errors.PermissionDenied")
	}

	instanceRepo := opts.instanceRepo

	instance, err := instanceRepo.Get(ctx, opts.DB(), database.WithCondition(instanceRepo.IDCondition(u.ID)))
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-j05Hdo", "Errors.Instance.Get")
	}

	u.ShouldSkipUpdate = instance.Name == u.Name

	return nil
}
