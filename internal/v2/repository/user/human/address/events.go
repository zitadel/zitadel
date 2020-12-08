package address

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	addressEventPrefix      = eventstore.EventType("user.human.address.")
	HumanAddressChangedType = addressEventPrefix + "changed"
)

type HumanAddressChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Country       string `json:"country,omitempty"`
	Locality      string `json:"locality,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Region        string `json:"region,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`
}

func (e *HumanAddressChangedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanAddressChangedEvent) Data() interface{} {
	return e
}

func NewHumanAddressChangedEvent(
	ctx context.Context,
	current *HumanAddressWriteModel,
	country,
	locality,
	postalCode,
	region,
	streetAddress string,
) *HumanAddressChangedEvent {
	e := &HumanAddressChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanAddressChangedType,
		),
	}

	if current.Country != country {
		e.Country = country
	}
	if current.Locality != locality {
		e.Locality = locality
	}
	if current.PostalCode != postalCode {
		e.PostalCode = postalCode
	}
	if current.Region != region {
		e.Region = region
	}
	if current.StreetAddress != streetAddress {
		e.StreetAddress = streetAddress
	}

	return e
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
