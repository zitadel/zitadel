package group

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	groupEventTypePrefix = eventstore.EventType("group.")

	AggregateType    = "group"
	AggregateVersion = "v1" // todo: review
)
