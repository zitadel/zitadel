package model

import (
	"encoding/json"

	"github.com/zitadel/logging"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Address struct {
	es_models.ObjectRoot

	Country       string `json:"country,omitempty"`
	Locality      string `json:"locality,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Region        string `json:"region,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`
}

func (a *Address) Changes(changed *Address) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if a.Country != changed.Country {
		changes["country"] = changed.Country
	}
	if a.Locality != changed.Locality {
		changes["locality"] = changed.Locality
	}
	if a.PostalCode != changed.PostalCode {
		changes["postalCode"] = changed.PostalCode
	}
	if a.Region != changed.Region {
		changes["region"] = changed.Region
	}
	if a.StreetAddress != changed.StreetAddress {
		changes["streetAddress"] = changed.StreetAddress
	}
	return changes
}

func (u *Human) appendUserAddressChangedEvent(event *es_models.Event) error {
	if u.Address == nil {
		u.Address = new(Address)
	}
	return u.Address.setData(event)
}

func (a *Address) setData(event *es_models.Event) error {
	a.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-clos0").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-so92s", "could not unmarshal event")
	}
	return nil
}
