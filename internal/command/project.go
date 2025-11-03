package command

import (
	"context"
	"strings"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddProject struct {
	models.ObjectRoot

	Name                   string
	ProjectRoleAssertion   bool
	ProjectRoleCheck       bool
	HasProjectCheck        bool
	PrivateLabelingSetting domain.PrivateLabelingSetting
}

func (p *AddProject) IsValid() error {
	if p.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-fmq7bqQX1s", "Errors.ResourceOwnerMissing")
	}
	if p.Name == "" {
		return zerrors.ThrowInvalidArgument(nil, "PROJECT-IOVCC", "Errors.Project.Invalid")
	}
	return nil
}

func (c *Commands) AddProject(ctx context.Context, add *AddProject) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err := add.IsValid(); err != nil {
		return nil, err
	}

	if add.AggregateID == "" {
		add.AggregateID, err = c.idGenerator.Next()
		if err != nil {
			return nil, err
		}
	}
	wm, err := c.getProjectWriteModelByID(ctx, add.AggregateID, add.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if isProjectStateExists(wm.State) {
		return nil, zerrors.ThrowAlreadyExists(nil, "COMMAND-opamwu", "Errors.Project.AlreadyExisting")
	}
	if err := c.checkPermissionCreateProject(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, err
	}

	events := []eventstore.Command{
		project.NewProjectAddedEvent(
			ctx,
			ProjectAggregateFromWriteModelWithCTX(ctx, &wm.WriteModel),
			add.Name,
			add.ProjectRoleAssertion,
			add.ProjectRoleCheck,
			add.HasProjectCheck,
			add.PrivateLabelingSetting),
	}
	postCommit, err := c.projectCreatedMilestone(ctx, &events)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	postCommit(ctx)
	err = AppendAndReduce(wm, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
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

func (c *Commands) checkProjectExists(ctx context.Context, projectID, resourceOwner string) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if !authz.GetFeatures(ctx).ShouldUseImprovedPerformance(feature.ImprovedPerformanceTypeProject) {
		return c.checkProjectExistsOld(ctx, projectID, resourceOwner)
	}

	agg, state, err := c.projectAggregateByID(ctx, projectID, resourceOwner)
	if err != nil || !state.Valid() {
		return "", zerrors.ThrowPreconditionFailed(err, "COMMA-VCnwD", "Errors.Project.NotFound")
	}
	return agg.ResourceOwner, nil
}

type ChangeProject struct {
	models.ObjectRoot

	Name                   *string
	ProjectRoleAssertion   *bool
	ProjectRoleCheck       *bool
	HasProjectCheck        *bool
	PrivateLabelingSetting *domain.PrivateLabelingSetting
}

func (p *ChangeProject) IsValid() error {
	if p.AggregateID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-4m9vS", "Errors.Project.Invalid")
	}
	if p.Name != nil && *p.Name == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-4m9vS", "Errors.Project.Invalid")
	}
	return nil
}

func (c *Commands) ChangeProject(ctx context.Context, change *ChangeProject) (_ *domain.ObjectDetails, err error) {
	if err := change.IsValid(); err != nil {
		return nil, err
	}

	existing, err := c.getProjectWriteModelByID(ctx, change.AggregateID, change.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if !isProjectStateExists(existing.State) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.NotFound")
	}
	if err := c.checkPermissionUpdateProject(ctx, existing.ResourceOwner, existing.AggregateID); err != nil {
		return nil, err
	}

	changedEvent := existing.NewChangedEvent(
		ctx,
		ProjectAggregateFromWriteModelWithCTX(ctx, &existing.WriteModel),
		change.Name,
		change.ProjectRoleAssertion,
		change.ProjectRoleCheck,
		change.HasProjectCheck,
		change.PrivateLabelingSetting)
	if changedEvent == nil {
		return writeModelToObjectDetails(&existing.WriteModel), nil
	}
	err = c.pushAppendAndReduce(ctx, existing, changedEvent)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existing.WriteModel), nil
}

