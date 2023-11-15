package model

import es_models "github.com/zitadel/zitadel/v2/internal/eventstore/v1/models"

type ProjectRole struct {
	es_models.ObjectRoot

	Key         string
	DisplayName string
	Group       string
}

func (p *ProjectRole) IsValid() bool {
	return p.AggregateID != "" && p.Key != ""
}
