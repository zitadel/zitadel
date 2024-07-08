package command

import (
	"context"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddProjectWithID(ctx context.Context, project *domain.Project, resourceOwner, projectID string) (_ *domain.Project, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	existingProject, err := c.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingProject.State != domain.ProjectStateUnspecified {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-opamwu", "Errors.Project.AlreadyExisting")
	}
	return c.addProjectWithID(ctx, project, resourceOwner, projectID)
}

func (c *Commands) AddProject(ctx context.Context, project *domain.Project, resourceOwner, ownerUserID string) (_ *domain.Project, err error) {
	if !project.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-IOVCC", "Errors.Project.Invalid")
	}

	projectID, err := id_generator.Next()
	if err != nil {
		return nil, err
	}

	return c.addProjectWithIDWithOwner(ctx, project, resourceOwner, ownerUserID, projectID)
}

func (c *Commands) addProjectWithID(ctx context.Context, projectAdd *domain.Project, resourceOwner, projectID string) (_ *domain.Project, err error) {
	projectAdd.AggregateID = projectID
	addedProject := NewProjectWriteModel(projectAdd.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedProject.WriteModel)

	events := []eventstore.Command{
		project.NewProjectAddedEvent(
			ctx,
			projectAgg,
			projectAdd.Name,
			projectAdd.ProjectRoleAssertion,
			projectAdd.ProjectRoleCheck,
			projectAdd.HasProjectCheck,
			projectAdd.PrivateLabelingSetting),
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedProject, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return projectWriteModelToProject(addedProject), nil
}

func (c *Commands) addProjectWithIDWithOwner(ctx context.Context, projectAdd *domain.Project, resourceOwner, ownerUserID, projectID string) (_ *domain.Project, err error) {
	if !projectAdd.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-IOVCC", "Errors.Project.Invalid")
	}
	projectAdd.AggregateID = projectID
	addedProject := NewProjectWriteModel(projectAdd.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedProject.WriteModel)

	projectRole := domain.RoleProjectOwner
	events := []eventstore.Command{
		project.NewProjectAddedEvent(
			ctx,
			projectAgg,
			projectAdd.Name,
			projectAdd.ProjectRoleAssertion,
			projectAdd.ProjectRoleCheck,
			projectAdd.HasProjectCheck,
			projectAdd.PrivateLabelingSetting),
		project.NewProjectMemberAddedEvent(ctx, projectAgg, ownerUserID, projectRole),
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedProject, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return projectWriteModelToProject(addedProject), nil
}

func AddProjectCommand(
	a *project.Aggregate,
	name string,
	owner string,
	projectRoleAssertion bool,
	projectRoleCheck bool,
	hasProjectCheck bool,
	privateLabelingSetting domain.PrivateLabelingSetting,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if name = strings.TrimSpace(name); name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "PROJE-C01yo", "Errors.Invalid.Argument")
		}
		if !privateLabelingSetting.Valid() {
			return nil, zerrors.ThrowInvalidArgument(nil, "PROJE-AO52V", "Errors.Invalid.Argument")
		}
		if owner == "" {
			return nil, zerrors.ThrowPreconditionFailed(nil, "PROJE-hzxwo", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				project.NewProjectAddedEvent(ctx, &a.Aggregate,
					name,
					projectRoleAssertion,
					projectRoleCheck,
					hasProjectCheck,
					privateLabelingSetting,
				),
				project.NewProjectMemberAddedEvent(ctx, &a.Aggregate,
					owner,
					domain.RoleProjectOwner),
			}, nil
		}, nil
	}
}

func projectWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer, projectID, resourceOwner string) (project *ProjectWriteModel, err error) {
	project = NewProjectWriteModel(projectID, resourceOwner)
	events, err := filter(ctx, project.Query())
	if err != nil {
		return nil, err
	}

	project.AppendEvents(events...)
	if err := project.Reduce(); err != nil {
		return nil, err
	}

	return project, nil
}

func (c *Commands) getProjectByID(ctx context.Context, projectID, resourceOwner string) (_ *domain.Project, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	projectWriteModel, err := c.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if projectWriteModel.State == domain.ProjectStateUnspecified || projectWriteModel.State == domain.ProjectStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "PROJECT-Gd2hh", "Errors.Project.NotFound")
	}
	return projectWriteModelToProject(projectWriteModel), nil
}

