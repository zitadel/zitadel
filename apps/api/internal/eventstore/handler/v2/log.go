package handler

import (
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func (h *Handler) log() *logging.Entry {
	return logging.WithFields("projection", h.projection.Name())
}

func (h *Handler) logFailure(fail *failure) *logging.Entry {
	return h.log().WithField("sequence", fail.sequence).
		WithField("instance", fail.instance).
		WithField("aggregate", fail.aggregateID)
}

func (h *Handler) logEvent(event eventstore.Event) *logging.Entry {
	return h.log().WithField("sequence", event.Sequence()).
		WithField("instance", event.Aggregate().InstanceID).
		WithField("aggregate", event.Aggregate().Type)
}
