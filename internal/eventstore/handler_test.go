package eventstore_test

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/query"
)

type OrgHandler struct {
	ctx context.Context
	eventstore.ReadModelHandler
	client *sql.DB

	shouldPush chan bool
	lock       sync.Mutex
	orgs       []*query.OrgReadModel
	pushSet    bool

	currentSequence uint64
}

func NewOrgHandler(
	ctx context.Context,
	es *eventstore.Eventstore,
	client *sql.DB,
	queue *eventstore.JobQueue,
	requeueAfter time.Duration,
) *OrgHandler {
	h := &OrgHandler{
		ctx:              ctx,
		ReadModelHandler: *eventstore.NewReadModelHandler(es, queue, requeueAfter),
		client:           client,
		shouldPush:       make(chan bool, 1),
	}
	go h.Process(ctx)

	return h
}

func (h *OrgHandler) Process(ctx context.Context) {
	for {
		// workaround to priorice cancel and events before push
		select {
		case <-ctx.Done():
			if h.pushSet {
				h.push()
			}
			logging.Log("EVENT-XG5Og").Info("stop processing")
			return
		case event := <-h.Handler.EventQueue:
			h.process(event)
			continue
		default:
			//continue to lower prio select
		}
		// if not canceled and no events push is allowed
		select {
		case <-ctx.Done():
			logging.Log("EVENT-XG5Og").Info("stop processing")
			return
		case event := <-h.Handler.EventQueue:
			h.process(event)
		case <-h.shouldPush:
			h.push()
		}
	}
}

func (h *OrgHandler) process(event eventstore.EventReader) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.currentSequence != event.PreviousSequence() {
		//TODO: load current sequence of the view

	}

	org := h.getOrg(event.Aggregate().ID)
	org.AppendEvents(event)
	err := org.Reduce()
	if err != nil {
		return err
	}

	if !h.pushSet {
		h.pushSet = true
		h.shouldPush <- true
	}

	return nil
}

func (h *OrgHandler) push() error {
	h.lock.Lock()
	defer h.lock.Unlock()

	//TODO: lock table
	//TODO: defer unlock table

	stmts := make([]eventstore.Statement, 0, len(h.orgs))
	for _, org := range h.orgs {
		stmts = append(stmts, org.Statements()...)
	}

	prepareds := make([]func(context.Context, *sql.Tx) error, len(stmts))
	for i, stmt := range stmts {
		prepareds[i] = stmt.Prepare()
	}

	tx, err := h.client.Begin()
	if err != nil {
		return err
	}

	for _, prepared := range prepareds {
		err = prepared(h.ctx, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		//TODO: should we trigger a bulk load then?
		return err
	}

	h.pushSet = false

	return nil
}

func (h *OrgHandler) getOrg(orgID string) *query.OrgReadModel {
	for _, org := range h.orgs {
		if org.AggregateID == orgID {
			return org
		}
	}
	org := query.NewOrgReadModel(orgID)
	h.orgs = append(h.orgs, org)
	return org
}
