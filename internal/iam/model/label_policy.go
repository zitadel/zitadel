package model

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type LabelPolicy struct {
	models.ObjectRoot

	State          PolicyState
	Default        bool
	PrimaryColor   string
	SecondaryColor string
}

func (p *LabelPolicy) IsValid() bool {
	return p.ObjectRoot.AggregateID != ""
}
