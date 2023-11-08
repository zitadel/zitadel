package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

const (
	maxLimit = 1000
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

	return admin_pb.EventsToPb(ctx, events)
}

func (s *Server) ListEventTypes(ctx context.Context, in *admin_pb.ListEventTypesRequest) (*admin_pb.ListEventTypesResponse, error) {
	eventTypes := s.query.SearchEventTypes(ctx)
	return admin_pb.EventTypesToPb(eventTypes), nil
}

func (s *Server) ListAggregateTypes(ctx context.Context, in *admin_pb.ListAggregateTypesRequest) (*admin_pb.ListAggregateTypesResponse, error) {
	aggregateTypes := s.query.SearchAggregateTypes(ctx)
	return admin_pb.AggregateTypesToPb(aggregateTypes), nil
}

func eventRequestToFilter(ctx context.Context, req *admin_pb.ListEventsRequest) (*eventstore.SearchQueryBuilder, error) {
	eventTypes := make([]eventstore.EventType, len(req.EventTypes))
	for i, eventType := range req.EventTypes {
		eventTypes[i] = eventstore.EventType(eventType)
	}
	aggregateIDs := make([]string, 0, 1)
	if req.AggregateId != "" {
		aggregateIDs = append(aggregateIDs, req.AggregateId)
	}
	aggregateTypes := make([]eventstore.AggregateType, len(req.AggregateTypes))
	for i, aggregateType := range req.AggregateTypes {
		aggregateTypes[i] = eventstore.AggregateType(aggregateType)
	}
	limit := uint64(req.Limit)
	if limit == 0 || limit > maxLimit {
		limit = maxLimit
	}

	builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderDesc().
		InstanceID(authz.GetInstance(ctx).InstanceID()).
		Limit(limit).
		AwaitOpenTransactions().
		ResourceOwner(req.ResourceOwner).
		EditorUser(req.EditorUserId).
		SequenceGreater(req.Sequence)

	if len(aggregateIDs) > 0 || len(aggregateTypes) > 0 || len(eventTypes) > 0 {
		builder.AddQuery().
			AggregateIDs(aggregateIDs...).
			AggregateTypes(aggregateTypes...).
			EventTypes(eventTypes...).
			Builder()
	}

	if req.GetAsc() {
		builder.OrderAsc()
		builder.CreationDateAfter(req.CreationDate.AsTime())
	} else {
		builder.CreationDateBefore(req.CreationDate.AsTime())
	}

	return builder, nil
}
