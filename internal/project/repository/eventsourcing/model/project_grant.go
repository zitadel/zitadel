package model

import (
	"encoding/json"
	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/project/model"
	"reflect"
)

type ProjectGrant struct {
	es_models.ObjectRoot
	State        int32                 `json:"-"`
	GrantID      string                `json:"grantId,omitempty"`
	GrantedOrgID string                `json:"grantedOrgId,omitempty"`
	RoleKeys     []string              `json:"roleKeys,omitempty"`
	Members      []*ProjectGrantMember `json:"-"`
}

type ProjectGrantID struct {
	es_models.ObjectRoot
	GrantID string `json:"grantId"`
}

func GetProjectGrant(grants []*ProjectGrant, id string) (int, *ProjectGrant) {
	for i, g := range grants {
		if g.GrantID == id {
			return i, g
		}
	}
	return -1, nil
}

func GetProjectGrantByOrgID(grants []*ProjectGrant, resourceOwner string) (int, *ProjectGrant) {
	for i, g := range grants {
		if g.GrantedOrgID == resourceOwner {
			return i, g
		}
	}
	return -1, nil
}

func (g *ProjectGrant) Changes(changed *ProjectGrant) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	changes["grantId"] = g.GrantID
	if !reflect.DeepEqual(g.RoleKeys, changed.RoleKeys) {
		changes["roleKeys"] = changed.RoleKeys
	}
	return changes
}

func GrantsToModel(grants []*ProjectGrant) []*model.ProjectGrant {
	convertedGrants := make([]*model.ProjectGrant, len(grants))
	for i, g := range grants {
		convertedGrants[i] = GrantToModel(g)
	}
	return convertedGrants
}

func GrantsFromModel(grants []*model.ProjectGrant) []*ProjectGrant {
	convertedGrants := make([]*ProjectGrant, len(grants))
	for i, g := range grants {
		convertedGrants[i] = GrantFromModel(g)
	}
	return convertedGrants
}

func GrantFromModel(grant *model.ProjectGrant) *ProjectGrant {
	members := GrantMembersFromModel(grant.Members)
	return &ProjectGrant{
		ObjectRoot:   grant.ObjectRoot,
		GrantID:      grant.GrantID,
		GrantedOrgID: grant.GrantedOrgID,
		State:        int32(grant.State),
		RoleKeys:     grant.RoleKeys,
		Members:      members,
	}
}

func GrantToModel(grant *ProjectGrant) *model.ProjectGrant {
	members := GrantMembersToModel(grant.Members)
	return &model.ProjectGrant{
		ObjectRoot:   grant.ObjectRoot,
		GrantID:      grant.GrantID,
		GrantedOrgID: grant.GrantedOrgID,
		State:        model.ProjectGrantState(grant.State),
		RoleKeys:     grant.RoleKeys,
		Members:      members,
	}
}

func (p *Project) appendAddGrantEvent(event *es_models.Event) error {
	grant := new(ProjectGrant)
	err := grant.getData(event)
	if err != nil {
		return err
	}
	grant.ObjectRoot.CreationDate = event.CreationDate
	p.Grants = append(p.Grants, grant)
	return nil
}

func (p *Project) appendChangeGrantEvent(event *es_models.Event) error {
	grant := new(ProjectGrant)
	err := grant.getData(event)
	if err != nil {
		return err
	}
	if i, g := GetProjectGrant(p.Grants, grant.GrantID); g != nil {
		p.Grants[i].getData(event)
	}
	return nil
}

func (p *Project) appendGrantStateEvent(event *es_models.Event, state model.ProjectGrantState) error {
	grant := new(ProjectGrant)
	err := grant.getData(event)
	if err != nil {
		return err
	}
	if i, g := GetProjectGrant(p.Grants, grant.GrantID); g != nil {
		g.State = int32(state)
		p.Grants[i] = g
	}
	return nil
}

func (p *Project) appendRemoveGrantEvent(event *es_models.Event) error {
	grant := new(ProjectGrant)
	err := grant.getData(event)
	if err != nil {
		return err
	}

	if i, g := GetProjectGrant(p.Grants, grant.GrantID); g != nil {
		p.Grants[i] = p.Grants[len(p.Grants)-1]
		p.Grants[len(p.Grants)-1] = nil
		p.Grants = p.Grants[:len(p.Grants)-1]
	}
	return nil
}

func (g *ProjectGrant) getData(event *es_models.Event) error {
	g.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, g); err != nil {
		logging.Log("EVEN-4h6gd").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
