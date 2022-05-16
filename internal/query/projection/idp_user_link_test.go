package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
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
				event: getEvent(testEvent(
					repository.EventType(user.UserIDPLinkAddedType),
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
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IDPUserLinkTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.idp_user_links (idp_id, user_id, external_user_id, creation_date, change_date, sequence, resource_owner, display_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"agg-id",
								"external-user-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
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
				event: getEvent(testEvent(
					repository.EventType(user.UserIDPLinkRemovedType),
					user.AggregateType,
					[]byte(`{
	"idpConfigId": "idp-config-id",
    "userId": "external-user-id"
}`),
				), user.UserIDPLinkRemovedEventMapper),
			},
			reduce: (&idpUserLinkProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IDPUserLinkTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.idp_user_links WHERE (idp_id = $1) AND (user_id = $2) AND (external_user_id = $3)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"agg-id",
								"external-user-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceCascadeRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserIDPLinkCascadeRemovedType),
					user.AggregateType,
					[]byte(`{
	"idpConfigId": "idp-config-id",
    "userId": "external-user-id"
}`),
				), user.UserIDPLinkCascadeRemovedEventMapper),
			},
			reduce: (&idpUserLinkProjection{}).reduceCascadeRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IDPUserLinkTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.idp_user_links WHERE (idp_id = $1) AND (user_id = $2) AND (external_user_id = $3)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"agg-id",
								"external-user-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceOrgRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgRemovedEventType),
					org.AggregateType,
					[]byte(`{}`),
				), org.OrgRemovedEventMapper),
			},
			reduce: (&idpUserLinkProjection{}).reduceOrgRemoved,
			want: wantReduce{
				aggregateType:    org.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IDPUserLinkTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.idp_user_links WHERE (resource_owner = $1)",
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
				event: getEvent(testEvent(
					repository.EventType(user.UserRemovedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.UserRemovedEventMapper),
			},
			reduce: (&idpUserLinkProjection{}).reduceUserRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IDPUserLinkTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.idp_user_links WHERE (user_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.IDPConfigRemovedEvent",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.IDPConfigRemovedEventType),
					org.AggregateType,
					[]byte(`{
						"idpConfigId": "idp-config-id"
					}`),
				), org.IDPConfigRemovedEventMapper),
			},
			reduce: (&idpUserLinkProjection{}).reduceIDPConfigRemoved,
			want: wantReduce{
				aggregateType:    org.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IDPUserLinkTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.idp_user_links WHERE (idp_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name: "iam.IDPConfigRemovedEvent",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.IDPConfigRemovedEventType),
					iam.AggregateType,
					[]byte(`{
						"idpConfigId": "idp-config-id"
					}`),
				), iam.IDPConfigRemovedEventMapper),
			},
			reduce: (&idpUserLinkProjection{}).reduceIDPConfigRemoved,
			want: wantReduce{
				aggregateType:    iam.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IDPUserLinkTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.idp_user_links WHERE (idp_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"ro-id",
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
