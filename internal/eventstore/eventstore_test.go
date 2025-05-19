package eventstore

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/service"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// testEvent implements the Event interface
type testEvent struct {
	BaseEvent

	description         string
	shouldCheckPrevious bool
	data                func() interface{}
}

func newTestEvent(id, description string, data func() interface{}, checkPrevious bool) *testEvent {
	return &testEvent{
		description:         description,
		data:                data,
		shouldCheckPrevious: checkPrevious,
		BaseEvent: *NewBaseEventForPush(
			service.WithService(authz.NewMockContext("instanceID", "resourceOwner", "editorUser"), "editorService"),
			NewAggregate(authz.NewMockContext("zitadel", "caos", "adlerhurst"), id, "test.aggregate", "v1"),
			"test.event",
		),
	}
}

func (e *testEvent) Payload() interface{} {
	return e.data()
}

func (e *testEvent) UniqueConstraints() []*UniqueConstraint {
	return nil
}

func (e *testEvent) Assets() []*Asset {
	return nil
}

func testFilterMapper(event Event) (Event, error) {
	if event == nil {
		return newTestEvent("testID", "hodor", nil, false), nil
	}
	return &testEvent{description: "hodor", BaseEvent: *BaseEventFromRepo(event)}, nil
}

func Test_eventstore_RegisterFilterEventMapper(t *testing.T) {
	type fields struct {
		eventMapper map[EventType]eventTypeInterceptors
	}
	type args struct {
		eventType EventType
		mapper    func(Event) (Event, error)
	}
	type res struct {
		event       Event
		mapperCount int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
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
				mapperCount: 0,
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
				event:       newTestEvent("testID", "hodor", nil, false),
				mapperCount: 1,
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
				event:       newTestEvent("testID", "hodor", nil, false),
				mapperCount: 2,
			},
		},
		{
			name: "existing interceptor existing filter mapper",
			fields: fields{
				eventMapper: map[EventType]eventTypeInterceptors{
					"event.type": {
						eventMapper: func(Event) (Event, error) {
							return nil, zerrors.ThrowUnimplemented(nil, "V2-1qPvn", "unimplemented")
						},
					},
				},
			},
			args: args{
				eventType: "new.event",
				mapper:    testFilterMapper,
			},
			res: res{
				event:       newTestEvent("testID", "hodor", nil, false),
				mapperCount: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			eventInterceptors = tt.fields.eventMapper
			RegisterFilterEventMapper("test", tt.args.eventType, tt.args.mapper)
			if len(eventInterceptors) != tt.res.mapperCount {
				t.Errorf("unexpected mapper count: want %d, got %d", tt.res.mapperCount, len(eventInterceptors))
				return
			}

			if tt.res.mapperCount == 0 {
				return
			}

			mapper := eventInterceptors[tt.args.eventType]
			event, err := mapper.eventMapper(nil)
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			if !reflect.DeepEqual(tt.res.event, event) {
				t.Errorf("events should be deep equal. \ngot %#v\nwant %#v", event, tt.res.event)
			}
		})
	}
}

