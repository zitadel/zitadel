package projection

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	idpTemplateInsertStmt = `INSERT INTO projections.idp_templates6` +
		` (id, creation_date, change_date, sequence, resource_owner, instance_id, state, name, owner_type, type, is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update, auto_linking)` +
		` VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	idpTemplateUpdateMinimalStmt = `UPDATE projections.idp_templates6 SET (is_creation_allowed, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)`
	idpTemplateUpdateStmt        = `UPDATE projections.idp_templates6 SET (name, is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update, auto_linking, change_date, sequence)` +
		` = ($1, $2, $3, $4, $5, $6, $7, $8) WHERE (id = $9) AND (instance_id = $10)`
)

func TestIDPTemplateProjection_reducesRemove(t *testing.T) {
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
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(IDPInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_templates6 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceOwnerRemoved",
			reduce: (&idpTemplateProjection{}).reduceOwnerRemoved,
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
							expectedStmt: "DELETE FROM projections.idp_templates6 WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceIDPRemoved",
			reduce: (&idpTemplateProjection{}).reduceIDPRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.IDPRemovedEventType,
						org.AggregateType,
						[]byte(`{
	"id": "idp-id"
}`),
					), org.IDPRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_templates6 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceIDPConfigRemoved",
			reduce: (&idpTemplateProjection{}).reduceIDPConfigRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.IDPConfigRemovedEventType,
						org.AggregateType,
						[]byte(`{
	"idpConfigId": "idp-id"
}`),
					), org.IDPConfigRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_templates6 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
		})
	}
}

func TestIDPTemplateProjection_reducesOAuth(t *testing.T) {
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
			name: "instance reduceOAuthIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.OAuthIDPAddedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"authorizationEndpoint": "auth",
	"tokenEndpoint": "token",
 	"userEndpoint": "user",
	"scopes": ["profile"],
	"idAttribute": "id-attribute",
	"usePKCE": false,
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.OAuthIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOAuthIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"custom-zitadel-instance",
								domain.IdentityProviderTypeSystem,
								domain.IDPTypeOAuth,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_oauth2 (idp_id, instance_id, client_id, client_secret, authorization_endpoint, token_endpoint, user_endpoint, scopes, id_attribute, use_pkce) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								"auth",
								"token",
								"user",
								database.TextArray[string]{"profile"},
								"id-attribute",
								false,
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceOAuthIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.OAuthIDPAddedEventType,
						org.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"authorizationEndpoint": "auth",
	"tokenEndpoint": "token",
 	"userEndpoint": "user",
	"scopes": ["profile"],
	"idAttribute": "id-attribute",
	"usePKCE": true,
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), org.OAuthIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOAuthIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"custom-zitadel-instance",
								domain.IdentityProviderTypeOrg,
								domain.IDPTypeOAuth,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_oauth2 (idp_id, instance_id, client_id, client_secret, authorization_endpoint, token_endpoint, user_endpoint, scopes, id_attribute, use_pkce) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								"auth",
								"token",
								"user",
								database.TextArray[string]{"profile"},
								"id-attribute",
								true,
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOAuthIDPChanged minimal",
			args: args{
				event: getEvent(
					testEvent(
						instance.OAuthIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"isCreationAllowed": true,
	"clientId": "id"
}`),
					), instance.OAuthIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOAuthIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateMinimalStmt,
							expectedArgs: []interface{}{
								true,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_oauth2 SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
			name: "instance reduceOAuthIDPChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.OAuthIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"authorizationEndpoint": "auth",
	"tokenEndpoint": "token",
 	"userEndpoint": "user",
	"scopes": ["profile"],
	"idAttribute": "id-attribute",
	"usePKCE": true,
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.OAuthIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOAuthIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateStmt,
							expectedArgs: []interface{}{
								"custom-zitadel-instance",
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_oauth2 SET (client_id, client_secret, authorization_endpoint, token_endpoint, user_endpoint, scopes, id_attribute, use_pkce) = ($1, $2, $3, $4, $5, $6, $7, $8) WHERE (idp_id = $9) AND (instance_id = $10)",
							expectedArgs: []interface{}{
								"client_id",
								anyArg{},
								"auth",
								"token",
								"user",
								database.TextArray[string]{"profile"},
								"id-attribute",
								true,
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
		})
	}
}

