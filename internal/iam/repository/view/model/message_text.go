package model

import (
	"time"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"

	es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	MessageTextKeyAggregateID     = "aggregate_id"
	MessageTextKeyMessageTextType = "message_text_type"
	MessageTextKeyLanguage        = "language"
)

type MessageTextView struct {
	AggregateID  string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
	State        int32     `json:"-" gorm:"column:message_text_state"`

	MessageTextType string `json:"-" gorm:"column:message_text_type;primary_key"`
	Language        string `json:"-" gorm:"column:language;primary_key"`
	Title           string `json:"-" gorm:"column:title"`
	PreHeader       string `json:"-" gorm:"column:pre_header"`
	Subject         string `json:"-" gorm:"column:subject"`
	Greeting        string `json:"-" gorm:"column:greeting"`
	Text            string `json:"-" gorm:"column:text"`
	ButtonText      string `json:"-" gorm:"column:button_text"`
	FooterText      string `json:"-" gorm:"column:footer_text"`
	Default         bool   `json:"-" gorm:"-"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func MessageTextViewFromModel(template *model.MessageTextView) *MessageTextView {
	return &MessageTextView{
		AggregateID:     template.AggregateID,
		Sequence:        template.Sequence,
		CreationDate:    template.CreationDate,
		ChangeDate:      template.ChangeDate,
		MessageTextType: template.MessageTextType,
		Language:        template.Language.String(),
		Title:           template.Title,
		PreHeader:       template.PreHeader,
		Subject:         template.Subject,
		Greeting:        template.Greeting,
		Text:            template.Text,
		ButtonText:      template.ButtonText,
		FooterText:      template.FooterText,
		Default:         template.Default,
	}
}

func messageTextsViewToModelArr(texts []*MessageTextView, defaultIn bool) []*model.MessageTextView {
	result := make([]*model.MessageTextView, len(texts))
	for i, r := range texts {
		r.Default = defaultIn
		result[i] = MessageTextViewToModel(r)
	}
	return result
}

func MessageTextViewToModel(template *MessageTextView) *model.MessageTextView {
	lang := language.Make(template.Language)
	return &model.MessageTextView{
		AggregateID:     template.AggregateID,
		Sequence:        template.Sequence,
		CreationDate:    template.CreationDate,
		ChangeDate:      template.ChangeDate,
		MessageTextType: template.MessageTextType,
		Language:        lang,
		Title:           template.Title,
		PreHeader:       template.PreHeader,
		Subject:         template.Subject,
		Greeting:        template.Greeting,
		Text:            template.Text,
		ButtonText:      template.ButtonText,
		FooterText:      template.FooterText,
		Default:         template.Default,
	}
}

func (i *MessageTextView) AppendEvent(event *models.Event) (err error) {
	i.Sequence = event.Sequence
	switch event.Type {
	case es_model.CustomTextSet, org_es_model.CustomTextSet:
		i.setRootData(event)
		customText := new(CustomTextView)
		err = customText.SetData(event)
		if err != nil {
			return err
		}
		if customText.Key == domain.MessageTitle {
			i.Title = customText.Text
		}
		if customText.Key == domain.MessagePreHeader {
			i.PreHeader = customText.Text
		}
		if customText.Key == domain.MessageSubject {
			i.Subject = customText.Text
		}
		if customText.Key == domain.MessageGreeting {
			i.Greeting = customText.Text
		}
		if customText.Key == domain.MessageText {
			i.Text = customText.Text
		}
		if customText.Key == domain.MessageButtonText {
			i.ButtonText = customText.Text
		}
		if customText.Key == domain.MessageFooterText {
			i.FooterText = customText.Text
		}
		i.ChangeDate = event.CreationDate
	case es_model.CustomTextRemoved, org_es_model.CustomTextRemoved:
		customText := new(CustomTextView)
		err = customText.SetData(event)
		if err != nil {
			return err
		}
		if customText.Key == domain.MessageTitle {
			i.Title = ""
		}
		if customText.Key == domain.MessagePreHeader {
			i.PreHeader = ""
		}
		if customText.Key == domain.MessageSubject {
			i.Subject = ""
		}
		if customText.Key == domain.MessageGreeting {
			i.Greeting = ""
		}
		if customText.Key == domain.MessageText {
			i.Text = ""
		}
		if customText.Key == domain.MessageButtonText {
			i.ButtonText = ""
		}
		if customText.Key == domain.MessageFooterText {
			i.FooterText = ""
		}
		i.ChangeDate = event.CreationDate
	case org_es_model.CustomTextMessageRemoved:
		i.State = int32(model.PolicyStateRemoved)
	}
	return err
}

func (r *MessageTextView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}