func Test_eventData(t *testing.T) {
	type args struct {
		event Command
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
				event: newTestEvent(
					"id",
					"hodor",
					func() interface{} {
						return []byte(`{"piff":"paff"}`)
					},
					false),
			},
			res: res{
				jsonText: []byte(`{"piff":"paff"}`),
				wantErr:  false,
			},
		},
		{
			name: "data as invalid json bytes",
			args: args{
				event: newTestEvent(
					"id",
					"hodor",
					func() interface{} {
						return []byte(`{"piffpaff"}`)
					},
					false),
			},
			res: res{
				jsonText: []byte(nil),
				wantErr:  true,
			},
		},
		{
			name: "data as struct",
			args: args{
				event: newTestEvent(
					"id",
					"hodor",
					func() interface{} {
						return struct {
							Piff string `json:"piff"`
						}{Piff: "paff"}
					},
					false),
			},
			res: res{
				jsonText: []byte(`{"piff":"paff"}`),
				wantErr:  false,
			},
		},
		{
			name: "data as ptr to struct",
			args: args{
				event: newTestEvent(
					"id",
					"hodor",
					func() interface{} {
						return &struct {
							Piff string `json:"piff"`
						}{Piff: "paff"}
					},
					false),
			},
			res: res{
				jsonText: []byte(`{"piff":"paff"}`),
				wantErr:  false,
			},
		},
		{
			name: "no data",
			args: args{
				event: newTestEvent(
					"id",
					"hodor",
					func() interface{} {
						return nil
					},
					false),
			},
			res: res{
				jsonText: []byte(nil),
				wantErr:  false,
			},
		},
		{
			name: "invalid because primitive",
			args: args{
				event: newTestEvent(
					"id",
					"hodor",
					func() interface{} {
						return ""
					},
					false),
			},
			res: res{
				jsonText: []byte(nil),
				wantErr:  true,
			},
		},
		{
			name: "invalid because pointer to primitive",
			args: args{
				event: newTestEvent(
					"id",
					"hodor",
					func() interface{} {
						var s string
						return &s
					},
					false),
			},
			res: res{
				jsonText: []byte(nil),
				wantErr:  true,
			},
		},
		{
			name: "invalid because invalid struct for json",
			args: args{
				event: newTestEvent(
					"id",
					"hodor",
					func() interface{} {
						return struct {
							Field chan string `json:"field"`
						}{}
					},
					false),
			},
			res: res{
				jsonText: []byte(nil),
				wantErr:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EventData(tt.args.event)
			if (err != nil) != tt.res.wantErr {
				t.Errorf("EventData() error = %v, wantErr %v", err, tt.res.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.res.jsonText) {
				t.Errorf("EventData() = %v, want %v", string(got), string(tt.res.jsonText))
			}
		})
	}
}

var _ Pusher = (*testPusher)(nil)

func (repo *testPusher) Client() *database.DB {
	return nil
}

type testPusher struct {
	events []Event
	errs   []error

	t *testing.T
}

func (repo *testPusher) Health(ctx context.Context) error {
	return nil
}

func (repo *testPusher) Push(_ context.Context, _ database.ContextQueryExecuter, commands ...Command) (events []Event, err error) {
	if len(repo.errs) != 0 {
		err, repo.errs = repo.errs[0], repo.errs[1:]
		return nil, err
	}

	if len(repo.events) != len(commands) {
		repo.t.Errorf("length of events unequal want: %d got %d", len(repo.events), len(commands))
		return nil, fmt.Errorf("")
	}

	events = make([]Event, len(commands))

	for i := 0; i < len(repo.events); i++ {

		var payload []byte
		switch p := commands[i].Payload().(type) {
		case []byte:
			payload = p
		case nil:
			// don't do anything
		default:
			payload, err = json.Marshal(p)
			if err != nil {
				return nil, err
			}
		}

		compareEvents(repo.t, repo.events[i], commands[i])
		events[i] = &BaseEvent{
			Seq:       uint64(i),
			Creation:  time.Now(),
			EventType: commands[i].Type(),
			Data:      payload,
			User:      commands[i].Creator(),
			Agg: &Aggregate{
				Version:       commands[i].Aggregate().Version,
				ID:            commands[i].Aggregate().ID,
				Type:          commands[i].Aggregate().Type,
				ResourceOwner: commands[i].Aggregate().ResourceOwner,
				InstanceID:    commands[i].Aggregate().InstanceID,
			},
		}
	}

	return events, nil
}

type testQuerier struct {
	events    []Event
	sequence  float64
	instances []string
	err       error
	t         *testing.T
}

func (repo *testQuerier) Health(ctx context.Context) error {
	return nil
}

func (repo *testQuerier) CreateInstance(ctx context.Context, instance string) error {
	return nil
}

func (repo *testQuerier) Filter(ctx context.Context, searchQuery *SearchQueryBuilder) ([]Event, error) {
	if repo.err != nil {
		return nil, repo.err
	}
	return repo.events, nil
}

