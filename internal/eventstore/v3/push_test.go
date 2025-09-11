package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"testing"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/postgres"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/mock"
	"github.com/zitadel/zitadel/internal/execution/target"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
)

func Test_mapCommands(t *testing.T) {
	type args struct {
		commands  []eventstore.Command
		sequences []*latestSequence
	}
	type want struct {
		events       []eventstore.Event
		placeHolders []string
		args         []any
		err          func(t *testing.T, err error)
		shouldPanic  bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no commands",
			args: args{
				commands:  []eventstore.Command{},
				sequences: []*latestSequence{},
			},
			want: want{
				events:       []eventstore.Event{},
				placeHolders: []string{},
				args:         []any{},
			},
		},
		{
			name: "one command",
			args: args{
				commands: []eventstore.Command{
					&mockCommand{
						aggregate: mockAggregate("V3-VEIvq"),
					},
				},
				sequences: []*latestSequence{
					{
						aggregate: mockAggregate("V3-VEIvq"),
						sequence:  0,
					},
				},
			},
			want: want{
				events: []eventstore.Event{
					mockEvent(
						mockAggregate("V3-VEIvq"),
						1,
						nil,
					),
				},
				placeHolders: []string{
					"($1, $2, $3, $4, $5, $6, $7, $8, $9, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp()), $10)",
				},
				args: []any{
					"instance",
					"ro",
					"type",
					"V3-VEIvq",
					uint16(1),
					"creator",
					"event.type",
					Payload(nil),
					uint64(1),
					0,
				},
				err: func(t *testing.T, err error) {},
			},
		},
		{
			name: "multiple commands same aggregate",
			args: args{
				commands: []eventstore.Command{
					&mockCommand{
						aggregate: mockAggregate("V3-VEIvq"),
					},
					&mockCommand{
						aggregate: mockAggregate("V3-VEIvq"),
					},
				},
				sequences: []*latestSequence{
					{
						aggregate: mockAggregate("V3-VEIvq"),
						sequence:  5,
					},
				},
			},
			want: want{
				events: []eventstore.Event{
					mockEvent(
						mockAggregate("V3-VEIvq"),
						6,
						nil,
					),
					mockEvent(
						mockAggregate("V3-VEIvq"),
						7,
						nil,
					),
				},
				placeHolders: []string{
					"($1, $2, $3, $4, $5, $6, $7, $8, $9, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp()), $10)",
					"($11, $12, $13, $14, $15, $16, $17, $18, $19, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp()), $20)",
				},
				args: []any{
					// first event
					"instance",
					"ro",
					"type",
					"V3-VEIvq",
					uint16(1),
					"creator",
					"event.type",
					Payload(nil),
					uint64(6),
					0,
					// second event
					"instance",
					"ro",
					"type",
					"V3-VEIvq",
					uint16(1),
					"creator",
					"event.type",
					Payload(nil),
					uint64(7),
					1,
				},
				err: func(t *testing.T, err error) {},
			},
		},
		{
			name: "one command per aggregate",
			args: args{
				commands: []eventstore.Command{
					&mockCommand{
						aggregate: mockAggregate("V3-VEIvq"),
					},
					&mockCommand{
						aggregate: mockAggregate("V3-IT6VN"),
					},
				},
				sequences: []*latestSequence{
					{
						aggregate: mockAggregate("V3-VEIvq"),
						sequence:  5,
					},
					{
						aggregate: mockAggregate("V3-IT6VN"),
						sequence:  0,
					},
				},
			},
			want: want{
				events: []eventstore.Event{
					mockEvent(
						mockAggregate("V3-VEIvq"),
						6,
						nil,
					),
					mockEvent(
						mockAggregate("V3-IT6VN"),
						1,
						nil,
					),
				},
				placeHolders: []string{
					"($1, $2, $3, $4, $5, $6, $7, $8, $9, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp()), $10)",
					"($11, $12, $13, $14, $15, $16, $17, $18, $19, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp()), $20)",
				},
				args: []any{
					// first event
					"instance",
					"ro",
					"type",
					"V3-VEIvq",
					uint16(1),
					"creator",
					"event.type",
					Payload(nil),
					uint64(6),
					0,
					// second event
					"instance",
					"ro",
					"type",
					"V3-IT6VN",
					uint16(1),
					"creator",
					"event.type",
					Payload(nil),
					uint64(1),
					1,
				},
				err: func(t *testing.T, err error) {},
			},
		},
		{
			name: "missing sequence",
			args: args{
				commands: []eventstore.Command{
					&mockCommand{
						aggregate: mockAggregate("V3-VEIvq"),
					},
				},
				sequences: []*latestSequence{},
			},
			want: want{
				events:       []eventstore.Event{},
				placeHolders: []string{},
				args:         []any{},
				err:          func(t *testing.T, err error) {},
				shouldPanic:  true,
			},
		},
	}
	for _, tt := range tests {
		if tt.want.err == nil {
			tt.want.err = func(t *testing.T, err error) {
				require.NoError(t, err)
			}
		}
		// is used to set the the [pushPlaceholderFmt]
		NewEventstore(&database.DB{Database: new(postgres.Config)})
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				cause := recover()
				assert.Equal(t, tt.want.shouldPanic, cause != nil)
			}()
			gotEvents, gotPlaceHolders, gotArgs, err := mapCommands(tt.args.commands, tt.args.sequences)
			tt.want.err(t, err)

			assert.ElementsMatch(t, tt.want.events, gotEvents)
			assert.ElementsMatch(t, tt.want.placeHolders, gotPlaceHolders)
			assert.ElementsMatch(t, tt.want.args, gotArgs)
		})
	}
}

