package eventstore

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
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

func columnsToQuery(cols []Column) (names []string, parameters []string, values []interface{}) {
	names = make([]string, len(cols))
	values = make([]interface{}, len(cols))
	parameters = make([]string, len(cols))
	for i, col := range cols {
		names[i] = col.Name
		values[i] = col.Value
		parameters[i] = "$" + strconv.Itoa(i+1)

	}
	return names, parameters, values
}

func columnsToWhere(cols []Column, paramOffset int) (wheres []string, values []interface{}) {
	wheres = make([]string, len(cols))
	values = make([]interface{}, len(cols))

	for i, col := range cols {
		wheres[i] = "(" + col.Name + " = " + strconv.Itoa(i+1+paramOffset) + ")"
		values[i] = col.Value
	}

	return wheres, values
}

type Statement interface {
	Prepare(ctx context.Context, tx *sql.Tx) func() (sql.Result, error)
}

type CreateStatement struct {
	TableName string
	Values    []Column
}

func (stmt *CreateStatement) Prepare(ctx context.Context, tx *sql.Tx) func() (sql.Result, error) {
	cols, params, args := columnsToQuery(stmt.Values)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", stmt.TableName, strings.Join(cols, ", "), strings.Join(params, ", "))
	statement, err := tx.PrepareContext(ctx, query)

	return func() (sql.Result, error) {
		if err != nil {
			return nil, err
		}
		return statement.ExecContext(ctx, args)
	}
}

type UpdateStatement struct {
	TableName string
	PK        []Column
	Values    []Column
}

func (stmt *UpdateStatement) Prepare(ctx context.Context, tx *sql.Tx) func() (sql.Result, error) {
	cols, params, args := columnsToQuery(stmt.Values)
	wheres, whereArgs := columnsToWhere(stmt.PK, len(params))
	args = append(args, whereArgs)
	query := fmt.Sprintf("UPDATE %s SET (%s) = (%s) WHERE %s", stmt.TableName, strings.Join(cols, ", "), strings.Join(params, ", "), strings.Join(wheres, " AND "))
	statement, err := tx.PrepareContext(ctx, query)

	return func() (sql.Result, error) {
		if err != nil {
			return nil, err
		}
		return statement.ExecContext(ctx, args)
	}
}

type DeleteStatement struct {
	TableName string
	PK        []Column
}

func (stmt *DeleteStatement) Prepare(ctx context.Context, tx *sql.Tx) func() (sql.Result, error) {
	wheres, args := columnsToWhere(stmt.PK, 0)
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", stmt.TableName, strings.Join(wheres, " AND "))
	statement, err := tx.PrepareContext(ctx, query)

	return func() (sql.Result, error) {
		if err != nil {
			return nil, err
		}
		return statement.ExecContext(ctx, args)
	}
}
