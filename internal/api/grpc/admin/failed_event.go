package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) ListFailedEvents(ctx context.Context, _ *admin_pb.ListFailedEventsRequest) (*admin_pb.ListFailedEventsResponse, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	failedEventsOld, err := s.administrator.GetFailedEvents(ctx, instanceID)
	if err != nil {
		return nil, err
	}
	convertedOld := FailedEventsViewToPb(failedEventsOld)
	instanceIDQuery, err := query.NewFailedEventInstanceIDSearchQuery(instanceID)
	if err != nil {
		return nil, err
	}
	failedEvents, err := s.query.SearchFailedEvents(ctx, &query.FailedEventSearchQueries{
		Queries: []query.SearchQuery{instanceIDQuery},
	})
	if err != nil {
		return nil, err
	}
	convertedNew := FailedEventsToPb(s.database, failedEvents)
	return &admin_pb.ListFailedEventsResponse{Result: append(convertedOld, convertedNew...)}, nil
}

func (s *Server) RemoveFailedEvent(ctx context.Context, req *admin_pb.RemoveFailedEventRequest) (*admin_pb.RemoveFailedEventResponse, error) {
	var err error
	if req.Database != s.database {
		err = s.administrator.RemoveFailedEvent(ctx, RemoveFailedEventRequestToModel(ctx, req))
	} else {
		err = s.query.RemoveFailedEvent(ctx, req.ViewName, authz.GetInstance(ctx).InstanceID(), req.FailedSequence)
	}
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveFailedEventResponse{}, nil
}