func (repo *testQuerier) FilterToReducer(ctx context.Context, searchQuery *SearchQueryBuilder, reduce Reducer) error {
	if repo.err != nil {
		return repo.err
	}
	for _, event := range repo.events {
		if err := reduce(event); err != nil {
			return err
		}
	}
	return nil
}

func (repo *testQuerier) LatestSequence(ctx context.Context, queryFactory *SearchQueryBuilder) (float64, error) {
	if repo.err != nil {
		return 0, repo.err
	}
	return repo.sequence, nil
}

func (repo *testQuerier) InstanceIDs(ctx context.Context, queryFactory *SearchQueryBuilder) ([]string, error) {
	if repo.err != nil {
		return nil, repo.err
	}
	return repo.instances, nil
}

func (*testQuerier) Client() *database.DB {
	return nil
}

func TestEventstore_Push(t *testing.T) {
	type args struct {
		events []Command
	}
	type fields struct {
		maxRetries  int
		pusher      *testPusher
		eventMapper map[EventType]func(Event) (Event, error)
	}
	type res struct {
		wantErr bool
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		res    res
	}{
		{
			name: "one aggregate one event",
			args: args{
				events: []Command{
					newTestEvent(
						"1",
						"",
						func() interface{} {
							return []byte(nil)
						},
						false),
				},
			},
			fields: fields{
				maxRetries: 1,
				pusher: &testPusher{
					t: t,
					events: []Event{
						&BaseEvent{
							Agg: &Aggregate{
								ID:            "1",
								Type:          "test.aggregate",
								ResourceOwner: "caos",
								InstanceID:    "zitadel",
								Version:       "v1",
							},
							Data:      []byte(nil),
							User:      "editorUser",
							EventType: "test.event",
						},
					},
				},
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(e Event) (Event, error) {
						return &testEvent{
							BaseEvent: BaseEvent{
								Agg: &Aggregate{
									Type: e.Aggregate().Type,
								},
							},
						}, nil
					},
				},
			},
		},
		{
			name: "one aggregate one event, retry disabled",
			args: args{
				events: []Command{
					newTestEvent(
						"1",
						"",
						func() interface{} {
							return []byte(nil)
						},
						false),
				},
			},
			fields: fields{
				maxRetries: 0,
				pusher: &testPusher{
					t: t,
					events: []Event{
						&BaseEvent{
							Agg: &Aggregate{
								ID:            "1",
								Type:          "test.aggregate",
								ResourceOwner: "caos",
								InstanceID:    "zitadel",
								Version:       "v1",
							},
							Data:      []byte(nil),
							User:      "editorUser",
							EventType: "test.event",
						},
					},
				},
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(e Event) (Event, error) {
						return &testEvent{
							BaseEvent: BaseEvent{
								Agg: &Aggregate{
									Type: e.Aggregate().Type,
								},
							},
						}, nil
					},
				},
			},
		},
		{
			name: "one aggregate multiple events",
			args: args{
				events: []Command{
					newTestEvent(
						"1",
						"",
						func() interface{} {
							return []byte(nil)
						},
						false),
					newTestEvent(
						"1",
						"",
						func() interface{} {
							return []byte(nil)
						},
						false),
				},
			},
			fields: fields{
				maxRetries: 1,
				pusher: &testPusher{
					t: t,
					events: []Event{
						&BaseEvent{
							Agg: &Aggregate{
								ID:            "1",
								Type:          "test.aggregate",
								ResourceOwner: "caos",
								InstanceID:    "zitadel",
								Version:       "v1",
							},
							Data:      []byte(nil),
							User:      "editorUser",
							EventType: "test.event",
						},
						&BaseEvent{
							Agg: &Aggregate{
								ID:            "1",
								Type:          "test.aggregate",
								ResourceOwner: "caos",
								InstanceID:    "zitadel",
								Version:       "v1",
							},
							Data:      []byte(nil),
							User:      "editorUser",
							EventType: "test.event",
						},
					},
				},
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(e Event) (Event, error) {
						return &testEvent{
							BaseEvent: BaseEvent{
								Agg: &Aggregate{
									Type: e.Aggregate().Type,
								},
							},
						}, nil
					},
				},
			},
			res: res{
				wantErr: false,
			},
		},
		{
			name: "multiple aggregates",
			args: args{
				events: []Command{
					newTestEvent(
						"1",
						"",
						func() interface{} {
							return []byte(nil)
						},
						false),
					newTestEvent(
						"1",
						"",
						func() interface{} {
							return []byte(nil)
						},
						false),
					newTestEvent(
						"2",
						"",
						func() interface{} {
							return []byte(nil)
						},
						true),
				},
			},
			fields: fields{
				maxRetries: 1,
				pusher: &testPusher{
					t: t,
					events: combineEventLists(
						[]Event{
							&BaseEvent{
								Agg: &Aggregate{
									ID:            "1",
									Type:          "test.aggregate",
									ResourceOwner: "caos",
									InstanceID:    "zitadel",
									Version:       "v1",
								},
								Data:      []byte(nil),
								User:      "editorUser",
								EventType: "test.event",
							},
							&BaseEvent{
								Agg: &Aggregate{
									ID:            "1",
									Type:          "test.aggregate",
									ResourceOwner: "caos",
									InstanceID:    "zitadel",
									Version:       "v1",
								},
								Data:      []byte(nil),
								User:      "editorUser",
								EventType: "test.event",
							},
						},
						[]Event{
							&BaseEvent{
								Agg: &Aggregate{
									ID:            "2",
									Type:          "test.aggregate",
									ResourceOwner: "caos",
									InstanceID:    "zitadel",
									Version:       "v1",
								},
								Data:      []byte(nil),
								User:      "editorUser",
								EventType: "test.event",
							},
						},
					),
				},
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(e Event) (Event, error) {
						return &testEvent{
							BaseEvent: BaseEvent{
								Agg: &Aggregate{
									Type: e.Aggregate().Type,
								},
							},
						}, nil
					},
				},
			},
			res: res{
				wantErr: false,
			},
		},
		{
			name: "push fails",
			args: args{
				events: []Command{
					newTestEvent(
						"1",
						"",
						func() interface{} {
							return []byte(nil)
						},
						false),
				},
			},
			fields: fields{
				maxRetries: 1,
				pusher: &testPusher{
					t:    t,
					errs: []error{zerrors.ThrowInternal(nil, "V2-qaa4S", "test err")},
				},
			},
			res: res{
				wantErr: true,
			},
		},
		{
			name: "aggreagtes to events mapping fails",
			args: args{
				events: []Command{
					newTestEvent(
						"1",
						"",
						func() interface{} {
							return `{"data":""`
						},
						false),
				},
			},
			fields: fields{
				maxRetries: 1,
				pusher: &testPusher{
					t:    t,
					errs: []error{zerrors.ThrowInternal(nil, "V2-qaa4S", "test err")},
				},
			},
			res: res{
				wantErr: true,
			},
		},
		{
			name: "retry succeeds",
			args: args{
				events: []Command{
					newTestEvent(
						"1",
						"",
						func() interface{} {
							return []byte(nil)
						},
						false),
				},
			},
			fields: fields{
				maxRetries: 1,
				pusher: &testPusher{
					t: t,
					events: []Event{
						&BaseEvent{
							Agg: &Aggregate{
								ID:            "1",
								Type:          "test.aggregate",
								ResourceOwner: "caos",
								InstanceID:    "zitadel",
								Version:       "v1",
							},
							Data:      []byte(nil),
							User:      "editorUser",
							EventType: "test.event",
						},
					},
					errs: []error{
						zerrors.ThrowInternal(&pgconn.PgError{
							ConstraintName: "events2_pkey",
							Code:           "23505",
						}, "foo-err", "Errors.Internal"),
					},
				},
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(e Event) (Event, error) {
						return &testEvent{
							BaseEvent: BaseEvent{
								Agg: &Aggregate{
									Type: e.Aggregate().Type,
								},
							},
						}, nil
					},
				},
			},
		},
		{
			name: "retry fails",
			args: args{
				events: []Command{
					newTestEvent(
						"1",
						"",
						func() interface{} {
							return []byte(nil)
						},
						false),
				},
			},
			fields: fields{
				maxRetries: 1,
				pusher: &testPusher{
					t: t,
					events: []Event{
						&BaseEvent{
							Agg: &Aggregate{
								ID:            "1",
								Type:          "test.aggregate",
								ResourceOwner: "caos",
								InstanceID:    "zitadel",
								Version:       "v1",
							},
							Data:      []byte(nil),
							User:      "editorUser",
							EventType: "test.event",
						},
					},
					errs: []error{
						zerrors.ThrowInternal(&pgconn.PgError{
							ConstraintName: "events2_pkey",
							Code:           "23505",
						}, "foo-err", "Errors.Internal"),
						zerrors.ThrowInternal(&pgconn.PgError{
							ConstraintName: "events2_pkey",
							Code:           "23505",
						}, "foo-err", "Errors.Internal"),
					},
				},
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(e Event) (Event, error) {
						return &testEvent{
							BaseEvent: BaseEvent{
								Agg: &Aggregate{
									Type: e.Aggregate().Type,
								},
							},
						}, nil
					},
				},
			},
			res: res{wantErr: true},
		},
		{
			name: "retry disabled",
			args: args{
				events: []Command{
					newTestEvent(
						"1",
						"",
						func() interface{} {
							return []byte(nil)
						},
						false),
				},
			},
			fields: fields{
				maxRetries: 0,
				pusher: &testPusher{
					t: t,
					events: []Event{
						&BaseEvent{
							Agg: &Aggregate{
								ID:            "1",
								Type:          "test.aggregate",
								ResourceOwner: "caos",
								InstanceID:    "zitadel",
								Version:       "v1",
							},
							Data:      []byte(nil),
							User:      "editorUser",
							EventType: "test.event",
						},
					},
					errs: []error{
						zerrors.ThrowInternal(&pgconn.PgError{
							ConstraintName: "events2_pkey",
							Code:           "23505",
						}, "foo-err", "Errors.Internal"),
					},
				},
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(e Event) (Event, error) {
						return &testEvent{
							BaseEvent: BaseEvent{
								Agg: &Aggregate{
									Type: e.Aggregate().Type,
								},
							},
						}, nil
					},
				},
			},
			res: res{wantErr: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventInterceptors = map[EventType]eventTypeInterceptors{}
			es := &Eventstore{
				maxRetries: tt.fields.maxRetries,
				pusher:     tt.fields.pusher,
			}
			for eventType, mapper := range tt.fields.eventMapper {
				RegisterFilterEventMapper("test", eventType, mapper)
			}
			if len(eventInterceptors) != len(tt.fields.eventMapper) {
				t.Errorf("register event mapper failed expected mapper amount: %d, got: %d", len(tt.fields.eventMapper), len(eventInterceptors))
				t.FailNow()
			}

			_, err := es.Push(context.Background(), tt.args.events...)
			if (err != nil) != tt.res.wantErr {
				t.Errorf("Eventstore.aggregatesToEvents() error = %v, wantErr %v", err, tt.res.wantErr)
			}
		})
	}
}

