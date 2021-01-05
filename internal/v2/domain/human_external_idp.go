package domain

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type ExternalIDP struct {
	es_models.ObjectRoot

	IDPConfigID string
	UserID      string
	DisplayName string
}
