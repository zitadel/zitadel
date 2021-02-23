package domain

import "github.com/caos/zitadel/internal/eventstore/v1/models"

type MailText struct {
	models.ObjectRoot

	State        PolicyState
	Default      bool
	MailTextType string
	Language     string
	Title        string
	PreHeader    string
	Subject      string
	Greeting     string
	Text         string
	ButtonText   string
}

func (m *MailText) IsValid() bool {
	return m.MailTextType != "" && m.Language != "" && m.Title != "" && m.PreHeader != "" && m.Subject != "" && m.Greeting != "" && m.Text != "" && m.ButtonText != ""
}
