package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type ProjectWriteModel struct {
	eventstore.WriteModel

	Name                   string
	ProjectRoleAssertion   bool
	ProjectRoleCheck       bool
	HasProjectCheck        bool
	PrivateLabelingSetting domain.PrivateLabelingSetting
	State                  domain.ProjectState
}

func NewProjectWriteModel(projectID string, resourceOwner string) *ProjectWriteModel {
	return &ProjectWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *ProjectWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *project.ProjectAddedEvent:
			wm.Name = e.Name
			wm.ProjectRoleAssertion = e.ProjectRoleAssertion
			wm.ProjectRoleCheck = e.ProjectRoleCheck
			wm.HasProjectCheck = e.HasProjectCheck
			wm.PrivateLabelingSetting = e.PrivateLabelingSetting
			wm.State = domain.ProjectStateActive
		case *project.ProjectChangeEvent:
			if e.Name != nil {
				wm.Name = *e.Name
			}
			if e.ProjectRoleAssertion != nil {
				wm.ProjectRoleAssertion = *e.ProjectRoleAssertion
			}
			if e.ProjectRoleCheck != nil {
				wm.ProjectRoleCheck = *e.ProjectRoleCheck
			}
			if e.HasProjectCheck != nil {
				wm.HasProjectCheck = *e.HasProjectCheck
			}
			if e.PrivateLabelingSetting != nil {
				wm.PrivateLabelingSetting = *e.PrivateLabelingSetting
			}
		case *project.ProjectDeactivatedEvent:
			if wm.State == domain.ProjectStateRemoved {
				continue
			}
			wm.State = domain.ProjectStateInactive
		case *project.ProjectReactivatedEvent:
			if wm.State == domain.ProjectStateRemoved {
				continue
			}
			wm.State = domain.ProjectStateActive
		case *project.ProjectRemovedEvent:
			wm.State = domain.ProjectStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *ProjectWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(project.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(project.ProjectAddedType,
			project.ProjectChangedType,
			project.ProjectDeactivatedType,
			project.ProjectReactivatedType,
			project.ProjectRemovedType).
		Builder()
}

func (wm *ProjectWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name *string,
	projectRoleAssertion,
	projectRoleCheck,
	hasProjectCheck *bool,
	privateLabelingSetting *domain.PrivateLabelingSetting,
) *project.ProjectChangeEvent {
	changes := make([]project.ProjectChanges, 0)

	oldName := ""
	if name != nil && wm.Name != *name {
		oldName = wm.Name
		changes = append(changes, project.ChangeName(*name))
	}
	if projectRoleAssertion != nil && wm.ProjectRoleAssertion != *projectRoleAssertion {
		changes = append(changes, project.ChangeProjectRoleAssertion(*projectRoleAssertion))
	}
	if projectRoleCheck != nil && wm.ProjectRoleCheck != *projectRoleCheck {
		changes = append(changes, project.ChangeProjectRoleCheck(*projectRoleCheck))
	}
	if hasProjectCheck != nil && wm.HasProjectCheck != *hasProjectCheck {
		changes = append(changes, project.ChangeHasProjectCheck(*hasProjectCheck))
	}
	if privateLabelingSetting != nil && wm.PrivateLabelingSetting != *privateLabelingSetting {
		changes = append(changes, project.ChangePrivateLabelingSetting(*privateLabelingSetting))
	}
	if len(changes) == 0 {
		return nil
	}
	return project.NewProjectChangeEvent(ctx, aggregate, oldName, changes)
}

func isProjectStateExists(state domain.ProjectState) bool {
	return !slices.Contains([]domain.ProjectState{domain.ProjectStateRemoved, domain.ProjectStateUnspecified}, state)
}

func ProjectAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, project.AggregateType, project.AggregateVersion)
}

func ProjectAggregateFromWriteModelWithCTX(ctx context.Context, wm *eventstore.WriteModel) *eventstore.Aggregate {
	return project.AggregateFromWriteModel(ctx, wm)
}
