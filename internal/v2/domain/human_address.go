package domain

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type Address struct {
	es_models.ObjectRoot

	Country       string
	Locality      string
	PostalCode    string
	Region        string
	StreetAddress string
}
