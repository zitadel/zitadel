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
						PushIntent: &testIntent{
							currentSequence: eventstore.SequenceIgnore(),
						},
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
						PushIntent: &testIntent{
							currentSequence: eventstore.SequenceIgnore(),
						},
					},
					{
						sequence: 1,
						PushIntent: &testIntent{
							currentSequence: nil,
						},
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
						PushIntent: &testIntent{
							currentSequence: eventstore.SequenceMatches(0),
						},
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
						PushIntent: &testIntent{
							currentSequence: eventstore.SequenceMatches(2),
						},
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
						PushIntent: &testIntent{
							currentSequence: eventstore.SequenceAtLeast(0),
						},
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
						PushIntent: &testIntent{
							currentSequence: eventstore.SequenceAtLeast(2),
						},
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
