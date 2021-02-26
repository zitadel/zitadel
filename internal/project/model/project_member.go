package model

import es_models "github.com/caos/zitadel/internal/eventstore/v1/models"

type ProjectMember struct {
	es_models.ObjectRoot

	UserID string
	Roles  []string
}

func NewProjectMember(projectID, userID string) *ProjectMember {
	return &ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: projectID}, UserID: userID}
}

func (p *ProjectMember) IsValid() bool {
	return p.AggregateID != "" && p.UserID != "" && len(p.Roles) != 0
}
