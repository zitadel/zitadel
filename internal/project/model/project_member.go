package model

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type ProjectMember struct {
	es_models.ObjectRoot

	UserID string
	Roles  []string
}

func NewProjectMember(projectID, userID string) *ProjectMember {
	return &ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: projectID}, UserID: userID}
}

func (p *ProjectMember) IsValid() bool {
	if p.ID == "" || p.UserID == "" || len(p.Roles) == 0 {
		return false
	}
	return true
}
