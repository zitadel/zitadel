package command

import (
	"context"
	"reflect"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddOrgMemberCommand(a *org.Aggregate, userID string, roles ...string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if userID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-4Mlfs", "Errors.Invalid.Argument")
		}
		if len(roles) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "V2-PfYhb", "Errors.Invalid.Argument")
		}

		if len(domain.CheckForInvalidRoles(roles, domain.OrgRolePrefix, c.zitadelRoles)) > 0 && len(domain.CheckForInvalidRoles(roles, domain.RoleSelfManagementGlobal, c.zitadelRoles)) > 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "Org-4N8es", "Errors.Org.MemberInvalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
				if exists, err := ExistsUser(ctx, filter, userID, ""); err != nil || !exists {
					return nil, zerrors.ThrowPreconditionFailed(err, "ORG-GoXOn", "Errors.User.NotFound")
				}
				if isMember, err := IsOrgMember(ctx, filter, a.ID, userID); err != nil || isMember {
					return nil, zerrors.ThrowAlreadyExists(err, "ORG-poWwe", "Errors.Org.Member.AlreadyExists")
				}
				return []eventstore.Command{org.NewMemberAddedEvent(ctx, &a.Aggregate, userID, roles...)}, nil
			},
			nil
	}
}

func IsOrgMember(ctx context.Context, filter preparation.FilterToQueryReducer, orgID, userID string) (isMember bool, err error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(orgID).
		OrderAsc().
		AddQuery().
		AggregateIDs(orgID).
		AggregateTypes(org.AggregateType).
		EventTypes(
			org.MemberAddedEventType,
			org.MemberRemovedEventType,
			org.MemberCascadeRemovedEventType,
		).Builder())
	if err != nil {
		return false, err
	}

	for _, event := range events {
		switch e := event.(type) {
		case *org.MemberAddedEvent:
			if e.UserID == userID {
				isMember = true
			}
		case *org.MemberRemovedEvent:
			if e.UserID == userID {
				isMember = false
			}
		case *org.MemberCascadeRemovedEvent:
			if e.UserID == userID {
				isMember = false
			}
		}
	}

	return isMember, nil
}

func (c *Commands) AddOrgMember(ctx context.Context, orgID, userID string, roles ...string) (*domain.Member, error) {
	orgAgg := org.NewAggregate(orgID)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.AddOrgMemberCommand(orgAgg, userID, roles...))
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	addedMember := NewOrgMemberWriteModel(orgID, userID)
	err = AppendAndReduce(addedMember, events...)
	if err != nil {
		return nil, err
	}
	return memberWriteModelToMember(&addedMember.MemberWriteModel), nil
}

func (c *Commands) addOrgMember(ctx context.Context, orgAgg *eventstore.Aggregate, addedMember *OrgMemberWriteModel, member *domain.Member) (eventstore.Command, error) {
	if !member.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-W8m4l", "Errors.Org.MemberInvalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.OrgRolePrefix, c.zitadelRoles)) > 0 && len(domain.CheckForInvalidRoles(member.Roles, domain.RoleSelfManagementGlobal, c.zitadelRoles)) > 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-4N8es", "Errors.Org.MemberInvalid")
	}
	err := c.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.State == domain.MemberStateActive {
		return nil, zerrors.ThrowAlreadyExists(nil, "Org-PtXi1", "Errors.Org.Member.AlreadyExists")
	}

	return org.NewMemberAddedEvent(ctx, orgAgg, member.UserID, member.Roles...), nil
}

// ChangeOrgMember updates an existing member
func (c *Commands) ChangeOrgMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	if !member.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-LiaZi", "Errors.Org.MemberInvalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.OrgRolePrefix, c.zitadelRoles)) > 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "IAM-m9fG8", "Errors.Org.MemberInvalid")
	}

	existingMember, err := c.orgMemberWriteModelByID(ctx, member.AggregateID, member.UserID)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "Org-LiaZi", "Errors.Org.Member.RolesNotChanged")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingMember.MemberWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewMemberChangedEvent(ctx, orgAgg, member.UserID, member.Roles...))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&existingMember.MemberWriteModel), nil
}

func (c *Commands) RemoveOrgMember(ctx context.Context, orgID, userID string) (*domain.ObjectDetails, error) {
	m, err := c.orgMemberWriteModelByID(ctx, orgID, userID)
	if err != nil && !zerrors.IsNotFound(err) {
		return nil, err
	}
	if zerrors.IsNotFound(err) {
		// empty response because we have no data that match the request
		return &domain.ObjectDetails{}, nil
	}

	orgAgg := OrgAggregateFromWriteModel(&m.MemberWriteModel.WriteModel)
	removeEvent := c.removeOrgMember(ctx, orgAgg, userID, false)
	pushedEvents, err := c.eventstore.Push(ctx, removeEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(m, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&m.WriteModel), nil
}

func (c *Commands) removeOrgMember(ctx context.Context, orgAgg *eventstore.Aggregate, userID string, cascade bool) eventstore.Command {
	if cascade {
		return org.NewMemberCascadeRemovedEvent(
			ctx,
			orgAgg,
			userID)
	} else {
		return org.NewMemberRemovedEvent(ctx, orgAgg, userID)
	}
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
		return nil, zerrors.ThrowNotFound(nil, "Org-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
