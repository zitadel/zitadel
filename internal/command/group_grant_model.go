package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/groupgrant"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type GroupGrantWriteModel struct {
	eventstore.WriteModel

	GroupID        string
	ProjectID      string
	ProjectGrantID string
	RoleKeys       []string
	State          domain.GroupGrantState
}

func NewGroupGrantWriteModel(groupGrantID string, resourceOwner string) *GroupGrantWriteModel {
	return &GroupGrantWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   groupGrantID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *GroupGrantWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *groupgrant.GroupGrantAddedEvent:
			wm.GroupID = e.GroupID
			wm.ProjectID = e.ProjectID
			wm.ProjectGrantID = e.ProjectGrantID
			wm.RoleKeys = e.RoleKeys
			wm.State = domain.GroupGrantStateActive
		case *groupgrant.GroupGrantChangedEvent:
			wm.RoleKeys = e.RoleKeys
		case *groupgrant.GroupGrantCascadeChangedEvent:
			wm.RoleKeys = e.RoleKeys
		case *groupgrant.GroupGrantDeactivatedEvent:
			if wm.State == domain.GroupGrantStateRemoved {
				continue
			}
			wm.State = domain.GroupGrantStateInactive
		case *groupgrant.GroupGrantReactivatedEvent:
			if wm.State == domain.GroupGrantStateRemoved {
				continue
			}
			wm.State = domain.GroupGrantStateActive
		case *groupgrant.GroupGrantRemovedEvent:
			wm.State = domain.GroupGrantStateRemoved
		case *groupgrant.GroupGrantCascadeRemovedEvent:
			wm.State = domain.GroupGrantStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GroupGrantWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(groupgrant.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(groupgrant.GroupGrantAddedType,
			groupgrant.GroupGrantChangedType,
			groupgrant.GroupGrantCascadeChangedType,
			groupgrant.GroupGrantDeactivatedType,
			groupgrant.GroupGrantReactivatedType,
			groupgrant.GroupGrantRemovedType,
			groupgrant.GroupGrantCascadeRemovedType).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func GroupGrantAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, groupgrant.AggregateType, groupgrant.AggregateVersion)
}

type GroupGrantPreConditionReadModel struct {
	eventstore.WriteModel

	GroupID            string
	ProjectID          string
	ProjectGrantID     string
	ResourceOwner      string
	GroupExists        bool
	ProjectExists      bool
	ProjectGrantExists bool
	ExistingRoleKeys   []string
}

func NewGroupGrantPreConditionReadModel(userID, projectID, projectGrantID, resourceOwner string) *GroupGrantPreConditionReadModel {
	return &GroupGrantPreConditionReadModel{
		GroupID:        userID,
		ProjectID:      projectID,
		ProjectGrantID: projectGrantID,
		ResourceOwner:  resourceOwner,
	}
}

func (wm *GroupGrantPreConditionReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *group.GroupAddedEvent:
			wm.GroupExists = true
		case *group.GroupRemovedEvent:
			wm.GroupExists = false
		case *project.ProjectAddedEvent:
			if wm.ProjectGrantID == "" && wm.ResourceOwner == e.Aggregate().ResourceOwner {
				wm.ProjectExists = true
			}
		case *project.ProjectRemovedEvent:
			wm.ProjectExists = false
		case *project.GrantAddedEvent:
			if wm.ProjectGrantID == e.GrantID && wm.ResourceOwner == e.GrantedOrgID {
				wm.ProjectGrantExists = true
				wm.ExistingRoleKeys = e.RoleKeys
			}
		case *project.GrantChangedEvent:
			if wm.ProjectGrantID == e.GrantID {
				wm.ExistingRoleKeys = e.RoleKeys
			}
		case *project.GrantRemovedEvent:
			if wm.ProjectGrantID == e.GrantID {
				wm.ProjectGrantExists = false
				wm.ExistingRoleKeys = []string{}
			}
		case *project.RoleAddedEvent:
			if wm.ProjectGrantID != "" {
				continue
			}
			wm.ExistingRoleKeys = append(wm.ExistingRoleKeys, e.Key)
		case *project.RoleRemovedEvent:
			if wm.ProjectGrantID != "" {
				continue
			}
			for i, key := range wm.ExistingRoleKeys {
				if key == e.Key {
					copy(wm.ExistingRoleKeys[i:], wm.ExistingRoleKeys[i+1:])
					wm.ExistingRoleKeys[len(wm.ExistingRoleKeys)-1] = ""
					wm.ExistingRoleKeys = wm.ExistingRoleKeys[:len(wm.ExistingRoleKeys)-1]
					continue
				}
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GroupGrantPreConditionReadModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.GroupID).
		EventTypes(
			group.GroupAddedType,
			group.GroupRemovedType).
		Or().
		AggregateTypes(project.AggregateType).
		AggregateIDs(wm.ProjectID).
		EventTypes(
			project.ProjectAddedType,
			project.ProjectRemovedType,
			project.GrantAddedType,
			project.GrantChangedType,
			project.GrantRemovedType,
			project.RoleAddedType,
			project.RoleRemovedType).
		Builder()
	return query
}
