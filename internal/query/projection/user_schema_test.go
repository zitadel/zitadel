package projection

import (
	"encoding/json"
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/user/schema"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestUserSchemaProjection_reduces(t *testing.T) {
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
			name: "reduceCreated",
			args: args{
				event: getEvent(
					testEvent(
						schema.CreatedType,
						schema.AggregateType,
						[]byte(`{"schemaType": "type", "schema": {"$schema":"urn:zitadel:schema:v1","properties":{"name":{"type":"string","urn:zitadel:schema:permission":{"self":"rw"}}},"type":"object"}, "possibleAuthenticators": [1,2]}`),
					), eventstore.GenericEventMapper[schema.CreatedEvent]),
			},
			reduce: (&userSchemaProjection{}).reduceCreated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("user_schema"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.user_schemas (id, change_date, sequence, instance_id, state, type, revision, schema, possible_authenticators) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								uint64(15),
								"instance-id",
								domain.UserSchemaStateActive,
								"type",
								1,
								json.RawMessage(`{"$schema":"urn:zitadel:schema:v1","properties":{"name":{"type":"string","urn:zitadel:schema:permission":{"self":"rw"}}},"type":"object"}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername, domain.AuthenticatorTypePassword},
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUpdated",
			args: args{
				event: getEvent(
					testEvent(
						schema.CreatedType,
						schema.AggregateType,
						[]byte(`{"schemaType": "type", "schema": {"$schema":"urn:zitadel:schema:v1","properties":{"name":{"type":"string","urn:zitadel:schema:permission":{"self":"rw"}}},"type":"object"}, "possibleAuthenticators": [1,2]}`),
					), eventstore.GenericEventMapper[schema.UpdatedEvent]),
			},
			reduce: (&userSchemaProjection{}).reduceUpdated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("user_schema"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.user_schemas SET (change_date, sequence, type, schema, revision, possible_authenticators) = ($1, $2, $3, $4, revision + $5, $6) WHERE (id = $7) AND (instance_id = $8)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"type",
								json.RawMessage(`{"$schema":"urn:zitadel:schema:v1","properties":{"name":{"type":"string","urn:zitadel:schema:permission":{"self":"rw"}}},"type":"object"}`),
								1,
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername, domain.AuthenticatorTypePassword},
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceDeactivated",
			args: args{
				event: getEvent(
					testEvent(
						schema.DeactivatedType,
						schema.AggregateType,
						nil,
					), eventstore.GenericEventMapper[schema.DeactivatedEvent]),
			},
			reduce: (&userSchemaProjection{}).reduceDeactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("user_schema"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.user_schemas SET (change_date, sequence, state) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.UserSchemaStateInactive,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceReactivated",
			args: args{
				event: getEvent(
					testEvent(
						schema.ReactivatedType,
						schema.AggregateType,
						nil,
					), eventstore.GenericEventMapper[schema.ReactivatedEvent]),
			},
			reduce: (&userSchemaProjection{}).reduceReactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("user_schema"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.user_schemas SET (change_date, sequence, state) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.UserSchemaStateActive,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceDeleted",
			args: args{
				event: getEvent(
					testEvent(
						schema.DeletedType,
						schema.AggregateType,
						nil,
					), eventstore.GenericEventMapper[schema.DeletedEvent]),
			},
			reduce: (&userSchemaProjection{}).reduceDeleted,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("user_schema"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_schemas WHERE (id = $1) AND (instance_id = $2)",
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
			reduce: reduceInstanceRemovedHelper(UserSchemaInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_schemas WHERE (instance_id = $1)",
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
			assertReduce(t, got, err, UserSchemaTable, tt.want)
		})
	}
}