func TestIDPTemplateProjection_reducesAzureAD(t *testing.T) {
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
			name: "instance reduceAzureADIDPAdded minimal",
			args: args{
				event: getEvent(
					testEvent(
						instance.AzureADIDPAddedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
	"client_id": "client_id",
	"client_secret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    }
}`),
					), instance.AzureADIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceAzureADIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"name",
								domain.IdentityProviderTypeSystem,
								domain.IDPTypeAzureAD,
								false,
								false,
								false,
								false,
								domain.AutoLinkingOptionUnspecified,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_azure (idp_id, instance_id, client_id, client_secret, scopes, tenant, is_email_verified) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.TextArray[string](nil),
								"",
								false,
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceAzureADIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.AzureADIDPAddedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
	"client_id": "client_id",
	"client_secret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"tenant": "tenant",
	"isEmailVerified": true,
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.AzureADIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceAzureADIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"name",
								domain.IdentityProviderTypeSystem,
								domain.IDPTypeAzureAD,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_azure (idp_id, instance_id, client_id, client_secret, scopes, tenant, is_email_verified) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
								"tenant",
								true,
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceAzureADIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.AzureADIDPAddedEventType,
						org.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
	"client_id": "client_id",
	"client_secret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"tenant": "tenant",
	"isEmailVerified": true,
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), org.AzureADIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceAzureADIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"name",
								domain.IdentityProviderTypeOrg,
								domain.IDPTypeAzureAD,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_azure (idp_id, instance_id, client_id, client_secret, scopes, tenant, is_email_verified) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
								"tenant",
								true,
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceAzureADIDPChanged minimal",
			args: args{
				event: getEvent(
					testEvent(
						instance.AzureADIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"isCreationAllowed": true,
	"client_id": "id"
}`),
					), instance.AzureADIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceAzureADIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateMinimalStmt,
							expectedArgs: []interface{}{
								true,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_azure SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
			name: "instance reduceAzureADIDPChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.AzureADIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
	"client_id": "client_id",
	"client_secret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"tenant": "tenant",
	"isEmailVerified": true,
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.AzureADIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceAzureADIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateStmt,
							expectedArgs: []interface{}{
								"name",
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_azure SET (client_id, client_secret, scopes, tenant, is_email_verified) = ($1, $2, $3, $4, $5) WHERE (idp_id = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
								"tenant",
								true,
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
		})
	}
}

