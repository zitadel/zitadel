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

func NewHumanAddressChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanAddressChangedEvent {
	return &HumanAddressChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanAddressChangedType,
		),
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
