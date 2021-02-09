package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type MailTextWriteModel struct {
	eventstore.WriteModel

	MailTextType string
	Language     string
	Title        string
	PreHeader    string
	Subject      string
	Greeting     string
	Text         string
	ButtonText   string

	State domain.PolicyState
}

func (wm *MailTextWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.MailTextAddedEvent:
			if wm.MailTextType != e.MailTextType || wm.Language != e.Language {
				continue
			}
			wm.Title = e.Title
			wm.PreHeader = e.PreHeader
			wm.Subject = e.Subject
			wm.Greeting = e.Greeting
			wm.Text = e.Text
			wm.ButtonText = e.ButtonText
			wm.State = domain.PolicyStateActive
		case *policy.MailTextChangedEvent:
			if wm.MailTextType != e.MailTextType || wm.Language != e.Language {
				continue
			}
			if e.Title != nil {
				wm.Title = *e.Title
			}
			if e.PreHeader != nil {
				wm.PreHeader = *e.PreHeader
			}
			if e.Subject != nil {
				wm.Subject = *e.Subject
			}
			if e.Greeting != nil {
				wm.Greeting = *e.Greeting
			}
			if e.Text != nil {
				wm.Text = *e.Text
			}
			if e.ButtonText != nil {
				wm.ButtonText = *e.ButtonText
			}
		case *policy.MailTextRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
