package command

import (
	"context"
	"reflect"

	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

func (r *CommandSide) AddOrgMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	addedMember := NewOrgMemberWriteModel(member.AggregateID, member.UserID)
	orgAgg := OrgAggregateFromWriteModel(&addedMember.WriteModel)
	err := r.addOrgMember(ctx, orgAgg, addedMember, member)
	if err != nil {
		return nil, err
	}

	err = r.eventstore.PushAggregate(ctx, addedMember, orgAgg)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&addedMember.MemberWriteModel), nil
}

func (r *CommandSide) addOrgMember(ctx context.Context, orgAgg *org.Aggregate, addedMember *OrgMemberWriteModel, member *domain.Member) error {
	//TODO: check if roles valid

	if !member.IsValid() {
		return caos_errs.ThrowPreconditionFailed(nil, "Org-W8m4l", "Errors.Org.MemberInvalid")
	}

	err := r.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return err
	}
	if addedMember.IsActive {
		return errors.ThrowAlreadyExists(nil, "Org-PtXi1", "Errors.Org.Member.AlreadyExists")
	}

	orgAgg.PushEvents(org.NewMemberAddedEvent(ctx, member.UserID, member.Roles...))

	return nil
}

//ChangeOrgMember updates an existing member
func (r *CommandSide) ChangeOrgMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-LiaZi", "Errors.Org.MemberInvalid")
	}

	existingMember, err := r.orgMemberWriteModelByID(ctx, member.AggregateID, member.UserID)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-LiaZi", "Errors.Org.Member.RolesNotChanged")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingMember.MemberWriteModel.WriteModel)
	orgAgg.PushEvents(org.NewMemberChangedEvent(ctx, member.UserID, member.Roles...))

	events, err := r.eventstore.PushAggregates(ctx, orgAgg)
	if err != nil {
		return nil, err
	}

	existingMember.AppendEvents(events...)
	if err = existingMember.Reduce(); err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&existingMember.MemberWriteModel), nil
}

func (r *CommandSide) RemoveOrgMember(ctx context.Context, orgID, userID string) error {
	m, err := r.orgMemberWriteModelByID(ctx, orgID, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if errors.IsNotFound(err) {
		return nil
	}

	orgAgg := OrgAggregateFromWriteModel(&m.MemberWriteModel.WriteModel)
	orgAgg.PushEvents(org.NewMemberRemovedEvent(ctx, userID))

	return r.eventstore.PushAggregate(ctx, m, orgAgg)
}

func (r *CommandSide) orgMemberWriteModelByID(ctx context.Context, orgID, userID string) (member *OrgMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewOrgMemberWriteModel(orgID, userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if !writeModel.IsActive {
		return nil, errors.ThrowNotFound(nil, "Org-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
