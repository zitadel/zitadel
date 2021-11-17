package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/caos/zitadel/internal/domain"
	errs "github.com/caos/zitadel/internal/errors"
)

func Test_OrgIAMPolicyPrepares(t *testing.T) {
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
			name:    "prepareOrgIAMPolicyQuery no result",
			prepare: prepareOrgIAMPolicyQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT zitadel.projections.org_iam_policies.id,`+
						` zitadel.projections.org_iam_policies.sequence,`+
						` zitadel.projections.org_iam_policies.creation_date,`+
						` zitadel.projections.org_iam_policies.change_date,`+
						` zitadel.projections.org_iam_policies.resource_owner,`+
						` zitadel.projections.org_iam_policies.user_login_must_be_domain,`+
						` zitadel.projections.org_iam_policies.is_default,`+
						` zitadel.projections.org_iam_policies.state`+
						` FROM zitadel.projections.org_iam_policies`),
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
			object: (*OrgIAMPolicy)(nil),
		},
		{
			name:    "prepareOrgIAMPolicyQuery found",
			prepare: prepareOrgIAMPolicyQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(`SELECT zitadel.projections.org_iam_policies.id,`+
						` zitadel.projections.org_iam_policies.sequence,`+
						` zitadel.projections.org_iam_policies.creation_date,`+
						` zitadel.projections.org_iam_policies.change_date,`+
						` zitadel.projections.org_iam_policies.resource_owner,`+
						` zitadel.projections.org_iam_policies.user_login_must_be_domain,`+
						` zitadel.projections.org_iam_policies.is_default,`+
						` zitadel.projections.org_iam_policies.state`+
						` FROM zitadel.projections.org_iam_policies`),
					[]string{
						"id",
						"sequence",
						"creation_date",
						"change_date",
						"resource_owner",
						"user_login_must_be_domain",
						"is_default",
						"state",
					},
					[]driver.Value{
						"pol-id",
						uint64(20211109),
						testNow,
						testNow,
						"ro",
						true,
						true,
						domain.PolicyStateActive,
					},
				),
			},
			object: &OrgIAMPolicy{
				ID:                    "pol-id",
				CreationDate:          testNow,
				ChangeDate:            testNow,
				Sequence:              20211109,
				ResourceOwner:         "ro",
				State:                 domain.PolicyStateActive,
				UserLoginMustBeDomain: true,
				IsDefault:             true,
			},
		},
		{
			name:    "prepareOrgIAMPolicyQuery sql err",
			prepare: prepareOrgIAMPolicyQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT zitadel.projections.org_iam_policies.id,`+
						` zitadel.projections.org_iam_policies.sequence,`+
						` zitadel.projections.org_iam_policies.creation_date,`+
						` zitadel.projections.org_iam_policies.change_date,`+
						` zitadel.projections.org_iam_policies.resource_owner,`+
						` zitadel.projections.org_iam_policies.user_login_must_be_domain,`+
						` zitadel.projections.org_iam_policies.is_default,`+
						` zitadel.projections.org_iam_policies.state`+
						` FROM zitadel.projections.org_iam_policies`),
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
