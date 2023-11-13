package event

import (
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
	eventpb "github.com/zitadel/zitadel/pkg/grpc/event"
	"github.com/zitadel/zitadel/pkg/grpc/message"
)

func EventsToPb(events []*query.Event) (response []*eventpb.Event, err error) {
	response = make([]*eventpb.Event, len(events))

	for i, event := range events {
		response[i], err = EventToPb(event)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

func EventToPb(event *query.Event) (response *eventpb.Event, err error) {
	var payload *structpb.Struct
	if len(event.Payload) > 0 {
		payload = new(structpb.Struct)
		if err := payload.UnmarshalJSON(event.Payload); err != nil {
			return nil, errors.ThrowInternal(err, "ADMIN-eaimD", "Errors.Internal")
		}
	}
	return &eventpb.Event{
		Editor: &eventpb.Editor{
			UserId:      event.Editor.ID,
			DisplayName: event.Editor.DisplayName,
			Service:     event.Editor.Service,
		},
		Aggregate: &eventpb.Aggregate{
			Id:            event.Aggregate.ID,
			Type:          AggregateTypeToPb(string(event.Aggregate.Type)),
			ResourceOwner: event.Aggregate.ResourceOwner,
		},
		Sequence:     event.Sequence,
		CreationDate: timestamppb.New(event.CreationDate),
		Payload:      payload,
		Type:         EventTypeToPb(event.Type),
	}, nil
}

func EventTypeToPb(typ string) *eventpb.EventType {
	return &eventpb.EventType{
		Type:      typ,
		Localized: message.NewLocalizedEventType(typ),
	}
}

func AggregateTypeToPb(typ string) *eventpb.AggregateType {
	return &eventpb.AggregateType{
		Type:      typ,
		Localized: message.NewLocalizedAggregateType(typ),
	}
}