func TestEventstore_FilterEvents(t *testing.T) {
	type args struct {
		query *SearchQueryBuilder
	}
	type fields struct {
		repo        *testQuerier
		eventMapper map[EventType]func(Event) (Event, error)
	}
	type res struct {
		wantErr bool
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		res    res
	}{
		{
			name: "no events",
			args: args{
				query: &SearchQueryBuilder{
					columns: ColumnsEvent,
					queries: []*SearchQuery{
						{
							builder:        &SearchQueryBuilder{},
							aggregateTypes: []AggregateType{"no.aggregates"},
						},
					},
				},
			},
			fields: fields{
				repo: &testQuerier{
					events: []Event{},
					t:      t,
				},
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(e Event) (Event, error) {
						return &testEvent{}, nil
					},
				},
			},
			res: res{
				wantErr: false,
			},
		},
		{
			name: "repo error",
			args: args{
				query: &SearchQueryBuilder{
					columns: ColumnsEvent,
					queries: []*SearchQuery{
						{
							builder:        &SearchQueryBuilder{},
							aggregateTypes: []AggregateType{"no.aggregates"},
						},
					},
				},
			},
			fields: fields{
				repo: &testQuerier{
					t:   t,
					err: zerrors.ThrowInternal(nil, "V2-RfkBa", "test err"),
				},
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(e Event) (Event, error) {
						return &testEvent{}, nil
					},
				},
			},
			res: res{
				wantErr: true,
			},
		},
		{
			name: "found events",
			args: args{
				query: &SearchQueryBuilder{
					columns: ColumnsEvent,
					queries: []*SearchQuery{
						{
							builder:        &SearchQueryBuilder{},
							aggregateTypes: []AggregateType{"test.aggregate"},
						},
					},
				},
			},
			fields: fields{
				repo: &testQuerier{
					events: []Event{
						&BaseEvent{
							Agg: &Aggregate{
								ID: "test.aggregate",
							},
							EventType: "test.event",
						},
					},
					t: t,
				},
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(e Event) (Event, error) {
						return &testEvent{}, nil
					},
				},
			},
			res: res{
				wantErr: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventInterceptors = map[EventType]eventTypeInterceptors{}
			es := &Eventstore{
				querier: tt.fields.repo,
			}

			for eventType, mapper := range tt.fields.eventMapper {
				RegisterFilterEventMapper("test", eventType, mapper)
			}
			if len(eventInterceptors) != len(tt.fields.eventMapper) {
				t.Errorf("register event mapper failed expected mapper amount: %d, got: %d", len(tt.fields.eventMapper), len(eventInterceptors))
				t.FailNow()
			}

			_, err := es.Filter(context.Background(), tt.args.query)
			if (err != nil) != tt.res.wantErr {
				t.Errorf("Eventstore.aggregatesToEvents() error = %v, wantErr %v", err, tt.res.wantErr)
			}
		})
	}
}

