package eventstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	NotificationChannelSize = 100 // To be discussed

	notificationChannelName   = "zitadel_events"
	notificationListenQuery   = "listen " + notificationChannelName
	notificationUnlistenQuery = "unlisten " + notificationChannelName
)

type subscriptions struct {
	background context.Context
	cancelBg   context.CancelFunc

	db         *sql.DB
	eventTypes map[eventstore.EventType][]chan<- *eventstore.Notification
	mutex      sync.RWMutex
	waitGroup  sync.WaitGroup
}

func newSubscriptions(db *sql.DB) (*subscriptions, error) {
	s := &subscriptions{
		db:         db,
		eventTypes: make(map[eventstore.EventType][]chan<- *eventstore.Notification),
	}
	s.background, s.cancelBg = context.WithCancel(context.Background())
	s.waitGroup.Add(1)

	err := make(chan error)
	go s.listen(err)
	return s, <-err
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

// Add creates a notification channel with buffer size [NotificationChannelSize].
// When the buffer is full, payloads will be dropped and logged.
// The returned chan will be closed when the subscriptions are closed.
func (s *subscriptions) Add(eventTypes ...eventstore.EventType) <-chan *eventstore.Notification {
	ch := make(chan *eventstore.Notification, NotificationChannelSize)
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
		select {
		case ch <- payload:
		default:
			logging.WithFields("payload", payload).Warn("skipped push to full channel")
		}
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
// Rebuilds connections and waits for notification
// until the background context is cancelled.
// An error from the first connection setup is pushed on `ec`.
// If the first connection was successful `nil` is pushed on `ec`.
func (s *subscriptions) listen(ec chan<- error) {
	defer s.waitGroup.Done()

	var once sync.Once
	sendFirstError := func(err error) {
		once.Do(func() {
			ec <- err
			if err != nil {
				s.cancelBg() // terminate the listener
			}
		})
	}

	for {
		if err := s.background.Err(); err != nil {
			sendFirstError(err)
		}

		ctx, cancel := context.WithTimeout(s.background, time.Second)
		defer cancel()

		conn, err := s.db.Conn(ctx)
		if err != nil {
			logging.WithError(err).Error("subscription listener connection")
			sendFirstError(err)
			continue
		}

		err = conn.Raw(func(driverConn any) (err error) {
			defer conn.Close()
			conn, ok := driverConn.(*pgx.Conn)
			if !ok {
				return fmt.Errorf("wrong connection type %T expected %T", driverConn, conn)
			}
			_, err = conn.Exec(ctx, notificationListenQuery)
			if err != nil {
				return err
			}
			defer func() {
				_, err = conn.Exec(ctx, notificationUnlistenQuery)
			}()
			sendFirstError(nil) // setup went ok, tell the caller.
			s.notifyAll()
			return s.waitForNotifications(conn)
		})
		sendFirstError(err)
		select {
		case <-s.background.Done():
			logging.WithError(err).Info("subscription listener terminated")
			return
		case <-time.After(time.Second):
			logging.WithError(err).Error("subscription listener restarting")
		}
	}
}

// waitForNotifications waits and pushes every notification received from conn,
// until the background context is canceled and/or the connection returns an error.
// The returned error is never nil.
func (s *subscriptions) waitForNotifications(conn *pgx.Conn) error {
	logging.Info("wait for notification started")
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
