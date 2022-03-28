package projection

import (
	"testing"
	"time"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
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
						"passwordlessType": 1,
						"passwordCheckLifetime": 10000000,
						"externalLoginCheckLifetime": 10000000,
						"mfaInitSkipLifetime": 10000000,
						"secondFactorCheckLifetime": 10000000,
						"multiFactorCheckLifetime": 10000000
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
							expectedStmt: "INSERT INTO projections.login_policies (aggregate_id, instance_id, creation_date, change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, passwordless_type, is_default, hide_password_reset, password_check_lifetime, external_login_check_lifetime, mfa_init_skip_lifetime, second_factor_check_lifetime, multi_factor_check_lifetime) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
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
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
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
						"passwordCheckLifetime": 10000000,
						"externalLoginCheckLifetime": 10000000,
						"mfaInitSkipLifetime": 10000000,
						"secondFactorCheckLifetime": 10000000,
						"multiFactorCheckLifetime": 10000000
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
							expectedStmt: "UPDATE projections.login_policies SET (change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, passwordless_type, hide_password_reset, password_check_lifetime, external_login_check_lifetime, mfa_init_skip_lifetime, second_factor_check_lifetime, multi_factor_check_lifetime) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) WHERE (aggregate_id = $14)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								true,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
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
							expectedStmt: "UPDATE projections.login_policies SET (change_date, sequence, multi_factors) = ($1, $2, array_append(multi_factors, $3)) WHERE (aggregate_id = $4)",
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
							expectedStmt: "UPDATE projections.login_policies SET (change_date, sequence, multi_factors) = ($1, $2, array_remove(multi_factors, $3)) WHERE (aggregate_id = $4)",
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
							expectedStmt: "DELETE FROM projections.login_policies WHERE (aggregate_id = $1)",
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
							expectedStmt: "UPDATE projections.login_policies SET (change_date, sequence, second_factors) = ($1, $2, array_append(second_factors, $3)) WHERE (aggregate_id = $4)",
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
							expectedStmt: "UPDATE projections.login_policies SET (change_date, sequence, second_factors) = ($1, $2, array_remove(second_factors, $3)) WHERE (aggregate_id = $4)",
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
			name:   "instance.reduceLoginPolicyAdded",
			reduce: (&LoginPolicyProjection{}).reduceLoginPolicyAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.LoginPolicyAddedEventType),
					instance.AggregateType,
					[]byte(`{
						"allowUsernamePassword": true,
						"allowRegister": true,
						"allowExternalIdp": false,
						"forceMFA": false,
						"hidePasswordReset": true,
						"passwordlessType": 1,
						"passwordCheckLifetime": 10000000,
						"externalLoginCheckLifetime": 10000000,
						"mfaInitSkipLifetime": 10000000,
						"secondFactorCheckLifetime": 10000000,
						"multiFactorCheckLifetime": 10000000
			}`),
				), instance.LoginPolicyAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.login_policies (aggregate_id, instance_id, creation_date, change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, passwordless_type, is_default, hide_password_reset, password_check_lifetime, external_login_check_lifetime, mfa_init_skip_lifetime, second_factor_check_lifetime, multi_factor_check_lifetime) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
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
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
							},
						},
					},
				},
			},
		},
		{
			name:   "instance.reduceLoginPolicyChanged",
			reduce: (&LoginPolicyProjection{}).reduceLoginPolicyChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.LoginPolicyChangedEventType),
					instance.AggregateType,
					[]byte(`{
			"allowUsernamePassword": true,
			"allowRegister": true,
			"allowExternalIdp": true,
			"forceMFA": true,
			"hidePasswordReset": true,
			"passwordlessType": 1
			}`),
				), instance.LoginPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies SET (change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, passwordless_type, hide_password_reset) = ($1, $2, $3, $4, $5, $6, $7, $8) WHERE (aggregate_id = $9)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								true,
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance.reduceMFAAdded",
			reduce: (&LoginPolicyProjection{}).reduceMFAAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.LoginPolicyMultiFactorAddedEventType),
					instance.AggregateType,
					[]byte(`{
		"mfaType": 1
		}`),
				), instance.MultiFactorAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies SET (change_date, sequence, multi_factors) = ($1, $2, array_append(multi_factors, $3)) WHERE (aggregate_id = $4)",
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
			name:   "instance.reduceMFARemoved",
			reduce: (&LoginPolicyProjection{}).reduceMFARemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.LoginPolicyMultiFactorRemovedEventType),
					instance.AggregateType,
					[]byte(`{
			"mfaType": 1
			}`),
				), instance.MultiFactorRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies SET (change_date, sequence, multi_factors) = ($1, $2, array_remove(multi_factors, $3)) WHERE (aggregate_id = $4)",
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
			name:   "instance.reduce2FAAdded",
			reduce: (&LoginPolicyProjection{}).reduce2FAAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.LoginPolicySecondFactorAddedEventType),
					instance.AggregateType,
					[]byte(`{
			"mfaType": 2
			}`),
				), instance.SecondFactorAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies SET (change_date, sequence, second_factors) = ($1, $2, array_append(second_factors, $3)) WHERE (aggregate_id = $4)",
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
			name:   "instance.reduce2FARemoved",
			reduce: (&LoginPolicyProjection{}).reduce2FARemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.LoginPolicySecondFactorRemovedEventType),
					instance.AggregateType,
					[]byte(`{
			"mfaType": 2
			}`),
				), instance.SecondFactorRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       LoginPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies SET (change_date, sequence, second_factors) = ($1, $2, array_remove(second_factors, $3)) WHERE (aggregate_id = $4)",
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
