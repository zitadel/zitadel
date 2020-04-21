package grpc

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/golang/protobuf/ptypes"
)

func projectMemberFromModel(member *proj_model.ProjectMember) *ProjectMember {
	creationDate, err := ptypes.TimestampProto(member.CreationDate)
	logging.Log("GRPC-kd8re").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(member.ChangeDate)
	logging.Log("GRPC-dlei3").OnError(err).Debug("unable to parse timestamp")

	return &ProjectMember{
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     member.Sequence,
		UserId:       member.UserID,
		Roles:        member.Roles,
	}
}

func projectMemberAddToModel(member *ProjectMemberAdd) *proj_model.ProjectMember {
	return &proj_model.ProjectMember{
		ObjectRoot: models.ObjectRoot{
			AggregateID: member.Id,
		},
		UserID: member.UserId,
		Roles:  member.Roles,
	}
}

func projectMemberChangeToModel(member *ProjectMemberChange) *proj_model.ProjectMember {
	return &proj_model.ProjectMember{
		ObjectRoot: models.ObjectRoot{
			AggregateID: member.Id,
		},
		UserID: member.UserId,
		Roles:  member.Roles,
	}
}
