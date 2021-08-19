package handler

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
)

type Projectioner interface {
	Ensure(ctx context.Context)
}

type MapEvent func(eventstore.EventReader) handler.Statement

type Projection struct {
	name     string
	mapEvent MapEvent

	triggerMu   sync.Mutex
	triggerSet  time.Time
	execTrigger *time.Timer
	retryAfter  time.Duration

	es             *eventstore.Eventstore
	aggregateTypes []eventstore.AggregateType

	interval time.Duration

	dbClient *sql.DB
}

func NewProjection(
	name string,
	mapper MapEvent,
	es *eventstore.Eventstore,
	dbClient *sql.DB,
	aggregateTypes ...eventstore.AggregateType,
) Projectioner {
	p := &Projection{
		name:           name,
		mapEvent:       mapper,
		aggregateTypes: aggregateTypes,
		es:             es,
		interval:       30 * time.Second,
		retryAfter:     2 * time.Second,
		execTrigger:    time.NewTimer(0),
		dbClient:       dbClient,
	}

	//unitialized timer
	//https://github.com/golang/go/issues/12721
	<-p.execTrigger.C

	return p
}

func (p *Projection) Ensure(ctx context.Context) {
	errs := make(chan error)
	queue := make(chan eventstore.EventReader, 100)

	locker := NewMutexLocker()
	subscriber := NewSubscriptionHandler(p.aggregateTypes)
	subscriber.Query(ctx, queue)
	iterator := NewIterationHandler(IterationHandlerConfig{
		Eventstore: p.es,
		Interval:   p.interval,
		PreSteps: []PreStep{
			p.lockPrsStep(locker),
		},
		PostSteps: []PostStep{
			p.triggerExecPostStep(),
		},
	})
	iterator.Query(ctx, queue)
	sqlExec := NewSQLExecuter(SQLExecuterConfig{
		Client:         p.dbClient,
		ProjectionName: p.name,
		PostSteps: []PostStep{
			p.unlockPostStep(locker),
		},
	})

	stmts := []handler.Statement{}
	for {
		select {
		case <-ctx.Done():
			subscriber.Cancel()
		case e := <-queue:
			stmts = append(stmts, p.mapEvent(e))
		default:
			select {
			case <-ctx.Done():
				subscriber.Cancel()
			case e := <-queue:
				stmts = append(stmts, p.mapEvent(e))
			case <-p.execTrigger.C:
				var err error
				stmts, err = sqlExec.Execute(ctx, stmts)
				if err != nil {
					logging.LogWithFields("V2-o10Ls", "projection", p.name).WithError(err).Warn("exec failed")
					p.setTrigger(p.retryAfter)
				}
			}
		}
	}
}

func (p *Projection) cancelOnErr(ctx context.Context, errs <-chan error, cancel func()) {
	select {
	case <-ctx.Done():
		return
	case err := <-errs:
		logging.LogWithFields("V2-Q1qyy", "projection", p.name).WithError(err).Info("received err")
		cancel()
	}
}

func (p *Projection) lockPrsStep(locker Locker) func() error {
	return func() error {
		err := locker.Lock()
		logging.LogWithFields("V2-DFBBD", "projection", p.name).OnError(err).Debug("lock failed")
		return err
	}
}

func (p *Projection) unlockPostStep(locker Locker) func() error {
	return func() error {
		locker.Unlock()
		return nil
	}
}

func (p *Projection) triggerExecPostStep() func() error {
	return func() error {
		p.setTrigger(0)
		return nil
	}
}

func (p *Projection) setTrigger(d time.Duration) {
	p.triggerMu.Lock()
	defer p.triggerMu.Unlock()

	if p.triggerSet.Before(time.Now().Add(d)) {
		p.execTrigger.Reset(d)
		p.triggerSet = time.Now().Add(d)
	}
}
