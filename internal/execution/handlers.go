package execution

import (
	"context"
	_ "embed"
	"encoding/json"
	"strings"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
)

const (
	ExecutionHandlerTable     = "projections.execution_handler"
	ExecutionInstanceID       = "instance_id"
	ExecutionResourceOwner    = "resource_owner"
	ExecutionAggregateType    = "aggregate_type"
	ExecutionAggregateVersion = "aggregate_version"
	ExecutionAggregateID      = "aggregate_id"
	ExecutionSequence         = "sequence"
	ExecutionEventType        = "event_type"
	ExecutionCreatedAt        = "create_at"
	ExecutionPosition         = "position"
	ExecutionEventUserIDCol   = "user_id"
	ExecutionEventDataCol     = "event_data"
	ExecutionTargetsDataCol   = "targets_data"
)

type executionsHandler struct {
	es    *eventstore.Eventstore
	query Queries
}

func NewExecutionsHandler(
	ctx context.Context,
	config handler.Config,
	es *eventstore.Eventstore,
	query Queries,
) *handler.Handler {
	return handler.NewHandler(ctx, &config, &executionsHandler{es: es, query: query})
}

func (u *executionsHandler) Name() string {
	return ExecutionHandlerTable
}

func (*executionsHandler) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(ExecutionInstanceID, handler.ColumnTypeText),
			handler.NewColumn(ExecutionResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(ExecutionAggregateType, handler.ColumnTypeText),
			handler.NewColumn(ExecutionAggregateVersion, handler.ColumnTypeText),
			handler.NewColumn(ExecutionAggregateID, handler.ColumnTypeText),
			handler.NewColumn(ExecutionSequence, handler.ColumnTypeInt64),
			handler.NewColumn(ExecutionCreatedAt, handler.ColumnTypeText),
			handler.NewColumn(ExecutionPosition, handler.ColumnTypeText),
			handler.NewColumn(ExecutionEventType, handler.ColumnTypeText),
			handler.NewColumn(ExecutionEventUserIDCol, handler.ColumnTypeText),
			handler.NewColumn(ExecutionEventDataCol, handler.ColumnTypeJSONB),
			handler.NewColumn(ExecutionTargetsDataCol, handler.ColumnTypeJSONB),
		},
			handler.NewPrimaryKey(ExecutionInstanceID, ExecutionResourceOwner, ExecutionAggregateID, ExecutionPosition),
		),
	)
}

func (u *executionsHandler) Reducers() []handler.AggregateReducer {
	eventTypes := u.es.EventTypes()

	aggList := make(map[eventstore.AggregateType][]eventstore.EventType)
	for _, eventType := range eventTypes {
		aggType := eventstore.AggregateTypeFromEventType(eventstore.EventType(eventType))
		agg, ok := aggList[aggType]
		if !ok {
			aggList[aggType] = []eventstore.EventType{eventstore.EventType(eventType)}
		} else {
			found := false
			for _, aggEventType := range agg {
				if aggEventType == eventstore.EventType(eventType) {
					found = true
				}
			}
			if !found {
				agg = append(agg, eventstore.EventType(eventType))
			}
		}
	}

	aggReducers := make([]handler.AggregateReducer, len(aggList))
	for aggType, agg := range aggList {
		eventReducers := make([]handler.EventReducer, len(agg))
		for _, eventType := range agg {
			eventReducers = append(eventReducers,
				handler.EventReducer{
					Event:  eventType,
					Reduce: u.reduce,
				},
			)
		}
		aggReducers = append(aggReducers,
			handler.AggregateReducer{
				Aggregate:     aggType,
				EventReducers: eventReducers,
			},
		)
	}
	return aggReducers
}

func groupFromEventType(s string) string {
	parts := strings.Split(s, "/")
	return parts[1]
}

func idsForEventType(eventType string) []string {
	return []string{exec_repo.ID(domain.ExecutionTypeEvent, eventType), exec_repo.ID(domain.ExecutionTypeEvent, groupFromEventType(eventType)), exec_repo.IDAll(domain.ExecutionTypeEvent)}
}

