package command

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/policy"
)

type CustomMessageTextReadModel struct {
	eventstore.WriteModel

	MessageTextType string
	Language        language.Tag
	Title           string
	PreHeader       string
	Subject         string
	Greeting        string
	Text            string
	ButtonText      string
	FooterText      string

	State domain.PolicyState
}

func (wm *CustomMessageTextReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.CustomTextSetEvent:
			if e.Template != wm.MessageTextType || wm.Language != e.Language {
				continue
			}
			if e.Key == domain.MessageSubject {
				wm.Subject = e.Text
			}
			if e.Key == domain.MessageTitle {
				wm.Title = e.Text
			}
			if e.Key == domain.MessagePreHeader {
				wm.PreHeader = e.Text
			}
			if e.Key == domain.MessageText {
				wm.Text = e.Text
			}
			if e.Key == domain.MessageGreeting {
				wm.Greeting = e.Text
			}
			if e.Key == domain.MessageButtonText {
				wm.ButtonText = e.Text
			}
			if e.Key == domain.MessageFooterText {
				wm.FooterText = e.Text
			}
			wm.State = domain.PolicyStateActive
		case *policy.CustomTextRemovedEvent:
			if e.Key != wm.MessageTextType || wm.Language != e.Language {
				continue
			}
			if e.Key == domain.MessageSubject {
				wm.Subject = ""
			}
			if e.Key == domain.MessageTitle {
				wm.Title = ""
			}
			if e.Key == domain.MessagePreHeader {
				wm.PreHeader = ""
			}
			if e.Key == domain.MessageText {
				wm.Text = ""
			}
			if e.Key == domain.MessageGreeting {
				wm.Greeting = ""
			}
			if e.Key == domain.MessageButtonText {
				wm.ButtonText = ""
			}
			if e.Key == domain.MessageFooterText {
				wm.FooterText = ""
			}
		case *policy.CustomTextTemplateRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

type CustomMessageTemplatesReadModel struct {
	eventstore.WriteModel
	CustomMessageTemplate map[string]*CustomText
}

type CustomText struct {
	Template   string
	Language   language.Tag
	Title      string
	PreHeader  string
	Subject    string
	Greeting   string
	Text       string
	ButtonText string
	FooterText string
	State      domain.PolicyState
}

func (wm *CustomMessageTemplatesReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.CustomTextSetEvent:
			if _, ok := wm.CustomMessageTemplate[e.Template+e.Language.String()]; !ok {
				wm.CustomMessageTemplate[e.Template+e.Language.String()] = &CustomText{Language: e.Language, Template: e.Template}
			}
			if e.Key == domain.MessageSubject {
				wm.CustomMessageTemplate[e.Template+e.Language.String()].Subject = e.Text
			}
			if e.Key == domain.MessageTitle {
				wm.CustomMessageTemplate[e.Template+e.Language.String()].Title = e.Text
			}
			if e.Key == domain.MessagePreHeader {
				wm.CustomMessageTemplate[e.Template+e.Language.String()].PreHeader = e.Text
			}
			if e.Key == domain.MessageText {
				wm.CustomMessageTemplate[e.Template+e.Language.String()].Text = e.Text
			}
			if e.Key == domain.MessageGreeting {
				wm.CustomMessageTemplate[e.Template+e.Language.String()].Greeting = e.Text
			}
			if e.Key == domain.MessageButtonText {
				wm.CustomMessageTemplate[e.Template+e.Language.String()].ButtonText = e.Text
			}
			if e.Key == domain.MessageFooterText {
				wm.CustomMessageTemplate[e.Template+e.Language.String()].FooterText = e.Text
			}
			wm.CustomMessageTemplate[e.Template+e.Language.String()].State = domain.PolicyStateActive
		case *policy.CustomTextRemovedEvent:
			if _, ok := wm.CustomMessageTemplate[e.Template+e.Language.String()]; !ok {
				wm.CustomMessageTemplate[e.Template+e.Language.String()] = new(CustomText)
			}
			if e.Key == domain.MessageSubject {
				wm.CustomMessageTemplate[e.Template+e.Language.String()].Subject = ""
			}
			if e.Key == domain.MessageTitle {
				wm.CustomMessageTemplate[e.Template+e.Language.String()].Title = ""
			}
			if e.Key == domain.MessagePreHeader {
				wm.CustomMessageTemplate[e.Template+e.Language.String()].PreHeader = ""
			}
			if e.Key == domain.MessageText {
				wm.CustomMessageTemplate[e.Template+e.Language.String()].Text = ""
			}
			if e.Key == domain.MessageGreeting {
				wm.CustomMessageTemplate[e.Template+e.Language.String()].Greeting = ""
			}
			if e.Key == domain.MessageButtonText {
				wm.CustomMessageTemplate[e.Template+e.Language.String()].ButtonText = ""
			}
			if e.Key == domain.MessageFooterText {
				wm.CustomMessageTemplate[e.Template+e.Language.String()].FooterText = ""
			}
		case *policy.CustomTextTemplateRemovedEvent:
			if _, ok := wm.CustomMessageTemplate[e.Template+e.Language.String()]; ok {
				delete(wm.CustomMessageTemplate, e.Template)
			}
		}
	}
	return wm.WriteModel.Reduce()
}
