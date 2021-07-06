package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/user"
)

type HumanAddressWriteModel struct {
	eventstore.WriteModel

	Country       string
	Locality      string
	PostalCode    string
	Region        string
	StreetAddress string

	State domain.AddressState
}

func NewHumanAddressWriteModel(userID, resourceOwner string) *HumanAddressWriteModel {
	return &HumanAddressWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanAddressWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.Country = e.Country
			wm.Locality = e.Locality
			wm.PostalCode = e.PostalCode
			wm.Region = e.Region
			wm.StreetAddress = e.StreetAddress
			wm.State = domain.AddressStateActive
		case *user.HumanRegisteredEvent:
			wm.Country = e.Country
			wm.Locality = e.Locality
			wm.PostalCode = e.PostalCode
			wm.Region = e.Region
			wm.StreetAddress = e.StreetAddress
			wm.State = domain.AddressStateActive
		case *user.HumanAddressChangedEvent:
			if e.Country != nil {
				wm.Country = *e.Country
			}
			if e.Locality != nil {
				wm.Locality = *e.Locality
			}
			if e.PostalCode != nil {
				wm.PostalCode = *e.PostalCode
			}
			if e.Region != nil {
				wm.Region = *e.Region
			}
			if e.StreetAddress != nil {
				wm.StreetAddress = *e.StreetAddress
			}
		case *user.UserRemovedEvent:
			wm.State = domain.AddressStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanAddressWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.UserV1AddedType,
			user.UserV1RegisteredType,
			user.UserV1AddressChangedType,
			user.HumanAddedType,
			user.HumanRegisteredType,
			user.HumanAddressChangedType,
			user.UserRemovedType).
		Builder()
}

func (wm *HumanAddressWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	country,
	locality,
	postalCode,
	region,
	streetAddress string,
) (*user.HumanAddressChangedEvent, bool, error) {
	changes := make([]user.AddressChanges, 0)
	var err error

	if wm.Country != country {
		changes = append(changes, user.ChangeCountry(country))
	}
	if wm.Locality != locality {
		changes = append(changes, user.ChangeLocality(locality))
	}
	if wm.PostalCode != postalCode {
		changes = append(changes, user.ChangePostalCode(postalCode))
	}
	if wm.Region != region {
		changes = append(changes, user.ChangeRegion(region))
	}
	if wm.StreetAddress != streetAddress {
		changes = append(changes, user.ChangeStreetAddress(streetAddress))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := user.NewAddressChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