func (u *executionsHandler) reduce(e eventstore.Event) (*handler.Statement, error) {
	ctx := HandlerContext(e.Aggregate())

	targets, err := u.query.TargetsByExecutionID(ctx, idsForEventType(string(e.Type())))
	if err != nil {
		return nil, err
	}

	ee, err := NewEventExecution(e, targets)
	if err != nil {
		return nil, err
	}

	return handler.NewCreateStatement(
		e,
		ee.Columns(),
	), nil
}

type EventExecution struct {
	Aggregate   *eventstore.Aggregate
	Sequence    uint64
	EventType   eventstore.EventType
	CreatedAt   time.Time
	Position    float64
	UserID      string
	EventData   []byte
	TargetsData []byte
}

func NewEventExecution(e eventstore.Event, targets []*query.ExecutionTarget) (*EventExecution, error) {
	targetsData, err := json.Marshal(targets)
	if err != nil {
		return nil, err
	}
	return &EventExecution{
		Aggregate:   e.Aggregate(),
		Sequence:    e.Sequence(),
		EventType:   e.Type(),
		CreatedAt:   e.CreatedAt(),
		Position:    e.Position(),
		UserID:      e.Creator(),
		EventData:   e.DataAsBytes(),
		TargetsData: targetsData,
	}, nil
}

func (e *EventExecution) Columns() []handler.Column {
	return []handler.Column{
		handler.NewCol(ExecutionInstanceID, e.Aggregate.InstanceID),
		handler.NewCol(ExecutionResourceOwner, e.Aggregate.ResourceOwner),
		handler.NewCol(ExecutionAggregateType, e.Aggregate.Type),
		handler.NewCol(ExecutionAggregateVersion, e.Aggregate.Version),
		handler.NewCol(ExecutionAggregateID, e.Aggregate.ID),
		handler.NewCol(ExecutionSequence, e.Sequence),
		handler.NewCol(ExecutionEventType, e.EventType),
		handler.NewCol(ExecutionCreatedAt, e.CreatedAt),
		handler.NewCol(ExecutionPosition, e.Position),
		handler.NewCol(ExecutionEventUserIDCol, e.UserID),
		handler.NewCol(ExecutionEventDataCol, e.EventData),
		handler.NewCol(ExecutionTargetsDataCol, e.TargetsData),
	}
}

func (e *EventExecution) Targets() ([]Target, error) {
	var execTargets []*query.ExecutionTarget
	if err := json.Unmarshal(e.TargetsData, execTargets); err != nil {
		return nil, err
	}
	targets := make([]Target, len(execTargets))
	for i, target := range execTargets {
		targets[i] = target
	}
	return targets, nil
}

func (e *EventExecution) ContextInfo() *ContextInfoEvent {
	return &ContextInfoEvent{
		AggregateID:   e.Aggregate.ID,
		AggregateType: string(e.Aggregate.Type),
		ResourceOwner: e.Aggregate.ResourceOwner,
		InstanceID:    e.Aggregate.InstanceID,
		Version:       string(e.Aggregate.Version),
		Sequence:      e.Sequence,
		EventType:     string(e.EventType),
		CreatedAt:     e.CreatedAt.Format(time.RFC3339Nano),
		Position:      e.Position,
		UserID:        e.UserID,
		EventPayload:  e.EventData,
	}
}

func (e *EventExecution) WithLogFields(entry *logging.Entry) *logging.Entry {
	return entry.
		WithField("instanceID", e.Aggregate.InstanceID).
		WithField("resourceOwner", e.Aggregate.ResourceOwner).
		WithField("aggregateType", e.Aggregate.Type).
		WithField("aggregateVersion", e.Aggregate.Version).
		WithField("aggregateID", e.Aggregate.ID).
		WithField("sequence", e.Sequence).
		WithField("eventType", e.EventType)
}
