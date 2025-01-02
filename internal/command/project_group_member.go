package command

import (
	"context"
	"reflect"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddProjectGroupMember(ctx context.Context, member *domain.GroupMember, resourceOwner string) (_ *domain.GroupMember, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	addedMember := NewProjectGroupMemberWriteModel(member.AggregateID, member.GroupID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedMember.WriteModel)
	event, err := c.addProjectGroupMember(ctx, projectAgg, addedMember, member)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return projectGroupMemberWriteModelToMember(addedMember), nil
}

func (c *Commands) addProjectGroupMember(ctx context.Context, projectAgg *eventstore.Aggregate, addedMember *ProjectGroupMemberWriteModel, member *domain.GroupMember) (_ eventstore.Command, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if !member.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-W8m4l", "Errors.Project.Member.Invalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.ProjectRolePrefix, c.zitadelRoles)) > 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-3m9ds", "Errors.Project.Member.Invalid")
	}

	err = c.checkUserExists(ctx, addedMember.GroupID, "")
	if err != nil {
		return nil, err
	}
	err = c.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.State == domain.GroupMemberStateActive {
		return nil, zerrors.ThrowAlreadyExists(nil, "PROJECT-PtXi1", "Errors.Project.Member.AlreadyExists")
	}

	return project.NewProjectMemberAddedEvent(ctx, projectAgg, member.GroupID, member.Roles...), nil
}

// ChangeProjectMember updates an existing member
func (c *Commands) ChangeProjectGroupMember(ctx context.Context, member *domain.GroupMember, resourceOwner string) (*domain.GroupMember, error) {
	if !member.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-KJaZi", "Errors.Project.Group.Member.Invalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.ProjectRolePrefix, c.zitadelRoles)) > 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-4n9d", "Errors.Project.Group.Member.Invalid")
	}

	existingMember, err := c.projectGroupMemberWriteModelByID(ctx, member.AggregateID, member.GroupID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-KJaZi", "Errors.Project.Group.Member.RolesNotChanged")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingMember.ProjectGroupMemberWrite.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, project.NewProjectMemberChangedEvent(ctx, projectAgg, member.GroupID, member.Roles...))
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return projectGroupMemberWriteModelToMember(existingMember), nil
}

func (c *Commands) RemoveProjectGroupMember(ctx context.Context, projectID, groupID, resourceOwner string) (*domain.ObjectDetails, error) {
	if projectID == "" || groupID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-77bLd", "Errors.Project.Group.Member.Invalid")
	}
	m, err := c.projectGroupMemberWriteModelByID(ctx, projectID, groupID, resourceOwner)
	if err != nil && !zerrors.IsNotFound(err) {
		return nil, err
	}
	if zerrors.IsNotFound(err) {
		// empty response because we have no data that match the request
		return &domain.ObjectDetails{}, nil
	}

	projectAgg := ProjectAggregateFromWriteModel(&m.ProjectGroupMemberWrite.WriteModel)
	removeEvent := c.removeProjectGroupMember(ctx, projectAgg, groupID, false)
	pushedEvents, err := c.eventstore.Push(ctx, removeEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(m, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&m.WriteModel), nil
}

func (c *Commands) removeProjectGroupMember(ctx context.Context, projectAgg *eventstore.Aggregate, groupID string, cascade bool) eventstore.Command {
	if cascade {
		return project.NewProjectGroupMemberCascadeRemovedEvent(
			ctx,
			projectAgg,
			groupID)
	} else {
		return project.NewProjectGroupMemberRemovedEvent(ctx, projectAgg, groupID)
	}
}

func (c *Commands) projectGroupMemberWriteModelByID(ctx context.Context, projectID, groupID, resourceOwner string) (member *ProjectGroupMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewProjectGroupMemberWriteModel(projectID, groupID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.GroupMemberStateUnspecified || writeModel.State == domain.GroupMemberStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
