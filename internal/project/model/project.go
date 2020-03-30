package model

import (
	es_pkg "github.com/caos/zitadel/internal/eventstore/pkg"
	"github.com/caos/zitadel/internal/model"
)

type Project struct {
	es_pkg.ObjectRoot

	State ProjectState
	Name  string
}

type ProjectState model.Enum

var states = []string{"Active", "Inactive"}

func NewProject(id string) *Project {
	return &Project{ObjectRoot: es_pkg.ObjectRoot{ID: id}, State: Active}
}

func (p *Project) Changes(changed *Project) map[string]interface{} {
	changes := make(map[string]interface{}, 2)
	if changed.Name != "" && p.Name != changed.Name {
		changes["name"] = changed.Name
	}
	return changes
}
