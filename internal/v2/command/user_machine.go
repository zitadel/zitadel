package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) AddMachine(ctx context.Context, orgID, username string, machine *domain.Machine) (*domain.Machine, error) {
	if !machine.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-5M0ds", "Errors.User.Invalid")
	}
	userID, err := r.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	//TODO: Check Unique username
	machine.AggregateID = userID
	orgIAMPolicy, err := r.GetOrgIAMPolicy(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if !orgIAMPolicy.UserLoginMustBeDomain {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-6M0ds", "Errors.User.Invalid")
	}

	addedMachine := NewMachineWriteModel(machine.AggregateID)
	userAgg := UserAggregateFromWriteModel(&addedMachine.WriteModel)
	userAgg.PushEvents(
		user.NewMachineAddedEvent(
			ctx,
			username,
			machine.Name,
			machine.Description,
		),
	)
	return writeModelToMachine(addedMachine), nil
}

func (r *CommandSide) ChangeMachine(ctx context.Context, machine *domain.Machine) (*domain.Machine, error) {
	existingUser, err := r.machineWriteModelByID(ctx, machine.AggregateID)
	if err != nil {
		return nil, err
	}
	if existingUser.UserState == domain.UserStateDeleted || existingUser.UserState == domain.UserStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-5M0od", "Errors.User.NotFound")
	}

	changedEvent, hasChanged := existingUser.NewChangedEvent(ctx, machine.Name, machine.Description)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9fs", "Errors.User.Email.NotChanged")
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	userAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingUser, userAgg)
	if err != nil {
		return nil, err
	}
	return writeModelToMachine(existingUser), nil
}

func (r *CommandSide) machineWriteModelByID(ctx context.Context, userID string) (writeModel *MachineWriteModel, err error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-5M0ds", "Errors.User.UserIDMissing")
	}
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewMachineWriteModel(userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
