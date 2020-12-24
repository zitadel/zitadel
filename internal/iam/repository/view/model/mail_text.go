package model

import (
	"encoding/json"
	"time"

	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"

	es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	MailTextKeyAggregateID  = "aggregate_id"
	MailTextKeyMailTextType = "mail_text_type"
	MailTextKeyLanguage     = "language"
)

type MailTextView struct {
	AggregateID  string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
	State        int32     `json:"-" gorm:"column:mail_text_state"`

	MailTextType string `json:"mailTextType" gorm:"column:mail_text_type;primary_key"`
	Language     string `json:"language" gorm:"column:language;primary_key"`
	Title        string `json:"title" gorm:"column:title"`
	PreHeader    string `json:"preHeader" gorm:"column:pre_header"`
	Subject      string `json:"subject" gorm:"column:subject"`
	Greeting     string `json:"greeting" gorm:"column:greeting"`
	Text         string `json:"text" gorm:"column:text"`
	ButtonText   string `json:"buttonText" gorm:"column:button_text"`
	Default      bool   `json:"-" gorm:"-"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func MailTextViewFromModel(template *model.MailTextView) *MailTextView {
	return &MailTextView{
		AggregateID:  template.AggregateID,
		Sequence:     template.Sequence,
		CreationDate: template.CreationDate,
		ChangeDate:   template.ChangeDate,
		MailTextType: template.MailTextType,
		Language:     template.Language,
		Title:        template.Title,
		PreHeader:    template.PreHeader,
		Subject:      template.Subject,
		Greeting:     template.Greeting,
		Text:         template.Text,
		ButtonText:   template.ButtonText,
		Default:      template.Default,
	}
}

func MailTextsViewToModel(textsIn []*MailTextView, defaultIn bool) *model.MailTextsView {
	return &model.MailTextsView{
		Texts: mailTextsViewToModelArr(textsIn, defaultIn),
	}
}

func mailTextsViewToModelArr(texts []*MailTextView, defaultIn bool) []*model.MailTextView {
	result := make([]*model.MailTextView, len(texts))
	for i, r := range texts {
		r.Default = defaultIn
		result[i] = MailTextViewToModel(r)
	}
	return result
}

func MailTextViewToModel(template *MailTextView) *model.MailTextView {
	return &model.MailTextView{
		AggregateID:  template.AggregateID,
		Sequence:     template.Sequence,
		CreationDate: template.CreationDate,
		ChangeDate:   template.ChangeDate,
		MailTextType: template.MailTextType,
		Language:     template.Language,
		Title:        template.Title,
		PreHeader:    template.PreHeader,
		Subject:      template.Subject,
		Greeting:     template.Greeting,
		Text:         template.Text,
		ButtonText:   template.ButtonText,
		Default:      template.Default,
	}
}

func (i *MailTextView) AppendEvent(event *models.Event) (err error) {
	i.Sequence = event.Sequence
	switch event.Type {
	case es_model.MailTextAdded, org_es_model.MailTextAdded:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		err = i.SetData(event)
	case es_model.MailTextChanged, org_es_model.MailTextChanged:
		i.ChangeDate = event.CreationDate
		err = i.SetData(event)
	}
	return err
}

func (r *MailTextView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}

func (r *MailTextView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("MODEL-UFqAG").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-5CVaR", "Could not unmarshal data")
	}
	return nil
}
