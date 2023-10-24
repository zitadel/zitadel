package projection

import (
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestMilestonesProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.Event
	}
	now := time.Now()
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
					instance.InstanceAddedEventType,
					instance.AggregateType,
					[]byte(`{}`),
					now,
				), instance.InstanceAddedEventMapper),
			},
			reduce: (&milestoneProjection{}).reduceInstanceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.milestones (instance_id, type, reached_date) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								"instance-id",
								milestone.InstanceCreated,
								now,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.milestones (instance_id, type) VALUES ($1, $2)",
							expectedArgs: []interface{}{
								"instance-id",
								milestone.AuthenticationSucceededOnInstance,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.milestones (instance_id, type) VALUES ($1, $2)",
							expectedArgs: []interface{}{
								"instance-id",
								milestone.ProjectCreated,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.milestones (instance_id, type) VALUES ($1, $2)",
							expectedArgs: []interface{}{
								"instance-id",
								milestone.ApplicationCreated,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.milestones (instance_id, type) VALUES ($1, $2)",
							expectedArgs: []interface{}{
								"instance-id",
								milestone.AuthenticationSucceededOnApplication,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.milestones (instance_id, type) VALUES ($1, $2)",
							expectedArgs: []interface{}{
								"instance-id",
								milestone.InstanceDeleted,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceInstancePrimaryDomainSet",
			args: args{
				event: getEvent(testEvent(
					instance.InstanceDomainPrimarySetEventType,
					instance.AggregateType,
					[]byte(`{"domain": "my.domain"}`),
				), instance.DomainPrimarySetEventMapper),
			},
			reduce: (&milestoneProjection{}).reduceInstanceDomainPrimarySet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.milestones SET primary_domain = $1 WHERE (instance_id = $2) AND (last_pushed_date IS NULL)",
							expectedArgs: []interface{}{
								"my.domain",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectAdded",
			args: args{
				event: getEvent(timedTestEvent(
					project.ProjectAddedType,
					project.AggregateType,
					[]byte(`{}`),
					now,
				), project.ProjectAddedEventMapper),
			},
			reduce: (&milestoneProjection{}).reduceProjectAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.milestones SET reached_date = $1 WHERE (instance_id = $2) AND (type = $3) AND (reached_date IS NULL)",
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
			name: "reduceApplicationAdded",
			args: args{
				event: getEvent(timedTestEvent(
					project.ApplicationAddedType,
					project.AggregateType,
					[]byte(`{}`),
					now,
				), project.ApplicationAddedEventMapper),
			},
			reduce: (&milestoneProjection{}).reduceApplicationAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.milestones SET reached_date = $1 WHERE (instance_id = $2) AND (type = $3) AND (reached_date IS NULL)",
							expectedArgs: []interface{}{
								now,
								"instance-id",
								milestone.ApplicationCreated,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceOIDCConfigAdded user event",
			args: args{
				event: getEvent(testEvent(
					project.OIDCConfigAddedType,
					project.AggregateType,
					[]byte(`{}`),
				), project.OIDCConfigAddedEventMapper),
			},
			reduce: (&milestoneProjection{}).reduceOIDCConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer:      &testExecuter{},
			},
		},
		{
			name: "reduceOIDCConfigAdded system event",
			args: args{
				event: getEvent(toSystemEvent(testEvent(
					project.OIDCConfigAddedType,
					project.AggregateType,
					[]byte(`{"clientId": "client-id"}`),
				)), project.OIDCConfigAddedEventMapper),
			},
			reduce: (&milestoneProjection{}).reduceOIDCConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.milestones SET ignore_client_ids = array_append(ignore_client_ids, $1) WHERE (instance_id = $2) AND (type = $3) AND (reached_date IS NULL)",
							expectedArgs: []interface{}{
								"client-id",
								"instance-id",
								milestone.AuthenticationSucceededOnApplication,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAPIConfigAdded user event",
			args: args{
				event: getEvent(testEvent(
					project.APIConfigAddedType,
					project.AggregateType,
					[]byte(`{}`),
				), project.APIConfigAddedEventMapper),
			},
			reduce: (&milestoneProjection{}).reduceAPIConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer:      &testExecuter{},
			},
		},
		{
			name: "reduceAPIConfigAdded system event",
			args: args{
				event: getEvent(toSystemEvent(testEvent(
					project.APIConfigAddedType,
					project.AggregateType,
					[]byte(`{"clientId": "client-id"}`),
				)), project.APIConfigAddedEventMapper),
			},
			reduce: (&milestoneProjection{}).reduceAPIConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.milestones SET ignore_client_ids = array_append(ignore_client_ids, $1) WHERE (instance_id = $2) AND (type = $3) AND (reached_date IS NULL)",
							expectedArgs: []interface{}{
								"client-id",
								"instance-id",
								milestone.AuthenticationSucceededOnApplication,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserTokenAdded",
			args: args{
				event: getEvent(timedTestEvent(
					user.UserTokenAddedType,
					user.AggregateType,
					[]byte(`{"applicationId": "client-id"}`),
					now,
				), user.UserTokenAddedEventMapper),
			},
			reduce: (&milestoneProjection{}).reduceUserTokenAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("user"),
				sequence:      15,
				executer: &testExecuter{
					// TODO: This can be optimized to only use one statement with OR
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.milestones SET reached_date = $1 WHERE (instance_id = $2) AND (type = $3) AND (reached_date IS NULL)",
							expectedArgs: []interface{}{
								now,
								"instance-id",
								milestone.AuthenticationSucceededOnInstance,
							},
						},
						{
							expectedStmt: "UPDATE projections.milestones SET reached_date = $1 WHERE (instance_id = $2) AND (type = $3) AND (NOT (ignore_client_ids @> $4)) AND (reached_date IS NULL)",
							expectedArgs: []interface{}{
								now,
								"instance-id",
								milestone.AuthenticationSucceededOnApplication,
								database.TextArray[string]{"client-id"},
							},
						},
					},
				},
			},
		},
		{
			name: "reduceInstanceRemoved",
			args: args{
				event: getEvent(timedTestEvent(
					instance.InstanceRemovedEventType,
					instance.AggregateType,
					[]byte(`{}`),
					now,
				), instance.InstanceRemovedEventMapper),
			},
			reduce: (&milestoneProjection{}).reduceInstanceRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.milestones SET reached_date = $1 WHERE (instance_id = $2) AND (type = $3) AND (reached_date IS NULL)",
							expectedArgs: []interface{}{
								now,
								"instance-id",
								milestone.InstanceDeleted,
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
				), milestone.PushedEventMapper),
			},
			reduce: (&milestoneProjection{}).reduceMilestonePushed,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("milestone"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.milestones SET last_pushed_date = $1 WHERE (instance_id = $2) AND (type = $3)",
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
			name: "reduceMilestonePushed instance deleted milestone",
			args: args{
				event: getEvent(testEvent(
					milestone.PushedEventType,
					milestone.AggregateType,
					[]byte(`{"type": "InstanceDeleted"}`),
				), milestone.PushedEventMapper),
			},
			reduce: (&milestoneProjection{}).reduceMilestonePushed,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("milestone"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.milestones WHERE (instance_id = $1)",
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
			if !errors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}
			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, MilestonesProjectionTable, tt.want)
		})
	}
}
