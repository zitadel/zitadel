package domain

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type UserGrant struct {
	es_models.ObjectRoot

	State          UserGrantState
	UserID         string
	ProjectID      string
	ProjectGrantID string
	RoleKeys       []string
}

type UserGrantState int32

const (
	UserGrantStateUnspecified UserGrantState = iota
	UserGrantStateActive
	UserGrantStateInactive
	UserGrantStateRemoved
)

func (u *UserGrant) IsValid() bool {
	return u.ProjectID != "" && u.UserID != ""
}