func TestEventstore_LatestSequence(t *testing.T) {
	type args struct {
		query *SearchQueryBuilder
	}
	type fields struct {
		repo *testQuerier
	}
	type res struct {
		wantErr bool
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		res    res
	}{
		{
			name: "no events",
			args: args{
				query: &SearchQueryBuilder{
					columns: ColumnsMaxSequence,
					queries: []*SearchQuery{
						{
							builder:        &SearchQueryBuilder{},
							aggregateTypes: []AggregateType{"no.aggregates"},
						},
					},
				},
			},
			fields: fields{
				repo: &testQuerier{
					events: []Event{},
					t:      t,
				},
			},
			res: res{
				wantErr: false,
			},
		},
		{
			name: "repo error",
			args: args{
				query: &SearchQueryBuilder{
					columns: ColumnsMaxSequence,
					queries: []*SearchQuery{
						{
							builder:        &SearchQueryBuilder{},
							aggregateTypes: []AggregateType{"no.aggregates"},
						},
					},
				},
			},
			fields: fields{
				repo: &testQuerier{
					t:   t,
					err: zerrors.ThrowInternal(nil, "V2-RfkBa", "test err"),
				},
			},
			res: res{
				wantErr: true,
			},
		},
		{
			name: "found events",
			args: args{
				query: &SearchQueryBuilder{
					columns: ColumnsMaxSequence,
					queries: []*SearchQuery{
						{
							builder:        &SearchQueryBuilder{},
							aggregateTypes: []AggregateType{"test.aggregate"},
						},
					},
				},
			},
			fields: fields{
				repo: &testQuerier{
					// sequence: time.Now(),
					t: t,
				},
			},
			res: res{
				wantErr: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &Eventstore{
				querier: tt.fields.repo,
			}

			_, err := es.LatestSequence(context.Background(), tt.args.query)
			if (err != nil) != tt.res.wantErr {
				t.Errorf("Eventstore.aggregatesToEvents() error = %v, wantErr %v", err, tt.res.wantErr)
			}
		})
	}
}

