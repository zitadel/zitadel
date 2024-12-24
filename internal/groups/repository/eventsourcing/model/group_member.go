package model

import es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"

type GroupMember struct {
	es_models.ObjectRoot

	UserID string
}
