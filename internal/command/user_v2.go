package command

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
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

func (c *Commands) checkPermissionUpdateUser(ctx context.Context, resourceOwner, userID string) error {
	if userID != "" && userID == authz.GetCtxData(ctx).UserID {
		return nil
	}
	if err := c.checkPermission(ctx, domain.PermissionUserWrite, resourceOwner, userID); err != nil {
		return err
	}
	return nil
}

func (c *Commands) checkPermissionUpdateUserCredentials(ctx context.Context, resourceOwner, userID string) error {
	if userID != "" && userID == authz.GetCtxData(ctx).UserID {
		return nil
	}
	if err := c.checkPermission(ctx, domain.PermissionUserCredentialWrite, resourceOwner, userID); err != nil {
		return err
	}
	return nil
}

func (c *Commands) checkPermissionDeleteUser(ctx context.Context, resourceOwner, userID string) error {
	if userID != "" && userID == authz.GetCtxData(ctx).UserID {
		return nil
	}
	if err := c.checkPermission(ctx, domain.PermissionUserDelete, resourceOwner, userID); err != nil {
		return err
	}
	return nil
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

func (c *Commands) RemoveUserV2(ctx context.Context, userID string, cascadingUserMemberships []*CascadingMembership, cascadingGrantIDs ...string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-vaipl7s13l", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userRemoveWriteModel(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingUser.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-bd4ir1mblj", "Errors.User.NotFound")
	}
	if err := c.checkPermissionDeleteUser(ctx, existingUser.ResourceOwner, existingUser.AggregateID); err != nil {
		return nil, err
	}

	domainPolicy, err := c.domainPolicyWriteModel(ctx, existingUser.ResourceOwner)
	if err != nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-l40ykb3xh2", "Errors.Org.DomainPolicy.NotExisting")
	}
	var events []eventstore.Command
	events = append(events, user.NewUserRemovedEvent(ctx, &existingUser.Aggregate().Aggregate, existingUser.UserName, existingUser.IDPLinks, domainPolicy.UserLoginMustBeDomain))

	for _, grantID := range cascadingGrantIDs {
		removeEvent, _, err := c.removeUserGrant(ctx, grantID, "", true)
		if err != nil {
			logging.WithFields("usergrantid", grantID).WithError(err).Warn("could not cascade remove role on user grant")
			continue
		}
		events = append(events, removeEvent)
	}

	if len(cascadingUserMemberships) > 0 {
		membershipEvents, err := c.removeUserMemberships(ctx, cascadingUserMemberships)
		if err != nil {
			return nil, err
		}
		events = append(events, membershipEvents...)
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingUser, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingUser.WriteModel), nil
}

func (c *Commands) userRemoveWriteModel(ctx context.Context, userID string) (writeModel *UserV2WriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserRemoveWriteModel(userID, "")
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
