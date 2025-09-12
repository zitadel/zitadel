package command

import (
	"context"
	"reflect"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddProjectGrant struct {
	es_models.ObjectRoot

	GrantID      string
	GrantedOrgID string
	RoleKeys     []string
}

func (p *AddProjectGrant) IsValid() error {
	if p.AggregateID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-FYRnWEzBzV", "Errors.Project.Grant.Invalid")
	}
	if p.GrantedOrgID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-PPhHpWGRAE", "Errors.Project.Grant.Invalid")
	}
	return nil
}

func (c *Commands) AddProjectGrant(ctx context.Context, grant *AddProjectGrant) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err := grant.IsValid(); err != nil {
		return nil, err
	}

	if grant.GrantID == "" {
		grant.GrantID, err = c.idGenerator.Next()
		if err != nil {
			return nil, err
		}
	}

	projectResourceOwner, err := c.checkProjectGrantPreCondition(ctx, grant.AggregateID, grant.GrantedOrgID, grant.ResourceOwner, grant.RoleKeys)
	if err != nil {
		return nil, err
	}
	// if there is no resourceowner provided then use the resourceowner of the project
	if grant.ResourceOwner == "" {
		grant.ResourceOwner = projectResourceOwner
	}
	if err := c.checkPermissionUpdateProjectGrant(ctx, grant.ResourceOwner, grant.AggregateID, grant.GrantID); err != nil {
		return nil, err
	}

	wm := NewProjectGrantWriteModel(grant.GrantID, grant.GrantedOrgID, grant.AggregateID, grant.ResourceOwner)
	// error if provided resourceowner is not equal to the resourceowner of the project or the project grant is for the same organization
	if projectResourceOwner != wm.ResourceOwner || wm.ResourceOwner == grant.GrantedOrgID {
		return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-ckUpbvboAH", "Errors.Project.Grant.Invalid")
	}
	if err := c.pushAppendAndReduce(ctx,
		wm,
		project.NewGrantAddedEvent(ctx,
			ProjectAggregateFromWriteModelWithCTX(ctx, &wm.WriteModel),
			grant.GrantID,
			grant.GrantedOrgID,
			grant.RoleKeys),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

type ChangeProjectGrant struct {
	es_models.ObjectRoot

	GrantID      string
	GrantedOrgID string
	RoleKeys     []string
}

func (c *Commands) ChangeProjectGrant(ctx context.Context, grant *ChangeProjectGrant, cascadeUserGrantIDs ...string) (_ *domain.ObjectDetails, err error) {
	if grant.GrantID == "" && grant.GrantedOrgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-1j83s", "Errors.IDMissing")
	}
	existingGrant, err := c.projectGrantWriteModelByID(ctx, grant.GrantID, grant.GrantedOrgID, grant.AggregateID, grant.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingGrant.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.Project.Grant.NotFound")
	}

	if err := c.checkPermissionUpdateProjectGrant(ctx, existingGrant.ResourceOwner, existingGrant.AggregateID, existingGrant.GrantID); err != nil {
		return nil, err
	}
	projectResourceOwner, err := c.checkProjectGrantPreCondition(ctx, existingGrant.AggregateID, existingGrant.GrantedOrgID, existingGrant.ResourceOwner, grant.RoleKeys)
	if err != nil {
		return nil, err
	}
	// error if provided resourceowner is not equal to the resourceowner of the project
	if existingGrant.ResourceOwner != projectResourceOwner {
		return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-q1BhA68RBC", "Errors.Project.Grant.Invalid")
	}

	// return if there are no changes to the project grant roles
	if reflect.DeepEqual(existingGrant.RoleKeys, grant.RoleKeys) {
		return writeModelToObjectDetails(&existingGrant.WriteModel), nil
	}

	events := []eventstore.Command{
		project.NewGrantChangedEvent(ctx,
			ProjectAggregateFromWriteModelWithCTX(ctx, &existingGrant.WriteModel),
			existingGrant.GrantID,
			grant.RoleKeys,
		),
	}

	removedRoles := domain.GetRemovedRoles(existingGrant.RoleKeys, grant.RoleKeys)
	if len(removedRoles) == 0 {
		pushedEvents, err := c.eventstore.Push(ctx, events...)
		if err != nil {
			return nil, err
		}
		err = AppendAndReduce(existingGrant, pushedEvents...)
		if err != nil {
			return nil, err
		}
		return writeModelToObjectDetails(&existingGrant.WriteModel), nil
	}

	for _, userGrantID := range cascadeUserGrantIDs {
		event, err := c.removeRoleFromUserGrant(ctx, userGrantID, removedRoles, true)
		if err != nil {
			continue
		}
		events = append(events, event)
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGrant.WriteModel), nil
}

