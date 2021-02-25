package admin

import (
	"context"

	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) ListFailedEvents(ctx context.Context, req *admin_pb.ListFailedEventsRequest) (*admin_pb.ListFailedEventsResponse, error) {
	failedEvents, err := s.administrator.GetFailedEvents(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListFailedEventsResponse{Result: FailedEventsToPb(failedEvents)}, nil
}

func (s *Server) RemoveFailedEvent(ctx context.Context, req *admin_pb.RemoveFailedEventRequest) (*admin_pb.RemoveFailedEventResponse, error) {
	err := s.administrator.RemoveFailedEvent(ctx, RemoveFailedEventRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveFailedEventResponse{}, nil
}
