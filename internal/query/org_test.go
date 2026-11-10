package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
	db_mock "github.com/zitadel/zitadel/internal/database/mock"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	orgUniqueQuery = "SELECT COUNT(*) = 0 FROM projections.orgs1 LEFT JOIN projections.org_domains2 ON projections.orgs1.id = projections.org_domains2.org_id AND projections.orgs1.instance_id = projections.org_domains2.instance_id WHERE (projections.org_domains2.is_verified = $1 AND projections.orgs1.instance_id = $2 AND (projections.org_domains2.domain ILIKE $3 OR projections.orgs1.name ILIKE $4) AND projections.orgs1.org_state <> $5)"
	orgUniqueCols  = []string{"is_unique"}

	prepareOrgsQueryStmt = `SELECT projections.orgs1.id,` +
		` projections.orgs1.creation_date,` +
		` projections.orgs1.change_date,` +
		` projections.orgs1.resource_owner,` +
		` projections.orgs1.org_state,` +
		` projections.orgs1.sequence,` +
		` projections.orgs1.name,` +
		` projections.orgs1.primary_domain,` +
		` COUNT(*) OVER ()` +
		` FROM projections.orgs1`
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

	prepareOrgQueryStmt = `SELECT projections.orgs1.id,` +
		` projections.orgs1.creation_date,` +
		` projections.orgs1.change_date,` +
		` projections.orgs1.resource_owner,` +
		` projections.orgs1.org_state,` +
		` projections.orgs1.sequence,` +
		` projections.orgs1.instance_id,` +
		` projections.orgs1.name,` +
		` projections.orgs1.primary_domain` +
		` FROM projections.orgs1`
	prepareOrgQueryCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"org_state",
		"sequence",
		"instance_id",
		"name",
		"primary_domain",
	}

	prepareOrgUniqueStmt = `SELECT COUNT(*) = 0` +
		` FROM projections.orgs1` +
		` LEFT JOIN projections.org_domains2 ON projections.orgs1.id = projections.org_domains2.org_id AND projections.orgs1.instance_id = projections.org_domains2.instance_id`
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
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Orgs)(nil),
		},
		{
			name:    "prepareOrgQuery no result",
			prepare: prepareOrgQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					regexp.QuoteMeta(prepareOrgQueryStmt),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsNotFound(err) {
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
						"instance-id",
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
				instanceID:    "instance-id",
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
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Org)(nil),
		},
		{
			name:    "prepareOrgUniqueQuery no result",
			prepare: prepareOrgUniqueQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					regexp.QuoteMeta(prepareOrgUniqueStmt),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsInternal(err) {
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
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
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
				err:      zerrors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		client, mock, err := sqlmock.New(
			sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
			sqlmock.ValueConverterOption(new(db_mock.TypeConverter)),
		)
		if err != nil {
			t.Fatalf("unable to mock db: %v", err)
		}
		if tt.want.sqlExpectations != nil {
			tt.want.sqlExpectations(mock)
		}

		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				client: &database.DB{
					DB: client,
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

func TestOrg_orgsCheckPermission(t *testing.T) {
	type want struct {
		orgs []*Org
	}
	tests := []struct {
		name        string
		want        want
		orgs        *Orgs
		permissions []string
	}{
		{
			"permissions for all",
			want{
				orgs: []*Org{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			&Orgs{
				Orgs: []*Org{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			[]string{"first", "second", "third"},
		},
		{
			"permissions for one, first",
			want{
				orgs: []*Org{
					{ID: "first"},
				},
			},
			&Orgs{
				Orgs: []*Org{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			[]string{"first"},
		},
		{
			"permissions for one, second",
			want{
				orgs: []*Org{
					{ID: "second"},
				},
			},
			&Orgs{
				Orgs: []*Org{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			[]string{"second"},
		},
		{
			"permissions for one, third",
			want{
				orgs: []*Org{
					{ID: "third"},
				},
			},
			&Orgs{
				Orgs: []*Org{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			[]string{"third"},
		},
		{
			"permissions for two, first third",
			want{
				orgs: []*Org{
					{ID: "first"}, {ID: "third"},
				},
			},
			&Orgs{
				Orgs: []*Org{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			[]string{"first", "third"},
		},
		{
			"permissions for two, second third",
			want{
				orgs: []*Org{
					{ID: "second"}, {ID: "third"},
				},
			},
			&Orgs{
				Orgs: []*Org{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			[]string{"second", "third"},
		},
		{
			"no permissions",
			want{
				orgs: []*Org{},
			},
			&Orgs{
				Orgs: []*Org{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkPermission := func(ctx context.Context, permission, orgID, resourceID string) (err error) {
				for _, perm := range tt.permissions {
					if resourceID == perm {
						return nil
					}
				}
				return errors.New("failed")
			}
			orgsCheckPermission(context.Background(), tt.orgs, checkPermission)
			require.Equal(t, tt.want.orgs, tt.orgs.Orgs)
		})
	}
}
