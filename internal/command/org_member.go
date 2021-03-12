package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"reflect"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddOrgMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	addedMember := NewOrgMemberWriteModel(member.AggregateID, member.UserID)
	orgAgg := OrgAggregateFromWriteModel(&addedMember.WriteModel)
	event, err := c.addOrgMember(ctx, orgAgg, addedMember, member)
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

func (c *Commands) addOrgMember(ctx context.Context, orgAgg *eventstore.Aggregate, addedMember *OrgMemberWriteModel, member *domain.Member) (eventstore.EventPusher, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-W8m4l", "Errors.Org.MemberInvalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.OrgRolePrefix, c.zitadelRoles)) > 0 {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-3m9fs", "Errors.Org.MemberInvalid")
	}
	err := c.checkUserExists(ctx, addedMember.UserID, "")
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "ORG-9cmsd", "Errors.User.NotFound")
	}
	err = c.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.State == domain.MemberStateActive {
		return nil, errors.ThrowAlreadyExists(nil, "Org-PtXi1", "Errors.Org.Member.AlreadyExists")
	}

	return org.NewMemberAddedEvent(ctx, orgAgg, member.UserID, member.Roles...), nil
}

//ChangeOrgMember updates an existing member
func (c *Commands) ChangeOrgMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-LiaZi", "Errors.Org.MemberInvalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.OrgRolePrefix, c.zitadelRoles)) > 0 {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-m9fG8", "Errors.Org.MemberInvalid")
	}

	existingMember, err := c.orgMemberWriteModelByID(ctx, member.AggregateID, member.UserID)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-LiaZi", "Errors.Org.Member.RolesNotChanged")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingMember.MemberWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewMemberChangedEvent(ctx, orgAgg, member.UserID, member.Roles...))
	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&existingMember.MemberWriteModel), nil
}

func (c *Commands) RemoveOrgMember(ctx context.Context, orgID, userID string) (*domain.ObjectDetails, error) {
	m, err := c.orgMemberWriteModelByID(ctx, orgID, userID)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if errors.IsNotFound(err) {
		return nil, nil
	}

	orgAgg := OrgAggregateFromWriteModel(&m.MemberWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewMemberRemovedEvent(ctx, orgAgg, userID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(m, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&m.WriteModel), nil
}

func (c *Commands) orgMemberWriteModelByID(ctx context.Context, orgID, userID string) (member *OrgMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewOrgMemberWriteModel(orgID, userID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, errors.ThrowNotFound(nil, "Org-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
