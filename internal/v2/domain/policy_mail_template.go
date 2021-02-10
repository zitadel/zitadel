package domain

import "github.com/caos/zitadel/internal/eventstore/models"

type MailTemplate struct {
	models.ObjectRoot

	State    PolicyState
	Default  bool
	Template []byte
}

func (m *MailTemplate) IsValid() bool {
	return m.Template != nil
}
