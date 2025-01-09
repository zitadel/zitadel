package samlrequest

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	AggregateType    = "saml_request"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(id, instanceID string) *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:          AggregateType,
			Version:       AggregateVersion,
			ID:            id,
			ResourceOwner: instanceID,
			InstanceID:    instanceID,
		},
	}
}
