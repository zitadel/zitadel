package model

import (
	"encoding/json"
	"reflect"

	"github.com/zitadel/logging"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
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

func (g *ProjectGrant) Changes(changed *ProjectGrant) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	changes["grantId"] = g.GrantID
	if !reflect.DeepEqual(g.RoleKeys, changed.RoleKeys) {
		changes["roleKeys"] = changed.RoleKeys
	}
	return changes
}

func (g *ProjectGrant) getData(event *es_models.Event) error {
	g.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, g); err != nil {
		logging.Log("EVEN-4h6gd").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
