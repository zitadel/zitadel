package grpc

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/golang/protobuf/ptypes"
)

func projectGrantFromModel(grant *proj_model.ProjectGrant) *ProjectGrant {
	creationDate, err := ptypes.TimestampProto(grant.CreationDate)
	logging.Log("GRPC-8d73s").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(grant.ChangeDate)
	logging.Log("GRPC-dlso3").OnError(err).Debug("unable to parse timestamp")

	return &ProjectGrant{
		Id:           grant.GrantID,
		State:        projectGrantStateFromModel(grant.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		GrantedOrgId: grant.GrantedOrgID,
		RoleKeys:     grant.RoleKeys,
		Sequence:     grant.Sequence,
	}
}

func projectGrantCreateToModel(grant *ProjectGrantCreate) *proj_model.ProjectGrant {
	return &proj_model.ProjectGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: grant.ProjectId,
		},
		GrantedOrgID: grant.GrantedOrgId,
		RoleKeys:     grant.RoleKeys,
	}
}

func projectGrantUpdateToModel(grant *ProjectGrantUpdate) *proj_model.ProjectGrant {
	return &proj_model.ProjectGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: grant.ProjectId,
		},
		GrantID:  grant.Id,
		RoleKeys: grant.RoleKeys,
	}
}

func projectGrantStateFromModel(state proj_model.ProjectGrantState) ProjectGrantState {
	switch state {
	case proj_model.PROJECTGRANTSTATE_ACTIVE:
		return ProjectGrantState_PROJECTGRANTSTATE_ACTIVE
	case proj_model.PROJECTGRANTSTATE_INACTIVE:
		return ProjectGrantState_PROJECTGRANTSTATE_INACTIVE
	default:
		return ProjectGrantState_PROJECTGRANTSTATE_UNSPECIFIED
	}
}
