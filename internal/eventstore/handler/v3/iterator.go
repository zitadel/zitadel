package handler

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/id"
)

type Iterator struct {
	client *sql.DB
	es     *eventstore.Eventstore
	pusher *Pusher

	bulkLimit  uint64
	aggregates []eventstore.AggregateType

	t              *time.Timer
	interval       time.Duration
	projectionName string
	lockerID       string
	reducers       map[eventstore.EventType]handler.Reduce
	iterationPool  chan func()
}

type IteratorConfig struct {
	Client     *sql.DB
	Eventstore *eventstore.Eventstore
	Interval   time.Duration
	Pool       chan func()

	BulkLimit uint64

	pusher         *Pusher
	projectionName string
	reducers       []handler.AggregateReducer
}

func NewIterator(config IteratorConfig) *Iterator {
	workerName, err := os.Hostname()
	if err != nil || workerName == "" {
		workerName, err = id.SonyFlakeGenerator.Next()
		logging.Log("SPOOL-bdO56").OnError(err).Panic("unable to generate lockID")
	}

	aggregateTypes := make([]eventstore.AggregateType, 0, len(config.reducers))
	reducers := make(map[eventstore.EventType]handler.Reduce, len(config.reducers))
	for _, aggReducer := range config.reducers {
		aggregateTypes = append(aggregateTypes, aggReducer.Aggregate)
		for _, eventReducer := range aggReducer.EventRedusers {
			reducers[eventReducer.Event] = eventReducer.Reduce
		}
	}

	return &Iterator{
		client:         config.Client,
		es:             config.Eventstore,
		pusher:         config.pusher,
		interval:       config.Interval,
		projectionName: config.projectionName,
		lockerID:       workerName,
		reducers:       reducers,
		t:              time.NewTimer(0),
		iterationPool:  config.Pool,
		bulkLimit:      config.BulkLimit,
		aggregates:     aggregateTypes,
	}
}

func (i *Iterator) Process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-i.t.C:
			if i.iterationPool == nil {
				i.iterate(ctx)
				continue
			}
			i.iterationPool <- func() { i.iterate(ctx) }
		}
	}
}

func (i *Iterator) iterate(ctx context.Context) {
	defer func() {
		i.t.Reset(i.interval)
	}()

	iterCtx, cancel := context.WithCancel(ctx)
	errs := make(chan error, 1)

	go i.lock(iterCtx, errs)
	err := <-errs
	if err != nil {
		cancel()
		return
	}

	i.errHandler(iterCtx, cancel, errs)

	events, err := i.query(iterCtx)
	if err != nil {
		cancel()
		i.unlock()
		return
	}
	i.prepareStmts(events)
	i.triggerPush(iterCtx)
	i.unlock()
	cancel()
}

func (i *Iterator) errHandler(ctx context.Context, cancel func(), errs <-chan error) {
	go func() {
		select {
		case err, ok := <-errs:
			if !ok {
				return
			}
			if err != nil {
				cancel()
				logging.LogWithFields("V3-3oXUV", "projection", i.projectionName).WithError(err).Warn("error occurred")
				return
			}
		case <-ctx.Done():
			return
		}
	}()
}

const (
	lockStmt = "INSERT INTO zitadel.projections.locks" +
		" (locker_id, locked_until, projection_name) VALUES ($1, now()+$2::INTERVAL, $3)" +
		" ON CONFLICT (projection_name)" +
		" DO UPDATE SET locker_id = $1, locked_until = now()+$2::INTERVAL" +
		" WHERE zitadel.projections.locks.projection_name = $3 AND (zitadel.projections.locks.locker_id = $1 OR zitadel.projections.locks.locked_until < now())"
)

func (i *Iterator) lock(ctx context.Context, errs chan<- error) {
	t := time.NewTimer(0)
	for {
		select {
		case <-ctx.Done():
			close(errs)
			return
		case <-t.C:
			//ensure same locker will be able to redo lock
			err := i.setLock(i.interval + 500*time.Millisecond)
			errs <- err
			t.Reset(i.interval)
		}
	}
}

func (i *Iterator) unlock() {
	err := i.setLock(0)
	logging.LogWithFields("V3-hhLZk", "projection", i.projectionName).OnError(err).Info("unlock failed")
}

func (i *Iterator) setLock(interval time.Duration) error {
	res, err := i.client.Exec(lockStmt, i.lockerID, interval.Seconds(), i.projectionName)
	if err != nil {
		return errors.ThrowInternal(err, "CRDB-uaDoR", "unable to execute lock")
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.ThrowAlreadyExists(nil, "CRDB-mmi4J", "projection already locked")
	}
	return nil
}

func (i *Iterator) query(ctx context.Context) ([]eventstore.EventReader, error) {
	sequences, err := i.currentSequence()
	if err != nil {
		return nil, err
	}

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).Limit(i.bulkLimit)
	for _, aggregate := range i.aggregates {
		seq := sequences[aggregate]
		query.AddQuery().AggregateTypes(aggregate).SequenceGreater(seq)
	}

	return i.es.FilterEvents(ctx, query)
}

const (
	currentSequenceStmt = `SELECT current_sequence, aggregate_type FROM zitadel.projections.current_sequences WHERE projection_name = $1`
)

func (i *Iterator) currentSequence() (sequences map[eventstore.AggregateType]uint64, err error) {
	sequences = make(map[eventstore.AggregateType]uint64)
	rows, err := i.client.Query(currentSequenceStmt, i.projectionName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			aggregateType eventstore.AggregateType
			sequence      uint64
		)

		err = rows.Scan(&sequence, &aggregateType)
		if err != nil {
			return nil, errors.ThrowInternal(err, "CRDB-qFZIy", "scan failed")
		}

		sequences[aggregateType] = sequence
	}

	if err = rows.Close(); err != nil {
		return nil, errors.ThrowInternal(err, "CRDB-N4hCw", "close rows failed")
	}

	if err = rows.Err(); err != nil {
		return nil, errors.ThrowInternal(err, "CRDB-9H4QV", "errors in scanning rows")
	}

	return sequences, nil
}

func (i *Iterator) prepareStmts(events []eventstore.EventReader) {
	stmts := make([]handler.Statement, 0, len(events))
	for _, event := range events {
		reduce, ok := i.reducers[event.Type()]
		if !ok {
			stmts = append(stmts, NewNoOpStatement(event))
			continue
		}
		additionalStmts, err := reduce(event)
		logging.LogWithFields("V3-bqjBP", "projection", i.projectionName, "seq", event.Sequence()).OnError(err).Fatal("reduce failed")
		stmts = append(stmts, additionalStmts...)
	}
	i.pusher.appendStmts(stmts...)
}

func (i *Iterator) triggerPush(ctx context.Context) {
	err := i.pusher.push(ctx)
	logging.LogWithFields("V3-j61pi", "projection", i.projectionName).OnError(err).Info("push failed")
}
