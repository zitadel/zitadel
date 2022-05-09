package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func TestLoginPolicyProjection_reduces(t *testing.T) {
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
			name: "org.reduceLoginPolicyAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LoginPolicyAddedEventType),
					org.AggregateType,
					[]byte(`{
	"allowUsernamePassword": true,
	"allowRegister": true,
	"allowExternalIdp": false,
	"forceMFA": false,
	"hidePasswordReset": true,
	"passwordlessType": 1
}`),
				), org.LoginPolicyAddedEventMapper),
			},
			reduce: (&LoginPolicyProjection{}).reduceLoginPolicyAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.login_policies (aggregate_id, creation_date, change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, passwordless_type, is_default, hide_password_reset, ignore_unknown_usernames, default_redirect_uri) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
								true,
								true,
								false,
								false,
								domain.PasswordlessTypeAllowed,
								false,
								true,
								false,
								"",
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceLoginPolicyAdded ignoreUnknownUsernames",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LoginPolicyAddedEventType),
					org.AggregateType,
					[]byte(`{
	"allowUsernamePassword": true,
	"allowRegister": true,
	"allowExternalIdp": false,
	"forceMFA": false,
	"hidePasswordReset": true,
	"passwordlessType": 1,
	"ignoreUnknownUsernames": true,
	"defaultRedirectURI": "https://example.com/redirect"
}`),
				), org.LoginPolicyAddedEventMapper),
			},
			reduce: (&LoginPolicyProjection{}).reduceLoginPolicyAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.login_policies (aggregate_id, creation_date, change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, passwordless_type, is_default, hide_password_reset, ignore_unknown_usernames, default_redirect_uri) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
								true,
								true,
								false,
								false,
								domain.PasswordlessTypeAllowed,
								false,
								true,
								true,
								"https://example.com/redirect",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceLoginPolicyChanged",
			reduce: (&LoginPolicyProjection{}).reduceLoginPolicyChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LoginPolicyChangedEventType),
					org.AggregateType,
					[]byte(`{
	"allowUsernamePassword": true,
	"allowRegister": true,
	"allowExternalIdp": true,
	"forceMFA": true,
	"hidePasswordReset": true,
	"passwordlessType": 1,
	"ignoreUnknownUsernames": true,
	"defaultRedirectURI": "https://example.com/redirect"
}`),
				), org.LoginPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_policies SET (change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, passwordless_type, hide_password_reset, ignore_unknown_usernames, default_redirect_uri) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) WHERE (aggregate_id = $11)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								true,
								true,
								"https://example.com/redirect",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceMFAAdded",
			reduce: (&LoginPolicyProjection{}).reduceMFAAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LoginPolicyMultiFactorAddedEventType),
					org.AggregateType,
					[]byte(`{
	"mfaType": 1
}`),
				), org.MultiFactorAddedEventEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_policies SET (change_date, sequence, multi_factors) = ($1, $2, array_append(multi_factors, $3)) WHERE (aggregate_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.MultiFactorTypeU2FWithPIN,
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceMFARemoved",
			reduce: (&LoginPolicyProjection{}).reduceMFARemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LoginPolicyMultiFactorRemovedEventType),
					org.AggregateType,
					[]byte(`{
			"mfaType": 1
			}`),
				), org.MultiFactorRemovedEventEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_policies SET (change_date, sequence, multi_factors) = ($1, $2, array_remove(multi_factors, $3)) WHERE (aggregate_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.MultiFactorTypeU2FWithPIN,
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceLoginPolicyRemoved",
			reduce: (&LoginPolicyProjection{}).reduceLoginPolicyRemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LoginPolicyRemovedEventType),
					org.AggregateType,
					nil,
				), org.LoginPolicyRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.login_policies WHERE (aggregate_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduce2FAAdded",
			reduce: (&LoginPolicyProjection{}).reduce2FAAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LoginPolicySecondFactorAddedEventType),
					org.AggregateType,
					[]byte(`{
			"mfaType": 2
			}`),
				), org.SecondFactorAddedEventEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_policies SET (change_date, sequence, second_factors) = ($1, $2, array_append(second_factors, $3)) WHERE (aggregate_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SecondFactorTypeU2F,
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduce2FARemoved",
			reduce: (&LoginPolicyProjection{}).reduce2FARemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LoginPolicySecondFactorRemovedEventType),
					org.AggregateType,
					[]byte(`{
			"mfaType": 2
			}`),
				), org.SecondFactorRemovedEventEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_policies SET (change_date, sequence, second_factors) = ($1, $2, array_remove(second_factors, $3)) WHERE (aggregate_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SecondFactorTypeU2F,
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceLoginPolicyAdded",
			reduce: (&LoginPolicyProjection{}).reduceLoginPolicyAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LoginPolicyAddedEventType),
					iam.AggregateType,
					[]byte(`{
			"allowUsernamePassword": true,
			"allowRegister": true,
			"allowExternalIdp": false,
			"forceMFA": false,
			"hidePasswordReset": true,
			"passwordlessType": 1
			}`),
				), iam.LoginPolicyAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.login_policies (aggregate_id, creation_date, change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, passwordless_type, is_default, hide_password_reset, ignore_unknown_usernames, default_redirect_uri) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
								true,
								true,
								false,
								false,
								domain.PasswordlessTypeAllowed,
								true,
								true,
								false,
								"",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceLoginPolicyAdded ignoreUnknownUsernames",
			reduce: (&LoginPolicyProjection{}).reduceLoginPolicyAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LoginPolicyAddedEventType),
					iam.AggregateType,
					[]byte(`{
			"allowUsernamePassword": true,
			"allowRegister": true,
			"allowExternalIdp": false,
			"forceMFA": false,
			"hidePasswordReset": true,
			"passwordlessType": 1,
			"ignoreUnknownUsernames": true,
			"defaultRedirectURI": "https://example.com/redirect"
			}`),
				), iam.LoginPolicyAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.login_policies (aggregate_id, creation_date, change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, passwordless_type, is_default, hide_password_reset, ignore_unknown_usernames, default_redirect_uri) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
								true,
								true,
								false,
								false,
								domain.PasswordlessTypeAllowed,
								true,
								true,
								true,
								"https://example.com/redirect",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceLoginPolicyChanged",
			reduce: (&LoginPolicyProjection{}).reduceLoginPolicyChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LoginPolicyChangedEventType),
					iam.AggregateType,
					[]byte(`{
			"allowUsernamePassword": true,
			"allowRegister": true,
			"allowExternalIdp": true,
			"forceMFA": true,
			"hidePasswordReset": true,
			"passwordlessType": 1,
			"defaultRedirectURI": "https://example.com/redirect"
			}`),
				), iam.LoginPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_policies SET (change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, passwordless_type, hide_password_reset, default_redirect_uri) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE (aggregate_id = $10)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								true,
								"https://example.com/redirect",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceMFAAdded",
			reduce: (&LoginPolicyProjection{}).reduceMFAAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LoginPolicyMultiFactorAddedEventType),
					iam.AggregateType,
					[]byte(`{
		"mfaType": 1
		}`),
				), iam.MultiFactorAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_policies SET (change_date, sequence, multi_factors) = ($1, $2, array_append(multi_factors, $3)) WHERE (aggregate_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.MultiFactorTypeU2FWithPIN,
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceMFARemoved",
			reduce: (&LoginPolicyProjection{}).reduceMFARemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LoginPolicyMultiFactorRemovedEventType),
					iam.AggregateType,
					[]byte(`{
			"mfaType": 1
			}`),
				), iam.MultiFactorRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_policies SET (change_date, sequence, multi_factors) = ($1, $2, array_remove(multi_factors, $3)) WHERE (aggregate_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.MultiFactorTypeU2FWithPIN,
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduce2FAAdded",
			reduce: (&LoginPolicyProjection{}).reduce2FAAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LoginPolicySecondFactorAddedEventType),
					iam.AggregateType,
					[]byte(`{
			"mfaType": 2
			}`),
				), iam.SecondFactorAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_policies SET (change_date, sequence, second_factors) = ($1, $2, array_append(second_factors, $3)) WHERE (aggregate_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SecondFactorTypeU2F,
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduce2FARemoved",
			reduce: (&LoginPolicyProjection{}).reduce2FARemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LoginPolicySecondFactorRemovedEventType),
					iam.AggregateType,
					[]byte(`{
			"mfaType": 2
			}`),
				), iam.SecondFactorRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.login_policies SET (change_date, sequence, second_factors) = ($1, $2, array_remove(second_factors, $3)) WHERE (aggregate_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SecondFactorTypeU2F,
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
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, tt.want)
		})
	}
}
