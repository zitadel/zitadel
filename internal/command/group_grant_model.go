package command

import (
	"slices"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/groupgrant"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type GroupGrantWriteModel struct {
	eventstore.WriteModel

	GroupID        string
	ProjectID      string
	ProjectGrantID string
	RoleKeys       []string
	State          domain.GroupGrantState
}

func NewGroupGrantWriteModel(groupGrantID, resourceOwner string) *GroupGrantWriteModel {
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
			wm.ResourceOwner = e.Aggregate().ResourceOwner
		case *groupgrant.GroupGrantChangedEvent:
			wm.RoleKeys = e.RoleKeys
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
		EventTypes(
			groupgrant.GroupGrantAddedType,
			groupgrant.GroupGrantChangedType,
			groupgrant.GroupGrantRemovedType,
			groupgrant.GroupGrantCascadeRemovedType,
		).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func GroupGrantAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, groupgrant.AggregateType, groupgrant.AggregateVersion)
}

// GroupGrantPreConditionReadModel validates that the group exists in the organization
// and that the project (or project grant) provides the requested roles
type GroupGrantPreConditionReadModel struct {
	eventstore.WriteModel

	GroupID                 string
	ProjectID               string
	ProjectGrantID          string
	ResourceOwner           string
	ProjectResourceOwner    string
	FoundGrantID            string
	GroupExists             bool
	ProjectExists           bool
	ExistingRoleKeysProject []string
	ExistingRoleKeysGrant   []string
}

func NewGroupGrantPreConditionReadModel(groupID, projectID, projectGrantID, resourceOwner string) *GroupGrantPreConditionReadModel {
	return &GroupGrantPreConditionReadModel{
		GroupID:                 groupID,
		ProjectID:               projectID,
		ProjectGrantID:          projectGrantID,
		ResourceOwner:           resourceOwner,
		ExistingRoleKeysGrant:   make([]string, 0),
		ExistingRoleKeysProject: make([]string, 0),
	}
}

func (wm *GroupGrantPreConditionReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *group.GroupAddedEvent:
			if e.Aggregate().ResourceOwner == wm.ResourceOwner {
				wm.GroupExists = true
			}
		case *group.GroupRemovedEvent:
			wm.GroupExists = false
		case *project.ProjectAddedEvent:
			if projectExistsOnOrganization(wm.ResourceOwner, e.Aggregate().ResourceOwner) {
				wm.ProjectExists = true
			}
			wm.ProjectResourceOwner = e.Aggregate().ResourceOwner
		case *project.ProjectRemovedEvent:
			wm.ProjectExists = false
		case *project.GrantAddedEvent:
			if projectGrantExistsOnOrganization(wm.ProjectGrantID, wm.ResourceOwner, e.GrantID, e.GrantedOrgID) {
				wm.ExistingRoleKeysGrant = e.RoleKeys
				wm.FoundGrantID = e.GrantID
			}
		case *project.GrantChangedEvent:
			if wm.FoundGrantID == e.GrantID {
				wm.ExistingRoleKeysGrant = e.RoleKeys
			}
		case *project.GrantRemovedEvent:
			if wm.FoundGrantID == e.GrantID {
				wm.ExistingRoleKeysGrant = []string{}
				wm.FoundGrantID = ""
			}
		case *project.RoleAddedEvent:
			wm.ExistingRoleKeysProject = append(wm.ExistingRoleKeysProject, e.Key)
		case *project.RoleRemovedEvent:
			wm.ExistingRoleKeysProject = slices.DeleteFunc(wm.ExistingRoleKeysProject, func(key string) bool {
				return key == e.Key
			})
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GroupGrantPreConditionReadModel) existingRoles() []string {
	if wm.FoundGrantID != "" {
		return wm.ExistingRoleKeysGrant
	}
	return wm.ExistingRoleKeysProject
}

func (wm *GroupGrantPreConditionReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(group.AggregateType).
		AggregateIDs(wm.GroupID).
		EventTypes(
			group.GroupAddedEventType,
			group.GroupRemovedEventType,
		).
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
			project.RoleRemovedType,
		).
		Builder()
}
