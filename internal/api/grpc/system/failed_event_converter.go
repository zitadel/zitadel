package system

import (
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/view/model"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func FailedEventsViewToPb(failedEvents []*model.FailedEvent) []*system_pb.FailedEvent {
	events := make([]*system_pb.FailedEvent, len(failedEvents))
	for i, failedEvent := range failedEvents {
		events[i] = FailedEventViewToPb(failedEvent)
	}
	return events
}

func FailedEventViewToPb(failedEvent *model.FailedEvent) *system_pb.FailedEvent {
	return &system_pb.FailedEvent{
		Database:       failedEvent.Database,
		ViewName:       failedEvent.ViewName,
		FailedSequence: failedEvent.FailedSequence,
		FailureCount:   failedEvent.FailureCount,
		ErrorMessage:   failedEvent.ErrMsg,
	}
}

func FailedEventsToPb(failedEvents *query.FailedEvents) []*system_pb.FailedEvent {
	events := make([]*system_pb.FailedEvent, len(failedEvents.FailedEvents))
	for i, failedEvent := range failedEvents.FailedEvents {
		events[i] = FailedEventToPb(failedEvent)
	}
	return events
}

func FailedEventToPb(failedEvent *query.FailedEvent) *system_pb.FailedEvent {
	return &system_pb.FailedEvent{
		Database:       "zitadel",
		ViewName:       failedEvent.ProjectionName,
		FailedSequence: failedEvent.FailedSequence,
		FailureCount:   failedEvent.FailureCount,
		ErrorMessage:   failedEvent.Error,
	}
}

func RemoveFailedEventRequestToModel(req *system_pb.RemoveFailedEventRequest) *model.FailedEvent {
	return &model.FailedEvent{
		Database:       req.Database,
		ViewName:       req.ViewName,
		FailedSequence: req.FailedSequence,
	}
}
