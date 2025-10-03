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

	InstanceDomains []string
	InstanceName    string
}

func NewDeleteInstanceCommand(instanceID string) *DeleteInstanceCommand {
	return &DeleteInstanceCommand{ID: instanceID}
}

// Events implements Commander.
func (d *DeleteInstanceCommand) Events(ctx context.Context, opts *CommandOpts) ([]eventstore.Command, error) {
	milestoneAggregate := milestone.NewInstanceAggregate(d.ID)
	instanceAggregate := instance.NewAggregate(d.ID)

	return []eventstore.Command{
		instance.NewInstanceRemovedEvent(ctx, &instanceAggregate.Aggregate, d.InstanceName, d.InstanceDomains),
		milestone.NewReachedEvent(ctx, milestoneAggregate, milestone.InstanceDeleted),
	}, nil
}

// Execute implements Commander.
func (d *DeleteInstanceCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	closeFunc, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}

	defer func() { err = closeFunc(ctx, err) }()

	instanceRepo := opts.instanceRepo.LoadDomains()

	instanceToDelete, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.IDCondition(d.ID)))
	if err != nil {
		return err
	}

	d.InstanceDomains = make([]string, len(instanceToDelete.Domains))
	for i, domain := range instanceToDelete.Domains {
		d.InstanceDomains[i] = domain.Domain
	}

	d.InstanceName = instanceToDelete.Name

	deletedRows, err := instanceRepo.Delete(ctx, pool, instanceToDelete.ID)
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

// String implements Commander.
func (d *DeleteInstanceCommand) String() string {
	return "DeleteInstanceCommand"
}

// Validate implements Commander.
func (d *DeleteInstanceCommand) Validate(ctx context.Context, opts *CommandOpts) (err error) {
	if strings.TrimSpace(d.ID) == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-VpQ9lF", "Errors.Invalid.Argument")
	}

	return nil
}

var _ Commander = (*DeleteInstanceCommand)(nil)
