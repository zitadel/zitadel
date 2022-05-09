package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/lib/pq"

	"github.com/zitadel/zitadel/internal/domain"
	errs "github.com/zitadel/zitadel/internal/errors"
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
					regexp.QuoteMeta(`SELECT zitadel.projections.login_policies.aggregate_id,`+
						` zitadel.projections.login_policies.creation_date,`+
						` zitadel.projections.login_policies.change_date,`+
						` zitadel.projections.login_policies.sequence,`+
						` zitadel.projections.login_policies.allow_register,`+
						` zitadel.projections.login_policies.allow_username_password,`+
						` zitadel.projections.login_policies.allow_external_idps,`+
						` zitadel.projections.login_policies.force_mfa,`+
						` zitadel.projections.login_policies.second_factors,`+
						` zitadel.projections.login_policies.multi_factors,`+
						` zitadel.projections.login_policies.passwordless_type,`+
						` zitadel.projections.login_policies.is_default,`+
						` zitadel.projections.login_policies.hide_password_reset,`+
						` zitadel.projections.login_policies.ignore_unknown_usernames,`+
						` zitadel.projections.login_policies.default_redirect_uri`+
						` FROM zitadel.projections.login_policies`),
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
					regexp.QuoteMeta(`SELECT zitadel.projections.login_policies.aggregate_id,`+
						` zitadel.projections.login_policies.creation_date,`+
						` zitadel.projections.login_policies.change_date,`+
						` zitadel.projections.login_policies.sequence,`+
						` zitadel.projections.login_policies.allow_register,`+
						` zitadel.projections.login_policies.allow_username_password,`+
						` zitadel.projections.login_policies.allow_external_idps,`+
						` zitadel.projections.login_policies.force_mfa,`+
						` zitadel.projections.login_policies.second_factors,`+
						` zitadel.projections.login_policies.multi_factors,`+
						` zitadel.projections.login_policies.passwordless_type,`+
						` zitadel.projections.login_policies.is_default,`+
						` zitadel.projections.login_policies.hide_password_reset,`+
						` zitadel.projections.login_policies.ignore_unknown_usernames,`+
						` zitadel.projections.login_policies.default_redirect_uri`+
						` FROM zitadel.projections.login_policies`),
					[]string{
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
						"default_redirect_uri",
					},
					[]driver.Value{
						"ro",
						testNow,
						testNow,
						uint64(20211109),
						true,
						true,
						true,
						true,
						pq.Int32Array{int32(domain.SecondFactorTypeOTP)},
						pq.Int32Array{int32(domain.MultiFactorTypeU2FWithPIN)},
						domain.PasswordlessTypeAllowed,
						true,
						true,
						true,
						"https://example.com/redirect",
					},
				),
			},
			object: &LoginPolicy{
				OrgID:                  "ro",
				CreationDate:           testNow,
				ChangeDate:             testNow,
				Sequence:               20211109,
				AllowRegister:          true,
				AllowUsernamePassword:  true,
				AllowExternalIDPs:      true,
				ForceMFA:               true,
				SecondFactors:          []domain.SecondFactorType{domain.SecondFactorTypeOTP},
				MultiFactors:           []domain.MultiFactorType{domain.MultiFactorTypeU2FWithPIN},
				PasswordlessType:       domain.PasswordlessTypeAllowed,
				IsDefault:              true,
				HidePasswordReset:      true,
				IgnoreUnknownUsernames: true,
				DefaultRedirectURI:     "https://example.com/redirect",
			},
		},
		{
			name:    "prepareLoginPolicyQuery sql err",
			prepare: prepareLoginPolicyQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT zitadel.projections.login_policies.aggregate_id,`+
						` zitadel.projections.login_policies.creation_date,`+
						` zitadel.projections.login_policies.change_date,`+
						` zitadel.projections.login_policies.sequence,`+
						` zitadel.projections.login_policies.allow_register,`+
						` zitadel.projections.login_policies.allow_username_password,`+
						` zitadel.projections.login_policies.allow_external_idps,`+
						` zitadel.projections.login_policies.force_mfa,`+
						` zitadel.projections.login_policies.second_factors,`+
						` zitadel.projections.login_policies.multi_factors,`+
						` zitadel.projections.login_policies.passwordless_type,`+
						` zitadel.projections.login_policies.is_default,`+
						` zitadel.projections.login_policies.hide_password_reset,`+
						` zitadel.projections.login_policies.ignore_unknown_usernames,`+
						` zitadel.projections.login_policies.default_redirect_uri`+
						` FROM zitadel.projections.login_policies`),
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
					regexp.QuoteMeta(`SELECT zitadel.projections.login_policies.second_factors`+
						` FROM zitadel.projections.login_policies`),
					[]string{
						"second_factors",
					},
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
					regexp.QuoteMeta(`SELECT zitadel.projections.login_policies.second_factors`+
						` FROM zitadel.projections.login_policies`),
					[]string{
						"second_factors",
					},
					[]driver.Value{
						pq.Int32Array{int32(domain.SecondFactorTypeOTP)},
					},
				),
			},
			object: &SecondFactors{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Factors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
			},
		},
		{
			name:    "prepareLoginPolicy2FAsQuery found no factors",
			prepare: prepareLoginPolicy2FAsQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(`SELECT zitadel.projections.login_policies.second_factors`+
						` FROM zitadel.projections.login_policies`),
					[]string{
						"second_factors",
					},
					[]driver.Value{
						pq.Int32Array{},
					},
				),
			},
			object: &SecondFactors{Factors: []domain.SecondFactorType{}},
		},
		{
			name:    "prepareLoginPolicy2FAsQuery sql err",
			prepare: prepareLoginPolicy2FAsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT zitadel.projections.login_policies.second_factors`+
						` FROM zitadel.projections.login_policies`),
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
					regexp.QuoteMeta(`SELECT zitadel.projections.login_policies.multi_factors`+
						` FROM zitadel.projections.login_policies`),
					[]string{
						"multi_factors",
					},
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
					regexp.QuoteMeta(`SELECT zitadel.projections.login_policies.multi_factors`+
						` FROM zitadel.projections.login_policies`),
					[]string{
						"multi_factors",
					},
					[]driver.Value{
						pq.Int32Array{int32(domain.MultiFactorTypeU2FWithPIN)},
					},
				),
			},
			object: &MultiFactors{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Factors: []domain.MultiFactorType{domain.MultiFactorTypeU2FWithPIN},
			},
		},
		{
			name:    "prepareLoginPolicyMFAsQuery found no factors",
			prepare: prepareLoginPolicyMFAsQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(`SELECT zitadel.projections.login_policies.multi_factors`+
						` FROM zitadel.projections.login_policies`),
					[]string{
						"multi_factors",
					},
					[]driver.Value{
						pq.Int32Array{},
					},
				),
			},
			object: &MultiFactors{Factors: []domain.MultiFactorType{}},
		},
		{
			name:    "prepareLoginPolicyMFAsQuery sql err",
			prepare: prepareLoginPolicyMFAsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT zitadel.projections.login_policies.multi_factors`+
						` FROM zitadel.projections.login_policies`),
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
