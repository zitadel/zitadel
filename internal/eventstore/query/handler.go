package query

import (
	"context"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
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
	QueryLimit() uint64

	AggregateTypes() []models.AggregateType
	CurrentSequence(*models.Event) (uint64, error)
	Eventstore() eventstore.Eventstore
}

func ReduceEvent(handler Handler, event *models.Event) {
	currentSequence, err := handler.CurrentSequence(event)
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

	processedSequences := map[models.AggregateType]uint64{}

	for _, unprocessedEvent := range unprocessedEvents {
		currentSequence, err := handler.CurrentSequence(unprocessedEvent)
		if err != nil {
			logging.Log("HANDL-BmpkC").WithError(err).Warn("unable to get current sequence")
			return
		}
		_, ok := processedSequences[unprocessedEvent.AggregateType]
		if !ok {
			processedSequences[unprocessedEvent.AggregateType] = currentSequence
		}
		if processedSequences[unprocessedEvent.AggregateType] != currentSequence {
			logging.LogWithFields("QUERY-DOYVN",
				"processed", processedSequences[unprocessedEvent.AggregateType],
				"current", currentSequence).
				Warn("sequence not matching")
			return
		}

		err = handler.Reduce(unprocessedEvent)
		logging.LogWithFields("HANDL-V42TI", "seq", unprocessedEvent.Sequence).OnError(err).Warn("reduce failed")
		processedSequences[unprocessedEvent.AggregateType] = unprocessedEvent.Sequence
	}
	if len(unprocessedEvents) == eventLimit {
		logging.LogWithFields("QUERY-BSqe9", "seq", event.Sequence).Warn("didnt process event")
		return
	}
	err = handler.Reduce(event)
	logging.LogWithFields("HANDL-wQDL2", "seq", event.Sequence).OnError(err).Warn("reduce failed")
}
