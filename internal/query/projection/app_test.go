package projection

import (
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestAppProjection_reduces(t *testing.T) {
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
			name: "project reduceAppAdded",
			args: args{
				event: getEvent(
					testEvent(
						project.ApplicationAddedType,
						project.AggregateType,
						[]byte(`{
			"appId": "app-id",
			"name": "my-app"
		}`),
					),
					project.ApplicationAddedEventMapper,
				),
			},
			reduce: (&appProjection{}).reduceAppAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.apps6 (id, name, project_id, creation_date, change_date, resource_owner, instance_id, state, sequence) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"app-id",
								"my-app",
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								domain.AppStateActive,
								uint64(15),
							},
						},
					},
				},
			},
		},
		{
			name: "project reduceAppChanged",
			args: args{
				event: getEvent(
					testEvent(
						project.ApplicationChangedType,
						project.AggregateType,
						[]byte(`{
			"appId": "app-id",
			"name": "my-app"
		}`),
					), project.ApplicationChangedEventMapper),
			},
			reduce: (&appProjection{}).reduceAppChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.apps6 SET (name, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								"my-app",
								anyArg{},
								uint64(15),
								"app-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project reduceAppChanged no change",
			args: args{
				event: getEvent(
					testEvent(
						project.ApplicationChangedType,
						project.AggregateType,
						[]byte(`{
			"appId": "app-id"
		}`),
					), project.ApplicationChangedEventMapper),
			},
			reduce: (&appProjection{}).reduceAppChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{},
				},
			},
		},
		{
			name: "project reduceAppDeactivated",
			args: args{
				event: getEvent(
					testEvent(
						project.ApplicationDeactivatedType,
						project.AggregateType,
						[]byte(`{
			"appId": "app-id"
		}`),
					), project.ApplicationDeactivatedEventMapper),
			},
			reduce: (&appProjection{}).reduceAppDeactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.apps6 SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								domain.AppStateInactive,
								anyArg{},
								uint64(15),
								"app-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project reduceAppReactivated",
			args: args{
				event: getEvent(
					testEvent(
						project.ApplicationReactivatedType,
						project.AggregateType,
						[]byte(`{
			"appId": "app-id"
		}`),
					), project.ApplicationReactivatedEventMapper),
			},
			reduce: (&appProjection{}).reduceAppReactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.apps6 SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								domain.AppStateActive,
								anyArg{},
								uint64(15),
								"app-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project reduceAppRemoved",
			args: args{
				event: getEvent(
					testEvent(
						project.ApplicationRemovedType,
						project.AggregateType,
						[]byte(`{
			"appId": "app-id"
		}`),
					), project.ApplicationRemovedEventMapper),
			},
			reduce: (&appProjection{}).reduceAppRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.apps6 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"app-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project reduceProjectRemoved",
			args: args{
				event: getEvent(
					testEvent(
						project.ProjectRemovedType,
						project.AggregateType,
						[]byte(`{}`),
					), project.ProjectRemovedEventMapper),
			},
			reduce: (&appProjection{}).reduceProjectRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.apps6 WHERE (project_id = $1) AND (instance_id = $2)",
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
			reduce: reduceInstanceRemovedHelper(AppColumnInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.apps6 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project reduceAPIConfigAdded",
			args: args{
				event: getEvent(
					testEvent(
						project.APIConfigAddedType,
						project.AggregateType,
						[]byte(`{
		            "appId": "app-id",
					"clientId": "client-id",
					"clientSecret": {},
				    "authMethodType": 1
				}`),
					), project.APIConfigAddedEventMapper),
			},
			reduce: (&appProjection{}).reduceAPIConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.apps6_api_configs (app_id, instance_id, client_id, client_secret, auth_method) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"app-id",
								"instance-id",
								"client-id",
								anyArg{},
								domain.APIAuthMethodTypePrivateKeyJWT,
							},
						},
						{
							expectedStmt: "UPDATE projections.apps6 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"app-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project reduceAPIConfigChanged",
			args: args{
				event: getEvent(
					testEvent(
						project.APIConfigChangedType,
						project.AggregateType,
						[]byte(`{
		            "appId": "app-id",
					"clientId": "client-id",
					"clientSecret": {},
				    "authMethodType": 1
				}`),
					), project.APIConfigChangedEventMapper),
			},
			reduce: (&appProjection{}).reduceAPIConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.apps6_api_configs SET (client_secret, auth_method) = ($1, $2) WHERE (app_id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								domain.APIAuthMethodTypePrivateKeyJWT,
								"app-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.apps6 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"app-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project reduceAPIConfigChanged noop",
			args: args{
				event: getEvent(
					testEvent(
						project.APIConfigChangedType,
						project.AggregateType,
						[]byte(`{
		            "appId": "app-id"
				}`),
					), project.APIConfigChangedEventMapper),
			},
			reduce: (&appProjection{}).reduceAPIConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{},
				},
			},
		},
		{
			name: "project reduceAPIConfigSecretChanged",
			args: args{
				event: getEvent(
					testEvent(
						project.APIConfigSecretChangedType,
						project.AggregateType,
						[]byte(`{
                        "appId": "app-id",
                        "client_secret": {}
		}`),
					), project.APIConfigSecretChangedEventMapper),
			},
			reduce: (&appProjection{}).reduceAPIConfigSecretChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.apps6_api_configs SET client_secret = $1 WHERE (app_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								"app-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.apps6 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"app-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project reduceOIDCConfigAdded",
			args: args{
				event: getEvent(
					testEvent(
						project.OIDCConfigAddedType,
						project.AggregateType,
						[]byte(`{
                        "oidcVersion": 0,
                        "appId": "app-id",
                        "clientId": "client-id",
                        "clientSecret": {},
                        "redirectUris": ["redirect.one.ch", "redirect.two.ch"],
                        "responseTypes": [1,2],
                        "grantTypes": [1,2],
                        "applicationType": 2,
                        "authMethodType": 2,
                        "postLogoutRedirectUris": ["logout.one.ch", "logout.two.ch"],
                        "devMode": true,
                        "accessTokenType": 1,
                        "accessTokenRoleAssertion": true,
                        "idTokenRoleAssertion": true,
                        "idTokenUserinfoAssertion": true,
                        "clockSkew": 1000,
                        "additionalOrigins": ["origin.one.ch", "origin.two.ch"],
						"skipNativeAppSuccessPage": true
		}`),
					), project.OIDCConfigAddedEventMapper),
			},
			reduce: (&appProjection{}).reduceOIDCConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.apps6_oidc_configs (app_id, instance_id, version, client_id, client_secret, redirect_uris, response_types, grant_types, application_type, auth_method_type, post_logout_redirect_uris, is_dev_mode, access_token_type, access_token_role_assertion, id_token_role_assertion, id_token_userinfo_assertion, clock_skew, additional_origins, skip_native_app_success_page) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)",
							expectedArgs: []interface{}{
								"app-id",
								"instance-id",
								domain.OIDCVersionV1,
								"client-id",
								anyArg{},
								database.TextArray[string]{"redirect.one.ch", "redirect.two.ch"},
								database.Array[domain.OIDCResponseType]{1, 2},
								database.Array[domain.OIDCGrantType]{1, 2},
								domain.OIDCApplicationTypeNative,
								domain.OIDCAuthMethodTypeNone,
								database.TextArray[string]{"logout.one.ch", "logout.two.ch"},
								true,
								domain.OIDCTokenTypeJWT,
								true,
								true,
								true,
								1 * time.Microsecond,
								database.TextArray[string]{"origin.one.ch", "origin.two.ch"},
								true,
							},
						},
						{
							expectedStmt: "UPDATE projections.apps6 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"app-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project reduceOIDCConfigChanged",
			args: args{
				event: getEvent(
					testEvent(
						project.OIDCConfigChangedType,
						project.AggregateType,
						[]byte(`{
                        "oidcVersion": 0,
                        "appId": "app-id",
                        "redirectUris": ["redirect.one.ch", "redirect.two.ch"],
                        "responseTypes": [1,2],
                        "grantTypes": [1,2],
                        "applicationType": 2,
                        "authMethodType": 2,
                        "postLogoutRedirectUris": ["logout.one.ch", "logout.two.ch"],
                        "devMode": true,
                        "accessTokenType": 1,
                        "accessTokenRoleAssertion": true,
                        "idTokenRoleAssertion": true,
                        "idTokenUserinfoAssertion": true,
                        "clockSkew": 1000,
                        "additionalOrigins": ["origin.one.ch", "origin.two.ch"],
						"skipNativeAppSuccessPage": true

		}`),
					), project.OIDCConfigChangedEventMapper),
			},
			reduce: (&appProjection{}).reduceOIDCConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.apps6_oidc_configs SET (version, redirect_uris, response_types, grant_types, application_type, auth_method_type, post_logout_redirect_uris, is_dev_mode, access_token_type, access_token_role_assertion, id_token_role_assertion, id_token_userinfo_assertion, clock_skew, additional_origins, skip_native_app_success_page) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) WHERE (app_id = $16) AND (instance_id = $17)",
							expectedArgs: []interface{}{
								domain.OIDCVersionV1,
								database.TextArray[string]{"redirect.one.ch", "redirect.two.ch"},
								database.Array[domain.OIDCResponseType]{1, 2},
								database.Array[domain.OIDCGrantType]{1, 2},
								domain.OIDCApplicationTypeNative,
								domain.OIDCAuthMethodTypeNone,
								database.TextArray[string]{"logout.one.ch", "logout.two.ch"},
								true,
								domain.OIDCTokenTypeJWT,
								true,
								true,
								true,
								1 * time.Microsecond,
								database.TextArray[string]{"origin.one.ch", "origin.two.ch"},
								true,
								"app-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.apps6 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"app-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project reduceOIDCConfigChanged noop",
			args: args{
				event: getEvent(
					testEvent(
						project.OIDCConfigChangedType,
						project.AggregateType,
						[]byte(`{
                        "appId": "app-id"
		}`),
					), project.OIDCConfigChangedEventMapper),
			},
			reduce: (&appProjection{}).reduceOIDCConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{},
				},
			},
		},
		{
			name: "project reduceOIDCConfigSecretChanged",
			args: args{
				event: getEvent(
					testEvent(
						project.OIDCConfigSecretChangedType,
						project.AggregateType,
						[]byte(`{
                        "appId": "app-id",
                        "client_secret": {}
		}`),
					), project.OIDCConfigSecretChangedEventMapper),
			},
			reduce: (&appProjection{}).reduceOIDCConfigSecretChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.apps6_oidc_configs SET client_secret = $1 WHERE (app_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								"app-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.apps6 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"app-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project.reduceOwnerRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			reduce: (&appProjection{}).reduceOwnerRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.apps6 WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
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
			assertReduce(t, got, err, AppProjectionTable, tt.want)
		})
	}
}
