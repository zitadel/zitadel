package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
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

func (e *HumanAddressChangedEvent) Payload() interface{} {
	return e
}

func (e *HumanAddressChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewAddressChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []AddressChanges,
) (*HumanAddressChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "USER-3n8fs", "Errors.NoChangesFound")
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

func HumanAddressChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	addressChanged := &HumanAddressChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(addressChanged)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-5M0pd", "unable to unmarshal human address changed")
	}

	return addressChanged, nil
}
