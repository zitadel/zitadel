package model

import "github.com/zitadel/zitadel/internal/eventstore/v1/models"

const (
	KeyPairAggregate models.AggregateType = "key_pair"

	KeyPairAdded models.EventType = "key_pair.added"
)
