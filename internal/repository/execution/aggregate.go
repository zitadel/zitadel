package execution

import (
	"strings"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	AggregateType    = "execution"
	AggregateVersion = "v1"
)

func NewAggregate(aggrID, instanceID string) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            aggrID,
		Type:          AggregateType,
		ResourceOwner: instanceID,
		InstanceID:    instanceID,
		Version:       AggregateVersion,
	}
}

func ID(executionType domain.ExecutionType, value string) string {
	return strings.Join([]string{executionType.String(), value}, ".")
}

func IDAll(executionType domain.ExecutionType) string {
	return executionType.String()
}
