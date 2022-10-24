package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type Machine struct {
	models.ObjectRoot

	Username    string
	Name        string
	Description string
}

func (m *Machine) content() error {
	if m.ResourceOwner == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-xiown2", "Errors.ResourceOwnerMissing")
	}
	if m.AggregateID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-p0p2mi", "Errors.User.UserIDMissing")
	}
	/* not necessary for change
	if m.Username == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-bm9Ds", "Errors.User.Invalid")
	}*/
	if m.Name == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-bs9Ds", "Errors.User.Invalid")
	}
	return nil
}

func (c *Commands) AddMachine(ctx context.Context, machine *Machine) (*domain.ObjectDetails, error) {
	if machine.AggregateID == "" {
		userID, err := c.idGenerator.Next()
		if err != nil {
			return nil, err
		}
		machine.AggregateID = userID
	}
	domainPolicy, err := c.getOrgDomainPolicy(ctx, machine.ResourceOwner)
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-3M9fs", "Errors.Org.DomainPolicy.NotFound")
	}

	validation := prepareAddUserMachine(machine, domainPolicy)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func prepareAddUserMachine(machine *Machine, domainPolicy *domain.DomainPolicy) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if err := machine.content(); err != nil {
			return nil, err
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getMachineWriteModelByID(ctx, filter, machine.AggregateID, machine.ResourceOwner)
			if err != nil {
				return nil, err
			}
			if isUserStateExists(writeModel.UserState) {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-k2una", "Errors.User.AlreadyExisting")
			}
			return []eventstore.Command{
				user.NewMachineAddedEvent(
					ctx,
					UserAggregateFromWriteModel(&writeModel.WriteModel),
					machine.Username,
					machine.Name,
					machine.Description,
					domainPolicy.UserLoginMustBeDomain,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) ChangeMachine(ctx context.Context, machine *Machine) (*domain.ObjectDetails, error) {
	validation := prepareChangeUserMachine(machine)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func prepareChangeUserMachine(machine *Machine) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if err := machine.content(); err != nil {
			return nil, err
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getMachineWriteModelByID(ctx, filter, machine.AggregateID, machine.ResourceOwner)
			if err != nil {
				return nil, err
			}
			if !isUserStateExists(writeModel.UserState) {
				return nil, caos_errs.ThrowNotFound(nil, "COMMAND-5M0od", "Errors.User.NotFound")
			}

			event, hasChanged, err := writeModel.NewChangedEvent(ctx, UserAggregateFromWriteModel(&writeModel.WriteModel), machine.Name, machine.Description)
			if err != nil {
				return nil, err
			}
			if !hasChanged {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2n8vs", "Errors.User.NotChanged")
			}

			return []eventstore.Command{
				event,
			}, nil
		}, nil
	}
}

func getMachineWriteModelByID(ctx context.Context, filter preparation.FilterToQueryReducer, userID, resourceOwner string) (_ *MachineWriteModel, err error) {
	writeModel := NewMachineWriteModel(userID, resourceOwner)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}
