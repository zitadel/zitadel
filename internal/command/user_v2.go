package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) LockUserV2(ctx context.Context, userID string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-agz3eczifm", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.userStateWriteModel(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-450yxuqrh1", "Errors.User.NotFound")
	}
	if !hasUserState(existingHuman.UserState, domain.UserStateActive, domain.UserStateInitial) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-lgws8wtsqf", "Errors.User.ShouldBeActiveOrInitial")
	}

	if err := c.checkPermissionUpdateUser(ctx, existingHuman.ResourceOwner, existingHuman.AggregateID); err != nil {
		return nil, err
	}

	if err := c.pushAppendAndReduce(ctx, existingHuman, user.NewUserLockedEvent(ctx, &existingHuman.Aggregate().Aggregate)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingHuman.WriteModel), nil
}

func (c *Commands) UnlockUserV2(ctx context.Context, userID string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-a9ld4xckax", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.userStateWriteModel(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-x377t913pw", "Errors.User.NotFound")
	}
	if !hasUserState(existingHuman.UserState, domain.UserStateLocked) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-olb9vb0oca", "Errors.User.NotLocked")
	}
	if err := c.checkPermissionUpdateUser(ctx, existingHuman.ResourceOwner, existingHuman.AggregateID); err != nil {
		return nil, err
	}

	if err := c.pushAppendAndReduce(ctx, existingHuman, user.NewUserUnlockedEvent(ctx, &existingHuman.Aggregate().Aggregate)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingHuman.WriteModel), nil
}

func (c *Commands) DeactivateUserV2(ctx context.Context, userID string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-78iiirat8y", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.userStateWriteModel(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-5gp2p62iin", "Errors.User.NotFound")
	}
	if isUserStateInitial(existingHuman.UserState) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-gvx4kct9r2", "Errors.User.CantDeactivateInitial")
	}
	if isUserStateInactive(existingHuman.UserState) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-5gunjw0cd7", "Errors.User.AlreadyInactive")
	}
	if err := c.checkPermissionUpdateUser(ctx, existingHuman.ResourceOwner, existingHuman.AggregateID); err != nil {
		return nil, err
	}

	if err := c.pushAppendAndReduce(ctx, existingHuman, user.NewUserDeactivatedEvent(ctx, &existingHuman.Aggregate().Aggregate)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingHuman.WriteModel), nil
}

func (c *Commands) ReactivateUserV2(ctx context.Context, userID string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-0nx1ie38fw", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.userStateWriteModel(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-9hy5kzbuk6", "Errors.User.NotFound")
	}
	if !isUserStateInactive(existingHuman.UserState) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-s5qqcz97hf", "Errors.User.NotInactive")
	}
	if err := c.checkPermissionUpdateUser(ctx, existingHuman.ResourceOwner, existingHuman.AggregateID); err != nil {
		return nil, err
	}

	if err := c.pushAppendAndReduce(ctx, existingHuman, user.NewUserReactivatedEvent(ctx, &existingHuman.Aggregate().Aggregate)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingHuman.WriteModel), nil
}

func (c *Commands) userStateWriteModel(ctx context.Context, userID string) (writeModel *UserV2WriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserStateWriteModel(userID, "")
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
