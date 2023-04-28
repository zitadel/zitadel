package system

import (
	"context"

	"github.com/zitadel/zitadel/internal/query"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func (s *Server) ListFailedEvents(ctx context.Context, _ *system_pb.ListFailedEventsRequest) (*system_pb.ListFailedEventsResponse, error) {
	failedEvents, err := s.query.SearchFailedEvents(ctx, new(query.FailedEventSearchQueries))
	if err != nil {
		return nil, err
	}
	return &system_pb.ListFailedEventsResponse{Result: FailedEventsToPb(s.database, failedEvents)}, nil
}

func (s *Server) RemoveFailedEvent(ctx context.Context, req *system_pb.RemoveFailedEventRequest) (*system_pb.RemoveFailedEventResponse, error) {
	err := s.query.RemoveFailedEvent(ctx, req.ViewName, req.InstanceId, req.FailedSequence)
	if err != nil {
		return nil, err
	}
	return &system_pb.RemoveFailedEventResponse{}, nil
}
