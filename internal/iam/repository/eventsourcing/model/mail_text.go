package model

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type MailText struct {
	models.ObjectRoot
	State        int32 `json:"-"`
	MailTextType string
	Language     string
	Title        string
	PreHeader    string
	Subject      string
	Greeting     string
	Text         string
	ButtonText   string
}

func GetDefaultMailText(mailTexts []*MailText, mailTextType string, language string) (int, *MailText) {
	for i, m := range mailTexts {
		if m.MailTextType == mailTextType && m.Language == language {
			return i, m
		}
	}
	return -1, nil
}

func MailTextsToModel(members []*MailText) []*iam_model.MailText {
	convertedMailTexts := make([]*iam_model.MailText, len(members))
	for i, m := range members {
		convertedMailTexts[i] = MailTextToModel(m)
	}
	return convertedMailTexts
}

func MailTextToModel(policy *MailText) *iam_model.MailText {
	return &iam_model.MailText{
		ObjectRoot:   policy.ObjectRoot,
		State:        iam_model.PolicyState(policy.State),
		MailTextType: policy.MailTextType,
		Language:     policy.Language,
		Title:        policy.Title,
		PreHeader:    policy.PreHeader,
		Subject:      policy.Subject,
		Greeting:     policy.Greeting,
		Text:         policy.Text,
		ButtonText:   policy.ButtonText,
	}
}

func MailTextsFromModel(members []*iam_model.MailText) []*MailText {
	convertedMailTexts := make([]*MailText, len(members))
	for i, m := range members {
		convertedMailTexts[i] = MailTextFromModel(m)
	}
	return convertedMailTexts
}

func MailTextFromModel(policy *iam_model.MailText) *MailText {
	return &MailText{
		ObjectRoot:   policy.ObjectRoot,
		State:        int32(policy.State),
		MailTextType: policy.MailTextType,
		Language:     policy.Language,
		Title:        policy.Title,
		PreHeader:    policy.PreHeader,
		Subject:      policy.Subject,
		Greeting:     policy.Greeting,
		Text:         policy.Text,
		ButtonText:   policy.ButtonText,
	}
}

func (p *MailText) Changes(changed *MailText) map[string]interface{} {
	changes := make(map[string]interface{}, 2)

	if changed.MailTextType != p.MailTextType {
		changes["mailTextType"] = changed.MailTextType
	}

	if changed.Language != p.Language {
		changes["language"] = changed.Language
	}

	if changed.Title != p.Title {
		changes["title"] = changed.Title
	}

	if changed.PreHeader != p.PreHeader {
		changes["preHeader"] = changed.PreHeader
	}

	if changed.Subject != p.Subject {
		changes["subject"] = changed.Subject
	}

	if changed.Greeting != p.Greeting {
		changes["greeting"] = changed.Greeting
	}

	if changed.Text != p.Text {
		changes["text"] = changed.Text
	}

	if changed.ButtonText != p.ButtonText {
		changes["buttonText"] = changed.ButtonText
	}

	return changes
}

func (i *IAM) appendAddMailTextEvent(event *es_models.Event) error {
	mailText := &MailText{}
	err := mailText.SetDataLabel(event)
	if err != nil {
		return err
	}
	mailText.ObjectRoot.CreationDate = event.CreationDate
	i.DefaultMailTexts = append(i.DefaultMailTexts, mailText)
	return nil
}

func (i *IAM) appendChangeMailTextEvent(event *es_models.Event) error {
	mailText := &MailText{}
	err := mailText.SetDataLabel(event)
	if err != nil {
		return err
	}
	if n, m := GetDefaultMailText(i.DefaultMailTexts, mailText.MailTextType, mailText.Language); m != nil {
		i.DefaultMailTexts[n] = mailText
	}
	return nil
}

func (p *MailText) SetDataLabel(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "MODEL-3FUV5", "unable to unmarshal data")
	}
	return nil
}
