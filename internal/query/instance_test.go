package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/text/language"
)

var (
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
		` FROM (SELECT DISTINCT projections.instances.id, COUNT(*) OVER () FROM projections.instances` +
		` LEFT JOIN projections.instance_domains ON projections.instances.id = projections.instance_domains.instance_id) AS f` +
		` LEFT JOIN projections.instances ON f.id = projections.instances.id` +
		` LEFT JOIN projections.instance_domains ON f.id = projections.instance_domains.instance_id`
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
		prepare        any
		additionalArgs []reflect.Value
		want           want
		object         any
	}{
		{
			name: "prepareInstancesQuery no result",
			prepare: func() (sq.SelectBuilder, func(*sql.Rows) (*Instances, error)) {
				filter, query, scan := prepareInstancesQuery()
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
			prepare: func() (sq.SelectBuilder, func(*sql.Rows) (*Instances, error)) {
				filter, query, scan := prepareInstancesQuery()
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
			prepare: func() (sq.SelectBuilder, func(*sql.Rows) (*Instances, error)) {
				filter, query, scan := prepareInstancesQuery()
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
			prepare: func() (sq.SelectBuilder, func(*sql.Rows) (*Instances, error)) {
				filter, query, scan := prepareInstancesQuery()
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, tt.additionalArgs...)
		})
	}
}
