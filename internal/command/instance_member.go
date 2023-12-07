package command

import (
	"context"
	"reflect"

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
				if exists, err := ExistsUser(ctx, filter, userID, ""); err != nil || !exists {
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

func (c *Commands) AddInstanceMember(ctx context.Context, userID string, roles ...string) (*domain.Member, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.AddInstanceMemberCommand(instanceAgg, userID, roles...))
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	addedMember := NewInstanceMemberWriteModel(ctx, userID)
	err = AppendAndReduce(addedMember, events...)
	if err != nil {
		return nil, err
	}
	return memberWriteModelToMember(&addedMember.MemberWriteModel), nil
}

// ChangeInstanceMember updates an existing member
func (c *Commands) ChangeInstanceMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	if !member.IsIAMValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTANCE-LiaZi", "Errors.IAM.MemberInvalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.IAMRolePrefix, c.zitadelRoles)) > 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTANCE-3m9fs", "Errors.IAM.MemberInvalid")
	}

	existingMember, err := c.instanceMemberWriteModelByID(ctx, member.UserID)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "INSTANCE-LiaZi", "Errors.IAM.Member.RolesNotChanged")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingMember.MemberWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewMemberChangedEvent(ctx, instanceAgg, member.UserID, member.Roles...))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&existingMember.MemberWriteModel), nil
}

func (c *Commands) RemoveInstanceMember(ctx context.Context, userID string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTANCE-LiaZi", "Errors.IDMissing")
	}
	memberWriteModel, err := c.instanceMemberWriteModelByID(ctx, userID)
	if err != nil && !zerrors.IsNotFound(err) {
		return nil, err
	}
	if zerrors.IsNotFound(err) {
		// empty response because we have no data that match the request
		return &domain.ObjectDetails{}, nil
	}

	instanceAgg := InstanceAggregateFromWriteModel(&memberWriteModel.MemberWriteModel.WriteModel)
	removeEvent := c.removeInstanceMember(ctx, instanceAgg, userID, false)
	pushedEvents, err := c.eventstore.Push(ctx, removeEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(memberWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&memberWriteModel.MemberWriteModel.WriteModel), nil
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

func (c *Commands) instanceMemberWriteModelByID(ctx context.Context, userID string) (member *InstanceMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewInstanceMemberWriteModel(ctx, userID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
