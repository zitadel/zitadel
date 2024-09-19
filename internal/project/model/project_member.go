package model

import es_models "github.com/zitadel/zitadel/v2/internal/eventstore/v1/models"

type ProjectMember struct {
	es_models.ObjectRoot

	UserID string
	Roles  []string
}
