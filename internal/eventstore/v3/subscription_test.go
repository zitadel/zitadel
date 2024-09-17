package eventstore

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func Test_subscriptions_Close(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan *eventstore.Notification, 1)
	s := &subscriptions{
		background: ctx,
		cancelBg:   cancel,
		eventTypes: map[eventstore.EventType][]chan<- *eventstore.Notification{
			"foo": {ch},
		},
	}
	s.Close()
	assert.True(t, s.mutex.TryLock())
	_, ok := <-ch
	assert.False(t, ok)
}

func Test_subscriptions_Add(t *testing.T) {
	s := &subscriptions{
		eventTypes: make(map[eventstore.EventType][]chan<- *eventstore.Notification),
	}
	ch := s.Add("foo", "bar")
	assert.True(t, s.mutex.TryLock())
	assert.Equal(t, NotificationChannelSize, cap(ch))
	assert.Len(t, s.eventTypes, 2)
	for _, chans := range s.eventTypes {
		assert.Len(t, chans, 1)
	}
}

func Test_subscriptions_GetSubscribedEvents(t *testing.T) {
	s := &subscriptions{
		eventTypes: map[eventstore.EventType][]chan<- *eventstore.Notification{
			"foo": nil,
			"bar": nil,
		},
	}
	arg := []eventstore.Event{
		&event{typ: "foo"},
		&event{typ: "bar"},
		&event{typ: "spanac"},
	}
	got := s.GetSubscribedEvents(arg)
	assert.True(t, s.mutex.TryLock())
	assert.Equal(t, arg[:2], got)
}

func Test_subscriptions_push(t *testing.T) {
	const eventType eventstore.EventType = "foo"
	type args struct {
		payload *eventstore.Notification
	}
	tests := []struct {
		name string
		ch   chan *eventstore.Notification
		args args
		want *eventstore.Notification
	}{
		{
			name: "send payload",
			ch:   make(chan *eventstore.Notification, 1),
			args: args{&eventstore.Notification{
				EventType: eventType,
				Position:  decimal.NewFromInt(123),
			}},
			want: &eventstore.Notification{
				EventType: eventType,
				Position:  decimal.NewFromInt(123),
			},
		},
		{
			name: "blocked channel",
			ch:   make(chan *eventstore.Notification),
			args: args{&eventstore.Notification{
				EventType: eventType,
				Position:  decimal.NewFromInt(123),
			}},
			want: nil,
		},
		{
			name: "unknown event type",
			ch:   make(chan *eventstore.Notification, 1),
			args: args{&eventstore.Notification{
				EventType: "bar",
				Position:  decimal.NewFromInt(123),
			}},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &subscriptions{
				eventTypes: map[eventstore.EventType][]chan<- *eventstore.Notification{
					eventType: {tt.ch},
				},
			}

			s.push(tt.args.payload)
			assert.True(t, s.mutex.TryLock())

			var got *eventstore.Notification
			select {
			case got = <-tt.ch:
			default:
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_subscriptions_notifyAll(t *testing.T) {
	ch := make(chan *eventstore.Notification, 2)
	s := &subscriptions{
		eventTypes: map[eventstore.EventType][]chan<- *eventstore.Notification{
			"foo": {ch},
			"bar": {ch},
		},
	}
	s.notifyAll()
	assert.True(t, s.mutex.TryLock())
	close(ch)

	want := []*eventstore.Notification{
		{EventType: "foo"},
		{EventType: "bar"},
	}

	got := make([]*eventstore.Notification, 0, 2)
	for n := range ch {
		got = append(got, n)
	}
	assert.Equal(t, want, got)
}

func Test_buildPgNotifyQuery(t *testing.T) {
	type args struct {
		events []eventstore.Event
	}
	tests := []struct {
		name      string
		args      args
		wantQuery string
		wantArgs  []any
		wantOk    bool
	}{
		{
			name:      "nil events",
			args:      args{nil},
			wantQuery: "",
			wantArgs:  nil,
			wantOk:    false,
		},
		{
			name: "1 event",
			args: args{[]eventstore.Event{
				&event{
					typ:      "foo",
					position: decimal.NewFromInt(1),
				},
			}},
			wantQuery: "SELECT pg_notify($1, $2);",
			wantArgs: []any{
				notificationChannelName,
				"{\"event_type\":\"foo\",\"position\":\"1\"}",
			},
			wantOk: true,
		},
		{
			name: "multiple events",
			args: args{[]eventstore.Event{
				&event{
					typ:      "foo",
					position: decimal.NewFromInt(1),
				},
				&event{
					typ:      "bar",
					position: decimal.NewFromInt(2),
				},
				&event{
					typ:      "spanac",
					position: decimal.NewFromInt(3),
				},
			}},
			wantQuery: "SELECT pg_notify($1, $2), pg_notify($1, $3), pg_notify($1, $4);",
			wantArgs: []any{
				notificationChannelName,
				"{\"event_type\":\"foo\",\"position\":\"1\"}",
				"{\"event_type\":\"bar\",\"position\":\"2\"}",
				"{\"event_type\":\"spanac\",\"position\":\"3\"}",
			},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotArgs, gotOk := buildPgNotifyQuery(tt.args.events)
			assert.Equal(t, tt.wantQuery, gotQuery)
			assert.Equal(t, tt.wantArgs, gotArgs)
			assert.Equal(t, tt.wantOk, gotOk)
		})
	}
}
