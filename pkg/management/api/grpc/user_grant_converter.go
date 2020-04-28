package grpc

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/golang/protobuf/ptypes"
)

func usergrantFromModel(grant *grant_model.UserGrant) *UserGrant {
	creationDate, err := ptypes.TimestampProto(grant.CreationDate)
	logging.Log("GRPC-ki9ds").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(grant.ChangeDate)
	logging.Log("GRPC-sl9ew").OnError(err).Debug("unable to parse timestamp")

	converted := &UserGrant{
		Id:           grant.AggregateID,
		UserId:       grant.UserID,
		State:        usergrantStateFromModel(grant.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     grant.Sequence,
		ProjectId:    grant.ProjectID,
		RoleKeys:     grant.RoleKeys,
	}
	return converted
}

func userGrantCreateToModel(u *UserGrantCreate) *grant_model.UserGrant {
	grant := &grant_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.UserId},
		UserID:     u.UserId,
		ProjectID:  u.ProjectId,
		RoleKeys:   u.RoleKeys,
	}
	return grant
}

func userGrantUpdateToModel(u *UserGrantUpdate) *grant_model.UserGrant {
	grant := &grant_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.Id},
		RoleKeys:   u.RoleKeys,
	}
	return grant
}

func projectUserGrantUpdateToModel(u *ProjectUserGrantUpdate) *grant_model.UserGrant {
	grant := &grant_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.Id},
		RoleKeys:   u.RoleKeys,
	}
	return grant
}

func projectGrantUserGrantCreateToModel(u *ProjectGrantUserGrantCreate) *grant_model.UserGrant {
	grant := &grant_model.UserGrant{
		UserID:    u.UserId,
		ProjectID: u.ProjectId,
		RoleKeys:  u.RoleKeys,
	}
	return grant
}

func projectGrantUserGrantUpdateToModel(u *ProjectGrantUserGrantUpdate) *grant_model.UserGrant {
	grant := &grant_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.Id},
		RoleKeys:   u.RoleKeys,
	}
	return grant
}

func usergrantStateFromModel(state grant_model.UserGrantState) UserGrantState {
	switch state {
	case grant_model.USERGRANTSTATE_ACTIVE:
		return UserGrantState_USERGRANTSTATE_ACTIVE
	case grant_model.USERGRANTSTATE_INACTIVE:
		return UserGrantState_USERGRANTSTATE_INACTIVE
	default:
		return UserGrantState_USERGRANTSTATE_UNSPECIFIED
	}
}