func (c *Commands) projectAggregateByID(ctx context.Context, projectID, resourceOwner string) (*eventstore.Aggregate, domain.ProjectState, error) {
	result, err := c.projectState(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, domain.ProjectStateUnspecified, zerrors.ThrowNotFound(err, "COMMA-NDQoF", "Errors.Project.NotFound")
	}
	if len(result) == 0 {
		_ = projection.ProjectGrantFields.Trigger(ctx)
		result, err = c.projectState(ctx, projectID, resourceOwner)
		if err != nil || len(result) == 0 {
			return nil, domain.ProjectStateUnspecified, zerrors.ThrowNotFound(err, "COMMA-U1nza", "Errors.Project.NotFound")
		}
	}

	var state domain.ProjectState
	err = result[0].Value.Unmarshal(&state)
	if err != nil {
		return nil, state, zerrors.ThrowNotFound(err, "COMMA-o4n6F", "Errors.Project.NotFound")
	}
	return &result[0].Aggregate, state, nil
}

func (c *Commands) projectState(ctx context.Context, projectID, resourceOwner string) ([]*eventstore.SearchResult, error) {
	return c.eventstore.Search(
		ctx,
		map[eventstore.FieldType]any{
			eventstore.FieldTypeObjectType:     project.ProjectSearchType,
			eventstore.FieldTypeObjectID:       projectID,
			eventstore.FieldTypeObjectRevision: project.ProjectObjectRevision,
			eventstore.FieldTypeFieldName:      project.ProjectStateSearchField,
			eventstore.FieldTypeResourceOwner:  resourceOwner,
		},
	)
}

func (c *Commands) checkProjectExists(ctx context.Context, projectID, resourceOwner string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if !authz.GetFeatures(ctx).ShouldUseImprovedPerformance(feature.ImprovedPerformanceTypeProject) {
		return c.checkProjectExistsOld(ctx, projectID, resourceOwner)
	}

	_, state, err := c.projectAggregateByID(ctx, projectID, resourceOwner)
	if err != nil || !state.Valid() {
		return zerrors.ThrowPreconditionFailed(err, "COMMA-VCnwD", "Errors.Project.NotFound")
	}
	return nil
}

func (c *Commands) ChangeProject(ctx context.Context, projectChange *domain.Project, resourceOwner string) (*domain.Project, error) {
	if !projectChange.IsValid() || projectChange.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4m9vS", "Errors.Project.Invalid")
	}

	if !authz.GetFeatures(ctx).ShouldUseImprovedPerformance(feature.ImprovedPerformanceTypeProject) {
		return c.changeProjectOld(ctx, projectChange, resourceOwner)
	}

	existingProject, err := c.getProjectWriteModelByID(ctx, projectChange.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingProject.State == domain.ProjectStateUnspecified || existingProject.State == domain.ProjectStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.NotFound")
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingProject.WriteModel)
	changedEvent, hasChanged, err := existingProject.NewChangedEvent(
		ctx,
		projectAgg,
		projectChange.Name,
		projectChange.ProjectRoleAssertion,
		projectChange.ProjectRoleCheck,
		projectChange.HasProjectCheck,
		projectChange.PrivateLabelingSetting)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-2M0fs", "Errors.NoChangesFound")
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingProject, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return projectWriteModelToProject(existingProject), nil
}

func (c *Commands) DeactivateProject(ctx context.Context, projectID string, resourceOwner string) (*domain.ObjectDetails, error) {
	if projectID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-88iF0", "Errors.Project.ProjectIDMissing")
	}

	if !authz.GetFeatures(ctx).ShouldUseImprovedPerformance(feature.ImprovedPerformanceTypeProject) {
		return c.deactivateProjectOld(ctx, projectID, resourceOwner)
	}

	projectAgg, state, err := c.projectAggregateByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if state == domain.ProjectStateUnspecified || state == domain.ProjectStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-112M9", "Errors.Project.NotFound")
	}
	if state != domain.ProjectStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-mki55", "Errors.Project.NotActive")
	}

	pushedEvents, err := c.eventstore.Push(ctx, project.NewProjectDeactivatedEvent(ctx, projectAgg))
	if err != nil {
		return nil, err
	}

	existingProject, err := c.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingProject.State == domain.ProjectStateUnspecified || existingProject.State == domain.ProjectStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-112M9", "Errors.Project.NotFound")
	}
	if existingProject.State != domain.ProjectStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-mki55", "Errors.Project.NotActive")
	}

	return &domain.ObjectDetails{
		ResourceOwner: pushedEvents[0].Aggregate().ResourceOwner,
		Sequence:      pushedEvents[0].Sequence(),
		EventDate:     pushedEvents[0].CreatedAt(),
	}, nil
}

