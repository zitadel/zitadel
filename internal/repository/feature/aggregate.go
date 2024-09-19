package feature

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

const (
	eventTypePrefix = eventstore.EventType("feature.")
	setSuffix       = ".set"
)

const (
	AggregateType    = "feature"
	AggregateVersion = "v1"
)
