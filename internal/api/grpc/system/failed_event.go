package system

import (
	"context"

	"github.com/caos/zitadel/internal/query"
	system_pb "github.com/caos/zitadel/pkg/grpc/system"
)

func (s *Server) ListFailedEvents(ctx context.Context, req *system_pb.ListFailedEventsRequest) (*system_pb.ListFailedEventsResponse, error) {
	failedEventsOld, err := s.administrator.GetFailedEvents(ctx)
	if err != nil {
		return nil, err
	}
	convertedOld := FailedEventsViewToPb(failedEventsOld)

	failedEvents, err := s.query.SearchFailedEvents(ctx, new(query.FailedEventSearchQueries))
	if err != nil {
		return nil, err
	}
	convertedNew := FailedEventsToPb(failedEvents)
	convertedOld = append(convertedOld, convertedNew...)
	return &system_pb.ListFailedEventsResponse{Result: convertedOld}, nil
}

func (s *Server) RemoveFailedEvent(ctx context.Context, req *system_pb.RemoveFailedEventRequest) (*system_pb.RemoveFailedEventResponse, error) {
	var err error
	if req.Database != "zitadel" {
		err = s.administrator.RemoveFailedEvent(ctx, RemoveFailedEventRequestToModel(req))
	} else {
		err = s.query.RemoveFailedEvent(ctx, req.ViewName, req.FailedSequence)
	}
	if err != nil {
		return nil, err
	}
	return &system_pb.RemoveFailedEventResponse{}, nil
}
