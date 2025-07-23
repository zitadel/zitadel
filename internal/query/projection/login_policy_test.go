package projection

import (
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
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
			name: "org reduceLoginPolicyAdded without forceMFALocalOnly",
			args: args{
				event: getEvent(
					testEvent(
						org.LoginPolicyAddedEventType,
						org.AggregateType,
						[]byte(`{
						"allowUsernamePassword": true,
						"allowRegister": true,
						"allowExternalIdp": false,
						"forceMFA": false,
						"hidePasswordReset": true,
						"ignoreUnknownUsernames": true,
						"allowDomainDiscovery": true,
						"disableLoginWithEmail": true,
						"disableLoginWithPhone": true,
						"passwordlessType": 1,
						"defaultRedirectURI": "https://example.com/redirect",
						"passwordCheckLifetime": 10000000,
						"externalLoginCheckLifetime": 10000000,
						"mfaInitSkipLifetime": 10000000,
						"secondFactorCheckLifetime": 10000000,
						"multiFactorCheckLifetime": 10000000,
						"enableRegistrationCaptcha": false,
						"enableLoginCaptcha": false,
						"captchaType": 0,
						"captchaSiteKey": "",
						"captchaSecretKey": ""
					}`),
					), org.LoginPolicyAddedEventMapper),
			},
			reduce: (&loginPolicyProjection{}).reduceLoginPolicyAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.login_policies5 (aggregate_id, instance_id, creation_date, change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, force_mfa_local_only, passwordless_type, is_default, hide_password_reset, ignore_unknown_usernames, allow_domain_discovery, disable_login_with_email, disable_login_with_phone, default_redirect_uri, password_check_lifetime, external_login_check_lifetime, mfa_init_skip_lifetime, second_factor_check_lifetime, multi_factor_check_lifetime, enable_registration_captcha, enable_login_captcha, captcha_type, captcha_site_key, captcha_secret_key) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)",
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
								false,
								domain.PasswordlessTypeAllowed,
								false,
								true,
								true,
								true,
								true,
								true,
								"https://example.com/redirect",
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								false,
								false,
								domain.CaptchaTypeDisabled,
								"",
								"",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceLoginPolicyAdded",
			args: args{
				event: getEvent(testEvent(
					org.LoginPolicyAddedEventType,
					org.AggregateType,
					[]byte(`{
						"allowUsernamePassword": true,
						"allowRegister": true,
						"allowExternalIdp": true,
						"forceMFA": true,
						"forceMFALocalOnly": true,
						"hidePasswordReset": true,
						"ignoreUnknownUsernames": true,
						"allowDomainDiscovery": true,
						"disableLoginWithEmail": true,
						"disableLoginWithPhone": true,
						"passwordlessType": 1,
						"defaultRedirectURI": "https://example.com/redirect",
						"passwordCheckLifetime": 10000000,
						"externalLoginCheckLifetime": 10000000,
						"mfaInitSkipLifetime": 10000000,
						"secondFactorCheckLifetime": 10000000,
						"multiFactorCheckLifetime": 10000000,
						"enableRegistrationCaptcha": false,
						"enableLoginCaptcha": false,
						"captchaType": 0,
						"captchaSiteKey": "",
						"captchaSecretKey": ""
					}`),
				), org.LoginPolicyAddedEventMapper),
			},
			reduce: (&loginPolicyProjection{}).reduceLoginPolicyAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.login_policies5 (aggregate_id, instance_id, creation_date, change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, force_mfa_local_only, passwordless_type, is_default, hide_password_reset, ignore_unknown_usernames, allow_domain_discovery, disable_login_with_email, disable_login_with_phone, default_redirect_uri, password_check_lifetime, external_login_check_lifetime, mfa_init_skip_lifetime, second_factor_check_lifetime, multi_factor_check_lifetime, enable_registration_captcha, enable_login_captcha, captcha_type, captcha_site_key, captcha_secret_key) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								false,
								true,
								true,
								true,
								true,
								true,
								"https://example.com/redirect",
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								false,
								false,
								domain.CaptchaTypeDisabled,
								"",
								"",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceLoginPolicyChanged",
			reduce: (&loginPolicyProjection{}).reduceLoginPolicyChanged,
			args: args{
				event: getEvent(
					testEvent(
						org.LoginPolicyChangedEventType,
						org.AggregateType,
						[]byte(`{
						"allowUsernamePassword": true,
						"allowRegister": true,
						"allowExternalIdp": true,
						"forceMFA": true,
						"forceMFALocalOnly": true,
						"hidePasswordReset": true,
						"ignoreUnknownUsernames": true,
						"allowDomainDiscovery": true,
						"disableLoginWithEmail": true,
						"disableLoginWithPhone": true,
						"passwordlessType": 1,
						"defaultRedirectURI": "https://example.com/redirect",
						"passwordCheckLifetime": 10000000,
						"externalLoginCheckLifetime": 10000000,
						"mfaInitSkipLifetime": 10000000,
						"secondFactorCheckLifetime": 10000000,
						"multiFactorCheckLifetime": 10000000,
						"enableRegistrationCaptcha": false,
						"enableLoginCaptcha": false,
						"captchaType": 0,
						"captchaSiteKey": "",
						"captchaSecretKey": ""
					}`),
					), org.LoginPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies5 SET (change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, force_mfa_local_only, passwordless_type, hide_password_reset, ignore_unknown_usernames, allow_domain_discovery, disable_login_with_email, disable_login_with_phone, default_redirect_uri, password_check_lifetime, external_login_check_lifetime, mfa_init_skip_lifetime, second_factor_check_lifetime, multi_factor_check_lifetime, enable_registration_captcha, enable_login_captcha, captcha_type, captcha_site_key, captcha_secret_key) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24) WHERE (aggregate_id = $25) AND (instance_id = $26)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								true,
								true,
								true,
								true,
								true,
								"https://example.com/redirect",
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								false,
								false,
								domain.CaptchaTypeDisabled,
								"",
								"",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceMFAAdded",
			reduce: (&loginPolicyProjection{}).reduceMFAAdded,
			args: args{
				event: getEvent(
					testEvent(
						org.LoginPolicyMultiFactorAddedEventType,
						org.AggregateType,
						[]byte(`{
	"mfaType": 1
}`),
					), org.MultiFactorAddedEventEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies5 SET (change_date, sequence, multi_factors) = ($1, $2, array_append(multi_factors, $3)) WHERE (aggregate_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.MultiFactorTypeU2FWithPIN,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceMFARemoved",
			reduce: (&loginPolicyProjection{}).reduceMFARemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.LoginPolicyMultiFactorRemovedEventType,
						org.AggregateType,
						[]byte(`{
			"mfaType": 1
			}`),
					), org.MultiFactorRemovedEventEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies5 SET (change_date, sequence, multi_factors) = ($1, $2, array_remove(multi_factors, $3)) WHERE (aggregate_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.MultiFactorTypeU2FWithPIN,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceLoginPolicyRemoved",
			reduce: (&loginPolicyProjection{}).reduceLoginPolicyRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.LoginPolicyRemovedEventType,
						org.AggregateType,
						nil,
					), org.LoginPolicyRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.login_policies5 WHERE (aggregate_id = $1) AND (instance_id = $2)",
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
			name:   "org reduceSecondFactorAdded",
			reduce: (&loginPolicyProjection{}).reduceSecondFactorAdded,
			args: args{
				event: getEvent(
					testEvent(
						org.LoginPolicySecondFactorAddedEventType,
						org.AggregateType,
						[]byte(`{
			"mfaType": 2
			}`),
					), org.SecondFactorAddedEventEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies5 SET (change_date, sequence, second_factors) = ($1, $2, array_append(second_factors, $3)) WHERE (aggregate_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SecondFactorTypeU2F,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceSecondFactorRemoved",
			reduce: (&loginPolicyProjection{}).reduceSecondFactorRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.LoginPolicySecondFactorRemovedEventType,
						org.AggregateType,
						[]byte(`{
			"mfaType": 2
			}`),
					), org.SecondFactorRemovedEventEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies5 SET (change_date, sequence, second_factors) = ($1, $2, array_remove(second_factors, $3)) WHERE (aggregate_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SecondFactorTypeU2F,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceLoginPolicyAdded",
			reduce: (&loginPolicyProjection{}).reduceLoginPolicyAdded,
			args: args{
				event: getEvent(
					testEvent(
						instance.LoginPolicyAddedEventType,
						instance.AggregateType,
						[]byte(`{
						"allowUsernamePassword": true,
						"allowRegister": true,
						"allowExternalIdp": true,
						"forceMFA": true,
						"forceMFALocalOnly": true,
						"hidePasswordReset": true,
						"ignoreUnknownUsernames": true,
						"allowDomainDiscovery": true,
						"disableLoginWithEmail": true,
						"disableLoginWithPhone": true,
						"passwordlessType": 1,
						"defaultRedirectURI": "https://example.com/redirect",
						"passwordCheckLifetime": 10000000,
						"externalLoginCheckLifetime": 10000000,
						"mfaInitSkipLifetime": 10000000,
						"secondFactorCheckLifetime": 10000000,
						"multiFactorCheckLifetime": 10000000,
						"enableRegistrationCaptcha": false,
						"enableLoginCaptcha": false,
						"captchaType": 0,
						"captchaSiteKey": "",
						"captchaSecretKey": ""
			}`),
					), instance.LoginPolicyAddedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.login_policies5 (aggregate_id, instance_id, creation_date, change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, force_mfa_local_only, passwordless_type, is_default, hide_password_reset, ignore_unknown_usernames, allow_domain_discovery, disable_login_with_email, disable_login_with_phone, default_redirect_uri, password_check_lifetime, external_login_check_lifetime, mfa_init_skip_lifetime, second_factor_check_lifetime, multi_factor_check_lifetime, enable_registration_captcha, enable_login_captcha, captcha_type, captcha_site_key, captcha_secret_key) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								true,
								true,
								true,
								true,
								true,
								true,
								"https://example.com/redirect",
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								false,
								false,
								domain.CaptchaTypeDisabled,
								"",
								"",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceLoginPolicyChanged",
			reduce: (&loginPolicyProjection{}).reduceLoginPolicyChanged,
			args: args{
				event: getEvent(
					testEvent(
						instance.LoginPolicyChangedEventType,
						instance.AggregateType,
						[]byte(`{
			"allowUsernamePassword": true,
			"allowRegister": true,
			"allowExternalIdp": true,
			"forceMFA": true,
			"forceMFALocalOnly": true,
			"hidePasswordReset": true,
			"ignoreUnknownUsernames": true,
			"allowDomainDiscovery": true,
			"disableLoginWithEmail": true,
			"disableLoginWithPhone": true,
			"passwordlessType": 1,
			"defaultRedirectURI": "https://example.com/redirect",
			"enableRegistrationCaptcha": false,
			"enableLoginCaptcha": false,
			"captchaType": 0,
			"captchaSiteKey": "",
			"captchaSecretKey": ""
			}`),
					), instance.LoginPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies5 SET (change_date, sequence, allow_register, allow_username_password, allow_external_idps, force_mfa, force_mfa_local_only, passwordless_type, hide_password_reset, ignore_unknown_usernames, allow_domain_discovery, disable_login_with_email, disable_login_with_phone, default_redirect_uri, enable_registration_captcha, enable_login_captcha, captcha_type, captcha_site_key, captcha_secret_key) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) WHERE (aggregate_id = $20) AND (instance_id = $21)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								true,
								true,
								true,
								true,
								true,
								"https://example.com/redirect",
								false,
								false,
								domain.CaptchaTypeDisabled,
								"",
								"",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceMFAAdded",
			reduce: (&loginPolicyProjection{}).reduceMFAAdded,
			args: args{
				event: getEvent(
					testEvent(
						instance.LoginPolicyMultiFactorAddedEventType,
						instance.AggregateType,
						[]byte(`{
		"mfaType": 1
		}`),
					), instance.MultiFactorAddedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies5 SET (change_date, sequence, multi_factors) = ($1, $2, array_append(multi_factors, $3)) WHERE (aggregate_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.MultiFactorTypeU2FWithPIN,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceMFARemoved",
			reduce: (&loginPolicyProjection{}).reduceMFARemoved,
			args: args{
				event: getEvent(
					testEvent(
						instance.LoginPolicyMultiFactorRemovedEventType,
						instance.AggregateType,
						[]byte(`{
			"mfaType": 1
			}`),
					), instance.MultiFactorRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies5 SET (change_date, sequence, multi_factors) = ($1, $2, array_remove(multi_factors, $3)) WHERE (aggregate_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.MultiFactorTypeU2FWithPIN,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceSecondFactorAdded u2f",
			reduce: (&loginPolicyProjection{}).reduceSecondFactorAdded,
			args: args{
				event: getEvent(
					testEvent(
						instance.LoginPolicySecondFactorAddedEventType,
						instance.AggregateType,
						[]byte(`{
			"mfaType": 2
			}`),
					), instance.SecondFactorAddedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies5 SET (change_date, sequence, second_factors) = ($1, $2, array_append(second_factors, $3)) WHERE (aggregate_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SecondFactorTypeU2F,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceSecondFactorRemoved  u2f",
			reduce: (&loginPolicyProjection{}).reduceSecondFactorRemoved,
			args: args{
				event: getEvent(
					testEvent(
						instance.LoginPolicySecondFactorRemovedEventType,
						instance.AggregateType,
						[]byte(`{
			"mfaType": 2
			}`),
					), instance.SecondFactorRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies5 SET (change_date, sequence, second_factors) = ($1, $2, array_remove(second_factors, $3)) WHERE (aggregate_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SecondFactorTypeU2F,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceSecondFactorAdded otp email",
			reduce: (&loginPolicyProjection{}).reduceSecondFactorAdded,
			args: args{
				event: getEvent(testEvent(
					instance.LoginPolicySecondFactorAddedEventType,
					instance.AggregateType,
					[]byte(`{
			"mfaType": 3
			}`),
				), instance.SecondFactorAddedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies5 SET (change_date, sequence, second_factors) = ($1, $2, array_append(second_factors, $3)) WHERE (aggregate_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SecondFactorTypeOTPEmail,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceSecondFactorRemoved otp email",
			reduce: (&loginPolicyProjection{}).reduceSecondFactorRemoved,
			args: args{
				event: getEvent(testEvent(
					instance.LoginPolicySecondFactorRemovedEventType,
					instance.AggregateType,
					[]byte(`{
			"mfaType": 3
			}`),
				), instance.SecondFactorRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.login_policies5 SET (change_date, sequence, second_factors) = ($1, $2, array_remove(second_factors, $3)) WHERE (aggregate_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SecondFactorTypeOTPEmail,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceOwnerRemoved",
			reduce: (&loginPolicyProjection{}).reduceOwnerRemoved,
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
							expectedStmt: "DELETE FROM projections.login_policies5 WHERE (instance_id = $1) AND (aggregate_id = $2)",
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
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(LoginPolicyInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.login_policies5 WHERE (instance_id = $1)",
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
			assertReduce(t, got, err, LoginPolicyTable, tt.want)
		})
	}
}
