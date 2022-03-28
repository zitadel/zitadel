package command

import (
	"context"
	"reflect"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddInstanceMember(ctx context.Context, instanceID string, member *domain.Member) (*domain.Member, error) {
	if member.UserID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-Mf83b", "Errors.IAM.MemberInvalid")
	}
	addedMember := NewInstanceMemberWriteModel(instanceID, member.UserID)
	instanceAgg := InstanceAggregateFromWriteModel(&addedMember.MemberWriteModel.WriteModel)
	err := c.checkUserExists(ctx, addedMember.UserID, "")
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "INSTANCE-5N9vs", "Errors.User.NotFound")
	}
	event, err := c.addInstanceMember(ctx, instanceAgg, addedMember, member)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedMember, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return memberWriteModelToMember(&addedMember.MemberWriteModel), nil
}

func (c *Commands) addInstanceMember(ctx context.Context, instanceAgg *eventstore.Aggregate, addedMember *InstanceMemberWriteModel, member *domain.Member) (eventstore.Command, error) {
	if !member.IsIAMValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-GR34U", "Errors.IAM.MemberInvalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.IAMRolePrefix, c.zitadelRoles)) > 0 {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-4m0fS", "Errors.IAM.MemberInvalid")
	}
	err := c.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.State == domain.MemberStateActive {
		return nil, errors.ThrowAlreadyExists(nil, "INSTANCE-sdgQ4", "Errors.IAM.Member.AlreadyExists")
	}

	return instance.NewMemberAddedEvent(ctx, instanceAgg, member.UserID, member.Roles...), nil
}

//ChangeInstanceMember updates an existing member
func (c *Commands) ChangeInstanceMember(ctx context.Context, instanceID string, member *domain.Member) (*domain.Member, error) {
	if !member.IsIAMValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-LiaZi", "Errors.IAM.MemberInvalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.IAMRolePrefix, c.zitadelRoles)) > 0 {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-3m9fs", "Errors.IAM.MemberInvalid")
	}

	existingMember, err := c.instanceMemberWriteModelByID(ctx, instanceID, member.UserID)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "INSTANCE-LiaZi", "Errors.IAM.Member.RolesNotChanged")
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

func (c *Commands) RemoveInstanceMember(ctx context.Context, instanceID, userID string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-LiaZi", "Errors.IDMissing")
	}
	memberWriteModel, err := c.instanceMemberWriteModelByID(ctx, instanceID, userID)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if errors.IsNotFound(err) {
		return nil, nil
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

func (c *Commands) instanceMemberWriteModelByID(ctx context.Context, instanceID, userID string) (member *InstanceMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewInstanceMemberWriteModel(instanceID, userID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, errors.ThrowNotFound(nil, "INSTANCE-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
