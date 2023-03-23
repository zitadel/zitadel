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
	"github.com/zitadel/zitadel/internal/repository/org"
)

var (
	idpTemplateInsertStmt = `INSERT INTO projections.idp_templates4` +
		` (id, creation_date, change_date, sequence, resource_owner, instance_id, state, name, owner_type, type, is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update)` +
		` VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	idpTemplateUpdateMinimalStmt = `UPDATE projections.idp_templates4 SET (is_creation_allowed, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)`
	idpTemplateUpdateStmt        = `UPDATE projections.idp_templates4 SET (name, is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update, change_date, sequence)` +
		` = ($1, $2, $3, $4, $5, $6, $7) WHERE (id = $8) AND (instance_id = $9)`
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
				event: getEvent(testEvent(
					repository.EventType(instance.InstanceRemovedEventType),
					instance.AggregateType,
					nil,
				), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(IDPInstanceIDCol),
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_templates4 WHERE (instance_id = $1)",
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
				event: getEvent(testEvent(
					repository.EventType(org.OrgRemovedEventType),
					org.AggregateType,
					nil,
				), org.OrgRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates4 SET (change_date, sequence, owner_removed) = ($1, $2, $3) WHERE (instance_id = $4) AND (resource_owner = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
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
				event: getEvent(testEvent(
					repository.EventType(org.IDPRemovedEventType),
					org.AggregateType,
					[]byte(`{
	"id": "idp-id"
}`),
				), org.IDPRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.idp_templates4 WHERE (id = $1) AND (instance_id = $2)",
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
			if !errors.IsErrorInvalidArgument(err) {
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
				event: getEvent(testEvent(
					repository.EventType(instance.OAuthIDPAddedEventType),
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
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true
}`),
				), instance.OAuthIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOAuthIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_oauth2 (idp_id, instance_id, client_id, client_secret, authorization_endpoint, token_endpoint, user_endpoint, scopes, id_attribute) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								"auth",
								"token",
								"user",
								database.StringArray{"profile"},
								"id-attribute",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceOAuthIDPAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OAuthIDPAddedEventType),
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
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true
}`),
				), org.OAuthIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOAuthIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_oauth2 (idp_id, instance_id, client_id, client_secret, authorization_endpoint, token_endpoint, user_endpoint, scopes, id_attribute) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								"auth",
								"token",
								"user",
								database.StringArray{"profile"},
								"id-attribute",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOAuthIDPChanged minimal",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.OAuthIDPChangedEventType),
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
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							expectedStmt: "UPDATE projections.idp_templates4_oauth2 SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
				event: getEvent(testEvent(
					repository.EventType(instance.OAuthIDPChangedEventType),
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
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true
}`),
				), instance.OAuthIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOAuthIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates4_oauth2 SET (client_id, client_secret, authorization_endpoint, token_endpoint, user_endpoint, scopes, id_attribute) = ($1, $2, $3, $4, $5, $6, $7) WHERE (idp_id = $8) AND (instance_id = $9)",
							expectedArgs: []interface{}{
								"client_id",
								anyArg{},
								"auth",
								"token",
								"user",
								database.StringArray{"profile"},
								"id-attribute",
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
			if !errors.IsErrorInvalidArgument(err) {
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
				event: getEvent(testEvent(
					repository.EventType(instance.AzureADIDPAddedEventType),
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
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_azure (idp_id, instance_id, client_id, client_secret, scopes, tenant, is_email_verified) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.StringArray(nil),
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
				event: getEvent(testEvent(
					repository.EventType(instance.AzureADIDPAddedEventType),
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
	"isAutoUpdate": true
}`),
				), instance.AzureADIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceAzureADIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_azure (idp_id, instance_id, client_id, client_secret, scopes, tenant, is_email_verified) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.StringArray{"profile"},
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
				event: getEvent(testEvent(
					repository.EventType(org.AzureADIDPAddedEventType),
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
	"isAutoUpdate": true
}`),
				), org.AzureADIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceAzureADIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_azure (idp_id, instance_id, client_id, client_secret, scopes, tenant, is_email_verified) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								database.StringArray{"profile"},
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
				event: getEvent(testEvent(
					repository.EventType(instance.AzureADIDPChangedEventType),
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
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							expectedStmt: "UPDATE projections.idp_templates4_azure SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
				event: getEvent(testEvent(
					repository.EventType(instance.AzureADIDPChangedEventType),
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
	"isAutoUpdate": true
}`),
				), instance.AzureADIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceAzureADIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates4_azure SET (client_id, client_secret, scopes, tenant, is_email_verified) = ($1, $2, $3, $4, $5) WHERE (idp_id = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								"client_id",
								anyArg{},
								database.StringArray{"profile"},
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
			if !errors.IsErrorInvalidArgument(err) {
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
				event: getEvent(testEvent(
					repository.EventType(instance.GitHubIDPAddedEventType),
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
	"isAutoUpdate": true
}`),
				), instance.GitHubIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitHubIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_github (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
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
			name: "org reduceGitHubIDPAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.GitHubIDPAddedEventType),
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
	"isAutoUpdate": true
}`),
				), org.GitHubIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitHubIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_github (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
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
			name: "instance reduceGitHubIDPChanged minimal",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.GitHubIDPChangedEventType),
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
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							expectedStmt: "UPDATE projections.idp_templates4_github SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
				event: getEvent(testEvent(
					repository.EventType(instance.GitHubIDPChangedEventType),
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
	"isAutoUpdate": true
}`),
				), instance.GitHubIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitHubIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates4_github SET (client_id, client_secret, scopes) = ($1, $2, $3) WHERE (idp_id = $4) AND (instance_id = $5)",
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
			if !errors.IsErrorInvalidArgument(err) {
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
				event: getEvent(testEvent(
					repository.EventType(instance.GitHubEnterpriseIDPAddedEventType),
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
	"isAutoUpdate": true
}`),
				), instance.GitHubEnterpriseIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitHubEnterpriseIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_github_enterprise (idp_id, instance_id, client_id, client_secret, authorization_endpoint, token_endpoint, user_endpoint, scopes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								"auth",
								"token",
								"user",
								database.StringArray{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceGitHubEnterpriseIDPAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.GitHubEnterpriseIDPAddedEventType),
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
	"isAutoUpdate": true
}`),
				), org.GitHubEnterpriseIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitHubEnterpriseIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_github_enterprise (idp_id, instance_id, client_id, client_secret, authorization_endpoint, token_endpoint, user_endpoint, scopes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"client_id",
								anyArg{},
								"auth",
								"token",
								"user",
								database.StringArray{"profile"},
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceGitHubEnterpriseIDPChanged minimal",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.GitHubEnterpriseIDPChangedEventType),
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
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							expectedStmt: "UPDATE projections.idp_templates4_github_enterprise SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
				event: getEvent(testEvent(
					repository.EventType(instance.GitHubEnterpriseIDPChangedEventType),
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
	"isAutoUpdate": true
}`),
				), instance.GitHubEnterpriseIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitHubEnterpriseIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates4_github_enterprise SET (client_id, client_secret, authorization_endpoint, token_endpoint, user_endpoint, scopes) = ($1, $2, $3, $4, $5, $6) WHERE (idp_id = $7) AND (instance_id = $8)",
							expectedArgs: []interface{}{
								"client_id",
								anyArg{},
								"auth",
								"token",
								"user",
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
			if !errors.IsErrorInvalidArgument(err) {
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
				event: getEvent(testEvent(
					repository.EventType(instance.GitLabIDPAddedEventType),
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
	"isAutoUpdate": true
}`),
				), instance.GitLabIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitLabIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_gitlab (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
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
			name: "org reduceGitLabIDPAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.GitLabIDPAddedEventType),
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
	"isAutoUpdate": true
}`),
				), org.GitLabIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitLabIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_gitlab (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
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
			name: "instance reduceGitLabIDPChanged minimal",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.GitLabIDPChangedEventType),
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
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							expectedStmt: "UPDATE projections.idp_templates4_gitlab SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
				event: getEvent(testEvent(
					repository.EventType(instance.GitLabIDPChangedEventType),
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
	"isAutoUpdate": true
}`),
				), instance.GitLabIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitLabIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates4_gitlab SET (client_id, client_secret, scopes) = ($1, $2, $3) WHERE (idp_id = $4) AND (instance_id = $5)",
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
			if !errors.IsErrorInvalidArgument(err) {
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
				event: getEvent(testEvent(
					repository.EventType(instance.GitLabSelfHostedIDPAddedEventType),
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
	"isAutoUpdate": true
}`),
				), instance.GitLabSelfHostedIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitLabSelfHostedIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_gitlab_self_hosted (idp_id, instance_id, issuer, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5, $6)",
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
			name: "org reduceGitLabSelfHostedIDPAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.GitLabSelfHostedIDPAddedEventType),
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
	"isAutoUpdate": true
}`),
				), org.GitLabSelfHostedIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitLabSelfHostedIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_gitlab_self_hosted (idp_id, instance_id, issuer, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5, $6)",
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
			name: "instance reduceGitLabSelfHostedIDPChanged minimal",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.GitLabSelfHostedIDPChangedEventType),
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
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							expectedStmt: "UPDATE projections.idp_templates4_gitlab_self_hosted SET issuer = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
				event: getEvent(testEvent(
					repository.EventType(instance.GitLabSelfHostedIDPChangedEventType),
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
	"isAutoUpdate": true
}`),
				), instance.GitLabSelfHostedIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGitLabSelfHostedIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates4_gitlab_self_hosted SET (issuer, client_id, client_secret, scopes) = ($1, $2, $3, $4) WHERE (idp_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								"issuer",
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
			if !errors.IsErrorInvalidArgument(err) {
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
				event: getEvent(testEvent(
					repository.EventType(instance.GoogleIDPAddedEventType),
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_google (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
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
			name: "org reduceGoogleIDPAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.GoogleIDPAddedEventType),
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
	"isAutoUpdate": true
}`),
				), org.GoogleIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceGoogleIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_google (idp_id, instance_id, client_id, client_secret, scopes) VALUES ($1, $2, $3, $4, $5)",
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
	"clientId": "id"
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
							expectedStmt: "UPDATE projections.idp_templates4_google SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
							expectedStmt: idpTemplateUpdateStmt,
							expectedArgs: []interface{}{
								"name",
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
							expectedStmt: "UPDATE projections.idp_templates4_google SET (client_id, client_secret, scopes) = ($1, $2, $3) WHERE (idp_id = $4) AND (instance_id = $5)",
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
			if !errors.IsErrorInvalidArgument(err) {
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
				event: getEvent(testEvent(
					repository.EventType(instance.LDAPIDPAddedEventType),
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"host": "host",
	"port": "port",
	"tls": true,
	"baseDN": "base",
	"userObjectClass": "user",
	"userUniqueAttribute": "uid",
	"admin": "admin",
	"password": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
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
	"isAutoUpdate": true
}`),
				), instance.LDAPIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceLDAPIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_ldap (idp_id, instance_id, host, port, tls, base_dn, user_object_class, user_unique_attribute, admin, password, id_attribute, first_name_attribute, last_name_attribute, display_name_attribute, nick_name_attribute, preferred_username_attribute, email_attribute, email_verified, phone_attribute, phone_verified_attribute, preferred_language_attribute, avatar_url_attribute, profile_attribute) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"host",
								"port",
								true,
								"base",
								"user",
								"uid",
								"admin",
								anyArg{},
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
				event: getEvent(testEvent(
					repository.EventType(org.LDAPIDPAddedEventType),
					org.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"host": "host",
	"port": "port",
	"tls": true,
	"baseDN": "base",
	"userObjectClass": "user",
	"userUniqueAttribute": "uid",
	"admin": "admin",
	"password": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
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
	"isAutoUpdate": true
}`),
				), org.LDAPIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceLDAPIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_ldap (idp_id, instance_id, host, port, tls, base_dn, user_object_class, user_unique_attribute, admin, password, id_attribute, first_name_attribute, last_name_attribute, display_name_attribute, nick_name_attribute, preferred_username_attribute, email_attribute, email_verified, phone_attribute, phone_verified_attribute, preferred_language_attribute, avatar_url_attribute, profile_attribute) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"host",
								"port",
								true,
								"base",
								"user",
								"uid",
								"admin",
								anyArg{},
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
				event: getEvent(testEvent(
					repository.EventType(instance.LDAPIDPChangedEventType),
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"host": "host"
}`),
				), instance.LDAPIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceLDAPIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates4 SET (name, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								"custom-zitadel-instance",
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates4_ldap SET host = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"host",
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
				event: getEvent(testEvent(
					repository.EventType(instance.LDAPIDPChangedEventType),
					instance.AggregateType,
					[]byte(`{
	"id": "idp-id",
	"name": "custom-zitadel-instance",
	"host": "host",
	"port": "port",
	"tls": true,
	"baseDN": "base",
	"userObjectClass": "user",
	"userUniqueAttribute": "uid",
	"admin": "admin",
	"password": {
        "cryptoType": 0,
        "algorithm": "RSA-265",
        "keyId": "key-id"
    },
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
	"isAutoUpdate": true
}`),
				), instance.LDAPIDPChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceLDAPIDPChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
								anyArg{},
								uint64(15),
								"idp-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates4_ldap SET (host, port, tls, base_dn, user_object_class, user_unique_attribute, admin, password, id_attribute, first_name_attribute, last_name_attribute, display_name_attribute, nick_name_attribute, preferred_username_attribute, email_attribute, email_verified, phone_attribute, phone_verified_attribute, preferred_language_attribute, avatar_url_attribute, profile_attribute) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21) WHERE (idp_id = $22) AND (instance_id = $23)",
							expectedArgs: []interface{}{
								"host",
								"port",
								true,
								"base",
								"user",
								"uid",
								"admin",
								anyArg{},
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
				event: getEvent(testEvent(
					repository.EventType(org.OrgRemovedEventType),
					org.AggregateType,
					nil,
				), org.OrgRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates4 SET (change_date, sequence, owner_removed) = ($1, $2, $3) WHERE (instance_id = $4) AND (resource_owner = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
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
			if !errors.IsErrorInvalidArgument(err) {
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
				event: getEvent(testEvent(
					repository.EventType(instance.OIDCIDPAddedEventType),
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_oidc (idp_id, instance_id, issuer, client_id, client_secret, scopes, id_token_mapping) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"issuer",
								"client_id",
								anyArg{},
								database.StringArray{"profile"},
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
				event: getEvent(testEvent(
					repository.EventType(org.OIDCIDPAddedEventType),
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
	"isCreationAllowed": true,
	"isLinkingAllowed": true,
	"isAutoCreation": true,
	"isAutoUpdate": true
}`),
				), org.OIDCIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOIDCIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_oidc (idp_id, instance_id, issuer, client_id, client_secret, scopes, id_token_mapping) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"idp-id",
								"instance-id",
								"issuer",
								"client_id",
								anyArg{},
								database.StringArray{"profile"},
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
				event: getEvent(testEvent(
					repository.EventType(instance.OIDCIDPChangedEventType),
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
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							expectedStmt: "UPDATE projections.idp_templates4_oidc SET client_id = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
				event: getEvent(testEvent(
					repository.EventType(instance.OIDCIDPChangedEventType),
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
							expectedStmt: idpTemplateUpdateStmt,
							expectedArgs: []interface{}{
								"name",
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
							expectedStmt: "UPDATE projections.idp_templates4_oidc SET (client_id, client_secret, issuer, scopes, id_token_mapping) = ($1, $2, $3, $4, $5) WHERE (idp_id = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								"client_id",
								anyArg{},
								"issuer",
								database.StringArray{"profile"},
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
			if !errors.IsErrorInvalidArgument(err) {
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
			reduce: (&idpTemplateProjection{}).reduceOldConfigAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceOldConfigAdded",
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
			reduce: (&idpTemplateProjection{}).reduceOldConfigAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOldConfigChanged",
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
			reduce: (&idpTemplateProjection{}).reduceOldConfigChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates4 SET (name, is_auto_creation, change_date, sequence) = ($1, $2, $3, $4) WHERE (id = $5) AND (instance_id = $6)",
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
			reduce: (&idpTemplateProjection{}).reduceOldConfigChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates4 SET (name, is_auto_creation, change_date, sequence) = ($1, $2, $3, $4) WHERE (id = $5) AND (instance_id = $6)",
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
        "scopes": ["profile"],
        "idpDisplayNameMapping": 0,
        "usernameMapping": 1
        }`),
				), instance.IDPOIDCConfigAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOldOIDCConfigAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates4 SET (change_date, sequence, type) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.IDPTypeOIDC,
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_oidc (idp_id, instance_id, issuer, client_id, client_secret, scopes, id_token_mapping) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"instance-id",
								"issuer",
								"client-id",
								anyArg{},
								database.StringArray{"profile"},
								true,
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceOldOIDCConfigAdded",
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
        "scopes": ["profile"],
        "idpDisplayNameMapping": 0,
        "usernameMapping": 1
        }`),
				), org.IDPOIDCConfigAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOldOIDCConfigAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates4 SET (change_date, sequence, type) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.IDPTypeOIDC,
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_oidc (idp_id, instance_id, issuer, client_id, client_secret, scopes, id_token_mapping) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"idp-config-id",
								"instance-id",
								"issuer",
								"client-id",
								anyArg{},
								database.StringArray{"profile"},
								true,
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOldOIDCConfigChanged",
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
        "scopes": ["profile"],
        "idpDisplayNameMapping": 0,
        "usernameMapping": 1
        }`),
				), instance.IDPOIDCConfigChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOldOIDCConfigChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates4 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates4_oidc SET (client_id, client_secret, issuer, scopes) = ($1, $2, $3, $4) WHERE (idp_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								"client-id",
								anyArg{},
								"issuer",
								database.StringArray{"profile"},
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
        "scopes": ["profile"],
        "idpDisplayNameMapping": 0,
        "usernameMapping": 1
        }`),
				), org.IDPOIDCConfigChangedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceOldOIDCConfigChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates4 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates4_oidc SET (client_id, client_secret, issuer, scopes) = ($1, $2, $3, $4) WHERE (idp_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								"client-id",
								anyArg{},
								"issuer",
								database.StringArray{"profile"},
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
			reduce: (&idpTemplateProjection{}).reduceOldJWTConfigAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates4 SET (change_date, sequence, type) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.IDPTypeJWT,
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_jwt (idp_id, instance_id, issuer, jwt_endpoint, keys_endpoint, header_name) VALUES ($1, $2, $3, $4, $5, $6)",
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
			reduce: (&idpTemplateProjection{}).reduceOldJWTConfigAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates4 SET (change_date, sequence, type) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.IDPTypeJWT,
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_jwt (idp_id, instance_id, issuer, jwt_endpoint, keys_endpoint, header_name) VALUES ($1, $2, $3, $4, $5, $6)",
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
			reduce: (&idpTemplateProjection{}).reduceOldJWTConfigChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates4 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates4_jwt SET (jwt_endpoint, keys_endpoint, header_name, issuer) = ($1, $2, $3, $4) WHERE (idp_id = $5) AND (instance_id = $6)",
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
			reduce: (&idpTemplateProjection{}).reduceOldJWTConfigChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.idp_templates4 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"idp-config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.idp_templates4_jwt SET (jwt_endpoint, keys_endpoint, header_name, issuer) = ($1, $2, $3, $4) WHERE (idp_id = $5) AND (instance_id = $6)",
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
			if !errors.IsErrorInvalidArgument(err) {
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
				event: getEvent(testEvent(
					repository.EventType(instance.JWTIDPAddedEventType),
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_jwt (idp_id, instance_id, issuer, jwt_endpoint, keys_endpoint, header_name) VALUES ($1, $2, $3, $4, $5, $6)",
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
				event: getEvent(testEvent(
					repository.EventType(org.JWTIDPAddedEventType),
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
	"isAutoUpdate": true
}`),
				), org.JWTIDPAddedEventMapper),
			},
			reduce: (&idpTemplateProjection{}).reduceJWTIDPAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
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
							},
						},
						{
							expectedStmt: "INSERT INTO projections.idp_templates4_jwt (idp_id, instance_id, issuer, jwt_endpoint, keys_endpoint, header_name) VALUES ($1, $2, $3, $4, $5, $6)",
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
	"isCreationAllowed": true,
	"jwtEndpoint": "jwt"
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
							expectedStmt: "UPDATE projections.idp_templates4_jwt SET jwt_endpoint = $1 WHERE (idp_id = $2) AND (instance_id = $3)",
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
				event: getEvent(testEvent(
					repository.EventType(instance.JWTIDPChangedEventType),
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
							expectedStmt: "UPDATE projections.idp_templates4 SET (is_creation_allowed, is_linking_allowed, is_auto_creation, is_auto_update, change_date, sequence) = ($1, $2, $3, $4, $5, $6) WHERE (id = $7) AND (instance_id = $8)",
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
							expectedStmt: "UPDATE projections.idp_templates4_jwt SET (jwt_endpoint, keys_endpoint, header_name, issuer) = ($1, $2, $3, $4) WHERE (idp_id = $5) AND (instance_id = $6)",
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
			if !errors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, IDPTemplateTable, tt.want)
		})
	}
}
