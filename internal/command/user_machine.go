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

type AddMachine struct {
	Machine    *Machine
	Pat        *AddPat
	MachineKey *AddMachineKey
}

type Machine struct {
	models.ObjectRoot

	Username        string
	Name            string
	Description     string
	AccessTokenType domain.OIDCTokenType
}

func (m *Machine) IsZero() bool {
	return m.Username == "" && m.Name == ""
}

func AddMachineCommand(a *user.Aggregate, machine *Machine) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if a.ResourceOwner == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-xiown2", "Errors.ResourceOwnerMissing")
		}
		if a.ID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-p0p2mi", "Errors.User.UserIDMissing")
		}
		if machine.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-bs9Ds", "Errors.User.Invalid")
		}
		if machine.Username == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-bm9Ds", "Errors.User.Invalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getMachineWriteModel(ctx, a.ID, a.ResourceOwner, filter)
			if err != nil {
				return nil, err
			}
			if isUserStateExists(writeModel.UserState) {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-k2una", "Errors.User.AlreadyExisting")
			}
			domainPolicy, err := domainPolicyWriteModel(ctx, filter, a.ResourceOwner)
			if err != nil {
				return nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-3M9fs", "Errors.Org.DomainPolicy.NotFound")
			}
			return []eventstore.Command{
				user.NewMachineAddedEvent(ctx, &a.Aggregate, machine.Username, machine.Name, machine.Description, domainPolicy.UserLoginMustBeDomain, machine.AccessTokenType),
			}, nil
		}, nil
	}
}

func (c *Commands) AddMachine(ctx context.Context, machine *Machine) (*domain.ObjectDetails, error) {
	if machine.AggregateID == "" {
		userID, err := c.idGenerator.Next()
		if err != nil {
			return nil, err
		}
		machine.AggregateID = userID
	}

	agg := user.NewAggregate(machine.AggregateID, machine.ResourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, AddMachineCommand(agg, machine))
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

func (c *Commands) ChangeMachine(ctx context.Context, machine *Machine) (*domain.ObjectDetails, error) {
	agg := user.NewAggregate(machine.AggregateID, machine.ResourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, changeMachineCommand(agg, machine))
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

func changeMachineCommand(a *user.Aggregate, machine *Machine) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if a.ResourceOwner == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-xiown3", "Errors.ResourceOwnerMissing")
		}
		if a.ID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-p0p3mi", "Errors.User.UserIDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getMachineWriteModel(ctx, a.ID, a.ResourceOwner, filter)
			if err != nil {
				return nil, err
			}
			if !isUserStateExists(writeModel.UserState) {
				return nil, caos_errs.ThrowNotFound(nil, "COMMAND-5M0od", "Errors.User.NotFound")
			}
			changedEvent, hasChanged, err := writeModel.NewChangedEvent(ctx, &a.Aggregate, machine.Name, machine.Description, machine.AccessTokenType)
			if err != nil {
				return nil, err
			}
			if !hasChanged {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2n8vs", "Errors.User.NotChanged")
			}

			return []eventstore.Command{
				changedEvent,
			}, nil
		}, nil
	}
}

func getMachineWriteModel(ctx context.Context, userID, resourceOwner string, filter preparation.FilterToQueryReducer) (*MachineWriteModel, error) {
	writeModel := NewMachineWriteModel(userID, resourceOwner)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModel, nil
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}
