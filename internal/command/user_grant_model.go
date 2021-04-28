package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/repository/usergrant"
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
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, usergrant.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(usergrant.UserGrantAddedType,
			usergrant.UserGrantChangedType,
			usergrant.UserGrantCascadeChangedType,
			usergrant.UserGrantDeactivatedType,
			usergrant.UserGrantReactivatedType,
			usergrant.UserGrantRemovedType,
			usergrant.UserGrantCascadeRemovedType)
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

	UserID             string
	ProjectID          string
	ProjectGrantID     string
	UserExists         bool
	ProjectExists      bool
	ProjectGrantExists bool
	ExistingRoleKeys   []string
}

func NewUserGrantPreConditionReadModel(userID, projectID, projectGrantID string) *UserGrantPreConditionReadModel {
	return &UserGrantPreConditionReadModel{
		UserID:         userID,
		ProjectID:      projectID,
		ProjectGrantID: projectGrantID,
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
			wm.ProjectExists = true
		case *project.ProjectRemovedEvent:
			wm.ProjectExists = false
		case *project.GrantAddedEvent:
			if wm.ProjectGrantID == e.GrantID {
				wm.ProjectGrantExists = true
			}
			wm.ExistingRoleKeys = e.RoleKeys
		case *project.GrantChangedEvent:
			wm.ExistingRoleKeys = e.RoleKeys
		case *project.GrantRemovedEvent:
			if wm.ProjectGrantID == e.GrantID {
				wm.ProjectGrantExists = false
			}
			wm.ExistingRoleKeys = []string{}
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

func (wm *UserGrantPreConditionReadModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType, project.AggregateType).
		AggregateIDs(wm.UserID, wm.ProjectID).
		EventTypes(user.UserV1AddedType,
			user.HumanAddedType,
			user.UserV1RegisteredType,
			user.HumanRegisteredType,
			user.MachineAddedEventType,
			user.UserRemovedType,
			project.ProjectAddedType,
			project.ProjectRemovedType,
			project.GrantAddedType,
			project.GrantChangedType,
			project.GrantRemovedType,
			project.RoleAddedType,
			project.RoleRemovedType)
	return query
}
