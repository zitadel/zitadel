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
	ID              string
	ResourceOwner   string
	Username        *string
	Name            *string
	Description     *string
	AccessTokenType *domain.OIDCTokenType

	// Details are set after a successful execution of the command
	Details *domain.ObjectDetails

	Metadata []*domain.Metadata
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
	if h.AccessTokenType != nil {
		return true
	}
	if len(h.Metadata) > 0 {
		return true
	}
	return false
}

func (c *Commands) ChangeUserMachine(ctx context.Context, machine *ChangeMachine) (err error) {
	existingMachine, err := c.UserMachineWriteModel(
		ctx,
		machine.ID,
		machine.ResourceOwner,
		len(machine.Metadata) > 0,
	)
	if err != nil {
		return err
	}
	if machine.Changed() {
		if err := c.checkPermissionUpdateUser(ctx, existingMachine.ResourceOwner, existingMachine.AggregateID, true); err != nil {
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
	if machine.AccessTokenType != nil && *machine.AccessTokenType != existingMachine.AccessTokenType {
		machineChanges = append(machineChanges, user.ChangeAccessTokenType(*machine.AccessTokenType))
	}
	if len(machineChanges) > 0 {
		cmds = append(cmds, user.NewMachineChangedEvent(ctx, &existingMachine.Aggregate().Aggregate, machineChanges))
	}
	metadataCmds, err := c.updateUserMetadata(ctx, machine.Metadata, &existingMachine.Aggregate().Aggregate)
	if err != nil {
		return err
	}
	cmds = append(cmds, metadataCmds...)
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

func (c *Commands) updateUserMetadata(ctx context.Context, metadata []*domain.Metadata, aggregate *eventstore.Aggregate) ([]eventstore.Command, error) {
	if len(metadata) == 0 {
		return nil, nil
	}
	cmds := make([]eventstore.Command, 0, len(metadata))
	for _, md := range metadata {
		if md == nil {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-uAFkgS", "Errors.Metadata.Invalid")
		}
		// remove metadata if the value is empty
		if len(md.Value) == 0 {
			cmd, err := c.removeUserMetadata(ctx, aggregate, md.Key)
			if err != nil {
				return nil, err
			}
			cmds = append(cmds, cmd)
			continue
		}

		cmd, err := c.setUserMetadata(ctx, aggregate, md)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}
