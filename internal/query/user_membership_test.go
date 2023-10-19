package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	membershipsStmt = regexp.QuoteMeta(
		"SELECT members.user_id" +
			", members.roles" +
			", members.creation_date" +
			", members.change_date" +
			", members.sequence" +
			", members.resource_owner" +
			", members.org_id" +
			", members.id" +
			", members.project_id" +
			", members.grant_id" +
			", projections.project_grants3.granted_org_id" +
			", projections.projects3.name" +
			", projections.orgs1.name" +
			", projections.instances.name" +
			", COUNT(*) OVER ()" +
			" FROM (" +
			"SELECT members.user_id" +
			", members.roles" +
			", members.creation_date" +
			", members.change_date" +
			", members.sequence" +
			", members.resource_owner" +
			", members.instance_id" +
			", members.org_id" +
			", NULL::TEXT AS id" +
			", NULL::TEXT AS project_id" +
			", NULL::TEXT AS grant_id" +
			" FROM projections.org_members3 AS members" +
			" WHERE members.owner_removed = $1 AND members.user_owner_removed = $2" +
			" UNION ALL " +
			"SELECT members.user_id" +
			", members.roles" +
			", members.creation_date" +
			", members.change_date" +
			", members.sequence" +
			", members.resource_owner" +
			", members.instance_id" +
			", NULL::TEXT AS org_id" +
			", members.id" +
			", NULL::TEXT AS project_id" +
			", NULL::TEXT AS grant_id" +
			" FROM projections.instance_members3 AS members" +
			" WHERE members.owner_removed = $3 AND members.user_owner_removed = $4" +
			" UNION ALL " +
			"SELECT members.user_id" +
			", members.roles" +
			", members.creation_date" +
			", members.change_date" +
			", members.sequence" +
			", members.resource_owner" +
			", members.instance_id" +
			", NULL::TEXT AS org_id" +
			", NULL::TEXT AS id" +
			", members.project_id" +
			", NULL::TEXT AS grant_id" +
			" FROM projections.project_members3 AS members" +
			" WHERE members.owner_removed = $5 AND members.user_owner_removed = $6" +
			" UNION ALL " +
			"SELECT members.user_id" +
			", members.roles" +
			", members.creation_date" +
			", members.change_date" +
			", members.sequence" +
			", members.resource_owner" +
			", members.instance_id" +
			", NULL::TEXT AS org_id" +
			", NULL::TEXT AS id" +
			", members.project_id" +
			", members.grant_id" +
			" FROM projections.project_grant_members3 AS members" +
			" WHERE members.granted_org_removed = $7 AND members.owner_removed = $8 AND members.user_owner_removed = $9" +
			") AS members" +
			" LEFT JOIN projections.projects3 ON members.project_id = projections.projects3.id AND members.instance_id = projections.projects3.instance_id" +
			" LEFT JOIN projections.orgs1 ON members.org_id = projections.orgs1.id AND members.instance_id = projections.orgs1.instance_id" +
			" LEFT JOIN projections.project_grants3 ON members.grant_id = projections.project_grants3.grant_id AND members.instance_id = projections.project_grants3.instance_id" +
			" LEFT JOIN projections.instances ON members.instance_id = projections.instances.id" +
			` AS OF SYSTEM TIME '-1 ms'`)
	membershipCols = []string{
		"user_id",
		"roles",
		"creation_date",
		"change_date",
		"sequence",
		"resource_owner",
		"org_id",
		"instance_id",
		"project_id",
		"grant_id",
		"granted_org_id",
		"name", //project name
		"name", //org name
		"name", // instance name
		"count",
	}
)

