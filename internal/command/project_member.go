package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"reflect"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddProjectMember(ctx context.Context, member *domain.Member, resourceOwner string) (*domain.Member, error) {
	addedMember := NewProjectMemberWriteModel(member.AggregateID, member.UserID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedMember.WriteModel)
	event, err := c.addProjectMember(ctx, projectAgg, addedMember, member)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&addedMember.MemberWriteModel), nil
}

func (c *Commands) addProjectMember(ctx context.Context, projectAgg *eventstore.Aggregate, addedMember *ProjectMemberWriteModel, member *domain.Member) (eventstore.EventPusher, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "PROJECT-W8m4l", "Errors.Project.Member.Invalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.ProjectRolePrefix, c.zitadelRoles)) > 0 {
		return nil, caos_errs.ThrowInvalidArgument(nil, "PROJECT-3m9ds", "Errors.Project.Member.Invalid")
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
		return nil, errors.ThrowAlreadyExists(nil, "PROJECT-PtXi1", "Errors.Project.Member.AlreadyExists")
	}

	return project.NewProjectMemberAddedEvent(ctx, projectAgg, member.UserID, member.Roles...), nil
}

//ChangeProjectMember updates an existing member
func (c *Commands) ChangeProjectMember(ctx context.Context, member *domain.Member, resourceOwner string) (*domain.Member, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "PROJECT-LiaZi", "Errors.Project.Member.Invalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.ProjectRolePrefix, c.zitadelRoles)) > 0 {
		return nil, caos_errs.ThrowInvalidArgument(nil, "PROJECT-3m9d", "Errors.Project.Member.Invalid")
	}

	existingMember, err := c.projectMemberWriteModelByID(ctx, member.AggregateID, member.UserID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-LiaZi", "Errors.Project.Member.RolesNotChanged")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingMember.MemberWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, project.NewProjectMemberChangedEvent(ctx, projectAgg, member.UserID, member.Roles...))
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
		return nil, caos_errs.ThrowInvalidArgument(nil, "PROJECT-66mHd", "Errors.Project.Member.Invalid")
	}
	m, err := c.projectMemberWriteModelByID(ctx, projectID, userID, resourceOwner)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if errors.IsNotFound(err) {
		return nil, nil
	}

	projectAgg := ProjectAggregateFromWriteModel(&m.MemberWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, project.NewProjectMemberRemovedEvent(ctx, projectAgg, userID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(m, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&m.WriteModel), nil
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
		return nil, errors.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
