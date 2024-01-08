package command

import (
	"context"
	"reflect"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddUserGrant(ctx context.Context, usergrant *domain.UserGrant, resourceOwner string) (_ *domain.UserGrant, err error) {
	event, addedUserGrant, err := c.addUserGrant(ctx, usergrant, resourceOwner)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(addedUserGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return userGrantWriteModelToUserGrant(addedUserGrant), nil
}

func (c *Commands) addUserGrant(ctx context.Context, userGrant *domain.UserGrant, resourceOwner string) (command eventstore.Command, _ *UserGrantWriteModel, err error) {
	if !userGrant.IsValid() {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-kVfMa", "Errors.UserGrant.Invalid")
	}
	err = c.checkUserGrantPreCondition(ctx, userGrant, resourceOwner)
	if err != nil {
		return nil, nil, err
	}
	userGrant.AggregateID, err = c.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}

	addedUserGrant := NewUserGrantWriteModel(userGrant.AggregateID, resourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&addedUserGrant.WriteModel)
	command = usergrant.NewUserGrantAddedEvent(
		ctx,
		userGrantAgg,
		userGrant.UserID,
		userGrant.ProjectID,
		userGrant.ProjectGrantID,
		userGrant.RoleKeys,
	)
	return command, addedUserGrant, nil
}

func (c *Commands) ChangeUserGrant(ctx context.Context, userGrant *domain.UserGrant, resourceOwner string) (_ *domain.UserGrant, err error) {
	event, changedUserGrant, err := c.changeUserGrant(ctx, userGrant, resourceOwner, false)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(changedUserGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return userGrantWriteModelToUserGrant(changedUserGrant), nil
}

func (c *Commands) changeUserGrant(ctx context.Context, userGrant *domain.UserGrant, resourceOwner string, cascade bool) (_ eventstore.Command, _ *UserGrantWriteModel, err error) {
	if userGrant.AggregateID == "" {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-3M0sd", "Errors.UserGrant.Invalid")
	}
	existingUserGrant, err := c.userGrantWriteModelByID(ctx, userGrant.AggregateID, userGrant.ResourceOwner)
	if err != nil {
		return nil, nil, err
	}
	err = checkExplicitProjectPermission(ctx, existingUserGrant.ProjectGrantID, existingUserGrant.ProjectID)
	if err != nil {
		return nil, nil, err
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.UserGrant.NotFound")
	}
	if reflect.DeepEqual(existingUserGrant.RoleKeys, userGrant.RoleKeys) {
		return nil, nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Rs8fy", "Errors.UserGrant.NotChanged")
	}
	userGrant.ProjectID = existingUserGrant.ProjectID
	userGrant.ProjectGrantID = existingUserGrant.ProjectGrantID
	err = c.checkUserGrantPreCondition(ctx, userGrant, resourceOwner)
	if err != nil {
		return nil, nil, err
	}

	changedUserGrant := NewUserGrantWriteModel(userGrant.AggregateID, resourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&changedUserGrant.WriteModel)

	if cascade {
		return usergrant.NewUserGrantCascadeChangedEvent(ctx, userGrantAgg, userGrant.RoleKeys), existingUserGrant, nil
	}
	return usergrant.NewUserGrantChangedEvent(ctx, userGrantAgg, userGrant.RoleKeys), existingUserGrant, nil
}

func (c *Commands) removeRoleFromUserGrant(ctx context.Context, userGrantID string, roleKeys []string, cascade bool) (_ eventstore.Command, err error) {
	existingUserGrant, err := c.userGrantWriteModelByID(ctx, userGrantID, "")
	if err != nil {
		return nil, err
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.UserGrant.NotFound")
	}
	keyExists := false
	for i, key := range existingUserGrant.RoleKeys {
		for _, roleKey := range roleKeys {
			if key == roleKey {
				keyExists = true
				copy(existingUserGrant.RoleKeys[i:], existingUserGrant.RoleKeys[i+1:])
				existingUserGrant.RoleKeys[len(existingUserGrant.RoleKeys)-1] = ""
				existingUserGrant.RoleKeys = existingUserGrant.RoleKeys[:len(existingUserGrant.RoleKeys)-1]
				continue
			}
		}
	}
	if !keyExists {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-5m8g9", "Errors.UserGrant.RoleKeyNotFound")
	}
	changedUserGrant := NewUserGrantWriteModel(userGrantID, existingUserGrant.ResourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&changedUserGrant.WriteModel)

	if cascade {
		return usergrant.NewUserGrantCascadeChangedEvent(ctx, userGrantAgg, existingUserGrant.RoleKeys), nil
	}

	return usergrant.NewUserGrantChangedEvent(ctx, userGrantAgg, existingUserGrant.RoleKeys), nil
}

func (c *Commands) DeactivateUserGrant(ctx context.Context, grantID, resourceOwner string) (objectDetails *domain.ObjectDetails, err error) {
	if grantID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-N2OhG", "Errors.UserGrant.IDMissing")
	}

	existingUserGrant, err := c.userGrantWriteModelByID(ctx, grantID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.UserGrant.NotFound")
	}
	if existingUserGrant.State != domain.UserGrantStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-1S9gx", "Errors.UserGrant.NotActive")
	}
	err = checkExplicitProjectPermission(ctx, existingUserGrant.ProjectGrantID, existingUserGrant.ProjectID)
	if err != nil {
		return nil, err
	}

	deactivateUserGrant := NewUserGrantWriteModel(grantID, resourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&deactivateUserGrant.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, usergrant.NewUserGrantDeactivatedEvent(ctx, userGrantAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingUserGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingUserGrant.WriteModel), nil
}

func (c *Commands) ReactivateUserGrant(ctx context.Context, grantID, resourceOwner string) (objectDetails *domain.ObjectDetails, err error) {
	if grantID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Qxy8v", "Errors.UserGrant.IDMissing")
	}

	existingUserGrant, err := c.userGrantWriteModelByID(ctx, grantID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Lp0gs", "Errors.UserGrant.NotFound")
	}
	if existingUserGrant.State != domain.UserGrantStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-1ML0v", "Errors.UserGrant.NotInactive")
	}
	err = checkExplicitProjectPermission(ctx, existingUserGrant.ProjectGrantID, existingUserGrant.ProjectID)
	if err != nil {
		return nil, err
	}
	deactivateUserGrant := NewUserGrantWriteModel(grantID, resourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&deactivateUserGrant.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, usergrant.NewUserGrantReactivatedEvent(ctx, userGrantAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingUserGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingUserGrant.WriteModel), nil
}

