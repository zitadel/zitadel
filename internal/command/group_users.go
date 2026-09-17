package command

import (
	"context"
	"slices"
	"strings"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	repo "github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddUsersToGroup(ctx context.Context, groupID string, userIDs []string) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// precondition: check whether the group exists
	group, err := c.checkGroupExists(ctx, groupID, userIDs)
	if err != nil {
		return nil, err
	}

	// check whether the requester has permissions to add users to the group
	err = c.checkPermissionAddUserToGroup(ctx, group.ResourceOwner, group.AggregateID)
	if err != nil {
		return nil, err
	}

	// add the users to the group
	return c.addUsersToGroup(ctx, group)
}

func (c *Commands) RemoveUsersFromGroup(ctx context.Context, groupID string, userIDs []string) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// precondition: check whether the group exists
	group, err := c.checkGroupExists(ctx, groupID, userIDs)
	if err != nil {
		return nil, err
	}

	// check whether the requester has permissions to remove users from the group
	err = c.checkPermissionRemoveUserFromGroup(ctx, group.ResourceOwner, group.AggregateID)
	if err != nil {
		return nil, err
	}

	userIDsToRemove := group.getUserIDsToRemove()
	if len(userIDsToRemove) == 0 {
		// the userIDs are not present in the group; desired state achieved
		return writeModelToObjectDetails(&group.WriteModel), nil
	}

	// remove users from the group; emit one event per (group, user) pair so
	// each membership has its own creation / destruction record and the
	// eventstore unique constraint releases for that pair specifically
	groupAggregate := GroupAggregateFromWriteModel(ctx, &group.WriteModel)
	events := make([]eventstore.Command, 0, len(userIDsToRemove))
	for _, userID := range userIDsToRemove {
		events = append(events, repo.NewGroupUserRemovedEvent(ctx, groupAggregate, userID))
	}
	return c.pushAppendAndReduceDetails(ctx, group, events...)
}

func (c *Commands) addUsersToGroup(ctx context.Context, group *GroupWriteModel) (*domain.ObjectDetails, error) {
	userIDsToAdd := group.getUserIDsToAdd()
	if len(userIDsToAdd) == 0 {
		// no new users to add
		return writeModelToObjectDetails(&group.WriteModel), nil
	}

	// precondition: check that all users exist in the same organization as the group
	if err := c.checkUsersExist(ctx, userIDsToAdd, group.ResourceOwner); err != nil {
		return nil, err
	}

	// add users to the group; emit one event per (group, user) pair so the
	// eventstore can register a unique constraint per membership, matching
	// the pattern used by org / project / IAM MemberAddedEvent
	groupAggregate := GroupAggregateFromWriteModel(ctx, &group.WriteModel)
	events := make([]eventstore.Command, 0, len(userIDsToAdd))
	for _, userID := range userIDsToAdd {
		events = append(events, repo.NewGroupUserAddedEvent(ctx, groupAggregate, userID))
	}
	return c.pushAppendAndReduceDetails(ctx, group, events...)
}

// checkUsersExist verifies with a single eventstore query that every user exists
// in the given organization, reporting all missing user IDs at once instead of
// failing on the first one
func (c *Commands) checkUsersExist(ctx context.Context, userIDs []string, resourceOwner string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	wm := newUsersExistenceWriteModel(userIDs, resourceOwner)
	if err = c.eventstore.FilterToQueryReducer(ctx, wm); err != nil {
		return err
	}
	missing := wm.missingUserIDs()
	if len(missing) > 0 {
		return zerrors.ThrowPreconditionFailedf(nil, "CMDGRP-5jLqXz", "Errors.User.NotFound: %s", strings.Join(missing, ", "))
	}
	return nil
}

// usersExistenceWriteModel replays the lifecycle events of multiple users
// to verify their existence in a single eventstore query
type usersExistenceWriteModel struct {
	eventstore.WriteModel

	userIDs  []string
	existing map[string]bool
}

func newUsersExistenceWriteModel(userIDs []string, resourceOwner string) *usersExistenceWriteModel {
	return &usersExistenceWriteModel{
		WriteModel: eventstore.WriteModel{ResourceOwner: resourceOwner},
		userIDs:    userIDs,
		existing:   make(map[string]bool, len(userIDs)),
	}
}

func (wm *usersExistenceWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		OrderAsc().
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.userIDs...).
		EventTypes(
			user.UserV1AddedType,
			user.UserV1RegisteredType,
			user.HumanAddedType,
			user.HumanRegisteredType,
			user.MachineAddedEventType,
			user.UserRemovedType,
		).Builder()
}

func (wm *usersExistenceWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch event.(type) {
		case *user.HumanAddedEvent, *user.HumanRegisteredEvent, *user.MachineAddedEvent:
			wm.existing[event.Aggregate().ID] = true
		case *user.UserRemovedEvent:
			wm.existing[event.Aggregate().ID] = false
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *usersExistenceWriteModel) missingUserIDs() []string {
	missing := make([]string, 0, len(wm.userIDs))
	for _, userID := range wm.userIDs {
		if !wm.existing[userID] {
			missing = append(missing, userID)
		}
	}
	return missing
}

// getUserIDsToAdd returns the userIDs that are not already in the group
func (g *GroupWriteModel) getUserIDsToAdd() []string {
	userIDsToAdd := make([]string, 0)
	for _, userID := range g.UserIDs {
		if _, ok := g.existingUserIDs[userID]; !ok && !slices.Contains(userIDsToAdd, userID) {
			userIDsToAdd = append(userIDsToAdd, userID)
		}
	}
	return userIDsToAdd
}

// getUserIDsToRemove returns the userIDs that are in the group and should be removed
// if a userID is not in the group, the desired state has already been achieved
func (g *GroupWriteModel) getUserIDsToRemove() []string {
	userIDsToRemove := make([]string, 0)
	for _, userID := range g.UserIDs {
		if _, ok := g.existingUserIDs[userID]; ok && !slices.Contains(userIDsToRemove, userID) {
			userIDsToRemove = append(userIDsToRemove, userID)
		}
	}
	return userIDsToRemove
}

// removeUserFromGroups returns the events to remove a user from multiple groups.
// This is needed when a user is deleted and subsequently needs to be removed from all groups.
// Note: Ensure that the groupIDs are retrieved via SearchGroupUsers before calling this method
func (c *Commands) removeUserFromGroups(ctx context.Context, userID string, groupIDs []string, resourceOwner string) ([]eventstore.Command, error) {
	events := make([]eventstore.Command, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		events = append(
			events,
			repo.NewGroupUserRemovedEvent(
				ctx,
				&repo.NewAggregate(groupID, resourceOwner).Aggregate,
				userID,
			),
		)
	}
	return events, nil
}
