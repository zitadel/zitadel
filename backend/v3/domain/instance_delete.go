package domain

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type DeleteInstanceCommand struct {
	ID string `json:"id"`

	// InstanceDomains is public for testing purposes
	// do not use this field
	InstanceDomains []string
	// InstanceName is public for testing purposes
	// do not use this field
	InstanceName string
}

// RequiresTransaction implements [Transactional].
func (d *DeleteInstanceCommand) RequiresTransaction() {}

func NewDeleteInstanceCommand(instanceID string) *DeleteInstanceCommand {
	return &DeleteInstanceCommand{ID: instanceID}
}

// Events implements [Commander].
func (d *DeleteInstanceCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	milestoneAggregate := milestone.NewInstanceAggregate(d.ID)
	instanceAggregate := instance.NewAggregate(d.ID)

	return []eventstore.Command{
		instance.NewInstanceRemovedEvent(ctx, &instanceAggregate.Aggregate, d.InstanceName, d.InstanceDomains),
		milestone.NewReachedEvent(ctx, milestoneAggregate, milestone.InstanceDeleted),
	}, nil
}

// Execute implements [Commander].
func (d *DeleteInstanceCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	instanceRepo := opts.instanceRepo.LoadDomains()

	instanceToDelete, err := instanceRepo.Get(ctx, opts.DB(), database.WithCondition(instanceRepo.IDCondition(d.ID)))
	if err != nil {
		return err
	}

	d.InstanceDomains = make([]string, len(instanceToDelete.Domains))
	for i, domain := range instanceToDelete.Domains {
		d.InstanceDomains[i] = domain.Domain
	}

	d.InstanceName = instanceToDelete.Name

	deletedRows, err := instanceRepo.Delete(ctx, opts.DB(), instanceToDelete.ID)
	if err != nil {
		return err
	}

	if deletedRows > 1 {
		err = zerrors.ThrowInternalf(nil, "DOM-Od04Jx", "expecting 1 row deleted, got %d", deletedRows)
		return err
	}

	if deletedRows < 1 {
		err = zerrors.ThrowNotFound(nil, "DOM-daglwD", "instance not found")
	}

	return err
}

// String implements [Commander].
func (d *DeleteInstanceCommand) String() string {
	return "DeleteInstanceCommand"
}

// Validate implements [Commander].
func (d *DeleteInstanceCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if d.ID = strings.TrimSpace(d.ID); d.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-VpQ9lF", "Errors.Invalid.Argument")
	}

	if authZErr := opts.Permissions.CheckInstancePermission(ctx, InstanceWritePermission); authZErr != nil {
		return zerrors.ThrowPermissionDenied(authZErr, "DOM-Yz8f1X", "permission denied")
	}

	return nil
}

var (
	_ Commander     = (*DeleteInstanceCommand)(nil)
	_ Transactional = (*DeleteInstanceCommand)(nil)
)
