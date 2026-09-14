package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddInstanceMemberCommand(a *instance.Aggregate, userID string, roles ...string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if userID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INSTA-SDSfs", "Errors.Invalid.Argument")
		}
		if len(domain.CheckForInvalidRoles(roles, domain.IAMRolePrefix, c.zitadelRoles)) > 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "INSTANCE-4m0fS", "Errors.Instance.MemberInvalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
				if exists, err := ExistsUser(ctx, filter, userID, "", false); err != nil || !exists {
					return nil, zerrors.ThrowPreconditionFailed(err, "INSTA-GSXOn", "Errors.User.NotFound")
				}
				if isMember, err := IsInstanceMember(ctx, filter, a.ID, userID); err != nil || isMember {
					return nil, zerrors.ThrowAlreadyExists(err, "INSTA-pFDwe", "Errors.Instance.Member.AlreadyExists")
				}
				return []eventstore.Command{instance.NewMemberAddedEvent(ctx, &a.Aggregate, userID, roles...)}, nil
			},
			nil
	}
}

func IsInstanceMember(ctx context.Context, filter preparation.FilterToQueryReducer, instanceID, userID string) (isMember bool, err error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderAsc().
		AddQuery().
		AggregateIDs(instanceID).
		AggregateTypes(instance.AggregateType).
		EventTypes(
			instance.MemberAddedEventType,
			instance.MemberRemovedEventType,
			instance.MemberCascadeRemovedEventType,
		).Builder())
	if err != nil {
		return false, err
	}

	for _, event := range events {
		switch e := event.(type) {
		case *instance.MemberAddedEvent:
			if e.UserID == userID {
				isMember = true
			}
		case *instance.MemberRemovedEvent:
			if e.UserID == userID {
				isMember = false
			}
		case *instance.MemberCascadeRemovedEvent:
			if e.UserID == userID {
				isMember = false
			}
		}
	}

	return isMember, nil
}

type AddInstanceMember struct {
	InstanceID string
	UserID     string
	Roles      []string
}

// AddInstanceMember adds a user as a member of the instance with the given roles.
func (c *Commands) AddInstanceMember(ctx context.Context, member *AddInstanceMember) (*domain.ObjectDetails, error) {
	if err := c.checkPermissionUpdateInstanceMember(ctx, member.InstanceID); err != nil {
		return nil, err
	}
	return c.addInstanceMember(ctx, member)
}

// EnsureInstanceMemberRolesFromLogin makes sure the user holds the given instance member roles
// during the login v1 ZITADEL-IdP flow.
// It adds the membership if the user is not a member yet and otherwise adds the missing roles
// to the existing membership. Roles the user already holds are never removed.
//
// It intentionally bypasses the standard instance-member permission check.
// Authorization is established by the validated `urn:zitadel:iam:org:project:roles` token claim
// and the configured InstanceRolesInfo for the Zitadel provider.
// Do not use it for any user-initiated API path.
func (c *Commands) EnsureInstanceMemberRolesFromLogin(ctx context.Context, member *AddInstanceMember) (*domain.ObjectDetails, error) {
	if member.InstanceID == "" || member.UserID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTA-Ee7ai", "Errors.Instance.MemberInvalid")
	}
	granted := configuredRolesOnly(member.Roles, c.zitadelRoles)
	if len(granted) == 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTA-oo0Ai", "Errors.Instance.MemberInvalid")
	}
	existingMember, err := c.instanceMemberWriteModelByID(ctx, member.InstanceID, member.UserID)
	if err != nil {
		return nil, err
	}
	if !existingMember.State.Exists() {
		return c.addInstanceMember(ctx, &AddInstanceMember{
			InstanceID: member.InstanceID,
			UserID:     member.UserID,
			Roles:      granted,
		})
	}
	// the roles are merged so that roles granted elsewhere are retained.
	roles := missingRolesAdded(existingMember.Roles, granted)
	if len(roles) == len(existingMember.Roles) {
		return writeModelToObjectDetails(&existingMember.WriteModel), nil
	}
	pushedEvents, err := c.eventstore.Push(ctx,
		instance.NewMemberChangedEvent(ctx,
			InstanceAggregateFromWriteModel(&existingMember.WriteModel),
			member.UserID,
			roles...,
		),
	)
	if err != nil {
		return nil, err
	}
	if err = AppendAndReduce(existingMember, pushedEvents...); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingMember.WriteModel), nil
}

