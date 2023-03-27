package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/text/language"

	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	instanceQuery = `SELECT projections.instances.id,` +
		` projections.instances.creation_date,` +
		` projections.instances.change_date,` +
		` projections.instances.sequence,` +
		` projections.instances.default_org_id,` +
		` projections.instances.iam_project_id,` +
		` projections.instances.console_client_id,` +
		` projections.instances.console_app_id,` +
		` projections.instances.default_language` +
		` FROM projections.instances` +
		` AS OF SYSTEM TIME '-1 ms'`
	instanceCols = []string{
		"id",
		"creation_date",
		"change_date",
		"sequence",
		"default_org_id",
		"iam_project_id",
		"console_client_id",
		"console_app_id",
		"default_language",
	}
	instancesQuery = `SELECT f.count, f.id,` +
		` projections.instances.creation_date,` +
		` projections.instances.change_date,` +
		` projections.instances.sequence,` +
		` projections.instances.name,` +
		` projections.instances.default_org_id,` +
		` projections.instances.iam_project_id,` +
		` projections.instances.console_client_id,` +
		` projections.instances.console_app_id,` +
		` projections.instances.default_language,` +
		` projections.instance_domains.domain,` +
		` projections.instance_domains.is_primary,` +
		` projections.instance_domains.is_generated,` +
		` projections.instance_domains.creation_date,` +
		` projections.instance_domains.change_date, ` +
		` projections.instance_domains.sequence` +
		` FROM (SELECT projections.instances.id, COUNT(*) OVER () FROM projections.instances) AS f` +
		` LEFT JOIN projections.instances ON f.id = projections.instances.id` +
		` LEFT JOIN projections.instance_domains ON f.id = projections.instance_domains.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`
	instancesCols = []string{
		"count",
		"id",
		"creation_date",
		"change_date",
		"sequence",
		"name",
		"default_org_id",
		"iam_project_id",
		"console_client_id",
		"console_app_id",
		"default_language",
		"domain",
		"is_primary",
		"is_generated",
		"creation_date",
		"change_date",
		"sequence",
	}
)

