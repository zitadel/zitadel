package projection

import (
	"context"
	"testing"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
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
			reduce: (&instanceMemberProjection{
				StatementHandler: getStatementHandlerWithFilters(
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
				)(t)}).reduceAdded,
			want: wantReduce{
				aggregateType:    instance.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       InstanceMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.instance_members2 (user_id, user_resource_owner, owner_removed_user, roles, creation_date, change_date, sequence, resource_owner, instance_id, owner_removed, id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
							expectedArgs: []interface{}{
								"user-id",
								"org1",
								false,
								database.StringArray{"role"},
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								false,
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
			name: "org.OrgRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgRemovedEventType),
					org.AggregateType,
					[]byte(`{}`),
				), org.OrgRemovedEventMapper),
			},
			reduce: (&instanceMemberProjection{}).reduceUserOwnerRemoved,
			want: wantReduce{
				aggregateType:    org.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       InstanceMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.instance_members2 SET (change_date, sequence, owner_removed_user) = ($1, $2, $3) WHERE (user_resource_owner = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
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
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, tt.want)
		})
	}
}
