package query

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v1"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	eventLimit = 10000
)

type Handler interface {
	ViewModel() string
	EventQuery() (*models.SearchQuery, error)
	Reduce(*models.Event) error
	OnError(event *models.Event, err error) error
	OnSuccess() error
	MinimumCycleDuration() time.Duration
	LockDuration() time.Duration
	QueryLimit() uint64

	AggregateTypes() []models.AggregateType
	CurrentSequence() (uint64, error)
	Eventstore() v1.Eventstore
}

func ReduceEvent(handler Handler, event *models.Event) {
	currentSequence, err := handler.CurrentSequence()
	if err != nil {
		logging.Log("HANDL-BmpkC").WithError(err).Warn("unable to get current sequence")
		return
	}

	searchQuery := models.NewSearchQuery().
		AggregateTypeFilter(handler.AggregateTypes()...).
		SequenceBetween(currentSequence, event.Sequence).
		SetLimit(eventLimit)

	unprocessedEvents, err := handler.Eventstore().FilterEvents(context.Background(), searchQuery)
	if err != nil {
		logging.LogWithFields("HANDL-L6YH1", "seq", event.Sequence).Warn("filter failed")
		return
	}

	for _, unprocessedEvent := range unprocessedEvents {
		currentSequence, err := handler.CurrentSequence()
		if err != nil {
			logging.Log("HANDL-BmpkC").WithError(err).Warn("unable to get current sequence")
			return
		}
		if unprocessedEvent.Sequence < currentSequence {
			logging.LogWithFields("QUERY-DOYVN",
				"unprocessed", unprocessedEvent.Sequence,
				"current", currentSequence,
				"view", handler.ViewModel()).
				Warn("sequence not matching")
			return
		}

		err = handler.Reduce(unprocessedEvent)
		logging.LogWithFields("HANDL-V42TI", "seq", unprocessedEvent.Sequence).OnError(err).Warn("reduce failed")
	}
	if len(unprocessedEvents) == eventLimit {
		logging.LogWithFields("QUERY-BSqe9", "seq", event.Sequence).Warn("didnt process event")
		return
	}
	err = handler.Reduce(event)
	logging.LogWithFields("HANDL-wQDL2", "seq", event.Sequence).OnError(err).Warn("reduce failed")
}