// configuredRolesOnly returns the subset of roles that are configured as instance roles.
// The roles of a ZITADEL provider's claim originate from an external
// instance, whose role set may differ, so unknown roles are dropped.
func configuredRolesOnly(roles []string, zitadelRoles []authz.RoleMapping) []string {
	invalid := domain.CheckForInvalidRoles(roles, domain.IAMRolePrefix, zitadelRoles)
	return slices.DeleteFunc(slices.Clone(roles), func(role string) bool {
		return slices.Contains(invalid, role)
	})
}

// missingRolesAdded returns the existing roles with new roles to be added that are not already
// present, preserving the order of the existing roles.
func missingRolesAdded(existing, add []string) []string {
	roles := slices.Clone(existing)
	for _, role := range add {
		if !slices.Contains(roles, role) {
			roles = append(roles, role)
		}
	}
	return roles
}

// addInstanceMember performs the write without any permission check.
func (c *Commands) addInstanceMember(ctx context.Context, member *AddInstanceMember) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(member.InstanceID)
	//nolint:staticcheck
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.AddInstanceMemberCommand(instanceAgg, member.UserID, member.Roles...))
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	addedMember := NewInstanceMemberWriteModel(member.InstanceID, member.UserID)
	err = AppendAndReduce(addedMember, events...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&addedMember.WriteModel), nil
}

type ChangeInstanceMember struct {
	InstanceID string
	UserID     string
	Roles      []string
}

func (i *ChangeInstanceMember) IsValid(zitadelRoles []authz.RoleMapping) error {
	if i.InstanceID == "" || i.UserID == "" || len(i.Roles) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "INSTANCE-LiaZi", "Errors.Instance.MemberInvalid")
	}
	if len(domain.CheckForInvalidRoles(i.Roles, domain.IAMRolePrefix, zitadelRoles)) > 0 {
		return zerrors.ThrowInvalidArgument(nil, "INSTANCE-3m9fs", "Errors.Instance.MemberInvalid")
	}
	return nil
}

// ChangeInstanceMember updates an existing member
func (c *Commands) ChangeInstanceMember(ctx context.Context, member *ChangeInstanceMember) (*domain.ObjectDetails, error) {
	if err := member.IsValid(c.zitadelRoles); err != nil {
		return nil, err
	}

	existingMember, err := c.instanceMemberWriteModelByID(ctx, member.InstanceID, member.UserID)
	if err != nil {
		return nil, err
	}
	if !existingMember.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-D8JxR", "Errors.NotFound")
	}
	if err := c.checkPermissionUpdateInstanceMember(ctx, existingMember.AggregateID); err != nil {
		return nil, err
	}
	if slices.Compare(existingMember.Roles, member.Roles) == 0 {
		return writeModelToObjectDetails(&existingMember.WriteModel), nil
	}
	pushedEvents, err := c.eventstore.Push(ctx,
		instance.NewMemberChangedEvent(ctx,
			InstanceAggregateFromWriteModel(&existingMember.WriteModel),
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

func (c *Commands) RemoveInstanceMember(ctx context.Context, instanceID, userID string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTANCE-LiaZi", "Errors.IDMissing")
	}
	existingMember, err := c.instanceMemberWriteModelByID(ctx, instanceID, userID)
	if err != nil {
		return nil, err
	}
	if err := c.checkPermissionDeleteInstanceMember(ctx, instanceID); err != nil {
		return nil, err
	}
	if !existingMember.State.Exists() {
		return writeModelToObjectDetails(&existingMember.WriteModel), nil
	}

	pushedEvents, err := c.eventstore.Push(ctx,
		c.removeInstanceMember(ctx,
			InstanceAggregateFromWriteModel(&existingMember.WriteModel),
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

func (c *Commands) removeInstanceMember(ctx context.Context, instanceAgg *eventstore.Aggregate, userID string, cascade bool) eventstore.Command {
	if cascade {
		return instance.NewMemberCascadeRemovedEvent(
			ctx,
			instanceAgg,
			userID)
	} else {
		return instance.NewMemberRemovedEvent(ctx, instanceAgg, userID)
	}
}

func (c *Commands) instanceMemberWriteModelByID(ctx context.Context, instanceID, userID string) (member *InstanceMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewInstanceMemberWriteModel(instanceID, userID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}
