package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddOrgMemberCommand(member *AddOrgMember) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if err := member.IsValid(c.zitadelRoles); err != nil {
			return nil, err
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) (_ []eventstore.Command, err error) {
				ctx, span := tracing.NewSpan(ctx)
				defer func() { span.EndWithError(err) }()

				if exists, err := ExistsUser(ctx, filter, member.UserID, "", false); err != nil || !exists {
					return nil, zerrors.ThrowPreconditionFailed(err, "ORG-GoXOn", "Errors.User.NotFound")
				}
				if isMember, err := IsOrgMember(ctx, filter, member.OrgID, member.UserID); err != nil || isMember {
					return nil, zerrors.ThrowAlreadyExists(err, "ORG-poWwe", "Errors.Org.Member.AlreadyExists")
				}
				return []eventstore.Command{org.NewMemberAddedEvent(ctx, &org.NewAggregate(member.OrgID).Aggregate, member.UserID, member.Roles...)}, nil
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

type AddOrgMember struct {
	OrgID  string
	UserID string
	Roles  []string
}

func (m *AddOrgMember) IsValid(zitadelRoles []authz.RoleMapping) error {
	if m.UserID == "" || m.OrgID == "" || len(m.Roles) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "ORG-4Mlfs", "Errors.Invalid.Argument")
	}
	if len(domain.CheckForInvalidRoles(m.Roles, domain.OrgRolePrefix, zitadelRoles)) > 0 && len(domain.CheckForInvalidRoles(m.Roles, domain.RoleSelfManagementGlobal, zitadelRoles)) > 0 {
		return zerrors.ThrowInvalidArgument(nil, "Org-4N8es", "Errors.Org.MemberInvalid")
	}
	return nil
}

func (c *Commands) AddOrgMember(ctx context.Context, member *AddOrgMember) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	if err := c.checkOrgExists(ctx, member.OrgID); err != nil {
		return nil, err
	}
	if err := c.checkPermissionUpdateOrgMember(ctx, member.OrgID, member.OrgID); err != nil {
		return nil, err
	}
	//nolint:staticcheck
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.AddOrgMemberCommand(member))
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	addedMember := NewOrgMemberWriteModel(member.OrgID, member.UserID)
	err = AppendAndReduce(addedMember, events...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&addedMember.WriteModel), nil
}

type ChangeOrgMember struct {
	OrgID  string
	UserID string
	Roles  []string
}

func (c *ChangeOrgMember) IsValid(zitadelRoles []authz.RoleMapping) error {
	if c.OrgID == "" || c.UserID == "" || len(c.Roles) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "Org-LiaZi", "Errors.Org.MemberInvalid")
	}
	if len(domain.CheckForInvalidRoles(c.Roles, domain.OrgRolePrefix, zitadelRoles)) > 0 {
		return zerrors.ThrowInvalidArgument(nil, "IAM-m9fG8", "Errors.Org.MemberInvalid")
	}

	return nil
}

// ChangeOrgMember updates an existing member
func (c *Commands) ChangeOrgMember(ctx context.Context, member *ChangeOrgMember) (*domain.ObjectDetails, error) {
	if err := member.IsValid(c.zitadelRoles); err != nil {
		return nil, err
	}

	existingMember, err := c.orgMemberWriteModelByID(ctx, member.OrgID, member.UserID)
	if err != nil {
		return nil, err
	}
	if !existingMember.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "Org-D8JxR", "Errors.NotFound")
	}
	if err := c.checkPermissionUpdateOrgMember(ctx, existingMember.ResourceOwner, existingMember.AggregateID); err != nil {
		return nil, err
	}

	if slices.Compare(existingMember.Roles, member.Roles) == 0 {
		return writeModelToObjectDetails(&existingMember.WriteModel), nil
	}

	pushedEvents, err := c.eventstore.Push(ctx,
		org.NewMemberChangedEvent(ctx,
			OrgAggregateFromWriteModelWithCTX(ctx, &existingMember.WriteModel),
			member.UserID,
			member.Roles...,
		),
	)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&existingMember.WriteModel), nil
}

func (c *Commands) RemoveOrgMember(ctx context.Context, orgID, userID string) (*domain.ObjectDetails, error) {
	if orgID == "" || userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-LiaZi", "Errors.Org.MemberInvalid")
	}
	existingMember, err := c.orgMemberWriteModelByID(ctx, orgID, userID)
	if err != nil {
		return nil, err
	}
	if !existingMember.State.Exists() {
		return writeModelToObjectDetails(&existingMember.WriteModel), nil
	}
	if err := c.checkPermissionDeleteOrgMember(ctx, existingMember.ResourceOwner, existingMember.AggregateID); err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx,
		c.removeOrgMember(ctx,
			OrgAggregateFromWriteModelWithCTX(ctx, &existingMember.WriteModel),
			userID,
			false,
		),
	)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingMember.WriteModel), nil
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

	return writeModel, nil
}
