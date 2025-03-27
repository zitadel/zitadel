package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	loginPolicyQuery = `SELECT projections.login_policies5.aggregate_id,` +
		` projections.login_policies5.creation_date,` +
		` projections.login_policies5.change_date,` +
		` projections.login_policies5.sequence,` +
		` projections.login_policies5.allow_register,` +
		` projections.login_policies5.allow_username_password,` +
		` projections.login_policies5.allow_external_idps,` +
		` projections.login_policies5.force_mfa,` +
		` projections.login_policies5.force_mfa_local_only,` +
		` projections.login_policies5.second_factors,` +
		` projections.login_policies5.multi_factors,` +
		` projections.login_policies5.passwordless_type,` +
		` projections.login_policies5.is_default,` +
		` projections.login_policies5.hide_password_reset,` +
		` projections.login_policies5.ignore_unknown_usernames,` +
		` projections.login_policies5.allow_domain_discovery,` +
		` projections.login_policies5.disable_login_with_email,` +
		` projections.login_policies5.disable_login_with_phone,` +
		` projections.login_policies5.default_redirect_uri,` +
		` projections.login_policies5.password_check_lifetime,` +
		` projections.login_policies5.external_login_check_lifetime,` +
		` projections.login_policies5.mfa_init_skip_lifetime,` +
		` projections.login_policies5.second_factor_check_lifetime,` +
		` projections.login_policies5.multi_factor_check_lifetime` +
		` FROM projections.login_policies5`
	loginPolicyCols = []string{
		"aggregate_id",
		"creation_date",
		"change_date",
		"sequence",
		"allow_register",
		"allow_username_password",
		"allow_external_idps",
		"force_mfa",
		"force_mfa_local_only",
		"second_factors",
		"multi_factors",
		"passwordless_type",
		"is_default",
		"hide_password_reset",
		"ignore_unknown_usernames",
		"allow_domain_discovery",
		"disable_login_with_email",
		"disable_login_with_phone",
		"default_redirect_uri",
		"password_check_lifetime",
		"external_login_check_lifetime",
		"mfa_init_skip_lifetime",
		"second_factor_check_lifetime",
		"multi_factor_check_lifetime",
	}

	prepareLoginPolicy2FAsStmt = `SELECT projections.login_policies5.second_factors` +
		` FROM projections.login_policies5`
	prepareLoginPolicy2FAsCols = []string{
		"second_factors",
	}

	prepareLoginPolicyMFAsStmt = `SELECT projections.login_policies5.multi_factors` +
		` FROM projections.login_policies5`
	prepareLoginPolicyMFAsCols = []string{
		"multi_factors",
	}
)

