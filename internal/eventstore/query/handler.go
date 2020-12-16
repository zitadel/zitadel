package query

import (
	"context"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
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
	CurrentSequence() (uint64, error)
	Eventstore() eventstore.Eventstore
}

func ReduceEvent(handler Handler, event *models.Event) {
	sequence, err := handler.CurrentSequence()
	if err != nil {
		logging.Log("HANDL-BmpkC").WithError(err).Warn("unable to get current sequence")
		return
	}
	if event.PreviousSequence > sequence {
		searchQuery := models.NewSearchQuery().
			AggregateTypeFilter(handler.AggregateTypes()...).
			SequenceBetween(sequence, event.PreviousSequence)

		events, err := handler.Eventstore().FilterEvents(context.Background(), searchQuery)
		if err != nil {
			logging.LogWithFields("HANDL-L6YH1", "seq", event.Sequence).Warn("filter failed")
			return
		}
		for _, previousEvent := range events {
			//if other process already updated view
			if event.PreviousSequence > previousEvent.Sequence {
				continue
			}
			err = handler.Reduce(previousEvent)
			logging.LogWithFields("HANDL-V42TI", "seq", previousEvent.Sequence).OnError(err).Warn("reduce failed")
			return
		}
	} else if event.PreviousSequence < sequence {
		logging.LogWithFields("HANDL-w9Bdy", "previousSeq", event.PreviousSequence, "currentSeq", sequence).Debug("already processed")
		return
	}
	err = handler.Reduce(event)
	logging.LogWithFields("HANDL-wQDL2", "seq", event.Sequence).OnError(err).Warn("reduce failed")
}
