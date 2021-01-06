package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
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

	UserState domain.UserState
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
		case *user.HumanAddedEvent, *user.HumanRegisteredEvent:
			wm.AppendEvents(e)
		case *user.UserRemovedEvent:
			wm.AppendEvents(e)
		}
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
			wm.UserState = domain.UserStateActive
		case *user.HumanRegisteredEvent:
			wm.Country = e.Country
			wm.Locality = e.Locality
			wm.PostalCode = e.PostalCode
			wm.Region = e.Region
			wm.StreetAddress = e.StreetAddress
			wm.UserState = domain.UserStateActive
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
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
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
