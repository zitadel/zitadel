package address

import "github.com/caos/zitadel/internal/eventstore/v2"

type HumanAddressWriteModel struct {
	eventstore.WriteModel

	Country       string
	Locality      string
	PostalCode    string
	Region        string
	StreetAddress string
}
