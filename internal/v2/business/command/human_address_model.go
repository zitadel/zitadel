package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

type HumanAddressWriteModel struct {
	eventstore.WriteModel

	Country       string
	Locality      string
	PostalCode    string
	Region        string
	StreetAddress string
}

func NewHumanAddressWriteModel(userID string) *HumanAddressWriteModel {
	return &HumanAddressWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: userID,
		},
	}
}

func (wm *HumanAddressWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanAddressChangedEvent:
			wm.AppendEvents(e)
			//TODO: Handle relevant User Events (remove, etc)

		}
	}
}

func (wm *HumanAddressWriteModel) Reduce() error {
	//TODO: implement
	return nil
}

func (wm *HumanAddressWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID)
}

func (wm *HumanAddressWriteModel) NewChangedEvent(
	ctx context.Context,
	country,
	locality,
	postalCode,
	region,
	streetAddress string,
) (*user.HumanAddressChangedEvent, bool) {
	hasChanged := false
	changedEvent := user.NewHumanAddressChangedEvent(ctx)
	if wm.Country != country {
		hasChanged = true
		changedEvent.Country = country
	}
	if wm.Locality != locality {
		hasChanged = true
		changedEvent.Locality = locality
	}
	if wm.PostalCode != postalCode {
		hasChanged = true
		changedEvent.PostalCode = postalCode
	}
	if wm.Region != region {
		hasChanged = true
		changedEvent.Region = region
	}
	if wm.StreetAddress != streetAddress {
		hasChanged = true
		changedEvent.StreetAddress = streetAddress
	}
	return changedEvent, hasChanged
}
