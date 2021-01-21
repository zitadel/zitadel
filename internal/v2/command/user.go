package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

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
	userAgg.PushEvents(user.NewUsernameChangedEvent(ctx, existingUser.UserName, userName, orgIAMPolicy.UserLoginMustBeDomain))

	return r.eventstore.PushAggregate(ctx, existingUser, userAgg)
}

func (r *CommandSide) DeactivateUser(ctx context.Context, userID, resourceOwner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-m0gDf", "Errors.User.UserIDMissing")
	}
	existingUser, err := r.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-3M9ds", "Errors.User.NotFound")
	}
	if existingUser.UserState == domain.UserStateInactive {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M0sf", "Errors.User.AlreadyInactive")
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	userAgg.PushEvents(user.NewUserDeactivatedEvent(ctx))

	return r.eventstore.PushAggregate(ctx, existingUser, userAgg)
}

func (r *CommandSide) ReactivateUser(ctx context.Context, userID, resourceOwner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M9ds", "Errors.User.UserIDMissing")
	}
	existingUser, err := r.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-4M0sd", "Errors.User.NotFound")
	}
	if existingUser.UserState != domain.UserStateInactive {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M0sf", "Errors.User.NotInactive")
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	userAgg.PushEvents(user.NewUserReactivatedEvent(ctx))

	return r.eventstore.PushAggregate(ctx, existingUser, userAgg)
}

func (r *CommandSide) LockUser(ctx context.Context, userID, resourceOwner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M0sd", "Errors.User.UserIDMissing")
	}
	existingUser, err := r.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5M9fs", "Errors.User.NotFound")
	}
	if existingUser.UserState != domain.UserStateActive && existingUser.UserState != domain.UserStateInitial {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9fs", "Errors.User.ShouldBeActiveOrInitial")
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	userAgg.PushEvents(user.NewUserLockedEvent(ctx))

	return r.eventstore.PushAggregate(ctx, existingUser, userAgg)
}

func (r *CommandSide) UnlockUser(ctx context.Context, userID, resourceOwner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-M0dse", "Errors.User.UserIDMissing")
	}
	existingUser, err := r.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-M0dos", "Errors.User.NotFound")
	}
	if existingUser.UserState != domain.UserStateLocked {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0ds", "Errors.User.NotLocked")
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	userAgg.PushEvents(user.NewUserUnlockedEvent(ctx))

	return r.eventstore.PushAggregate(ctx, existingUser, userAgg)
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
	orgIAMPolicy, err := r.getOrgIAMPolicy(ctx, existingUser.ResourceOwner)
	if err != nil {
		return err
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	userAgg.PushEvents(user.NewUserRemovedEvent(ctx, existingUser.ResourceOwner, existingUser.UserName, orgIAMPolicy.UserLoginMustBeDomain))
	//TODO: remove user grants

	return r.eventstore.PushAggregate(ctx, existingUser, userAgg)
}

func (r *CommandSide) checkUserExists(ctx context.Context, userID, resourceOwner string) (bool, error) {
	userWriteModel, err := r.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return false, err
	}
	return userWriteModel.UserState != domain.UserStateUnspecified && userWriteModel.UserState != domain.UserStateDeleted, nil
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

func (r *CommandSide) userReadModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *UserWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserWriteModel(userID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
