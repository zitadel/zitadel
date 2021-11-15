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

func Test_OrgPrepares(t *testing.T) {
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
			name:    "prepareOrgsQuery no result",
			prepare: prepareOrgsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT zitadel.projections.orgs.id,`+
						` zitadel.projections.orgs.creation_date,`+
						` zitadel.projections.orgs.change_date,`+
						` zitadel.projections.orgs.resource_owner,`+
						` zitadel.projections.orgs.org_state,`+
						` zitadel.projections.orgs.sequence,`+
						` zitadel.projections.orgs.name,`+
						` zitadel.projections.orgs.primary_domain,`+
						` COUNT(*) OVER ()`+
						` FROM zitadel.projections.orgs`),
					nil,
					nil,
				),
			},
			object: &Orgs{Orgs: []*Org{}},
		},
		{
			name:    "prepareOrgsQuery one result",
			prepare: prepareOrgsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT zitadel.projections.orgs.id,`+
						` zitadel.projections.orgs.creation_date,`+
						` zitadel.projections.orgs.change_date,`+
						` zitadel.projections.orgs.resource_owner,`+
						` zitadel.projections.orgs.org_state,`+
						` zitadel.projections.orgs.sequence,`+
						` zitadel.projections.orgs.name,`+
						` zitadel.projections.orgs.primary_domain,`+
						` COUNT(*) OVER ()`+
						` FROM zitadel.projections.orgs`),
					[]string{
						"id",
						"creation_date",
						"change_date",
						"resource_owner",
						"org_state",
						"sequence",
						"name",
						"primary_domain",
						"count",
					},
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							"ro",
							domain.OrgStateActive,
							uint64(20211109),
							"org-name",
							"zitadel.ch",
						},
					},
				),
			},
			object: &Orgs{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Orgs: []*Org{
					{
						ID:            "id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.OrgStateActive,
						Sequence:      20211109,
						Name:          "org-name",
						Domain:        "zitadel.ch",
					},
				},
			},
		},
		{
			name:    "prepareOrgsQuery multiple result",
			prepare: prepareOrgsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT zitadel.projections.orgs.id,`+
						` zitadel.projections.orgs.creation_date,`+
						` zitadel.projections.orgs.change_date,`+
						` zitadel.projections.orgs.resource_owner,`+
						` zitadel.projections.orgs.org_state,`+
						` zitadel.projections.orgs.sequence,`+
						` zitadel.projections.orgs.name,`+
						` zitadel.projections.orgs.primary_domain,`+
						` COUNT(*) OVER ()`+
						` FROM zitadel.projections.orgs`),
					[]string{
						"id",
						"creation_date",
						"change_date",
						"resource_owner",
						"org_state",
						"sequence",
						"name",
						"primary_domain",
						"count",
					},
					[][]driver.Value{
						{
							"id-1",
							testNow,
							testNow,
							"ro",
							domain.OrgStateActive,
							uint64(20211108),
							"org-name-1",
							"zitadel.ch",
						},
						{
							"id-2",
							testNow,
							testNow,
							"ro",
							domain.OrgStateActive,
							uint64(20211108),
							"org-name-2",
							"caos.ch",
						},
					},
				),
			},
			object: &Orgs{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Orgs: []*Org{
					{
						ID:            "id-1",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.OrgStateActive,
						Sequence:      20211108,
						Name:          "org-name-1",
						Domain:        "zitadel.ch",
					},
					{
						ID:            "id-2",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.OrgStateActive,
						Sequence:      20211108,
						Name:          "org-name-2",
						Domain:        "caos.ch",
					},
				},
			},
		},
		{
			name:    "prepareOrgsQuery sql err",
			prepare: prepareOrgsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT zitadel.projections.orgs.id,`+
						` zitadel.projections.orgs.creation_date,`+
						` zitadel.projections.orgs.change_date,`+
						` zitadel.projections.orgs.resource_owner,`+
						` zitadel.projections.orgs.org_state,`+
						` zitadel.projections.orgs.sequence,`+
						` zitadel.projections.orgs.name,`+
						` zitadel.projections.orgs.primary_domain,`+
						` COUNT(*) OVER ()`+
						` FROM zitadel.projections.orgs`),
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
			name:    "prepareOrgQuery no result",
			prepare: prepareOrgQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT zitadel.projections.orgs.id,`+
						` zitadel.projections.orgs.creation_date,`+
						` zitadel.projections.orgs.change_date,`+
						` zitadel.projections.orgs.resource_owner,`+
						` zitadel.projections.orgs.org_state,`+
						` zitadel.projections.orgs.sequence,`+
						` zitadel.projections.orgs.name,`+
						` zitadel.projections.orgs.primary_domain`+
						` FROM zitadel.projections.orgs`),
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
			object: (*Org)(nil),
		},
		{
			name:    "prepareOrgQuery found",
			prepare: prepareOrgQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(`SELECT zitadel.projections.orgs.id,`+
						` zitadel.projections.orgs.creation_date,`+
						` zitadel.projections.orgs.change_date,`+
						` zitadel.projections.orgs.resource_owner,`+
						` zitadel.projections.orgs.org_state,`+
						` zitadel.projections.orgs.sequence,`+
						` zitadel.projections.orgs.name,`+
						` zitadel.projections.orgs.primary_domain`+
						` FROM zitadel.projections.orgs`),
					[]string{
						"id",
						"creation_date",
						"change_date",
						"resource_owner",
						"org_state",
						"sequence",
						"name",
						"primary_domain",
					},
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"ro",
						domain.OrgStateActive,
						uint64(20211108),
						"org-name",
						"zitadel.ch",
					},
				),
			},
			object: &Org{
				ID:            "id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				State:         domain.OrgStateActive,
				Sequence:      20211108,
				Name:          "org-name",
				Domain:        "zitadel.ch",
			},
		},
		{
			name:    "prepareOrgQuery sql err",
			prepare: prepareOrgQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT zitadel.projections.orgs.id,`+
						` zitadel.projections.orgs.creation_date,`+
						` zitadel.projections.orgs.change_date,`+
						` zitadel.projections.orgs.resource_owner,`+
						` zitadel.projections.orgs.org_state,`+
						` zitadel.projections.orgs.sequence,`+
						` zitadel.projections.orgs.name,`+
						` zitadel.projections.orgs.primary_domain`+
						` FROM zitadel.projections.orgs`),
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
			name:    "prepareOrgUniqueQuery no result",
			prepare: prepareOrgUniqueQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT COUNT(*) = 0`+
						` FROM zitadel.projections.orgs`),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsInternal(err) {
						return fmt.Errorf("err should be zitadel.Internal got: %w", err), false
					}
					return nil, true
				},
			},
			object: false,
		},
		{
			name:    "prepareOrgUniqueQuery found",
			prepare: prepareOrgUniqueQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(`SELECT COUNT(*) = 0`+
						` FROM zitadel.projections.orgs`),
					[]string{
						"count",
					},
					[]driver.Value{
						1,
					},
				),
			},
			object: true,
		},
		{
			name:    "prepareOrgUniqueQuery sql err",
			prepare: prepareOrgUniqueQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT COUNT(*) = 0`+
						` FROM zitadel.projections.orgs`),
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
