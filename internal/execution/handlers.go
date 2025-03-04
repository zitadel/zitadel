package execution

import (
	"context"
	"encoding/json"
	"slices"
	"strings"

	"github.com/riverqueue/river"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/queue"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
)

const (
	HandlerTable = "projections.execution_handler"
)

type eventHandlerQueries interface {
	TargetsByExecutionID(ctx context.Context, ids []string) (execution []*query.ExecutionTarget, err error)
}

type Queue interface {
	Insert(ctx context.Context, args river.JobArgs, opts ...queue.InsertOpt) error
}

type eventHandler struct {
	eventTypes                 []string
	aggregateTypeFromEventType func(typ eventstore.EventType) eventstore.AggregateType
	query                      eventHandlerQueries
	queue                      Queue
}

func NewEventHandler(
	ctx context.Context,
	config handler.Config,
	eventTypes []string,
	aggregateTypeFromEventType func(typ eventstore.EventType) eventstore.AggregateType,
	query eventHandlerQueries,
	queue Queue,
) *handler.Handler {
	return handler.NewHandler(ctx, &config, &eventHandler{
		eventTypes:                 eventTypes,
		aggregateTypeFromEventType: aggregateTypeFromEventType,
		query:                      query,
		queue:                      queue,
	})
}

func (u *eventHandler) Name() string {
	return HandlerTable
}

func (u *eventHandler) Reducers() []handler.AggregateReducer {
	aggList := make(map[eventstore.AggregateType][]eventstore.EventType)
	for _, eventType := range u.eventTypes {
		aggType := u.aggregateTypeFromEventType(eventstore.EventType(eventType))
		aggEventTypes, ok := aggList[aggType]
		if !ok {
			aggList[aggType] = []eventstore.EventType{eventstore.EventType(eventType)}
		} else {
			found := false
			for _, aggEventType := range aggEventTypes {
				if aggEventType == eventstore.EventType(eventType) {
					found = true
				}
			}
			if !found {
				aggList[aggType] = append(aggList[aggType], eventstore.EventType(eventType))
			}
		}
	}

	aggReducers := make([]handler.AggregateReducer, len(aggList))
	i := 0
	for aggType, aggEventTypes := range aggList {
		eventReducers := make([]handler.EventReducer, len(aggEventTypes))
		for j, eventType := range aggEventTypes {
			eventReducers[j] = handler.EventReducer{
				Event:  eventType,
				Reduce: u.reduce,
			}
		}
		aggReducers[i] = handler.AggregateReducer{
			Aggregate:     aggType,
			EventReducers: eventReducers,
		}
		i++
	}
	return aggReducers
}

func groupsFromEventType(s string) []string {
	parts := strings.Split(s, ".")
	groups := make([]string, len(parts))
	groupBase := ""
	for i, part := range parts {
		if groupBase == "" {
			groupBase = part
		} else {
			groupBase = groupBase + "." + part
		}

		if groupBase == s {
			groups[i] = groupBase
		} else {
			groups[i] = groupBase + ".*"
		}
	}
	// sort to end up with the most specific group first
	slices.SortFunc(groups, func(a, b string) int {
		return strings.Compare(a, b) * -1
	})
	return groups
}

func idsForEventType(eventType string) []string {
	ids := make([]string, 0)
	for _, group := range groupsFromEventType(eventType) {
		ids = append(ids,
			exec_repo.ID(domain.ExecutionTypeEvent, group),
		)
	}
	return append(ids,
		exec_repo.IDAll(domain.ExecutionTypeEvent),
	)
}

func (u *eventHandler) reduce(e eventstore.Event) (*handler.Statement, error) {
	ctx := HandlerContext(e.Aggregate())

	targets, err := u.query.TargetsByExecutionID(ctx, idsForEventType(string(e.Type())))
	if err != nil {
		return nil, err
	}

	// no execution from worker necessary
	if len(targets) == 0 {
		return handler.NewNoOpStatement(e), nil
	}

	return handler.NewStatement(e, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(e.Aggregate())
		req, err := NewRequest(e, targets)
		if err != nil {
			return err
		}
		return u.queue.Insert(ctx,
			req,
			queue.WithQueueName(exec_repo.QueueName),
		)
	}), nil
}

func NewRequest(e eventstore.Event, targets []*query.ExecutionTarget) (*exec_repo.Request, error) {
	targetsData, err := json.Marshal(targets)
	if err != nil {
		return nil, err
	}
	return &exec_repo.Request{
		Aggregate:   e.Aggregate(),
		Sequence:    e.Sequence(),
		EventType:   e.Type(),
		CreatedAt:   e.CreatedAt(),
		UserID:      e.Creator(),
		EventData:   e.DataAsBytes(),
		TargetsData: targetsData,
	}, nil
}
