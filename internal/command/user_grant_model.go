package command

import (
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
	ProjectGrantExists      bool
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
			if checkIfProjectNecessary(wm.ResourceOwner, e.Aggregate().ResourceOwner) {
				wm.ProjectExists = true
			}
			wm.ProjectResourceOwner = e.Aggregate().ResourceOwner
		case *project.ProjectRemovedEvent:
			if checkIfProjectNecessary(wm.ResourceOwner, e.Aggregate().ResourceOwner) {
				wm.ProjectExists = false
			}
		case *project.GrantAddedEvent:
			// grantID is empty to search for a granted project or provided AND has to be equal to the grantID of the granted project
			// AND there has to be an organization the project is granted to AND it can not be the organization belong to itself
			if checkIfProjectGrantNecessary(wm.ResourceOwner, e.Aggregate().ResourceOwner, e.GrantedOrgID, wm.ProjectGrantID, e.GrantID) {
				wm.ProjectGrantExists = true
				wm.ExistingRoleKeysGrant = e.RoleKeys
				// persist grantID without changing the original
				wm.FoundGrantID = e.GrantID
			}
		case *project.GrantChangedEvent:
			// grantedOrgID is empty as grantID is always set
			if checkIfProjectGrantNecessary(wm.ResourceOwner, e.Aggregate().ResourceOwner, "", wm.FoundGrantID, e.GrantID) {
				wm.ExistingRoleKeysGrant = e.RoleKeys
			}
		case *project.GrantRemovedEvent:
			// grantedOrgID is empty as grantID is always set
			if checkIfProjectGrantNecessary(wm.ResourceOwner, e.Aggregate().ResourceOwner, "", wm.FoundGrantID, e.GrantID) {
				wm.ProjectGrantExists = false
				wm.ExistingRoleKeysGrant = []string{}
			}
		case *project.RoleAddedEvent:
			if checkIfProjectRoleNecessary(wm.ResourceOwner, e.Aggregate().ResourceOwner) {
				wm.ExistingRoleKeysProject = append(wm.ExistingRoleKeysProject, e.Key)
			}
		case *project.RoleRemovedEvent:
			if checkIfProjectRoleNecessary(wm.ResourceOwner, e.Aggregate().ResourceOwner) {
				for i, key := range wm.ExistingRoleKeysProject {
					if key == e.Key {
						copy(wm.ExistingRoleKeysProject[i:], wm.ExistingRoleKeysProject[i+1:])
						wm.ExistingRoleKeysProject[len(wm.ExistingRoleKeysProject)-1] = ""
						wm.ExistingRoleKeysProject = wm.ExistingRoleKeysProject[:len(wm.ExistingRoleKeysProject)-1]
						continue
					}
				}
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func checkIfProjectNecessary(providedResourceOwner, projectResourceOwner string) bool {
	// organization of the project is used OR provided is organization of the project
	if providedResourceOwner == "" || providedResourceOwner == projectResourceOwner {
		return true
	}
	return false
}

func checkIfProjectGrantNecessary(providedResourceOwner, projectResourceOwner, grantedOrgID, providedProjectGrantID, projectGrantID string) bool {
	// grantID is empty to search for a granted project or provided AND has to be equal to the grantID of the granted project
	// AND there has to be an organization the project is granted to AND it can not be the organization belong to itself
	if (providedProjectGrantID == "" || providedProjectGrantID == projectGrantID) &&
		(providedResourceOwner != "" && providedResourceOwner == grantedOrgID && grantedOrgID != projectResourceOwner) {
		return true
	}
	return false
}

func checkIfProjectRoleNecessary(providedResourceOwner, projectResourceOwner string) bool {
	// organization of the project is used OR provided is organization of project
	if providedResourceOwner == "" || providedResourceOwner == projectResourceOwner {
		return true
	}
	return false
}

func (wm *UserGrantPreConditionReadModel) existingRoles(projectID, projectGrantID string) []string {
	// not the requested project or not found project grant
	if wm.ProjectID != projectID || wm.FoundGrantID != projectGrantID {
		return nil
	}
	// requested project
	if projectGrantID == "" {
		return wm.ExistingRoleKeysProject
	}
	return wm.ExistingRoleKeysGrant
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
