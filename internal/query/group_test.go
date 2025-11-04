package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	prepareGroupsStmt = `SELECT projections.groups1.id,` +
		` projections.groups1.name,` +
		` projections.groups1.description,` +
		` projections.groups1.creation_date,` +
		` projections.groups1.change_date,` +
		` projections.groups1.resource_owner,` +
		` projections.groups1.instance_id,` +
		` projections.groups1.sequence,` +
		` projections.groups1.state,` +
		` COUNT(*) OVER ()` +
		` FROM projections.groups1`

	prepareGroupsColumns = []string{
		"id",
		"name",
		"description",
		"creation_date",
		"change_date",
		"resource_owner",
		"instance_id",
		"sequence",
		"state",
		"count",
	}

	prepareGroupStmt = `SELECT projections.groups1.id,` +
		` projections.groups1.name,` +
		` projections.groups1.description,` +
		` projections.groups1.creation_date,` +
		` projections.groups1.change_date,` +
		` projections.groups1.resource_owner,` +
		` projections.groups1.instance_id,` +
		` projections.groups1.sequence,` +
		` projections.groups1.state` +
		` FROM projections.groups1`

	prepareGroupColumns = []string{
		"id",
		"name",
		"description",
		"creation_date",
		"change_date",
		"resource_owner",
		"instance_id",
		"sequence",
		"state",
	}
)

func Test_GroupPrepares(t *testing.T) {
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
			name:    "prepareGroupsQuery, no result",
			prepare: prepareGroupsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareGroupsStmt),
					nil,
					nil,
				),
			},
			object: &Groups{Groups: []*Group{}},
		},
		{
			name:    "prepareGroupsQuery, one result",
			prepare: prepareGroupsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareGroupsStmt),
					prepareGroupsColumns,
					[][]driver.Value{
						{
							"9090",
							"group1",
							"my new group",
							testNow,
							testNow,
							"org1",
							"instance1",
							1,
							domain.GroupStateActive,
						},
					},
				),
			},
			object: &Groups{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Groups: []*Group{
					{
						ID:            "9090",
						Name:          "group1",
						Description:   "my new group",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "org1",
						InstanceID:    "instance1",
						Sequence:      1,
						State:         domain.GroupStateActive,
					},
				},
			},
		},
		{
			name:    "prepareGroupsQuery, multiple results",
			prepare: prepareGroupsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareGroupsStmt),
					prepareGroupsColumns,
					[][]driver.Value{
						{
							"9091",
							"group1",
							"my first group",
							testNow,
							testNow,
							"org1",
							"instance1",
							1,
							domain.GroupStateActive,
						},
						{
							"9092",
							"group2",
							"my second group",
							testNow,
							testNow,
							"org1",
							"instance1",
							1,
							domain.GroupStateActive,
						},
					},
				),
			},
			object: &Groups{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Groups: []*Group{
					{
						ID:            "9091",
						Name:          "group1",
						Description:   "my first group",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "org1",
						InstanceID:    "instance1",
						Sequence:      1,
						State:         domain.GroupStateActive,
					},
					{
						ID:            "9092",
						Name:          "group2",
						Description:   "my second group",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "org1",
						InstanceID:    "instance1",
						Sequence:      1,
						State:         domain.GroupStateActive,
					},
				},
			},
		},
		{
			name:    "prepareGroupsQuery, sql err",
			prepare: prepareGroupsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareGroupsStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Groups)(nil),
		},
		{
			name:    "prepareGroupsQuery, no result",
			prepare: prepareGroupsQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					regexp.QuoteMeta(prepareGroupsStmt),
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
			object: &Groups{
				SearchResponse: SearchResponse{
					Count: 0,
				},
				Groups: []*Group{},
			},
		},
		{
			name:    "prepareGroupQuery, no result",
			prepare: prepareGroupQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					prepareGroupStmt,
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
			object: (*Group)(nil),
		},
		{
			name:    "prepareGroupQuery, sql err",
			prepare: prepareGroupQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareGroupStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Group)(nil),
		},
		{
			name:    "prepareGroupQuery, found",
			prepare: prepareGroupQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareGroupStmt),
					prepareGroupColumns,
					[]driver.Value{
						"9090",
						"group1",
						"my new group",
						testNow,
						testNow,
						"org1",
						"instance1",
						1,
						domain.GroupStateActive,
					},
				),
			},
			object: &Group{
				ID:            "9090",
				Name:          "group1",
				Description:   "my new group",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "org1",
				InstanceID:    "instance1",
				Sequence:      1,
				State:         domain.GroupStateActive,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

func Test_GroupsCheckPermission(t *testing.T) {
	tests := []struct {
		name        string
		want        []*Group
		groups      *Groups
		permissions []string
	}{
		{
			name: "no permissions",
			want: []*Group{},
			groups: &Groups{
				Groups: []*Group{
					{ID: "group1"}, {ID: "group2"}, {ID: "group3"},
				},
			},
			permissions: []string{},
		},
		{
			name: "permissions for all groups",
			want: []*Group{
				{ID: "group1"}, {ID: "group2"}, {ID: "group3"},
			},
			groups: &Groups{
				Groups: []*Group{
					{ID: "group1"}, {ID: "group2"}, {ID: "group3"},
				},
			},
			permissions: []string{"group1", "group2", "group3"},
		},
		{
			name: "permissions for group1",
			want: []*Group{
				{ID: "group1"},
			},
			groups: &Groups{
				Groups: []*Group{
					{ID: "group1"}, {ID: "group2"}, {ID: "group3"},
				},
			},
			permissions: []string{"group1"},
		},
		{
			name: "permissions for group2",
			want: []*Group{
				{ID: "group2"},
			},
			groups: &Groups{
				Groups: []*Group{
					{ID: "group1"}, {ID: "group2"}, {ID: "group3"},
				},
			},
			permissions: []string{"group2"},
		},
		{
			name: "permissions for group3",
			want: []*Group{
				{ID: "group3"},
			},
			groups: &Groups{
				Groups: []*Group{
					{ID: "group1"}, {ID: "group2"}, {ID: "group3"},
				},
			},
			permissions: []string{"group3"},
		},
		{
			name: "permissions for group1 and group2",
			want: []*Group{
				{ID: "group1"}, {ID: "group2"},
			},
			groups: &Groups{
				Groups: []*Group{
					{ID: "group1"}, {ID: "group2"}, {ID: "group3"},
				},
			},
			permissions: []string{"group1", "group2"},
		},
		{
			name: "permissions for group1 and group3",
			want: []*Group{
				{ID: "group1"}, {ID: "group3"},
			},
			groups: &Groups{
				Groups: []*Group{
					{ID: "group1"}, {ID: "group2"}, {ID: "group3"},
				},
			},
			permissions: []string{"group1", "group3"},
		},
		{
			name: "permissions for group2 and group3",
			want: []*Group{
				{ID: "group2"}, {ID: "group3"},
			},
			groups: &Groups{
				Groups: []*Group{
					{ID: "group1"}, {ID: "group2"}, {ID: "group3"},
				},
			},
			permissions: []string{"group2", "group3"},
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
			groupsCheckPermission(context.Background(), tt.groups, checkPermission)
			require.Equal(t, tt.want, tt.groups.Groups)
		})
	}

}
