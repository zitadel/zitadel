package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ChangeMachine struct {
	ID            string
	ResourceOwner string
	Username      *string
	Name          *string
	Description   *string

	// Details are set after a successful execution of the command
	Details *domain.ObjectDetails
}

func (h *ChangeMachine) Changed() bool {
	if h.Username != nil {
		return true
	}
	if h.Name != nil {
		return true
	}
	if h.Description != nil {
		return true
	}
	return false
}

func (c *Commands) ChangeUserMachine(ctx context.Context, machine *ChangeMachine) (err error) {
	existingMachine, err := c.UserMachineWriteModel(
		ctx,
		machine.ID,
		machine.ResourceOwner,
		false,
	)
	if err != nil {
		return err
	}
	if machine.Changed() {
		if err := c.checkPermissionUpdateUser(ctx, existingMachine.ResourceOwner, existingMachine.AggregateID); err != nil {
			return err
		}
	}

	cmds := make([]eventstore.Command, 0)
	if machine.Username != nil {
		cmds, err = c.changeUsername(ctx, cmds, existingMachine, *machine.Username)
		if err != nil {
			return err
		}
	}
	var machineChanges []user.MachineChanges
	if machine.Name != nil && *machine.Name != existingMachine.Name {
		machineChanges = append(machineChanges, user.ChangeName(*machine.Name))
	}
	if machine.Description != nil && *machine.Description != existingMachine.Description {
		machineChanges = append(machineChanges, user.ChangeDescription(*machine.Description))
	}
	if len(machineChanges) > 0 {
		cmds = append(cmds, user.NewMachineChangedEvent(ctx, &existingMachine.Aggregate().Aggregate, machineChanges))
	}
	if len(cmds) == 0 {
		machine.Details = writeModelToObjectDetails(&existingMachine.WriteModel)
		return nil
	}
	err = c.pushAppendAndReduce(ctx, existingMachine, cmds...)
	if err != nil {
		return err
	}
	machine.Details = writeModelToObjectDetails(&existingMachine.WriteModel)
	return nil
}

func (c *Commands) UserMachineWriteModel(ctx context.Context, userID, resourceOwner string, metadataWM bool) (writeModel *UserV2WriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	writeModel = NewUserMachineWriteModel(userID, resourceOwner, metadataWM)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(writeModel.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-ugjs0upun6", "Errors.User.NotFound")
	}
	return writeModel, nil
}
