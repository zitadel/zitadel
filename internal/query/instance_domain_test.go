package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
)

func Test_InstanceDomainPrepares(t *testing.T) {
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
			prepare: prepareInstanceDomainsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT projections.instance_domains.creation_date,`+
						` projections.instance_domains.change_date,`+
						` projections.instance_domains.sequence,`+
						` projections.instance_domains.domain,`+
						` projections.instance_domains.instance_id,`+
						` projections.instance_domains.is_generated,`+
						` COUNT(*) OVER ()`+
						` FROM projections.instance_domains`),
					nil,
					nil,
				),
			},
			object: &InstanceDomains{Domains: []*InstanceDomain{}},
		},
		{
			name:    "prepareDomainsQuery one result",
			prepare: prepareInstanceDomainsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT projections.instance_domains.creation_date,`+
						` projections.instance_domains.change_date,`+
						` projections.instance_domains.sequence,`+
						` projections.instance_domains.domain,`+
						` projections.instance_domains.instance_id,`+
						` projections.instance_domains.is_generated,`+
						` COUNT(*) OVER ()`+
						` FROM projections.instance_domains`),
					[]string{
						"creation_date",
						"change_date",
						"sequence",
						"domain",
						"instance_id",
						"is_generated",
						"count",
					},
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211109),
							"zitadel.ch",
							"inst-id",
							true,
						},
					},
				),
			},
			object: &InstanceDomains{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Domains: []*InstanceDomain{
					{
						CreationDate: testNow,
						ChangeDate:   testNow,
						Sequence:     20211109,
						Domain:       "zitadel.ch",
						InstanceID:   "inst-id",
						IsGenerated:  true,
					},
				},
			},
		},
		{
			name:    "prepareDomainsQuery multiple result",
			prepare: prepareInstanceDomainsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT projections.instance_domains.creation_date,`+
						` projections.instance_domains.change_date,`+
						` projections.instance_domains.sequence,`+
						` projections.instance_domains.domain,`+
						` projections.instance_domains.instance_id,`+
						` projections.instance_domains.is_generated,`+
						` COUNT(*) OVER ()`+
						` FROM projections.instance_domains`),
					[]string{
						"creation_date",
						"change_date",
						"sequence",
						"domain",
						"instance_id",
						"is_generated",
						"count",
					},
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211109),
							"zitadel.ch",
							"inst-id",
							true,
						},
						{
							testNow,
							testNow,
							uint64(20211109),
							"zitadel.com",
							"inst-id",
							false,
						},
					},
				),
			},
			object: &InstanceDomains{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Domains: []*InstanceDomain{
					{
						CreationDate: testNow,
						ChangeDate:   testNow,
						Sequence:     20211109,
						Domain:       "zitadel.ch",
						InstanceID:   "inst-id",
						IsGenerated:  true,
					},
					{
						CreationDate: testNow,
						ChangeDate:   testNow,
						Sequence:     20211109,
						Domain:       "zitadel.com",
						InstanceID:   "inst-id",
						IsGenerated:  false,
					},
				},
			},
		},
		{
			name:    "prepareDomainsQuery sql err",
			prepare: prepareInstanceDomainsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT projections.instance_domains.creation_date,`+
						` projections.instance_domains.change_date,`+
						` projections.instance_domains.sequence,`+
						` projections.instance_domains.domain,`+
						` projections.instance_domains.instance_id,`+
						` projections.instance_domains.is_generated,`+
						` COUNT(*) OVER ()`+
						` FROM projections.instance_domains`),
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
