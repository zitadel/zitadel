package model

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type OrgMember struct {
	es_models.ObjectRoot
	UserID string
	Roles  []string
}

func NewOrgMember(orgID, userID string) *OrgMember {
	return &OrgMember{ObjectRoot: es_models.ObjectRoot{AggregateID: orgID}, UserID: userID}
}

func (member *OrgMember) IsValid() bool {
	return member.AggregateID != "" && member.UserID != "" && len(member.Roles) != 0
}
