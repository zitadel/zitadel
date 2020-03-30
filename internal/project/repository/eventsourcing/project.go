package eventsourcing

import (
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/pkg"
	"github.com/caos/zitadel/internal/project/model"
	"github.com/sony/sonyflake"
	"strconv"
)

var idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})

type Project struct {
	pkg.ObjectRoot
	Name  string `json:"name,omitempty"`
	State int32  `json:"-"`
}

func ProjectFromModel(project *model.Project) *Project {
	return &Project{
		Name:  project.Name,
		State: model.ProjectStateToInt(project.State),
		ObjectRoot: pkg.ObjectRoot{
			ID:           project.ID,
			Sequence:     project.Sequence,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
		},
	}
}

func ProjectToModel(project *Project) *model.Project {
	return &model.Project{
		ObjectRoot: pkg.ObjectRoot{
			ID:           project.ID,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
			Sequence:     project.Sequence,
		},
		Name:  project.Name,
		State: model.ProjectStateFromInt(project.State),
	}
}

func ProjectByIDFilter(id string, latestSequence uint64) *models.FilterEventsRequest {
	return &models.FilterEventsRequest{
		LatestSequence: latestSequence,
		AggregateType:  model.ProjectAggregate,
		AggregateID:    id,
	}
}

func ProjectFilter(latestSequence uint64) *models.FilterEventsRequest {
	return &models.FilterEventsRequest{
		LatestSequence: latestSequence,
		AggregateType:  model.ProjectAggregate,
	}
}

func ProjectCreateEvents(project *Project) (projectAggregate *SaveAggregate) {
	var err error
	id, err := idGenerator.NextID()
	if err != nil {
		return nil
	}
	project.ID = strconv.FormatUint(id, 10)

	return createdProject(project)
}

func ProjectUpdateEvents(project *model.ProjectChange) *pkg.SaveAggregate {
	return updatedProject(project)
}

func createdProject(p *Project) *pkg.SaveAggregate {
	return pkg.NewSaveAggregate(p.ID, model.ProjectAggregate, 0,
		pkg.SaveEvent{Type: pkg.EventType(model.AddedProject), Payload: p},
	)
}

func updatedProject(p *model.ProjectChange) *pkg.SaveAggregate {
	return pkg.NewSaveAggregate(p.ID, model.ProjectAggregate, p.Sequence,
		pkg.SaveEvent{Type: model.ChangedProject, Payload: p.Payload},
	)
}

func ProjectDeactivateEvents(id string, sequence uint64) *pkg.SaveAggregate {
	return pkg.NewSaveAggregate(id, model.ProjectAggregate, sequence,
		pkg.SaveEvent{model.DeactivatedProject, nil},
	)
}

func ProjectReactivateEvents(id string, sequence uint64) *pkg.SaveAggregate {
	return pkg.NewSaveAggregate(id, model.ProjectAggregate, sequence,
		pkg.SaveEvent{model.ReactivatedProject, nil},
	)
}
