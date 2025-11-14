package projection

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestMilestonesProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.Event
	}
	now := time.Now()
	date, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	require.NoError(t, err)
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.Event) (*handler.Statement, error)
		want   wantReduce
	}{
		{
			name: "reduceInstanceAdded",
			args: args{
				event: getEvent(timedTestEvent(
					milestone.ReachedEventType,
					milestone.AggregateType,
					[]byte(`{"type": "InstanceCreated"}`),
					now,
					withVersion(milestone.AggregateVersion),
				), milestone.ReachedEventMapper),
			},
			reduce: (&milestoneProjection{}).reduceReached,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("milestone"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.milestones3 (instance_id, type, reached_date) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								"instance-id",
								milestone.InstanceCreated,
								now,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceInstanceAdded with reached date",
			args: args{
				event: getEvent(timedTestEvent(
					milestone.ReachedEventType,
					milestone.AggregateType,
					[]byte(`{"type": "InstanceCreated", "reachedDate":"2006-01-02T15:04:05Z"}`),
					now,
					withVersion(milestone.AggregateVersion),
				), milestone.ReachedEventMapper),
			},
			reduce: (&milestoneProjection{}).reduceReached,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("milestone"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.milestones3 (instance_id, type, reached_date) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								"instance-id",
								milestone.InstanceCreated,
								date,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMilestonePushed normal milestone",
			args: args{
				event: getEvent(timedTestEvent(
					milestone.PushedEventType,
					milestone.AggregateType,
					[]byte(`{"type": "ProjectCreated"}`),
					now,
					withVersion(milestone.AggregateVersion),
				), milestone.PushedEventMapper),
			},
			reduce: (&milestoneProjection{}).reducePushed,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("milestone"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.milestones3 SET last_pushed_date = $1 WHERE (instance_id = $2) AND (type = $3)",
							expectedArgs: []interface{}{
								now,
								"instance-id",
								milestone.ProjectCreated,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMilestonePushed normal milestone with pushed date",
			args: args{
				event: getEvent(timedTestEvent(
					milestone.PushedEventType,
					milestone.AggregateType,
					[]byte(`{"type": "ProjectCreated", "pushedDate":"2006-01-02T15:04:05Z"}`),
					now,
					withVersion(milestone.AggregateVersion),
				), milestone.PushedEventMapper),
			},
			reduce: (&milestoneProjection{}).reducePushed,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("milestone"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.milestones3 SET last_pushed_date = $1 WHERE (instance_id = $2) AND (type = $3)",
							expectedArgs: []interface{}{
								date,
								"instance-id",
								milestone.ProjectCreated,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMilestonePushed instance deleted milestone",
			args: args{
				event: getEvent(testEvent(
					milestone.PushedEventType,
					milestone.AggregateType,
					[]byte(`{"type": "InstanceDeleted"}`),
					withVersion(milestone.AggregateVersion),
				), milestone.PushedEventMapper),
			},
			reduce: (&milestoneProjection{}).reducePushed,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("milestone"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.milestones3 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"instance-id",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := baseEvent(t)
			got, err := tt.reduce(event)
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}
			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, MilestonesProjectionTable, tt.want)
		})
	}
}
