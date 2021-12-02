package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/lib/pq"
)

var (
	membershipsStmt = regexp.QuoteMeta(
		"SELECT user_id" +
			", roles" +
			", creation_date" +
			", change_date" +
			", sequence" +
			", resource_owner" +
			", org_id" +
			", iam_id" +
			", project_id" +
			", grant_id" +
			", COUNT(*) OVER ()" +
			" FROM (" +
			"SELECT zitadel.projections.org_members.user_id" +
			", zitadel.projections.org_members.roles" +
			", zitadel.projections.org_members.creation_date" +
			", zitadel.projections.org_members.change_date" +
			", zitadel.projections.org_members.sequence" +
			", zitadel.projections.org_members.resource_owner" +
			", zitadel.projections.org_members.org_id" +
			", NULL::STRING AS iam_id" +
			", NULL::STRING AS project_id" +
			", NULL::STRING AS grant_id" +
			" FROM zitadel.projections.org_members" +
			" UNION ALL " +
			"SELECT zitadel.projections.iam_members.user_id" +
			", zitadel.projections.iam_members.roles" +
			", zitadel.projections.iam_members.creation_date" +
			", zitadel.projections.iam_members.change_date" +
			", zitadel.projections.iam_members.sequence" +
			", zitadel.projections.iam_members.resource_owner" +
			", NULL::STRING AS org_id" +
			", zitadel.projections.iam_members.iam_id" +
			", NULL::STRING AS project_id" +
			", NULL::STRING AS grant_id" +
			" FROM zitadel.projections.iam_members" +
			" UNION ALL " +
			"SELECT zitadel.projections.project_members.user_id" +
			", zitadel.projections.project_members.roles" +
			", zitadel.projections.project_members.creation_date" +
			", zitadel.projections.project_members.change_date" +
			", zitadel.projections.project_members.sequence" +
			", zitadel.projections.project_members.resource_owner" +
			", NULL::STRING AS org_id" +
			", NULL::STRING AS iam_id" +
			", zitadel.projections.project_members.project_id" +
			", NULL::STRING AS grant_id" +
			" FROM zitadel.projections.project_members" +
			" UNION ALL " +
			"SELECT zitadel.projections.project_grant_members.user_id" +
			", zitadel.projections.project_grant_members.roles" +
			", zitadel.projections.project_grant_members.creation_date" +
			", zitadel.projections.project_grant_members.change_date" +
			", zitadel.projections.project_grant_members.sequence" +
			", zitadel.projections.project_grant_members.resource_owner" +
			", NULL::STRING AS org_id" +
			", NULL::STRING AS iam_id" +
			", zitadel.projections.project_grant_members.project_id" +
			", zitadel.projections.project_grant_members.grant_id" +
			" FROM zitadel.projections.project_grant_members" +
			") AS m")
	membershipCols = []string{
		"user_id",
		"roles",
		"creation_date",
		"change_date",
		"sequence",
		"resource_owner",
		"org_id",
		"iam_id",
		"project_id",
		"grant_id",
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
			prepare: prepareMembershipsQuery,
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
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							"org-id",
							nil,
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
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						Org:           &OrgMembership{OrgID: "org-id"},
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one iam member",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							"iam-id",
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
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						IAM:           &IAMMembership{IAMID: "iam-id"},
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one project member",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							nil,
							"project-id",
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
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						Project:       &ProjectMembership{ProjectID: "project-id"},
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one project grant member",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							nil,
							"project-id",
							"grant-id",
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
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						ProjectGrant: &ProjectGrantMembership{
							GrantID:   "grant-id",
							ProjectID: "project-id",
						},
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one for each member type",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							"org-id",
							nil,
							nil,
							nil,
						},
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							"iam-id",
							nil,
							nil,
						},
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							nil,
							"project-id",
							nil,
						},
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							nil,
							"project-id",
							"grant-id",
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
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						Org:           &OrgMembership{OrgID: "org-id"},
					},
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						IAM:           &IAMMembership{IAMID: "iam-id"},
					},
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						Project:       &ProjectMembership{ProjectID: "project-id"},
					},
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						ProjectGrant: &ProjectGrantMembership{
							ProjectID: "project-id",
							GrantID:   "grant-id",
						},
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery sql err",
			prepare: prepareMembershipsQuery,
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
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