func Test_InstancePrepares(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name           string
		prepare        interface{}
		additionalArgs []reflect.Value
		want           want
		object         interface{}
	}{
		{
			name:           "prepareInstanceQuery no result",
			additionalArgs: []reflect.Value{reflect.ValueOf("")},
			prepare:        prepareInstanceQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(instanceQuery),
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
			object: (*Instance)(nil),
		},
		{
			name:           "prepareInstanceQuery found",
			additionalArgs: []reflect.Value{reflect.ValueOf("")},
			prepare:        prepareInstanceQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(instanceQuery),
					instanceCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						uint64(20211108),
						"global-org-id",
						"project-id",
						"client-id",
						"app-id",
						"en",
					},
				),
			},
			object: &Instance{
				ID:           "id",
				CreationDate: testNow,
				ChangeDate:   testNow,
				Sequence:     20211108,
				DefaultOrgID: "global-org-id",
				IAMProjectID: "project-id",
				ConsoleID:    "client-id",
				ConsoleAppID: "app-id",
				DefaultLang:  language.English,
			},
		},
		{
			name:           "prepareInstanceQuery sql err",
			additionalArgs: []reflect.Value{reflect.ValueOf("")},
			prepare:        prepareInstanceQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(instanceQuery),
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
			name: "prepareInstancesQuery no result",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Instances, error)) {
				filter, query, scan := prepareInstancesQuery(ctx, db)
				return query(filter), scan
			},
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(instancesQuery),
					nil,
					nil,
				),
			},
			object: &Instances{Instances: []*Instance{}},
		},
		{
			name: "prepareInstancesQuery one result",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Instances, error)) {
				filter, query, scan := prepareInstancesQuery(ctx, db)
				return query(filter), scan
			},
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(instancesQuery),
					instancesCols,
					[][]driver.Value{
						{
							"1",
							"id",
							testNow,
							testNow,
							uint64(20211108),
							"test",
							"global-org-id",
							"project-id",
							"client-id",
							"app-id",
							"en",
							"test.zitadel.cloud",
							true,
							true,
							testNow,
							testNow,
							uint64(20211108),
						},
					},
				),
			},
			object: &Instances{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Instances: []*Instance{
					{
						ID:           "id",
						CreationDate: testNow,
						ChangeDate:   testNow,
						Sequence:     20211108,
						Name:         "test",
						DefaultOrgID: "global-org-id",
						IAMProjectID: "project-id",
						ConsoleID:    "client-id",
						ConsoleAppID: "app-id",
						DefaultLang:  language.English,
						Domains: []*InstanceDomain{
							{
								CreationDate: testNow,
								ChangeDate:   testNow,
								Sequence:     20211108,
								InstanceID:   "id",
								Domain:       "test.zitadel.cloud",
								IsGenerated:  true,
								IsPrimary:    true,
							},
						},
					},
				},
			},
		},
		{
			name: "prepareInstancesQuery multiple results",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Instances, error)) {
				filter, query, scan := prepareInstancesQuery(ctx, db)
				return query(filter), scan
			},
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(instancesQuery),
					instancesCols,
					[][]driver.Value{
						{
							2,
							"id",
							testNow,
							testNow,
							uint64(20211108),
							"test",
							"global-org-id",
							"project-id",
							"client-id",
							"app-id",
							"en",
							"test.zitadel.cloud",
							false,
							true,
							testNow,
							testNow,
							uint64(20211108),
						},
						{
							2,
							"id",
							testNow,
							testNow,
							uint64(20211108),
							"test",
							"global-org-id",
							"project-id",
							"client-id",
							"app-id",
							"en",
							"zitadel.cloud",
							true,
							false,
							testNow,
							testNow,
							uint64(20211108),
						},
						{
							2,
							"id2",
							testNow,
							testNow,
							uint64(20211108),
							"test2",
							"global-org-id",
							"project-id",
							"client-id",
							"app-id",
							"en",
							"test2.zitadel.cloud",
							true,
							true,
							testNow,
							testNow,
							uint64(20211108),
						},
					},
				),
			},
			object: &Instances{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Instances: []*Instance{
					{
						ID:           "id",
						CreationDate: testNow,
						ChangeDate:   testNow,
						Sequence:     20211108,
						Name:         "test",
						DefaultOrgID: "global-org-id",
						IAMProjectID: "project-id",
						ConsoleID:    "client-id",
						ConsoleAppID: "app-id",
						DefaultLang:  language.English,
						Domains: []*InstanceDomain{
							{
								CreationDate: testNow,
								ChangeDate:   testNow,
								Sequence:     20211108,
								Domain:       "test.zitadel.cloud",
								InstanceID:   "id",
								IsGenerated:  true,
								IsPrimary:    false,
							},
							{
								CreationDate: testNow,
								ChangeDate:   testNow,
								Sequence:     20211108,
								Domain:       "zitadel.cloud",
								InstanceID:   "id",
								IsGenerated:  false,
								IsPrimary:    true,
							},
						},
					}, {
						ID:           "id2",
						CreationDate: testNow,
						ChangeDate:   testNow,
						Sequence:     20211108,
						Name:         "test2",
						DefaultOrgID: "global-org-id",
						IAMProjectID: "project-id",
						ConsoleID:    "client-id",
						ConsoleAppID: "app-id",
						DefaultLang:  language.English,
						Domains: []*InstanceDomain{
							{
								CreationDate: testNow,
								ChangeDate:   testNow,
								Sequence:     20211108,
								Domain:       "test2.zitadel.cloud",
								InstanceID:   "id2",
								IsGenerated:  true,
								IsPrimary:    true,
							},
						},
					},
				},
			},
		},
		{
			name: "prepareInstancesQuery sql err",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Instances, error)) {
				filter, query, scan := prepareInstancesQuery(ctx, db)
				return query(filter), scan
			},
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(instancesQuery),
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, append(defaultPrepareArgs, tt.additionalArgs...)...)
		})
	}
}
