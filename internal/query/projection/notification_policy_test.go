package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func TestNotificationPolicyProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.Event
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.Event) (*handler.Statement, error)
		want   wantReduce
	}{
		{
			name: "org reduceAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.NotificationPolicyAddedEventType,
						org.AggregateType,
						[]byte(`{
						"passwordChange": true
}`),
					), org.NotificationPolicyAddedEventMapper),
			},
			reduce: (&notificationPolicyProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.notification_policies (creation_date, change_date, sequence, id, state, password_change, is_default, resource_owner, instance_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.PolicyStateActive,
								true,
								false,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceChanged",
			reduce: (&notificationPolicyProjection{}).reduceChanged,
			args: args{
				event: getEvent(
					testEvent(
						org.NotificationPolicyChangedEventType,
						org.AggregateType,
						[]byte(`{
						"passwordChange": true
		}`),
					), org.NotificationPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.notification_policies SET (change_date, sequence, password_change) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceRemoved",
			reduce: (&notificationPolicyProjection{}).reduceRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.NotificationPolicyRemovedEventType,
						org.AggregateType,
						nil,
					), org.NotificationPolicyRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.notification_policies WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		}, {
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(NotificationPolicyColumnInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.notification_policies WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceAdded",
			reduce: (&notificationPolicyProjection{}).reduceAdded,
			args: args{
				event: getEvent(
					testEvent(
						instance.NotificationPolicyAddedEventType,
						instance.AggregateType,
						[]byte(`{
						"passwordChange": true
					}`),
					), instance.NotificationPolicyAddedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.notification_policies (creation_date, change_date, sequence, id, state, password_change, is_default, resource_owner, instance_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.PolicyStateActive,
								true,
								true,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceChanged",
			reduce: (&notificationPolicyProjection{}).reduceChanged,
			args: args{
				event: getEvent(
					testEvent(
						instance.NotificationPolicyChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"passwordChange": true
					}`),
					), instance.NotificationPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.notification_policies SET (change_date, sequence, password_change) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceOwnerRemoved",
			reduce: (&notificationPolicyProjection{}).reduceOwnerRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.notification_policies WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
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

			if ok := errors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, NotificationPolicyProjectionTable, tt.want)
		})
	}
}
