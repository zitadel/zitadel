package model

import es_models "github.com/caos/zitadel/internal/eventstore/v1/models"

type Address struct {
	es_models.ObjectRoot

	Country       string
	Locality      string
	PostalCode    string
	Region        string
	StreetAddress string
}
