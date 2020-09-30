package eventstore

import (
	"reflect"
	"testing"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

type testEvent struct {
	description         string
	shouldCheckPrevious bool
}

func (e *testEvent) CheckPrevious() bool {
	return e.shouldCheckPrevious
}

func testPushMapper(Event) (*repository.Event, error) {
	return &repository.Event{AggregateID: "aggregateID"}, nil
}

func Test_eventstore_RegisterPushEventMapper(t *testing.T) {
	type fields struct {
		eventMapper map[EventType]eventTypeInterceptors
	}
	type args struct {
		eventType EventType
		mapper    func(Event) (*repository.Event, error)
	}
	type res struct {
		event *repository.Event
		isErr func(error) bool
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "no event type",
			args: args{
				eventType: "",
				mapper:    testPushMapper,
			},
			fields: fields{
				eventMapper: map[EventType]eventTypeInterceptors{},
			},
			res: res{
				isErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "no event mapper",
			args: args{
				eventType: "event.type",
				mapper:    nil,
			},
			fields: fields{
				eventMapper: map[EventType]eventTypeInterceptors{},
			},
			res: res{
				isErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "new interceptor",
			fields: fields{
				eventMapper: map[EventType]eventTypeInterceptors{},
			},
			args: args{
				eventType: "new.event",
				mapper:    testPushMapper,
			},
			res: res{
				event: &repository.Event{AggregateID: "aggregateID"},
			},
		},
		{
			name: "existing interceptor new push mapper",
			fields: fields{
				eventMapper: map[EventType]eventTypeInterceptors{
					"existing": {},
				},
			},
			args: args{
				eventType: "new.event",
				mapper:    testPushMapper,
			},
			res: res{
				event: &repository.Event{AggregateID: "aggregateID"},
			},
		},
		{
			name: "existing interceptor existing push mapper",
			fields: fields{
				eventMapper: map[EventType]eventTypeInterceptors{
					"existing": {
						pushMapper: func(Event) (*repository.Event, error) {
							return nil, errors.ThrowUnimplemented(nil, "V2-1qPvn", "unimplemented")
						},
					},
				},
			},
			args: args{
				eventType: "new.event",
				mapper:    testPushMapper,
			},
			res: res{
				event: &repository.Event{AggregateID: "aggregateID"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &Eventstore{
				eventMapper: tt.fields.eventMapper,
			}
			err := es.RegisterPushEventMapper(tt.args.eventType, tt.args.mapper)
			if (tt.res.isErr != nil && !tt.res.isErr(err)) || (tt.res.isErr == nil && err != nil) {
				t.Errorf("wrong error got: %v", err)
				return
			}
			if tt.res.isErr != nil {
				return
			}

			mapper := es.eventMapper[tt.args.eventType]
			event, err := mapper.pushMapper(nil)
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			if !reflect.DeepEqual(tt.res.event, event) {
				t.Errorf("events should be deep equal. \ngot %v\nwant %v", event, tt.res.event)
			}
		})
	}
}

func testFilterMapper(*repository.Event) (Event, error) {
	return &testEvent{description: "hodor"}, nil
}

func Test_eventstore_RegisterFilterEventMapper(t *testing.T) {
	type fields struct {
		eventMapper map[EventType]eventTypeInterceptors
	}
	type args struct {
		eventType EventType
		mapper    func(*repository.Event) (Event, error)
	}
	type res struct {
		event Event
		isErr func(error) bool
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "no event type",
			args: args{
				eventType: "",
				mapper:    testFilterMapper,
			},
			fields: fields{
				eventMapper: map[EventType]eventTypeInterceptors{},
			},
			res: res{
				isErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "no event mapper",
			args: args{
				eventType: "event.type",
				mapper:    nil,
			},
			fields: fields{
				eventMapper: map[EventType]eventTypeInterceptors{},
			},
			res: res{
				isErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "new interceptor",
			fields: fields{
				eventMapper: map[EventType]eventTypeInterceptors{},
			},
			args: args{
				eventType: "event.type",
				mapper:    testFilterMapper,
			},
			res: res{
				event: &testEvent{description: "hodor"},
			},
		},
		{
			name: "existing interceptor new filter mapper",
			fields: fields{
				eventMapper: map[EventType]eventTypeInterceptors{
					"event.type": {},
				},
			},
			args: args{
				eventType: "new.event",
				mapper:    testFilterMapper,
			},
			res: res{
				event: &testEvent{description: "hodor"},
			},
		},
		{
			name: "existing interceptor existing filter mapper",
			fields: fields{
				eventMapper: map[EventType]eventTypeInterceptors{
					"event.type": {
						filterMapper: func(*repository.Event) (Event, error) {
							return nil, errors.ThrowUnimplemented(nil, "V2-1qPvn", "unimplemented")
						},
					},
				},
			},
			args: args{
				eventType: "new.event",
				mapper:    testFilterMapper,
			},
			res: res{
				event: &testEvent{description: "hodor"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &Eventstore{
				eventMapper: tt.fields.eventMapper,
			}
			err := es.RegisterFilterEventMapper(tt.args.eventType, tt.args.mapper)
			if (tt.res.isErr != nil && !tt.res.isErr(err)) || (tt.res.isErr == nil && err != nil) {
				t.Errorf("wrong error got: %v", err)
				return
			}
			if tt.res.isErr != nil {
				return
			}

			mapper := es.eventMapper[tt.args.eventType]
			event, err := mapper.filterMapper(nil)
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			if !reflect.DeepEqual(tt.res.event, event) {
				t.Errorf("events should be deep equal. \ngot %v\nwant %v", event, tt.res.event)
			}
		})
	}
}
