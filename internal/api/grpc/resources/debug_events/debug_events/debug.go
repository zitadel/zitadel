package debug_events

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	resource_object "github.com/zitadel/zitadel/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	debug_events "github.com/zitadel/zitadel/pkg/grpc/resources/debug_events/v3alpha"
)

func (s *Server) CreateDebugEvents(ctx context.Context, req *debug_events.CreateDebugEventsRequest) (_ *debug_events.CreateDebugEventsResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	details, err := s.command.CreateDebugEvents(ctx, debugEventsFromRequest(req))
	if err != nil {
		return nil, err
	}
	return &debug_events.CreateDebugEventsResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_INSTANCE, authz.GetInstance(ctx).InstanceID()),
	}, nil
}

func (s *Server) GetDebugEventsStateById(ctx context.Context, req *debug_events.GetDebugEventsStateByIdRequest) (_ *debug_events.GetDebugEventsStateByIdResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	state, err := s.query.GetDebugEventsStateByID(ctx, req.GetId(), req.GetTriggerBulk())
	if err != nil {
		return nil, err
	}

	return &debug_events.GetDebugEventsStateByIdResponse{
		State: eventsStateToPB(state),
	}, nil
}
func (s *Server) ListDebugEventsStates(ctx context.Context, req *debug_events.ListDebugEventsStatesRequest) (_ *debug_events.ListDebugEventsStatesResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	states, err := s.query.ListDebugEventsStates(ctx, req.GetTriggerBulk())
	if err != nil {
		return nil, err
	}

	return &debug_events.ListDebugEventsStatesResponse{
		States: eventStatesToPB(states),
	}, nil
}
