package grpc

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/golang/protobuf/ptypes"
)

func projectGrantMemberFromModel(member *proj_model.ProjectGrantMember) *ProjectGrantMember {
	creationDate, err := ptypes.TimestampProto(member.CreationDate)
	logging.Log("GRPC-7du3s").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(member.ChangeDate)
	logging.Log("GRPC-8duew").OnError(err).Debug("unable to parse timestamp")

	return &ProjectGrantMember{
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     member.Sequence,
		UserId:       member.UserID,
		Roles:        member.Roles,
	}
}

func projectGrantMemberAddToModel(member *ProjectGrantMemberAdd) *proj_model.ProjectGrantMember {
	return &proj_model.ProjectGrantMember{
		ObjectRoot: models.ObjectRoot{
			AggregateID: member.ProjectId,
		},
		GrantID: member.GrantId,
		UserID:  member.UserId,
		Roles:   member.Roles,
	}
}

func projectGrantMemberChangeToModel(member *ProjectGrantMemberChange) *proj_model.ProjectGrantMember {
	return &proj_model.ProjectGrantMember{
		ObjectRoot: models.ObjectRoot{
			AggregateID: member.ProjectId,
		},
		GrantID: member.GrantId,
		UserID:  member.UserId,
		Roles:   member.Roles,
	}
}
