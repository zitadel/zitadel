package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type LabelPolicy struct {
	models.ObjectRoot

	State          PolicyState
	Default        bool
	PrimaryColor   string
	SecundaryColor string
}

func (p *LabelPolicy) IsValid() bool {
	return p.ObjectRoot.AggregateID != ""
}
