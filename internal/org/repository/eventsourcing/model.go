package eventsourcing

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
)

const (
	orgVersion = "v1"
)

type Org struct {
	es_models.ObjectRoot `json:"-"`

	Name   string `json:"name"`
	Domain string `json:"domain"`
	State  int32  `json:"-"`
}

func OrgFromModel(org *org_model.Org) *Org {
	return &Org{
		ObjectRoot: es_models.ObjectRoot{
			ID:           org.ID,
			Sequence:     org.Sequence,
			ChangeDate:   org.ChangeDate,
			CreationDate: org.CreationDate,
		},
		Domain: org.Domain,
		Name:   org.Name,
		State:  org_model.ProjectStateToInt(org.State),
	}
}

func OrgToModel(org *Org) *org_model.Org {
	return &org_model.Org{
		ObjectRoot: es_models.ObjectRoot{
			ID:           org.ID,
			Sequence:     org.Sequence,
			ChangeDate:   org.ChangeDate,
			CreationDate: org.CreationDate,
		},
		Domain: org.Domain,
		Name:   org.Name,
		State:  org_model.ProjectStateFromInt(org.State),
	}
}

func OrgFromEvents(org *Org, events ...*es_models.Event) (*Org, error) {
	if org == nil {
		org = new(Org)
	}

	return org, org.AppendEvents(events...)
}

func (o *Org) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		err := o.AppendEvent(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Org) AppendEvent(event *es_models.Event) error {
	o.ObjectRoot.AppendEvent(event)

	switch event.Type {
	case org_model.OrgAdded, org_model.OrgChanged:
		err := json.Unmarshal(event.Data, o)
		if err != nil {
			return errors.ThrowInternal(err, "EVENT-BpbQZ", "unable to unmarshal event")
		}
	case org_model.OrgDeactivated:
		o.State = org_model.ProjectStateToInt(org_model.Inactive)
	case org_model.OrgReactivated:
		o.State = org_model.ProjectStateToInt(org_model.Active)
	}

	return nil
}

func (o *Org) Changes(changed *Org) map[string]interface{} {
	changes := make(map[string]interface{}, 2)

	if changed.Name != "" && changed.Name != o.Name {
		changes["name"] = changed.Name
	}
	if changed.Domain != "" && changed.Domain != o.Domain {
		changes["domain"] = changed.Domain
	}

	return changes
}
