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
	prepareOrganizationSettingsListStmt = `SELECT projections.organization_settings.id,` +
		` projections.organization_settings.creation_date,` +
		` projections.organization_settings.change_date,` +
		` projections.organization_settings.resource_owner,` +
		` projections.organization_settings.sequence,` +
		` projections.organization_settings.user_uniqueness,` +
		` COUNT(*) OVER ()` +
		` FROM projections.organization_settings`
	prepareOrganizationSettingsListCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"user_uniqueness",
		"count",
	}
)

func Test_OrganizationSettingsListPrepares(t *testing.T) {
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
			name:    "prepareOrganizationSettingsListQuery no result",
			prepare: prepareOrganizationSettingsListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareOrganizationSettingsListStmt),
					nil,
					nil,
				),
			},
			object: &OrganizationSettingsList{OrganizationSettingsList: []*OrganizationSettings{}},
		},
		{
			name:    "prepareOrganizationSettingsListQuery one result",
			prepare: prepareOrganizationSettingsListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareOrganizationSettingsListStmt),
					prepareOrganizationSettingsListCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							"ro",
							uint64(20211108),
							true,
						},
					},
				),
			},
			object: &OrganizationSettingsList{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				OrganizationSettingsList: []*OrganizationSettings{
					{
						ID:                          "id",
						CreationDate:                testNow,
						ChangeDate:                  testNow,
						ResourceOwner:               "ro",
						Sequence:                    20211108,
						OrganizationScopedUsernames: true,
					},
				},
			},
		},
		{
			name:    "prepareOrganizationSettingsListQuery multiple result",
			prepare: prepareOrganizationSettingsListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareOrganizationSettingsListStmt),
					prepareOrganizationSettingsListCols,
					[][]driver.Value{
						{
							"id-1",
							testNow,
							testNow,
							"ro",
							uint64(20211108),
							true,
						},
						{
							"id-2",
							testNow,
							testNow,
							"ro",
							uint64(20211108),
							false,
						},
						{
							"id-3",
							testNow,
							testNow,
							"ro",
							uint64(20211108),
							true,
						},
					},
				),
			},
			object: &OrganizationSettingsList{
				SearchResponse: SearchResponse{
					Count: 3,
				},
				OrganizationSettingsList: []*OrganizationSettings{
					{
						ID:                          "id-1",
						CreationDate:                testNow,
						ChangeDate:                  testNow,
						ResourceOwner:               "ro",
						Sequence:                    20211108,
						OrganizationScopedUsernames: true,
					},
					{
						ID:                          "id-2",
						CreationDate:                testNow,
						ChangeDate:                  testNow,
						ResourceOwner:               "ro",
						Sequence:                    20211108,
						OrganizationScopedUsernames: false,
					},
					{
						ID:                          "id-3",
						CreationDate:                testNow,
						ChangeDate:                  testNow,
						ResourceOwner:               "ro",
						Sequence:                    20211108,
						OrganizationScopedUsernames: true,
					},
				},
			},
		},
		{
			name:    "prepareOrganizationSettingsListQuery sql err",
			prepare: prepareOrganizationSettingsListQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareOrganizationSettingsListStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*OrganizationSettingsList)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
