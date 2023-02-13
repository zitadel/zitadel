package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func TestIDPTemplateProjection_reducesOIDC(t *testing.T) {
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
			name: "instance reduceOIDCIDPAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.OIDCIDPAddedEventType),
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"issuer": "issuer",
	"client_id": "client_id",
	"client_secret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true
}`),
				), instance.OIDCIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOIDCIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.idps3 (id, creation_date, change_date, sequence, resource_owner, instance_id, state, name, owner_type, is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPConfigStateActive,
								"custom-zitadel-instance",
								domain.IdentityProviderTypeSystem,
								true,
								true,
								true,
								true,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idps3_oidc (idp_id, instance_id, issuer, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5, $6)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"issuer",
								"client_id",
								anyArg{},
								database.StringArray{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOIDCIDPChanged minimal",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.OIDCIDPChangedEventType),
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"issuer": "issuer"
}`),
				), instance.OIDCIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOIDCIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps3 SET (name, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								"custom-zitadel-instance",
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idps3_oidc SET issuer = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"issuer",
								"idp-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOIDCIDPChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.OIDCIDPChangedEventType),
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"issuer": "issuer",
	"client_id": "client_id",
	"client_secret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true
}`),
				), instance.OIDCIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOIDCIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps3 SET (name, is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update, change_date, sequence) = ($1, $2, $3, $4, $5, $6, $7) WHERE (id = $8) AND (instance_id = $9)",
							expectedArgs: []interface{}{
								"custom-zitadel-instance",
								true,
								true,
								true,
								true,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idps3_oidc SET (client_id, client_secret, issuer, scopes) = ($1, $2, $3, $4) WHERE (idp_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								"client_id",
								anyArg{},
								"issuer",
								database.StringArray{"profile"},
								"idp-id",
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
			assertReduce(t, got, err, IDPTable, tt.want)
		})
	}
}

func TestIDPTemplateProjection_reducesJWT(t *testing.T) {
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
			name: "instance reduceJWTIDPAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.JWTIDPAddedEventType),
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"issuer": "issuer",
	"jwtEndpoint": "jwt",
	"keysEndpoint": "keys",
	"headerName": "header",
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true
}`),
				), instance.JWTIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceJWTIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.idps3 (id, creation_date, change_date, sequence, resource_owner, instance_id, state, name, owner_type, is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPConfigStateActive,
								"custom-zitadel-instance",
								domain.IdentityProviderTypeSystem,
								true,
								true,
								true,
								true,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idps3_jwt (idp_id, instance_id, issuer, jwt_endpoint, keys_endpoint, header_name) VALUES ($1, $2, $3, $4, $5, $6)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"issuer",
								"jwt",
								"keys",
								"header",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceJWTIDPChanged minimal",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.JWTIDPChangedEventType),
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"issuer": "issuer"
}`),
				), instance.JWTIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceJWTIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps3 SET (name, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								"custom-zitadel-instance",
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idps3_jwt SET issuer = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"issuer",
								"idp-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceJWTIDPChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.JWTIDPChangedEventType),
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"issuer": "issuer",
	"jwtEndpoint": "jwt",
	"keysEndpoint": "keys",
	"headerName": "header",
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true
}`),
				), instance.JWTIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceJWTIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps3 SET (name, is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update, change_date, sequence) = ($1, $2, $3, $4, $5, $6, $7) WHERE (id = $8) AND (instance_id = $9)",
							expectedArgs: []interface{}{
								"custom-zitadel-instance",
								true,
								true,
								true,
								true,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idps3_jwt SET (jwt_endpoint, keys_endpoint, header_name, issuer) = ($1, $2, $3, $4) WHERE (idp_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								"jwt",
								"keys",
								"header",
								"issuer",
								"idp-id",
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
			assertReduce(t, got, err, IDPTable, tt.want)
		})
	}
}

func TestIDPTemplateProjection_reducesGoogle(t *testing.T) {
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
			name: "instance reduceGoogleIDPAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.GoogleIDPAddedEventType),
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"clientID": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true
}`),
				), instance.GoogleIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGoogleIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.idps3 (id, creation_date, change_date, sequence, resource_owner, instance_id, state, owner_type, is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPConfigStateActive,
								domain.IdentityProviderTypeSystem,
								true,
								true,
								true,
								true,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idps3_google (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.StringArray{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceGoogleIDPChanged minimal",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.GoogleIDPChangedEventType),
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"isCreationAllowed": true,
	"clientID": "id"
}`),
				), instance.GoogleIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGoogleIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps3 SET (is_creation_allowed, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								true,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idps3_google SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"id",
								"idp-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceGoogleIDPChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.GoogleIDPChangedEventType),
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"clientID": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true
}`),
				), instance.GoogleIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGoogleIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps3 SET (is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update, change_date, sequence) = ($1, $2, $3, $4, $5, $6) WHERE (id = $7) AND (instance_id = $8)",
							expectedArgs: []interface{}{
								true,
								true,
								true,
								true,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idps3_google SET (client_id, client_secret, scopes) = ($1, $2, $3) WHERE (idp_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								"client_id",
								anyArg{},
								database.StringArray{"profile"},
								"idp-id",
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
			assertReduce(t, got, err, IDPTable, tt.want)
		})
	}
}
