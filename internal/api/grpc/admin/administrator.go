package admin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	view_model "github.com/caos/zitadel/internal/view/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetViews(ctx context.Context, _ *empty.Empty) (_ *admin.Views, err error) {
	views, err := s.administrator.GetViews()
	if err != nil {
		return nil, err
	}
	return &admin.Views{Views: viewsFromModel(views)}, nil
}

func (s *Server) ClearView(ctx context.Context, viewID *admin.ViewID) (_ *empty.Empty, err error) {
	err = s.administrator.ClearView(ctx, viewID.Database, viewID.ViewName)
	return &empty.Empty{}, err
}

func (s *Server) GetFailedEvents(ctx context.Context, _ *empty.Empty) (_ *admin.FailedEvents, err error) {
	failedEvents, err := s.administrator.GetFailedEvents(ctx)
	if err != nil {
		return nil, err
	}
	return &admin.FailedEvents{FailedEvents: failedEventsFromModel(failedEvents)}, nil
}

func (s *Server) RemoveFailedEvent(ctx context.Context, failedEventID *admin.FailedEventID) (_ *empty.Empty, err error) {
	err = s.administrator.RemoveFailedEvent(ctx, &view_model.FailedEvent{Database: failedEventID.Database, ViewName: failedEventID.ViewName, FailedSequence: failedEventID.FailedSequence})
	return &empty.Empty{}, err
}
