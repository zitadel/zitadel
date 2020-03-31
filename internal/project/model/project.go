package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type Project struct {
	es_models.ObjectRoot

	State ProjectState
	Name  string
}

type ProjectState es_models.Enum

var states = []string{"Active", "Inactive"}

func NewProject(id string) *Project {
	return &Project{ObjectRoot: es_models.ObjectRoot{ID: id}, State: Active}
}

func (p *Project) Changes(changed *Project) map[string]interface{} {
	changes := make(map[string]interface{}, 2)
	if changed.Name != "" && p.Name != changed.Name {
		changes["name"] = changed.Name
	}
	return changes
}
