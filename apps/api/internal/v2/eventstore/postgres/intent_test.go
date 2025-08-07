package postgres

import (
	"testing"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

func Test_checkSequences(t *testing.T) {
	type args struct {
		intents []*intent
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ignore",
			args: args{
				intents: []*intent{
					{
						sequence: 1,
						PushAggregate: eventstore.NewPushAggregate(
							"", "", "",
							eventstore.IgnoreCurrentSequence(),
						),
					},
				},
			},
			want: true,
		},
		{
			name: "ignores",
			args: args{
				intents: []*intent{
					{
						sequence: 1,
						PushAggregate: eventstore.NewPushAggregate(
							"", "", "",
							eventstore.IgnoreCurrentSequence(),
						),
					},
					{
						sequence: 1,
						PushAggregate: eventstore.NewPushAggregate(
							"", "", "",
						),
					},
				},
			},
			want: true,
		},
		{
			name: "matches",
			args: args{
				intents: []*intent{
					{
						sequence: 0,
						PushAggregate: eventstore.NewPushAggregate(
							"", "", "",
							eventstore.CurrentSequenceMatches(0),
						),
					},
				},
			},
			want: true,
		},
		{
			name: "does not match",
			args: args{
				intents: []*intent{
					{
						sequence: 1,
						PushAggregate: eventstore.NewPushAggregate(
							"", "", "",
							eventstore.CurrentSequenceMatches(2),
						),
					},
				},
			},
			want: false,
		},
		{
			name: "at least",
			args: args{
				intents: []*intent{
					{
						sequence: 10,
						PushAggregate: eventstore.NewPushAggregate(
							"", "", "",
							eventstore.CurrentSequenceAtLeast(0),
						),
					},
				},
			},
			want: true,
		},
		{
			name: "at least too low",
			args: args{
				intents: []*intent{
					{
						sequence: 1,
						PushAggregate: eventstore.NewPushAggregate(
							"", "", "",
							eventstore.CurrentSequenceAtLeast(2),
						),
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkSequences(tt.args.intents); got != tt.want {
				t.Errorf("checkSequences() = %v, want %v", got, tt.want)
			}
		})
	}
}
