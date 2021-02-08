package command

import (
	"context"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
	"github.com/caos/zitadel/internal/v2/repository/usergrant"
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

	projectAgg.PushEvents(project.NewGrantAddedEvent(ctx, grant.GrantID, grant.GrantedOrgID, grant.AggregateID, grant.RoleKeys))

	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, addedGrant, projectAgg)
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

	projectAgg.PushEvents(project.NewGrantChangedEvent(ctx, grant.GrantID, grant.RoleKeys))

	removedRoles := domain.GetRemovedRoles(existingGrant.RoleKeys, grant.RoleKeys)
	if len(removedRoles) == 0 {
		err = r.eventstore.PushAggregate(ctx, existingGrant, projectAgg)
		if err != nil {
			return nil, err
		}
		return projectGrantWriteModelToProjectGrant(existingGrant), nil
	}

	aggregates := make([]eventstore.Aggregater, 0)
	aggregates = append(aggregates, projectAgg)
	for _, userGrantID := range cascadeUserGrantIDs {
		grantAgg, _, err := r.removeRoleFromUserGrant(ctx, userGrantID, removedRoles, true)
		if err != nil {
			continue
		}
		aggregates = append(aggregates, grantAgg)
	}
	resultEvents, err := r.eventstore.PushAggregates(ctx, aggregates...)
	if err != nil {
		return nil, err
	}
	existingGrant.AppendEvents(resultEvents...)
	err = existingGrant.Reduce()
	if err != nil {
		return nil, err
	}
	return projectGrantWriteModelToProjectGrant(existingGrant), nil
}

func (r *CommandSide) removeRoleFromProjectGrant(ctx context.Context, projectAgg *project.Aggregate, projectID, projectGrantID, roleKey string, cascade bool) (_ *ProjectGrantWriteModel, err error) {
	existingProjectGrant, err := r.projectGrantWriteModelByID(ctx, projectID, projectGrantID, "")
	if err != nil {
		return nil, err
	}
	if existingProjectGrant.State == domain.ProjectGrantStateUnspecified || existingProjectGrant.State == domain.ProjectGrantStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.Grant.NotFound")
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5m8g9", "Errors.Project.Grant.RoleKeyNotFound")
	}
	changedProjectGrant := NewProjectGrantWriteModel(projectGrantID, projectID, existingProjectGrant.ResourceOwner)

	if !cascade {
		projectAgg.PushEvents(
			project.NewGrantChangedEvent(ctx, projectGrantID, existingProjectGrant.RoleKeys),
		)
	} else {
		projectAgg.PushEvents(
			usergrant.NewUserGrantCascadeChangedEvent(ctx, existingProjectGrant.RoleKeys),
		)
	}

	return changedProjectGrant, nil
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

	projectAgg.PushEvents(project.NewGrantDeactivateEvent(ctx, grantID))
	return r.eventstore.PushAggregate(ctx, existingGrant, projectAgg)
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
	projectAgg.PushEvents(project.NewGrantReactivatedEvent(ctx, grantID))
	return r.eventstore.PushAggregate(ctx, existingGrant, projectAgg)
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
	aggregates := make([]eventstore.Aggregater, 0)
	projectAgg := ProjectAggregateFromWriteModel(&existingGrant.WriteModel)
	projectAgg.PushEvents(project.NewGrantRemovedEvent(ctx, grantID, existingGrant.GrantedOrgID, projectID))
	aggregates = append(aggregates, projectAgg)

	for _, userGrantID := range cascadeUserGrantIDs {
		grantAgg, _, err := r.removeUserGrant(ctx, userGrantID, "", true)
		if err != nil {
			logging.LogWithFields("COMMAND-3m8sG", "usergrantid", grantID).WithError(err).Warn("could not cascade remove user grant")
			continue
		}
		aggregates = append(aggregates, grantAgg)
	}
	_, err = r.eventstore.PushAggregates(ctx, aggregates...)
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
