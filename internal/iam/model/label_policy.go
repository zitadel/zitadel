package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type LabelPolicy struct {
	models.ObjectRoot

	State                 PolicyState
	Default               bool
	PrimaryColor          string
	SecondaryColor        string
	HideLoginPolicySuffix bool
}

func (p *LabelPolicy) IsValid() bool {
	return p.ObjectRoot.AggregateID != ""
}
