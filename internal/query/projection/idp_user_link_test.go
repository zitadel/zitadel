package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestIDPUserLinkProjection_reduces(t *testing.T) {
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
			name: "reduceAdded",
			args: args{
				event: getEvent(
					testEvent(
						user.UserIDPLinkAddedType,
						user.AggregateType,
						[]byte(`{
	"idpConfigId": "idp-config-id",
    "userId": "external-user-id",
    "displayName": "gigi@caos.ch" 
}`),
					), user.UserIDPLinkAddedEventMapper),
			},
			reduce: (&idpUserLinkProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.idp_user_links3 (idp_id, user_id, external_user_id, creation_date, change_date, sequence, resource_owner, instance_id, display_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"agg-id",
								"external-user-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								"gigi@caos.ch",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						user.UserIDPLinkRemovedType,
						user.AggregateType,
						[]byte(`{
	"idpConfigId": "idp-config-id",
    "userId": "external-user-id"
}`),
					), user.UserIDPLinkRemovedEventMapper),
			},
			reduce: (&idpUserLinkProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_user_links3 WHERE (idp_id = $1) AND (user_id = $2) AND (external_user_id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"agg-id",
								"external-user-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceCascadeRemoved",
			args: args{
				event: getEvent(
					testEvent(
						user.UserIDPLinkCascadeRemovedType,
						user.AggregateType,
						[]byte(`{
	"idpConfigId": "idp-config-id",
    "userId": "external-user-id"
}`),
					), user.UserIDPLinkCascadeRemovedEventMapper),
			},
			reduce: (&idpUserLinkProjection{}).reduceCascadeRemoved,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_user_links3 WHERE (idp_id = $1) AND (user_id = $2) AND (external_user_id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"agg-id",
								"external-user-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceOwnerRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						[]byte(`{}`),
					), org.OrgRemovedEventMapper),
			},
			reduce: (&idpUserLinkProjection{}).reduceOwnerRemoved,
			want: wantReduce{
				aggregateType: org.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_user_links3 WHERE (resource_owner = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
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
							expectedStmt: "DELETE FROM projections.idp_user_links3 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserRemoved",
			args: args{
				event: getEvent(
					testEvent(
						user.UserRemovedType,
						user.AggregateType,
						[]byte(`{}`),
					), user.UserRemovedEventMapper),
			},
			reduce: (&idpUserLinkProjection{}).reduceUserRemoved,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_user_links3 WHERE (user_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceExternalIDMigrated",
			args: args{
				event: getEvent(testEvent(
					user.UserIDPExternalIDMigratedType,
					user.AggregateType,
					[]byte(`{
	"idpConfigId": "idp-config-id",
    "previousId": "previous-id",
	"newId": "new-id"
}`),
				), eventstore.GenericEventMapper[user.UserIDPExternalIDMigratedEvent]),
			},
			reduce: (&idpUserLinkProjection{}).reduceExternalIDMigrated,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_user_links3 SET (change_date, sequence, external_user_id) = ($1, $2, $3) WHERE (idp_id = $4) AND (user_id = $5) AND (external_user_id = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"new-id",
								"idp-config-id",
								"agg-id",
								"previous-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceExternalUsernameChanged",
			args: args{
				event: getEvent(testEvent(
					user.UserIDPExternalUsernameChangedType,
					user.AggregateType,
					[]byte(`{
	"idpConfigId": "idp-config-id",
    "userId": "external-user-id",
	"username": "new-username"
}`),
				), eventstore.GenericEventMapper[user.UserIDPExternalUsernameEvent]),
			},
			reduce: (&idpUserLinkProjection{}).reduceExternalUsernameChanged,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_user_links3 SET (change_date, sequence, display_name) = ($1, $2, $3) WHERE (idp_id = $4) AND (user_id = $5) AND (external_user_id = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"new-username",
								"idp-config-id",
								"agg-id",
								"external-user-id",
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
			reduce: (&idpUserLinkProjection{}).reduceIDPConfigRemoved,
			want: wantReduce{
				aggregateType: org.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_user_links3 WHERE (idp_id = $1) AND (resource_owner = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
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
			reduce: (&idpUserLinkProjection{}).reduceIDPConfigRemoved,
			want: wantReduce{
				aggregateType: instance.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_user_links3 WHERE (idp_id = $1) AND (resource_owner = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"ro-id",
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
			if ok := zerrors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPUserLinkTable, tt.want)
		})
	}
}
