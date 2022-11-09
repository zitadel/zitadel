package v3

import (
	"context"
	"database/sql"
	"sync"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
)

var (
	once sync.Once
)

type IDProjection struct {
	Name    string
	reduces map[eventstore.AggregateType][]Reducer
	client  *sql.DB
	es      *eventstore.Eventstore
	check   *handler.Check
	version string
}

func NewConfig(client *sql.DB, es *eventstore.Eventstore) *Config {
	return &Config{
		client:     client,
		eventstore: es,
	}
}

type Config struct {
	client     *sql.DB
	eventstore *eventstore.Eventstore
	Check      *handler.Check
	Reduces    map[eventstore.AggregateType][]Reducer
	Version    string
}

type Reducer struct {
	Event          eventstore.EventType
	PreviousEvents func(*sql.Tx, eventstore.Event) (*eventstore.SearchQueryBuilder, error)
	Reduce         handler.Reduce
}

func StartSubscriptionIDProjection(ctx context.Context, name string, config Config) *IDProjection {
	p := New(name, config)

	err := p.Init(ctx)
	logging.OnError(err).WithField("projection", name).Fatal("unable to initialize projection")

	go func() {
		sub := p.subscribe()
		for {
			select {
			case <-ctx.Done():
				sub.Unsubscribe()
				return
			case e := <-sub.Events:
				err := p.Process(ctx, e)
				logging.WithFields("name", name).OnError(err).Error("error occured in reduce, stop processing")
			}
		}
	}()

	return p
}

func New(name string, config Config) *IDProjection {
	return &IDProjection{
		client:  config.client,
		es:      config.eventstore,
		Name:    name,
		reduces: config.Reduces,
		check:   config.Check,
		version: config.Version,
	}
}

func (p *IDProjection) Start() {
	ctx := context.TODO()
	go func() {
		sub := p.subscribe()
		for {
			select {
			case <-ctx.Done():
				sub.Unsubscribe()
				return
			case e := <-sub.Events:
				err := p.Process(ctx, e)
				logging.WithFields("name", p.Name).OnError(err).Error("error occured in reduce, stop processing")
			}
		}
	}()

	once.Do(func() {
		row := p.client.QueryRow("SELECT COUNT(*) FROM " + p.Name + " AS OF SYSTEM TIME '-1ms'")
		var count int
		err := row.Scan(&count)
		if err != nil || count > 0 {
			logging.OnError(err).WithField("name", p.Name).Error("unable to check projection")
			return
		}

		err = p.es.TriggerSubscriptions(ctx,
			eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
				OrderAsc().
				AddQuery().
				ExcludedInstanceID("").
				Builder())
		logging.OnError(err).WithField("name", p.Name).Error("unable to trigger subscribers")
	})
}

// Process updates the projection by the given events
func (p *IDProjection) Process(ctx context.Context, event eventstore.Event) error {
	// for event := range events {
	reducer, ok := p.eventReduce(event)
	if !ok {
		logging.WithFields("eventType", event.Type()).Info("no reducer registered")
		return nil
		// continue
	}

	tx, err := p.client.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var stmts handler.Statements

	if reducer.PreviousEvents != nil {
		previousEventsQuery, err := reducer.PreviousEvents(tx, event)
		if err != nil {
			tx.Rollback()
			return err
		}
		if previousEventsQuery != nil {
			previousEvents, err := p.es.Filter(ctx, previousEventsQuery)
			if err != nil {
				tx.Rollback()
				return err
			}
			stmts, err = p.reducePreviousEvents(previousEvents)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	stmt, err := reducer.Reduce(event)
	if err != nil {
		tx.Rollback()
		return err
	}

	stmts = append(stmts, *stmt)

	return p.execStmts(ctx, tx, stmts)
}

func (p *IDProjection) Init(ctx context.Context) error {
	// for _, check := range checks {
	if p.check == nil || p.check.IsNoop() {
		return nil
	}
	tx, err := p.client.BeginTx(ctx, nil)
	if err != nil {
		return errors.ThrowInternal(err, "V3-iSUtO", "begin failed")
	}
	for i, execute := range p.check.Executes {
		logging.WithFields("projection", p.Name, "execute", i).Debug("executing check")
		next, err := execute(p.client, p.Name)
		if err != nil {
			logging.OnError(tx.Rollback()).Debug("unable to rollback tx")
			return err
		}
		if !next {
			logging.WithFields("projection", p.Name, "execute", i).Debug("skipping next check")
			break
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	// }
	return nil
}

func (p *IDProjection) eventReduce(event eventstore.Event) (reduce Reducer, ok bool) {
	aggReduces, ok := p.reduces[event.Aggregate().Type]
	if !ok {
		return reduce, false
	}
	for _, r := range aggReduces {
		if r.Event == event.Type() {
			return r, true
		}
	}

	return reduce, false
}

func (p *IDProjection) subscribe() *eventstore.Subscription {
	queue := make(chan eventstore.Event, 100)
	types := map[eventstore.AggregateType][]eventstore.EventType{}
	for agg, reduces := range p.reduces {
		types[agg] = make([]eventstore.EventType, len(reduces))
		for i, reduce := range reduces {
			types[agg][i] = reduce.Event
		}
	}
	return eventstore.SubscribeEventTypes(queue, types)
}

func (p *IDProjection) reducePreviousEvents(events []eventstore.Event) (handler.Statements, error) {
	stmts := make(handler.Statements, 0, len(events))
	for _, event := range events {
		reducer, ok := p.eventReduce(event)
		if !ok {
			logging.WithFields("eventType", event.Type()).Info("no additional reducer registered")
			return nil, errors.ThrowInternal(nil, "V3-nkOs7", "no additional reducer registered")
		}
		stmt, err := reducer.Reduce(event)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, *stmt)
	}

	return stmts, nil
}

func (p *IDProjection) execStmts(ctx context.Context, tx *sql.Tx, stmts handler.Statements) error {
	for _, stmt := range stmts {
		if err := stmt.Execute(tx, p.Name); err != nil {
			logging.WithFields(
				"event", stmt.EventID,
				"instance", stmt.InstanceID,
			).WithError(err).Error("unable to execute statement")
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// Trigger is a noop function to fulfill previous trigger functions
func (p *IDProjection) Trigger(context.Context) {}
