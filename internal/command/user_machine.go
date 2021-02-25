package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddMachine(ctx context.Context, orgID string, machine *domain.Machine) (*domain.Machine, error) {
	if !machine.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-bm9Ds", "Errors.User.Invalid")
	}

	userID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	machine.AggregateID = userID

	orgIAMPolicy, err := c.getOrgIAMPolicy(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if !orgIAMPolicy.UserLoginMustBeDomain {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-6M0ds", "Errors.User.Invalid")
	}

	addedMachine := NewMachineWriteModel(machine.AggregateID, orgID)
	userAgg := UserAggregateFromWriteModel(&addedMachine.WriteModel)
	events, err := c.eventstore.PushEvents(ctx, user.NewMachineAddedEvent(
		ctx,
		userAgg,
		machine.Username,
		machine.Name,
		machine.Description,
		orgIAMPolicy.UserLoginMustBeDomain,
	))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedMachine, events...)
	if err != nil {
		return nil, err
	}
	return writeModelToMachine(addedMachine), nil
}

func (c *Commands) ChangeMachine(ctx context.Context, machine *domain.Machine) (*domain.Machine, error) {
	existingMachine, err := c.machineWriteModelByID(ctx, machine.AggregateID, machine.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingMachine.UserState) {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-5M0od", "Errors.User.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingMachine.WriteModel)
	changedEvent, hasChanged := existingMachine.NewChangedEvent(ctx, userAgg, machine.Name, machine.Description)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2n8vs", "Errors.User.NotChanged")
	}

	events, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMachine, events...)
	if err != nil {
		return nil, err
	}
	return writeModelToMachine(existingMachine), nil
}

//TODO: adlerhurst we should check userID on the same level, in user.go userID is checked in public funcs
func (c *Commands) machineWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *MachineWriteModel, err error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-0Plof", "Errors.User.UserIDMissing")
	}
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewMachineWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