type testReducer struct {
	t              *testing.T
	events         []Event
	expectedLength int
	err            error
}

func (r *testReducer) Reduce() error {
	r.t.Helper()
	if len(r.events) != r.expectedLength {
		r.t.Errorf("wrong amount of append events wanted: %d, got %d", r.expectedLength, len(r.events))
	}
	if r.err != nil {
		return r.err
	}
	return nil
}

func (r *testReducer) AppendEvents(e ...Event) {
	r.events = append(r.events, e...)
}

func TestEventstore_FilterToReducer(t *testing.T) {
	type args struct {
		query     *SearchQueryBuilder
		readModel reducer
	}
	type fields struct {
		repo        *testQuerier
		eventMapper map[EventType]func(Event) (Event, error)
	}
	type res struct {
		wantErr bool
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		res    res
	}{
		{
			name: "no events",
			args: args{
				query: &SearchQueryBuilder{
					columns: ColumnsEvent,
					queries: []*SearchQuery{
						{
							builder:        &SearchQueryBuilder{},
							aggregateTypes: []AggregateType{"no.aggregates"},
						},
					},
				},
				readModel: &testReducer{
					t:              t,
					expectedLength: 0,
				},
			},
			fields: fields{
				repo: &testQuerier{
					events: []Event{},
					t:      t,
				},
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(e Event) (Event, error) {
						return &testEvent{}, nil
					},
				},
			},
			res: res{
				wantErr: false,
			},
		},
		{
			name: "repo error",
			args: args{
				query: &SearchQueryBuilder{
					columns: ColumnsEvent,
					queries: []*SearchQuery{
						{
							builder:        &SearchQueryBuilder{},
							aggregateTypes: []AggregateType{"no.aggregates"},
						},
					},
				},
				readModel: &testReducer{
					t:              t,
					expectedLength: 0,
				},
			},
			fields: fields{
				repo: &testQuerier{
					t:   t,
					err: zerrors.ThrowInternal(nil, "V2-RfkBa", "test err"),
				},
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(e Event) (Event, error) {
						return &testEvent{}, nil
					},
				},
			},
			res: res{
				wantErr: true,
			},
		},
		{
			name: "found events",
			args: args{
				query: &SearchQueryBuilder{
					columns: ColumnsEvent,
					queries: []*SearchQuery{
						{
							builder:        &SearchQueryBuilder{},
							aggregateTypes: []AggregateType{"test.aggregate"},
						},
					},
				},
				readModel: &testReducer{
					t:              t,
					expectedLength: 1,
				},
			},
			fields: fields{
				repo: &testQuerier{
					events: []Event{
						&BaseEvent{
							Agg: &Aggregate{
								ID: "test.aggregate",
							},
							EventType: "test.event",
						},
					},
					t: t,
				},
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(e Event) (Event, error) {
						return &testEvent{}, nil
					},
				},
			},
		},
		{
			name: "append in reducer fails",
			args: args{
				query: &SearchQueryBuilder{
					columns: ColumnsEvent,
					queries: []*SearchQuery{
						{
							builder:        &SearchQueryBuilder{},
							aggregateTypes: []AggregateType{"test.aggregate"},
						},
					},
				},
				readModel: &testReducer{
					t:              t,
					err:            zerrors.ThrowInvalidArgument(nil, "V2-W06TG", "test err"),
					expectedLength: 1,
				},
			},
			fields: fields{
				repo: &testQuerier{
					events: []Event{
						&BaseEvent{
							Agg: &Aggregate{
								ID: "test.aggregate",
							},
							EventType: "test.event",
						},
					},
					t: t,
				},
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(e Event) (Event, error) {
						return &testEvent{}, nil
					},
				},
			},
			res: res{
				wantErr: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &Eventstore{
				querier: tt.fields.repo,
			}
			for eventType, mapper := range tt.fields.eventMapper {
				RegisterFilterEventMapper("test", eventType, mapper)
			}
			if len(eventInterceptors) != len(tt.fields.eventMapper) {
				t.Errorf("register event mapper failed expected mapper amount: %d, got: %d", len(tt.fields.eventMapper), len(eventInterceptors))
				t.FailNow()
			}

			err := es.FilterToReducer(context.Background(), tt.args.query, tt.args.readModel)
			if (err != nil) != tt.res.wantErr {
				t.Errorf("Eventstore.aggregatesToEvents() error = %v, wantErr %v", err, tt.res.wantErr)
			}
		})
	}
}

