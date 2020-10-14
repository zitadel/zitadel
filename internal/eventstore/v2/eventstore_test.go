package eventstore

import (
	"reflect"
	"testing"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

// testEvent implements the Event interface
type testEvent struct {
	description         string
	shouldCheckPrevious bool
	data                func() interface{}
}

func (e *testEvent) CheckPrevious() bool {
	return e.shouldCheckPrevious
}

func (e *testEvent) EditorService() string {
	return "editorService"
}
func (e *testEvent) EditorUser() string {
	return "editorUser"
}
func (e *testEvent) Type() EventType {
	return "test.event"
}
func (e *testEvent) Data() interface{} {
	return e.data()
}

func (e *testEvent) PreviousSequence() uint64 {
	return 0
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

func Test_eventData(t *testing.T) {
	type args struct {
		event Event
	}
	type res struct {
		jsonText []byte
		wantErr  bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "data as json bytes",
			args: args{
				event: &testEvent{
					data: func() interface{} {
						return []byte(`{"piff":"paff"}`)
					},
				},
			},
			res: res{
				jsonText: []byte(`{"piff":"paff"}`),
				wantErr:  false,
			},
		},
		{
			name: "data as invalid json bytes",
			args: args{
				event: &testEvent{
					data: func() interface{} {
						return []byte(`{"piffpaff"}`)
					},
				},
			},
			res: res{
				jsonText: []byte(nil),
				wantErr:  true,
			},
		},
		{
			name: "data as struct",
			args: args{
				event: &testEvent{
					data: func() interface{} {
						return struct {
							Piff string `json:"piff"`
						}{Piff: "paff"}
					},
				},
			},
			res: res{
				jsonText: []byte(`{"piff":"paff"}`),
				wantErr:  false,
			},
		},
		{
			name: "data as ptr to struct",
			args: args{
				event: &testEvent{
					data: func() interface{} {
						return &struct {
							Piff string `json:"piff"`
						}{Piff: "paff"}
					},
				},
			},
			res: res{
				jsonText: []byte(`{"piff":"paff"}`),
				wantErr:  false,
			},
		},
		{
			name: "no data",
			args: args{
				event: &testEvent{
					data: func() interface{} {
						return nil
					},
				},
			},
			res: res{
				jsonText: []byte(nil),
				wantErr:  false,
			},
		},
		{
			name: "invalid because primitive",
			args: args{
				event: &testEvent{
					data: func() interface{} {
						return ""
					},
				},
			},
			res: res{
				jsonText: []byte(nil),
				wantErr:  true,
			},
		},
		{
			name: "invalid because pointer to primitive",
			args: args{
				event: &testEvent{
					data: func() interface{} {
						var s string
						return &s
					},
				},
			},
			res: res{
				jsonText: []byte(nil),
				wantErr:  true,
			},
		},
		{
			name: "invalid because invalid struct for json",
			args: args{
				event: &testEvent{
					data: func() interface{} {
						return struct {
							Field chan string `json:"field"`
						}{}
					},
				},
			},
			res: res{
				jsonText: []byte(nil),
				wantErr:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := eventData(tt.args.event)
			if (err != nil) != tt.res.wantErr {
				t.Errorf("eventData() error = %v, wantErr %v", err, tt.res.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.res.jsonText) {
				t.Errorf("eventData() = %v, want %v", string(got), string(tt.res.jsonText))
			}
		})
	}
}
