package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type ProjectGroupMemberWrite struct {
	eventstore.WriteModel

	GroupID string
	Roles   []string

	State domain.GroupMemberState
}

type ProjectGroupMemberWriteModel struct {
	ProjectGroupMemberWrite
}

// AppendEvents implements eventstore.QueryReducer.
// Subtle: this method shadows the method (ProjectGroupMemberWrite).AppendEvents of ProjectGroupMemberWriteModel.ProjectGroupMemberWrite.
func (wm *ProjectGroupMemberWriteModel) AppendEvents(...eventstore.Event) {
	panic("unimplemented")
}

// Query implements eventstore.QueryReducer.
func (wm *ProjectGroupMemberWriteModel) Query() *eventstore.SearchQueryBuilder {
	panic("unimplemented")
}

// Reduce implements eventstore.QueryReducer.
// Subtle: this method shadows the method (ProjectGroupMemberWrite).Reduce of ProjectGroupMemberWriteModel.ProjectGroupMemberWrite.
func (wm *ProjectGroupMemberWriteModel) Reduce() error {
	panic("unimplemented")
}

func NewProjectGroupMemberWriteModel(projectID, groupID, resourceOwner string) *ProjectGroupMemberWriteModel {
	return &ProjectGroupMemberWriteModel{
		ProjectGroupMemberWrite{
			WriteModel: eventstore.WriteModel{
				AggregateID:   projectID,
				ResourceOwner: resourceOwner,
			},
			GroupID: groupID,
		},
	}
}

func (wm *ProjectGroupMemberWriteModel) AppendGroupEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *project.GroupMemberAddedEvent:
			if e.GroupID != wm.ProjectGroupMemberWrite.GroupID {
				continue
			}
			wm.ProjectGroupMemberWrite.AppendEvents(&e.GroupMemberAddedEvent)
		case *project.GroupMemberChangedEvent:
			if e.GroupID != wm.ProjectGroupMemberWrite.GroupID {
				continue
			}
			wm.ProjectGroupMemberWrite.AppendEvents(&e.GroupMemberChangedEvent)
		case *project.GroupMemberRemovedEvent:
			if e.GroupID != wm.ProjectGroupMemberWrite.GroupID {
				continue
			}
			wm.ProjectGroupMemberWrite.AppendEvents(&e.GroupMemberRemovedEvent)
		case *project.GroupMemberCascadeRemovedEvent:
			if e.GroupID != wm.ProjectGroupMemberWrite.GroupID {
				continue
			}
			wm.ProjectGroupMemberWrite.AppendEvents(&e.GroupMemberCascadeRemovedEvent)
		}
	}
}

func (wm *ProjectGroupMemberWriteModel) ReduceGroup() error {
	return wm.ProjectGroupMemberWrite.Reduce()
}

func (wm *ProjectGroupMemberWriteModel) QueryGroup() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(project.AggregateType).
		AggregateIDs(wm.ProjectGroupMemberWrite.AggregateID).
		EventTypes(project.GroupMemberAddedType,
			project.GroupMemberChangedType,
			project.GroupMemberRemovedType,
			project.GroupMemberCascadeRemovedType).
		Builder()
}