func Test_MembershipPrepares(t *testing.T) {
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
			name:    "prepareMembershipsQuery no result",
			prepare: prepareMembershipWrapper(false),
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					nil,
					nil,
				),
			},
			object: &Memberships{Memberships: []*Membership{}},
		},
		{
			name:    "prepareMembershipsQuery one org member",
			prepare: prepareMembershipWrapper(false),
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							database.TextArray[string]{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							"org-id",
							nil,
							nil,
							nil,
							nil,
							nil,
							"org-name",
							nil,
						},
					},
				),
			},
			object: &Memberships{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Memberships: []*Membership{
					{
						UserID:        "user-id",
						Roles:         database.TextArray[string]{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						Org:           &OrgMembership{OrgID: "org-id", Name: "org-name"},
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one instance member",
			prepare: prepareMembershipWrapper(false),
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							database.TextArray[string]{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							"iam-id",
							nil,
							nil,
							nil,
							nil,
							nil,
							"instance",
						},
					},
				),
			},
			object: &Memberships{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Memberships: []*Membership{
					{
						UserID:        "user-id",
						Roles:         database.TextArray[string]{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						IAM:           &IAMMembership{IAMID: "iam-id", Name: "instance"},
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one project member",
			prepare: prepareMembershipWrapper(false),
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							database.TextArray[string]{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							nil,
							"project-id",
							nil,
							nil,
							"project-name",
							nil,
							nil,
						},
					},
				),
			},
			object: &Memberships{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Memberships: []*Membership{
					{
						UserID:        "user-id",
						Roles:         database.TextArray[string]{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						Project:       &ProjectMembership{ProjectID: "project-id", Name: "project-name"},
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one project grant member",
			prepare: prepareMembershipWrapper(false),
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							database.TextArray[string]{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							nil,
							"project-id",
							"grant-id",
							"granted-org-id",
							"project-name",
							nil,
							nil,
						},
					},
				),
			},
			object: &Memberships{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Memberships: []*Membership{
					{
						UserID:        "user-id",
						Roles:         database.TextArray[string]{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						ProjectGrant: &ProjectGrantMembership{
							GrantID:      "grant-id",
							ProjectID:    "project-id",
							ProjectName:  "project-name",
							GrantedOrgID: "granted-org-id",
						},
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one for each member type",
			prepare: prepareMembershipWrapper(false),
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							database.TextArray[string]{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							"org-id",
							nil,
							nil,
							nil,
							nil,
							nil,
							"org-name",
							nil,
						},
						{
							"user-id",
							database.TextArray[string]{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							"iam-id",
							nil,
							nil,
							nil,
							nil,
							nil,
							"instance",
						},
						{
							"user-id",
							database.TextArray[string]{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							nil,
							"project-id",
							nil,
							nil,
							"project-name",
							nil,
							nil,
						},
						{
							"user-id",
							database.TextArray[string]{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							nil,
							"project-id",
							"grant-id",
							"granted-org-id",
							"project-name",
							nil,
							nil,
						},
					},
				),
			},
			object: &Memberships{
				SearchResponse: SearchResponse{
					Count: 4,
				},
				Memberships: []*Membership{
					{
						UserID:        "user-id",
						Roles:         database.TextArray[string]{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						Org:           &OrgMembership{OrgID: "org-id", Name: "org-name"},
					},
					{
						UserID:        "user-id",
						Roles:         database.TextArray[string]{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						IAM:           &IAMMembership{IAMID: "iam-id", Name: "instance"},
					},
					{
						UserID:        "user-id",
						Roles:         database.TextArray[string]{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						Project:       &ProjectMembership{ProjectID: "project-id", Name: "project-name"},
					},
					{
						UserID:        "user-id",
						Roles:         database.TextArray[string]{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						ProjectGrant: &ProjectGrantMembership{
							ProjectID:    "project-id",
							GrantID:      "grant-id",
							ProjectName:  "project-name",
							GrantedOrgID: "granted-org-id",
						},
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery sql err",
			prepare: prepareMembershipWrapper(false),
			want: want{
				sqlExpectations: mockQueryErr(
					membershipsStmt,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Memberships)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}

func prepareMembershipWrapper(withOwnerRemoved bool) func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Memberships, error)) {
	return func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Memberships, error)) {
		builder, _, fun := prepareMembershipsQuery(ctx, db, withOwnerRemoved, &MembershipSearchQuery{})
		return builder, fun
	}
}