func combineEventLists(lists ...[]Event) []Event {
	events := []Event{}
	for _, list := range lists {
		events = append(events, list...)
	}
	return events
}

func compareEvents(t *testing.T, want Event, got Command) {
	t.Helper()

	if want.Aggregate().ID != got.Aggregate().ID {
		t.Errorf("wrong aggregateID got %q want %q", got.Aggregate().ID, want.Aggregate().ID)
	}
	if want.Aggregate().Type != got.Aggregate().Type {
		t.Errorf("wrong aggregateType got %q want %q", got.Aggregate().Type, want.Aggregate().Type)
	}
	if !reflect.DeepEqual(want.DataAsBytes(), got.Payload()) {
		t.Errorf("wrong data got %s want %s", string(got.Payload().([]byte)), string(want.DataAsBytes()))
	}
	if want.Creator() != got.Creator() {
		t.Errorf("wrong creator got %q want %q", got.Creator(), want.Creator())
	}
	if want.Aggregate().ResourceOwner != got.Aggregate().ResourceOwner {
		t.Errorf("wrong resource owner got %q want %q", got.Aggregate().ResourceOwner, want.Aggregate().ResourceOwner)
	}
	if want.Type() != got.Type() {
		t.Errorf("wrong event type got %q want %q", got.Type(), want.Type())
	}
	if want.Revision() != got.Revision() {
		t.Errorf("wrong version got %q want %q", got.Revision(), want.Revision())
	}
}

