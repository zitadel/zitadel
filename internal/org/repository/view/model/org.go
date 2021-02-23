package model

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	org_model "github.com/caos/zitadel/internal/org/model"
)

const (
	OrgKeyOrgDomain     = "domain"
	OrgKeyOrgID         = "id"
	OrgKeyOrgName       = "name"
	OrgKeyResourceOwner = "resource_owner"
	OrgKeyState         = "org_state"
)

type OrgView struct {
	ID            string    `json:"-" gorm:"column:id;primary_key"`
	CreationDate  time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate    time.Time `json:"-" gorm:"column:change_date"`
	ResourceOwner string    `json:"-" gorm:"column:resource_owner"`
	State         int32     `json:"-" gorm:"column:org_state"`
	Sequence      uint64    `json:"-" gorm:"column:sequence"`

	Name   string `json:"name" gorm:"column:name"`
	Domain string `json:"domain" gorm:"column:domain"`
}

func OrgFromModel(org *org_model.OrgView) *OrgView {
	return &OrgView{
		ChangeDate:    org.ChangeDate,
		CreationDate:  org.CreationDate,
		ID:            org.ID,
		Name:          org.Name,
		ResourceOwner: org.ResourceOwner,
		Sequence:      org.Sequence,
		State:         int32(org.State),
	}
}

func OrgToModel(org *OrgView) *org_model.OrgView {
	return &org_model.OrgView{
		ChangeDate:    org.ChangeDate,
		CreationDate:  org.CreationDate,
		ID:            org.ID,
		Name:          org.Name,
		ResourceOwner: org.ResourceOwner,
		Sequence:      org.Sequence,
		State:         org_model.OrgState(org.State),
	}
}

func OrgsToModel(orgs []*OrgView) []*org_model.OrgView {
	modelOrgs := make([]*org_model.OrgView, len(orgs))

	for i, org := range orgs {
		modelOrgs[i] = OrgToModel(org)
	}

	return modelOrgs
}

func (o *OrgView) AppendEvent(event *es_models.Event) (err error) {
	switch event.Type {
	case model.OrgAdded:
		o.CreationDate = event.CreationDate
		o.State = int32(org_model.OrgStateActive)
		o.setRootData(event)
		err = o.SetData(event)
	case model.OrgChanged:
		o.setRootData(event)
		err = o.SetData(event)
	case model.OrgDeactivated:
		o.State = int32(org_model.OrgStateInactive)
	case model.OrgReactivated:
		o.State = int32(org_model.OrgStateActive)
	}
	return err
}

func (o *OrgView) setRootData(event *es_models.Event) {
	o.ChangeDate = event.CreationDate
	o.Sequence = event.Sequence
	o.ID = event.AggregateID
	o.ResourceOwner = event.ResourceOwner
}

func (o *OrgView) SetData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, o); err != nil {
		logging.Log("VIEW-5W7Op").WithError(err).Error("could not unmarshal event data")
		return errors.ThrowInternal(err, "VIEW-HZKME", "Could not unmarshal data")
	}
	return nil
}
