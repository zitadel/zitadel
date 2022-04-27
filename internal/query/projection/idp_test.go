package projection

import (
	"testing"

	"github.com/lib/pq"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func TestIDPProjection_reduces(t *testing.T) {
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
			name: "instance.reduceIDPAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.IDPConfigAddedEventType),
					instance.AggregateType,
					[]byte(`{
	"idpConfigId": "idp-config-id",
	"name": "custom-zitadel-instance",
	"idpType": 0,
	"stylingType": 0,
	"autoRegister": true
}`),
				), instance.IDPConfigAddedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.idps (id, creation_date, change_date, sequence, resource_owner, instance_id, state, name, styling_type, auto_register, owner_type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
							expectedArgs: []interface{}{
								"idp-config-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPConfigStateActive,
								"custom-zitadel-instance",
								domain.IDPConfigStylingTypeUnspecified,
								true,
								domain.IdentityProviderTypeSystem,
							},
						},
					},
				},
			},
		},
		{
			name: "instance.reduceIDPChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.IDPConfigChangedEventType),
					instance.AggregateType,
					[]byte(`{
	"idpConfigId": "idp-config-id",
	"name": "custom-zitadel-instance",
	"stylingType": 1,
	"autoRegister": true
}`),
				), instance.IDPConfigChangedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps SET (name, styling_type, auto_register, change_date, sequence) = ($1, $2, $3, $4, $5) WHERE (id = $6)",
							expectedArgs: []interface{}{
								"custom-zitadel-instance",
								domain.IDPConfigStylingTypeGoogle,
								true,
								anyArg{},
								uint64(15),
								"idp-config-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.reduceIDPDeactivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.IDPConfigDeactivatedEventType),
					instance.AggregateType,
					[]byte(`{
	"idpConfigId": "idp-config-id"
}`),
				), instance.IDPConfigDeactivatedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceIDPDeactivated,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								domain.IDPConfigStateInactive,
								anyArg{},
								uint64(15),
								"idp-config-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.reduceIDPReactivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.IDPConfigReactivatedEventType),
					instance.AggregateType,
					[]byte(`{
	"idpConfigId": "idp-config-id"
}`),
				), instance.IDPConfigReactivatedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceIDPReactivated,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								domain.IDPConfigStateActive,
								anyArg{},
								uint64(15),
								"idp-config-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.reduceIDPRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.IDPConfigRemovedEventType),
					instance.AggregateType,
					[]byte(`{
	"idpConfigId": "idp-config-id"
}`),
				), instance.IDPConfigRemovedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceIDPRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idps WHERE (id = $1)",
							expectedArgs: []interface{}{
								"idp-config-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.reduceOIDCConfigAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.IDPOIDCConfigAddedEventType),
					instance.AggregateType,
					[]byte(`{
	"idpConfigId": "idp-config-id",
	"clientId": "client-id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"issuer": "issuer",
	"authorizationEndpoint": "https://api.zitadel.ch/authorize",
    "tokenEndpoint": "https://api.zitadel.ch/token",
    "scopes": ["profile"],
    "idpDisplayNameMapping": 0,
    "usernameMapping": 1
}`),
				), instance.IDPOIDCConfigAddedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceOIDCConfigAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps SET (change_date, sequence, type) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.IDPConfigTypeOIDC,
								"idp-config-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idps_oidc_config (idp_id, client_id, client_secret, issuer, scopes, display_name_mapping, username_mapping, authorization_endpoint, token_endpoint) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"client-id",
								anyArg{},
								"issuer",
								pq.StringArray{"profile"},
								domain.OIDCMappingFieldUnspecified,
								domain.OIDCMappingFieldPreferredLoginName,
								"https://api.zitadel.ch/authorize",
								"https://api.zitadel.ch/token",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.reduceOIDCConfigChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.IDPOIDCConfigChangedEventType),
					instance.AggregateType,
					[]byte(`{
	"idpConfigId": "idp-config-id",
	"clientId": "client-id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"issuer": "issuer",
	"authorizationEndpoint": "https://api.zitadel.ch/authorize",
    "tokenEndpoint": "https://api.zitadel.ch/token",
    "scopes": ["profile"],
    "idpDisplayNameMapping": 0,
    "usernameMapping": 1
}`),
				), instance.IDPOIDCConfigChangedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceOIDCConfigChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"idp-config-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idps_oidc_config SET (client_id, client_secret, issuer, authorization_endpoint, token_endpoint, scopes, display_name_mapping, username_mapping) = ($1, $2, $3, $4, $5, $6, $7, $8) WHERE (idp_id = $9)",
							expectedArgs: []interface{}{
								"client-id",
								anyArg{},
								"issuer",
								"https://api.zitadel.ch/authorize",
								"https://api.zitadel.ch/token",
								pq.StringArray{"profile"},
								domain.OIDCMappingFieldUnspecified,
								domain.OIDCMappingFieldPreferredLoginName,
								"idp-config-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.reduceOIDCConfigChanged: no op",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.IDPOIDCConfigChangedEventType),
					instance.AggregateType,
					[]byte("{}"),
				), instance.IDPOIDCConfigChangedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceOIDCConfigChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{},
				},
			},
		},
		{
			name: "instance.reduceJWTConfigAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.IDPJWTConfigAddedEventType),
					instance.AggregateType,
					[]byte(`{
	"idpConfigId": "idp-config-id",
	"jwtEndpoint": "https://api.zitadel.ch/jwt",
	"issuer": "issuer",
    "keysEndpoint": "https://api.zitadel.ch/keys",
    "headerName": "hodor"
}`),
				), instance.IDPJWTConfigAddedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceJWTConfigAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps SET (change_date, sequence, type) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.IDPConfigTypeJWT,
								"idp-config-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idps_jwt_config (idp_id, endpoint, issuer, keys_endpoint, header_name) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"https://api.zitadel.ch/jwt",
								"issuer",
								"https://api.zitadel.ch/keys",
								"hodor",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.reduceJWTConfigChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.IDPJWTConfigChangedEventType),
					instance.AggregateType,
					[]byte(`{
	"idpConfigId": "idp-config-id",
	"jwtEndpoint": "https://api.zitadel.ch/jwt",
	"issuer": "issuer",
    "keysEndpoint": "https://api.zitadel.ch/keys",
    "headerName": "hodor"
}`),
				), instance.IDPJWTConfigChangedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceJWTConfigChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"idp-config-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idps_jwt_config SET (endpoint, issuer, keys_endpoint, header_name) = ($1, $2, $3, $4) WHERE (idp_id = $5)",
							expectedArgs: []interface{}{
								"https://api.zitadel.ch/jwt",
								"issuer",
								"https://api.zitadel.ch/keys",
								"hodor",
								"idp-config-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.reduceJWTConfigChanged: no op",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.IDPJWTConfigChangedEventType),
					instance.AggregateType,
					[]byte(`{}`),
				), instance.IDPJWTConfigChangedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceJWTConfigChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{},
				},
			},
		},
		{
			name: "org.reduceIDPAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.IDPConfigAddedEventType),
					org.AggregateType,
					[]byte(`{
        "idpConfigId": "idp-config-id",
        "name": "custom-zitadel-instance",
        "idpType": 0,
        "stylingType": 0,
        "autoRegister": true
        }`),
				), org.IDPConfigAddedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.idps (id, creation_date, change_date, sequence, resource_owner, instance_id, state, name, styling_type, auto_register, owner_type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
							expectedArgs: []interface{}{
								"idp-config-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPConfigStateActive,
								"custom-zitadel-instance",
								domain.IDPConfigStylingTypeUnspecified,
								true,
								domain.IdentityProviderTypeOrg,
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceIDPChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.IDPConfigChangedEventType),
					org.AggregateType,
					[]byte(`{
        "idpConfigId": "idp-config-id",
        "name": "custom-zitadel-instance",
        "stylingType": 1,
        "autoRegister": true
        }`),
				), org.IDPConfigChangedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps SET (name, styling_type, auto_register, change_date, sequence) = ($1, $2, $3, $4, $5) WHERE (id = $6)",
							expectedArgs: []interface{}{
								"custom-zitadel-instance",
								domain.IDPConfigStylingTypeGoogle,
								true,
								anyArg{},
								uint64(15),
								"idp-config-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceIDPDeactivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.IDPConfigDeactivatedEventType),
					org.AggregateType,
					[]byte(`{
        "idpConfigId": "idp-config-id"
        }`),
				), org.IDPConfigDeactivatedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceIDPDeactivated,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								domain.IDPConfigStateInactive,
								anyArg{},
								uint64(15),
								"idp-config-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceIDPReactivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.IDPConfigReactivatedEventType),
					org.AggregateType,
					[]byte(`{
        "idpConfigId": "idp-config-id"
        }`),
				), org.IDPConfigReactivatedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceIDPReactivated,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								domain.IDPConfigStateActive,
								anyArg{},
								uint64(15),
								"idp-config-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceIDPRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.IDPConfigRemovedEventType),
					org.AggregateType,
					[]byte(`{
        "idpConfigId": "idp-config-id"
        }`),
				), org.IDPConfigRemovedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceIDPRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idps WHERE (id = $1)",
							expectedArgs: []interface{}{
								"idp-config-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceOIDCConfigAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.IDPOIDCConfigAddedEventType),
					org.AggregateType,
					[]byte(`{
        "idpConfigId": "idp-config-id",
        "clientId": "client-id",
        "clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
        },
        "issuer": "issuer",
        "authorizationEndpoint": "https://api.zitadel.ch/authorize",
        "tokenEndpoint": "https://api.zitadel.ch/token",
        "scopes": ["profile"],
        "idpDisplayNameMapping": 0,
        "usernameMapping": 1
        }`),
				), org.IDPOIDCConfigAddedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceOIDCConfigAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps SET (change_date, sequence, type) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.IDPConfigTypeOIDC,
								"idp-config-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idps_oidc_config (idp_id, client_id, client_secret, issuer, scopes, display_name_mapping, username_mapping, authorization_endpoint, token_endpoint) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"client-id",
								anyArg{},
								"issuer",
								pq.StringArray{"profile"},
								domain.OIDCMappingFieldUnspecified,
								domain.OIDCMappingFieldPreferredLoginName,
								"https://api.zitadel.ch/authorize",
								"https://api.zitadel.ch/token",
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceOIDCConfigChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.IDPOIDCConfigChangedEventType),
					org.AggregateType,
					[]byte(`{
        "idpConfigId": "idp-config-id",
        "clientId": "client-id",
        "clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
        },
        "issuer": "issuer",
        "authorizationEndpoint": "https://api.zitadel.ch/authorize",
        "tokenEndpoint": "https://api.zitadel.ch/token",
        "scopes": ["profile"],
        "idpDisplayNameMapping": 0,
        "usernameMapping": 1
        }`),
				), org.IDPOIDCConfigChangedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceOIDCConfigChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"idp-config-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idps_oidc_config SET (client_id, client_secret, issuer, authorization_endpoint, token_endpoint, scopes, display_name_mapping, username_mapping) = ($1, $2, $3, $4, $5, $6, $7, $8) WHERE (idp_id = $9)",
							expectedArgs: []interface{}{
								"client-id",
								anyArg{},
								"issuer",
								"https://api.zitadel.ch/authorize",
								"https://api.zitadel.ch/token",
								pq.StringArray{"profile"},
								domain.OIDCMappingFieldUnspecified,
								domain.OIDCMappingFieldPreferredLoginName,
								"idp-config-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceOIDCConfigChanged: no op",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.IDPOIDCConfigChangedEventType),
					org.AggregateType,
					[]byte("{}"),
				), org.IDPOIDCConfigChangedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceOIDCConfigChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{},
				},
			},
		},
		{
			name: "org.reduceJWTConfigAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.IDPJWTConfigAddedEventType),
					org.AggregateType,
					[]byte(`{
        "idpConfigId": "idp-config-id",
        "jwtEndpoint": "https://api.zitadel.ch/jwt",
        "issuer": "issuer",
        "keysEndpoint": "https://api.zitadel.ch/keys",
        "headerName": "hodor"
        }`),
				), org.IDPJWTConfigAddedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceJWTConfigAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps SET (change_date, sequence, type) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.IDPConfigTypeJWT,
								"idp-config-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idps_jwt_config (idp_id, endpoint, issuer, keys_endpoint, header_name) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"https://api.zitadel.ch/jwt",
								"issuer",
								"https://api.zitadel.ch/keys",
								"hodor",
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceJWTConfigChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.IDPJWTConfigChangedEventType),
					org.AggregateType,
					[]byte(`{
        "idpConfigId": "idp-config-id",
        "jwtEndpoint": "https://api.zitadel.ch/jwt",
        "issuer": "issuer",
        "keysEndpoint": "https://api.zitadel.ch/keys",
        "headerName": "hodor"
        }`),
				), org.IDPJWTConfigChangedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceJWTConfigChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idps SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"idp-config-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idps_jwt_config SET (endpoint, issuer, keys_endpoint, header_name) = ($1, $2, $3, $4) WHERE (idp_id = $5)",
							expectedArgs: []interface{}{
								"https://api.zitadel.ch/jwt",
								"issuer",
								"https://api.zitadel.ch/keys",
								"hodor",
								"idp-config-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceJWTConfigChanged: no op",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.IDPJWTConfigChangedEventType),
					org.AggregateType,
					[]byte(`{}`),
				), org.IDPJWTConfigChangedEventMapper),
			},
			reduce: (&IDPProjection{}).reduceJWTConfigChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       IDPTable,
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
