package model

import es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"

type ProjectGrantGroupMember struct {
	es_models.ObjectRoot
	GrantID string
	GroupID string
	Roles   []string
}
