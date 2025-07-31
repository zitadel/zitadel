package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestIDPLoginPolicyLinkProjection_reduces(t *testing.T) {
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
			name: "iam reduceAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.LoginPolicyIDPProviderAddedEventType,
						instance.AggregateType,
						[]byte(`{
	"idpConfigId": "idp-config-id",
    "idpProviderType": 1
}`),
					), instance.IdentityProviderAddedEventMapper),
			},
			reduce: (&idpLoginPolicyLinkProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: instance.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.idp_login_policy_links5 (idp_id, aggregate_id, creation_date, change_date, sequence, resource_owner, instance_id, provider_type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []any{
								"idp-config-id",
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IdentityProviderTypeSystem,
							},
						},
					},
				},
			},
		},
		{
			name: "iam reduceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.LoginPolicyIDPProviderRemovedEventType,
						instance.AggregateType,
						[]byte(`{
	"idpConfigId": "idp-config-id",
    "idpProviderType": 1
}`),
					), instance.IdentityProviderRemovedEventMapper),
			},
			reduce: (&idpLoginPolicyLinkProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: instance.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_login_policy_links5 WHERE (idp_id = $1) AND (aggregate_id = $2) AND (instance_id = $3)",
							expectedArgs: []any{
								"idp-config-id",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "iam reduceCascadeRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.LoginPolicyIDPProviderCascadeRemovedEventType,
						instance.AggregateType,
						[]byte(`{
	"idpConfigId": "idp-config-id",
    "idpProviderType": 1
}`),
					), instance.IdentityProviderCascadeRemovedEventMapper),
			},
			reduce: (&idpLoginPolicyLinkProjection{}).reduceCascadeRemoved,
			want: wantReduce{
				aggregateType: instance.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_login_policy_links5 WHERE (idp_id = $1) AND (aggregate_id = $2) AND (instance_id = $3)",
							expectedArgs: []any{
								"idp-config-id",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.LoginPolicyIDPProviderAddedEventType,
						org.AggregateType,
						[]byte(`{
	"idpConfigId": "idp-config-id",
    "idpProviderType": 1
}`),
					), org.IdentityProviderAddedEventMapper),
			},
			reduce: (&idpLoginPolicyLinkProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: org.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.idp_login_policy_links5 (idp_id, aggregate_id, creation_date, change_date, sequence, resource_owner, instance_id, provider_type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []any{
								"idp-config-id",
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IdentityProviderTypeOrg,
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.LoginPolicyIDPProviderRemovedEventType,
						org.AggregateType,
						[]byte(`{
	"idpConfigId": "idp-config-id",
    "idpProviderType": 1
}`),
					), org.IdentityProviderRemovedEventMapper),
			},
			reduce: (&idpLoginPolicyLinkProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: org.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_login_policy_links5 WHERE (idp_id = $1) AND (aggregate_id = $2) AND (instance_id = $3)",
							expectedArgs: []any{
								"idp-config-id",
								"agg-id",
								"instance-id",
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
			reduce: reduceInstanceRemovedHelper(IDPUserLinkInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_login_policy_links5 WHERE (instance_id = $1)",
							expectedArgs: []any{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceCascadeRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.LoginPolicyIDPProviderCascadeRemovedEventType,
						org.AggregateType,
						[]byte(`{
	"idpConfigId": "idp-config-id",
    "idpProviderType": 1
}`),
					), org.IdentityProviderCascadeRemovedEventMapper),
			},
			reduce: (&idpLoginPolicyLinkProjection{}).reduceCascadeRemoved,
			want: wantReduce{
				aggregateType: org.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_login_policy_links5 WHERE (idp_id = $1) AND (aggregate_id = $2) AND (instance_id = $3)",
							expectedArgs: []any{
								"idp-config-id",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reducePolicyRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.LoginPolicyRemovedEventType,
						org.AggregateType,
						nil,
					), org.LoginPolicyRemovedEventMapper),
			},
			reduce: (&idpLoginPolicyLinkProjection{}).reducePolicyRemoved,
			want: wantReduce{
				aggregateType: org.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_login_policy_links5 WHERE (aggregate_id = $1) AND (instance_id = $2)",
							expectedArgs: []any{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org IDPConfigRemovedEvent",
			args: args{
				event: getEvent(
					testEvent(
						org.IDPConfigRemovedEventType,
						org.AggregateType,
						[]byte(`{
						"idpConfigId": "idp-config-id"
					}`),
					), org.IDPConfigRemovedEventMapper),
			},
			reduce: (&idpLoginPolicyLinkProjection{}).reduceIDPConfigRemoved,
			want: wantReduce{
				aggregateType: org.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_login_policy_links5 WHERE (idp_id = $1) AND (resource_owner = $2) AND (instance_id = $3)",
							expectedArgs: []any{
								"idp-config-id",
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "iam IDPConfigRemovedEvent",
			args: args{
				event: getEvent(
					testEvent(
						instance.IDPConfigRemovedEventType,
						instance.AggregateType,
						[]byte(`{
						"idpConfigId": "idp-config-id"
					}`),
					), instance.IDPConfigRemovedEventMapper),
			},
			reduce: (&idpLoginPolicyLinkProjection{}).reduceIDPConfigRemoved,
			want: wantReduce{
				aggregateType: instance.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_login_policy_links5 WHERE (idp_id = $1) AND (instance_id = $2)",
							expectedArgs: []any{
								"idp-config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org IDPRemovedEvent",
			args: args{
				event: getEvent(
					testEvent(
						org.IDPRemovedEventType,
						org.AggregateType,
						[]byte(`{
						"id": "id"
					}`),
					), org.IDPRemovedEventMapper),
			},
			reduce: (&idpLoginPolicyLinkProjection{}).reduceIDPRemoved,
			want: wantReduce{
				aggregateType: org.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_login_policy_links5 WHERE (idp_id = $1) AND (resource_owner = $2) AND (instance_id = $3)",
							expectedArgs: []any{
								"id",
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "iam IDPRemovedEvent",
			args: args{
				event: getEvent(
					testEvent(
						instance.IDPRemovedEventType,
						instance.AggregateType,
						[]byte(`{
						"id": "id"
					}`),
					), instance.IDPRemovedEventMapper),
			},
			reduce: (&idpLoginPolicyLinkProjection{}).reduceIDPRemoved,
			want: wantReduce{
				aggregateType: instance.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_login_policy_links5 WHERE (idp_id = $1) AND (instance_id = $2)",
							expectedArgs: []any{
								"id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceOwnerRemoved",
			reduce: (&idpLoginPolicyLinkProjection{}).reduceOwnerRemoved,
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
							expectedStmt: "DELETE FROM projections.idp_login_policy_links5 WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []any{
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
			if ok := zerrors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPLoginPolicyLinkTable, tt.want)
		})
	}
}
