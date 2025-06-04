package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	PermissionCheck PermissionCheck
}

func (m *Machine) IsZero() bool {
	return m.Username == "" && m.Name == ""
}

func AddMachineCommand(a *user.Aggregate, machine *Machine) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if a.ResourceOwner == "" && machine.PermissionCheck == nil {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-xiown3", "Errors.ResourceOwnerMissing")
		}
		if a.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-p0p2mi", "Errors.User.UserIDMissing")
		}
		if machine.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-bs9Ds", "Errors.User.Invalid")
		}
		if machine.Username == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-bm9Ds", "Errors.User.Invalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			ctx, span := tracing.NewSpan(ctx)
			defer func() { span.EndWithError(err) }()

			writeModel, err := getMachineWriteModel(ctx, a.ID, a.ResourceOwner, filter, machine.PermissionCheck)
			if err != nil {
				return nil, err
			}
			if isUserStateExists(writeModel.UserState) {
				return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-k2una", "Errors.User.AlreadyExisting")
			}
			domainPolicy, err := domainPolicyWriteModel(ctx, filter, a.ResourceOwner)
			if err != nil {
				return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-3M9fs", "Errors.Org.DomainPolicy.NotFound")
			}
			return []eventstore.Command{
				user.NewMachineAddedEvent(ctx, &a.Aggregate, machine.Username, machine.Name, machine.Description, domainPolicy.UserLoginMustBeDomain, machine.AccessTokenType),
			}, nil
		}, nil
	}
}

type addMachineOption func(context.Context, *Machine) error

func AddMachineWithUsernameToIDFallback() addMachineOption {
	return func(ctx context.Context, m *Machine) error {
		if m.Username == "" {
			m.Username = m.AggregateID
		}
		return nil
	}
}

func (c *Commands) AddMachine(ctx context.Context, machine *Machine, state *domain.UserState, check PermissionCheck, options ...addMachineOption) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if machine.AggregateID == "" {
		userID, err := c.idGenerator.Next()
		if err != nil {
			return nil, err
		}
		machine.AggregateID = userID
	}

	agg := user.NewAggregate(machine.AggregateID, machine.ResourceOwner)
	for _, option := range options {
		if err = option(ctx, machine); err != nil {
			return nil, err
		}
	}
	if check != nil {
		if err = check(machine.ResourceOwner, machine.AggregateID); err != nil {
			return nil, err
		}
	}
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, AddMachineCommand(agg, machine))
	if err != nil {
		return nil, err
	}

	if state != nil {
		var cmd eventstore.Command
		switch *state {
		case domain.UserStateInactive:
			cmd = user.NewUserDeactivatedEvent(ctx, &agg.Aggregate)
		case domain.UserStateLocked:
			cmd = user.NewUserLockedEvent(ctx, &agg.Aggregate)
		case domain.UserStateDeleted:
		// users are never imported if deleted
		case domain.UserStateActive:
		// added because of the linter
		case domain.UserStateSuspend:
		// added because of the linter
		case domain.UserStateInitial:
		// added because of the linter
		case domain.UserStateUnspecified:
			// added because of the linter
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}

	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

// Deprecated: use ChangeUserMachine instead
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
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func changeMachineCommand(a *user.Aggregate, machine *Machine) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if a.ResourceOwner == "" && machine.PermissionCheck == nil {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-xiown3", "Errors.ResourceOwnerMissing")
		}
		if a.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-p0p3mi", "Errors.User.UserIDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getMachineWriteModel(ctx, a.ID, a.ResourceOwner, filter, machine.PermissionCheck)
			if err != nil {
				return nil, err
			}
			if !isUserStateExists(writeModel.UserState) {
				return nil, zerrors.ThrowNotFound(nil, "COMMAND-5M0od", "Errors.User.NotFound")
			}
			changedEvent, hasChanged := writeModel.NewChangedEvent(ctx, &a.Aggregate, machine.Name, machine.Description, machine.AccessTokenType)
			if !hasChanged {
				return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-2n8vs", "Errors.User.NotChanged")
			}

			return []eventstore.Command{
				changedEvent,
			}, nil
		}, nil
	}
}

func getMachineWriteModel(ctx context.Context, userID, resourceOwner string, filter preparation.FilterToQueryReducer, permissionCheck PermissionCheck) (_ *MachineWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
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
	if permissionCheck != nil {
		if err := permissionCheck(writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
			return nil, err
		}
	}
	return writeModel, err
}
