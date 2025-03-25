package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddGroupMember(ctx context.Context, member *domain.Member, resourceOwner string) (_ *domain.Member, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	addedMember := NewGroupMemberWriteModel(member.AggregateID, member.UserID, resourceOwner)
	// addedGroupMember := NewUserGroupMemberWriteModel(member.AggregateID, member.UserID, resourceOwner)

	groupAgg := GroupAggregateFromWriteModel(&addedMember.WriteModel)
	// userGroupAgg := GroupAggregateFromWriteModel(&addedGroupMember.WriteModel)

	event, err := c.addGroupMember(ctx, groupAgg, addedMember, member)
	if err != nil {
		return nil, err
	}

	// userEvent, err := c.addUserGroupMember(ctx, userGroupAgg, addedGroupMember, member.AggregateID)
	// if err != nil {
	// 	return nil, err
	// }

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}

	// userPushedEvents, err := c.eventstore.Push(ctx, userEvent)
	// if err != nil {
	// 	return nil, err
	// }

	err = AppendAndReduce(addedMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	// err = AppendAndReduce(addedGroupMember, userPushedEvents...)
	// if err != nil {
	// 	return nil, err
	// }
	return groupMemberWriteModelToMember(&addedMember.MemberWriteModel), nil
}

func (c *Commands) addGroupMember(ctx context.Context, groupAgg *eventstore.Aggregate, addedMember *GroupMemberWriteModel, member *domain.Member) (_ eventstore.Command, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if !member.IsGroupMemberValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "GROUP-X9n3m", "Errors.Group.Member.Invalid")
	}

	err = c.checkUserExists(ctx, addedMember.UserID, "")
	if err != nil {
		return nil, err
	}
	err = c.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.State == domain.MemberStateActive {
		return nil, zerrors.ThrowAlreadyExists(nil, "GROUP-QvYJ2", "Errors.Group.Member.AlreadyExists")
	}

	return group.NewGroupMemberAddedEvent(ctx, groupAgg, member.UserID), nil
}

// func (c *Commands) addUserGroupMember(ctx context.Context, userGroupAgg *eventstore.Aggregate, addedGroup *UserGroupMemberWriteModel, group string) (_ eventstore.Command, err error) {
// 	ctx, span := tracing.NewSpan(ctx)
// 	defer func() { span.EndWithError(err) }()

// 	if group == "" {
// 		return nil, zerrors.ThrowInvalidArgument(nil, "GROUP-Y1n3m", "Errors.Group.Invalid")
// 	}

// 	err = c.checkGroupExists(ctx, addedGroup.GroupID, addedGroup.ResourceOwner)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = c.eventstore.FilterToQueryReducer(ctx, addedGroup)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if addedGroup.State == domain.GroupStateActive {
// 		return nil, zerrors.ThrowAlreadyExists(nil, "GROUP-QvYJ2", "Errors.Group.Member.AlreadyExists")
// 	}

// 	return user.NewUserGroupAddedEvent(ctx, userGroupAgg, group), nil
// }

// Need to update these to remove user from the users13 table and group_members table
// ChangeGroupMember updates an existing member
func (c *Commands) ChangeGroupMember(ctx context.Context, member *domain.Member, resourceOwner string) (*domain.Member, error) {
	if !member.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "GROUP-MjbZi", "Errors.Group.Member.Invalid")
	}

	existingMember, err := c.groupMemberWriteModelByID(ctx, member.AggregateID, member.UserID, resourceOwner)
	if err != nil {
		return nil, err
	}

	groupAgg := GroupAggregateFromWriteModel(&existingMember.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, group.NewGroupMemberChangedEvent(ctx, groupAgg, member.UserID))
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return groupMemberWriteModelToMember(&existingMember.MemberWriteModel), nil
}

func (c *Commands) RemoveGroupMember(ctx context.Context, groupID, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if groupID == "" || userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "GROUP-77nHd", "Errors.Group.Member.Invalid")
	}
	m, err := c.groupMemberWriteModelByID(ctx, groupID, userID, resourceOwner)
	if zerrors.IsNotFound(err) {
		// empty response because we have no data that match the request
		return &domain.ObjectDetails{}, err
	}

	// um, err := c.userGroupMemberWriteModelByID(ctx, groupID, userID, resourceOwner)
	// if zerrors.IsNotFound(err) {
	// 	return &domain.ObjectDetails{}, err
	// }

	groupAgg := GroupAggregateFromWriteModel(&m.WriteModel)
	// userGroupAgg := GroupAggregateFromWriteModel(&um.WriteModel)
	removeEvent := c.removeGroupMember(ctx, groupAgg, userID, false)

	// userEvents := c.removeUserGroupMember(ctx, userGroupAgg, groupID, false)

	pushedEvents, err := c.eventstore.Push(ctx, removeEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(m, pushedEvents...)
	if err != nil {
		return nil, err
	}

	// userPushedEvents, err := c.eventstore.Push(ctx, userEvents)
	// if err != nil {
	// 	return nil, err
	// }

	// err = AppendAndReduce(um, userPushedEvents...)
	// if err != nil {
	// 	return nil, err
	// }

	return writeModelToObjectDetails(&m.WriteModel), nil
}

func (c *Commands) removeGroupMember(ctx context.Context, groupAgg *eventstore.Aggregate, userID string, cascade bool) eventstore.Command {
	if cascade {
		return group.NewGroupMemberCascadeRemovedEvent(
			ctx,
			groupAgg,
			userID)
	} else {
		return group.NewGroupMemberRemovedEvent(ctx, groupAgg, userID)
	}
}

// func (c *Commands) removeUserGroupMember(ctx context.Context, userGroupAgg *eventstore.Aggregate, groupID string, cascade bool) eventstore.Command {
// 	if cascade {
// 		return user.NewUserGroupCascadeRemovedEvent(
// 			ctx,
// 			userGroupAgg,
// 			groupID)
// 	} else {
// 		return user.NewUserGroupRemovedEvent(ctx, userGroupAgg, groupID)
// 	}
// }

func (c *Commands) groupMemberWriteModelByID(ctx context.Context, groupID, userID, resourceOwner string) (member *GroupMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewGroupMemberWriteModel(groupID, userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	// if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
	// 	return nil, zerrors.ThrowNotFound(nil, "GROUP-E9KyS", "Errors.NotFound")
	// }

	return writeModel, nil
}

// func (c *Commands) userGroupMemberWriteModelByID(ctx context.Context, groupID, userID, resourceOwner string) (userGroup *UserGroupMemberWriteModel, err error) {
// 	ctx, span := tracing.NewSpan(ctx)
// 	defer func() { span.EndWithError(err) }()

// 	writeModel := NewUserGroupMemberWriteModel(groupID, userID, resourceOwner)
// 	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// if writeModel.State == domain.GroupStateUnspecified || writeModel.State == domain.GroupStateRemoved {
// 	// 	return nil, zerrors.ThrowNotFound(nil, "GROUP-F1KyS", "Errors.NotFound")
// 	// }

// 	return writeModel, nil
// }
