package command

import (
	"context"
	"reflect"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddProjectGrantMember(ctx context.Context, member *domain.ProjectGrantMember, resourceOwner string) (*domain.ProjectGrantMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-8fi7G", "Errors.Project.Member.Invalid")
	}
	err := c.checkUserExists(ctx, member.UserID, "")
	if err != nil {
		return nil, err
	}
	addedMember := NewProjectGrantMemberWriteModel(member.AggregateID, member.UserID, member.GrantID)
	err = c.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.State == domain.MemberStateActive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-16dVN", "Errors.Project.Member.AlreadyExists")
	}
	projectAgg := ProjectAggregateFromWriteModel(&addedMember.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(
		ctx,
		project.NewProjectGrantMemberAddedEvent(ctx, projectAgg, member.AggregateID, member.UserID, member.GrantID, member.Roles...))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToProjectGrantMember(addedMember), nil
}

//ChangeProjectGrantMember updates an existing member
func (c *Commands) ChangeProjectGrantMember(ctx context.Context, member *domain.ProjectGrantMember, resourceOwner string) (*domain.ProjectGrantMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-109fs", "Errors.Project.Member.Invalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.ProjectGrantRolePrefix, c.zitadelRoles)) > 0 {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-m0sDf", "Errors.Project.Member.Invalid")
	}

	existingMember, err := c.projectGrantMemberWriteModelByID(ctx, member.AggregateID, member.UserID, member.GrantID)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-2n8vx", "Errors.Project.Member.RolesNotChanged")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingMember.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(
		ctx,
		project.NewProjectGrantMemberChangedEvent(ctx, projectAgg, member.UserID, member.GrantID, member.Roles...))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToProjectGrantMember(existingMember), nil
}

func (c *Commands) RemoveProjectGrantMember(ctx context.Context, projectID, userID, grantID, resourceOwner string) error {
	m, err := c.projectGrantMemberWriteModelByID(ctx, projectID, userID, grantID)
	if err != nil {
		return err
	}

	projectAgg := ProjectAggregateFromWriteModel(&m.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, project.NewProjectGrantMemberRemovedEvent(ctx, projectAgg, projectID, userID, grantID))
	return err
}

func (c *Commands) projectGrantMemberWriteModelByID(ctx context.Context, projectID, userID, grantID string) (member *ProjectGrantMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewProjectGrantMemberWriteModel(projectID, userID, grantID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, errors.ThrowNotFound(nil, "PROJECT-37fug", "Errors.NotFound")
	}

	return writeModel, nil
}
