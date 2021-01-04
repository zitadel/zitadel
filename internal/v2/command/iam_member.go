package command

import (
	"context"
	"reflect"

	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) AddIAMMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-W8m4l", "Errors.IAM.MemberInvalid")
	}

	addedMember := NewIAMMemberWriteModel(member.AggregateID, member.UserID)
	err := r.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.IsActive {
		return nil, errors.ThrowAlreadyExists(nil, "IAM-PtXi1", "Errors.IAM.Member.AlreadyExists")
	}

	iamAgg := IAMAggregateFromWriteModel(&addedMember.MemberWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewMemberAddedEvent(ctx, member.UserID, member.Roles...))

	err = r.eventstore.PushAggregate(ctx, addedMember, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToMember(addedMember), nil
}

//ChangeIAMMember updates an existing member
func (r *CommandSide) ChangeIAMMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-LiaZi", "Errors.IAM.MemberInvalid")
	}

	existingMember, err := r.iamMemberWriteModelByID(ctx, member.AggregateID, member.UserID)
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

	return writeModelToMember(existingMember), nil
}

func (r *CommandSide) RemoveIAMMember(ctx context.Context, member *iam_model.IAMMember) error {
	m, err := r.iamMemberWriteModelByID(ctx, member.AggregateID, member.UserID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if errors.IsNotFound(err) {
		return nil
	}

	iamAgg := IAMAggregateFromWriteModel(&m.MemberWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewMemberRemovedEvent(ctx, member.UserID))

	return r.eventstore.PushAggregate(ctx, m, iamAgg)
}

func (r *CommandSide) iamMemberWriteModelByID(ctx context.Context, iamID, userID string) (member *IAMMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMMemberWriteModel(iamID, userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if !writeModel.IsActive {
		return nil, errors.ThrowNotFound(nil, "IAM-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
