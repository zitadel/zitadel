package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
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
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(user.HumanAddedType,
			user.HumanRegisteredType,
			user.HumanAddressChangedType,
			user.UserRemovedType)
}

func (wm *HumanAddressWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	country,
	locality,
	postalCode,
	region,
	streetAddress string,
) (*user.HumanAddressChangedEvent, bool) {
	hasChanged := false
	changedEvent := user.NewHumanAddressChangedEvent(ctx, aggregate)
	if wm.Country != country {
		hasChanged = true
		changedEvent.Country = &country
	}
	if wm.Locality != locality {
		hasChanged = true
		changedEvent.Locality = &locality
	}
	if wm.PostalCode != postalCode {
		hasChanged = true
		changedEvent.PostalCode = &postalCode
	}
	if wm.Region != region {
		hasChanged = true
		changedEvent.Region = &region
	}
	if wm.StreetAddress != streetAddress {
		hasChanged = true
		changedEvent.StreetAddress = &streetAddress
	}
	return changedEvent, hasChanged
}
