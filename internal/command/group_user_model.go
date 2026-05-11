package command

import (
	"slices"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
)

// GroupUsersWriteModel is the write-model for the user-membership of a group.
type GroupUsersWriteModel struct {
	eventstore.WriteModel

	State domain.GroupState

	existingUsers map[string]map[string]group.AttributeValue
}

func NewGroupUsersWriteModel(groupID, orgID string) *GroupUsersWriteModel {
	return &GroupUsersWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   groupID,
			ResourceOwner: orgID,
		},
		existingUsers: make(map[string]map[string]group.AttributeValue),
	}
}

func (g *GroupUsersWriteModel) GetWriteModel() *eventstore.WriteModel {
	return &g.WriteModel
}

func (g *GroupUsersWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(g.ResourceOwner).
		AddQuery().
		AggregateTypes(group.AggregateType).
		AggregateIDs(g.AggregateID).
		EventTypes(
			group.GroupAddedEventType,
			group.GroupRemovedEventType,
			group.GroupUsersAddedEventType,
			group.GroupUsersChangedEventType,
			group.GroupUsersRemovedEventType,
		).Builder()
}

func (g *GroupUsersWriteModel) Reduce() error {
	for _, event := range g.Events {
		switch e := event.(type) {
		case *group.GroupAddedEvent:
			g.State = domain.GroupStateActive
		case *group.GroupRemovedEvent:
			g.State = domain.GroupStateRemoved
			g.existingUsers = make(map[string]map[string]group.AttributeValue)
		case *group.GroupUsersAddedEvent:
			for _, u := range e.Users {
				g.existingUsers[u.UserID] = u.Attributes
			}
		case *group.GroupUserChangedEvent:
			if _, ok := g.existingUsers[e.UserID]; ok {
				g.existingUsers[e.UserID] = e.Attributes
			}
		case *group.GroupUsersRemovedEvent:
			for _, userID := range e.UserIDs {
				delete(g.existingUsers, userID)
			}
		}
	}
	return g.WriteModel.Reduce()
}

// UsersToAdd filters requested users down to those not yet in the group,
// de-duplicating repeated userIDs within the request.
func (g *GroupUsersWriteModel) UsersToAdd(requested []group.GroupUser) []group.GroupUser {
	out := make([]group.GroupUser, 0, len(requested))
	seen := make(map[string]struct{}, len(requested))
	for _, u := range requested {
		if _, exists := g.existingUsers[u.UserID]; exists {
			continue
		}
		if _, dup := seen[u.UserID]; dup {
			continue
		}
		seen[u.UserID] = struct{}{}
		out = append(out, u)
	}
	return out
}

// UserIDsToRemove returns the intersection of requested IDs and current members.
// IDs that are not members are silently dropped (desired state already achieved).
func (g *GroupUsersWriteModel) UserIDsToRemove(requested []string) []string {
	out := make([]string, 0, len(requested))
	for _, userID := range requested {
		if _, exists := g.existingUsers[userID]; !exists {
			continue
		}
		if slices.Contains(out, userID) {
			continue
		}
		out = append(out, userID)
	}
	return out
}
