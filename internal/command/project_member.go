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

func (c *Commands) AddProjectMember(ctx context.Context, member *domain.Member, resourceOwner string) (*domain.Member, error) {
	addedMember := NewProjectMemberWriteModel(member.AggregateID, member.UserID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedMember.WriteModel)
	event, err := c.addProjectMember(ctx, projectAgg, addedMember, member)
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

	return memberWriteModelToMember(&addedMember.MemberWriteModel), nil
}

func (c *Commands) addProjectMember(ctx context.Context, projectAgg *eventstore.Aggregate, addedMember *ProjectMemberWriteModel, member *domain.Member) (eventstore.Command, error) {
	if !member.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-W8m4l", "Errors.Project.Member.Invalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.ProjectRolePrefix, c.zitadelRoles)) > 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-3m9ds", "Errors.Project.Member.Invalid")
	}

	err := c.checkUserExists(ctx, addedMember.UserID, "")
	if err != nil {
		return nil, err
	}
	err = c.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.State == domain.MemberStateActive {
		return nil, zerrors.ThrowAlreadyExists(nil, "PROJECT-PtXi1", "Errors.Project.Member.AlreadyExists")
	}

	return project.NewProjectMemberAddedEvent(ctx, projectAgg, member.UserID, member.Roles...), nil
}

// ChangeProjectMember updates an existing member
func (c *Commands) ChangeProjectMember(ctx context.Context, member *domain.Member, resourceOwner string) (*domain.Member, error) {
	if !member.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-LiaZi", "Errors.Project.Member.Invalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.ProjectRolePrefix, c.zitadelRoles)) > 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-3m9d", "Errors.Project.Member.Invalid")
	}

	existingMember, err := c.projectMemberWriteModelByID(ctx, member.AggregateID, member.UserID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-LiaZi", "Errors.Project.Member.RolesNotChanged")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingMember.MemberWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, project.NewProjectMemberChangedEvent(ctx, projectAgg, member.UserID, member.Roles...))
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&existingMember.MemberWriteModel), nil
}

func (c *Commands) RemoveProjectMember(ctx context.Context, projectID, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if projectID == "" || userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-66mHd", "Errors.Project.Member.Invalid")
	}
	m, err := c.projectMemberWriteModelByID(ctx, projectID, userID, resourceOwner)
	if err != nil && !zerrors.IsNotFound(err) {
		return nil, err
	}
	if zerrors.IsNotFound(err) {
		// empty response because we have no data that match the request
		return &domain.ObjectDetails{}, nil
	}

	projectAgg := ProjectAggregateFromWriteModel(&m.MemberWriteModel.WriteModel)
	removeEvent := c.removeProjectMember(ctx, projectAgg, userID, false)
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

func (c *Commands) removeProjectMember(ctx context.Context, projectAgg *eventstore.Aggregate, userID string, cascade bool) eventstore.Command {
	if cascade {
		return project.NewProjectMemberCascadeRemovedEvent(
			ctx,
			projectAgg,
			userID)
	} else {
		return project.NewProjectMemberRemovedEvent(ctx, projectAgg, userID)
	}
}

func (c *Commands) projectMemberWriteModelByID(ctx context.Context, projectID, userID, resourceOwner string) (member *ProjectMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewProjectMemberWriteModel(projectID, userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