func TestEventstore_queueExecutions(t *testing.T) {
	events := []eventstore.Event{
		mockEventType(mockAggregate("TEST"), 1, []byte(`{"test":"test"}`), "ex.foo.bar"),
		mockEventType(mockAggregate("TEST"), 2, []byte("{}"), "ex.bar.foo"),
		mockEventType(mockAggregate("TEST"), 3, nil, "ex.removed"),
	}
	type args struct {
		ctx    context.Context
		tx     database.Tx
		events []eventstore.Event
	}
	tests := []struct {
		name    string
		queue   func(t *testing.T) eventstore.ExecutionQueue
		args    args
		wantErr bool
	}{
		{
			name: "incorrect Tx type, noop",
			queue: func(t *testing.T) eventstore.ExecutionQueue {
				mQueue := mock.NewMockExecutionQueue(gomock.NewController(t))
				return mQueue
			},
			args: args{
				ctx:    context.Background(),
				tx:     nil,
				events: events,
			},
			wantErr: false,
		},
		{
			name: "no events",
			queue: func(t *testing.T) eventstore.ExecutionQueue {
				mQueue := mock.NewMockExecutionQueue(gomock.NewController(t))
				return mQueue
			},
			args: args{
				ctx:    context.Background(),
				tx:     &sql.Tx{},
				events: []eventstore.Event{},
			},
			wantErr: false,
		},
		{
			name: "no router in Ctx",
			queue: func(t *testing.T) eventstore.ExecutionQueue {
				mQueue := mock.NewMockExecutionQueue(gomock.NewController(t))
				return mQueue
			},
			args: args{
				ctx:    context.Background(),
				tx:     &sql.Tx{},
				events: events,
			},
			wantErr: false,
		},
		{
			name: "not found in router",
			queue: func(t *testing.T) eventstore.ExecutionQueue {
				mQueue := mock.NewMockExecutionQueue(gomock.NewController(t))
				return mQueue
			},
			args: args{
				ctx: authz.WithExecutionRouter(
					context.Background(),
					target.NewRouter([]target.Target{
						{
							ExecutionID: "function/fooBar",
						},
					}),
				),
				tx:     &sql.Tx{},
				events: events,
			},
			wantErr: false,
		},
		{
			name: "event prefix",
			queue: func(t *testing.T) eventstore.ExecutionQueue {
				mQueue := mock.NewMockExecutionQueue(gomock.NewController(t))
				mQueue.EXPECT().InsertManyFastTx(
					gomock.Any(),
					gomock.Any(),
					[]river.JobArgs{
						mustNewRequest(t, events[0], []target.Target{{ExecutionID: "event"}}),
						mustNewRequest(t, events[1], []target.Target{{ExecutionID: "event"}}),
						mustNewRequest(t, events[2], []target.Target{{ExecutionID: "event"}}),
					},
					gomock.Any(),
				)
				return mQueue
			},
			args: args{
				ctx: authz.WithExecutionRouter(
					context.Background(),
					target.NewRouter([]target.Target{
						{ExecutionID: "function/fooBar"},
						{ExecutionID: "event"},
					}),
				),
				tx:     &sql.Tx{},
				events: events,
			},
			wantErr: false,
		},
		{
			name: "event wildcard and exact match",
			queue: func(t *testing.T) eventstore.ExecutionQueue {
				mQueue := mock.NewMockExecutionQueue(gomock.NewController(t))
				mQueue.EXPECT().InsertManyFastTx(
					gomock.Any(),
					gomock.Any(),
					[]river.JobArgs{
						mustNewRequest(t, events[0], []target.Target{{ExecutionID: "event/ex.foo.*"}}),
						mustNewRequest(t, events[2], []target.Target{{ExecutionID: "event/ex.removed"}}),
					},
					gomock.Any(),
				)
				return mQueue
			},
			args: args{
				ctx: authz.WithExecutionRouter(
					context.Background(),
					target.NewRouter([]target.Target{
						{ExecutionID: "function/fooBar"},
						{ExecutionID: "event/ex.foo.*"},
						{ExecutionID: "event/ex.removed"},
					}),
				),
				tx:     &sql.Tx{},
				events: events,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &Eventstore{
				queue: tt.queue(t),
			}
			err := es.queueExecutions(tt.args.ctx, tt.args.tx, tt.args.events)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func mustNewRequest(t *testing.T, e eventstore.Event, targets []target.Target) *exec_repo.Request {
	req, err := exec_repo.NewRequest(e, targets)
	require.NoError(t, err, "exec_repo.NewRequest")
	return req
}
