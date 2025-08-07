package admin

import (
	"context"

	event_grpc "github.com/zitadel/zitadel/internal/api/grpc/event"
	"github.com/zitadel/zitadel/internal/api/grpc/server/middleware"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	event_pb "github.com/zitadel/zitadel/pkg/grpc/event"
)

func EventTypesToPb(eventTypes []string) *ListEventTypesResponse {
	res := &ListEventTypesResponse{EventTypes: make([]*event_pb.EventType, len(eventTypes))}

	for i, eventType := range eventTypes {
		res.EventTypes[i] = event_grpc.EventTypeToPb(eventType)
	}

	return res
}

func AggregateTypesToPb(aggregateTypes []string) *ListAggregateTypesResponse {
	res := &ListAggregateTypesResponse{AggregateTypes: make([]*event_pb.AggregateType, len(aggregateTypes))}

	for i, aggregateType := range aggregateTypes {
		res.AggregateTypes[i] = event_grpc.AggregateTypeToPb(aggregateType)
	}

	return res
}

func EventsToPb(ctx context.Context, events []*query.Event) (_ *ListEventsResponse, err error) {
	_, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	res, err := event_grpc.EventsToPb(events)
	if err != nil {
		return nil, err
	}
	return &ListEventsResponse{
		Events: res,
	}, nil
}

func (resp *ListEventTypesResponse) Localizers() []middleware.Localizer {
	if resp == nil {
		return nil
	}

	localizers := make([]middleware.Localizer, len(resp.EventTypes))
	for i, eventType := range resp.EventTypes {
		localizers[i] = eventType.Localized
	}
	return localizers
}

func (resp *ListAggregateTypesResponse) Localizers() []middleware.Localizer {
	if resp == nil {
		return nil
	}

	localizers := make([]middleware.Localizer, len(resp.AggregateTypes))
	for i, aggregateType := range resp.AggregateTypes {
		localizers[i] = aggregateType.Localized
	}
	return localizers
}

func (resp *ListEventsResponse) Localizers() []middleware.Localizer {
	if resp == nil {
		return nil
	}

	localizers := make([]middleware.Localizer, 0, len(resp.Events)*2)
	for _, event := range resp.Events {
		localizers = append(localizers, event.Type.Localized, event.Aggregate.Type.Localized)
	}
	return localizers
}
