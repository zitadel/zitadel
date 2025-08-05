package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	preparePrivacyPolicyStmt = `SELECT projections.privacy_policies4.id,` +
		` projections.privacy_policies4.sequence,` +
		` projections.privacy_policies4.creation_date,` +
		` projections.privacy_policies4.change_date,` +
		` projections.privacy_policies4.resource_owner,` +
		` projections.privacy_policies4.privacy_link,` +
		` projections.privacy_policies4.tos_link,` +
		` projections.privacy_policies4.help_link,` +
		` projections.privacy_policies4.support_email,` +
		` projections.privacy_policies4.docs_link,` +
		` projections.privacy_policies4.custom_link,` +
		` projections.privacy_policies4.custom_link_text,` +
		` projections.privacy_policies4.is_default,` +
		` projections.privacy_policies4.state` +
		` FROM projections.privacy_policies4`
	preparePrivacyPolicyCols = []string{
		"id",
		"sequence",
		"creation_date",
		"change_date",
		"resource_owner",
		"privacy_link",
		"tos_link",
		"help_link",
		"support_email",
		"docs_link",
		"custom_link",
		"custom_link_text",
		"is_default",
		"state",
	}
)

func Test_PrivacyPolicyPrepares(t *testing.T) {
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
			name:    "preparePrivacyPolicyQuery no result",
			prepare: preparePrivacyPolicyQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					regexp.QuoteMeta(preparePrivacyPolicyStmt),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsNotFound(err) {
						return fmt.Errorf("err should be NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*PrivacyPolicy)(nil),
		},
		{
			name:    "preparePrivacyPolicyQuery found",
			prepare: preparePrivacyPolicyQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(preparePrivacyPolicyStmt),
					preparePrivacyPolicyCols,
					[]driver.Value{
						"pol-id",
						uint64(20211109),
						testNow,
						testNow,
						"ro",
						"privacy.ch",
						"tos.ch",
						"help.ch",
						"support@example.com",
						"zitadel.com/docs",
						"zitadel.com",
						"Zitadel",
						true,
						domain.PolicyStateActive,
					},
				),
			},
			object: &PrivacyPolicy{
				ID:             "pol-id",
				CreationDate:   testNow,
				ChangeDate:     testNow,
				Sequence:       20211109,
				ResourceOwner:  "ro",
				State:          domain.PolicyStateActive,
				PrivacyLink:    "privacy.ch",
				TOSLink:        "tos.ch",
				HelpLink:       "help.ch",
				SupportEmail:   "support@example.com",
				DocsLink:       "zitadel.com/docs",
				CustomLink:     "zitadel.com",
				CustomLinkText: "Zitadel",
				IsDefault:      true,
			},
		},
		{
			name:    "preparePrivacyPolicyQuery sql err",
			prepare: preparePrivacyPolicyQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(preparePrivacyPolicyStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*PrivacyPolicy)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
