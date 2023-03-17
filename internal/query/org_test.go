package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	errs "errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
)

var (
	orgUniqueQuery = "SELECT COUNT(*) = 0 FROM projections.orgs LEFT JOIN projections.org_domains2 ON projections.orgs.id = projections.org_domains2.org_id AND projections.orgs.instance_id = projections.org_domains2.instance_id AS OF SYSTEM TIME '-1 ms' WHERE (projections.org_domains2.is_verified = $1 AND projections.orgs.instance_id = $2 AND (projections.org_domains2.domain ILIKE $3 OR projections.orgs.name ILIKE $4) AND projections.orgs.org_state <> $5)"
	orgUniqueCols  = []string{"is_unique"}

	prepareOrgsQueryStmt = `SELECT projections.orgs.id,` +
		` projections.orgs.creation_date,` +
		` projections.orgs.change_date,` +
		` projections.orgs.resource_owner,` +
		` projections.orgs.org_state,` +
		` projections.orgs.sequence,` +
		` projections.orgs.name,` +
		` projections.orgs.primary_domain,` +
		` COUNT(*) OVER ()` +
		` FROM projections.orgs` +
		` AS OF SYSTEM TIME '-1 ms' `
	prepareOrgsQueryCols = []string{
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

	prepareOrgQueryStmt = `SELECT projections.orgs.id,` +
		` projections.orgs.creation_date,` +
		` projections.orgs.change_date,` +
		` projections.orgs.resource_owner,` +
		` projections.orgs.org_state,` +
		` projections.orgs.sequence,` +
		` projections.orgs.name,` +
		` projections.orgs.primary_domain` +
		` FROM projections.orgs` +
		` AS OF SYSTEM TIME '-1 ms' `
	prepareOrgQueryCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"org_state",
		"sequence",
		"name",
		"primary_domain",
	}

	prepareOrgUniqueStmt = `SELECT COUNT(*) = 0` +
		` FROM projections.orgs` +
		` LEFT JOIN projections.org_domains2 ON projections.orgs.id = projections.org_domains2.org_id AND projections.orgs.instance_id = projections.org_domains2.instance_id` +
		` AS OF SYSTEM TIME '-1 ms' `
	prepareOrgUniqueCols = []string{
		"count",
	}
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
					regexp.QuoteMeta(prepareOrgsQueryStmt),
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
					regexp.QuoteMeta(prepareOrgsQueryStmt),
					prepareOrgsQueryCols,
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
					regexp.QuoteMeta(prepareOrgsQueryStmt),
					prepareOrgsQueryCols,
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
					regexp.QuoteMeta(prepareOrgsQueryStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errs.Is(err, sql.ErrConnDone) {
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
					regexp.QuoteMeta(prepareOrgQueryStmt),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errors.IsNotFound(err) {
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
					regexp.QuoteMeta(prepareOrgQueryStmt),
					prepareOrgQueryCols,
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
					regexp.QuoteMeta(prepareOrgQueryStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errs.Is(err, sql.ErrConnDone) {
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
					regexp.QuoteMeta(prepareOrgUniqueStmt),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errors.IsInternal(err) {
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
					regexp.QuoteMeta(prepareOrgUniqueStmt),
					prepareOrgUniqueCols,
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
					regexp.QuoteMeta(prepareOrgUniqueStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errs.Is(err, sql.ErrConnDone) {
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

func TestQueries_IsOrgUnique(t *testing.T) {
	type args struct {
		name   string
		domain string
	}
	type want struct {
		err             func(error) bool
		sqlExpectations sqlExpectation
		isUnique        bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "existing domain",
			args: args{
				domain: "exists",
				name:   "",
			},
			want: want{
				isUnique:        false,
				sqlExpectations: mockQueries(orgUniqueQuery, orgUniqueCols, [][]driver.Value{{false}}, true, "", "exists", "", domain.OrgStateRemoved),
			},
		},
		{
			name: "existing name",
			args: args{
				domain: "",
				name:   "exists",
			},
			want: want{
				isUnique:        false,
				sqlExpectations: mockQueries(orgUniqueQuery, orgUniqueCols, [][]driver.Value{{false}}, true, "", "", "exists", domain.OrgStateRemoved),
			},
		},
		{
			name: "existing name and domain",
			args: args{
				domain: "exists",
				name:   "exists",
			},
			want: want{
				isUnique:        false,
				sqlExpectations: mockQueries(orgUniqueQuery, orgUniqueCols, [][]driver.Value{{false}}, true, "", "exists", "exists", domain.OrgStateRemoved),
			},
		},
		{
			name: "not existing",
			args: args{
				domain: "not-exists",
				name:   "not-exists",
			},
			want: want{
				isUnique:        true,
				sqlExpectations: mockQueries(orgUniqueQuery, orgUniqueCols, [][]driver.Value{{true}}, true, "", "not-exists", "not-exists", domain.OrgStateRemoved),
			},
		},
		{
			name: "no arg",
			args: args{
				domain: "",
				name:   "",
			},
			want: want{
				isUnique: false,
				err:      errors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		client, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("unable to mock db: %v", err)
		}
		if tt.want.sqlExpectations != nil {
			tt.want.sqlExpectations(mock)
		}

		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				client: &database.DB{
					DB:       client,
					Database: new(prepareDB),
				},
			}

			gotIsUnique, err := q.IsOrgUnique(context.Background(), tt.args.name, tt.args.domain)
			if (tt.want.err == nil && err != nil) || (err != nil && tt.want.err != nil && !tt.want.err(err)) {
				t.Errorf("Queries.IsOrgUnique() unexpected error = %v", err)
				return
			}
			if gotIsUnique != tt.want.isUnique {
				t.Errorf("Queries.IsOrgUnique() = %v, want %v", gotIsUnique, tt.want.isUnique)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectation was met: %v", err)
			}
		})

	}
}
