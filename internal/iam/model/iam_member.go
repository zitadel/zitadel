package model

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type IamMember struct {
	es_models.ObjectRoot

	UserID string
	Roles  []string
}

func NewIamMember(iamID, userID string) *IamMember {
	return &IamMember{ObjectRoot: es_models.ObjectRoot{AggregateID: iamID}, UserID: userID}
}

func (i *IamMember) IsValid() bool {
	return i.AggregateID != "" && i.UserID != "" && len(i.Roles) != 0
}
