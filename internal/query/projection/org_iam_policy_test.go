package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func TestOrgIAMPolicyProjection_reduces(t *testing.T) {
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
			name: "org.reduceAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgIAMPolicyAddedEventType),
					org.AggregateType,
					[]byte(`{
						"userLoginMustBeDomain": true
}`),
				), org.OrgIAMPolicyAddedEventMapper),
			},
			reduce: (&orgIAMPolicyProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       OrgIAMPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.org_iam_policies (creation_date, change_date, sequence, id, state, user_login_must_be_domain, is_default, resource_owner) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.PolicyStateActive,
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
			name:   "org.reduceChanged",
			reduce: (&orgIAMPolicyProjection{}).reduceChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgIAMPolicyChangedEventType),
					org.AggregateType,
					[]byte(`{
						"userLoginMustBeDomain": true
		}`),
				), org.OrgIAMPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       OrgIAMPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.org_iam_policies SET (change_date, sequence, user_login_must_be_domain) = ($1, $2, $3) WHERE (id = $4)",
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
		{
			name:   "org.reduceRemoved",
			reduce: (&orgIAMPolicyProjection{}).reduceRemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgIAMPolicyRemovedEventType),
					org.AggregateType,
					nil,
				), org.OrgIAMPolicyRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       OrgIAMPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.org_iam_policies WHERE (id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceAdded",
			reduce: (&orgIAMPolicyProjection{}).reduceAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.OrgIAMPolicyAddedEventType),
					iam.AggregateType,
					[]byte(`{
						"userLoginMustBeDomain": true
					}`),
				), iam.OrgIAMPolicyAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       OrgIAMPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.org_iam_policies (creation_date, change_date, sequence, id, state, user_login_must_be_domain, is_default, resource_owner) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.PolicyStateActive,
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
			name:   "iam.reduceChanged",
			reduce: (&orgIAMPolicyProjection{}).reduceChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.OrgIAMPolicyChangedEventType),
					iam.AggregateType,
					[]byte(`{
						"userLoginMustBeDomain": true
					}`),
				), iam.OrgIAMPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       OrgIAMPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.org_iam_policies SET (change_date, sequence, user_login_must_be_domain) = ($1, $2, $3) WHERE (id = $4)",
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