func TestEventstore_mapEvents(t *testing.T) {
	type fields struct {
		eventMapper map[EventType]func(Event) (Event, error)
	}
	type args struct {
		events []Event
	}
	type res struct {
		events  []Event
		wantErr bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "mapping failed",
			args: args{
				events: []Event{
					&BaseEvent{
						EventType: "test.event",
					},
				},
			},
			fields: fields{
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(Event) (Event, error) {
						return nil, zerrors.ThrowInternal(nil, "V2-8FbQk", "test err")
					},
				},
			},
			res: res{
				wantErr: true,
			},
		},
		{
			name: "mapping succeeded",
			args: args{
				events: []Event{
					&BaseEvent{
						EventType: "test.event",
					},
				},
			},
			fields: fields{
				eventMapper: map[EventType]func(Event) (Event, error){
					"test.event": func(Event) (Event, error) {
						return &testEvent{}, nil
					},
				},
			},
			res: res{
				events: []Event{
					&testEvent{},
				},
				wantErr: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &Eventstore{}
			for eventType, mapper := range tt.fields.eventMapper {
				RegisterFilterEventMapper("test", eventType, mapper)
			}
			if len(eventInterceptors) != len(tt.fields.eventMapper) {
				t.Errorf("register event mapper failed expected mapper amount: %d, got: %d", len(tt.fields.eventMapper), len(eventInterceptors))
				t.FailNow()
			}

			gotMappedEvents, err := es.mapEvents(tt.args.events)
			if (err != nil) != tt.res.wantErr {
				t.Errorf("Eventstore.mapEvents() error = %v, wantErr %v", err, tt.res.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMappedEvents, tt.res.events) {
				t.Errorf("Eventstore.mapEvents() = %v, want %v", gotMappedEvents, tt.res.events)
			}
		})
	}
}
