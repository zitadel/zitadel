package grpc

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/golang/protobuf/ptypes"
)

func projectFromModel(project *proj_model.Project) *Project {
	creationDate, err := ptypes.TimestampProto(project.CreationDate)
	logging.Log("GRPC-iejs3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(project.ChangeDate)
	logging.Log("GRPC-di7rw").OnError(err).Debug("unable to parse timestamp")

	return &Project{
		Id:           project.ID,
		State:        projectStateFromModel(project.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Name:         project.Name,
		Sequence:     project.Sequence,
	}
}

func projectStateFromModel(state proj_model.ProjectState) ProjectState {
	switch state {
	case proj_model.Active:
		return ProjectState_PROJECTSTATE_ACTIVE
	case proj_model.Inactive:
		return ProjectState_PROJECTSTATE_INACTIVE
	default:
		return ProjectState_PROJECTSTATE_UNSPECIFIED
	}
}

func projectUpdateToModel(project *ProjectUpdateRequest) *proj_model.Project {
	return &proj_model.Project{
		ObjectRoot: models.ObjectRoot{
			ID: project.Id,
		},
		Name: project.Name,
	}
}
