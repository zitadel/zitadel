package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
)

var (
	prepareOrgDomainsStmt = `SELECT projections.org_domains2.creation_date,` +
		` projections.org_domains2.change_date,` +
		` projections.org_domains2.sequence,` +
		` projections.org_domains2.domain,` +
		` projections.org_domains2.org_id,` +
		` projections.org_domains2.is_verified,` +
		` projections.org_domains2.is_primary,` +
		` projections.org_domains2.validation_type,` +
		` COUNT(*) OVER ()` +
		` FROM projections.org_domains2` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareOrgDomainsCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"org_state",
		"sequence",
		"name",
		"primary_domain",
		"count",
	}
)

func Test_OrgDomainPrepares(t *testing.T) {
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
			name:    "prepareDomainsQuery no result",
			prepare: prepareDomainsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareOrgDomainsStmt),
					nil,
					nil,
				),
			},
			object: &Domains{Domains: []*Domain{}},
		},
		{
			name:    "prepareDomainsQuery one result",
			prepare: prepareDomainsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareOrgDomainsStmt),
					prepareOrgDomainsCols,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211109),
							"zitadel.ch",
							"ro",
							true,
							true,
							domain.OrgDomainValidationTypeHTTP,
						},
					},
				),
			},
			object: &Domains{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Domains: []*Domain{
					{
						CreationDate:   testNow,
						ChangeDate:     testNow,
						Sequence:       20211109,
						Domain:         "zitadel.ch",
						OrgID:          "ro",
						IsVerified:     true,
						IsPrimary:      true,
						ValidationType: domain.OrgDomainValidationTypeHTTP,
					},
				},
			},
		},
		{
			name:    "prepareDomainsQuery multiple result",
			prepare: prepareDomainsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareOrgDomainsStmt),
					prepareOrgDomainsCols,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211109),
							"zitadel.ch",
							"ro",
							true,
							true,
							domain.OrgDomainValidationTypeHTTP,
						},
						{
							testNow,
							testNow,
							uint64(20211109),
							"zitadel.ch",
							"ro",
							false,
							false,
							domain.OrgDomainValidationTypeDNS,
						},
					},
				),
			},
			object: &Domains{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Domains: []*Domain{
					{
						CreationDate:   testNow,
						ChangeDate:     testNow,
						Sequence:       20211109,
						Domain:         "zitadel.ch",
						OrgID:          "ro",
						IsVerified:     true,
						IsPrimary:      true,
						ValidationType: domain.OrgDomainValidationTypeHTTP,
					},
					{
						CreationDate:   testNow,
						ChangeDate:     testNow,
						Sequence:       20211109,
						Domain:         "zitadel.ch",
						OrgID:          "ro",
						IsVerified:     false,
						IsPrimary:      false,
						ValidationType: domain.OrgDomainValidationTypeDNS,
					},
				},
			},
		},
		{
			name:    "prepareDomainsQuery sql err",
			prepare: prepareDomainsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareOrgDomainsStmt),
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