func (c *Commands) DeactivateProject(ctx context.Context, projectID string, resourceOwner string) (*domain.ObjectDetails, error) {
	if projectID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-88iF0", "Errors.Project.ProjectIDMissing")
	}

	if !authz.GetFeatures(ctx).ShouldUseImprovedPerformance(feature.ImprovedPerformanceTypeProject) {
		return c.deactivateProjectOld(ctx, projectID, resourceOwner)
	}

	projectAgg, state, err := c.projectAggregateByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if !isProjectStateExists(state) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-112M9", "Errors.Project.NotFound")
	}
	if state != domain.ProjectStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-mki55", "Errors.Project.NotActive")
	}
	if err := c.checkPermissionUpdateProject(ctx, projectAgg.ResourceOwner, projectAgg.ID); err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, project.NewProjectDeactivatedEvent(ctx, projectAgg))
	if err != nil {
		return nil, err
	}

	return &domain.ObjectDetails{
		ResourceOwner: pushedEvents[0].Aggregate().ResourceOwner,
		Sequence:      pushedEvents[0].Sequence(),
		EventDate:     pushedEvents[0].CreatedAt(),
	}, nil
}

func (c *Commands) ReactivateProject(ctx context.Context, projectID string, resourceOwner string) (*domain.ObjectDetails, error) {
	if projectID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-3ihsF", "Errors.Project.ProjectIDMissing")
	}

	if !authz.GetFeatures(ctx).ShouldUseImprovedPerformance(feature.ImprovedPerformanceTypeProject) {
		return c.reactivateProjectOld(ctx, projectID, resourceOwner)
	}

	projectAgg, state, err := c.projectAggregateByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if !isProjectStateExists(state) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.NotFound")
	}
	if state != domain.ProjectStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-5M9bs", "Errors.Project.NotInactive")
	}
	if err := c.checkPermissionUpdateProject(ctx, projectAgg.ResourceOwner, projectAgg.ID); err != nil {
		return nil, err
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

// Deprecated: use commands.DeleteProject
func (c *Commands) RemoveProject(ctx context.Context, projectID, resourceOwner string, cascadingUserGrantIDs ...string) (*domain.ObjectDetails, error) {
	if projectID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-66hM9", "Errors.Project.ProjectIDMissing")
	}

	existingProject, err := c.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isProjectStateExists(existingProject.State) {
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

	events := []eventstore.Command{
		project.NewProjectRemovedEvent(ctx,
			ProjectAggregateFromWriteModelWithCTX(ctx, &existingProject.WriteModel),
			existingProject.Name,
			uniqueConstraints,
		),
	}

	for _, grantID := range cascadingUserGrantIDs {
		event, _, err := c.removeUserGrant(ctx, grantID, "", true, false, nil)
		if err != nil {
			logging.WithFields("id", "COMMAND-b8Djf", "usergrantid", grantID).WithError(err).Warn("could not cascade remove user grant")
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

func (c *Commands) DeleteProject(ctx context.Context, id, resourceOwner string, cascadingUserGrantIDs ...string) (time.Time, error) {
	if id == "" {
		return time.Time{}, zerrors.ThrowInvalidArgument(nil, "COMMAND-obqos2l3no", "Errors.IDMissing")
	}

	existing, err := c.getProjectWriteModelByID(ctx, id, resourceOwner)
	if err != nil {
		return time.Time{}, err
	}
	if !isProjectStateExists(existing.State) {
		return existing.WriteModel.ChangeDate, nil
	}
	if err := c.checkPermissionDeleteProject(ctx, existing.ResourceOwner, existing.AggregateID); err != nil {
		return time.Time{}, err
	}

	samlEntityIDsAgg, err := c.getSAMLEntityIdsWriteModelByProjectID(ctx, id, resourceOwner)
	if err != nil {
		return time.Time{}, err
	}

	uniqueConstraints := make([]*eventstore.UniqueConstraint, len(samlEntityIDsAgg.EntityIDs))
	for i, entityID := range samlEntityIDsAgg.EntityIDs {
		uniqueConstraints[i] = project.NewRemoveSAMLConfigEntityIDUniqueConstraint(entityID.EntityID)
	}
	events := []eventstore.Command{
		project.NewProjectRemovedEvent(ctx,
			ProjectAggregateFromWriteModelWithCTX(ctx, &existing.WriteModel),
			existing.Name,
			uniqueConstraints,
		),
	}
	for _, grantID := range cascadingUserGrantIDs {
		event, _, err := c.removeUserGrant(ctx, grantID, "", true, false, nil)
		if err != nil {
			logging.WithFields("id", "COMMAND-b8Djf", "usergrantid", grantID).WithError(err).Warn("could not cascade remove user grant")
			continue
		}
		events = append(events, event)
	}

	if err := c.pushAppendAndReduce(ctx, existing, events...); err != nil {
		return time.Time{}, err
	}
	return existing.WriteModel.ChangeDate, nil
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
