package admin

import (
	"github.com/caos/zitadel/internal/view/model"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func FailedEventsToPb(failedEvents []*model.FailedEvent) []*admin_pb.FailedEvent {
	events := make([]*admin_pb.FailedEvent, len(failedEvents))
	for i, failedEvent := range failedEvents {
		events[i] = FailedEventToPb(failedEvent)
	}
	return events
}

func FailedEventToPb(failedEvent *model.FailedEvent) *admin_pb.FailedEvent {
	return &admin_pb.FailedEvent{
		Database:       failedEvent.Database,
		ViewName:       failedEvent.ViewName,
		FailedSequence: failedEvent.FailedSequence,
		FailureCount:   failedEvent.FailureCount,
		ErrorMessage:   failedEvent.ErrMsg,
	}
}

func RemoveFailedEventRequestToModel(req *admin_pb.RemoveFailedEventRequest) *model.FailedEvent {
	return &model.FailedEvent{
		Database:       req.Database,
		ViewName:       req.ViewName,
		FailedSequence: req.FailedSequence,
	}
}
