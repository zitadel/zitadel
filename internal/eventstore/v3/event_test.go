package eventstore

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
			got, err := commandToEvent(tt.args.sequence, tt.args.command)

			tt.want.err(t, err)
			assert.Equal(t, tt.want.event, got)
		})
	}
}
