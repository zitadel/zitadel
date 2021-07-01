package model

import (
	"encoding/json"
	"time"

	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"

	es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	PrivacyKeyAggregateID = "aggregate_id"
)

type PrivacyPolicyView struct {
	AggregateID  string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
	State        int32     `json:"-" gorm:"column:lockout_policy_state"`

	TOSLink     string `json:"tosLink" gorm:"column:tos_link"`
	PrivacyLink string `json:"privacyLink" gorm:"column:privacy_link"`
	Default     bool   `json:"-" gorm:"-"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func PrivacyViewFromModel(policy *model.PrivacyPolicyView) *PrivacyPolicyView {
	return &PrivacyPolicyView{
		AggregateID:  policy.AggregateID,
		Sequence:     policy.Sequence,
		CreationDate: policy.CreationDate,
		ChangeDate:   policy.ChangeDate,
		TOSLink:      policy.TOSLink,
		PrivacyLink:  policy.PrivacyLink,
		Default:      policy.Default,
	}
}

func PrivacyViewToModel(policy *PrivacyPolicyView) *model.PrivacyPolicyView {
	return &model.PrivacyPolicyView{
		AggregateID:  policy.AggregateID,
		Sequence:     policy.Sequence,
		CreationDate: policy.CreationDate,
		ChangeDate:   policy.ChangeDate,
		TOSLink:      policy.TOSLink,
		PrivacyLink:  policy.PrivacyLink,
		Default:      policy.Default,
	}
}

func (i *PrivacyPolicyView) AppendEvent(event *models.Event) (err error) {
	i.Sequence = event.Sequence
	i.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.PrivacyPolicyAdded, org_es_model.PrivacyPolicyAdded:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		err = i.SetData(event)
	case es_model.PrivacyPolicyChanged, org_es_model.PrivacyPolicyChanged:
		err = i.SetData(event)
	}
	return err
}

func (r *PrivacyPolicyView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}

func (r *PrivacyPolicyView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-gHls0").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}
