package eventstore_test

import (
	"context"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/query"
)

type OrgHandler struct {
	eventstore.ReadModelHandler
	orgs []*query.OrgReadModel
}

func NewOrgHandler(
	es *eventstore.Eventstore,
	queue *eventstore.JobQueue,
	requeueAfter time.Duration,
) *OrgHandler {
	return &OrgHandler{
		ReadModelHandler: *eventstore.NewReadModelHandler(es, queue, requeueAfter),
	}
}

func (h *OrgHandler) Process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logging.Log("EVENT-XG5Og").Info("stop processing")
			return
		case event := <-h.Handler.EventQueue:
			if h.HasLocked {
				h.processBulk(event)
				continue
			}
			h.process(event)
		}
	}
}

func (h *OrgHandler) process(eventstore.EventReader) error {
	return nil
}

func (h *OrgHandler) processBulk(event eventstore.EventReader) error {
	org := h.getBulkOrg(event.Aggregate().ID)
	org.AppendEvents(event)

	if h.BulkUntil == event.Sequence() {
		h.endBulk()
	}
	return nil
}

func (h *OrgHandler) endBulk() {
	stmts := make([]string, 0, len(h.orgs))
	for _, org := range h.orgs {
		stmt, err := org.Reduce()
		if err != nil {
			//TODO: how to handle this error?
			logging.LogWithFields("EVENT-VbCAf", "orgID", org.AggregateID).Warn("reduce failed")
			continue
		}
		stmts = append(stmts, stmt...)
	}
	//append unlock statement to stmts
	//execute stmts

	_ = stmts

	h.BulkUntil = 0
	h.HasLocked = false

}

func (h *OrgHandler) getBulkOrg(orgID string) *query.OrgReadModel {
	for _, org := range h.orgs {
		if org.AggregateID == orgID {
			return org
		}
	}
	org := query.NewOrgReadModel(orgID)
	h.orgs = append(h.orgs, org)
	return org
}
