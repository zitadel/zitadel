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

func (r *CommandSide) AddProjectMember(ctx context.Context, member *domain.Member, resourceOwner string) (*domain.Member, error) {
	addedMember := NewProjectMemberWriteModel(member.AggregateID, member.UserID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedMember.WriteModel)
	err := r.addProjectMember(ctx, projectAgg, addedMember, member)
	if err != nil {
		return nil, err
	}

	err = r.eventstore.PushAggregate(ctx, addedMember, projectAgg)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&addedMember.MemberWriteModel), nil
}

func (r *CommandSide) addProjectMember(ctx context.Context, projectAgg *project.Aggregate, addedMember *ProjectMemberWriteModel, member *domain.Member) error {
	//TODO: check if roles valid

	if !member.IsValid() {
		return caos_errs.ThrowPreconditionFailed(nil, "PROJECT-W8m4l", "Errors.Project.MemberInvalid")
	}

	err := r.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return err
	}
	if addedMember.State == domain.MemberStateActive {
		return errors.ThrowAlreadyExists(nil, "PROJECT-PtXi1", "Errors.Project.Member.AlreadyExists")
	}

	projectAgg.PushEvents(project.NewProjectMemberAddedEvent(ctx, member.UserID, member.Roles...))

	return nil
}

//ChangeProjectMember updates an existing member
func (r *CommandSide) ChangeProjectMember(ctx context.Context, member *domain.Member, resourceOwner string) (*domain.Member, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-LiaZi", "Errors.Project.MemberInvalid")
	}

	existingMember, err := r.projectMemberWriteModelByID(ctx, member.AggregateID, member.UserID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-LiaZi", "Errors.Project.Member.RolesNotChanged")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingMember.MemberWriteModel.WriteModel)
	projectAgg.PushEvents(project.NewProjectMemberChangedEvent(ctx, member.UserID, member.Roles...))

	events, err := r.eventstore.PushAggregates(ctx, projectAgg)
	if err != nil {
		return nil, err
	}

	existingMember.AppendEvents(events...)
	if err = existingMember.Reduce(); err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&existingMember.MemberWriteModel), nil
}

func (r *CommandSide) RemoveProjectMember(ctx context.Context, projectID, userID, resourceOwner string) error {
	m, err := r.projectMemberWriteModelByID(ctx, projectID, userID, resourceOwner)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if errors.IsNotFound(err) {
		return nil
	}

	projectAgg := ProjectAggregateFromWriteModel(&m.MemberWriteModel.WriteModel)
	projectAgg.PushEvents(project.NewProjectMemberRemovedEvent(ctx, userID))

	return r.eventstore.PushAggregate(ctx, m, projectAgg)
}

func (r *CommandSide) projectMemberWriteModelByID(ctx context.Context, projectID, userID, resourceOwner string) (member *ProjectMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewProjectMemberWriteModel(projectID, userID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, errors.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
