package iam

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *Repository) AddMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-W8m4l", "Errors.IAM.MemberInvalid")
	}

	iam, err := r.iamByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}

	idx, _ := iam.Members.MemberByUserID(member.UserID)
	if idx > -1 {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-GPhuz", "Errors.IAM.MemberAlreadyExisting")
	}

	iamAgg := iam_repo.AggregateFromReadModel(iam).
		PushMemberAdded(ctx, member.UserID, member.Roles...)

	events, err := r.eventstore.PushAggregates(ctx, iamAgg)
	if err != nil {
		return nil, err
	}

	if err = iam.AppendAndReduce(events...); err != nil {
		return nil, err
	}

	_, addedMember := iam.Members.MemberByUserID(member.UserID)
	if member == nil {
		return nil, errors.ThrowInternal(nil, "IAM-nuoDN", "Errors.Internal")
	}
	return readModelToMember(addedMember), nil
}

//ChangeMember updates an existing member
func (r *Repository) ChangeMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error) {
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

func (r *Repository) RemoveMember(ctx context.Context, member *iam_model.IAMMember) error {
	iam, err := r.iamByID(ctx, member.AggregateID)
	if err != nil {
		return err
	}

	i, _ := iam.Members.MemberByUserID(member.UserID)
	if i == -1 {
		return nil
	}

	iamAgg := iam_repo.AggregateFromReadModel(iam).
		PushEvents(iam_repo.NewMemberRemovedEvent(ctx, member.UserID))

	events, err := r.eventstore.PushAggregates(ctx, iamAgg)
	if err != nil {
		return err
	}

	return iam.AppendAndReduce(events...)
}

func (r *Repository) MemberByID(ctx context.Context, iamID, userID string) (member *iam_repo.MemberReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	member = new(iam_repo.MemberReadModel)

	//query view

	query := eventstore.NewSearchQueryFactory(eventstore.ColumnsEvent, iam_repo.AggregateType).
		AggregateIDs(iamID).
		EventData(map[string]interface{}{
			"userId": userID,
		}).SequenceGreater(member.ProcessedSequence)

	err = r.eventstore.FilterToReducer(ctx, query, member)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (r *Repository) memberWriteModelByID(ctx context.Context, iamID, userID string) (member *iam_repo.MemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := iam_repo.NewMemberReadModel(iamID, userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.IsRemoved {
		return nil, errors.ThrowNotFound(nil, "IAM-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
