package eventstore

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/cockroach"
	"github.com/zitadel/zitadel/internal/eventstore"
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
					"($1, $2, $3, $4, $5, $6, $7, $8, $9, hlc_to_timestamp(cluster_logical_timestamp()), cluster_logical_timestamp(), $10)",
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
					"($1, $2, $3, $4, $5, $6, $7, $8, $9, hlc_to_timestamp(cluster_logical_timestamp()), cluster_logical_timestamp(), $10)",
					"($11, $12, $13, $14, $15, $16, $17, $18, $19, hlc_to_timestamp(cluster_logical_timestamp()), cluster_logical_timestamp(), $20)",
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
					"($1, $2, $3, $4, $5, $6, $7, $8, $9, hlc_to_timestamp(cluster_logical_timestamp()), cluster_logical_timestamp(), $10)",
					"($11, $12, $13, $14, $15, $16, $17, $18, $19, hlc_to_timestamp(cluster_logical_timestamp()), cluster_logical_timestamp(), $20)",
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
		NewEventstore(&database.DB{Database: new(cockroach.Config)})
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
