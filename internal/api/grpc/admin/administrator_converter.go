package admin

import (
	view_model "github.com/caos/zitadel/internal/view/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
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
	return &admin.View{
		Database: view.Database,
		ViewName: view.ViewName,
		Sequence: view.CurrentSequence,
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
