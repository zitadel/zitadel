package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddGroupMember(ctx context.Context, member *domain.Member, resourceOwner string) (_ *domain.Member, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	addedMember := NewGroupMemberWriteModel(member.AggregateID, member.UserID, resourceOwner)
	addedGroupMember := NewUserGroupMemberWriteModel(member.AggregateID, member.UserID, resourceOwner)

	groupAgg := GroupAggregateFromWriteModel(&addedMember.WriteModel)
	userGroupAgg := GroupAggregateFromWriteModel(&addedGroupMember.WriteModel)

	event, err := c.addGroupMember(ctx, groupAgg, addedMember, member)
	if err != nil {
		return nil, err
	}

	userEvent, err := c.addUserGroupMember(ctx, userGroupAgg, addedGroupMember, member.AggregateID)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}

	userPushedEvents, err := c.eventstore.Push(ctx, userEvent)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(addedMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(addedGroupMember, userPushedEvents...)
	if err != nil {
		return nil, err
	}
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

func (c *Commands) addUserGroupMember(ctx context.Context, userGroupAgg *eventstore.Aggregate, addedGroup *UserGroupMemberWriteModel, group string) (_ eventstore.Command, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if group == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "GROUP-Y1n3m", "Errors.Group.Invalid")
	}

	err = c.checkGroupExists(ctx, addedGroup.GroupID, addedGroup.ResourceOwner)
	if err != nil {
		return nil, err
	}
	err = c.eventstore.FilterToQueryReducer(ctx, addedGroup)
	if err != nil {
		return nil, err
	}
	if addedGroup.State == domain.GroupStateActive {
		return nil, zerrors.ThrowAlreadyExists(nil, "GROUP-QvYJ2", "Errors.Group.Member.AlreadyExists")
	}

	return user.NewUserGroupAddedEvent(ctx, userGroupAgg, group), nil
}

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
	if err != nil && !zerrors.IsNotFound(err) {
		return nil, err
	}
	if zerrors.IsNotFound(err) {
		// empty response because we have no data that match the request
		return &domain.ObjectDetails{}, nil
	}

	groupAgg := GroupAggregateFromWriteModel(&m.WriteModel)
	removeEvent := c.removeGroupMember(ctx, groupAgg, userID, false)
	pushedEvents, err := c.eventstore.Push(ctx, removeEvent)
	if err != nil {
		return nil, err
	}
	/*
		// Remove group_id entry from users table
		err = c.removeGroupIDFromUser(ctx, userID, groupID)
		if err != nil {
			return nil, err
		}
	*/
	err = AppendAndReduce(m, pushedEvents...)
	if err != nil {
		return nil, err
	}

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

func (c *Commands) groupMemberWriteModelByID(ctx context.Context, groupID, userID, resourceOwner string) (member *GroupMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewGroupMemberWriteModel(groupID, userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "GROUP-E9KyS", "Errors.NotFound")
	}

	return writeModel, nil
}

/*
func (c *Commands) addGroupIDToUser(ctx context.Context, userID, groupID string) error {
	ctx, span := tracing.NewSpan(ctx)
	defer span.End()

	// Fetch the existing user
	user, err := c.query.GetUserByID(ctx, false, userID)
	if err != nil {
		return err
	}

	// Add the groupID to the user's group_ids
	user.GroupIDs = append(user.GroupIDs, groupID)

	// Update the user in the database
	userAgg := eventstore.NewAggregate(ctx, user.ID, eventstore.AggregateType(es_user.UserGroupAddedType), "v1")
	evt := es_user.NewUserGroupAddedEvent(ctx, userAgg, user.GroupIDs...)
	_, err = c.eventstore.Push(ctx, evt)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) removeGroupIDFromUser(ctx context.Context, userID, groupID string) error {
	ctx, span := tracing.NewSpan(ctx)
	defer span.End()

	// Fetch the existing user
	user, err := c.query.GetUserByID(ctx, false, userID)
	if err != nil {
		return err
	}

	// Remove the groupID from the user's group_ids
	for i, id := range user.GroupIDs {
		if id == groupID {
			user.GroupIDs = append(user.GroupIDs[:i], user.GroupIDs[i+1:]...)
			break
		}
	}

	// Update the user in the database
	userAgg := eventstore.NewAggregate(ctx, user.ID, eventstore.AggregateType(es_user.UserGroupRemovedType), "v1")
	_, err = c.eventstore.Push(ctx, es_user.NewUserGroupRemovedEvent(ctx, userAgg, user.GroupIDs...))
	if err != nil {
		return err
	}

	return nil
}
*/
