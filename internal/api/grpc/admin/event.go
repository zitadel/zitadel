package admin

import (
	"context"
	"slices"
	"time"

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
	var fromTime, sinceTime, untilTime time.Time
	// We ignore the deprecation warning here because we still need to support the deprecated field.
	//nolint:staticcheck
	if creationDatePb := req.GetCreationDate(); creationDatePb != nil {
		fromTime = creationDatePb.AsTime()
	}
	if fromTimePb := req.GetFrom(); fromTimePb != nil {
		fromTime = fromTimePb.AsTime()
	}
	if timeRange := req.GetRange(); timeRange != nil {
		// If range is set, we ignore the from and the deprecated creation_date fields
		fromTime = time.Time{}
		if timeSincePb := timeRange.GetSince(); timeSincePb != nil {
			sinceTime = timeSincePb.AsTime()
		}
		if timeUntilPb := timeRange.GetUntil(); timeUntilPb != nil {
			untilTime = timeUntilPb.AsTime()
		}
	}
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
	if len(aggregateTypes) == 0 {
		aggregateTypes = aggregateTypesFromEventTypes(eventTypes)
	}
	aggregateTypes = slices.Compact(aggregateTypes)

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
		SequenceGreater(req.Sequence).
		CreationDateAfter(sinceTime).
		CreationDateBefore(untilTime)

	if len(aggregateIDs) > 0 || len(aggregateTypes) > 0 || len(eventTypes) > 0 {
		builder.AddQuery().
			AggregateIDs(aggregateIDs...).
			AggregateTypes(aggregateTypes...).
			EventTypes(eventTypes...).
			Builder()
	}

	if req.GetAsc() {
		builder.OrderAsc()
		builder.CreationDateAfter(fromTime)
	} else {
		builder.CreationDateBefore(fromTime)
	}
	return builder, nil
}

func aggregateTypesFromEventTypes(eventTypes []eventstore.EventType) []eventstore.AggregateType {
	aggregateTypes := make([]eventstore.AggregateType, 0, len(eventTypes))

	for _, eventType := range eventTypes {
		aggregateTypes = append(aggregateTypes, eventstore.AggregateTypeFromEventType(eventType))
	}

	return aggregateTypes
}