func TestIDPTemplateProjection_reducesGitHub(t *testing.T) {
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
			name: "instance reduceGitHubIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.GitHubIDPAddedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.GitHubIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitHubIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"name",
								domain.IdentityProviderTypeSystem,
								domain.IDPTypeGitHub,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_github (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceGitHubIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.GitHubIDPAddedEventType,
						org.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), org.GitHubIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitHubIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"name",
								domain.IdentityProviderTypeOrg,
								domain.IDPTypeGitHub,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_github (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceGitHubIDPChanged minimal",
			args: args{
				event: getEvent(
					testEvent(
						instance.GitHubIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"isCreationAllowed": true,
	"clientId": "id"
}`),
					), instance.GitHubIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitHubIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateMinimalStmt,
							expectedArgs: []interface{}{
								true,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_github SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
			name: "instance reduceGitHubIDPChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.GitHubIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.GitHubIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitHubIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateStmt,
							expectedArgs: []interface{}{
								"name",
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_github SET (client_id, client_secret, scopes) = ($1, $2, $3) WHERE (idp_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
		})
	}
}

func TestIDPTemplateProjection_reducesGitHubEnterprise(t *testing.T) {
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
			name: "instance reduceGitHubEnterpriseIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.GitHubEnterpriseIDPAddedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"authorizationEndpoint": "auth",
	"tokenEndpoint": "token",
 	"userEndpoint": "user",
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.GitHubEnterpriseIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitHubEnterpriseIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"name",
								domain.IdentityProviderTypeSystem,
								domain.IDPTypeGitHubEnterprise,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_github_enterprise (idp_id, instance_id, client_id, client_secret, authorization_endpoint, token_endpoint, user_endpoint, scopes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								"auth",
								"token",
								"user",
								database.TextArray[string]{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceGitHubEnterpriseIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.GitHubEnterpriseIDPAddedEventType,
						org.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"authorizationEndpoint": "auth",
	"tokenEndpoint": "token",
 	"userEndpoint": "user",
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), org.GitHubEnterpriseIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitHubEnterpriseIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"name",
								domain.IdentityProviderTypeOrg,
								domain.IDPTypeGitHubEnterprise,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_github_enterprise (idp_id, instance_id, client_id, client_secret, authorization_endpoint, token_endpoint, user_endpoint, scopes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								"auth",
								"token",
								"user",
								database.TextArray[string]{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceGitHubEnterpriseIDPChanged minimal",
			args: args{
				event: getEvent(
					testEvent(
						instance.GitHubEnterpriseIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"isCreationAllowed": true,
	"clientId": "id"
}`),
					), instance.GitHubEnterpriseIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitHubEnterpriseIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateMinimalStmt,
							expectedArgs: []interface{}{
								true,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_github_enterprise SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
			name: "instance reduceGitHubEnterpriseIDPChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.GitHubEnterpriseIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"authorizationEndpoint": "auth",
	"tokenEndpoint": "token",
 	"userEndpoint": "user",
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.GitHubEnterpriseIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitHubEnterpriseIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateStmt,
							expectedArgs: []interface{}{
								"name",
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_github_enterprise SET (client_id, client_secret, authorization_endpoint, token_endpoint, user_endpoint, scopes) = ($1, $2, $3, $4, $5, $6) WHERE (idp_id = $7) AND (instance_id = $8)",
							expectedArgs: []interface{}{
								"client_id",
								anyArg{},
								"auth",
								"token",
								"user",
								database.TextArray[string]{"profile"},
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
		})
	}
}

func TestIDPTemplateProjection_reducesGitLab(t *testing.T) {
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
			name: "instance reduceGitLabIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.GitLabIDPAddedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
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
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.GitLabIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitLabIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"",
								domain.IdentityProviderTypeSystem,
								domain.IDPTypeGitLab,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_gitlab (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceGitLabIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.GitLabIDPAddedEventType,
						org.AggregateType,
						[]byte(`{
	"id": "idp-id",
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
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), org.GitLabIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitLabIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"",
								domain.IdentityProviderTypeOrg,
								domain.IDPTypeGitLab,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_gitlab (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceGitLabIDPChanged minimal",
			args: args{
				event: getEvent(
					testEvent(
						instance.GitLabIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"isCreationAllowed": true,
	"client_id": "id"
}`),
					), instance.GitLabIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitLabIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateMinimalStmt,
							expectedArgs: []interface{}{
								true,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_gitlab SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
			name: "instance reduceGitLabIDPChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.GitLabIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
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
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.GitLabIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitLabIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateStmt,
							expectedArgs: []interface{}{
								"name",
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_gitlab SET (client_id, client_secret, scopes) = ($1, $2, $3) WHERE (idp_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
		})
	}
}

func TestIDPTemplateProjection_reducesGitLabSelfHosted(t *testing.T) {
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
			name: "instance reduceGitLabSelfHostedIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.GitLabSelfHostedIDPAddedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
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
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.GitLabSelfHostedIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitLabSelfHostedIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"name",
								domain.IdentityProviderTypeSystem,
								domain.IDPTypeGitLabSelfHosted,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_gitlab_self_hosted (idp_id, instance_id, issuer, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5, $6)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"issuer",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceGitLabSelfHostedIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.GitLabSelfHostedIDPAddedEventType,
						org.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
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
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), org.GitLabSelfHostedIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitLabSelfHostedIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"name",
								domain.IdentityProviderTypeOrg,
								domain.IDPTypeGitLabSelfHosted,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_gitlab_self_hosted (idp_id, instance_id, issuer, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5, $6)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"issuer",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceGitLabSelfHostedIDPChanged minimal",
			args: args{
				event: getEvent(
					testEvent(
						instance.GitLabSelfHostedIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"isCreationAllowed": true,
	"issuer": "issuer"
}`),
					), instance.GitLabSelfHostedIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitLabSelfHostedIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateMinimalStmt,
							expectedArgs: []interface{}{
								true,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_gitlab_self_hosted SET issuer = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
			name: "instance reduceGitLabSelfHostedIDPChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.GitLabSelfHostedIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
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
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.GitLabSelfHostedIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitLabSelfHostedIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateStmt,
							expectedArgs: []interface{}{
								"name",
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_gitlab_self_hosted SET (issuer, client_id, client_secret, scopes) = ($1, $2, $3, $4) WHERE (idp_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								"issuer",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
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
				event: getEvent(
					testEvent(
						instance.GoogleIDPAddedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.GoogleIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGoogleIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"",
								domain.IdentityProviderTypeSystem,
								domain.IDPTypeGoogle,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_google (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceGoogleIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.GoogleIDPAddedEventType,
						org.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), org.GoogleIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGoogleIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"",
								domain.IdentityProviderTypeOrg,
								domain.IDPTypeGoogle,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_google (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceGoogleIDPChanged minimal",
			args: args{
				event: getEvent(
					testEvent(
						instance.GoogleIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"isCreationAllowed": true,
	"clientId": "id"
}`),
					), instance.GoogleIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGoogleIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateMinimalStmt,
							expectedArgs: []interface{}{
								true,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_google SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
				event: getEvent(
					testEvent(
						instance.GoogleIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.GoogleIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGoogleIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateStmt,
							expectedArgs: []interface{}{
								"name",
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_google SET (client_id, client_secret, scopes) = ($1, $2, $3) WHERE (idp_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
		})
	}
}

func TestIDPTemplateProjection_reducesLDAP(t *testing.T) {
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
			name: "instance reduceLDAPIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.LDAPIDPAddedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"servers": ["server"],
	"startTls": false,
	"baseDN": "basedn",
	"bindDN": "binddn",
	"bindPassword": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"userBase": "user",
	"userObjectClasses": ["object"],
	"userFilters": ["filter"],
	"timeout": 30000000000,
	"rootCA": `+stringToJSONByte("certificate")+`,
	"idAttribute": "id",
	"firstNameAttribute": "first",
	"lastNameAttribute": "last",
	"displayNameAttribute": "display",
	"nickNameAttribute": "nickname",
	"preferredUsernameAttribute": "username",
	"emailAttribute": "email",
	"emailVerifiedAttribute": "email_verified",
	"phoneAttribute": "phone",
	"phoneVerifiedAttribute": "phone_verified",
	"preferredLanguageAttribute": "lang",
	"avatarURLAttribute": "avatar",
	"profileAttribute": "profile",
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.LDAPIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceLDAPIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"custom-zitadel-instance",
								domain.IdentityProviderTypeSystem,
								domain.IDPTypeLDAP,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_ldap2 (idp_id, instance_id, servers, start_tls, base_dn, bind_dn, bind_password, user_base, user_object_classes, user_filters, timeout, root_ca, id_attribute, first_name_attribute, last_name_attribute, display_name_attribute, nick_name_attribute, preferred_username_attribute, email_attribute, email_verified, phone_attribute, phone_verified_attribute, preferred_language_attribute, avatar_url_attribute, profile_attribute) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								database.TextArray[string]{"server"},
								false,
								"basedn",
								"binddn",
								anyArg{},
								"user",
								database.TextArray[string]{"object"},
								database.TextArray[string]{"filter"},
								time.Duration(30000000000),
								[]byte("certificate"),
								"id",
								"first",
								"last",
								"display",
								"nickname",
								"username",
								"email",
								"email_verified",
								"phone",
								"phone_verified",
								"lang",
								"avatar",
								"profile",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceLDAPIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.LDAPIDPAddedEventType,
						org.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"servers": ["server"],
	"startTls": false,
	"baseDN": "basedn",
	"bindDN": "binddn",
	"bindPassword": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"userBase": "user",
	"userObjectClasses": ["object"],
	"userFilters": ["filter"],
	"timeout": 30000000000,
	"rootCA": `+stringToJSONByte("certificate")+`,
	"idAttribute": "id",
	"firstNameAttribute": "first",
	"lastNameAttribute": "last",
	"displayNameAttribute": "display",
	"nickNameAttribute": "nickname",
	"preferredUsernameAttribute": "username",
	"emailAttribute": "email",
	"emailVerifiedAttribute": "email_verified",
	"phoneAttribute": "phone",
	"phoneVerifiedAttribute": "phone_verified",
	"preferredLanguageAttribute": "lang",
	"avatarURLAttribute": "avatar",
	"profileAttribute": "profile",
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), org.LDAPIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceLDAPIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"custom-zitadel-instance",
								domain.IdentityProviderTypeOrg,
								domain.IDPTypeLDAP,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_ldap2 (idp_id, instance_id, servers, start_tls, base_dn, bind_dn, bind_password, user_base, user_object_classes, user_filters, timeout, root_ca, id_attribute, first_name_attribute, last_name_attribute, display_name_attribute, nick_name_attribute, preferred_username_attribute, email_attribute, email_verified, phone_attribute, phone_verified_attribute, preferred_language_attribute, avatar_url_attribute, profile_attribute) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								database.TextArray[string]{"server"},
								false,
								"basedn",
								"binddn",
								anyArg{},
								"user",
								database.TextArray[string]{"object"},
								database.TextArray[string]{"filter"},
								time.Duration(30000000000),
								[]byte("certificate"),
								"id",
								"first",
								"last",
								"display",
								"nickname",
								"username",
								"email",
								"email_verified",
								"phone",
								"phone_verified",
								"lang",
								"avatar",
								"profile",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceLDAPIDPChanged minimal",
			args: args{
				event: getEvent(
					testEvent(
						instance.LDAPIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"baseDN": "basedn"
}`),
					), instance.LDAPIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceLDAPIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (name, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								"custom-zitadel-instance",
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_ldap2 SET base_dn = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"basedn",
								"idp-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceLDAPIDPChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.LDAPIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"servers": ["server"],
	"startTls": false,
	"baseDN": "basedn",
	"bindDN": "binddn",
	"bindPassword": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"userBase": "user",
	"userObjectClasses": ["object"],
	"userFilters": ["filter"],
	"timeout": 30000000000,
	"rootCA": `+stringToJSONByte("certificate")+`,
	"idAttribute": "id",
	"firstNameAttribute": "first",
	"lastNameAttribute": "last",
	"displayNameAttribute": "display",
	"nickNameAttribute": "nickname",
	"preferredUsernameAttribute": "username",
	"emailAttribute": "email",
	"emailVerifiedAttribute": "email_verified",
	"phoneAttribute": "phone",
	"phoneVerifiedAttribute": "phone_verified",
	"preferredLanguageAttribute": "lang",
	"avatarURLAttribute": "avatar",
	"profileAttribute": "profile",
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.LDAPIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceLDAPIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateStmt,
							expectedArgs: []interface{}{
								"custom-zitadel-instance",
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_ldap2 SET (servers, start_tls, base_dn, bind_dn, bind_password, user_base, user_object_classes, user_filters, timeout, root_ca, id_attribute, first_name_attribute, last_name_attribute, display_name_attribute, nick_name_attribute, preferred_username_attribute, email_attribute, email_verified, phone_attribute, phone_verified_attribute, preferred_language_attribute, avatar_url_attribute, profile_attribute) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23) WHERE (idp_id = $24) AND (instance_id = $25)",
							expectedArgs: []interface{}{
								database.TextArray[string]{"server"},
								false,
								"basedn",
								"binddn",
								anyArg{},
								"user",
								database.TextArray[string]{"object"},
								database.TextArray[string]{"filter"},
								time.Duration(30000000000),
								[]byte("certificate"),
								"id",
								"first",
								"last",
								"display",
								"nickname",
								"username",
								"email",
								"email_verified",
								"phone",
								"phone_verified",
								"lang",
								"avatar",
								"profile",
								"idp-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceOwnerRemoved",
			reduce: (&idpProjection{}).reduceOwnerRemoved,
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
							expectedStmt: "DELETE FROM projections.idp_templates6 WHERE (instance_id = $1) AND (resource_owner = $2)",
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
		})
	}
}

func TestIDPTemplateProjection_reducesApple(t *testing.T) {
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
			name: "instance reduceAppleIDPAdded",
			args: args{
				event: getEvent(testEvent(
					instance.AppleIDPAddedEventType,
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"clientId": "client_id",
	"teamId": "team_id",
	"keyId": "key_id",
	"privateKey": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["name"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
				), instance.AppleIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceAppleIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"",
								domain.IdentityProviderTypeSystem,
								domain.IDPTypeApple,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_apple (idp_id, instance_id, client_id, team_id, key_id, private_key, scopes) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								"team_id",
								"key_id",
								anyArg{},
								database.TextArray[string]{"name"},
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceAppleIDPAdded",
			args: args{
				event: getEvent(testEvent(
					org.AppleIDPAddedEventType,
					org.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"clientId": "client_id",
	"teamId": "team_id",
	"keyId": "key_id",
	"privateKey": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["name"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
				), org.AppleIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceAppleIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"",
								domain.IdentityProviderTypeOrg,
								domain.IDPTypeApple,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_apple (idp_id, instance_id, client_id, team_id, key_id, private_key, scopes) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								"team_id",
								"key_id",
								anyArg{},
								database.TextArray[string]{"name"},
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceAppleIDPChanged minimal",
			args: args{
				event: getEvent(testEvent(
					instance.AppleIDPChangedEventType,
					instance.AggregateType,
					[]byte(`{
			"id": "idp-id",
			"isCreationAllowed": true,
			"clientId": "id"
		}`),
				), instance.AppleIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceAppleIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateMinimalStmt,
							expectedArgs: []interface{}{
								true,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_apple SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
			name: "instance reduceAppleIDPChanged",
			args: args{
				event: getEvent(testEvent(
					instance.AppleIDPChangedEventType,
					instance.AggregateType,
					[]byte(`{
			"id": "idp-id",
			"name": "name",
			"clientId": "client_id",
			"teamId": "team_id",
			"keyId": "key_id",
			"privateKey": {
				"cryptoType": 0,
				"algorithm": "RSA-265",
				"keyId": "key-id"
			},
			"scopes": ["name"],
			"isCreationAllowed": true,
			"isLinkingAllowed": true,
			"isAutoCreation": true,
			"isAutoUpdate": true,
			"autoLinkingOption": 1
		}`),
				), instance.AppleIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceAppleIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateStmt,
							expectedArgs: []interface{}{
								"name",
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_apple SET (client_id, team_id, key_id, private_key, scopes) = ($1, $2, $3, $4, $5) WHERE (idp_id = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								"client_id",
								"team_id",
								"key_id",
								anyArg{},
								database.TextArray[string]{"name"},
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
		})
	}
}

func TestIDPTemplateProjection_reducesSAML(t *testing.T) {
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
			name: "instance reduceSAMLIDPAdded",
			args: args{
				event: getEvent(testEvent(
					instance.SAMLIDPAddedEventType,
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"metadata": `+stringToJSONByte("metadata")+`,
	"key": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"certificate": `+stringToJSONByte("certificate")+`,
	"binding": "binding",
	"nameIDFormat": 3,
	"transientMappingAttributeName": "customAttribute",
	"withSignedRequest": true,
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
				), instance.SAMLIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceSAMLIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"custom-zitadel-instance",
								domain.IdentityProviderTypeSystem,
								domain.IDPTypeSAML,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_saml (idp_id, instance_id, metadata, key, certificate, binding, with_signed_request, transient_mapping_attribute_name, name_id_format) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								[]byte("metadata"),
								anyArg{},
								anyArg{},
								"binding",
								true,
								"customAttribute",
								domain.SAMLNameIDFormatTransient,
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceSAMLIDPAdded",
			args: args{
				event: getEvent(testEvent(
					org.SAMLIDPAddedEventType,
					org.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"metadata": `+stringToJSONByte("metadata")+`,
	"key": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"certificate": `+stringToJSONByte("certificate")+`,
	"binding": "binding",
	"nameIDFormat": 3,
	"transientMappingAttributeName": "customAttribute",
	"withSignedRequest": true,
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
				), org.SAMLIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceSAMLIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"custom-zitadel-instance",
								domain.IdentityProviderTypeOrg,
								domain.IDPTypeSAML,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_saml (idp_id, instance_id, metadata, key, certificate, binding, with_signed_request, transient_mapping_attribute_name, name_id_format) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								[]byte("metadata"),
								anyArg{},
								anyArg{},
								"binding",
								true,
								"customAttribute",
								domain.SAMLNameIDFormatTransient,
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceSAMLIDPChanged minimal",
			args: args{
				event: getEvent(testEvent(
					instance.SAMLIDPChangedEventType,
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"binding": "binding"
}`),
				), instance.SAMLIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceSAMLIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (name, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								"custom-zitadel-instance",
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_saml SET binding = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"binding",
								"idp-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceSAMLIDPChanged",
			args: args{
				event: getEvent(testEvent(
					instance.SAMLIDPChangedEventType,
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"metadata": `+stringToJSONByte("metadata")+`,
	"key": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"certificate": `+stringToJSONByte("certificate")+`,
	"binding": "binding",
	"withSignedRequest": true,
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
				), instance.SAMLIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceSAMLIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateStmt,
							expectedArgs: []interface{}{
								"custom-zitadel-instance",
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_saml SET (metadata, key, certificate, binding, with_signed_request) = ($1, $2, $3, $4, $5) WHERE (idp_id = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								[]byte("metadata"),
								anyArg{},
								anyArg{},
								"binding",
								true,
								"idp-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceOwnerRemoved",
			reduce: (&idpProjection{}).reduceOwnerRemoved,
			args: args{
				event: getEvent(testEvent(
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
							expectedStmt: "DELETE FROM projections.idp_templates6 WHERE (instance_id = $1) AND (resource_owner = $2)",
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
		})
	}
}

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
				event: getEvent(
					testEvent(
						instance.OIDCIDPAddedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"issuer": "issuer",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"idTokenMapping": true,
	"usePKCE": true,
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.OIDCIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOIDCIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"",
								domain.IdentityProviderTypeSystem,
								domain.IDPTypeOIDC,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_oidc (idp_id, instance_id, issuer, client_id, client_secret, scopes, id_token_mapping, use_pkce) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"issuer",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
								true,
								true,
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceOIDCIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.OIDCIDPAddedEventType,
						org.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"issuer": "issuer",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"idTokenMapping": true,
	"usePKCE": true,
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), org.OIDCIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOIDCIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"",
								domain.IdentityProviderTypeOrg,
								domain.IDPTypeOIDC,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_oidc (idp_id, instance_id, issuer, client_id, client_secret, scopes, id_token_mapping, use_pkce) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"issuer",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
								true,
								true,
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOIDCIDPChanged minimal",
			args: args{
				event: getEvent(
					testEvent(
						instance.OIDCIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"isCreationAllowed": true,
	"clientId": "id"
}`),
					), instance.OIDCIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOIDCIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateMinimalStmt,
							expectedArgs: []interface{}{
								true,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_oidc SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
			name: "instance reduceOIDCIDPChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.OIDCIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"name": "name",
	"issuer": "issuer",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"idTokenMapping": true,
	"usePKCE": true,
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.OIDCIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOIDCIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateStmt,
							expectedArgs: []interface{}{
								"name",
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_oidc SET (client_id, client_secret, issuer, scopes, id_token_mapping, use_pkce) = ($1, $2, $3, $4, $5, $6) WHERE (idp_id = $7) AND (instance_id = $8)",
							expectedArgs: []interface{}{
								"client_id",
								anyArg{},
								"issuer",
								database.TextArray[string]{"profile"},
								true,
								true,
								"idp-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOIDCIDPMigratedAzureAD",
			args: args{
				event: getEvent(testEvent(
					instance.OIDCIDPMigratedAzureADEventType,
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "name",
	"client_id": "client_id",
	"client_secret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"tenant": "tenant",
	"isEmailVerified": true,
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
				), instance.OIDCIDPMigratedAzureADEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOIDCIDPMigratedAzureAD,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (change_date, sequence, name, type, is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update, auto_linking) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE (id = $10) AND (instance_id = $11)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"name",
								domain.IDPTypeAzureAD,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "DELETE FROM projections.idp_templates6_oidc WHERE (idp_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_azure (idp_id, instance_id, client_id, client_secret, scopes, tenant, is_email_verified) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
								"tenant",
								true,
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceOIDCIDPMigratedAzureAD",
			args: args{
				event: getEvent(testEvent(
					org.OIDCIDPMigratedAzureADEventType,
					org.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "name",
	"client_id": "client_id",
	"client_secret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"tenant": "tenant",
	"isEmailVerified": true,
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
				), org.OIDCIDPMigratedAzureADEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOIDCIDPMigratedAzureAD,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (change_date, sequence, name, type, is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update, auto_linking) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE (id = $10) AND (instance_id = $11)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"name",
								domain.IDPTypeAzureAD,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "DELETE FROM projections.idp_templates6_oidc WHERE (idp_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_azure (idp_id, instance_id, client_id, client_secret, scopes, tenant, is_email_verified) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
								"tenant",
								true,
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOIDCIDPMigratedGoogle",
			args: args{
				event: getEvent(testEvent(
					instance.OIDCIDPMigratedGoogleEventType,
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "name",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
				), instance.OIDCIDPMigratedGoogleEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOIDCIDPMigratedGoogle,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (change_date, sequence, name, type, is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update, auto_linking) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE (id = $10) AND (instance_id = $11)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"name",
								domain.IDPTypeGoogle,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "DELETE FROM projections.idp_templates6_oidc WHERE (idp_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_google (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceOIDCIDPMigratedGoogle",
			args: args{
				event: getEvent(testEvent(
					org.OIDCIDPMigratedGoogleEventType,
					org.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "name",
	"clientId": "client_id",
	"clientSecret": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
	"scopes": ["profile"],
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
				), org.OIDCIDPMigratedGoogleEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOIDCIDPMigratedGoogle,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (change_date, sequence, name, type, is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update, auto_linking) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE (id = $10) AND (instance_id = $11)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"name",
								domain.IDPTypeGoogle,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "DELETE FROM projections.idp_templates6_oidc WHERE (idp_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_google (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.TextArray[string]{"profile"},
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
		})
	}
}

func TestIDPTemplateProjection_reducesOldConfig(t *testing.T) {
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
			name: "instance reduceOldConfigAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.IDPConfigAddedEventType,
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
			reduce: (&idpTemplateProjection{}).reduceOldConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-config-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"custom-zitadel-instance",
								domain.IdentityProviderTypeSystem,
								domain.IDPTypeUnspecified,
								true,
								true,
								true,
								false,
								domain.AutoLinkingOptionUnspecified,
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceOldConfigAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.IDPConfigAddedEventType,
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
			reduce: (&idpTemplateProjection{}).reduceOldConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-config-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"custom-zitadel-instance",
								domain.IdentityProviderTypeOrg,
								domain.IDPTypeUnspecified,
								true,
								true,
								true,
								false,
								domain.AutoLinkingOptionUnspecified,
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOldConfigChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.IDPConfigChangedEventType,
						instance.AggregateType,
						[]byte(`{
        "idpConfigId": "idp-config-id",
        "name": "custom-zitadel-instance",
        "stylingType": 1,
        "autoRegister": true
        }`),
					), instance.IDPConfigChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOldConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (name, is_auto_creation, change_date, sequence) = ($1, $2, $3, $4) WHERE (id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								"custom-zitadel-instance",
								true,
								anyArg{},
								uint64(15),
								"idp-config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceOldConfigChanged",
			args: args{
				event: getEvent(
					testEvent(
						org.IDPConfigChangedEventType,
						org.AggregateType,
						[]byte(`{
        "idpConfigId": "idp-config-id",
        "name": "custom-zitadel-instance",
        "stylingType": 1,
        "autoRegister": true
        }`),
					), org.IDPConfigChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOldConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (name, is_auto_creation, change_date, sequence) = ($1, $2, $3, $4) WHERE (id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								"custom-zitadel-instance",
								true,
								anyArg{},
								uint64(15),
								"idp-config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOldOIDCConfigAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.IDPOIDCConfigAddedEventType,
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
        "scopes": ["profile"],
        "idpDisplayNameMapping": 0,
        "usernameMapping": 1
        }`),
					), instance.IDPOIDCConfigAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOldOIDCConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (change_date, sequence, type) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.IDPTypeOIDC,
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_oidc (idp_id, instance_id, issuer, client_id, client_secret, scopes, id_token_mapping, use_pkce) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"instance-id",
								"issuer",
								"client-id",
								anyArg{},
								database.TextArray[string]{"profile"},
								true,
								false,
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceOldOIDCConfigAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.IDPOIDCConfigAddedEventType,
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
        "scopes": ["profile"],
        "idpDisplayNameMapping": 0,
        "usernameMapping": 1
        }`),
					), org.IDPOIDCConfigAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOldOIDCConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (change_date, sequence, type) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.IDPTypeOIDC,
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_oidc (idp_id, instance_id, issuer, client_id, client_secret, scopes, id_token_mapping, use_pkce) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"instance-id",
								"issuer",
								"client-id",
								anyArg{},
								database.TextArray[string]{"profile"},
								true,
								false,
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOldOIDCConfigChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.IDPOIDCConfigChangedEventType,
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
        "scopes": ["profile"],
        "idpDisplayNameMapping": 0,
        "usernameMapping": 1
        }`),
					), instance.IDPOIDCConfigChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOldOIDCConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_oidc SET (client_id, client_secret, issuer, scopes) = ($1, $2, $3, $4) WHERE (idp_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								"client-id",
								anyArg{},
								"issuer",
								database.TextArray[string]{"profile"},
								"idp-config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceOldOIDCConfigChanged",
			args: args{
				event: getEvent(
					testEvent(
						org.IDPOIDCConfigChangedEventType,
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
        "scopes": ["profile"],
        "idpDisplayNameMapping": 0,
        "usernameMapping": 1
        }`),
					), org.IDPOIDCConfigChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOldOIDCConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_oidc SET (client_id, client_secret, issuer, scopes) = ($1, $2, $3, $4) WHERE (idp_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								"client-id",
								anyArg{},
								"issuer",
								database.TextArray[string]{"profile"},
								"idp-config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOldJWTConfigAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.IDPJWTConfigAddedEventType,
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
			reduce: (&idpTemplateProjection{}).reduceOldJWTConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (change_date, sequence, type) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.IDPTypeJWT,
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_jwt (idp_id, instance_id, issuer, jwt_endpoint, keys_endpoint, header_name) VALUES ($1, $2, $3, $4, $5, $6)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"instance-id",
								"issuer",
								"https://api.zitadel.ch/jwt",
								"https://api.zitadel.ch/keys",
								"hodor",
							},
						},
					},
				},
			},
		},

		{
			name: "org reduceOldJWTConfigAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.IDPJWTConfigAddedEventType,
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
			reduce: (&idpTemplateProjection{}).reduceOldJWTConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (change_date, sequence, type) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.IDPTypeJWT,
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_jwt (idp_id, instance_id, issuer, jwt_endpoint, keys_endpoint, header_name) VALUES ($1, $2, $3, $4, $5, $6)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"instance-id",
								"issuer",
								"https://api.zitadel.ch/jwt",
								"https://api.zitadel.ch/keys",
								"hodor",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOldJWTConfigChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.IDPJWTConfigChangedEventType,
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
			reduce: (&idpTemplateProjection{}).reduceOldJWTConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_jwt SET (jwt_endpoint, keys_endpoint, header_name, issuer) = ($1, $2, $3, $4) WHERE (idp_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								"https://api.zitadel.ch/jwt",
								"https://api.zitadel.ch/keys",
								"hodor",
								"issuer",
								"idp-config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceOldJWTConfigChanged",
			args: args{
				event: getEvent(
					testEvent(
						org.IDPJWTConfigChangedEventType,
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
			reduce: (&idpTemplateProjection{}).reduceOldJWTConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_jwt SET (jwt_endpoint, keys_endpoint, header_name, issuer) = ($1, $2, $3, $4) WHERE (idp_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								"https://api.zitadel.ch/jwt",
								"https://api.zitadel.ch/keys",
								"hodor",
								"issuer",
								"idp-config-id",
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
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
				event: getEvent(
					testEvent(
						instance.JWTIDPAddedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"issuer": "issuer",
	"jwtEndpoint": "jwt",
	"keysEndpoint": "keys",
	"headerName": "header",
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.JWTIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceJWTIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"",
								domain.IdentityProviderTypeSystem,
								domain.IDPTypeJWT,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_jwt (idp_id, instance_id, issuer, jwt_endpoint, keys_endpoint, header_name) VALUES ($1, $2, $3, $4, $5, $6)",
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
			name: "org reduceJWTIDPAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.JWTIDPAddedEventType,
						org.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"issuer": "issuer",
	"jwtEndpoint": "jwt",
	"keysEndpoint": "keys",
	"headerName": "header",
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), org.JWTIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceJWTIDPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateInsertStmt,
							expectedArgs: []interface{}{
								"idp-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.IDPStateActive,
								"",
								domain.IdentityProviderTypeOrg,
								domain.IDPTypeJWT,
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates6_jwt (idp_id, instance_id, issuer, jwt_endpoint, keys_endpoint, header_name) VALUES ($1, $2, $3, $4, $5, $6)",
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
				event: getEvent(
					testEvent(
						instance.JWTIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"isCreationAllowed": true,
	"jwtEndpoint": "jwt"
}`),
					), instance.JWTIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceJWTIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: idpTemplateUpdateMinimalStmt,
							expectedArgs: []interface{}{
								true,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_jwt SET jwt_endpoint = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"jwt",
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
				event: getEvent(
					testEvent(
						instance.JWTIDPChangedEventType,
						instance.AggregateType,
						[]byte(`{
	"id": "idp-id",
	"issuer": "issuer",
	"jwtEndpoint": "jwt",
	"keysEndpoint": "keys",
	"headerName": "header",
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true,
	"autoLinkingOption": 1
}`),
					), instance.JWTIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceJWTIDPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates6 SET (is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update, auto_linking, change_date, sequence) = ($1, $2, $3, $4, $5, $6, $7) WHERE (id = $8) AND (instance_id = $9)",
							expectedArgs: []interface{}{
								true,
								true,
								true,
								true,
								domain.AutoLinkingOptionUsername,
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates6_jwt SET (jwt_endpoint, keys_endpoint, header_name, issuer) = ($1, $2, $3, $4) WHERE (idp_id = $5) AND (instance_id = $6)",
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
		})
	}
}

func stringToJSONByte(data string) string {
	jsondata, _ := json.Marshal([]byte(data))
	return string(jsondata)
}
