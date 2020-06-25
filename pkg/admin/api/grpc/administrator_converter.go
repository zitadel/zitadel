package grpc

import (
	view_model "github.com/caos/zitadel/internal/view/model"
)

func viewsFromModel(views []*view_model.View) []*View {
	result := make([]*View, len(views))
	for i, view := range views {
		result[i] = viewFromModel(view)
	}

	return result
}

func failedEventsFromModel(failedEvents []*view_model.FailedEvent) []*FailedEvent {
	result := make([]*FailedEvent, len(failedEvents))
	for i, view := range failedEvents {
		result[i] = failedEventFromModel(view)
	}

	return result
}

func viewFromModel(view *view_model.View) *View {
	return &View{
		Database: view.Database,
		ViewName: view.ViewName,
		Sequence: view.CurrentSequence,
	}
}

func failedEventFromModel(failedEvent *view_model.FailedEvent) *FailedEvent {
	return &FailedEvent{
		Database:       failedEvent.Database,
		ViewName:       failedEvent.ViewName,
		FailedSequence: failedEvent.FailedSequence,
		FailureCount:   failedEvent.FailureCount,
		ErrorMessage:   failedEvent.ErrMsg,
	}
}