func (c *Commands) removeRoleFromProjectGrant(ctx context.Context, projectAgg *eventstore.Aggregate, projectID, projectGrantID, roleKey string, cascade bool) (_ eventstore.Command, _ *ProjectGrantWriteModel, err error) {
	existingProjectGrant, err := c.projectGrantWriteModelByID(ctx, projectGrantID, "", projectID, "")
	if err != nil {
		return nil, nil, err
	}
	if !existingProjectGrant.State.Exists() {
		return nil, nil, zerrors.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.Project.Grant.NotFound")
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
		return nil, nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-5m8g9", "Errors.Project.Grant.RoleKeyNotFound")
	}
	changedProjectGrant := NewProjectGrantWriteModel(projectGrantID, projectID, "", existingProjectGrant.ResourceOwner)

	if cascade {
		return project.NewGrantCascadeChangedEvent(ctx, projectAgg, projectGrantID, existingProjectGrant.RoleKeys), changedProjectGrant, nil
	}

	return project.NewGrantChangedEvent(ctx, projectAgg, projectGrantID, existingProjectGrant.RoleKeys), changedProjectGrant, nil
}

func (c *Commands) DeactivateProjectGrant(ctx context.Context, projectID, grantID, grantedOrgID, resourceOwner string) (details *domain.ObjectDetails, err error) {
	if (grantID == "" && grantedOrgID == "") || projectID == "" {
		return details, zerrors.ThrowInvalidArgument(nil, "PROJECT-p0s4V", "Errors.IDMissing")
	}

	projectResourceOwner, err := c.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	existingGrant, err := c.projectGrantWriteModelByID(ctx, grantID, grantedOrgID, projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	if !existingGrant.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.Project.Grant.NotFound")
	}
	// error if provided resourceowner is not equal to the resourceowner of the project
	if projectResourceOwner != existingGrant.ResourceOwner {
		return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-0l10S9OmZV", "Errors.Project.Grant.Invalid")
	}
	// return if project grant is already inactive
	if existingGrant.State == domain.ProjectGrantStateInactive {
		return writeModelToObjectDetails(&existingGrant.WriteModel), nil
	}
	// error if project grant is neither active nor inactive
	if existingGrant.State != domain.ProjectGrantStateActive {
		return details, zerrors.ThrowPreconditionFailed(nil, "PROJECT-47fu8", "Errors.Project.Grant.NotActive")
	}
	if err := c.checkPermissionUpdateProjectGrant(ctx, existingGrant.ResourceOwner, existingGrant.AggregateID, existingGrant.GrantID); err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx,
		project.NewGrantDeactivateEvent(ctx,
			ProjectAggregateFromWriteModelWithCTX(ctx, &existingGrant.WriteModel),
			existingGrant.GrantID,
		),
	)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGrant.WriteModel), nil
}

func (c *Commands) checkProjectGrantExists(ctx context.Context, grantID, grantedOrgID, projectID, resourceOwner string) (string, string, error) {
	existingGrant, err := c.projectGrantWriteModelByID(ctx, grantID, grantedOrgID, projectID, resourceOwner)
	if err != nil {
		return "", "", err
	}
	if !existingGrant.State.Exists() {
		return "", "", zerrors.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.Project.Grant.NotFound")
	}
	return existingGrant.GrantedOrgID, existingGrant.ResourceOwner, nil
}

