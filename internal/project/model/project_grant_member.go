package model

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type ProjectGrantMember struct {
	es_models.ObjectRoot
	GrantID string
	UserID  string
	Roles   []string
}

func NewProjectGrantMember(projectID, grantID, userID string) *ProjectGrantMember {
	return &ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{ID: projectID}, GrantID: grantID, UserID: userID}
}

func (p *ProjectGrantMember) IsValid() bool {
	return p.ID != "" && p.UserID != "" && len(p.Roles) != 0
}
