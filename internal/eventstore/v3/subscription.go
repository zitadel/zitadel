package eventstore

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type subscriptions struct {
	pool       *pgxpool.Pool
	eventTypes map[eventstore.EventType][]chan<- decimal.Decimal
	mutex      sync.RWMutex
	cancelWait context.CancelFunc
	waitGroup  sync.WaitGroup
}

func newSubscriptions(pool *pgxpool.Pool) *subscriptions {
	return &subscriptions{
		pool:       pool,
		eventTypes: make(map[eventstore.EventType][]chan<- decimal.Decimal),
	}
}

func (s *subscriptions) Close() {
	s.mutex.RLock()
	if s.cancelWait != nil {
		s.cancelWait()
	}
	s.mutex.RUnlock()
	s.waitGroup.Wait()
}

func (s *subscriptions) Add(ch chan<- decimal.Decimal, eventTypes ...eventstore.EventType) {
	s.mutex.Lock()
	if s.cancelWait != nil {
		s.cancelWait()
	}
	// waitForNotifications shutdown
	s.waitGroup.Wait()

	for _, typ := range eventTypes {
		s.eventTypes[typ] = append(s.eventTypes[typ], ch)
	}

	// create new context within the same Lock, so that next calls
	// to Add cancel the correct context.
	var ctx context.Context
	ctx, s.cancelWait = context.WithCancel(context.Background())
	s.mutex.Unlock()
	s.waitForNotifications(ctx)
}

// GetSubscribedEvents returns a subset of events which currently have a subscription.
func (s *subscriptions) GetSubscribedEvents(events []eventstore.Event) []eventstore.Event {
	out := make([]eventstore.Event, 0, len(events))

	s.mutex.RLock()
	for _, event := range events {
		if _, has := s.eventTypes[event.Type()]; has {
			out = append(out, event)
		}
	}
	s.mutex.RUnlock()

	return out
}

// Push a payload to all channels subscribed to eventType.
func (s *subscriptions) Push(eventType eventstore.EventType, payload decimal.Decimal) {
	s.mutex.RLock()
	for _, ch := range s.eventTypes[eventType] {
		ch <- payload
	}
	s.mutex.RUnlock()
}

// waitForNotifications hijacks a connection from the pool and LISTENs for all configured events.
// After it returns the listener is fully setup.
// A go routine will be created, waiting for notifications from the database.
// [subscriptions.cancelWait] needs to be called to cleanup the routine and the connection.
// Upon termination connections are closed and not returned to the pool.
// This is to prevent active LISTENing connections in the pool.
func (s *subscriptions) waitForNotifications(ctx context.Context) {
	conn, err := s.buildListenerConn(ctx)
	if err != nil {
		return
	}

	s.waitGroup.Add(1)
	go func() {
		defer func(ctx context.Context) {
			ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
			err := conn.Close(ctx)
			cancel()
			logging.OnError(err).Warn("notification connection close")
			s.waitGroup.Done()
		}(ctx)

		for {
			notification, err := conn.WaitForNotification(ctx)
			if errors.Is(err, context.Canceled) {
				break // intended termination
			}
			if err != nil {
				logging.WithError(err).Error("eventstore wait for notifications")
				s.waitForNotifications(ctx) // rebuild
				return
			}

			payload, err := decimal.NewFromString(notification.Payload)
			logging.WithError(err).Error("invalid notification payload")
			s.Push(eventstore.EventType(notification.Channel), payload)
		}
	}()
}

// buildListenerConn acquires a connection and issue LISTEN commands for all event types.
// On error it retries indefinitely unless the context is canceled.
// The returned conn is hijacked and should be closed by the caller.
func (s *subscriptions) buildListenerConn(ctx context.Context) (*pgx.Conn, error) {
	var queries strings.Builder

	s.mutex.RLock()
	for typ := range s.eventTypes {
		fmt.Fprintf(&queries, "listen %s;\n", typ)
	}
	s.mutex.RUnlock()

	var conn *pgxpool.Conn
	for {
		var err error
		localCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		conn, err = s.pool.Acquire(localCtx)
		cancel()
		if err != nil {
			logging.WithError(err).Error("acquire notification connection")
			select {
			case <-ctx.Done():
				return nil, err
			case <-time.After(time.Second):
				continue // try again
			}
		}

		_, err = conn.Exec(ctx, queries.String())
		if err != nil {
			conn.Release()
			logging.WithError(err).Error("set notification listeners")
			select {
			case <-ctx.Done():
				return nil, err
			case <-time.After(time.Second):
				continue // try again
			}
		}
		break
	}

	// Hijack so the pool knows it's never getting this connection back.
	return conn.Hijack(), nil
}

// buildPgNotifyQuery builds a single query with multiple pg_notify, one for each event.
func buildPgNotifyQuery(events []eventstore.Event) (query string, args []any, ok bool) {
	if len(events) == 0 {
		return "", nil, false
	}

	notifies := make([]string, len(events))
	args = make([]any, 0, len(events)*2)

	for i, event := range events {
		notifies[i] = fmt.Sprintf("pg_notify($%d, $%d)", i*2+1, i*2+2)
		args = append(args, event.Type(), event.Position())
	}

	return fmt.Sprintf("SELECT %s;", strings.Join(notifies, ", ")), args, true
}
