package projection

import (
	"context"
	"testing"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
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
			name: "instance MemberAddedType",
			args: args{
				event: getEvent(
					testEvent(
						instance.MemberAddedEventType,
						instance.AggregateType,
						[]byte(`{
					"userId": "user-id",
					"roles": ["role"]
				}`),
					), instance.MemberAddedEventMapper),
			},
			reduce: (&instanceMemberProjection{
				es: newMockEventStore().appendFilterResponse(
					[]eventstore.Event{
						user.NewHumanAddedEvent(context.Background(),
							&user.NewAggregate("user-id", "org1").Aggregate,
							"username1",
							"firstname1",
							"lastname1",
							"nickname1",
							"displayname1",
							language.German,
							domain.GenderMale,
							"email1",
							true,
						),
					},
				),
			}).reduceAdded,
			want: wantReduce{
				aggregateType: instance.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.instance_members4 (user_id, user_resource_owner, roles, creation_date, change_date, sequence, resource_owner, instance_id, id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"user-id",
								"org1",
								database.TextArray[string]{"role"},
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
			name: "instance MemberChangedType",
			args: args{
				event: getEvent(
					testEvent(
						instance.MemberChangedEventType,
						instance.AggregateType,
						[]byte(`{
					"userId": "user-id",
					"roles": ["role", "changed"]
				}`),
					), instance.MemberChangedEventMapper),
			},
			reduce: (&instanceMemberProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType: instance.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.instance_members4 SET (roles, change_date, sequence) = ($1, $2, $3) WHERE (instance_id = $4) AND (user_id = $5)",
							expectedArgs: []interface{}{
								database.TextArray[string]{"role", "changed"},
								anyArg{},
								uint64(15),
								"instance-id",
								"user-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance MemberCascadeRemovedType",
			args: args{
				event: getEvent(
					testEvent(
						instance.MemberCascadeRemovedEventType,
						instance.AggregateType,
						[]byte(`{
					"userId": "user-id"
				}`),
					), instance.MemberCascadeRemovedEventMapper),
			},
			reduce: (&instanceMemberProjection{}).reduceCascadeRemoved,
			want: wantReduce{
				aggregateType: instance.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_members4 WHERE (instance_id = $1) AND (user_id = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"user-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance MemberRemovedType",
			args: args{
				event: getEvent(
					testEvent(
						instance.MemberRemovedEventType,
						instance.AggregateType,
						[]byte(`{
					"userId": "user-id"
				}`),
					), instance.MemberRemovedEventMapper),
			},
			reduce: (&instanceMemberProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: instance.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_members4 WHERE (instance_id = $1) AND (user_id = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"user-id",
							},
						},
					},
				},
			},
		},
		{
			name: "user UserRemoved",
			args: args{
				event: getEvent(
					testEvent(
						user.UserRemovedType,
						user.AggregateType,
						[]byte(`{}`),
					), user.UserRemovedEventMapper),
			},
			reduce: (&instanceMemberProjection{}).reduceUserRemoved,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_members4 WHERE (instance_id = $1) AND (user_id = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.OrgRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			reduce: (&instanceMemberProjection{}).reduceUserOwnerRemoved,
			want: wantReduce{
				aggregateType: org.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_members4 WHERE (instance_id = $1) AND (user_resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(MemberInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_members4 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
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
			if ok := zerrors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, InstanceMemberProjectionTable, tt.want)
		})
	}
}
