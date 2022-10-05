package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type AddMachine struct {
	Machine                  *domain.Machine
	Pat                      bool
	PatExpirationDate        time.Time
	PatScopes                []string
	MachineKey               bool
	MachineKeyType           domain.AuthNKeyType
	MachineKeyExpirationDate time.Time
}

func AddMachineCommand(a *user.Aggregate, machine *domain.Machine) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if !machine.IsValid() {
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
			domainPolicy, err := domainPolicyWriteModel(ctx, filter)
			if err != nil {
				return nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-3M9fs", "Errors.Org.DomainPolicy.NotFound")
			}
			if !domainPolicy.UserLoginMustBeDomain {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M0dd", "Errors.User.Invalid")
			}

			return []eventstore.Command{
				user.NewMachineAddedEvent(ctx, &a.Aggregate, machine.Username, machine.Name, machine.Description, domainPolicy.UserLoginMustBeDomain),
			}, nil
		}, nil
	}
}

func (c *Commands) AddMachine(ctx context.Context, orgID string, machine *domain.Machine) (*domain.Machine, error) {
	userID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	return c.addMachineWithID(ctx, orgID, userID, machine)
}

func (c *Commands) AddMachineWithID(ctx context.Context, orgID string, userID string, machine *domain.Machine) (*domain.Machine, error) {
	return c.addMachineWithID(ctx, orgID, userID, machine)
}

func (c *Commands) addMachineWithID(ctx context.Context, orgID string, userID string, machine *domain.Machine) (*domain.Machine, error) {
	agg := user.NewAggregate(userID, orgID)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, AddMachineCommand(agg, machine))
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}

	return &domain.Machine{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   agg.ID,
			Sequence:      events[len(events)-1].Sequence(),
			CreationDate:  events[len(events)-1].CreationDate(),
			ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
			InstanceID:    events[len(events)-1].Aggregate().InstanceID,
		},
		Username:    machine.Username,
		Name:        machine.Name,
		Description: machine.Description,
		State:       machine.State,
	}, nil
}

func (c *Commands) ChangeMachine(ctx context.Context, machine *domain.Machine) (*domain.Machine, error) {
	agg := user.NewAggregate(machine.AggregateID, machine.ResourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, changeMachineCommand(agg, machine))
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}

	return &domain.Machine{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   agg.ID,
			Sequence:      events[len(events)-1].Sequence(),
			CreationDate:  events[len(events)-1].CreationDate(),
			ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
			InstanceID:    events[len(events)-1].Aggregate().InstanceID,
		},
		Username:    machine.Username,
		Name:        machine.Name,
		Description: machine.Description,
	}, nil
}

func changeMachineCommand(a *user.Aggregate, machine *domain.Machine) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if !machine.IsValid() {
			return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-bm9Ds", "Errors.User.Invalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getMachineWriteModel(ctx, a.ID, a.ResourceOwner, filter)
			if err != nil {
				return nil, err
			}
			if !isUserStateExists(writeModel.UserState) {
				return nil, caos_errs.ThrowNotFound(nil, "COMMAND-5M0od", "Errors.User.NotFound")
			}

			domainPolicy, err := domainPolicyWriteModel(ctx, filter)
			if err != nil {
				return nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-3M9fs", "Errors.Org.DomainPolicy.NotFound")
			}
			if !domainPolicy.UserLoginMustBeDomain {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M0dd", "Errors.User.Invalid")
			}

			changedEvent, hasChanged, err := writeModel.NewChangedEvent(ctx, &a.Aggregate, machine.Name, machine.Description)
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
