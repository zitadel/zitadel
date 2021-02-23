package model

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type MailTemplate struct {
	models.ObjectRoot

	State    PolicyState
	Default  bool
	Template []byte
}

func (p *MailTemplate) IsValid() bool {
	return p.ObjectRoot.AggregateID != ""
}
