package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"reflect"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/usergrant"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (r *CommandSide) AddUserGrant(ctx context.Context, usergrant *domain.UserGrant, resourceOwner string) (_ *domain.UserGrant, err error) {
	event, addedUserGrant, err := r.addUserGrant(ctx, usergrant, resourceOwner)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := r.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(addedUserGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return userGrantWriteModelToUserGrant(addedUserGrant), nil
}

func (r *CommandSide) addUserGrant(ctx context.Context, userGrant *domain.UserGrant, resourceOwner string) (pusher eventstore.EventPusher, _ *UserGrantWriteModel, err error) {
	err = checkExplicitProjectPermission(ctx, userGrant.ProjectGrantID, userGrant.ProjectID)
	if err != nil {
		return nil, nil, err
	}
	if !userGrant.IsValid() {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0fs", "Errors.UserGrant.Invalid")
	}
	err = r.checkUserExists(ctx, userGrant.UserID, "")
	if err != nil {
		return nil, nil, err
	}
	err = r.checkProjectExists(ctx, userGrant.ProjectID, resourceOwner)
	if err != nil {
		return nil, nil, err
	}
	userGrant.AggregateID, err = r.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}

	addedUserGrant := NewUserGrantWriteModel(userGrant.AggregateID, resourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&addedUserGrant.WriteModel)
	pusher = usergrant.NewUserGrantAddedEvent(
		ctx,
		userGrantAgg,
		userGrant.UserID,
		userGrant.ProjectID,
		userGrant.ProjectGrantID,
		userGrant.RoleKeys,
	)
	return pusher, addedUserGrant, nil
}

func (r *CommandSide) ChangeUserGrant(ctx context.Context, userGrant *domain.UserGrant, resourceOwner string) (_ *domain.UserGrant, err error) {
	event, changedUserGrant, err := r.changeUserGrant(ctx, userGrant, resourceOwner, false)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := r.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(changedUserGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return userGrantWriteModelToUserGrant(changedUserGrant), nil
}

func (r *CommandSide) changeUserGrant(ctx context.Context, userGrant *domain.UserGrant, resourceOwner string, cascade bool) (_ eventstore.EventPusher, _ *UserGrantWriteModel, err error) {
	err = checkExplicitProjectPermission(ctx, userGrant.ProjectGrantID, userGrant.ProjectID)
	if err != nil {
		return nil, nil, err
	}
	if userGrant.IsValid() {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M0sd", "Errors.UserGrant.Invalid")
	}

	existingUserGrant, err := r.userGrantWriteModelByID(ctx, userGrant.AggregateID, userGrant.ResourceOwner)
	if err != nil {
		return nil, nil, err
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return nil, nil, caos_errs.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.UserGrant.NotFound")
	}
	if reflect.DeepEqual(existingUserGrant.RoleKeys, userGrant.RoleKeys) {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Rs8fy", "Errors.UserGrant.NotChanged")
	}

	changedUserGrant := NewUserGrantWriteModel(userGrant.AggregateID, resourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&changedUserGrant.WriteModel)

	if cascade {
		return usergrant.NewUserGrantCascadeChangedEvent(ctx, userGrantAgg, userGrant.RoleKeys), existingUserGrant, nil
	}
	return usergrant.NewUserGrantChangedEvent(ctx, userGrantAgg, userGrant.RoleKeys), existingUserGrant, nil
}

func (r *CommandSide) removeRoleFromUserGrant(ctx context.Context, userGrantID string, roleKeys []string, cascade bool) (_ eventstore.EventPusher, err error) {
	existingUserGrant, err := r.userGrantWriteModelByID(ctx, userGrantID, "")
	if err != nil {
		return nil, err
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.UserGrant.NotFound")
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5m8g9", "Errors.UserGrant.RoleKeyNotFound")
	}
	changedUserGrant := NewUserGrantWriteModel(userGrantID, "")
	userGrantAgg := UserGrantAggregateFromWriteModel(&changedUserGrant.WriteModel)

	if cascade {
		return usergrant.NewUserGrantCascadeChangedEvent(ctx, userGrantAgg, existingUserGrant.RoleKeys), nil
	}

	return usergrant.NewUserGrantChangedEvent(ctx, userGrantAgg, existingUserGrant.RoleKeys), nil
}

func (r *CommandSide) DeactivateUserGrant(ctx context.Context, grantID, resourceOwner string) (err error) {
	if grantID == "" || resourceOwner == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-M0dsf", "Errors.UserGrant.IDMissing")
	}

	existingUserGrant, err := r.userGrantWriteModelByID(ctx, grantID, resourceOwner)
	if err != nil {
		return err
	}
	err = checkExplicitProjectPermission(ctx, existingUserGrant.ProjectGrantID, existingUserGrant.ProjectID)
	if err != nil {
		return err
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.UserGrant.NotFound")
	}
	if existingUserGrant.State != domain.UserGrantStateActive {
		return caos_errs.ThrowNotFound(nil, "COMMAND-1S9gx", "Errors.UserGrant.NotActive")
	}

	deactivateUserGrant := NewUserGrantWriteModel(grantID, resourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&deactivateUserGrant.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, usergrant.NewUserGrantDeactivatedEvent(ctx, userGrantAgg))
	return err
}

