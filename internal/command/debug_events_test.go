package command

import (
	"io"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/debug_events"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_CreateDebugEvents(t *testing.T) {
	ctx := authz.NewMockContextWithPermissions("instance1", "org1", "user1", nil)
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		dbe *DebugEvents
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr error
	}{
		{
			name: "filter error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{&DebugEvents{
				AggregateID: "dbg1",
				Events: []DebugEvent{
					DebugEventAdded{
						ProjectionSleep: time.Millisecond,
						Blob:            gu.Ptr("a"),
					},
				},
			}},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "already exists",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							debug_events.NewAddedEvent(
								ctx, debug_events.NewAggregate("dbg1", "instance1"),
								time.Millisecond, gu.Ptr("a"),
							),
						),
					),
				),
			},
			args: args{&DebugEvents{
				AggregateID: "dbg1",
				Events: []DebugEvent{
					DebugEventAdded{
						ProjectionSleep: time.Millisecond,
						Blob:            gu.Ptr("a"),
					},
				},
			}},
			wantErr: zerrors.ThrowAlreadyExists(nil, "COMMAND-Aex6j", "debug aggregate already exists"),
		},
		{
			name: "double added event, already exists",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{&DebugEvents{
				AggregateID: "dbg1",
				Events: []DebugEvent{
					DebugEventAdded{
						ProjectionSleep: time.Millisecond,
						Blob:            gu.Ptr("a"),
					},
					DebugEventAdded{
						ProjectionSleep: time.Millisecond,
						Blob:            gu.Ptr("a"),
					},
				},
			}},
			wantErr: zerrors.ThrowAlreadyExists(nil, "COMMAND-Aex6j", "debug aggregate already exists"),
		},
		{
			name: "changed event, not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{&DebugEvents{
				AggregateID: "dbg1",
				Events: []DebugEvent{
					DebugEventChanged{
						ProjectionSleep: time.Millisecond,
						Blob:            gu.Ptr("a"),
					},
				},
			}},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-Thie6", "debug aggregate not found"),
		},
		{
			name: "removed event, not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{&DebugEvents{
				AggregateID: "dbg1",
				Events: []DebugEvent{
					DebugEventRemoved{
						ProjectionSleep: time.Millisecond,
					},
				},
			}},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-Ohna9", "debug aggregate not found"),
		},
		{
			name: "changed after removed event, not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							debug_events.NewAddedEvent(
								ctx, debug_events.NewAggregate("dbg1", "instance1"),
								time.Millisecond, gu.Ptr("a"),
							),
						),
					),
				),
			},
			args: args{&DebugEvents{
				AggregateID: "dbg1",
				Events: []DebugEvent{
					DebugEventRemoved{
						ProjectionSleep: time.Millisecond,
					},
					DebugEventChanged{
						ProjectionSleep: time.Millisecond,
						Blob:            gu.Ptr("a"),
					},
				},
			}},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-Thie6", "debug aggregate not found"),
		},
		{
			name: "double removed event, not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							debug_events.NewAddedEvent(
								ctx, debug_events.NewAggregate("dbg1", "instance1"),
								time.Millisecond, gu.Ptr("a"),
							),
						),
					),
				),
			},
			args: args{&DebugEvents{
				AggregateID: "dbg1",
				Events: []DebugEvent{
					DebugEventRemoved{
						ProjectionSleep: time.Millisecond,
					},
					DebugEventRemoved{
						ProjectionSleep: time.Millisecond,
					},
				},
			}},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-Ohna9", "debug aggregate not found"),
		},
		{
			name: "added, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						debug_events.NewAddedEvent(
							ctx, debug_events.NewAggregate("dbg1", "instance1"),
							time.Millisecond, gu.Ptr("a"),
						),
					),
				),
			},
			args: args{&DebugEvents{
				AggregateID: "dbg1",
				Events: []DebugEvent{
					DebugEventAdded{
						ProjectionSleep: time.Millisecond,
						Blob:            gu.Ptr("a"),
					},
				},
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "changed, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							debug_events.NewAddedEvent(
								ctx, debug_events.NewAggregate("dbg1", "instance1"),
								time.Millisecond, gu.Ptr("a"),
							),
						),
					),
					expectPush(
						debug_events.NewChangedEvent(
							ctx, debug_events.NewAggregate("dbg1", "instance1"),
							time.Millisecond, gu.Ptr("b"),
						),
					),
				),
			},
			args: args{&DebugEvents{
				AggregateID: "dbg1",
				Events: []DebugEvent{
					DebugEventChanged{
						ProjectionSleep: time.Millisecond,
						Blob:            gu.Ptr("b"),
					},
				},
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "removed, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							debug_events.NewAddedEvent(
								ctx, debug_events.NewAggregate("dbg1", "instance1"),
								time.Millisecond, gu.Ptr("a"),
							),
						),
					),
					expectPush(
						debug_events.NewRemovedEvent(
							ctx, debug_events.NewAggregate("dbg1", "instance1"),
							time.Millisecond,
						),
					),
				),
			},
			args: args{&DebugEvents{
				AggregateID: "dbg1",
				Events: []DebugEvent{
					DebugEventRemoved{
						ProjectionSleep: time.Millisecond,
					},
				},
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "added, changed, changed, removed ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						debug_events.NewAddedEvent(
							ctx, debug_events.NewAggregate("dbg1", "instance1"),
							time.Millisecond, gu.Ptr("a"),
						),
						debug_events.NewChangedEvent(
							ctx, debug_events.NewAggregate("dbg1", "instance1"),
							time.Millisecond, gu.Ptr("b"),
						),
						debug_events.NewChangedEvent(
							ctx, debug_events.NewAggregate("dbg1", "instance1"),
							time.Millisecond, gu.Ptr("c"),
						),
						debug_events.NewRemovedEvent(
							ctx, debug_events.NewAggregate("dbg1", "instance1"),
							time.Millisecond,
						),
					),
				),
			},
			args: args{&DebugEvents{
				AggregateID: "dbg1",
				Events: []DebugEvent{
					DebugEventAdded{
						ProjectionSleep: time.Millisecond,
						Blob:            gu.Ptr("a"),
					},
					DebugEventChanged{
						ProjectionSleep: time.Millisecond,
						Blob:            gu.Ptr("b"),
					},
					DebugEventChanged{
						ProjectionSleep: time.Millisecond,
						Blob:            gu.Ptr("c"),
					},
					DebugEventRemoved{
						ProjectionSleep: time.Millisecond,
					},
				},
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := c.CreateDebugEvents(ctx, tt.args.dbe)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