func (c *Commands) RemoveUserGrant(ctx context.Context, grantID, resourceOwner string) (objectDetails *domain.ObjectDetails, err error) {
	event, existingUserGrant, err := c.removeUserGrant(ctx, grantID, resourceOwner, false)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingUserGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingUserGrant.WriteModel), nil
}

func (c *Commands) BulkRemoveUserGrant(ctx context.Context, grantIDs []string, resourceOwner string) (err error) {
	if len(grantIDs) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-5M0sd", "Errors.UserGrant.IDMissing")
	}
	events := make([]eventstore.Command, len(grantIDs))
	for i, grantID := range grantIDs {
		event, _, err := c.removeUserGrant(ctx, grantID, resourceOwner, false)
		if err != nil {
			return err
		}
		events[i] = event
	}
	_, err = c.eventstore.Push(ctx, events...)
	return err
}

func (c *Commands) removeUserGrant(ctx context.Context, grantID, resourceOwner string, cascade bool) (_ eventstore.Command, writeModel *UserGrantWriteModel, err error) {
	if grantID == "" {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-J9sc5", "Errors.UserGrant.IDMissing")
	}

	existingUserGrant, err := c.userGrantWriteModelByID(ctx, grantID, resourceOwner)
	if err != nil {
		return nil, nil, err
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-1My0t", "Errors.UserGrant.NotFound")
	}
	if !cascade {
		err = checkExplicitProjectPermission(ctx, existingUserGrant.ProjectGrantID, existingUserGrant.ProjectID)
		if err != nil {
			return nil, nil, err
		}
	}

	removeUserGrant := NewUserGrantWriteModel(grantID, existingUserGrant.ResourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&removeUserGrant.WriteModel)
	if cascade {
		return usergrant.NewUserGrantCascadeRemovedEvent(
			ctx,
			userGrantAgg,
			existingUserGrant.UserID,
			existingUserGrant.ProjectID,
			existingUserGrant.ProjectGrantID), existingUserGrant, nil
	}
	return usergrant.NewUserGrantRemovedEvent(
		ctx,
		userGrantAgg,
		existingUserGrant.UserID,
		existingUserGrant.ProjectID,
		existingUserGrant.ProjectGrantID), existingUserGrant, nil
}

func (c *Commands) userGrantWriteModelByID(ctx context.Context, userGrantID, resourceOwner string) (writeModel *UserGrantWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserGrantWriteModel(userGrantID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) checkUserGrantPreCondition(ctx context.Context, usergrant *domain.UserGrant, resourceOwner string) error {
	preConditions := NewUserGrantPreConditionReadModel(usergrant.UserID, usergrant.ProjectID, usergrant.ProjectGrantID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, preConditions)
	if err != nil {
		return err
	}
	if !preConditions.UserExists {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-4f8sg", "Errors.User.NotFound")
	}
	if usergrant.ProjectGrantID == "" && !preConditions.ProjectExists {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-3n77S", "Errors.Project.NotFound")
	}
	if usergrant.ProjectGrantID != "" && !preConditions.ProjectGrantExists {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-4m9ff", "Errors.Project.Grant.NotFound")
	}
	if usergrant.HasInvalidRoles(preConditions.ExistingRoleKeys) {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-mm9F4", "Errors.Project.Role.NotFound")
	}
	return nil
}
