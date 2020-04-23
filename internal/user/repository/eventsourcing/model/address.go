package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
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

func AddressFromModel(address *model.Address) *Address {
	return &Address{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  address.ObjectRoot.AggregateID,
			Sequence:     address.Sequence,
			ChangeDate:   address.ChangeDate,
			CreationDate: address.CreationDate,
		},
		Country:       address.Country,
		Locality:      address.Locality,
		PostalCode:    address.PostalCode,
		Region:        address.Region,
		StreetAddress: address.StreetAddress,
	}
}

func AddressToModel(address *Address) *model.Address {
	return &model.Address{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  address.ObjectRoot.AggregateID,
			Sequence:     address.Sequence,
			ChangeDate:   address.ChangeDate,
			CreationDate: address.CreationDate,
		},
		Country:       address.Country,
		Locality:      address.Locality,
		PostalCode:    address.PostalCode,
		Region:        address.Region,
		StreetAddress: address.StreetAddress,
	}
}

func InitCodeToModel(code *InitUserCode) *model.InitUserCode {
	return &model.InitUserCode{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  code.ObjectRoot.AggregateID,
			Sequence:     code.Sequence,
			ChangeDate:   code.ChangeDate,
			CreationDate: code.CreationDate,
		},
		Expiry: code.Expiry,
		Code:   code.Code,
	}
}
