package projection

import (
	"testing"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestIDPUserLinkProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.EventReader
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.EventReader) (*handler.Statement, error)
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
			reduce: (&IDPUserLinkProjection{}).reduceAdded,
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
			reduce: (&IDPUserLinkProjection{}).reduceRemoved,
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
			reduce: (&IDPUserLinkProjection{}).reduceCascadeRemoved,
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
