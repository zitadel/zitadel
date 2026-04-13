package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type GenerateMachineSecret struct {
	PermissionCheck PermissionCheck
	ClientSecret    string
}

func (c *Commands) GenerateMachineSecret(ctx context.Context, userID string, resourceOwner string, set *GenerateMachineSecret) (*domain.ObjectDetails, error) {
	agg := user.NewAggregate(userID, resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareGenerateMachineSecret(agg, set)) //nolint:staticcheck
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

func (c *Commands) prepareGenerateMachineSecret(a *user.Aggregate, set *GenerateMachineSecret) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if a.ResourceOwner == "" && set.PermissionCheck == nil {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-0qp2hus", "Errors.ResourceOwnerMissing")
		}
		if a.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-bzoqjs", "Errors.User.UserIDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getMachineWriteModel(ctx, a.ID, a.ResourceOwner, filter, set.PermissionCheck)
			if err != nil {
				return nil, err
			}
			if !isUserStateExists(writeModel.UserState) {
				return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-x8910n", "Errors.User.NotExisting")
			}
			encodedHash, plain, err := c.newHashedSecret(ctx, filter)
			if err != nil {
				return nil, err
			}
			set.ClientSecret = plain

			return []eventstore.Command{
				user.NewMachineSecretSetEvent(ctx, &a.Aggregate, encodedHash),
			}, nil
		}, nil
	}
}

func (c *Commands) RemoveMachineSecret(ctx context.Context, userID string, resourceOwner string, permissionCheck PermissionCheck) (*domain.ObjectDetails, error) {
	agg := user.NewAggregate(userID, resourceOwner)
	//nolint:staticcheck
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareRemoveMachineSecret(agg, permissionCheck))
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

func prepareRemoveMachineSecret(a *user.Aggregate, check PermissionCheck) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if a.ResourceOwner == "" && check == nil {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-x0992n", "Errors.ResourceOwnerMissing")
		}
		if a.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-bzosjs", "Errors.User.UserIDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getMachineWriteModel(ctx, a.ID, a.ResourceOwner, filter, check)
			if err != nil {
				return nil, err
			}
			if !isUserStateExists(writeModel.UserState) {
				return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-x7s802", "Errors.User.NotExisting")
			}
			if writeModel.HashedSecret == "" {
				return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-coi82n", "Errors.User.Machine.Secret.NotExisting")
			}
			return []eventstore.Command{
				user.NewMachineSecretRemovedEvent(ctx, &a.Aggregate),
			}, nil
		}, nil
	}
}

func (c *Commands) MachineSecretCheckSucceeded(ctx context.Context, userID, resourceOwner, updated string) {
	agg := user.NewAggregate(userID, resourceOwner)
	if updated != "" {
		cmds := []eventstore.Command{user.NewMachineSecretHashUpdatedEvent(ctx, &agg.Aggregate, updated)}
		c.asyncPush(ctx, cmds...)
	}
}
