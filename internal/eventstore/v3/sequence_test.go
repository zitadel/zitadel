package eventstore

import (
	"context"
	_ "embed"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
)

func Test_searchSequence(t *testing.T) {
	sequence := &latestSequence{
		aggregate: mockAggregate("V3-p1BWC"),
		sequence:  1,
	}
	type args struct {
		sequences     []*latestSequence
		aggregateType eventstore.AggregateType
		aggregateID   string
		instanceID    string
	}
	tests := []struct {
		name string
		args args
		want *latestSequence
	}{
		{
			name: "type missmatch",
			args: args{
				sequences: []*latestSequence{
					sequence,
				},
				aggregateType: "wrong",
				aggregateID:   "V3-p1BWC",
				instanceID:    "instance",
			},
			want: nil,
		},
		{
			name: "id missmatch",
			args: args{
				sequences: []*latestSequence{
					sequence,
				},
				aggregateType: "type",
				aggregateID:   "wrong",
				instanceID:    "instance",
			},
			want: nil,
		},
		{
			name: "instance missmatch",
			args: args{
				sequences: []*latestSequence{
					sequence,
				},
				aggregateType: "type",
				aggregateID:   "V3-p1BWC",
				instanceID:    "wrong",
			},
			want: nil,
		},
		{
			name: "match",
			args: args{
				sequences: []*latestSequence{
					sequence,
				},
				aggregateType: "type",
				aggregateID:   "V3-p1BWC",
				instanceID:    "instance",
			},
			want: sequence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchSequence(tt.args.sequences, tt.args.aggregateType, tt.args.aggregateID, tt.args.instanceID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commandsToSequences(t *testing.T) {
	aggregate := mockAggregate("V3-MKHTF")
	type args struct {
		ctx      context.Context
		commands []eventstore.Command
	}
	tests := []struct {
		name string
		args args
		want []*latestSequence
	}{
		{
			name: "no command",
			args: args{
				ctx:      context.Background(),
				commands: []eventstore.Command{},
			},
			want: []*latestSequence{},
		},
		{
			name: "one command",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					&mockCommand{
						aggregate: aggregate,
					},
				},
			},
			want: []*latestSequence{
				{
					aggregate: aggregate,
				},
			},
		},
		{
			name: "two commands same aggregate",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					&mockCommand{
						aggregate: aggregate,
					},
					&mockCommand{
						aggregate: aggregate,
					},
				},
			},
			want: []*latestSequence{
				{
					aggregate: aggregate,
				},
			},
		},
		{
			name: "two commands different aggregates",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					&mockCommand{
						aggregate: aggregate,
					},
					&mockCommand{
						aggregate: mockAggregate("V3-cZkCy"),
					},
				},
			},
			want: []*latestSequence{
				{
					aggregate: aggregate,
				},
				{
					aggregate: mockAggregate("V3-cZkCy"),
				},
			},
		},
		{
			name: "instance set in command",
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "V3-ANV4p"),
				commands: []eventstore.Command{
					&mockCommand{
						aggregate: &eventstore.Aggregate{
							ID:            "V3-bF0Sa",
							Type:          "type",
							ResourceOwner: "to",
							InstanceID:    "instance",
							Version:       "v1",
						},
					},
				},
			},
			want: []*latestSequence{
				{
					aggregate: &eventstore.Aggregate{
						ID:            "V3-bF0Sa",
						Type:          "type",
						ResourceOwner: "to",
						InstanceID:    "instance",
						Version:       "v1",
					},
				},
			},
		},
		{
			name: "instance from context",
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "V3-ANV4p"),
				commands: []eventstore.Command{
					&mockCommand{
						aggregate: &eventstore.Aggregate{
							ID:            "V3-bF0Sa",
							Type:          "type",
							ResourceOwner: "to",
							Version:       "v1",
						},
					},
				},
			},
			want: []*latestSequence{
				{
					aggregate: &eventstore.Aggregate{
						ID:            "V3-bF0Sa",
						Type:          "type",
						ResourceOwner: "to",
						InstanceID:    "V3-ANV4p",
						Version:       "v1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := commandsToSequences(tt.args.ctx, tt.args.commands)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func Test_sequencesToSql(t *testing.T) {
	tests := []struct {
		name           string
		arg            []*latestSequence
		wantConditions []string
		wantArgs       []any
	}{
		{
			name:           "no sequence",
			arg:            []*latestSequence{},
			wantConditions: []string{},
			wantArgs:       []any{},
		},
		{
			name: "one",
			arg: []*latestSequence{
				{
					aggregate: mockAggregate("V3-SbpGB"),
				},
			},
			wantConditions: []string{
				`(SELECT instance_id, aggregate_type, aggregate_id, "sequence" FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 AND aggregate_id = $3 ORDER BY "sequence" DESC LIMIT 1)`,
			},
			wantArgs: []any{
				"instance",
				eventstore.AggregateType("type"),
				"V3-SbpGB",
			},
		},
		{
			name: "multiple",
			arg: []*latestSequence{
				{
					aggregate: mockAggregate("V3-SbpGB"),
				},
				{
					aggregate: mockAggregate("V3-0X3yt"),
				},
			},
			wantConditions: []string{
				`(SELECT instance_id, aggregate_type, aggregate_id, "sequence" FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 AND aggregate_id = $3 ORDER BY "sequence" DESC LIMIT 1)`,
				`(SELECT instance_id, aggregate_type, aggregate_id, "sequence" FROM eventstore.events2 WHERE instance_id = $4 AND aggregate_type = $5 AND aggregate_id = $6 ORDER BY "sequence" DESC LIMIT 1)`,
			},
			wantArgs: []any{
				"instance",
				eventstore.AggregateType("type"),
				"V3-SbpGB",
				"instance",
				eventstore.AggregateType("type"),
				"V3-0X3yt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConditions, gotArgs := sequencesToSql(tt.arg)
			if !reflect.DeepEqual(gotConditions, tt.wantConditions) {
				t.Errorf("sequencesToSql() gotConditions = %v, want %v", gotConditions, tt.wantConditions)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("sequencesToSql() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
