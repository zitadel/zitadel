package command

import (
	"context"
	"reflect"

	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) AddIAMMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	addedMember := NewIAMMemberWriteModel(member.UserID)
	iamAgg := IAMAggregateFromWriteModel(&addedMember.MemberWriteModel.WriteModel)
	err := r.addIAMMember(ctx, iamAgg, addedMember, member)
	if err != nil {
		return nil, err
	}

	err = r.eventstore.PushAggregate(ctx, addedMember, iamAgg)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&addedMember.MemberWriteModel), nil
}

func (r *CommandSide) addIAMMember(ctx context.Context, iamAgg *iam_repo.Aggregate, addedMember *IAMMemberWriteModel, member *domain.Member) error {
	//TODO: check if roles valid

	if !member.IsValid() {
		return caos_errs.ThrowPreconditionFailed(nil, "IAM-GR34U", "Errors.IAM.MemberInvalid")
	}

	err := r.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return err
	}
	if addedMember.State == domain.MemberStateActive {
		return errors.ThrowAlreadyExists(nil, "IAM-sdgQ4", "Errors.IAM.Member.AlreadyExists")
	}

	iamAgg.PushEvents(iam_repo.NewMemberAddedEvent(ctx, member.UserID, member.Roles...))

	return nil
}

//ChangeIAMMember updates an existing member
func (r *CommandSide) ChangeIAMMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-LiaZi", "Errors.IAM.MemberInvalid")
	}

	existingMember, err := r.iamMemberWriteModelByID(ctx, member.UserID)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-LiaZi", "Errors.IAM.Member.RolesNotChanged")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingMember.MemberWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewMemberChangedEvent(ctx, member.UserID, member.Roles...))

	events, err := r.eventstore.PushAggregates(ctx, iamAgg)
	if err != nil {
		return nil, err
	}

	existingMember.AppendEvents(events...)
	if err = existingMember.Reduce(); err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&existingMember.MemberWriteModel), nil
}

func (r *CommandSide) RemoveIAMMember(ctx context.Context, userID string) error {
	m, err := r.iamMemberWriteModelByID(ctx, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if errors.IsNotFound(err) {
		return nil
	}

	iamAgg := IAMAggregateFromWriteModel(&m.MemberWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewMemberRemovedEvent(ctx, userID))

	return r.eventstore.PushAggregate(ctx, m, iamAgg)
}

func (r *CommandSide) iamMemberWriteModelByID(ctx context.Context, userID string) (member *IAMMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMMemberWriteModel(userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, errors.ThrowNotFound(nil, "IAM-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
