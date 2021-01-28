package domain

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type ProjectGrant struct {
	es_models.ObjectRoot

	GrantID      string
	GrantedOrgID string
	State        ProjectGrantState
	RoleKeys     []string
}

type ProjectGrantIDs struct {
	ProjectID string
	GrantID   string
}

type ProjectGrantState int32

const (
	ProjectGrantStateUnspecified ProjectGrantState = iota
	ProjectGrantStateActive
	ProjectGrantStateInactive
	ProjectGrantStateRemoved
)

func (p *ProjectGrant) IsValid() bool {
	return p.GrantedOrgID != ""
}