func (c *Commands) ReactivateProjectGrant(ctx context.Context, projectID, grantID, grantedOrgID, resourceOwner string) (details *domain.ObjectDetails, err error) {
	if (grantID == "" && grantedOrgID == "") || projectID == "" {
		return details, zerrors.ThrowInvalidArgument(nil, "PROJECT-p0s4V", "Errors.IDMissing")
	}

	projectResourceOwner, err := c.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	existingGrant, err := c.projectGrantWriteModelByID(ctx, grantID, grantedOrgID, projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	if !existingGrant.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.Project.Grant.NotFound")
	}
	// error if provided resourceowner is not equal to the resourceowner of the project
	if projectResourceOwner != existingGrant.ResourceOwner {
		return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-byscAarAST", "Errors.Project.Grant.Invalid")
	}
	// return if project grant is already active
	if existingGrant.State == domain.ProjectGrantStateActive {
		return writeModelToObjectDetails(&existingGrant.WriteModel), nil
	}
	// error if project grant is neither active nor inactive
	if existingGrant.State != domain.ProjectGrantStateInactive {
		return details, zerrors.ThrowPreconditionFailed(nil, "PROJECT-47fu8", "Errors.Project.Grant.NotInactive")
	}
	if err := c.checkPermissionUpdateProjectGrant(ctx, existingGrant.ResourceOwner, existingGrant.AggregateID, existingGrant.GrantID); err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx,
		project.NewGrantReactivatedEvent(ctx,
			ProjectAggregateFromWriteModelWithCTX(ctx, &existingGrant.WriteModel),
			existingGrant.GrantID,
		),
	)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGrant.WriteModel), nil
}

// Deprecated: use commands.DeleteProjectGrant
func (c *Commands) RemoveProjectGrant(ctx context.Context, projectID, grantID, resourceOwner string, cascadeUserGrantIDs ...string) (details *domain.ObjectDetails, err error) {
	if grantID == "" || projectID == "" {
		return details, zerrors.ThrowInvalidArgument(nil, "PROJECT-1m9fJ", "Errors.IDMissing")
	}
	existingGrant, err := c.projectGrantWriteModelByID(ctx, grantID, "", projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	if !existingGrant.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.Project.Grant.NotFound")
	}
	if err := c.checkPermissionDeleteProjectGrant(ctx, existingGrant.ResourceOwner, existingGrant.AggregateID, existingGrant.GrantID); err != nil {
		return nil, err
	}
	events := make([]eventstore.Command, 0)
	events = append(events, project.NewGrantRemovedEvent(ctx,
		ProjectAggregateFromWriteModelWithCTX(ctx, &existingGrant.WriteModel),
		existingGrant.GrantID,
		existingGrant.GrantedOrgID,
	))

	for _, userGrantID := range cascadeUserGrantIDs {
		event, _, err := c.removeUserGrant(ctx, userGrantID, "", true, true, nil)
		if err != nil {
			logging.WithFields("id", "COMMAND-3m8sG", "usergrantid", grantID).WithError(err).Warn("could not cascade remove user grant")
			continue
		}
		if event != nil {
			events = append(events, event)
		}
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGrant.WriteModel), nil
}

