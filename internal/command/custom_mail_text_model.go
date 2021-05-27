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
			if !strings.HasSuffix(e.Key, wm.MessageTextType) || wm.Language != e.Language {
				continue
			}
			if e.Key == wm.MessageTextType+domain.MailSubject {
				wm.Subject = e.Text
			}
			if e.Key == wm.MessageTextType+domain.MailTitle {
				wm.Title = e.Text
			}
			if e.Key == wm.MessageTextType+domain.MailPreHeader {
				wm.PreHeader = e.Text
			}
			if e.Key == wm.MessageTextType+domain.MailGreeting {
				wm.Greeting = e.Text
			}
			if e.Key == wm.MessageTextType+domain.MailButtonText {
				wm.ButtonText = e.Text
			}
			if e.Key == wm.MessageTextType+domain.MailFooterText {
				wm.FooterText = e.Text
			}
			wm.State = domain.PolicyStateActive
		case *policy.CustomTextRemovedEvent:
			if !strings.HasSuffix(e.Key, wm.MessageTextType) || wm.Language != e.Language {
				continue
			}
			if e.Key == wm.MessageTextType+domain.MailSubject {
				wm.Subject = ""
			}
			if e.Key == wm.MessageTextType+domain.MailTitle {
				wm.Title = ""
			}
			if e.Key == wm.MessageTextType+domain.MailPreHeader {
				wm.PreHeader = ""
			}
			if e.Key == wm.MessageTextType+domain.MailGreeting {
				wm.Greeting = ""
			}
			if e.Key == wm.MessageTextType+domain.MailButtonText {
				wm.ButtonText = ""
			}
			if e.Key == wm.MessageTextType+domain.MailFooterText {
				wm.FooterText = ""
			}
		}
	}
	return wm.WriteModel.Reduce()
}
