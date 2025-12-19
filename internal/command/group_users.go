package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	groupusers "github.com/zitadel/zitadel/internal/repository/group_users"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddGroupUser(ctx context.Context, user *domain.GroupUser, resourceOwner string) (_ *domain.GroupUser, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	addedUser := NewGroupUserWriteModel(user.AggregateID, user.UserID, resourceOwner, user.Attributes)

	groupAgg := GroupAggregateFromWriteModel(&addedUser.WriteModel)

	event, err := c.addGroupUser(ctx, groupAgg, addedUser, user)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(addedUser, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return groupUserWriteModelToUser(addedUser), nil
}

func (c *Commands) addGroupUser(ctx context.Context, groupAgg *eventstore.Aggregate, addedUser *GroupUserWriteModel, member *domain.GroupUser) (_ eventstore.Command, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if !member.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "GROUP-X9n3m", "Errors.Group.User.Invalid")
	}

	err = c.checkUserExists(ctx, addedUser.UserID, "")
	if err != nil {
		return nil, err
	}
	err = c.eventstore.FilterToQueryReducer(ctx, addedUser)
	if err != nil {
		return nil, err
	}
	if addedUser.State == domain.GroupUserStateActive {
		return nil, zerrors.ThrowAlreadyExists(nil, "GROUP-QvYJ2", "Errors.Group.User.AlreadyExists")
	}

	return groupusers.NewGroupUserAddedEvent(ctx, groupAgg, member.UserID, member.Attributes...), nil
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
func (c *Commands) ChangeGroupUser(ctx context.Context, user *domain.GroupUser, resourceOwner string) (*domain.GroupUser, error) {
	if !user.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "GROUP-MjbZi", "Errors.Group.User.Invalid")
	}

	existingUser, err := c.groupUserWriteModelByID(ctx, user.AggregateID, user.UserID, resourceOwner, user.Attributes)
	if err != nil {
		return nil, err
	}

	groupAgg := GroupAggregateFromWriteModel(&existingUser.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, groupusers.NewGroupUserChangedEvent(ctx, groupAgg, user.UserID, user.Attributes...))
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(existingUser, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return groupUserWriteModelToUser(existingUser), nil
}

func (c *Commands) RemoveGroupUser(ctx context.Context, groupID, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if groupID == "" || userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "GROUP-77nHd", "Errors.Group.User.Invalid")
	}
	m, err := c.groupUserWriteModelByID(ctx, groupID, userID, resourceOwner, nil)
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
	removeEvent := c.removeGroupUser(ctx, groupAgg, userID, false)

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

func (c *Commands) removeGroupUser(ctx context.Context, groupAgg *eventstore.Aggregate, userID string, cascade bool) eventstore.Command {
	if cascade {
		return groupusers.NewGroupUserCascadeRemovedEvent(
			ctx,
			groupAgg,
			userID)
	} else {
		return groupusers.NewGroupUserRemovedEvent(ctx, groupAgg, userID)
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

func (c *Commands) groupUserWriteModelByID(ctx context.Context, groupID, userID, resourceOwner string, attributes []string) (member *GroupUserWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewGroupUserWriteModel(groupID, userID, resourceOwner, attributes)
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
