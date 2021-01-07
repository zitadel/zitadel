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

type AddressState int32

const (
	AddressStateUnspecified AddressState = iota
	AddressStateActive
	AddressStateRemoved

	addressStateCount
)

func (s AddressState) Valid() bool {
	return s >= 0 && s < addressStateCount
}
