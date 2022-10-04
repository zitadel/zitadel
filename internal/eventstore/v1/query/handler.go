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
	EventQuery(instanceIDs ...string) (*models.SearchQuery, error)
	Reduce(*models.Event) error
	OnError(event *models.Event, err error) error
	OnSuccess() error
	MinimumCycleDuration() time.Duration
	LockDuration() time.Duration
	QueryLimit() uint64

	AggregateTypes() []models.AggregateType
	CurrentCreationDate(instanceID string) (time.Time, error)
	Eventstore() v1.Eventstore

	Subscription() *v1.Subscription
}

func ReduceEvent(handler Handler, event *models.Event) {
	defer func() {
		err := recover()

		if err != nil {
			handler.Subscription().Unsubscribe()
			logging.WithFields(
				"cause", err,
				"stack", string(debug.Stack()),
				"event", event.ID,
				"instance", event.InstanceID,
			).Error("reduce panicked")
		}
	}()
	currentCreationDate, err := handler.CurrentCreationDate(event.InstanceID)
	if err != nil {
		logging.WithError(err).Warn("unable to get current sequence")
		return
	}

	searchQuery := models.NewSearchQuery().
		AddQuery().
		AggregateTypeFilter(handler.AggregateTypes()...).
		CreationDateBetweenFilter(currentCreationDate, event.CreationDate).
		InstanceIDFilter(event.InstanceID).
		SearchQuery().
		SetLimit(eventLimit)

	unprocessedEvents, err := handler.Eventstore().FilterEvents(context.Background(), searchQuery)
	if err != nil {
		logging.WithFields("eventId", event.ID).Warn("filter failed")
		return
	}

	for _, unprocessedEvent := range unprocessedEvents {
		currentCreationDate, err := handler.CurrentCreationDate(unprocessedEvent.InstanceID)
		if err != nil {
			logging.WithError(err).Warn("unable to get current sequence")
			return
		}
		if unprocessedEvent.CreationDate.Before(currentCreationDate) {
			logging.WithFields(
				"unprocessed", unprocessedEvent.ID,
				"current", currentCreationDate,
				"view", handler.ViewModel()).
				Warn("sequence not matching")
			return
		}

		err = handler.Reduce(unprocessedEvent)
		logging.WithFields("event", unprocessedEvent.ID).OnError(err).Warn("reduce failed")
	}
	if len(unprocessedEvents) == eventLimit {
		logging.WithFields("event", event.ID).Warn("didnt process event")
		return
	}
	err = handler.Reduce(event)
	logging.WithFields("event", event.ID).OnError(err).Warn("reduce failed")
}