func (c *Commands) ReactivateProject(ctx context.Context, projectID string, resourceOwner string) (*domain.ObjectDetails, error) {
	if projectID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-3ihsF", "Errors.Project.ProjectIDMissing")
	}

	if !authz.GetFeatures(ctx).ShouldUseImprovedPerformance(feature.ImprovedPerformanceTypeProject) {
		return c.reactivateProjectOld(ctx, projectID, resourceOwner)
	}

	projectAgg, state, err := c.projectAggregateByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if state == domain.ProjectStateUnspecified || state == domain.ProjectStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.NotFound")
	}

	if state != domain.ProjectStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-5M9bs", "Errors.Project.NotInactive")
	}

	existingProject, err := c.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingProject.State == domain.ProjectStateUnspecified || existingProject.State == domain.ProjectStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.NotFound")
	}
	if existingProject.State != domain.ProjectStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-5M9bs", "Errors.Project.NotInactive")
	}

	pushedEvents, err := c.eventstore.Push(ctx, project.NewProjectReactivatedEvent(ctx, projectAgg))
	if err != nil {
		return nil, err
	}

	return &domain.ObjectDetails{
		ResourceOwner: pushedEvents[0].Aggregate().ResourceOwner,
		Sequence:      pushedEvents[0].Sequence(),
		EventDate:     pushedEvents[0].CreatedAt(),
	}, nil
}

func (c *Commands) RemoveProject(ctx context.Context, projectID, resourceOwner string, cascadingUserGrantIDs ...string) (*domain.ObjectDetails, error) {
	if projectID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-66hM9", "Errors.Project.ProjectIDMissing")
	}

	if !authz.GetFeatures(ctx).ShouldUseImprovedPerformance(feature.ImprovedPerformanceTypeProject) {
		return c.removeProjectOld(ctx, projectID, resourceOwner)
	}

	existingProject, err := c.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingProject.State == domain.ProjectStateUnspecified || existingProject.State == domain.ProjectStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.NotFound")
	}

	samlEntityIDsAgg, err := c.getSAMLEntityIdsWriteModelByProjectID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	uniqueConstraints := make([]*eventstore.UniqueConstraint, len(samlEntityIDsAgg.EntityIDs))
	for i, entityID := range samlEntityIDsAgg.EntityIDs {
		uniqueConstraints[i] = project.NewRemoveSAMLConfigEntityIDUniqueConstraint(entityID.EntityID)
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingProject.WriteModel)
	events := []eventstore.Command{
		project.NewProjectRemovedEvent(ctx, projectAgg, existingProject.Name, uniqueConstraints),
	}

	for _, grantID := range cascadingUserGrantIDs {
		event, _, err := c.removeUserGrant(ctx, grantID, "", true)
		if err != nil {
			logging.LogWithFields("COMMAND-b8Djf", "usergrantid", grantID).WithError(err).Warn("could not cascade remove user grant")
			continue
		}
		events = append(events, event)
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingProject, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingProject.WriteModel), nil
}

func (c *Commands) getProjectWriteModelByID(ctx context.Context, projectID, resourceOwner string) (_ *ProjectWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	projectWriteModel := NewProjectWriteModel(projectID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, projectWriteModel)
	if err != nil {
		return nil, err
	}
	return projectWriteModel, nil
}

func (c *Commands) getSAMLEntityIdsWriteModelByProjectID(ctx context.Context, projectID, resourceOwner string) (*SAMLEntityIDsWriteModel, error) {
	samlEntityIDsAgg := NewSAMLEntityIDsWriteModel(projectID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, samlEntityIDsAgg)
	if err != nil {
		return nil, err
	}
	return samlEntityIDsAgg, nil
}
