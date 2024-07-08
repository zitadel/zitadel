package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) ChangeUsername(ctx context.Context, orgID, userID, userName string) (*domain.ObjectDetails, error) {
	userName = strings.TrimSpace(userName)
	if orgID == "" || userID == "" || userName == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-2N9fs", "Errors.IDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}

	if !isUserStateExists(existingUser.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-5N9ds", "Errors.User.NotFound")
	}

	if existingUser.UserName == userName {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-6m9gs", "Errors.User.UsernameNotChanged")
	}

	domainPolicy, err := c.domainPolicyWriteModel(ctx, orgID)
	if err != nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-38fnu", "Errors.Org.DomainPolicy.NotExisting")
	}
	if err = c.userValidateDomain(ctx, orgID, userName, domainPolicy.UserLoginMustBeDomain); err != nil {
		return nil, err
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)

	pushedEvents, err := c.eventstore.Push(ctx,
		user.NewUsernameChangedEvent(ctx, userAgg, existingUser.UserName, userName, domainPolicy.UserLoginMustBeDomain))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingUser, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingUser.WriteModel), nil
}

func (c *Commands) DeactivateUser(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-m0gDf", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingUser.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9ds", "Errors.User.NotFound")
	}
	if isUserStateInitial(existingUser.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-ke0fw", "Errors.User.CantDeactivateInitial")
	}
	if isUserStateInactive(existingUser.UserState) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-5M0sf", "Errors.User.AlreadyInactive")
	}

	pushedEvents, err := c.eventstore.Push(ctx,
		user.NewUserDeactivatedEvent(ctx, UserAggregateFromWriteModel(&existingUser.WriteModel)))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingUser, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingUser.WriteModel), nil
}

func (c *Commands) ReactivateUser(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4M9ds", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingUser.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-4M0sd", "Errors.User.NotFound")
	}
	if !isUserStateInactive(existingUser.UserState) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-6M0sf", "Errors.User.NotInactive")
	}

	pushedEvents, err := c.eventstore.Push(ctx,
		user.NewUserReactivatedEvent(ctx, UserAggregateFromWriteModel(&existingUser.WriteModel)))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingUser, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingUser.WriteModel), nil
}

func (c *Commands) LockUser(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-2M0sd", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingUser.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-5M9fs", "Errors.User.NotFound")
	}
	if !hasUserState(existingUser.UserState, domain.UserStateActive, domain.UserStateInitial) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-3NN8v", "Errors.User.ShouldBeActiveOrInitial")
	}

	pushedEvents, err := c.eventstore.Push(ctx,
		user.NewUserLockedEvent(ctx, UserAggregateFromWriteModel(&existingUser.WriteModel)))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingUser, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingUser.WriteModel), nil
}

func (c *Commands) UnlockUser(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-M0dse", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingUser.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-M0dos", "Errors.User.NotFound")
	}
	if !hasUserState(existingUser.UserState, domain.UserStateLocked) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-4M0ds", "Errors.User.NotLocked")
	}

	pushedEvents, err := c.eventstore.Push(ctx,
		user.NewUserUnlockedEvent(ctx, UserAggregateFromWriteModel(&existingUser.WriteModel)))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingUser, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingUser.WriteModel), nil
}

func (c *Commands) RemoveUser(ctx context.Context, userID, resourceOwner string, cascadingUserMemberships []*CascadingMembership, cascadingGrantIDs ...string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-2M0ds", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingUser.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-m9od", "Errors.User.NotFound")
	}

	domainPolicy, err := c.domainPolicyWriteModel(ctx, existingUser.ResourceOwner)
	if err != nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-3M9fs", "Errors.Org.DomainPolicy.NotExisting")
	}
	var events []eventstore.Command
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	events = append(events, user.NewUserRemovedEvent(ctx, userAgg, existingUser.UserName, existingUser.IDPLinks, domainPolicy.UserLoginMustBeDomain))

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

func (c *Commands) RevokeAccessToken(ctx context.Context, userID, orgID, tokenID string) (*domain.ObjectDetails, error) {
	removeEvent, accessTokenWriteModel, err := c.removeAccessToken(ctx, userID, orgID, tokenID)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, removeEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(accessTokenWriteModel, events...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&accessTokenWriteModel.WriteModel), nil
}

