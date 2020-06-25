package grpc

import (
	"context"
	view_model "github.com/caos/zitadel/internal/view/model"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetViews(ctx context.Context, _ *empty.Empty) (_ *Views, err error) {
	views, err := s.administrator.GetViews(ctx)
	if err != nil {
		return nil, err
	}
	return &Views{Views: viewsFromModel(views)}, nil
}

func (s *Server) ClearView(ctx context.Context, viewID *ViewID) (_ *empty.Empty, err error) {
	err = s.administrator.ClearView(ctx, viewID.Database, viewID.ViewName)
	return &empty.Empty{}, err
}

func (s *Server) GetFailedEvents(ctx context.Context, _ *empty.Empty) (_ *FailedEvents, err error) {
	failedEvents, err := s.administrator.GetFailedEvents(ctx)
	if err != nil {
		return nil, err
	}
	return &FailedEvents{FailedEvents: failedEventsFromModel(failedEvents)}, nil
}

func (s *Server) RemoveFailedEvent(ctx context.Context, failedEventID *FailedEventID) (_ *empty.Empty, err error) {
	err = s.administrator.RemoveFailedEvent(ctx, &view_model.FailedEvent{Database: failedEventID.Database, ViewName: failedEventID.ViewName, FailedSequence: failedEventID.FailedSequence})
	return &empty.Empty{}, err
}
