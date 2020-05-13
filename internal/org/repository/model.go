package repository

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type Org struct {
	es_models.ObjectRoot

	Name   string
	Domain string
}
