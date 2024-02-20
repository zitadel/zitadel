package execution

import (
	"strings"

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

const (
	grpcPrefix     = "grpc"
	eventPrefix    = "event"
	functionPrefix = "func"
)

func IDFromGRPC(method string) string {
	return strings.Join([]string{grpcPrefix, method}, ".")
}

func IDFromGRPCAll() string {
	return grpcPrefix
}

func IDFromEvent(event string) string {
	return strings.Join([]string{eventPrefix, event}, ".")
}

func IDFromEventAll() string {
	return eventPrefix
}

func IDFromFunction(name string) string {
	return strings.Join([]string{functionPrefix, name}, ".")
}
