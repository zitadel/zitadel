package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"reflect"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddIAMMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	addedMember := NewIAMMemberWriteModel(member.UserID)
	iamAgg := IAMAggregateFromWriteModel(&addedMember.MemberWriteModel.WriteModel)
	event, err := c.addIAMMember(ctx, iamAgg, addedMember, member)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedMember, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return memberWriteModelToMember(&addedMember.MemberWriteModel), nil
}

func (c *Commands) addIAMMember(ctx context.Context, iamAgg *eventstore.Aggregate, addedMember *IAMMemberWriteModel, member *domain.Member) (eventstore.EventPusher, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-GR34U", "Errors.IAM.MemberInvalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.IAMRolePrefix, c.zitadelRoles)) > 0 {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-4m0fS", "Errors.IAM.MemberInvalid")
	}

	err := c.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.State == domain.MemberStateActive {
		return nil, errors.ThrowAlreadyExists(nil, "IAM-sdgQ4", "Errors.IAM.Member.AlreadyExists")
	}

	return iam_repo.NewMemberAddedEvent(ctx, iamAgg, member.UserID, member.Roles...), nil
}

//ChangeIAMMember updates an existing member
func (c *Commands) ChangeIAMMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-LiaZi", "Errors.IAM.MemberInvalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.IAMRolePrefix, c.zitadelRoles)) > 0 {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-3m9fs", "Errors.IAM.MemberInvalid")
	}

	existingMember, err := c.iamMemberWriteModelByID(ctx, member.UserID)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-LiaZi", "Errors.IAM.Member.RolesNotChanged")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingMember.MemberWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewMemberChangedEvent(ctx, iamAgg, member.UserID, member.Roles...))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&existingMember.MemberWriteModel), nil
}

func (c *Commands) RemoveIAMMember(ctx context.Context, userID string) error {
	m, err := c.iamMemberWriteModelByID(ctx, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if errors.IsNotFound(err) {
		return nil
	}

	iamAgg := IAMAggregateFromWriteModel(&m.MemberWriteModel.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, iam_repo.NewMemberRemovedEvent(ctx, iamAgg, userID))
	return err
}

func (c *Commands) iamMemberWriteModelByID(ctx context.Context, userID string) (member *IAMMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMMemberWriteModel(userID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, errors.ThrowNotFound(nil, "IAM-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
