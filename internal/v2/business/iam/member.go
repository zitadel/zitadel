package iam

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *Repository) AddIAMMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error) {
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
		// PushEvents(iam_repo.NewMemberAddedEvent(ctx, member.UserID, member.Roles...))

	events, err := r.eventstore.PushAggregates(ctx, iamAgg)
	if err != nil {
		return nil, err
	}

	if err = iam.AppendAndReduce(events...); err != nil {
		return nil, err
	}

	_, addedMember := iam.Members.MemberByUserID(member.UserID)
	if member == nil {
		return nil, errors.ThrowInternal(nil, "IAM-nuoDN", "member not saved")
	}
	return readModelToMember(addedMember), nil
}

func (r *Repository) ChangeIAMMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-LiaZi", "Errors.IAM.MemberInvalid")
	}

	iam, err := r.iamByID(ctx, member.AggregateID)
	if err != nil {
		return nil, err
	}

	existingMember, err := r.memberWriteModelByID(ctx, member.AggregateID, member.UserID)
	if err != nil {
		return nil, err
	}

	iamAgg := iam_repo.AggregateFromReadModel(iam).
		PushMemberChanged(ctx, existingMember, nil)

	events, err := r.eventstore.PushAggregates(ctx, iamAgg)
	if err != nil {
		return nil, err
	}

	if err = iam.AppendAndReduce(events...); err != nil {
		return nil, err
	}

	_, addedMember := iam.Members.MemberByUserID(member.UserID)
	if member == nil {
		return nil, errors.ThrowInternal(nil, "IAM-E5nTQ", "member not saved")
	}
	return readModelToMember(addedMember), nil
}

func (r *Repository) RemoveIAMMember(ctx context.Context, member *iam_model.IAMMember) error {
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
