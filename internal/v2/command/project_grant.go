package command

import (
	"context"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
	"reflect"
)

func (r *CommandSide) AddProjectGrant(ctx context.Context, grant *domain.ProjectGrant, resourceOwner string) (_ *domain.ProjectGrant, err error) {
	if !grant.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-Bff2g", "Errors.Project.Grant.Invalid")
	}
	grant.GrantID, err = r.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	err = r.checkProjectExists(ctx, grant.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	err = r.checkOrgExists(ctx, grant.GrantedOrgID)
	if err != nil {
		return nil, err
	}
	addedGrant := NewProjectGrantWriteModel(grant.GrantID, grant.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedGrant.WriteModel)
	pushedEvents, err := r.eventstore.PushEvents(
		ctx,
		project.NewGrantAddedEvent(ctx, projectAgg, grant.GrantID, grant.GrantedOrgID, grant.AggregateID, grant.RoleKeys))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return projectGrantWriteModelToProjectGrant(addedGrant), nil
}

func (r *CommandSide) ChangeProjectGrant(ctx context.Context, grant *domain.ProjectGrant, resourceOwner string, cascadeUserGrantIDs ...string) (_ *domain.ProjectGrant, err error) {
	if grant.GrantID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-1j83s", "Errors.IDMissing")
	}
	err = r.checkProjectExists(ctx, grant.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	existingGrant, err := r.projectGrantWriteModelByID(ctx, grant.GrantID, grant.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingGrant.WriteModel)

	if reflect.DeepEqual(existingGrant.RoleKeys, grant.RoleKeys) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-0o0pL", "Errors.NoChangesFoundc")
	}

	events := []eventstore.EventPusher{
		project.NewGrantChangedEvent(ctx, projectAgg, grant.GrantID, grant.RoleKeys),
	}

	removedRoles := domain.GetRemovedRoles(existingGrant.RoleKeys, grant.RoleKeys)
	if len(removedRoles) == 0 {
		pushedEvents, err := r.eventstore.PushEvents(ctx, events...)
		if err != nil {
			return nil, err
		}
		err = AppendAndReduce(existingGrant, pushedEvents...)
		if err != nil {
			return nil, err
		}
		return projectGrantWriteModelToProjectGrant(existingGrant), nil
	}

	for _, userGrantID := range cascadeUserGrantIDs {
		event, err := r.removeRoleFromUserGrant(ctx, userGrantID, removedRoles, true)
		if err != nil {
			continue
		}
		events = append(events, event)
	}
	pushedEvents, err := r.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return projectGrantWriteModelToProjectGrant(existingGrant), nil
}

func (r *CommandSide) removeRoleFromProjectGrant(ctx context.Context, projectAgg *eventstore.Aggregate, projectID, projectGrantID, roleKey string, cascade bool) (_ eventstore.EventPusher, _ *ProjectGrantWriteModel, err error) {
	existingProjectGrant, err := r.projectGrantWriteModelByID(ctx, projectID, projectGrantID, "")
	if err != nil {
		return nil, nil, err
	}
	if existingProjectGrant.State == domain.ProjectGrantStateUnspecified || existingProjectGrant.State == domain.ProjectGrantStateRemoved {
		return nil, nil, caos_errs.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.Grant.NotFound")
	}
	keyExists := false
	for i, key := range existingProjectGrant.RoleKeys {
		if key == roleKey {
			keyExists = true
			copy(existingProjectGrant.RoleKeys[i:], existingProjectGrant.RoleKeys[i+1:])
			existingProjectGrant.RoleKeys[len(existingProjectGrant.RoleKeys)-1] = ""
			existingProjectGrant.RoleKeys = existingProjectGrant.RoleKeys[:len(existingProjectGrant.RoleKeys)-1]
			continue
		}
	}
	if !keyExists {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5m8g9", "Errors.Project.Grant.RoleKeyNotFound")
	}
	changedProjectGrant := NewProjectGrantWriteModel(projectGrantID, projectID, existingProjectGrant.ResourceOwner)

	if cascade {
		return project.NewGrantCascadeChangedEvent(ctx, projectAgg, projectGrantID, existingProjectGrant.RoleKeys), changedProjectGrant, nil
	}

	return project.NewGrantChangedEvent(ctx, projectAgg, projectGrantID, existingProjectGrant.RoleKeys), changedProjectGrant, nil
}

func (r *CommandSide) DeactivateProjectGrant(ctx context.Context, projectID, grantID, resourceOwner string) (err error) {
	if grantID == "" || projectID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "PROJECT-p0s4V", "Errors.IDMissing")
	}
	err = r.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return err
	}
	existingGrant, err := r.projectGrantWriteModelByID(ctx, grantID, projectID, resourceOwner)
	if err != nil {
		return err
	}
	if existingGrant.State != domain.ProjectGrantStateActive {
		return caos_errs.ThrowPreconditionFailed(nil, "PROJECT-47fu8", "Errors.Project.Grant.NotActive")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingGrant.WriteModel)

	_, err = r.eventstore.PushEvents(ctx, project.NewGrantDeactivateEvent(ctx, projectAgg, grantID))
	return err
}

func (r *CommandSide) ReactivateProjectGrant(ctx context.Context, projectID, grantID, resourceOwner string) (err error) {
	if grantID == "" || projectID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "PROJECT-p0s4V", "Errors.IDMissing")
	}
	err = r.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return err
	}
	existingGrant, err := r.projectGrantWriteModelByID(ctx, grantID, projectID, resourceOwner)
	if err != nil {
		return err
	}
	if existingGrant.State != domain.ProjectGrantStateInactive {
		return caos_errs.ThrowPreconditionFailed(nil, "PROJECT-47fu8", "Errors.Project.Grant.NotInactive")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingGrant.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, project.NewGrantReactivatedEvent(ctx, projectAgg, grantID))
	return err
}

func (r *CommandSide) RemoveProjectGrant(ctx context.Context, projectID, grantID, resourceOwner string, cascadeUserGrantIDs ...string) (err error) {
	if grantID == "" || projectID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "PROJECT-1m9fJ", "Errors.IDMissing")
	}
	err = r.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return err
	}
	existingGrant, err := r.projectGrantWriteModelByID(ctx, grantID, projectID, resourceOwner)
	if err != nil {
		return err
	}
	events := make([]eventstore.EventPusher, 0)
	projectAgg := ProjectAggregateFromWriteModel(&existingGrant.WriteModel)
	events = append(events, project.NewGrantRemovedEvent(ctx, projectAgg, grantID, existingGrant.GrantedOrgID, projectID))

	for _, userGrantID := range cascadeUserGrantIDs {
		event, err := r.removeUserGrant(ctx, userGrantID, "", true)
		if err != nil {
			logging.LogWithFields("COMMAND-3m8sG", "usergrantid", grantID).WithError(err).Warn("could not cascade remove user grant")
			continue
		}
		events = append(events, event)
	}
	_, err = r.eventstore.PushEvents(ctx, events...)
	return err
}

func (r *CommandSide) projectGrantWriteModelByID(ctx context.Context, grantID, projectID, resourceOwner string) (member *ProjectGrantWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewProjectGrantWriteModel(grantID, projectID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.ProjectGrantStateUnspecified || writeModel.State == domain.ProjectGrantStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.Project.Grant.NotFound")
	}

	return writeModel, nil
}
