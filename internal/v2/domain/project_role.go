package domain

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type ProjectRole struct {
	models.ObjectRoot

	Key         string
	DisplayName string
	Group       string
}

func NewProjectRole(projectID, key string) *ProjectRole {
	return &ProjectRole{ObjectRoot: models.ObjectRoot{AggregateID: projectID}, Key: key}
}

func (p *ProjectRole) IsValid() bool {
	return p.AggregateID != "" && p.Key != ""
}
