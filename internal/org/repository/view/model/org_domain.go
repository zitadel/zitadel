package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/org/model"
	es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

const (
	OrgDomainKeyOrgID    = "org_id"
	OrgDomainKeyDomain   = "domain"
	OrgDomainKeyVerified = "verified"
	OrgDomainKeyPrimary  = "primary_domain"
)

type OrgDomainView struct {
	Domain         string `json:"domain" gorm:"column:domain;primary_key"`
	OrgID          string `json:"-" gorm:"column:org_id;primary_key"`
	Verified       bool   `json:"-" gorm:"column:verified"`
	Primary        bool   `json:"-" gorm:"column:primary_domain"`
	ValidationType int32  `json:"validationType" gorm:"column:validation_type"`
	Sequence       uint64 `json:"-" gorm:"column:sequence"`

	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
}

func OrgDomainViewFromModel(domain *model.OrgDomainView) *OrgDomainView {
	return &OrgDomainView{
		OrgID:          domain.OrgID,
		Domain:         domain.Domain,
		Primary:        domain.Primary,
		Verified:       domain.Verified,
		ValidationType: int32(domain.ValidationType),
		CreationDate:   domain.CreationDate,
		ChangeDate:     domain.ChangeDate,
	}
}

func OrgDomainToModel(domain *OrgDomainView) *model.OrgDomainView {
	return &model.OrgDomainView{
		OrgID:          domain.OrgID,
		Domain:         domain.Domain,
		Primary:        domain.Primary,
		Verified:       domain.Verified,
		ValidationType: model.OrgDomainValidationType(domain.ValidationType),
		CreationDate:   domain.CreationDate,
		ChangeDate:     domain.ChangeDate,
	}
}

func OrgDomainsToModel(domain []*OrgDomainView) []*model.OrgDomainView {
	result := make([]*model.OrgDomainView, len(domain))
	for i, r := range domain {
		result[i] = OrgDomainToModel(r)
	}
	return result
}

func (d *OrgDomainView) AppendEvent(event *models.Event) (err error) {
	d.Sequence = event.Sequence
	d.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.OrgDomainAdded:
		d.setRootData(event)
		d.CreationDate = event.CreationDate
		err = d.SetData(event)
	case es_model.OrgDomainVerificationAdded:
		err = d.SetData(event)
	case es_model.OrgDomainVerified:
		d.Verified = true
	case es_model.OrgDomainPrimarySet:
		d.Primary = true
	}
	return err
}

func (r *OrgDomainView) setRootData(event *models.Event) {
	r.OrgID = event.AggregateID
}

func (r *OrgDomainView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-sj4Sf").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lub6s", "Could not unmarshal data")
	}
	return nil
}
