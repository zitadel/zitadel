package command

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddProjectRole struct {
	models.ObjectRoot

	Key         string
	DisplayName string
	Group       string
}

func (p *AddProjectRole) IsValid() bool {
	return p.AggregateID != "" && p.Key != ""
}

func (c *Commands) AddProjectRole(ctx context.Context, projectRole *AddProjectRole) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	projectResourceOwner, err := c.checkProjectExists(ctx, projectRole.AggregateID, projectRole.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if projectRole.ResourceOwner == "" {
		projectRole.ResourceOwner = projectResourceOwner
	}
	if err := c.checkPermissionWriteProjectRole(ctx, projectRole.ResourceOwner, projectRole.AggregateID); err != nil {
		return nil, err
	}

	roleWriteModel := NewProjectRoleWriteModelWithKey(projectRole.Key, projectRole.AggregateID, projectRole.ResourceOwner)
	if roleWriteModel.ResourceOwner != projectResourceOwner {
		return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-RLB4UpqQSd", "Errors.Project.Role.Invalid")
	}

	projectAgg := ProjectAggregateFromWriteModelWithCTX(ctx, &roleWriteModel.WriteModel)
	events, err := c.addProjectRoles(ctx, projectAgg, projectRole)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(roleWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&roleWriteModel.WriteModel), nil
}

func (c *Commands) checkPermissionWriteProjectRole(ctx context.Context, orgID, projectID string) error {
	return c.checkPermission(ctx, domain.PermissionProjectRoleWrite, orgID, projectID)
}

func (c *Commands) BulkAddProjectRole(ctx context.Context, projectID, resourceOwner string, projectRoles []*AddProjectRole) (details *domain.ObjectDetails, err error) {
	projectResourceOwner, err := c.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}
	for _, projectRole := range projectRoles {
		if projectRole.ResourceOwner == "" {
			projectRole.ResourceOwner = projectResourceOwner
		}
		if err := c.checkPermissionWriteProjectRole(ctx, projectRole.ResourceOwner, projectID); err != nil {
			return nil, err
		}
		if projectRole.ResourceOwner != projectResourceOwner {
			return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-9ZXtdaJKJJ", "Errors.Project.Role.Invalid")
		}
	}

	roleWriteModel := NewProjectRoleWriteModel(projectID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModelWithCTX(ctx, &roleWriteModel.WriteModel)
	events, err := c.addProjectRoles(ctx, projectAgg, projectRoles...)
	if err != nil {
		return details, err
	}
	return c.pushAppendAndReduceDetails(ctx, roleWriteModel, events...)
}

func (c *Commands) addProjectRoles(ctx context.Context, projectAgg *eventstore.Aggregate, projectRoles ...*AddProjectRole) ([]eventstore.Command, error) {
	var events []eventstore.Command
	for _, projectRole := range projectRoles {
		if projectRole.ResourceOwner != projectAgg.ResourceOwner {
			return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-4Q2WjlbHvc", "Errors.Project.Role.Invalid")
		}
		if !projectRole.IsValid() {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4m9vS", "Errors.Project.Role.Invalid")
		}
		events = append(events, project.NewRoleAddedEvent(
			ctx,
			projectAgg,
			projectRole.Key,
			projectRole.DisplayName,
			projectRole.Group,
		))
	}

	return events, nil
}

type ChangeProjectRole struct {
	models.ObjectRoot

	Key         string
	DisplayName string
	Group       string
}

func (p *ChangeProjectRole) IsValid() bool {
	return p.AggregateID != "" && p.Key != ""
}

func (c *Commands) ChangeProjectRole(ctx context.Context, projectRole *ChangeProjectRole) (_ *domain.ObjectDetails, err error) {
	if !projectRole.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-2ilfW", "Errors.Project.Invalid")
	}
	projectResourceOwner, err := c.checkProjectExists(ctx, projectRole.AggregateID, projectRole.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if projectRole.ResourceOwner == "" {
		projectRole.ResourceOwner = projectResourceOwner
	}
	if err := c.checkPermissionWriteProjectRole(ctx, projectRole.ResourceOwner, projectRole.AggregateID); err != nil {
		return nil, err
	}

	existingRole, err := c.getProjectRoleWriteModelByID(ctx, projectRole.Key, projectRole.AggregateID, projectRole.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingRole.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-vv8M9", "Errors.Project.Role.NotExisting")
	}
	if existingRole.ResourceOwner != projectResourceOwner {
		return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-3MizLWveMf", "Errors.Project.Role.Invalid")
	}

	projectAgg := ProjectAggregateFromWriteModelWithCTX(ctx, &existingRole.WriteModel)

	changeEvent, changed, err := existingRole.NewProjectRoleChangedEvent(ctx, projectAgg, projectRole.Key, projectRole.DisplayName, projectRole.Group)
	if err != nil {
		return nil, err
	}
	if !changed {
		return writeModelToObjectDetails(&existingRole.WriteModel), nil
	}

	return c.pushAppendAndReduceDetails(ctx, existingRole, changeEvent)
}

func (c *Commands) RemoveProjectRole(ctx context.Context, projectID, key, resourceOwner string, cascadingProjectGrantIds []string, cascadeUserGrantIDs ...string) (details *domain.ObjectDetails, err error) {
	if projectID == "" || key == "" {
		return details, zerrors.ThrowInvalidArgument(nil, "COMMAND-fl9eF", "Errors.Project.Role.Invalid")
	}
	existingRole, err := c.getProjectRoleWriteModelByID(ctx, key, projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	// return if project role is not existing
	if !existingRole.State.Exists() {
		return writeModelToObjectDetails(&existingRole.WriteModel), nil
	}
	if err := c.checkPermissionDeleteProjectRole(ctx, existingRole.ResourceOwner, projectID); err != nil {
		return nil, err
	}
	projectAgg := ProjectAggregateFromWriteModelWithCTX(ctx, &existingRole.WriteModel)
	events := []eventstore.Command{
		project.NewRoleRemovedEvent(ctx, projectAgg, key),
	}

	for _, projectGrantID := range cascadingProjectGrantIds {
		event, _, err := c.removeRoleFromProjectGrant(ctx, projectAgg, projectID, projectGrantID, key, true)
		if err != nil {
			logging.LogWithFields("COMMAND-6n77g", "projectgrantid", projectGrantID).WithError(err).Warn("could not cascade remove role from project grant")
			continue
		}
		events = append(events, event)
	}

	for _, grantID := range cascadeUserGrantIDs {
		event, err := c.removeRoleFromUserGrant(ctx, grantID, []string{key}, true)
		if err != nil {
			logging.LogWithFields("COMMAND-mK0of", "usergrantid", grantID).WithError(err).Warn("could not cascade remove role on user grant")
			continue
		}
		events = append(events, event)
	}

	return c.pushAppendAndReduceDetails(ctx, existingRole, events...)
}

func (c *Commands) checkPermissionDeleteProjectRole(ctx context.Context, orgID, projectID string) error {
	return c.checkPermission(ctx, domain.PermissionProjectRoleDelete, orgID, projectID)
}

func (c *Commands) getProjectRoleWriteModelByID(ctx context.Context, key, projectID, resourceOwner string) (*ProjectRoleWriteModel, error) {
	projectRoleWriteModel := NewProjectRoleWriteModelWithKey(key, projectID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, projectRoleWriteModel)
	if err != nil {
		return nil, err
	}
	return projectRoleWriteModel, nil
}