func (r *CommandSide) ReactivateUserGrant(ctx context.Context, grantID, resourceOwner string) (err error) {
	if grantID == "" || resourceOwner == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Qxy8v", "Errors.UserGrant.IDMissing")
	}

	existingUserGrant, err := r.userGrantWriteModelByID(ctx, grantID, resourceOwner)
	if err != nil {
		return err
	}
	err = checkExplicitProjectPermission(ctx, existingUserGrant.ProjectGrantID, existingUserGrant.ProjectID)
	if err != nil {
		return err
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-Lp0gs", "Errors.UserGrant.NotFound")
	}
	if existingUserGrant.State != domain.UserGrantStateInactive {
		return caos_errs.ThrowNotFound(nil, "COMMAND-1ML0v", "Errors.UserGrant.NotInactive")
	}

	deactivateUserGrant := NewUserGrantWriteModel(grantID, resourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&deactivateUserGrant.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, usergrant.NewUserGrantReactivatedEvent(ctx, userGrantAgg))
	return err
}

func (r *CommandSide) RemoveUserGrant(ctx context.Context, grantID, resourceOwner string) (err error) {
	event, err := r.removeUserGrant(ctx, grantID, resourceOwner, false)
	if err != nil {
		return nil
	}

	_, err = r.eventstore.PushEvents(ctx, event)
	return err
}

func (r *CommandSide) BulkRemoveUserGrant(ctx context.Context, grantIDs []string, resourceOwner string) (err error) {
	events := make([]eventstore.EventPusher, len(grantIDs))
	for i, grantID := range grantIDs {
		event, err := r.removeUserGrant(ctx, grantID, resourceOwner, false)
		if err != nil {
			return nil
		}
		events[i] = event
	}
	_, err = r.eventstore.PushEvents(ctx, events...)
	return err
}

func (r *CommandSide) removeUserGrant(ctx context.Context, grantID, resourceOwner string, cascade bool) (_ eventstore.EventPusher, err error) {
	if grantID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-J9sc5", "Errors.UserGrant.IDMissing")
	}

	existingUserGrant, err := r.userGrantWriteModelByID(ctx, grantID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !cascade {
		err = checkExplicitProjectPermission(ctx, existingUserGrant.ProjectGrantID, existingUserGrant.ProjectID)
		if err != nil {
			return nil, err
		}
	}

	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-1My0t", "Errors.UserGrant.NotFound")
	}

	removeUserGrant := NewUserGrantWriteModel(grantID, resourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&removeUserGrant.WriteModel)
	if cascade {
		return usergrant.NewUserGrantCascadeRemovedEvent(
			ctx,
			userGrantAgg,
			existingUserGrant.UserID,
			existingUserGrant.ProjectID,
			existingUserGrant.ProjectGrantID), nil
	}
	return usergrant.NewUserGrantRemovedEvent(
		ctx,
		userGrantAgg,
		existingUserGrant.UserID,
		existingUserGrant.ProjectID,
		existingUserGrant.ProjectGrantID), nil
}
func (r *CommandSide) userGrantWriteModelByID(ctx context.Context, userGrantID, resourceOwner string) (writeModel *UserGrantWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserGrantWriteModel(userGrantID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
