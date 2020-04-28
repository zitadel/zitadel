package model

import (
	"encoding/json"
	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/usergrant/model"
	"reflect"
)

const (
	UserGrantVersion = "v1"
)

type UserGrant struct {
	es_models.ObjectRoot

	State     int32    `json:"-"`
	UserID    string   `json:"userId,omitempty"`
	ProjectID string   `json:"projectId,omitempty"`
	RoleKeys  []string `json:"roleKeys,omitempty"`
}

type UserGrantID struct {
	es_models.ObjectRoot
	GrantID string `json:"grantId"`
}

func (g *UserGrant) Changes(changed *UserGrant) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if !reflect.DeepEqual(g.RoleKeys, changed.RoleKeys) {
		changes["roleKeys"] = changed.RoleKeys
	}
	return changes
}

func UserGrantFromModel(grant *model.UserGrant) *UserGrant {
	return &UserGrant{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  grant.ObjectRoot.AggregateID,
			Sequence:     grant.Sequence,
			ChangeDate:   grant.ChangeDate,
			CreationDate: grant.CreationDate,
		},
		UserID:    grant.UserID,
		ProjectID: grant.ProjectID,
		State:     int32(grant.State),
		RoleKeys:  grant.RoleKeys,
	}
}

func UserGrantToModel(grant *UserGrant) *model.UserGrant {
	return &model.UserGrant{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  grant.AggregateID,
			ChangeDate:   grant.ChangeDate,
			CreationDate: grant.CreationDate,
			Sequence:     grant.Sequence,
		},
		UserID:    grant.UserID,
		ProjectID: grant.ProjectID,
		State:     model.UserGrantState(grant.State),
		RoleKeys:  grant.RoleKeys,
	}
}

func (g *UserGrant) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := g.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (g *UserGrant) AppendEvent(event *es_models.Event) error {
	g.ObjectRoot.AppendEvent(event)
	switch event.Type {
	case model.UserGrantAdded,
		model.UserGrantChanged:
		return g.setData(event)
	case model.UserGrantDeactivated:
		return g.appendGrantStateEvent(model.USERGRANTSTATE_INACTIVE)
	case model.UserGrantReactivated:
		return g.appendGrantStateEvent(model.USERGRANTSTATE_ACTIVE)
	case model.UserGrantRemoved:
		return g.appendGrantStateEvent(model.USERGRANTSTATE_REMOVED)
	}
	return nil
}

func (g *UserGrant) appendAddGrantEvent(event *es_models.Event) error {
	return g.setData(event)
}

func (g *UserGrant) appendChangeGrantEvent(event *es_models.Event) error {
	return g.setData(event)
}

func (g *UserGrant) appendGrantStateEvent(state model.UserGrantState) error {
	g.State = int32(state)
	return nil
}

func (g *UserGrant) setData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, g); err != nil {
		logging.Log("EVEN-lso9x").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
