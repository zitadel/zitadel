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
	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	loginPolicyQuery = `SELECT projections.login_policies4.aggregate_id,` +
		` projections.login_policies4.creation_date,` +
		` projections.login_policies4.change_date,` +
		` projections.login_policies4.sequence,` +
		` projections.login_policies4.allow_register,` +
		` projections.login_policies4.allow_username_password,` +
		` projections.login_policies4.allow_external_idps,` +
		` projections.login_policies4.force_mfa,` +
		` projections.login_policies4.second_factors,` +
		` projections.login_policies4.multi_factors,` +
		` projections.login_policies4.passwordless_type,` +
		` projections.login_policies4.is_default,` +
		` projections.login_policies4.hide_password_reset,` +
		` projections.login_policies4.ignore_unknown_usernames,` +
		` projections.login_policies4.allow_domain_discovery,` +
		` projections.login_policies4.disable_login_with_email,` +
		` projections.login_policies4.disable_login_with_phone,` +
		` projections.login_policies4.default_redirect_uri,` +
		` projections.login_policies4.password_check_lifetime,` +
		` projections.login_policies4.external_login_check_lifetime,` +
		` projections.login_policies4.mfa_init_skip_lifetime,` +
		` projections.login_policies4.second_factor_check_lifetime,` +
		` projections.login_policies4.multi_factor_check_lifetime` +
		` FROM projections.login_policies4` +
		` AS OF SYSTEM TIME '-1 ms'`
	loginPolicyCols = []string{
		"aggregate_id",
		"creation_date",
		"change_date",
		"sequence",
		"allow_register",
		"allow_username_password",
		"allow_external_idps",
		"force_mfa",
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

	prepareLoginPolicy2FAsStmt = `SELECT projections.login_policies4.second_factors` +
		` FROM projections.login_policies4` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareLoginPolicy2FAsCols = []string{
		"second_factors",
	}

	prepareLoginPolicyMFAsStmt = `SELECT projections.login_policies4.multi_factors` +
		` FROM projections.login_policies4` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareLoginPolicyMFAsCols = []string{
		"multi_factors",
	}
)

func Test_LoginPolicyPrepares(t *testing.T) {
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
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(loginPolicyQuery),
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
						database.EnumArray[domain.SecondFactorType]{domain.SecondFactorTypeOTP},
						database.EnumArray[domain.MultiFactorType]{domain.MultiFactorTypeU2FWithPIN},
						domain.PasswordlessTypeAllowed,
						true,
						true,
						true,
						true,
						true,
						true,
						"https://example.com/redirect",
						time.Hour * 2,
						time.Hour * 2,
						time.Hour * 2,
						time.Hour * 2,
						time.Hour * 2,
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
				SecondFactors:              database.EnumArray[domain.SecondFactorType]{domain.SecondFactorTypeOTP},
				MultiFactors:               database.EnumArray[domain.MultiFactorType]{domain.MultiFactorTypeU2FWithPIN},
				PasswordlessType:           domain.PasswordlessTypeAllowed,
				IsDefault:                  true,
				HidePasswordReset:          true,
				IgnoreUnknownUsernames:     true,
				AllowDomainDiscovery:       true,
				DisableLoginWithEmail:      true,
				DisableLoginWithPhone:      true,
				DefaultRedirectURI:         "https://example.com/redirect",
				PasswordCheckLifetime:      time.Hour * 2,
				ExternalLoginCheckLifetime: time.Hour * 2,
				MFAInitSkipLifetime:        time.Hour * 2,
				SecondFactorCheckLifetime:  time.Hour * 2,
				MultiFactorCheckLifetime:   time.Hour * 2,
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
			object: nil,
		},
		{
			name:    "prepareLoginPolicy2FAsQuery no result",
			prepare: prepareLoginPolicy2FAsQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareLoginPolicy2FAsStmt),
					prepareLoginPolicy2FAsCols,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsNotFound(err) {
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
						database.EnumArray[domain.SecondFactorType]{domain.SecondFactorTypeOTP},
					},
				),
			},
			object: &SecondFactors{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Factors: database.EnumArray[domain.SecondFactorType]{domain.SecondFactorTypeOTP},
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
						database.EnumArray[domain.SecondFactorType]{},
					},
				),
			},
			object: &SecondFactors{Factors: database.EnumArray[domain.SecondFactorType]{}},
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
			object: nil,
		},
		{
			name:    "prepareLoginPolicyMFAsQuery no result",
			prepare: prepareLoginPolicyMFAsQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareLoginPolicyMFAsStmt),
					prepareLoginPolicyMFAsCols,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsNotFound(err) {
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
						database.EnumArray[domain.MultiFactorType]{domain.MultiFactorTypeU2FWithPIN},
					},
				),
			},
			object: &MultiFactors{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Factors: database.EnumArray[domain.MultiFactorType]{domain.MultiFactorTypeU2FWithPIN},
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
						database.EnumArray[domain.MultiFactorType]{},
					},
				),
			},
			object: &MultiFactors{Factors: database.EnumArray[domain.MultiFactorType]{}},
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
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
