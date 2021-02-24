package admin

import (
	"github.com/caos/logging"
	view_model "github.com/caos/zitadel/internal/view/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
)

func viewsFromModel(views []*view_model.View) []*admin.View {
	result := make([]*admin.View, len(views))
	for i, view := range views {
		result[i] = viewFromModel(view)
	}

	return result
}

func failedEventsFromModel(failedEvents []*view_model.FailedEvent) []*admin.FailedEvent {
	result := make([]*admin.FailedEvent, len(failedEvents))
	for i, view := range failedEvents {
		result[i] = failedEventFromModel(view)
	}

	return result
}

func viewFromModel(view *view_model.View) *admin.View {
	eventTimestamp, err := ptypes.TimestampProto(view.EventTimestamp)
	logging.Log("GRPC-KSo03").OnError(err).Debug("unable to parse timestamp")
	lastSpool, err := ptypes.TimestampProto(view.LastSuccessfulSpoolerRun)
	logging.Log("GRPC-0oP87").OnError(err).Debug("unable to parse timestamp")

	return &admin.View{
		Database:                 view.Database,
		ViewName:                 view.ViewName,
		ProcessedSequence:        view.CurrentSequence,
		EventTimestamp:           eventTimestamp,
		LastSuccessfulSpoolerRun: lastSpool,
	}
}

func failedEventFromModel(failedEvent *view_model.FailedEvent) *admin.FailedEvent {
	return &admin.FailedEvent{
		Database:       failedEvent.Database,
		ViewName:       failedEvent.ViewName,
		FailedSequence: failedEvent.FailedSequence,
		FailureCount:   failedEvent.FailureCount,
		ErrorMessage:   failedEvent.ErrMsg,
	}
}
