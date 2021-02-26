package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

type ExternalIDP struct {
	es_models.ObjectRoot

	IDPConfigID string
	UserID      string
	DisplayName string
}

func (idp *ExternalIDP) IsValid() bool {
	return idp.AggregateID != "" && idp.IDPConfigID != "" && idp.UserID != ""
}
