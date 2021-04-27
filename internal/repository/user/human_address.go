package user

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	addressEventPrefix      = humanEventPrefix + "address."
	HumanAddressChangedType = addressEventPrefix + "changed"
)

type HumanAddressChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Country       *string `json:"country,omitempty"`
	Locality      *string `json:"locality,omitempty"`
	PostalCode    *string `json:"postalCode,omitempty"`
	Region        *string `json:"region,omitempty"`
	StreetAddress *string `json:"streetAddress,omitempty"`
}

func (e *HumanAddressChangedEvent) Data() interface{} {
	return e
}

func (e *HumanAddressChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *HumanAddressChangedEvent) Assets() []*eventstore.Asset {
	return nil
}

func NewAddressChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []AddressChanges,
) (*HumanAddressChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "USER-3n8fs", "Errors.NoChangesFound")
	}
	changeEvent := &HumanAddressChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanAddressChangedType,
		),
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type AddressChanges func(event *HumanAddressChangedEvent)

func ChangeCountry(country string) func(event *HumanAddressChangedEvent) {
	return func(e *HumanAddressChangedEvent) {
		e.Country = &country
	}
}

func ChangeLocality(locality string) func(event *HumanAddressChangedEvent) {
	return func(e *HumanAddressChangedEvent) {
		e.Locality = &locality
	}
}

func ChangePostalCode(code string) func(event *HumanAddressChangedEvent) {
	return func(e *HumanAddressChangedEvent) {
		e.PostalCode = &code
	}
}

func ChangeRegion(region string) func(event *HumanAddressChangedEvent) {
	return func(e *HumanAddressChangedEvent) {
		e.Region = &region
	}
}

func ChangeStreetAddress(street string) func(event *HumanAddressChangedEvent) {
	return func(e *HumanAddressChangedEvent) {
		e.StreetAddress = &street
	}
}

func HumanAddressChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	addressChanged := &HumanAddressChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, addressChanged)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5M0pd", "unable to unmarshal human address changed")
	}

	return addressChanged, nil
}
