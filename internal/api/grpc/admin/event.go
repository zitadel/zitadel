package admin

import (
	"context"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) ListEvents(ctx context.Context, in *admin_pb.ListEventsRequest) (*admin_pb.ListEventsResponse, error) {
	filter, err := eventRequestToFilter(ctx, in)
	if err != nil {
		return nil, err
	}
	events, err := s.query.SearchEvents(ctx, filter)
	if err != nil {
		return nil, err
	}
	return convertEventsToResponse(events)
}

func (s *Server) ListEventTypes(ctx context.Context, in *admin_pb.ListEventTypesRequest) (*admin_pb.ListEventTypesResponse, error) {
	return &admin_pb.ListEventTypesResponse{
		EventTypes: s.query.SearchEventTypes(ctx),
	},nil
}

func (s *Server) ListAggregateTypes(ctx context.Context, in *admin_pb.ListAggregateTypesRequest) (*admin_pb.ListAggregateTypesResponse, error) {
	return &admin_pb.ListAggregateTypesResponse{
		AggregateTypes: s.query.SearchAggregateTypes(ctx),
	},nil
}

func eventRequestToFilter(ctx context.Context, req *admin_pb.ListEventsRequest) (*eventstore.SearchQueryBuilder, error) {
	eventTypes := make([]eventstore.EventType, len(req.EventTypes))
	for i, eventType := range req.EventTypes {
		eventTypes[i] = eventstore.EventType(eventType)
	}

	builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderDesc().
		InstanceID(authz.GetInstance(ctx).InstanceID()).
		Limit(uint64(req.Limit)).
		ResourceOwner(req.ResourceOwner).
		EditorUser(req.EditorUserId).
		AddQuery().
		AggregateIDs(req.AggregateId).
		AggregateTypes(eventstore.AggregateType(req.AggregateType)).
		EventTypes(eventTypes...).
		CreationDateAfter(req.CreationDate.AsTime()).
		SequenceGreater(req.Sequence).
		Builder()

	if req.Asc {
		builder.OrderAsc()
	}

	return builder, nil
}

func convertEventsToResponse(events []*query.Event) (response *admin_pb.ListEventsResponse, err error) {
	response = &admin_pb.ListEventsResponse{
		Events: make([]*admin_pb.Event, len(events)),
	}

	for i, event := range events {
		response.Events[i], err = convertEvent(event)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

func convertEvent(event *query.Event) (*admin_pb.Event, error) {
	var payload *structpb.Struct
	if len(event.Payload) > 0 {
		payload = new(structpb.Struct)
		if err := payload.UnmarshalJSON(event.Payload); err != nil {
			return nil, errors.ThrowInternal(err, "ADMIN-eaimD", "Errors.Internal")
		}
	}
	return &admin_pb.Event{
		Editor: &admin_pb.EventEditor{
			UserId:      event.Editor.ID,
			DisplayName: event.Editor.DisplayName,
			Service:     event.Editor.Service,
		},
		Aggregate: &admin_pb.EventAggregate{
			Id:            event.Aggregate.ID,
			Type:          string(event.Aggregate.Type),
			ResourceOwner: event.Aggregate.ResourceOwner,
		},
		Sequence:     event.Sequence,
		CreationDate: timestamppb.New(event.CreationDate),
		Payload:      payload,
	}, nil
}