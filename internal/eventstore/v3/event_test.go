package eventstore

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
)

func Test_commandToEvent(t *testing.T) {
	payload := struct {
		ID string
	}{
		ID: "test",
	}
	payloadMarshalled, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal of payload failed: %v", err)
	}
	type args struct {
		command eventstore.Command
	}
	type want struct {
		event *event
		err   func(t *testing.T, err error)
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no payload",
			args: args{
				command: &mockCommand{
					aggregate: mockAggregate("V3-Red9I"),
					payload:   nil,
				},
			},
			want: want{
				event: mockEvent(
					mockAggregate("V3-Red9I"),
					0,
					nil,
				).(*event),
			},
		},
		{
			name: "struct payload",
			args: args{
				command: &mockCommand{
					aggregate: mockAggregate("V3-Red9I"),
					payload:   payload,
				},
			},
			want: want{
				event: mockEvent(
					mockAggregate("V3-Red9I"),
					0,
					payloadMarshalled,
				).(*event),
			},
		},
		{
			name: "pointer payload",
			args: args{
				command: &mockCommand{
					aggregate: mockAggregate("V3-Red9I"),
					payload:   &payload,
				},
			},
			want: want{
				event: mockEvent(
					mockAggregate("V3-Red9I"),
					0,
					payloadMarshalled,
				).(*event),
			},
		},
		{
			name: "invalid payload",
			args: args{
				command: &mockCommand{
					aggregate: mockAggregate("V3-Red9I"),
					payload:   func() {},
				},
			},
			want: want{
				err: func(t *testing.T, err error) {
					assert.Error(t, err)
				},
			},
		},
	}
	for _, tt := range tests {
		if tt.want.err == nil {
			tt.want.err = func(t *testing.T, err error) {
				require.NoError(t, err)
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := commandToEvent(tt.args.command)

			tt.want.err(t, err)
			if tt.want.event == nil {
				assert.Nil(t, got)
				return
			}
			assert.Equal(t, tt.want.event, got)
		})
	}
}

func Test_commandToEventOld(t *testing.T) {
	payload := struct {
		ID string
	}{
		ID: "test",
	}
	payloadMarshalled, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal of payload failed: %v", err)
	}
	type args struct {
		sequence *latestSequence
		command  eventstore.Command
	}
	type want struct {
		event *event
		err   func(t *testing.T, err error)
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no payload",
			args: args{
				sequence: &latestSequence{
					aggregate: mockAggregate("V3-Red9I"),
					sequence:  0,
				},
				command: &mockCommand{
					aggregate: mockAggregate("V3-Red9I"),
					payload:   nil,
				},
			},
			want: want{
				event: mockEvent(
					mockAggregate("V3-Red9I"),
					0,
					nil,
				).(*event),
			},
		},
		{
			name: "struct payload",
			args: args{
				sequence: &latestSequence{
					aggregate: mockAggregate("V3-Red9I"),
					sequence:  0,
				},
				command: &mockCommand{
					aggregate: mockAggregate("V3-Red9I"),
					payload:   payload,
				},
			},
			want: want{
				event: mockEvent(
					mockAggregate("V3-Red9I"),
					0,
					payloadMarshalled,
				).(*event),
			},
		},
		{
			name: "pointer payload",
			args: args{
				sequence: &latestSequence{
					aggregate: mockAggregate("V3-Red9I"),
					sequence:  0,
				},
				command: &mockCommand{
					aggregate: mockAggregate("V3-Red9I"),
					payload:   &payload,
				},
			},
			want: want{
				event: mockEvent(
					mockAggregate("V3-Red9I"),
					0,
					payloadMarshalled,
				).(*event),
			},
		},
		{
			name: "invalid payload",
			args: args{
				sequence: &latestSequence{
					aggregate: mockAggregate("V3-Red9I"),
					sequence:  0,
				},
				command: &mockCommand{
					aggregate: mockAggregate("V3-Red9I"),
					payload:   func() {},
				},
			},
			want: want{
				err: func(t *testing.T, err error) {
					assert.Error(t, err)
				},
			},
		},
	}
	for _, tt := range tests {
		if tt.want.err == nil {
			tt.want.err = func(t *testing.T, err error) {
				require.NoError(t, err)
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := commandToEventOld(tt.args.sequence, tt.args.command)

			tt.want.err(t, err)
			assert.Equal(t, tt.want.event, got)
		})
	}
}

