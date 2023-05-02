package query

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/zitadel/logging"

	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

const (
	eventLimit = 10000
)

type Handler interface {
	ViewModel() string
	EventQuery(ctx context.Context, instanceIDs []string) (*models.SearchQuery, error)
	Reduce(*models.Event) error
	OnError(event *models.Event, err error) error
	OnSuccess(instanceIDs []string) error
	MinimumCycleDuration() time.Duration
	LockDuration() time.Duration
	QueryLimit() uint64

	AggregateTypes() []models.AggregateType
	CurrentSequence(ctx context.Context, instanceID string) (uint64, error)
	Eventstore() v1.Eventstore

	Subscription() *v1.Subscription
}

func ReduceEvent(ctx context.Context, handler Handler, event *models.Event) {
	defer func() {
		err := recover()

		if err != nil {
			handler.Subscription().Unsubscribe()
			logging.WithFields(
				"cause", err,
				"stack", string(debug.Stack()),
				"sequence", event.Sequence,
				"instance", event.InstanceID,
			).Error("reduce panicked")
		}
	}()
	currentSequence, err := handler.CurrentSequence(ctx, event.InstanceID)
	if err != nil {
		logging.WithError(err).Warn("unable to get current sequence")
		return
	}

	searchQuery := models.NewSearchQuery().
		AddQuery().
		AggregateTypeFilter(handler.AggregateTypes()...).
		SequenceBetween(currentSequence, event.Sequence).
		InstanceIDFilter(event.InstanceID).
		SearchQuery().
		SetLimit(eventLimit)

	unprocessedEvents, err := handler.Eventstore().FilterEvents(ctx, searchQuery)
	if err != nil {
		logging.WithFields("sequence", event.Sequence).Warn("filter failed")
		return
	}

	for _, unprocessedEvent := range unprocessedEvents {
		currentSequence, err := handler.CurrentSequence(ctx, unprocessedEvent.InstanceID)
		if err != nil {
			logging.WithError(err).Warn("unable to get current sequence")
			return
		}
		if unprocessedEvent.Sequence < currentSequence {
			logging.WithFields(
				"unprocessed", unprocessedEvent.Sequence,
				"current", currentSequence,
				"view", handler.ViewModel()).
				Warn("sequence not matching")
			return
		}

		err = handler.Reduce(unprocessedEvent)
		logging.WithFields("sequence", unprocessedEvent.Sequence).OnError(err).Warn("reduce failed")
	}
	if len(unprocessedEvents) == eventLimit {
		logging.WithFields("sequence", event.Sequence).Warn("didnt process event")
		return
	}
	err = handler.Reduce(event)
	logging.WithFields("sequence", event.Sequence).OnError(err).Warn("reduce failed")
}
