package admin

import (
	event_grpc "github.com/zitadel/zitadel/internal/api/grpc/event"
	"github.com/zitadel/zitadel/internal/api/grpc/server/middleware"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/message"
)

func EventTypesToPb(eventTypes []string) *ListEventTypesResponse {
	res := &ListEventTypesResponse{EventTypes: make([]*message.LocalizedMessage, len(eventTypes))}

	for i, eventType := range eventTypes {
		res.EventTypes[i] = message.NewLocalizedEventType(eventType)
	}

	return res
}

func AggregateTypesToPb(eventTypes []string) *ListAggregateTypesResponse {
	res := &ListAggregateTypesResponse{AggregateTypes: make([]*message.LocalizedMessage, len(eventTypes))}

	for i, eventType := range eventTypes {
		res.AggregateTypes[i] = message.NewLocalizedAggregateType(eventType)
	}

	return res
}

func EventsToPb(events []*query.Event) (*ListEventsResponse, error) {
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
		localizers[i] = eventType
	}
	return localizers
}

func (resp *ListAggregateTypesResponse) Localizers() []middleware.Localizer {
	if resp == nil {
		return nil
	}

	localizers := make([]middleware.Localizer, len(resp.AggregateTypes))
	for i, eventType := range resp.AggregateTypes {
		localizers[i] = eventType
	}
	return localizers
}

func (resp *ListEventsResponse) Localizers() []middleware.Localizer {
	if resp == nil {
		return nil
	}

	localizers := make([]middleware.Localizer, 0, len(resp.Events)*2)
	for _, event := range resp.Events {
		localizers = append(localizers, event.Type, event.Aggregate.Type)
	}
	return localizers
}
