package command

import (
	"strings"

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
			if e.Key != wm.MessageTextType || wm.Language != e.Language {
				continue
			}

			if strings.HasSuffix(e.Key, domain.MailSubject) {
				wm.Subject = e.Text
			}
			if strings.HasSuffix(e.Key, domain.MailTitle) {
				wm.Title = e.Text
			}
			if strings.HasSuffix(e.Key, domain.MailPreHeader) {
				wm.PreHeader = e.Text
			}
			if strings.HasSuffix(e.Key, domain.MailGreeting) {
				wm.Greeting = e.Text
			}
			if strings.HasSuffix(e.Key, domain.MailButtonText) {
				wm.ButtonText = e.Text
			}
			if strings.HasSuffix(e.Key, domain.MailFooterText) {
				wm.FooterText = e.Text
			}
			wm.State = domain.PolicyStateActive
		case *policy.CustomTextRemovedEvent:
			if e.Key != wm.MessageTextType || wm.Language != e.Language {
				continue
			}
			if strings.HasSuffix(e.Key, domain.MailSubject) {
				wm.Subject = ""
			}
			if strings.HasSuffix(e.Key, domain.MailTitle) {
				wm.Title = ""
			}
			if strings.HasSuffix(e.Key, domain.MailPreHeader) {
				wm.PreHeader = ""
			}
			if strings.HasSuffix(e.Key, domain.MailGreeting) {
				wm.Greeting = ""
			}
			if strings.HasSuffix(e.Key, domain.MailButtonText) {
				wm.ButtonText = ""
			}
			if strings.HasSuffix(e.Key, domain.MailFooterText) {
				wm.FooterText = ""
			}
		case *policy.CustomTextMessageRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
