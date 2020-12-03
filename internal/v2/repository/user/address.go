package user

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	addressEventPrefix      = humanEventPrefix + "address."
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
	return false
}

func (e *HumanAddressChangedEvent) Data() interface{} {
	return e
}

func NewHumanAddressChangedEvent(
	base *eventstore.BaseEvent,
	country,
	locality,
	postalCode,
	region,
	streetAddress string) *HumanAddressChangedEvent {
	return &HumanAddressChangedEvent{
		BaseEvent:     *base,
		Country:       country,
		Locality:      locality,
		PostalCode:    postalCode,
		Region:        region,
		StreetAddress: streetAddress,
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
