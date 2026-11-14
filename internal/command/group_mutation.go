package command

import (
	"context"
	"maps"
	"slices"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	repo "github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// GroupMutation describes an ordered change to a group's name or membership.
// A non-nil Members replaces the entire membership, including an empty slice.
type GroupMutation struct {
	Name          *string
	ExternalID    *GroupExternalID
	Members       *[]string
	AddMembers    []string
	RemoveMembers []string
}

// maxGroupMutations bounds the work done while the instance lock is held.
const maxGroupMutations = 1000

// MutateGroup validates all mutations before persisting their events in one transaction.
func (c *Commands) MutateGroup(ctx context.Context, groupID, orgID string, mutations []GroupMutation) (*domain.ObjectDetails, error) {
	if len(mutations) > maxGroupMutations {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "CMDGRP-mut3", "at most %d group mutations are allowed", maxGroupMutations)
	}
	return c.retryGroupChange(ctx, func(ctx context.Context) (*domain.ObjectDetails, error) {
		return c.mutateGroupOnce(ctx, groupID, orgID, mutations)
	})
}

func (c *Commands) mutateGroupOnce(ctx context.Context, groupID, orgID string, mutations []GroupMutation) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "CMDGRP-mut2", "Errors.Group.MissingOrganizationID")
	}
	group, err := c.getGroupWriteModelByID(ctx, groupID, orgID, []string{})
	if err != nil {
		return nil, err
	}
	if !group.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "CMDGRP-mut1", "Errors.Group.NotFound")
	}
	name := group.Name
	externalIDs := maps.Clone(group.ExternalIDs)
	members := maps.Clone(group.existingUserIDs)
	if err := c.validateGroupMutations(ctx, groupID, orgID, mutations); err != nil {
		return nil, err
	}
	for _, mutation := range mutations {
		if mutation.ExternalID != nil {
			externalIDs[mutation.ExternalID.Scope] = mutation.ExternalID.Value
		}
		if mutation.Name != nil {
			name = *mutation.Name
		}
		if mutation.Members != nil {
			clear(members)
			for _, id := range *mutation.Members {
				members[id] = struct{}{}
			}
		}
		for _, id := range mutation.AddMembers {
			members[id] = struct{}{}
		}
		for _, id := range mutation.RemoveMembers {
			delete(members, id)
		}
	}
	aggregate := GroupAggregateFromWriteModel(ctx, &group.WriteModel)
	events := make([]eventstore.Command, 0)
	if changed := group.NewChangedEvent(ctx, aggregate, &name, nil); changed != nil {
		events = append(events, changed)
	}
	for _, id := range slices.Sorted(maps.Keys(members)) {
		if _, exists := group.existingUserIDs[id]; !exists {
			events = append(events, repo.NewGroupUserAddedEvent(ctx, aggregate, id))
		}
	}
	for _, id := range slices.Sorted(maps.Keys(group.existingUserIDs)) {
		if _, exists := members[id]; !exists {
			events = append(events, repo.NewGroupUserRemovedEvent(ctx, aggregate, id))
		}
	}
	events = append(events, group.externalIDChanges(ctx, externalIDs)...)
	if len(events) == 0 {
		return writeModelToObjectDetails(&group.WriteModel), nil
	}
	return c.pushAppendAndReduceDetails(ctx, group, events...)
}

// validateGroupMutations runs each permission check and the user existence query
// at most once, because they run while the instance lock is held.
func (c *Commands) validateGroupMutations(ctx context.Context, groupID, orgID string, mutations []GroupMutation) error {
	var update, add, remove bool
	var added []string
	for _, mutation := range mutations {
		if err := (&UpdateGroup{Name: mutation.Name}).IsValid(); err != nil {
			return err
		}
		if err := mutation.ExternalID.IsValid(); err != nil {
			return err
		}
		update = update || mutation.Name != nil || mutation.ExternalID != nil
		add = add || mutation.Members != nil || len(mutation.AddMembers) > 0
		remove = remove || mutation.Members != nil || len(mutation.RemoveMembers) > 0
		added = append(added, mutation.AddMembers...)
		if mutation.Members != nil {
			added = append(added, (*mutation.Members)...)
		}
	}
	if update {
		if err := c.checkPermissionUpdateGroup(ctx, orgID, groupID); err != nil {
			return err
		}
	}
	if add {
		if err := c.checkPermissionAddUserToGroup(ctx, orgID, groupID); err != nil {
			return err
		}
	}
	if len(added) > 0 {
		slices.Sort(added)
		if err := c.checkUsersExist(ctx, slices.Compact(added), orgID); err != nil {
			return err
		}
	}
	if remove {
		return c.checkPermissionRemoveUserFromGroup(ctx, orgID, groupID)
	}
	return nil
}
