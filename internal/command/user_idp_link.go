package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddUserIDPLink(ctx context.Context, userID, resourceOwner string, link *AddLink) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-03j8f", "Errors.IDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingUser.UserState) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-vzktar7b7f", "Errors.User.NotFound")
	}
	if userID != authz.GetCtxData(ctx).UserID {
		if err := c.checkPermission(ctx, domain.PermissionUserWrite, existingUser.ResourceOwner, existingUser.AggregateID); err != nil {
			return nil, err
		}
	}
	//nolint:staticcheck
	event, err := addLink(ctx, c.eventstore.Filter, user.NewAggregate(existingUser.AggregateID, existingUser.ResourceOwner), link)
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func (c *Commands) BulkAddedUserIDPLinks(ctx context.Context, userID, resourceOwner string, links []*domain.UserIDPLink) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-03j8f", "Errors.IDMissing")
	}
	if len(links) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-Ek9s", "Errors.User.ExternalIDP.MinimumExternalIDPNeeded")
	}

	if err := c.checkUserExists(ctx, userID, resourceOwner); err != nil {
		return err
	}

	events := make([]eventstore.Command, len(links))
	for i, link := range links {
		linkWriteModel := NewUserIDPLinkWriteModel(userID, link.IDPConfigID, link.ExternalUserID, resourceOwner)
		userAgg := UserAggregateFromWriteModel(&linkWriteModel.WriteModel)

		events[i], err = c.addUserIDPLink(ctx, userAgg, link, true)
		if err != nil {
			return err
		}
	}

	_, err = c.eventstore.Push(ctx, events...)
	return err
}

func (c *Commands) addUserIDPLink(ctx context.Context, human *eventstore.Aggregate, link *domain.UserIDPLink, linkToExistingUser bool) (eventstore.Command, error) {
	if link.AggregateID != "" && human.ID != link.AggregateID {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-33M0g", "Errors.IDMissing")
	}
	if !link.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-6m9Kd", "Errors.User.ExternalIDP.Invalid")
	}
	idpWriteModel, err := IDPProviderWriteModel(ctx, c.eventstore.Filter, link.IDPConfigID)
	if err != nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-39nfs", "Errors.IDPConfig.NotExisting")
	}
	// IDP user will either be linked or created on a new user
	// Therefore we need to either check if linking is allowed or creation:
	if linkToExistingUser && !idpWriteModel.GetProviderOptions().IsLinkingAllowed {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-Sfee2", "Errors.ExternalIDP.LinkingNotAllowed")
	}
	if !linkToExistingUser && !idpWriteModel.GetProviderOptions().IsCreationAllowed {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-SJI3g", "Errors.ExternalIDP.CreationNotAllowed")
	}
	return user.NewUserIDPLinkAddedEvent(ctx, human, link.IDPConfigID, link.DisplayName, link.ExternalUserID), nil

}

func (c *Commands) RemoveUserIDPLink(ctx context.Context, link *domain.UserIDPLink) (*domain.ObjectDetails, error) {
	event, linkWriteModel, err := c.removeUserIDPLink(ctx, link, false)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(linkWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&linkWriteModel.WriteModel), nil
}

func (c *Commands) removeUserIDPLink(ctx context.Context, link *domain.UserIDPLink, cascade bool) (eventstore.Command, *UserIDPLinkWriteModel, error) {
	if !link.IsValid() || link.AggregateID == "" {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-3M9ds", "Errors.IDMissing")
	}

	existingLink, err := c.userIDPLinkWriteModelByID(ctx, link.AggregateID, link.IDPConfigID, link.ExternalUserID, link.ResourceOwner)
	if err != nil {
		return nil, nil, err
	}
	if existingLink.State == domain.UserIDPLinkStateUnspecified || existingLink.State == domain.UserIDPLinkStateRemoved {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-1M9xR", "Errors.User.ExternalIDP.NotFound")
	}
	if existingLink.AggregateID != authz.GetCtxData(ctx).UserID {
		if err := c.checkPermission(ctx, domain.PermissionUserWrite, existingLink.ResourceOwner, existingLink.AggregateID); err != nil {
			return nil, nil, err
		}
	}
	userAgg := UserAggregateFromWriteModel(&existingLink.WriteModel)
	if cascade {
		return user.NewUserIDPLinkCascadeRemovedEvent(ctx, userAgg, link.IDPConfigID, link.ExternalUserID), existingLink, nil
	}
	return user.NewUserIDPLinkRemovedEvent(ctx, userAgg, link.IDPConfigID, link.ExternalUserID), existingLink, nil
}

func (c *Commands) UserIDPLoginChecked(ctx context.Context, orgID, userID string, authRequest *domain.AuthRequest) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-5n8sM", "Errors.IDMissing")
	}

	existingHuman, err := c.getHumanWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingHuman.UserState == domain.UserStateUnspecified || existingHuman.UserState == domain.UserStateDeleted {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-dn88J", "Errors.User.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingHuman.WriteModel)
	_, err = c.eventstore.Push(ctx, user.NewUserIDPCheckSucceededEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
	return err
}

func (c *Commands) MigrateUserIDP(ctx context.Context, userID, orgID, idpConfigID, previousID, newID string) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-Sn3l1", "Errors.IDMissing")
	}

	writeModel, err := c.userIDPLinkWriteModelByID(ctx, userID, idpConfigID, previousID, orgID)
	if err != nil {
		return err
	}
	if writeModel.State != domain.UserIDPLinkStateActive {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-KJH2o", "Errors.User.ExternalIDP.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&writeModel.WriteModel)
	_, err = c.eventstore.Push(ctx, user.NewUserIDPExternalIDMigratedEvent(ctx, userAgg, idpConfigID, previousID, newID))
	return err
}

func (c *Commands) UpdateUserIDPLinkUsername(ctx context.Context, userID, orgID, idpConfigID, externalID, newUsername string) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-SFegz", "Errors.IDMissing")
	}

	writeModel, err := c.userIDPLinkWriteModelByID(ctx, userID, idpConfigID, externalID, orgID)
	if err != nil {
		return err
	}
	if writeModel.State != domain.UserIDPLinkStateActive {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-DGhre", "Errors.User.ExternalIDP.NotFound")
	}
	if writeModel.DisplayName == newUsername {
		return nil
	}

	userAgg := UserAggregateFromWriteModel(&writeModel.WriteModel) //nolint:contextcheck
	_, err = c.eventstore.Push(ctx, user.NewUserIDPExternalUsernameEvent(ctx, userAgg, idpConfigID, externalID, newUsername))
	return err
}

func (c *Commands) userIDPLinkWriteModelByID(ctx context.Context, userID, idpConfigID, externalUserID, resourceOwner string) (writeModel *UserIDPLinkWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserIDPLinkWriteModel(userID, idpConfigID, externalUserID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
