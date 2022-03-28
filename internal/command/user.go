package command

import (
	"context"
	"fmt"
	"time"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/query"

	"github.com/caos/zitadel/internal/eventstore/v1/models"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) ChangeUsername(ctx context.Context, instanceID, orgID, userID, userName string) (*domain.ObjectDetails, error) {
	if orgID == "" || userID == "" || userName == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-2N9fs", "Errors.IDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}

	if !isUserStateExists(existingUser.UserState) {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-5N9ds", "Errors.User.NotFound")
	}

	if existingUser.UserName == userName {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6m9gs", "Errors.User.UsernameNotChanged")
	}

	domainPolicy, err := c.getOrgDomainPolicy(ctx, instanceID, orgID)
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-38fnu", "Errors.Org.DomainPolicy.NotExisting")
	}

	if err := CheckDomainPolicyForUserName(userName, domainPolicy); err != nil {
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
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-m0gDf", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingUser.UserState) {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-3M9ds", "Errors.User.NotFound")
	}
	if isUserStateInitial(existingUser.UserState) {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-ke0fw", "Errors.User.CantDeactivateInitial")
	}
	if isUserStateInactive(existingUser.UserState) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M0sf", "Errors.User.AlreadyInactive")
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
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-4M9ds", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingUser.UserState) {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-4M0sd", "Errors.User.NotFound")
	}
	if !isUserStateInactive(existingUser.UserState) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M0sf", "Errors.User.NotInactive")
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
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-2M0sd", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingUser.UserState) {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-5M9fs", "Errors.User.NotFound")
	}
	if !hasUserState(existingUser.UserState, domain.UserStateActive, domain.UserStateInitial) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3NN8v", "Errors.User.ShouldBeActiveOrInitial")
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
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-M0dse", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingUser.UserState) {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-M0dos", "Errors.User.NotFound")
	}
	if !hasUserState(existingUser.UserState, domain.UserStateLocked) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0ds", "Errors.User.NotLocked")
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

func (c *Commands) RemoveUser(ctx context.Context, instanceID, userID, resourceOwner string, cascadingUserMemberships []*query.Membership, cascadingGrantIDs ...string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-2M0ds", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingUser.UserState) {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-m9od", "Errors.User.NotFound")
	}

	domainPolicy, err := c.getOrgDomainPolicy(ctx, instanceID, existingUser.ResourceOwner)
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-3M9fs", "Errors.Org.DomainPolicy.NotExisting")
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

func (c *Commands) AddUserToken(ctx context.Context, orgID, agentID, clientID, userID string, audience, scopes []string, lifetime time.Duration) (*domain.Token, error) {
	if userID == "" { //do not check for empty orgID (JWT Profile requests won't provide it, so service user requests fail)
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-Dbge4", "Errors.IDMissing")
	}
	userWriteModel := NewUserWriteModel(userID, orgID)
	event, accessToken, err := c.addUserToken(ctx, userWriteModel, agentID, clientID, "", audience, scopes, lifetime)
	if err != nil {
		return nil, err
	}
	_, err = c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
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

func (c *Commands) addUserToken(ctx context.Context, userWriteModel *UserWriteModel, agentID, clientID, refreshTokenID string, audience, scopes []string, lifetime time.Duration) (*user.UserTokenAddedEvent, *domain.Token, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, userWriteModel)
	if err != nil {
		return nil, nil, err
	}
	if !isUserStateExists(userWriteModel.UserState) {
		return nil, nil, caos_errs.ThrowNotFound(nil, "COMMAND-1d6Gg", "Errors.User.NotFound")
	}

	audience = domain.AddAudScopeToAudience(audience, scopes)

	preferredLanguage := ""
	existingHuman, err := c.getHumanWriteModelByID(ctx, userWriteModel.AggregateID, userWriteModel.ResourceOwner)
	if existingHuman != nil {
		preferredLanguage = existingHuman.PreferredLanguage.String()
	}
	expiration := time.Now().UTC().Add(lifetime)
	tokenID, err := c.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}

	userAgg := UserAggregateFromWriteModel(&userWriteModel.WriteModel)
	return user.NewUserTokenAddedEvent(ctx, userAgg, tokenID, clientID, agentID, preferredLanguage, refreshTokenID, audience, scopes, expiration),
		&domain.Token{
			ObjectRoot: models.ObjectRoot{
				AggregateID: userWriteModel.AggregateID,
			},
			TokenID:           tokenID,
			UserAgentID:       agentID,
			ApplicationID:     clientID,
			RefreshTokenID:    refreshTokenID,
			Audience:          audience,
			Scopes:            scopes,
			Expiration:        expiration,
			PreferredLanguage: preferredLanguage,
		}, nil
}

func (c *Commands) removeAccessToken(ctx context.Context, userID, orgID, tokenID string) (*user.UserTokenRemovedEvent, *UserAccessTokenWriteModel, error) {
	if userID == "" || orgID == "" || tokenID == "" {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-Dng42", "Errors.IDMissing")
	}
	refreshTokenWriteModel := NewUserAccessTokenWriteModel(userID, orgID, tokenID)
	err := c.eventstore.FilterToQueryReducer(ctx, refreshTokenWriteModel)
	if err != nil {
		return nil, nil, err
	}
	if refreshTokenWriteModel.UserState != domain.UserStateActive {
		return nil, nil, caos_errs.ThrowNotFound(nil, "COMMAND-BF4hd", "Errors.User.AccessToken.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&refreshTokenWriteModel.WriteModel)
	return user.NewUserTokenRemovedEvent(ctx, userAgg, tokenID), refreshTokenWriteModel, nil
}

func (c *Commands) userDomainClaimed(ctx context.Context, instanceID, userID string) (events []eventstore.Command, _ *UserWriteModel, err error) {
	existingUser, err := c.userWriteModelByID(ctx, userID, "")
	if err != nil {
		return nil, nil, err
	}
	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return nil, nil, caos_errs.ThrowNotFound(nil, "COMMAND-ii9K0", "Errors.User.NotFound")
	}
	changedUserGrant := NewUserWriteModel(userID, existingUser.ResourceOwner)
	userAgg := UserAggregateFromWriteModel(&changedUserGrant.WriteModel)

	domainPolicy, err := c.getOrgDomainPolicy(ctx, instanceID, existingUser.ResourceOwner)
	if err != nil {
		return nil, nil, err
	}

	id, err := c.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}
	return []eventstore.Command{
		user.NewDomainClaimedEvent(
			ctx,
			userAgg,
			fmt.Sprintf("%s@temporary.%s", id, c.iamDomain),
			existingUser.UserName,
			domainPolicy.UserLoginMustBeDomain),
	}, changedUserGrant, nil
}

func (c *Commands) UserDomainClaimedSent(ctx context.Context, orgID, userID string) (err error) {
	if userID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-5m0fs", "Errors.IDMissing")
	}
	existingUser, err := c.userWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingUser.UserState) {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5m9gK", "Errors.User.NotFound")
	}

	_, err = c.eventstore.Push(ctx,
		user.NewDomainClaimedSentEvent(ctx, UserAggregateFromWriteModel(&existingUser.WriteModel)))
	return err
}

func (c *Commands) checkUserExists(ctx context.Context, userID, resourceOwner string) error {
	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingUser.UserState) {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0fs", "Errors.User.NotFound")
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