func Test_LoginPolicyPrepares(t *testing.T) {
	duration := 2 * time.Hour
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
			name:    "prepareLoginPolicyQuery no result",
			prepare: prepareLoginPolicyQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					regexp.QuoteMeta(loginPolicyQuery),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*LoginPolicy)(nil),
		},
		{
			name:    "prepareLoginPolicyQuery found",
			prepare: prepareLoginPolicyQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(loginPolicyQuery),
					loginPolicyCols,
					[]driver.Value{
						"ro",
						testNow,
						testNow,
						uint64(20211109),
						true,
						true,
						true,
						true,
						true,
						database.NumberArray[domain.SecondFactorType]{domain.SecondFactorTypeTOTP},
						database.NumberArray[domain.MultiFactorType]{domain.MultiFactorTypeU2FWithPIN},
						domain.PasswordlessTypeAllowed,
						true,
						true,
						true,
						true,
						true,
						true,
						"https://example.com/redirect",
						&duration,
						&duration,
						&duration,
						&duration,
						&duration,
					},
				),
			},
			object: &LoginPolicy{
				OrgID:                      "ro",
				CreationDate:               testNow,
				ChangeDate:                 testNow,
				Sequence:                   20211109,
				AllowRegister:              true,
				AllowUsernamePassword:      true,
				AllowExternalIDPs:          true,
				ForceMFA:                   true,
				ForceMFALocalOnly:          true,
				SecondFactors:              database.NumberArray[domain.SecondFactorType]{domain.SecondFactorTypeTOTP},
				MultiFactors:               database.NumberArray[domain.MultiFactorType]{domain.MultiFactorTypeU2FWithPIN},
				PasswordlessType:           domain.PasswordlessTypeAllowed,
				IsDefault:                  true,
				HidePasswordReset:          true,
				IgnoreUnknownUsernames:     true,
				AllowDomainDiscovery:       true,
				DisableLoginWithEmail:      true,
				DisableLoginWithPhone:      true,
				DefaultRedirectURI:         "https://example.com/redirect",
				PasswordCheckLifetime:      database.Duration(duration),
				ExternalLoginCheckLifetime: database.Duration(duration),
				MFAInitSkipLifetime:        database.Duration(duration),
				SecondFactorCheckLifetime:  database.Duration(duration),
				MultiFactorCheckLifetime:   database.Duration(duration),
			},
		},
		{
			name:    "prepareLoginPolicyQuery sql err",
			prepare: prepareLoginPolicyQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(loginPolicyQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*LoginPolicy)(nil),
		},
		{
			name:    "prepareLoginPolicy2FAsQuery no result",
			prepare: prepareLoginPolicy2FAsQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					regexp.QuoteMeta(prepareLoginPolicy2FAsStmt),
					prepareLoginPolicy2FAsCols,
					nil,
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*SecondFactors)(nil),
		},
		{
			name:    "prepareLoginPolicy2FAsQuery found",
			prepare: prepareLoginPolicy2FAsQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareLoginPolicy2FAsStmt),
					prepareLoginPolicy2FAsCols,
					[]driver.Value{
						database.NumberArray[domain.SecondFactorType]{domain.SecondFactorTypeTOTP},
					},
				),
			},
			object: &SecondFactors{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Factors: database.NumberArray[domain.SecondFactorType]{domain.SecondFactorTypeTOTP},
			},
		},
		{
			name:    "prepareLoginPolicy2FAsQuery found no factors",
			prepare: prepareLoginPolicy2FAsQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareLoginPolicy2FAsStmt),
					prepareLoginPolicy2FAsCols,
					[]driver.Value{
						database.NumberArray[domain.SecondFactorType]{},
					},
				),
			},
			object: &SecondFactors{Factors: database.NumberArray[domain.SecondFactorType]{}},
		},
		{
			name:    "prepareLoginPolicy2FAsQuery sql err",
			prepare: prepareLoginPolicy2FAsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareLoginPolicy2FAsStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*SecondFactors)(nil),
		},
		{
			name:    "prepareLoginPolicyMFAsQuery no result",
			prepare: prepareLoginPolicyMFAsQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					regexp.QuoteMeta(prepareLoginPolicyMFAsStmt),
					prepareLoginPolicyMFAsCols,
					nil,
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*MultiFactors)(nil),
		},
		{
			name:    "prepareLoginPolicyMFAsQuery found",
			prepare: prepareLoginPolicyMFAsQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareLoginPolicyMFAsStmt),
					prepareLoginPolicyMFAsCols,
					[]driver.Value{
						database.NumberArray[domain.MultiFactorType]{domain.MultiFactorTypeU2FWithPIN},
					},
				),
			},
			object: &MultiFactors{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Factors: database.NumberArray[domain.MultiFactorType]{domain.MultiFactorTypeU2FWithPIN},
			},
		},
		{
			name:    "prepareLoginPolicyMFAsQuery found no factors",
			prepare: prepareLoginPolicyMFAsQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareLoginPolicyMFAsStmt),
					prepareLoginPolicyMFAsCols,
					[]driver.Value{
						database.NumberArray[domain.MultiFactorType]{},
					},
				),
			},
			object: &MultiFactors{Factors: database.NumberArray[domain.MultiFactorType]{}},
		},
		{
			name:    "prepareLoginPolicyMFAsQuery sql err",
			prepare: prepareLoginPolicyMFAsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareLoginPolicyMFAsStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*MultiFactors)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
