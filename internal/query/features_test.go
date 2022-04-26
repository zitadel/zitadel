package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	errs "github.com/zitadel/zitadel/internal/errors"
)

func Test_FeaturesPrepares(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareFeaturesQuery no result",
			prepare: prepareFeaturesQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(`SELECT zitadel.projections.features.aggregate_id,`+
						` zitadel.projections.features.change_date,`+
						` zitadel.projections.features.sequence,`+
						` zitadel.projections.features.is_default,`+
						` zitadel.projections.features.tier_name,`+
						` zitadel.projections.features.tier_description,`+
						` zitadel.projections.features.state,`+
						` zitadel.projections.features.state_description,`+
						` zitadel.projections.features.audit_log_retention,`+
						` zitadel.projections.features.login_policy_factors,`+
						` zitadel.projections.features.login_policy_idp,`+
						` zitadel.projections.features.login_policy_passwordless,`+
						` zitadel.projections.features.login_policy_registration,`+
						` zitadel.projections.features.login_policy_username_login,`+
						` zitadel.projections.features.login_policy_password_reset,`+
						` zitadel.projections.features.password_complexity_policy,`+
						` zitadel.projections.features.label_policy_private_label,`+
						` zitadel.projections.features.label_policy_watermark,`+
						` zitadel.projections.features.custom_domain,`+
						` zitadel.projections.features.privacy_policy,`+
						` zitadel.projections.features.metadata_user,`+
						` zitadel.projections.features.custom_text_message,`+
						` zitadel.projections.features.custom_text_login,`+
						` zitadel.projections.features.lockout_policy,`+
						` zitadel.projections.features.actions_allowed,`+
						` zitadel.projections.features.max_actions`+
						` FROM zitadel.projections.features`),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Features)(nil),
		},
		{
			name:    "prepareFeaturesQuery found",
			prepare: prepareFeaturesQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(`SELECT zitadel.projections.features.aggregate_id,`+
						` zitadel.projections.features.change_date,`+
						` zitadel.projections.features.sequence,`+
						` zitadel.projections.features.is_default,`+
						` zitadel.projections.features.tier_name,`+
						` zitadel.projections.features.tier_description,`+
						` zitadel.projections.features.state,`+
						` zitadel.projections.features.state_description,`+
						` zitadel.projections.features.audit_log_retention,`+
						` zitadel.projections.features.login_policy_factors,`+
						` zitadel.projections.features.login_policy_idp,`+
						` zitadel.projections.features.login_policy_passwordless,`+
						` zitadel.projections.features.login_policy_registration,`+
						` zitadel.projections.features.login_policy_username_login,`+
						` zitadel.projections.features.login_policy_password_reset,`+
						` zitadel.projections.features.password_complexity_policy,`+
						` zitadel.projections.features.label_policy_private_label,`+
						` zitadel.projections.features.label_policy_watermark,`+
						` zitadel.projections.features.custom_domain,`+
						` zitadel.projections.features.privacy_policy,`+
						` zitadel.projections.features.metadata_user,`+
						` zitadel.projections.features.custom_text_message,`+
						` zitadel.projections.features.custom_text_login,`+
						` zitadel.projections.features.lockout_policy,`+
						` zitadel.projections.features.actions_allowed,`+
						` zitadel.projections.features.max_actions`+
						` FROM zitadel.projections.features`),
					[]string{
						"aggregate_id",
						"change_date",
						"sequence",
						"is_default",
						"tier_name",
						"tier_description",
						"state",
						"state_description",
						"audit_log_retention",
						"login_policy_factors",
						"login_policy_idp",
						"login_policy_passwordless",
						"login_policy_registration",
						"login_policy_username_login",
						"login_policy_password_reset",
						"password_complexity_policy",
						"label_policy_private_label",
						"label_policy_watermark",
						"custom_domain",
						"privacy_policy",
						"metadata_user",
						"custom_text_message",
						"custom_text_login",
						"lockout_policy",
						"actions_allowed",
						"max_actions",
					},
					[]driver.Value{
						"aggregate-id",
						testNow,
						uint64(20211115),
						true,
						"tier-name",
						"tier-description",
						1,
						"state-description",
						uint(604800000000000), // 7days in nanoseconds
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						domain.ActionsMaxAllowed,
						10,
					},
				),
			},
			object: &Features{
				AggregateID:              "aggregate-id",
				ChangeDate:               testNow,
				Sequence:                 20211115,
				IsDefault:                true,
				TierName:                 "tier-name",
				TierDescription:          "tier-description",
				State:                    domain.FeaturesStateActive,
				StateDescription:         "state-description",
				AuditLogRetention:        7 * 24 * time.Hour,
				LoginPolicyFactors:       true,
				LoginPolicyIDP:           true,
				LoginPolicyPasswordless:  true,
				LoginPolicyRegistration:  true,
				LoginPolicyUsernameLogin: true,
				LoginPolicyPasswordReset: true,
				PasswordComplexityPolicy: true,
				LabelPolicyPrivateLabel:  true,
				LabelPolicyWatermark:     true,
				CustomDomain:             true,
				PrivacyPolicy:            true,
				MetadataUser:             true,
				CustomTextMessage:        true,
				CustomTextLogin:          true,
				LockoutPolicy:            true,
				ActionsAllowed:           domain.ActionsMaxAllowed,
				MaxActions:               10,
			},
		},
		{
			name:    "prepareFeaturesQuery found with empty",
			prepare: prepareFeaturesQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(`SELECT zitadel.projections.features.aggregate_id,`+
						` zitadel.projections.features.change_date,`+
						` zitadel.projections.features.sequence,`+
						` zitadel.projections.features.is_default,`+
						` zitadel.projections.features.tier_name,`+
						` zitadel.projections.features.tier_description,`+
						` zitadel.projections.features.state,`+
						` zitadel.projections.features.state_description,`+
						` zitadel.projections.features.audit_log_retention,`+
						` zitadel.projections.features.login_policy_factors,`+
						` zitadel.projections.features.login_policy_idp,`+
						` zitadel.projections.features.login_policy_passwordless,`+
						` zitadel.projections.features.login_policy_registration,`+
						` zitadel.projections.features.login_policy_username_login,`+
						` zitadel.projections.features.login_policy_password_reset,`+
						` zitadel.projections.features.password_complexity_policy,`+
						` zitadel.projections.features.label_policy_private_label,`+
						` zitadel.projections.features.label_policy_watermark,`+
						` zitadel.projections.features.custom_domain,`+
						` zitadel.projections.features.privacy_policy,`+
						` zitadel.projections.features.metadata_user,`+
						` zitadel.projections.features.custom_text_message,`+
						` zitadel.projections.features.custom_text_login,`+
						` zitadel.projections.features.lockout_policy,`+
						` zitadel.projections.features.actions_allowed,`+
						` zitadel.projections.features.max_actions`+
						` FROM zitadel.projections.features`),
					[]string{
						"aggregate_id",
						"change_date",
						"sequence",
						"is_default",
						"tier_name",
						"tier_description",
						"state",
						"state_description",
						"audit_log_retention",
						"login_policy_factors",
						"login_policy_idp",
						"login_policy_passwordless",
						"login_policy_registration",
						"login_policy_username_login",
						"login_policy_password_reset",
						"password_complexity_policy",
						"label_policy_private_label",
						"label_policy_watermark",
						"custom_domain",
						"privacy_policy",
						"metadata_user",
						"custom_text_message",
						"custom_text_login",
						"lockout_policy",
						"actions_allowed",
						"max_actions",
					},
					[]driver.Value{
						"aggregate-id",
						testNow,
						uint64(20211115),
						true,
						nil,
						nil,
						1,
						nil,
						uint(604800000000000), // 7days in nanoseconds
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						domain.ActionsMaxAllowed,
						10,
					},
				),
			},
			object: &Features{
				AggregateID:              "aggregate-id",
				ChangeDate:               testNow,
				Sequence:                 20211115,
				IsDefault:                true,
				TierName:                 "",
				TierDescription:          "",
				State:                    domain.FeaturesStateActive,
				StateDescription:         "",
				AuditLogRetention:        7 * 24 * time.Hour,
				LoginPolicyFactors:       true,
				LoginPolicyIDP:           true,
				LoginPolicyPasswordless:  true,
				LoginPolicyRegistration:  true,
				LoginPolicyUsernameLogin: true,
				LoginPolicyPasswordReset: true,
				PasswordComplexityPolicy: true,
				LabelPolicyPrivateLabel:  true,
				LabelPolicyWatermark:     true,
				CustomDomain:             true,
				PrivacyPolicy:            true,
				MetadataUser:             true,
				CustomTextMessage:        true,
				CustomTextLogin:          true,
				LockoutPolicy:            true,
				ActionsAllowed:           domain.ActionsMaxAllowed,
				MaxActions:               10,
			},
		},
		{
			name:    "prepareFeaturesQuery sql err",
			prepare: prepareFeaturesQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT zitadel.projections.features.aggregate_id,`+
						` zitadel.projections.features.change_date,`+
						` zitadel.projections.features.sequence,`+
						` zitadel.projections.features.is_default,`+
						` zitadel.projections.features.tier_name,`+
						` zitadel.projections.features.tier_description,`+
						` zitadel.projections.features.state,`+
						` zitadel.projections.features.state_description,`+
						` zitadel.projections.features.audit_log_retention,`+
						` zitadel.projections.features.login_policy_factors,`+
						` zitadel.projections.features.login_policy_idp,`+
						` zitadel.projections.features.login_policy_passwordless,`+
						` zitadel.projections.features.login_policy_registration,`+
						` zitadel.projections.features.login_policy_username_login,`+
						` zitadel.projections.features.login_policy_password_reset,`+
						` zitadel.projections.features.password_complexity_policy,`+
						` zitadel.projections.features.label_policy_private_label,`+
						` zitadel.projections.features.label_policy_watermark,`+
						` zitadel.projections.features.custom_domain,`+
						` zitadel.projections.features.privacy_policy,`+
						` zitadel.projections.features.metadata_user,`+
						` zitadel.projections.features.custom_text_message,`+
						` zitadel.projections.features.custom_text_login,`+
						` zitadel.projections.features.lockout_policy,`+
						` zitadel.projections.features.actions_allowed,`+
						` zitadel.projections.features.max_actions`+
						` FROM zitadel.projections.features`),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