func Test_commandsToEvents(t *testing.T) {
	ctx := context.Background()
	payload := struct {
		ID string
	}{
		ID: "test",
	}
	payloadMarshalled, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal of payload failed: %v", err)
	}
	type args struct {
		ctx  context.Context
		cmds []eventstore.Command
	}
	type want struct {
		events   []eventstore.Event
		commands []*command
		err      func(t *testing.T, err error)
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no commands",
			args: args{
				ctx:  ctx,
				cmds: nil,
			},
			want: want{
				events:   []eventstore.Event{},
				commands: []*command{},
				err: func(t *testing.T, err error) {
					require.NoError(t, err)
				},
			},
		},
		{
			name: "single command no payload",
			args: args{
				ctx: ctx,
				cmds: []eventstore.Command{
					&mockCommand{
						aggregate: mockAggregate("V3-Red9I"),
						payload:   nil,
					},
				},
			},
			want: want{
				events: []eventstore.Event{
					mockEvent(
						mockAggregate("V3-Red9I"),
						0,
						nil,
					),
				},
				commands: []*command{
					{
						InstanceID:    "instance",
						AggregateType: "type",
						AggregateID:   "V3-Red9I",
						Owner:         "ro",
						CommandType:   "event.type",
						Revision:      1,
						Payload:       nil,
						Creator:       "creator",
					},
				},
				err: func(t *testing.T, err error) {
					require.NoError(t, err)
				},
			},
		},
		{
			name: "single command no instance id",
			args: args{
				ctx: authz.WithInstanceID(ctx, "instance from ctx"),
				cmds: []eventstore.Command{
					&mockCommand{
						aggregate: mockAggregateWithInstance("V3-Red9I", ""),
						payload:   nil,
					},
				},
			},
			want: want{
				events: []eventstore.Event{
					mockEvent(
						mockAggregateWithInstance("V3-Red9I", "instance from ctx"),
						0,
						nil,
					),
				},
				commands: []*command{
					{
						InstanceID:    "instance from ctx",
						AggregateType: "type",
						AggregateID:   "V3-Red9I",
						Owner:         "ro",
						CommandType:   "event.type",
						Revision:      1,
						Payload:       nil,
						Creator:       "creator",
					},
				},
				err: func(t *testing.T, err error) {
					require.NoError(t, err)
				},
			},
		},
		{
			name: "single command with payload",
			args: args{
				ctx: ctx,
				cmds: []eventstore.Command{
					&mockCommand{
						aggregate: mockAggregate("V3-Red9I"),
						payload:   payload,
					},
				},
			},
			want: want{
				events: []eventstore.Event{
					mockEvent(
						mockAggregate("V3-Red9I"),
						0,
						payloadMarshalled,
					),
				},
				commands: []*command{
					{
						InstanceID:    "instance",
						AggregateType: "type",
						AggregateID:   "V3-Red9I",
						Owner:         "ro",
						CommandType:   "event.type",
						Revision:      1,
						Payload:       payloadMarshalled,
						Creator:       "creator",
					},
				},
				err: func(t *testing.T, err error) {
					require.NoError(t, err)
				},
			},
		},
		{
			name: "multiple commands",
			args: args{
				ctx: ctx,
				cmds: []eventstore.Command{
					&mockCommand{
						aggregate: mockAggregate("V3-Red9I"),
						payload:   payload,
					},
					&mockCommand{
						aggregate: mockAggregate("V3-Red9I"),
						payload:   nil,
					},
				},
			},
			want: want{
				events: []eventstore.Event{
					mockEvent(
						mockAggregate("V3-Red9I"),
						0,
						payloadMarshalled,
					),
					mockEvent(
						mockAggregate("V3-Red9I"),
						0,
						nil,
					),
				},
				commands: []*command{
					{
						InstanceID:    "instance",
						AggregateType: "type",
						AggregateID:   "V3-Red9I",
						CommandType:   "event.type",
						Revision:      1,
						Payload:       payloadMarshalled,
						Creator:       "creator",
						Owner:         "ro",
					},
					{
						InstanceID:    "instance",
						AggregateType: "type",
						AggregateID:   "V3-Red9I",
						CommandType:   "event.type",
						Revision:      1,
						Payload:       nil,
						Creator:       "creator",
						Owner:         "ro",
					},
				},
				err: func(t *testing.T, err error) {
					require.NoError(t, err)
				},
			},
		},
		{
			name: "invalid command",
			args: args{
				ctx: ctx,
				cmds: []eventstore.Command{
					&mockCommand{
						aggregate: mockAggregate("V3-Red9I"),
						payload:   func() {},
					},
				},
			},
			want: want{
				events:   nil,
				commands: nil,
				err: func(t *testing.T, err error) {
					assert.Error(t, err)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEvents, gotCommands, err := commandsToEvents(tt.args.ctx, tt.args.cmds)

			tt.want.err(t, err)
			assert.Equal(t, tt.want.events, gotEvents)
			require.Len(t, gotCommands, len(tt.want.commands))
			for i, wantCommand := range tt.want.commands {
				assertCommand(t, wantCommand, gotCommands[i])
			}
		})
	}
}

func assertCommand(t *testing.T, want, got *command) {
	t.Helper()
	assert.Equal(t, want.CommandType, got.CommandType)
	assert.Equal(t, want.Payload, got.Payload)
	assert.Equal(t, want.Creator, got.Creator)
	assert.Equal(t, want.Owner, got.Owner)
	assert.Equal(t, want.AggregateID, got.AggregateID)
	assert.Equal(t, want.AggregateType, got.AggregateType)
	assert.Equal(t, want.InstanceID, got.InstanceID)
	assert.Equal(t, want.Revision, got.Revision)
}
