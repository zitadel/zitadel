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
	GrantID                 string
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
		// resourceowner has to be either empty, resourceowner of the project or grantedOrdID of the projectgrant
		ResourceOwner: resourceOwner,
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
			if wm.ResourceOwner == "" || wm.ResourceOwner == e.Aggregate().ResourceOwner {
				wm.ProjectExists = true
			}
			wm.ProjectResourceOwner = e.Aggregate().ResourceOwner
		case *project.ProjectRemovedEvent:
			// only necessary if project existing
			if wm.ProjectResourceOwner == e.Aggregate().ResourceOwner {
				wm.ProjectExists = false
			}
		case *project.GrantAddedEvent:
			// grantID can be empty or provided and equal
			// and resourceowner has to be provided and equal to an organization the project is granted to, and not own resourceowner
			if (wm.ProjectGrantID == "" || wm.ProjectGrantID == e.GrantID) &&
				wm.ResourceOwner != "" && wm.ResourceOwner == e.GrantedOrgID && e.GrantedOrgID != e.Aggregate().ResourceOwner {
				wm.ProjectGrantExists = true
				wm.ExistingRoleKeysGrant = e.RoleKeys
				// persist grantID without changing the original
				wm.GrantID = e.GrantID
			}
		case *project.GrantChangedEvent:
			// grantID is filled by grant added event, and guarantees that project exists
			// only applicable if change for the same projectgrant
			if wm.GrantID != "" && wm.GrantID == e.GrantID {
				wm.ExistingRoleKeysGrant = e.RoleKeys
			}
		case *project.GrantRemovedEvent:
			// grantID is filled by grant added event, and guarantees that project exists
			// only applicable if change for the same projectgrant
			if wm.GrantID != "" && wm.GrantID == e.GrantID {
				wm.ProjectGrantExists = false
				wm.ExistingRoleKeysGrant = []string{}
				// remove grantID as not longer necessary
				wm.GrantID = ""
			}
		case *project.RoleAddedEvent:
			// ignore if project grant is requested or resourceowner is not equal to the resourceowner of the project
			if wm.ProjectGrantID != "" || (wm.ResourceOwner != "" && wm.ResourceOwner != e.Aggregate().ResourceOwner) {
				continue
			}
			wm.ExistingRoleKeysProject = append(wm.ExistingRoleKeysProject, e.Key)
		case *project.RoleRemovedEvent:
			// ignore if project grant is requested or resourceowner is not equal to the resourceowner of the project
			if wm.ProjectGrantID != "" || (wm.ResourceOwner != "" && wm.ResourceOwner != e.Aggregate().ResourceOwner) {
				continue
			}
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
	return wm.WriteModel.Reduce()
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
