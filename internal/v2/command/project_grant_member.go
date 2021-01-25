package command

import (
	"context"
	"reflect"

	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

func (r *CommandSide) AddProjectGrantMember(ctx context.Context, member *domain.ProjectGrantMember, resourceOwner string) (*domain.ProjectGrantMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-8fi7G", "Errors.Project.Member.Invalid")
	}
	addedMember := NewProjectGrantMemberWriteModel(member.AggregateID, member.UserID, member.GrantID, resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.State == domain.MemberStateActive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-16dVN", "Errors.Project.Member.AlreadyExists")
	}
	projectAgg := ProjectAggregateFromWriteModel(&addedMember.WriteModel)
	projectAgg.PushEvents(project.NewProjectGrantMemberAddedEvent(ctx, member.AggregateID, member.UserID, member.GrantID, member.Roles...))

	err = r.eventstore.PushAggregate(ctx, addedMember, projectAgg)
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

	existingMember, err := r.projectGrantMemberWriteModelByID(ctx, member.AggregateID, member.UserID, member.GrantID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-2n8vx", "Errors.Project.Member.RolesNotChanged")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingMember.WriteModel)
	projectAgg.PushEvents(project.NewProjectGrantMemberChangedEvent(ctx, member.UserID, member.GrantID, member.Roles...))

	events, err := r.eventstore.PushAggregates(ctx, projectAgg)
	if err != nil {
		return nil, err
	}

	existingMember.AppendEvents(events...)
	if err = existingMember.Reduce(); err != nil {
		return nil, err
	}

	return memberWriteModelToProjectGrantMember(existingMember), nil
}

func (r *CommandSide) RemoveProjectGrantMember(ctx context.Context, projectID, userID, grantID, resourceOwner string) error {
	m, err := r.projectGrantMemberWriteModelByID(ctx, projectID, userID, grantID, resourceOwner)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if errors.IsNotFound(err) {
		return nil
	}

	projectAgg := ProjectAggregateFromWriteModel(&m.WriteModel)
	projectAgg.PushEvents(project.NewProjectGrantMemberRemovedEvent(ctx, projectID, userID, grantID))

	return r.eventstore.PushAggregate(ctx, m, projectAgg)
}

func (r *CommandSide) projectGrantMemberWriteModelByID(ctx context.Context, projectID, userID, grantID, resourceOwner string) (member *ProjectGrantMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewProjectGrantMemberWriteModel(projectID, userID, grantID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, errors.ThrowNotFound(nil, "PROJECT-37fug", "Errors.NotFound")
	}

	return writeModel, nil
}
