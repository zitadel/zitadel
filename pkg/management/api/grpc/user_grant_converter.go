package grpc

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/golang/protobuf/ptypes"
)

func usergrantFromModel(grant *usr_model.UserGrant) *UserGrant {
	creationDate, err := ptypes.TimestampProto(grant.CreationDate)
	logging.Log("GRPC-ki9ds").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(grant.ChangeDate)
	logging.Log("GRPC-sl9ew").OnError(err).Debug("unable to parse timestamp")

	converted := &UserGrant{
		Id:           grant.GrantID,
		UserId:       grant.AggregateID,
		State:        usergrantStateFromModel(grant.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     grant.Sequence,
		ProjectId:    grant.ProjectID,
		RoleKeys:     grant.RoleKeys,
	}
	return converted
}

func userGrantCreateToModel(u *UserGrantCreate) *usr_model.UserGrant {
	grant := &usr_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.UserId},
		ProjectID:  u.ProjectId,
		RoleKeys:   u.RoleKeys,
	}
	return grant
}

func userGrantUpdateToModel(u *UserGrantUpdate) *usr_model.UserGrant {
	grant := &usr_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.UserId},
		GrantID:    u.Id,
		RoleKeys:   u.RoleKeys,
	}
	return grant
}

func projectUserGrantUpdateToModel(u *ProjectUserGrantUpdate) *usr_model.UserGrant {
	grant := &usr_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.UserId},
		GrantID:    u.Id,
		RoleKeys:   u.RoleKeys,
	}
	return grant
}

func projectGrantUserGrantCreateToModel(u *ProjectGrantUserGrantCreate) *usr_model.UserGrant {
	grant := &usr_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.UserId},
		ProjectID:  u.ProjectId,
		RoleKeys:   u.RoleKeys,
	}
	return grant
}

func projectGrantUserGrantUpdateToModel(u *ProjectGrantUserGrantUpdate) *usr_model.UserGrant {
	grant := &usr_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.UserId},
		GrantID:    u.Id,
		RoleKeys:   u.RoleKeys,
	}
	return grant
}

func usergrantStateFromModel(state usr_model.UserGrantState) UserGrantState {
	switch state {
	case usr_model.USERGRANTSTATE_ACTIVE:
		return UserGrantState_USERGRANTSTATE_ACTIVE
	case usr_model.USERGRANTSTATE_INACTIVE:
		return UserGrantState_USERGRANTSTATE_INACTIVE
	default:
		return UserGrantState_USERGRANTSTATE_UNSPECIFIED
	}
}
