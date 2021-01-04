package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) AddUser(ctx context.Context, user *usr_model.User) (*usr_model.User, error) {
	if !user.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2N9fs", "Errors.User.Invalid")
	}

	if user.Human != nil {
		human, err := r.AddHuman(ctx, user.ResourceOwner, user.UserName, user.Human)
		if err != nil {
			return nil, err
		}
		return &usr_model.User{UserName: user.UserName, Human: human}, nil
	} else if user.Machine != nil {

	}
	return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-8K0df", "Errors.User.TypeUndefined")
}

func (r *CommandSide) DeactivateUser(ctx context.Context, userID string) (*usr_model.User, error) {
	existingUser, err := r.userWriteModelByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existingUser.UserState != domain.UserStateUnspecified || existingUser.UserState != domain.UserStateDeleted {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-3M9ds", "Errors.User.NotFound")
	}
	if existingUser.UserState == domain.UserStateInactive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M0sf", "Errors.User.AlreadyInactive")
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	userAgg.PushEvents(user.NewUserDeactivatedEvent(ctx))

	err = r.eventstore.PushAggregate(ctx, existingUser, userAgg)
	if err != nil {
		return nil, err
	}
	return writeModelToUser(existingUser), nil
}

func (r *CommandSide) ReactivateUser(ctx context.Context, userID string) (*usr_model.User, error) {
	existingUser, err := r.userWriteModelByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existingUser.UserState != domain.UserStateUnspecified || existingUser.UserState != domain.UserStateDeleted {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-4M0sd", "Errors.User.NotFound")
	}
	if existingUser.UserState != domain.UserStateInactive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M0sf", "Errors.User.NotInactive")
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	userAgg.PushEvents(user.NewUserReactivatedEvent(ctx))

	err = r.eventstore.PushAggregate(ctx, existingUser, userAgg)
	if err != nil {
		return nil, err
	}
	return writeModelToUser(existingUser), nil
}

func (r *CommandSide) LockUser(ctx context.Context, userID string) (*usr_model.User, error) {
	existingUser, err := r.userWriteModelByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existingUser.UserState != domain.UserStateUnspecified || existingUser.UserState != domain.UserStateDeleted {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-5M9fs", "Errors.User.NotFound")
	}
	if existingUser.UserState != domain.UserStateActive && existingUser.UserState != domain.UserStateInitial {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9fs", "Errors.User.ShouldBeActiveOrInitial")
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	userAgg.PushEvents(user.NewUserLockedEvent(ctx))

	err = r.eventstore.PushAggregate(ctx, existingUser, userAgg)
	if err != nil {
		return nil, err
	}
	return writeModelToUser(existingUser), nil
}

func (r *CommandSide) UnlockUser(ctx context.Context, userID string) (*usr_model.User, error) {
	existingUser, err := r.userWriteModelByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existingUser.UserState != domain.UserStateUnspecified || existingUser.UserState != domain.UserStateDeleted {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-M0dos", "Errors.User.NotFound")
	}
	if existingUser.UserState != domain.UserStateLocked {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0ds", "Errors.User.NotLocked")
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	userAgg.PushEvents(user.NewUserUnlockedEvent(ctx))

	err = r.eventstore.PushAggregate(ctx, existingUser, userAgg)
	if err != nil {
		return nil, err
	}
	return writeModelToUser(existingUser), nil
}

func (r *CommandSide) userWriteModelByID(ctx context.Context, userID string) (writeModel *UserWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserWriteModel(userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
