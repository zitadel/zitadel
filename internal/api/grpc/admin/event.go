package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object"
)

func (s *Server) ListEvents(ctx context.Context, in *admin_pb.ListEventsRequest) (*admin_pb.ListEventsResponse, error) {
	filter, err := eventRequestToFilter(in)
	events, err := s.query.SearchEvents(ctx, filter)
	if err != nil {
		return nil, err
	}
	return convertEventsToResponse(events), nil
}

func eventRequestToFilter(req *admin_pb.ListEventsRequest) (*)

func convertEventsToResponse(events []eventstore.Event) *admin_pb.ListEventsResponse {

	return &admin_pb.ListEventsResponse{
		Details: &object_pb.ListDetails{

		},
	}
}