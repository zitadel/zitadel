package feature

import (
	v2 "github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
)

const (
	AggregateType   = string(v2.AggregateType)
	eventTypePrefix = AggregateType + "."
)
