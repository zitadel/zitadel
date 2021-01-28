package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/usergrant"
	"reflect"
)

func (r *CommandSide) AddUserGrant(ctx context.Context, usergrant *domain.UserGrant, resourceOwner string) (_ *domain.UserGrant, err error) {
	userGrantAgg, addedUserGrant, err := r.addUserGrant(ctx, usergrant, resourceOwner)
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, addedUserGrant, userGrantAgg)
	if err != nil {
		return nil, err
	}

	return userGrantWriteModelToUserGrant(addedUserGrant), nil
}

func (r *CommandSide) addUserGrant(ctx context.Context, userGrant *domain.UserGrant, resourceOwner string) (_ *usergrant.Aggregate, _ *UserGrantWriteModel, err error) {
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

	userGrantAgg.PushEvents(
		usergrant.NewUserGrantAddedEvent(
			ctx,
			resourceOwner,
			userGrant.UserID,
			userGrant.ProjectID,
			userGrant.ProjectGrantID,
			userGrant.RoleKeys,
		),
	)
	return userGrantAgg, addedUserGrant, nil
}

func (r *CommandSide) ChangeUserGrant(ctx context.Context, userGrant *domain.UserGrant, resourceOwner string) (_ *domain.UserGrant, err error) {
	userGrantAgg, addedUserGrant, err := r.changeUserGrant(ctx, userGrant, resourceOwner, false)
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, addedUserGrant, userGrantAgg)
	if err != nil {
		return nil, err
	}

	return userGrantWriteModelToUserGrant(addedUserGrant), nil
}

func (r *CommandSide) changeUserGrant(ctx context.Context, userGrant *domain.UserGrant, resourceOwner string, cascade bool) (_ *usergrant.Aggregate, _ *UserGrantWriteModel, err error) {
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

	if !cascade {
		userGrantAgg.PushEvents(
			usergrant.NewUserGrantChangedEvent(ctx, userGrant.RoleKeys),
		)
	} else {
		userGrantAgg.PushEvents(
			usergrant.NewUserGrantCascadeChangedEvent(ctx, userGrant.RoleKeys),
		)
	}

	return userGrantAgg, changedUserGrant, nil
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
	userGrantAgg.PushEvents(
		usergrant.NewUserGrantDeactivatedEvent(ctx),
	)

	return r.eventstore.PushAggregate(ctx, deactivateUserGrant, userGrantAgg)
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
	userGrantAgg.PushEvents(
		usergrant.NewUserGrantReactivatedEvent(ctx),
	)

	return r.eventstore.PushAggregate(ctx, deactivateUserGrant, userGrantAgg)
}

func (r *CommandSide) RemoveUserGrant(ctx context.Context, grantID, resourceOwner string) (err error) {
	userGrantAgg, removeUserGrant, err := r.removeUserGrant(ctx, grantID, resourceOwner, false)
	if err != nil {
		return nil
	}

	return r.eventstore.PushAggregate(ctx, removeUserGrant, userGrantAgg)
}

func (r *CommandSide) BulkRemoveUserGrant(ctx context.Context, grantIDs []string, resourceOwner string) (err error) {
	aggregates := make([]eventstore.Aggregater, len(grantIDs))
	for i, grantID := range grantIDs {
		userGrantAgg, _, err := r.removeUserGrant(ctx, grantID, resourceOwner, false)
		if err != nil {
			return nil
		}
		aggregates[i] = userGrantAgg
	}
	_, err = r.eventstore.PushAggregates(ctx, aggregates...)
	return err
}

func (r *CommandSide) removeUserGrant(ctx context.Context, grantID, resourceOwner string, cascade bool) (_ *usergrant.Aggregate, _ *UserGrantWriteModel, err error) {
	if grantID == "" || resourceOwner == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-J9sc5", "Errors.UserGrant.IDMissing")
	}

	existingUserGrant, err := r.userGrantWriteModelByID(ctx, grantID, resourceOwner)
	if err != nil {
		return nil, nil, err
	}
	err = checkExplicitProjectPermission(ctx, existingUserGrant.ProjectGrantID, existingUserGrant.ProjectID)
	if err != nil {
		return nil, nil, err
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return nil, nil, caos_errs.ThrowNotFound(nil, "COMMAND-1My0t", "Errors.UserGrant.NotFound")
	}

	//TODO: Remove Uniqueness
	removeUserGrant := NewUserGrantWriteModel(grantID, resourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&removeUserGrant.WriteModel)
	if !cascade {
		userGrantAgg.PushEvents(
			usergrant.NewUserGrantRemovedEvent(ctx, existingUserGrant.ResourceOwner, existingUserGrant.UserID, existingUserGrant.ProjectID),
		)
	} else {
		userGrantAgg.PushEvents(
			usergrant.NewUserGrantCascadeRemovedEvent(ctx, existingUserGrant.ResourceOwner, existingUserGrant.UserID, existingUserGrant.ProjectID),
		)
	}

	return userGrantAgg, removeUserGrant, nil
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
