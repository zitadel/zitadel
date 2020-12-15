package command

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) AddMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-W8m4l", "Errors.IAM.MemberInvalid")
	}

	addedMember := iam_repo.NewMemberWriteModel(member.AggregateID, member.UserID)
	err := r.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.IsActive {
		return nil, errors.ThrowAlreadyExists(nil, "IAM-PtXi1", "Errors.IAM.Member.AlreadyExists")
	}

	iamAgg := iam_repo.AggregateFromWriteModel(&addedMember.WriteModel.WriteModel).
		PushMemberAdded(ctx, member.UserID, member.Roles...)

	err = r.eventstore.PushAggregate(ctx, addedMember, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToMember(addedMember), nil
}

//ChangeMember updates an existing member
func (r *CommandSide) ChangeMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-LiaZi", "Errors.IAM.MemberInvalid")
	}

	existingMember, err := r.memberWriteModelByID(ctx, member.AggregateID, member.UserID)
	if err != nil {
		return nil, err
	}

	iam := iam_repo.AggregateFromWriteModel(&existingMember.WriteModel.WriteModel).
		PushMemberChangedFromExisting(ctx, existingMember, member.Roles...)

	events, err := r.eventstore.PushAggregates(ctx, iam)
	if err != nil {
		return nil, err
	}

	existingMember.AppendEvents(events...)
	if err = existingMember.Reduce(); err != nil {
		return nil, err
	}

	return writeModelToMember(existingMember), nil
}

func (r *CommandSide) RemoveMember(ctx context.Context, member *iam_model.IAMMember) error {
	m, err := r.memberWriteModelByID(ctx, member.AggregateID, member.UserID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if errors.IsNotFound(err) {
		return nil
	}

	iamAgg := iam_repo.AggregateFromWriteModel(&m.WriteModel.WriteModel).
		PushEvents(iam_repo.NewMemberRemovedEvent(ctx, member.UserID))

	return r.eventstore.PushAggregate(ctx, m, iamAgg)
}

func (r *CommandSide) memberWriteModelByID(ctx context.Context, iamID, userID string) (member *iam_repo.MemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := iam_repo.NewMemberWriteModel(iamID, userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if !writeModel.IsActive {
		return nil, errors.ThrowNotFound(nil, "IAM-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
