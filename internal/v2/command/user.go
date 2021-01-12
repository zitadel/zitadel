package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) getUser(ctx context.Context, userID, resourceOwner string) (*domain.User, error) {
	writeModel, err := r.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if writeModel.UserState == domain.UserStateUnspecified || writeModel.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-B8dQw", "Errors.User.NotFound")
	}
	return writeModelToUser(writeModel), nil
}

func (r *CommandSide) AddUser(ctx context.Context, orgID string, user *domain.User) (*domain.User, error) {
	if !user.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2N9fs", "Errors.User.Invalid")
	}

	if user.Human != nil {
		human, err := r.AddHuman(ctx, orgID, user.UserName, user.Human)
		if err != nil {
			return nil, err
		}
		return &domain.User{UserName: user.UserName, Human: human}, nil
	} else if user.Machine != nil {
		machine, err := r.AddMachine(ctx, orgID, user.UserName, user.Machine)
		if err != nil {
			return nil, err
		}
		return &domain.User{UserName: user.UserName, Machine: machine}, nil
	}
	return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-8K0df", "Errors.User.TypeUndefined")
}

func (r *CommandSide) RegisterUser(ctx context.Context, orgID string, user *domain.User) (*domain.User, error) {
	if !user.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2N9fs", "Errors.User.Invalid")
	}

	if user.Human != nil {
		human, err := r.RegisterHuman(ctx, orgID, user.UserName, user.Human, nil)
		if err != nil {
			return nil, err
		}
		return &domain.User{UserName: user.UserName, Human: human}, nil
	}
	return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-8K0df", "Errors.User.TypeUndefined")
}

func (r *CommandSide) ChangeUsername(ctx context.Context, orgID, userID, userName string) error {
	if orgID == "" || userID == "" || userName == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2N9fs", "Errors.IDMissing")
	}
	existingUser, err := r.userWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5N9ds", "Errors.User.NotFound")
	}
	if existingUser.UserName == userName {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9fs", "Errors.User.UsernameNotChanged")
	}

	orgIAMPolicy, err := r.getOrgIAMPolicy(ctx, orgID)
	if err != nil {
		return err
	}
	if err := CheckOrgIAMPolicyForUserName(userName, orgIAMPolicy); err != nil {
		return err
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	userAgg.PushEvents(user.NewUsernameChangedEvent(ctx, userName))
	//TODO: Check Uniqueness
	//TODO: release old username, set new unique username

	return r.eventstore.PushAggregate(ctx, existingUser, userAgg)
}

func (r *CommandSide) DeactivateUser(ctx context.Context, userID, resourceOwner string) (*domain.User, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-m0gDf", "Errors.User.UserIDMissing")
	}
	existingUser, err := r.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-3M9ds", "Errors.User.NotFound")
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

func (r *CommandSide) ReactivateUser(ctx context.Context, userID, resourceOwner string) (*domain.User, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M9ds", "Errors.User.UserIDMissing")
	}
	existingUser, err := r.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-4M0sd", "Errors.User.NotFound")
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

func (r *CommandSide) LockUser(ctx context.Context, userID, resourceOwner string) (*domain.User, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M0sd", "Errors.User.UserIDMissing")
	}
	existingUser, err := r.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-5M9fs", "Errors.User.NotFound")
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

func (r *CommandSide) UnlockUser(ctx context.Context, userID, resourceOwner string) (*domain.User, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-M0dse", "Errors.User.UserIDMissing")
	}
	existingUser, err := r.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-M0dos", "Errors.User.NotFound")
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

func (r *CommandSide) RemoveUser(ctx context.Context, userID, resourceOwner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M0ds", "Errors.User.UserIDMissing")
	}
	existingUser, err := r.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5M0od", "Errors.User.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	userAgg.PushEvents(user.NewUserRemovedEvent(ctx))
	//TODO: release unqie username
	//TODO: remove user grants

	return r.eventstore.PushAggregate(ctx, existingUser, userAgg)
}

func (r *CommandSide) userWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *UserWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserWriteModel(userID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
