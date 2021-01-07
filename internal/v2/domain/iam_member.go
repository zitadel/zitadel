package domain

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type Member struct {
	es_models.ObjectRoot

	UserID string
	Roles  []string
}

func (i *Member) IsValid() bool {
	return i.AggregateID != "" && i.UserID != "" && len(i.Roles) != 0
}
