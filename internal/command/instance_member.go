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
			return nil, zerrors.ThrowInvalidArgument(nil, "INSTANCE-4m0fS", "Errors.IAM.MemberInvalid")
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

func (c *Commands) AddInstanceMember(ctx context.Context, member *AddInstanceMember) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(member.InstanceID)
	if err := c.checkPermissionUpdateInstanceMember(ctx, member.InstanceID); err != nil {
		return nil, err
	}
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
		return zerrors.ThrowInvalidArgument(nil, "INSTANCE-LiaZi", "Errors.IAM.MemberInvalid")
	}
	if len(domain.CheckForInvalidRoles(i.Roles, domain.IAMRolePrefix, zitadelRoles)) > 0 {
		return zerrors.ThrowInvalidArgument(nil, "INSTANCE-3m9fs", "Errors.IAM.MemberInvalid")
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
