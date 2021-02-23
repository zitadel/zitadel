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

func (r *CommandSide) AddProjectGrantMember(ctx context.Context, member *domain.ProjectGrantMember, resourceOwner string) (*domain.ProjectGrantMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-8fi7G", "Errors.Project.Member.Invalid")
	}
	err := r.checkUserExists(ctx, member.UserID, "")
	if err != nil {
		return nil, err
	}
	addedMember := NewProjectGrantMemberWriteModel(member.AggregateID, member.UserID, member.GrantID)
	err = r.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.State == domain.MemberStateActive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-16dVN", "Errors.Project.Member.AlreadyExists")
	}
	projectAgg := ProjectAggregateFromWriteModel(&addedMember.WriteModel)
	pushedEvents, err := r.eventstore.PushEvents(
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
func (r *CommandSide) ChangeProjectGrantMember(ctx context.Context, member *domain.ProjectGrantMember, resourceOwner string) (*domain.ProjectGrantMember, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-109fs", "Errors.Project.Member.Invalid")
	}

	existingMember, err := r.projectGrantMemberWriteModelByID(ctx, member.AggregateID, member.UserID, member.GrantID)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-2n8vx", "Errors.Project.Member.RolesNotChanged")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingMember.WriteModel)
	pushedEvents, err := r.eventstore.PushEvents(
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

func (r *CommandSide) RemoveProjectGrantMember(ctx context.Context, projectID, userID, grantID, resourceOwner string) error {
	m, err := r.projectGrantMemberWriteModelByID(ctx, projectID, userID, grantID)
	if err != nil {
		return err
	}

	projectAgg := ProjectAggregateFromWriteModel(&m.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, project.NewProjectGrantMemberRemovedEvent(ctx, projectAgg, projectID, userID, grantID))
	return err
}

func (r *CommandSide) projectGrantMemberWriteModelByID(ctx context.Context, projectID, userID, grantID string) (member *ProjectGrantMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewProjectGrantMemberWriteModel(projectID, userID, grantID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, errors.ThrowNotFound(nil, "PROJECT-37fug", "Errors.NotFound")
	}

	return writeModel, nil
}
