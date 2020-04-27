package model

import (
	"encoding/json"
	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"reflect"
)

type UserGrant struct {
	es_models.ObjectRoot

	State     int32    `json:"-"`
	GrantID   string   `json:"grantId,omitempty"`
	ProjectID string   `json:"projectId,omitempty"`
	RoleKeys  []string `json:"roleKeys,omitempty"`
}

type UserGrantID struct {
	es_models.ObjectRoot
	GrantID string `json:"grantId"`
}

func GetUserGrant(grants []*UserGrant, grantId string) (int, *UserGrant) {
	for i, g := range grants {
		if g.GrantID == grantId {
			return i, g
		}
	}
	return -1, nil
}

func (g *UserGrant) Changes(changed *UserGrant) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	changes["grantId"] = g.GrantID
	if !reflect.DeepEqual(g.RoleKeys, changed.RoleKeys) {
		changes["roleKeys"] = changed.RoleKeys
	}
	return changes
}

func GrantsToModel(grants []*UserGrant) []*model.UserGrant {
	convertedGrants := make([]*model.UserGrant, len(grants))
	for i, g := range grants {
		convertedGrants[i] = GrantToModel(g)
	}
	return convertedGrants
}

func GrantsFromModel(grants []*model.UserGrant) []*UserGrant {
	convertedGrants := make([]*UserGrant, len(grants))
	for i, g := range grants {
		convertedGrants[i] = GrantFromModel(g)
	}
	return convertedGrants
}

func GrantFromModel(grant *model.UserGrant) *UserGrant {
	return &UserGrant{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  grant.ObjectRoot.AggregateID,
			Sequence:     grant.Sequence,
			ChangeDate:   grant.ChangeDate,
			CreationDate: grant.CreationDate,
		},
		GrantID:   grant.GrantID,
		ProjectID: grant.ProjectID,
		State:     int32(grant.State),
		RoleKeys:  grant.RoleKeys,
	}
}

func GrantToModel(grant *UserGrant) *model.UserGrant {
	return &model.UserGrant{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  grant.AggregateID,
			ChangeDate:   grant.ChangeDate,
			CreationDate: grant.CreationDate,
			Sequence:     grant.Sequence,
		},
		GrantID:   grant.GrantID,
		ProjectID: grant.ProjectID,
		State:     model.UserGrantState(grant.State),
		RoleKeys:  grant.RoleKeys,
	}
}

func (u *User) appendAddGrantEvent(event *es_models.Event) error {
	grant := new(UserGrant)
	err := grant.getData(event)
	if err != nil {
		return err
	}
	grant.ObjectRoot.CreationDate = event.CreationDate
	u.Grants = append(u.Grants, grant)
	return nil
}

func (u *User) appendChangeGrantEvent(event *es_models.Event) error {
	grant := new(UserGrant)
	err := grant.getData(event)
	if err != nil {
		return err
	}
	if i, g := GetUserGrant(u.Grants, grant.GrantID); g != nil {
		u.Grants[i].getData(event)
	}
	return nil
}

func (u *User) appendGrantStateEvent(event *es_models.Event, state model.UserGrantState) error {
	grant := new(UserGrant)
	err := grant.getData(event)
	if err != nil {
		return err
	}
	if i, g := GetUserGrant(u.Grants, grant.GrantID); g != nil {
		g.State = int32(state)
		u.Grants[i] = g
	}
	return nil
}

func (u *User) appendRemoveGrantEvent(event *es_models.Event) error {
	grant := new(UserGrant)
	err := grant.getData(event)
	if err != nil {
		return err
	}

	if i, g := GetUserGrant(u.Grants, grant.GrantID); g != nil {
		u.Grants[i] = u.Grants[len(u.Grants)-1]
		u.Grants[len(u.Grants)-1] = nil
		u.Grants = u.Grants[:len(u.Grants)-1]
	}
	return nil
}

func (g *UserGrant) getData(event *es_models.Event) error {
	g.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, g); err != nil {
		logging.Log("EVEN-lso9x").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
