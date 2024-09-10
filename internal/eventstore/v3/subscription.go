package eventstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	notificationChannelName   = "zitadel_events"
	notificationListenQuery   = "listen " + notificationChannelName
	notificationUnlistenQuery = "unlisten " + notificationChannelName
)

type subscriptions struct {
	background context.Context
	cancelBg   context.CancelFunc

	pool       *pgxpool.Pool
	eventTypes map[eventstore.EventType][]chan<- *eventstore.Notification
	mutex      sync.RWMutex
	waitGroup  sync.WaitGroup
}

func newSubscriptions(pool *pgxpool.Pool) *subscriptions {
	s := &subscriptions{
		pool:       pool,
		eventTypes: make(map[eventstore.EventType][]chan<- *eventstore.Notification),
	}
	s.background, s.cancelBg = context.WithCancel(context.Background())
	s.waitGroup.Add(1)
	go s.listen()
	return s
}

func (s *subscriptions) Close() {
	s.cancelBg()
	s.waitGroup.Wait()

	s.mutex.Lock()
	for _, channels := range s.eventTypes {
		for _, ch := range channels {
			close(ch)
		}
	}
	s.eventTypes = map[eventstore.EventType][]chan<- *eventstore.Notification{}
	s.mutex.Unlock()
}

// Add creates a notification channel.
// The returned chan will be closed when the subscriptions are closed.
func (s *subscriptions) Add(eventTypes ...eventstore.EventType) <-chan *eventstore.Notification {
	ch := make(chan *eventstore.Notification) // TBD: Do we need a buffered channel?
	s.mutex.Lock()
	for _, typ := range eventTypes {
		s.eventTypes[typ] = append(s.eventTypes[typ], ch)
	}
	s.mutex.Unlock()
	return ch
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
func (s *subscriptions) push(payload *eventstore.Notification) {
	s.mutex.RLock()
	for _, ch := range s.eventTypes[payload.EventType] {
		ch <- payload
	}
	s.mutex.RUnlock()
}

// notifyAll sends an [eventstore.Notification] to all subscribers.
func (s *subscriptions) notifyAll() {
	s.mutex.RLock()
	for typ, channels := range s.eventTypes {
		for _, ch := range channels {
			ch <- &eventstore.Notification{
				EventType: typ,
			}
		}
	}
	s.mutex.RUnlock()
}

// listen for zitadel_events on a notification channel.
// Rebuilds connections and waiter until the background context is cancelled.
func (s *subscriptions) listen() {
	for {
		err := s.waitForNotifications()
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			logging.WithError(err).Info("subscription listener terminated")
			return
		}
		logging.WithError(err).Error("subscription listener restarting")
		time.Sleep(time.Second)
	}
}

// waitForNotifications hijacks a connection from the pool and LISTENs on the zitadel_events notification channel,
// until the background context is canceled or the connection returns an error.
// Upon termination the connection is closed and not returned to the pool.
// This is to prevent active LISTENing connections in the pool.
// The returned error is never nil.
func (s *subscriptions) waitForNotifications() error {
	ctx, cancel := context.WithTimeout(s.background, 5*time.Second)
	conn, err := s.getListenerConn(ctx)
	cancel()
	if err != nil {
		return err
	}
	defer func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		_, err := conn.Exec(ctx, notificationUnlistenQuery)
		err = errors.Join(err, conn.Close(ctx))
		cancel()
		logging.OnError(err).Warn("subscription listener cleanup")
	}(s.background)

	s.notifyAll()
	logging.Info("subscription listener started")

	for {
		notification, err := conn.WaitForNotification(s.background)
		if err != nil {
			return err
		}
		var payload *eventstore.Notification
		err = json.Unmarshal([]byte(notification.Payload), &payload)
		logging.OnError(err).Error("invalid notification payload")
		s.push(payload)
	}
}

// getListenerConn acquires a connection and issue LISTEN commands for all event types.
// On error it retries until the context is expires.
// The returned conn is hijacked and should be closed by the caller.
func (s *subscriptions) getListenerConn(ctx context.Context) (*pgx.Conn, error) {
	for {
		conn, err := func() (*pgxpool.Conn, error) {
			conn, err := s.pool.Acquire(ctx)
			if err != nil {
				return nil, err
			}
			_, err = conn.Exec(ctx, notificationListenQuery)
			if err != nil {
				conn.Release()
				return nil, err
			}
			return conn, nil
		}()
		if err != nil {
			logging.WithError(err).Error("get listener connection")
			select {
			case <-ctx.Done():
				return nil, err
			case <-time.After(time.Second):
				continue // try again
			}
		}
		// Hijack so the pool knows it's never getting this connection back.
		return conn.Hijack(), nil
	}
}

// buildPgNotifyQuery builds a single query with multiple pg_notify, one for each event.
func buildPgNotifyQuery(events []eventstore.Event) (query string, args []any, ok bool) {
	if len(events) == 0 {
		return "", nil, false
	}

	notifies := make([]string, len(events))
	args = make([]any, 1, len(events)*2)
	args[0] = notificationChannelName

	for i, event := range events {
		payload, err := json.Marshal(eventstore.Notification{
			EventType: event.Type(),
			Position:  event.Position(),
		})
		if err != nil {
			logging.WithError(err).WithField("event", event).Error("build notify payload")
			continue
		}

		notifies[i] = fmt.Sprintf("pg_notify($1, $%d)", i+2)
		args = append(args, string(payload))
		ok = true
	}

	return fmt.Sprintf("SELECT %s;", strings.Join(notifies, ", ")), args, ok
}
