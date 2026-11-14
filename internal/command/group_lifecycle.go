package command

import (
	"context"
	"maps"
	"slices"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/eventstore"
	repo "github.com/zitadel/zitadel/internal/repository/group"
)

// Called under the instance lifecycle lock. Projections cannot be used here:
// they may omit a membership whose add transaction has already committed.
func orgGroupWriteModels(ctx context.Context, filter preparation.FilterToQueryReducer, orgID string) (map[string]*GroupWriteModel, error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).ResourceOwner(orgID).OrderAsc().AddQuery().AggregateTypes(repo.AggregateType).Builder())
	if err != nil {
		return nil, err
	}
	groups := make(map[string]*GroupWriteModel)
	for _, event := range events {
		id := event.Aggregate().ID
		if groups[id] == nil {
			groups[id] = NewGroupWriteModel(id, orgID, []string{})
		}
		groups[id].AppendEvents(event)
	}
	for _, group := range groups {
		if err := group.Reduce(); err != nil {
			return nil, err
		}
	}
	return groups, nil
}

func groupRemoval(ctx context.Context, group *GroupWriteModel) *repo.GroupRemovedEvent {
	return repo.NewGroupRemovedEvent(ctx, GroupAggregateFromWriteModel(ctx, &group.WriteModel), group.Name, group.ExternalIDs).
		WithMembers(slices.Sorted(maps.Keys(group.existingUserIDs)))
}

func orgGroupRemovals(ctx context.Context, filter preparation.FilterToQueryReducer, orgID string) ([]eventstore.Command, error) {
	groups, err := orgGroupWriteModels(ctx, filter, orgID)
	if err != nil {
		return nil, err
	}
	var events []eventstore.Command
	for _, id := range slices.Sorted(maps.Keys(groups)) {
		if groups[id].State.Exists() {
			events = append(events, groupRemoval(ctx, groups[id]))
		}
	}
	return events, nil
}

func (c *Commands) removeUserGroupMemberships(ctx context.Context, userID, orgID string) ([]eventstore.Command, error) {
	//nolint:staticcheck // cascade discovery must replay authoritative state in this transaction.
	groups, err := orgGroupWriteModels(ctx, c.eventstore.Filter, orgID)
	if err != nil {
		return nil, err
	}
	var events []eventstore.Command
	for _, id := range slices.Sorted(maps.Keys(groups)) {
		group := groups[id]
		if _, member := group.existingUserIDs[userID]; group.State.Exists() && member {
			events = append(events, repo.NewGroupUserRemovedEvent(ctx, GroupAggregateFromWriteModel(ctx, &group.WriteModel), userID))
		}
	}
	return events, nil
}
