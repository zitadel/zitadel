package projection

import (
	"testing"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestLoginNameProjection_reduces(t *testing.T) {
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
			name: "user.HumanAddedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanAddedType),
					user.AggregateType,
					[]byte(`{
					"userName": "human-added"
				}`),
				), user.HumanAddedEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceUserCreated,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.login_names_users (id, user_name, resource_owner) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								"agg-id",
								"human-added",
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name: "user.HumanRegisteredType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanRegisteredType),
					user.AggregateType,
					[]byte(`{
					"userName": "human-registered"
				}`),
				), user.HumanRegisteredEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceUserCreated,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.login_names_users (id, user_name, resource_owner) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								"agg-id",
								"human-registered",
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name: "user.MachineAddedEventType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.MachineAddedEventType),
					user.AggregateType,
					[]byte(`{
					"userName": "machine-added"
				}`),
				), user.MachineAddedEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceUserCreated,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.login_names_users (id, user_name, resource_owner) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								"agg-id",
								"machine-added",
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name: "user.UserRemovedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserRemovedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.UserRemovedEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceUserRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.login_names_users WHERE (id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "user.UserUserNameChangedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserUserNameChangedType),
					user.AggregateType,
					[]byte(`{
					"userName": "changed"
				}`),
				), user.UsernameChangedEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceUserNameChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_names_users SET (user_name) = ($1) WHERE (id = $2)",
							expectedArgs: []interface{}{
								"changed",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "user.UserDomainClaimedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserDomainClaimedType),
					user.AggregateType,
					[]byte(`{
					"userName": "claimed"
				}`),
				), user.DomainClaimedEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceUserDomainClaimed,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_names_users SET (user_name) = ($1) WHERE (id = $2)",
							expectedArgs: []interface{}{
								"claimed",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.OrgIAMPolicyAddedEventType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgIAMPolicyAddedEventType),
					user.AggregateType,
					[]byte(`{
					"userLoginMustBeDomain": true
				}`),
				), org.OrgIAMPolicyAddedEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceOrgIAMPolicyAdded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.login_names_policies (must_be_domain, is_default, resource_owner) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								true,
								false,
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.OrgIAMPolicyChangedEventType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgIAMPolicyChangedEventType),
					user.AggregateType,
					[]byte(`{
					"userLoginMustBeDomain": false
				}`),
				), org.OrgIAMPolicyChangedEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceOrgIAMPolicyChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_names_policies SET (must_be_domain) = ($1) WHERE (resource_owner = $2)",
							expectedArgs: []interface{}{
								false,
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.OrgIAMPolicyChangedEventType no change",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgIAMPolicyChangedEventType),
					user.AggregateType,
					[]byte(`{}`),
				), org.OrgIAMPolicyChangedEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceOrgIAMPolicyChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{},
				},
			},
		},
		{
			name: "org.OrgIAMPolicyRemovedEventType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgIAMPolicyRemovedEventType),
					user.AggregateType,
					[]byte(`{}`),
				), org.OrgIAMPolicyRemovedEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceOrgIAMPolicyRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.login_names_policies WHERE (resource_owner = $1)",
							expectedArgs: []interface{}{
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.OrgDomainVerifiedEventType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgDomainVerifiedEventType),
					user.AggregateType,
					[]byte(`{
						"domain": "verified"
					}`),
				), org.DomainVerifiedEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceDomainVerified,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.login_names_domains (name, resource_owner) VALUES ($1, $2)",
							expectedArgs: []interface{}{
								"verified",
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.OrgDomainRemovedEventType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgDomainRemovedEventType),
					user.AggregateType,
					[]byte(`{
						"domain": "remove"
					}`),
				), org.DomainRemovedEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceDomainRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.login_names_domains WHERE (name = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"remove",
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.OrgDomainPrimarySetEventType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgDomainPrimarySetEventType),
					user.AggregateType,
					[]byte(`{
						"domain": "primary"
					}`),
				), org.DomainPrimarySetEventMapper),
			},
			reduce: (&LoginNameProjection{}).reducePrimaryDomainSet,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_names_domains SET (is_primary) = ($1) WHERE (resource_owner = $2) AND (is_primary = $3)",
							expectedArgs: []interface{}{
								false,
								"ro-id",
								true,
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.login_names_domains SET (is_primary) = ($1) WHERE (name = $2) AND (resource_owner = $3)",
							expectedArgs: []interface{}{
								true,
								"primary",
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name: "iam.OrgIAMPolicyAddedEventType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.OrgIAMPolicyAddedEventType),
					user.AggregateType,
					[]byte(`{
					"userLoginMustBeDomain": true
				}`),
				), instance.OrgIAMPolicyAddedEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceOrgIAMPolicyAdded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.login_names_policies (must_be_domain, is_default, resource_owner) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								true,
								true,
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name: "iam.OrgIAMPolicyChangedEventType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.OrgIAMPolicyChangedEventType),
					user.AggregateType,
					[]byte(`{
					"userLoginMustBeDomain": false
				}`),
				), instance.OrgIAMPolicyChangedEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceOrgIAMPolicyChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_names_policies SET (must_be_domain) = ($1) WHERE (resource_owner = $2)",
							expectedArgs: []interface{}{
								false,
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name: "iam.OrgIAMPolicyChangedEventType no change",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.OrgIAMPolicyChangedEventType),
					user.AggregateType,
					[]byte(`{}`),
				), instance.OrgIAMPolicyChangedEventMapper),
			},
			reduce: (&LoginNameProjection{}).reduceOrgIAMPolicyChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       LoginNameProjectionTable,
				executer: &testExecuter{
					executions: []execution{},
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
