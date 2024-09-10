package keypair

import "github.com/zitadel/zitadel/internal/repository/keypair"

const (
	AggregateType   = string(keypair.AggregateType)
	eventTypePrefix = AggregateType + "."
)
