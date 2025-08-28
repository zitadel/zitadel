package command

import (
	"slices"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
)

type UserGrantWriteModel struct {
	eventstore.WriteModel

	UserID         string
	ProjectID      string
	ProjectGrantID string
	RoleKeys       []string
	State          domain.UserGrantState
}

func NewUserGrantWriteModel(userGrantID string, resourceOwner string) *UserGrantWriteModel {
	return &UserGrantWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userGrantID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *UserGrantWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *usergrant.UserGrantAddedEvent:
			wm.UserID = e.UserID
			wm.ProjectID = e.ProjectID
			wm.ProjectGrantID = e.ProjectGrantID
			wm.RoleKeys = e.RoleKeys
			wm.State = domain.UserGrantStateActive
			wm.ResourceOwner = e.Aggregate().ResourceOwner
		case *usergrant.UserGrantChangedEvent:
			wm.RoleKeys = e.RoleKeys
		case *usergrant.UserGrantCascadeChangedEvent:
			wm.RoleKeys = e.RoleKeys
		case *usergrant.UserGrantDeactivatedEvent:
			if wm.State == domain.UserGrantStateRemoved {
				continue
			}
			wm.State = domain.UserGrantStateInactive
		case *usergrant.UserGrantReactivatedEvent:
			if wm.State == domain.UserGrantStateRemoved {
				continue
			}
			wm.State = domain.UserGrantStateActive
		case *usergrant.UserGrantRemovedEvent:
			wm.State = domain.UserGrantStateRemoved
		case *usergrant.UserGrantCascadeRemovedEvent:
			wm.State = domain.UserGrantStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UserGrantWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(usergrant.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(usergrant.UserGrantAddedType,
			usergrant.UserGrantChangedType,
			usergrant.UserGrantCascadeChangedType,
			usergrant.UserGrantDeactivatedType,
			usergrant.UserGrantReactivatedType,
			usergrant.UserGrantRemovedType,
			usergrant.UserGrantCascadeRemovedType).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func UserGrantAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, usergrant.AggregateType, usergrant.AggregateVersion)
}

type UserGrantPreConditionReadModel struct {
	eventstore.WriteModel

	UserID                  string
	ProjectID               string
	ProjectResourceOwner    string
	ProjectGrantID          string
	FoundGrantID            string
	ResourceOwner           string
	UserExists              bool
	ProjectExists           bool
	ExistingRoleKeysProject []string
	ExistingRoleKeysGrant   []string
}

func NewUserGrantPreConditionReadModel(userID, projectID, projectGrantID string, resourceOwner string) *UserGrantPreConditionReadModel {
	return &UserGrantPreConditionReadModel{
		UserID:    userID,
		ProjectID: projectID,
		// ProjectGrantID can be empty, if grantedOrgID is in the resourceowner
		ProjectGrantID: projectGrantID,
		// resourceowner is either empty to use the project organization
		// or filled with the project organization
		// or filled with the organization the project is granted to
		ResourceOwner:           resourceOwner,
		ExistingRoleKeysGrant:   make([]string, 0),
		ExistingRoleKeysProject: make([]string, 0),
	}
}

func (wm *UserGrantPreConditionReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.UserExists = true
		case *user.HumanRegisteredEvent:
			wm.UserExists = true
		case *user.MachineAddedEvent:
			wm.UserExists = true
		case *user.UserRemovedEvent:
			wm.UserExists = false
		case *project.ProjectAddedEvent:
			if projectExistsOnOrganization(wm.ResourceOwner, e.Aggregate().ResourceOwner) {
				wm.ProjectExists = true
			}
			// We store the organization of the project for later checks, e.g. in case of a project grant
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
			slices.DeleteFunc(wm.ExistingRoleKeysProject, func(key string) bool {
				return key == e.Key
			})
		}
	}
	return wm.WriteModel.Reduce()
}

func projectExistsOnOrganization(requiredOrganization, projectResourceOwner string) bool {
	// Depending on the API, a request can either require a project to be part of a specific organization
	// or not. In the former case, the project must belong to the required organization.
	// In the latter case, it is sufficient that the project exists at all, since the user will be granted
	// automatically in the organization the project belongs to.
	return requiredOrganization == "" || requiredOrganization == projectResourceOwner
}

func projectGrantExistsOnOrganization(requiredGrantID, requiredOrganization, projectGrantID, grantedOrganization string) bool {
	// Depending on the API, a request can either require a project grant (ID) and/or an organization (ID),
	// where the project must be granted to.
	return (requiredGrantID == "" || requiredGrantID == projectGrantID) &&
		(requiredOrganization == "" || requiredOrganization == grantedOrganization)
}

func (wm *UserGrantPreConditionReadModel) existingRoles() []string {
	if wm.FoundGrantID != "" {
		return wm.ExistingRoleKeysGrant
	}
	return wm.ExistingRoleKeysProject
}

func (wm *UserGrantPreConditionReadModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.UserID).
		EventTypes(
			user.UserV1AddedType,
			user.HumanAddedType,
			user.UserV1RegisteredType,
			user.HumanRegisteredType,
			user.MachineAddedEventType,
			user.UserRemovedType).
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
