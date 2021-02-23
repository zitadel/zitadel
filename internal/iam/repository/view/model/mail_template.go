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
	MailTemplateKeyAggregateID = "aggregate_id"
)

type MailTemplateView struct {
	AggregateID  string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
	State        int32     `json:"-" gorm:"column:mail_template_state"`

	Template []byte `json:"template" gorm:"column:template"`
	Default  bool   `json:"-" gorm:"-"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func MailTemplateViewFromModel(template *model.MailTemplateView) *MailTemplateView {
	return &MailTemplateView{
		AggregateID:  template.AggregateID,
		Sequence:     template.Sequence,
		CreationDate: template.CreationDate,
		ChangeDate:   template.ChangeDate,
		Template:     template.Template,
		Default:      template.Default,
	}
}

func MailTemplateViewToModel(template *MailTemplateView) *model.MailTemplateView {
	return &model.MailTemplateView{
		AggregateID:  template.AggregateID,
		Sequence:     template.Sequence,
		CreationDate: template.CreationDate,
		ChangeDate:   template.ChangeDate,
		Template:     template.Template,
		Default:      template.Default,
	}
}

func (i *MailTemplateView) AppendEvent(event *models.Event) (err error) {
	i.Sequence = event.Sequence
	i.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.MailTemplateAdded, org_es_model.MailTemplateAdded:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		err = i.SetData(event)
	case es_model.MailTemplateChanged, org_es_model.MailTemplateChanged:
		i.ChangeDate = event.CreationDate
		err = i.SetData(event)
	}
	return err
}

func (r *MailTemplateView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}

func (r *MailTemplateView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("MODEL-YDZmZ").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-sKWwO", "Could not unmarshal data")
	}
	return nil
}
