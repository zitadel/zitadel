package admin

import (
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/view/model"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func FailedEventsViewToPb(failedEvents []*model.FailedEvent) []*admin_pb.FailedEvent {
	events := make([]*admin_pb.FailedEvent, len(failedEvents))
	for i, failedEvent := range failedEvents {
		events[i] = FailedEventViewToPb(failedEvent)
	}
	return events
}

func FailedEventViewToPb(failedEvent *model.FailedEvent) *admin_pb.FailedEvent {
	return &admin_pb.FailedEvent{
		Database:       failedEvent.Database,
		ViewName:       failedEvent.ViewName,
		FailedSequence: failedEvent.FailedSequence,
		FailureCount:   failedEvent.FailureCount,
		ErrorMessage:   failedEvent.ErrMsg,
	}
}

func FailedEventsToPb(failedEvents *query.FailedEvents) []*admin_pb.FailedEvent {
	events := make([]*admin_pb.FailedEvent, len(failedEvents.FailedEvents))
	for i, failedEvent := range failedEvents.FailedEvents {
		events[i] = FailedEventToPb(failedEvent)
	}
	return events
}

func FailedEventToPb(failedEvent *query.FailedEvent) *admin_pb.FailedEvent {
	return &admin_pb.FailedEvent{
		Database:       "zitadel",
		ViewName:       failedEvent.ProjectionName,
		FailedSequence: failedEvent.FailedSequence,
		FailureCount:   failedEvent.FailureCount,
		ErrorMessage:   failedEvent.Error,
	}
}

func RemoveFailedEventRequestToModel(req *admin_pb.RemoveFailedEventRequest) *model.FailedEvent {
	return &model.FailedEvent{
		Database:       req.Database,
		ViewName:       req.ViewName,
		FailedSequence: req.FailedSequence,
	}
}