func (c *Commands) removeAccessToken(ctx context.Context, userID, orgID, tokenID string) (*user.UserTokenRemovedEvent, *UserAccessTokenWriteModel, error) {
	if userID == "" || orgID == "" || tokenID == "" {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Dng42", "Errors.IDMissing")
	}
	refreshTokenWriteModel := NewUserAccessTokenWriteModel(userID, orgID, tokenID)
	err := c.eventstore.FilterToQueryReducer(ctx, refreshTokenWriteModel)
	if err != nil {
		return nil, nil, err
	}
	if refreshTokenWriteModel.UserState != domain.UserStateActive {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-BF4hd", "Errors.User.AccessToken.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&refreshTokenWriteModel.WriteModel)
	return user.NewUserTokenRemovedEvent(ctx, userAgg, tokenID), refreshTokenWriteModel, nil
}

func (c *Commands) userDomainClaimed(ctx context.Context, userID string) (events []eventstore.Command, _ *UserWriteModel, err error) {
	existingUser, err := c.userWriteModelByID(ctx, userID, "")
	if err != nil {
		return nil, nil, err
	}
	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-ii9K0", "Errors.User.NotFound")
	}
	changedUserGrant := NewUserWriteModel(userID, existingUser.ResourceOwner)
	userAgg := UserAggregateFromWriteModel(&changedUserGrant.WriteModel)

	domainPolicy, err := c.domainPolicyWriteModel(ctx, existingUser.ResourceOwner)
	if err != nil {
		return nil, nil, err
	}

	id, err := id_generator.Next()
	if err != nil {
		return nil, nil, err
	}
	return []eventstore.Command{
		user.NewDomainClaimedEvent(
			ctx,
			userAgg,
			fmt.Sprintf("%s@temporary.%s", id, authz.GetInstance(ctx).RequestedDomain()),
			existingUser.UserName,
			domainPolicy.UserLoginMustBeDomain),
	}, changedUserGrant, nil
}

func (c *Commands) prepareUserDomainClaimed(ctx context.Context, filter preparation.FilterToQueryReducer, userID string) (*user.DomainClaimedEvent, error) {
	userWriteModel, err := userWriteModelByID(ctx, filter, userID, "")
	if err != nil {
		return nil, err
	}
	if !userWriteModel.UserState.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-ii9K0", "Errors.User.NotFound")
	}
	domainPolicy, err := domainPolicyWriteModel(ctx, filter, userWriteModel.ResourceOwner)
	if err != nil {
		return nil, err
	}
	userAgg := UserAggregateFromWriteModel(&userWriteModel.WriteModel)

	id, err := id_generator.Next()
	if err != nil {
		return nil, err
	}

	return user.NewDomainClaimedEvent(
		ctx,
		userAgg,
		fmt.Sprintf("%s@temporary.%s", id, authz.GetInstance(ctx).RequestedDomain()),
		userWriteModel.UserName,
		domainPolicy.UserLoginMustBeDomain), nil
}

func (c *Commands) UserDomainClaimedSent(ctx context.Context, orgID, userID string) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-5m0fs", "Errors.IDMissing")
	}
	existingUser, err := c.userWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingUser.UserState) {
		return zerrors.ThrowNotFound(nil, "COMMAND-5m9gK", "Errors.User.NotFound")
	}

	_, err = c.eventstore.Push(ctx,
		user.NewDomainClaimedSentEvent(ctx, UserAggregateFromWriteModel(&existingUser.WriteModel)))
	return err
}

func (c *Commands) checkUserExists(ctx context.Context, userID, resourceOwner string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingUser.UserState) {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-uXHNj", "Errors.User.NotFound")
	}
	return nil
}

func (c *Commands) userWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *UserWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func ExistsUser(ctx context.Context, filter preparation.FilterToQueryReducer, id, resourceOwner string) (exists bool, err error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(resourceOwner).
		OrderAsc().
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(id).
		EventTypes(
			user.HumanRegisteredType,
			user.UserV1RegisteredType,
			user.HumanAddedType,
			user.UserV1AddedType,
			user.MachineAddedEventType,
			user.UserRemovedType,
		).Builder())
	if err != nil {
		return false, err
	}

	for _, event := range events {
		switch event.(type) {
		case *user.HumanRegisteredEvent, *user.HumanAddedEvent, *user.MachineAddedEvent:
			exists = true
		case *user.UserRemovedEvent:
			exists = false
		}
	}

	return exists, nil
}

func (c *Commands) newUserInitCode(ctx context.Context, filter preparation.FilterToQueryReducer, alg crypto.EncryptionAlgorithm) (*EncryptedCode, error) {
	return c.newEncryptedCode(ctx, filter, domain.SecretGeneratorTypeInitCode, alg)
}

func userWriteModelByID(ctx context.Context, filter preparation.FilterToQueryReducer, userID, resourceOwner string) (*UserWriteModel, error) {
	user := NewUserWriteModel(userID, resourceOwner)
	events, err := filter(ctx, user.Query())
	if err != nil {
		return nil, err
	}
	user.AppendEvents(events...)
	err = user.Reduce()
	return user, err
}
