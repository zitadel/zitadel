package command

import (
	"context"
	"maps"
	"slices"
	"strings"

	"github.com/zitadel/zitadel/internal/eventstore"
	repo "github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func groupRemoval(ctx context.Context, group *GroupWriteModel) *repo.GroupRemovedEvent {
	return repo.NewGroupRemovedEvent(ctx, GroupAggregateFromWriteModel(ctx, &group.WriteModel), group.Name, group.ExternalIDs).
		WithMembers(slices.Sorted(maps.Keys(group.existingUserIDs)))
}

// Lifecycle fields commit with events; cascades under the instance lock need no history replay.
func orgGroupRemovals(ctx context.Context, search func(context.Context, ...map[eventstore.FieldType]any) ([]*eventstore.SearchResult, error), orgID string) ([]eventstore.Command, error) {
	fields, err := search(ctx, map[eventstore.FieldType]any{
		eventstore.FieldTypeResourceOwner: orgID,
		eventstore.FieldTypeAggregateType: repo.AggregateType,
		eventstore.FieldTypeObjectType:    repo.FieldObjectType,
	})
	if err != nil {
		return nil, err
	}
	type groupState struct {
		name          string
		members, keys []string
	}
	groups := make(map[string]*groupState)
	for _, field := range fields {
		id := field.Aggregate.ID
		if groups[id] == nil {
			groups[id] = new(groupState)
		}
		state := groups[id]
		switch field.FieldName {
		case repo.FieldName:
			if err := field.Value.Unmarshal(&state.name); err != nil {
				return nil, err
			}
		case repo.FieldMember:
			state.members = append(state.members, field.Object.ID)
		case repo.FieldExternalIDKey:
			state.keys = append(state.keys, field.Object.ID)
		}
	}
	var events []eventstore.Command
	for _, id := range slices.Sorted(maps.Keys(groups)) {
		state := groups[id]
		if state.name == "" {
			return nil, zerrors.ThrowInternal(nil, "GROUP-lifecycle-fields", "Missing group lifecycle state")
		}
		slices.Sort(state.members)
		slices.Sort(state.keys)
		aggregate := eventstore.NewAggregate(ctx, id, repo.AggregateType, repo.AggregateVersion, eventstore.WithResourceOwner(orgID))
		events = append(events, repo.NewGroupRemovedEvent(ctx, aggregate, state.name).WithMembers(state.members).WithExternalIDKeys(state.keys))
	}
	return events, nil
}

func (c *Commands) removeUserGroupMemberships(ctx context.Context, userID, orgID string) ([]eventstore.Command, error) {
	fields, err := c.eventstore.Search(ctx, map[eventstore.FieldType]any{
		eventstore.FieldTypeResourceOwner: orgID,
		eventstore.FieldTypeAggregateType: repo.AggregateType,
		eventstore.FieldTypeObjectType:    repo.FieldObjectType,
		eventstore.FieldTypeFieldName:     repo.FieldMember,
		eventstore.FieldTypeValue:         userID,
	})
	if err != nil {
		return nil, err
	}
	slices.SortFunc(fields, func(a, b *eventstore.SearchResult) int { return strings.Compare(a.Aggregate.ID, b.Aggregate.ID) })
	events := make([]eventstore.Command, 0, len(fields))
	for _, field := range fields {
		aggregate := eventstore.NewAggregate(ctx, field.Aggregate.ID, repo.AggregateType, repo.AggregateVersion, eventstore.WithResourceOwner(orgID))
		events = append(events, repo.NewGroupUserRemovedEvent(ctx, aggregate, userID))
	}
	return events, nil
}
