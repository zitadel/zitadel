package instance

import "github.com/zitadel/zitadel/v2/internal/repository/instance"

const (
	AggregateType   = string(instance.AggregateType)
	eventTypePrefix = AggregateType + "."
)