func (c *Commands) DeleteProjectGrant(ctx context.Context, projectID, grantID, grantedOrgID, resourceOwner string, cascadeUserGrantIDs ...string) (details *domain.ObjectDetails, err error) {
	if (grantID == "" && grantedOrgID == "") || projectID == "" {
		return details, zerrors.ThrowInvalidArgument(nil, "PROJECT-1m9fJ", "Errors.IDMissing")
	}
	existingGrant, err := c.projectGrantWriteModelByID(ctx, grantID, grantedOrgID, projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	// return if project grant does not exist, or was removed already
	if !existingGrant.State.Exists() {
		return writeModelToObjectDetails(&existingGrant.WriteModel), nil
	}
	if err := c.checkPermissionDeleteProjectGrant(ctx, existingGrant.ResourceOwner, existingGrant.AggregateID, existingGrant.GrantID); err != nil {
		return nil, err
	}
	events := make([]eventstore.Command, 0)
	events = append(events, project.NewGrantRemovedEvent(ctx,
		ProjectAggregateFromWriteModelWithCTX(ctx, &existingGrant.WriteModel),
		existingGrant.GrantID,
		existingGrant.GrantedOrgID,
	),
	)

	for _, userGrantID := range cascadeUserGrantIDs {
		event, _, err := c.removeUserGrant(ctx, userGrantID, "", true, true, nil)
		if err != nil {
			logging.WithFields("id", "COMMAND-3m8sG", "usergrantid", grantID).WithError(err).Warn("could not cascade remove user grant")
			continue
		}
		if event != nil {
			events = append(events, event)
		}
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGrant.WriteModel), nil
}

func (c *Commands) projectGrantWriteModelByID(ctx context.Context, grantID, grantedOrgID, projectID, resourceOwner string) (member *ProjectGrantWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewProjectGrantWriteModel(grantID, grantedOrgID, projectID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) checkProjectGrantPreCondition(ctx context.Context, projectID, grantedOrgID, resourceOwner string, roles []string) (string, error) {
	if !authz.GetFeatures(ctx).ShouldUseImprovedPerformance(feature.ImprovedPerformanceTypeProjectGrant) {
		return c.checkProjectGrantPreConditionOld(ctx, projectID, grantedOrgID, resourceOwner, roles)
	}
	projectResourceOwner, existingRoleKeys, err := c.searchProjectGrantState(ctx, projectID, grantedOrgID, resourceOwner)
	if err != nil {
		return "", err
	}

	if domain.HasInvalidRoles(existingRoleKeys, roles) {
		return "", zerrors.ThrowPreconditionFailed(err, "COMMAND-6m9gd", "Errors.Project.Role.NotFound")
	}
	return projectResourceOwner, nil
}

func (c *Commands) searchProjectGrantState(ctx context.Context, projectID, grantedOrgID, resourceOwner string) (_ string, existingRoleKeys []string, err error) {
	projectStateQuery := map[eventstore.FieldType]any{
		eventstore.FieldTypeAggregateType: project.AggregateType,
		eventstore.FieldTypeAggregateID:   projectID,
		eventstore.FieldTypeFieldName:     project.ProjectStateSearchField,
		eventstore.FieldTypeObjectType:    project.ProjectSearchType,
	}
	grantedOrgQuery := map[eventstore.FieldType]any{
		eventstore.FieldTypeAggregateType: org.AggregateType,
		eventstore.FieldTypeAggregateID:   grantedOrgID,
		eventstore.FieldTypeFieldName:     org.OrgStateSearchField,
		eventstore.FieldTypeObjectType:    org.OrgSearchType,
	}
	roleQuery := map[eventstore.FieldType]any{
		eventstore.FieldTypeAggregateType: project.AggregateType,
		eventstore.FieldTypeAggregateID:   projectID,
		eventstore.FieldTypeFieldName:     project.ProjectRoleKeySearchField,
		eventstore.FieldTypeObjectType:    project.ProjectRoleSearchType,
	}

	// as resourceowner is not always provided, it has to be separately
	if resourceOwner != "" {
		projectStateQuery[eventstore.FieldTypeResourceOwner] = resourceOwner
		roleQuery[eventstore.FieldTypeResourceOwner] = resourceOwner
	}

	results, err := c.eventstore.Search(
		ctx,
		projectStateQuery,
		grantedOrgQuery,
		roleQuery,
	)
	if err != nil {
		return "", nil, err
	}

	var (
		existsProject                bool
		existingProjectResourceOwner string
		existsGrantedOrg             bool
	)

	for _, result := range results {
		switch result.Object.Type {
		case project.ProjectRoleSearchType:
			var role string
			err := result.Value.Unmarshal(&role)
			if err != nil {
				return "", nil, err
			}
			existingRoleKeys = append(existingRoleKeys, role)
		case org.OrgSearchType:
			var state domain.OrgState
			err := result.Value.Unmarshal(&state)
			if err != nil {
				return "", nil, err
			}
			existsGrantedOrg = state.Valid() && state != domain.OrgStateRemoved
		case project.ProjectSearchType:
			var state domain.ProjectState
			err := result.Value.Unmarshal(&state)
			if err != nil {
				return "", nil, err
			}
			existsProject = state.Valid() && state != domain.ProjectStateRemoved
			existingProjectResourceOwner = result.Aggregate.ResourceOwner
		}
	}

	if !existsProject {
		return "", nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-m9gsd", "Errors.Project.NotFound")
	}
	if !existsGrantedOrg {
		return "", nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-3m9gg", "Errors.Org.NotFound")
	}
	return existingProjectResourceOwner, existingRoleKeys, nil
}
