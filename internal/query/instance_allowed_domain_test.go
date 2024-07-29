package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
)

var (
	prepareInstanceAllowedDomainsStmt = `SELECT projections.instance_allowed_domains.creation_date,` +
		` projections.instance_allowed_domains.change_date,` +
		` projections.instance_allowed_domains.sequence,` +
		` projections.instance_allowed_domains.domain,` +
		` projections.instance_allowed_domains.instance_id,` +
		` COUNT(*) OVER ()` +
		` FROM projections.instance_allowed_domains` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareInstanceAllowedDomainsCols = []string{
		"creation_date",
		"change_date",
		"sequence",
		"domain",
		"instance_id",
		"count",
	}
)

func Test_InstanceAllowedDomainPrepares(t *testing.T) {
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
			name:    "prepareInstanceAllowedDomainsQuery no result",
			prepare: prepareInstanceAllowedDomainsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareInstanceAllowedDomainsStmt),
					nil,
					nil,
				),
			},
			object: &InstanceAllowedDomains{Domains: []*InstanceAllowedDomain{}},
		},
		{
			name:    "prepareInstanceAllowedDomainsQuery one result",
			prepare: prepareInstanceAllowedDomainsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareInstanceAllowedDomainsStmt),
					prepareInstanceAllowedDomainsCols,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211109),
							"zitadel.ch",
							"inst-id",
						},
					},
				),
			},
			object: &InstanceAllowedDomains{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Domains: []*InstanceAllowedDomain{
					{
						CreationDate: testNow,
						ChangeDate:   testNow,
						Sequence:     20211109,
						Domain:       "zitadel.ch",
						InstanceID:   "inst-id",
					},
				},
			},
		},
		{
			name:    "prepareInstanceAllowedDomainsQuery multiple result",
			prepare: prepareInstanceAllowedDomainsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareInstanceAllowedDomainsStmt),
					prepareInstanceAllowedDomainsCols,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211109),
							"zitadel.ch",
							"inst-id",
						},
						{
							testNow,
							testNow,
							uint64(20211109),
							"zitadel.com",
							"inst-id",
						},
					},
				),
			},
			object: &InstanceAllowedDomains{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Domains: []*InstanceAllowedDomain{
					{
						CreationDate: testNow,
						ChangeDate:   testNow,
						Sequence:     20211109,
						Domain:       "zitadel.ch",
						InstanceID:   "inst-id",
					},
					{
						CreationDate: testNow,
						ChangeDate:   testNow,
						Sequence:     20211109,
						Domain:       "zitadel.com",
						InstanceID:   "inst-id",
					},
				},
			},
		},
		{
			name:    "prepareInstanceAllowedDomainsQuery sql err",
			prepare: prepareInstanceAllowedDomainsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareInstanceAllowedDomainsStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Domains)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
