package command

import (
	"context"
	"reflect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func (c *Commands) AddInstanceMember(ctx context.Context, userID string, roles ...string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	return c.processWithFirst(ctx, c.PrepareAddInstanceMember(instanceAgg, userID, roles...))
}

func (c *Commands) PrepareAddInstanceMember(a *instance.Aggregate, userID string, roles ...string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if userID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INSTA-SDSfs", "Errors.Invalid.Argument")
		}
		if len(domain.CheckForInvalidRoles(roles, domain.IAMRolePrefix, c.zitadelRoles)) > 0 {
			return nil, errors.ThrowInvalidArgument(nil, "INSTANCE-4m0fS", "Errors.IAM.MemberInvalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
				if exists, err := ExistsUser(ctx, filter, userID, ""); err != nil || !exists {
					return nil, errors.ThrowPreconditionFailed(err, "INSTA-GSXOn", "Errors.User.NotFound")
				}
				if isMember, err := IsInstanceMember(ctx, filter, a.ID, userID); err != nil || isMember {
					return nil, errors.ThrowAlreadyExists(err, "INSTA-pFDwe", "Errors.Instance.Member.AlreadyExists")
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

// ChangeInstanceMember updates an existing member
func (c *Commands) ChangeInstanceMember(ctx context.Context, member *domain.Member) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	return c.processWithFirst(ctx, c.prepareChangeInstanceMember(instanceAgg, member))
}

func (c *Commands) prepareChangeInstanceMember(a *instance.Aggregate, member *domain.Member) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if !member.IsIAMValid() {
			return nil, errors.ThrowInvalidArgument(nil, "INSTANCE-LiaZi", "Errors.IAM.MemberInvalid")
		}
		if len(domain.CheckForInvalidRoles(member.Roles, domain.IAMRolePrefix, c.zitadelRoles)) > 0 {
			return nil, errors.ThrowInvalidArgument(nil, "INSTANCE-3m9fs", "Errors.IAM.MemberInvalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
				writeModel, err := getInstanceMemberWriteModel(ctx, filter, member.UserID)
				if err != nil {
					return nil, err
				}
				if !isMemberExisting(writeModel.State) {
					return nil, errors.ThrowNotFound(nil, "INSTANCE-D8JxR", "Errors.NotFound")
				}
				if reflect.DeepEqual(writeModel.Roles, member.Roles) {
					return nil, errors.ThrowPreconditionFailed(nil, "INSTANCE-LiaZi", "Errors.IAM.Member.RolesNotChanged")
				}
				return []eventstore.Command{
					instance.NewMemberChangedEvent(ctx, &a.Aggregate, member.UserID, member.Roles...),
				}, nil
			},
			nil
	}
}

func (c *Commands) RemoveInstanceMember(ctx context.Context, userID string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	return c.processWithFirst(ctx, c.prepareRemoveInstanceMember(instanceAgg, userID))
}

func (c *Commands) prepareRemoveInstanceMember(a *instance.Aggregate, userID string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if userID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INSTANCE-LiaZi", "Errors.IDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
				writeModel, err := getInstanceMemberWriteModel(ctx, filter, userID)
				if err != nil {
					return nil, err
				}
				if !isMemberExisting(writeModel.State) {
					return nil, errors.ThrowNotFound(nil, "INSTANCE-98201h", "Errors.NotFound")
				}
				return []eventstore.Command{
					c.removeInstanceMember(ctx, &a.Aggregate, userID, false),
				}, nil
			},
			nil
	}
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

func getInstanceMemberWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer, userID string) (*InstanceMemberWriteModel, error) {
	writeModel := NewInstanceMemberWriteModel(ctx, userID)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModel, nil
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}
