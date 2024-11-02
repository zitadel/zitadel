package instance

import "github.com/zitadel/zitadel/internal/repository/instance"

const (
	AggregateType   = string(instance.AggregateType)
	eventTypePrefix = AggregateType + "."
)
