package eventstore

import (
	"fmt"
	"strings"
	"sync"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type subscriptions struct {
	mutex      sync.RWMutex
	eventTypes map[eventstore.EventType][]chan<- decimal.Decimal
}

func newSubscriptions() *subscriptions {
	return &subscriptions{
		eventTypes: make(map[eventstore.EventType][]chan<- decimal.Decimal),
	}
}

func (s *subscriptions) Add(ch chan<- decimal.Decimal, eventTypes ...eventstore.EventType) {
	s.mutex.Lock()
	for _, typ := range eventTypes {
		s.eventTypes[typ] = append(s.eventTypes[typ], ch)
	}
	s.mutex.Unlock()
}

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
