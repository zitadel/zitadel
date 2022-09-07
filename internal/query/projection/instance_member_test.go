package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestInstanceMemberProjection_reduces(t *testing.T) {
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
			name: "instance.MemberAddedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.MemberAddedEventType),
					instance.AggregateType,
					[]byte(`{
					"userId": "user-id",
					"roles": ["role"]
				}`),
				), instance.MemberAddedEventMapper),
			},
			reduce: (&instanceMemberProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    instance.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       InstanceMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.instance_members2 (user_id, roles, creation_date, change_date, sequence, resource_owner, instance_id, id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"user-id",
								database.StringArray{"role"},
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.MemberChangedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.MemberChangedEventType),
					instance.AggregateType,
					[]byte(`{
					"userId": "user-id",
					"roles": ["role", "changed"]
				}`),
				), instance.MemberChangedEventMapper),
			},
			reduce: (&instanceMemberProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType:    instance.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       InstanceMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.instance_members2 SET (roles, change_date, sequence) = ($1, $2, $3) WHERE (user_id = $4)",
							expectedArgs: []interface{}{
								database.StringArray{"role", "changed"},
								anyArg{},
								uint64(15),
								"user-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.MemberCascadeRemovedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.MemberCascadeRemovedEventType),
					instance.AggregateType,
					[]byte(`{
					"userId": "user-id"
				}`),
				), instance.MemberCascadeRemovedEventMapper),
			},
			reduce: (&instanceMemberProjection{}).reduceCascadeRemoved,
			want: wantReduce{
				aggregateType:    instance.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       InstanceMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_members2 WHERE (user_id = $1)",
							expectedArgs: []interface{}{
								"user-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.MemberRemovedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.MemberRemovedEventType),
					instance.AggregateType,
					[]byte(`{
					"userId": "user-id"
				}`),
				), instance.MemberRemovedEventMapper),
			},
			reduce: (&instanceMemberProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType:    instance.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       InstanceMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_members2 WHERE (user_id = $1)",
							expectedArgs: []interface{}{
								"user-id",
							},
						},
					},
				},
			},
		},
		{
			name: "user.UserRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserRemovedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.UserRemovedEventMapper),
			},
			reduce: (&instanceMemberProjection{}).reduceUserRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       InstanceMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_members2 WHERE (user_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.reduceInstanceRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.InstanceRemovedEventType),
					instance.AggregateType,
					[]byte(`{"name": "Name"}`),
				), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(MemberInstanceID),
			want: wantReduce{
				projection:       InstanceMemberProjectionTable,
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_members2 WHERE (instance_id = $1)",
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
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, tt.want)
		})
	}
}
