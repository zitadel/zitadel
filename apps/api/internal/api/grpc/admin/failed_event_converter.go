package admin

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func FailedEventsToPb(database string, failedEvents *query.FailedEvents) []*admin_pb.FailedEvent {
	events := make([]*admin_pb.FailedEvent, len(failedEvents.FailedEvents))
	for i, failedEvent := range failedEvents.FailedEvents {
		events[i] = FailedEventToPb(database, failedEvent)
	}
	return events
}

func FailedEventToPb(database string, failedEvent *query.FailedEvent) *admin_pb.FailedEvent {
	var lastFailed *timestamppb.Timestamp
	if !failedEvent.LastFailed.IsZero() {
		lastFailed = timestamppb.New(failedEvent.LastFailed)
	}
	return &admin_pb.FailedEvent{
		Database:       database,
		ViewName:       failedEvent.ProjectionName,
		FailedSequence: failedEvent.FailedSequence,
		FailureCount:   failedEvent.FailureCount,
		ErrorMessage:   failedEvent.Error,
		LastFailed:     lastFailed,
	}
}
