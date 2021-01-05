package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
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
