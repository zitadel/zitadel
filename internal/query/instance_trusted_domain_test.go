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
	prepareInstanceTrustedDomainsStmt = `SELECT projections.instance_trusted_domains.creation_date,` +
		` projections.instance_trusted_domains.change_date,` +
		` projections.instance_trusted_domains.sequence,` +
		` projections.instance_trusted_domains.domain,` +
		` projections.instance_trusted_domains.instance_id,` +
		` COUNT(*) OVER ()` +
		` FROM projections.instance_trusted_domains`
	prepareInstanceTrustedDomainsCols = []string{
		"creation_date",
		"change_date",
		"sequence",
		"domain",
		"instance_id",
		"count",
	}
)

func Test_InstanceTrustedDomainPrepares(t *testing.T) {
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
			name:    "prepareInstanceTrustedDomainsQuery no result",
			prepare: prepareInstanceTrustedDomainsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareInstanceTrustedDomainsStmt),
					nil,
					nil,
				),
			},
			object: &InstanceTrustedDomains{Domains: []*InstanceTrustedDomain{}},
		},
		{
			name:    "prepareInstanceTrustedDomainsQuery one result",
			prepare: prepareInstanceTrustedDomainsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareInstanceTrustedDomainsStmt),
					prepareInstanceTrustedDomainsCols,
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
			object: &InstanceTrustedDomains{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Domains: []*InstanceTrustedDomain{
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
			name:    "prepareInstanceTrustedDomainsQuery multiple result",
			prepare: prepareInstanceTrustedDomainsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareInstanceTrustedDomainsStmt),
					prepareInstanceTrustedDomainsCols,
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
			object: &InstanceTrustedDomains{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Domains: []*InstanceTrustedDomain{
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
			name:    "prepareInstanceTrustedDomainsQuery sql err",
			prepare: prepareInstanceTrustedDomainsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareInstanceTrustedDomainsStmt),
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
