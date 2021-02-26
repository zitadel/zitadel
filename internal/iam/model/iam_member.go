package model

import es_models "github.com/caos/zitadel/internal/eventstore/v1/models"

type IAMMember struct {
	es_models.ObjectRoot

	UserID string
	Roles  []string
}

func NewIAMMember(iamID, userID string) *IAMMember {
	return &IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: iamID}, UserID: userID}
}

func (i *IAMMember) IsValid() bool {
	return i.AggregateID != "" && i.UserID != "" && len(i.Roles) != 0
}
