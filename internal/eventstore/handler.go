package eventstore

import (
	"context"
	"fmt"
	"time"

	"github.com/caos/logging"
)

type Handler struct {
	Eventstore *Eventstore
	JobQueue   *JobQueue
	Sub        *Subscription
	EventQueue chan EventReader
}

func NewHandler(
	eventstore *Eventstore,
	queue *JobQueue,
) *Handler {
	h := Handler{
		Eventstore: eventstore,
		JobQueue:   queue,
		//TODO: how huge should the queue be?
		EventQueue: make(chan EventReader, 100),
	}
	//TODO: start handler for EventQueue

	return &h
}

func (h Handler) Subscribe(aggregates ...AggregateType) {
	h.Sub = Subscribe(h.EventQueue, aggregates...)
}

type ReadModelHandler struct {
	Handler
	RequeueAfter time.Duration
	Timer        *time.Timer
	HasLocked    bool
	BulkUntil    uint64
}

func NewReadModelHandler(
	eventstore *Eventstore,
	queue *JobQueue,
	requeueAfter time.Duration,
) *ReadModelHandler {
	return &ReadModelHandler{
		Handler:      *NewHandler(eventstore, queue),
		RequeueAfter: requeueAfter,
		// first requeue is instant on startup
		Timer: time.NewTimer(0),
	}
}

func (h *ReadModelHandler) ResetTimer() {
	h.Timer.Reset(h.RequeueAfter)
}

func (h *ReadModelHandler) Lock() (currentSequence uint64, err error) {
	return 0, nil
}

func (h *ReadModelHandler) Queue(ctx context.Context, query *SearchQueryBuilder) error {
	latestSeq, err := h.Handler.JobQueue.Queue(ctx, query, h.EventQueue)
	if err != nil {
		return err
	}
	if latestSeq > query.eventSequence {
		h.BulkUntil = latestSeq
	}
	return nil
}

func (h *ReadModelHandler) Process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logging.Log("EVENT-XG5Og").Info("stop processing")
			return
		case event := <-h.Handler.EventQueue:
			func(event EventReader) {
				logging.LogWithFields("EVENT-Sls2l", "seq", event.Sequence()).Info("event received")
			}(event)
		}
	}
}

type Column struct {
	Name  string
	Value interface{}
}

type Statement interface {
	ToSQL() (stmt string, args []interface{})
}

type CreateStatement struct {
	TableName string
	Values    []Column
}

func (stmt *CreateStatement) ToSQL() (query string, args []interface{}) {
	// db := sql.OpenDB(nil)
	return fmt.Sprintf("INSERT INTO %s", stmt.TableName), nil
}

type UpdateStatement struct {
	PK     []Column
	Values []Column
}

func (stmt *UpdateStatement) ToSQL() (string, []interface{}) {
	return "", nil
}

type DeleteStatement struct {
	PK []Column
}

func (stmt *DeleteStatement) ToSQL() (string, []interface{}) {
	return "", nil
}
