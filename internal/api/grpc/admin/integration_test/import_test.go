//go:build integration

package admin_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	v1 "github.com/zitadel/zitadel/pkg/grpc/v1"
)

func TestServer_ImportData(t *testing.T) {
	orgIDs := generateIDs(10)
	projectIDs := generateIDs(10)
	userIDs := generateIDs(10)
	grantIDs := generateIDs(10)

	tests := []struct {
		name    string
		req     *admin.ImportDataRequest
		want    *admin.ImportDataResponse
		wantErr bool
	}{
		{
			name: "success",
			req: &admin.ImportDataRequest{
				Data: &admin.ImportDataRequest_DataOrgs{
					DataOrgs: &admin.ImportDataOrg{
						Orgs: []*admin.DataOrg{
							{
								OrgId: orgIDs[0],
								Org: &management.AddOrgRequest{
									Name: integration.OrganizationName(),
								},
								Projects: []*v1.DataProject{
									{
										ProjectId: projectIDs[0],
										Project: &management.AddProjectRequest{
											Name:                 integration.ProjectName(),
											ProjectRoleAssertion: true,
										},
									},
									{
										ProjectId: projectIDs[1],
										Project: &management.AddProjectRequest{
											Name:                 integration.ProjectName(),
											ProjectRoleAssertion: false,
										},
									},
								},
								ProjectRoles: []*management.AddProjectRoleRequest{
									{
										ProjectId:   projectIDs[0],
										RoleKey:     "role1",
										DisplayName: "role1",
									},
									{
										ProjectId:   projectIDs[0],
										RoleKey:     "role2",
										DisplayName: "role2",
									},
									{
										ProjectId:   projectIDs[1],
										RoleKey:     "role3",
										DisplayName: "role3",
									},
									{
										ProjectId:   projectIDs[1],
										RoleKey:     "role4",
										DisplayName: "role4",
									},
								},
								HumanUsers: []*v1.DataHumanUser{
									{
										UserId: userIDs[0],
										User: &management.ImportHumanUserRequest{
											UserName: integration.Username(),
											Profile: &management.ImportHumanUserRequest_Profile{
												FirstName:         integration.FirstName(),
												LastName:          integration.LastName(),
												DisplayName:       integration.Username(),
												PreferredLanguage: integration.Language(),
											},
											Email: &management.ImportHumanUserRequest_Email{
												Email:           integration.Email(),
												IsEmailVerified: true,
											},
										},
									},
									{
										UserId: userIDs[1],
										User: &management.ImportHumanUserRequest{
											UserName: integration.Username(),
											Profile: &management.ImportHumanUserRequest_Profile{
												FirstName:         integration.FirstName(),
												LastName:          integration.LastName(),
												DisplayName:       integration.Username(),
												PreferredLanguage: integration.Language(),
											},
											Email: &management.ImportHumanUserRequest_Email{
												Email:           integration.Email(),
												IsEmailVerified: true,
											},
										},
									},
								},
								ProjectGrants: []*v1.DataProjectGrant{
									{
										GrantId: grantIDs[0],
										ProjectGrant: &management.AddProjectGrantRequest{
											ProjectId:    projectIDs[0],
											GrantedOrgId: orgIDs[1],
											RoleKeys:     []string{"role1", "role2"},
										},
									},
									{
										GrantId: grantIDs[1],
										ProjectGrant: &management.AddProjectGrantRequest{
											ProjectId:    projectIDs[1],
											GrantedOrgId: orgIDs[1],
											RoleKeys:     []string{"role3", "role4"},
										},
									},
									{
										GrantId: grantIDs[2],
										ProjectGrant: &management.AddProjectGrantRequest{
											ProjectId:    projectIDs[0],
											GrantedOrgId: orgIDs[2],
											RoleKeys:     []string{"role1", "role2"},
										},
									},
									{
										GrantId: grantIDs[3],
										ProjectGrant: &management.AddProjectGrantRequest{
											ProjectId:    projectIDs[1],
											GrantedOrgId: orgIDs[2],
											RoleKeys:     []string{"role3", "role4"},
										},
									},
								},
							},
							{
								OrgId: orgIDs[1],
								Org: &management.AddOrgRequest{
									Name: integration.OrganizationName(),
								},
								UserGrants: []*management.AddUserGrantRequest{
									{
										UserId:         userIDs[0],
										ProjectId:      projectIDs[0],
										ProjectGrantId: grantIDs[0],
									},
									{
										UserId:         userIDs[0],
										ProjectId:      projectIDs[1],
										ProjectGrantId: grantIDs[1],
									},
								},
							},
							{
								OrgId: orgIDs[2],
								Org: &management.AddOrgRequest{
									Name: integration.OrganizationName(),
								},
								UserGrants: []*management.AddUserGrantRequest{
									{
										UserId:         userIDs[1],
										ProjectId:      projectIDs[0],
										ProjectGrantId: grantIDs[2],
									},
									{
										UserId:         userIDs[1],
										ProjectId:      projectIDs[1],
										ProjectGrantId: grantIDs[3],
									},
								},
							},
						},
					},
				},
				Timeout: time.Minute.String(),
			},
			want: &admin.ImportDataResponse{
				Success: &admin.ImportDataSuccess{
					Orgs: []*admin.ImportDataSuccessOrg{
						{
							OrgId:      orgIDs[0],
							ProjectIds: projectIDs[0:2],
							ProjectRoles: []string{
								projectIDs[0] + "_role1",
								projectIDs[0] + "_role2",
								projectIDs[1] + "_role3",
								projectIDs[1] + "_role4",
							},
							HumanUserIds: userIDs[0:2],
							ProjectGrants: []*admin.ImportDataSuccessProjectGrant{
								{
									GrantId:   grantIDs[0],
									ProjectId: projectIDs[0],
									OrgId:     orgIDs[1],
								},
								{
									GrantId:   grantIDs[1],
									ProjectId: projectIDs[1],
									OrgId:     orgIDs[1],
								},
								{
									GrantId:   grantIDs[2],
									ProjectId: projectIDs[0],
									OrgId:     orgIDs[2],
								},
								{
									GrantId:   grantIDs[3],
									ProjectId: projectIDs[1],
									OrgId:     orgIDs[2],
								},
							},
						},
						{
							OrgId: orgIDs[1],
							UserGrants: []*admin.ImportDataSuccessUserGrant{
								{
									ProjectId: projectIDs[0],
									UserId:    userIDs[0],
								},
								{
									UserId:    userIDs[0],
									ProjectId: projectIDs[1],
								},
							},
						},
						{
							OrgId: orgIDs[2],
							UserGrants: []*admin.ImportDataSuccessUserGrant{
								{
									ProjectId: projectIDs[0],
									UserId:    userIDs[1],
								},
								{
									UserId:    userIDs[1],
									ProjectId: projectIDs[1],
								},
							},
						},
					},
				},
			},
		},
		{
			name: "duplicate project grant error",
			req: &admin.ImportDataRequest{
				Data: &admin.ImportDataRequest_DataOrgs{
					DataOrgs: &admin.ImportDataOrg{
						Orgs: []*admin.DataOrg{
							{
								OrgId: orgIDs[4],
								Org: &management.AddOrgRequest{
									Name: integration.OrganizationName(),
								},
							},
							{
								OrgId: orgIDs[3],
								Org: &management.AddOrgRequest{
									Name: integration.OrganizationName(),
								},
								Projects: []*v1.DataProject{
									{
										ProjectId: projectIDs[2],
										Project: &management.AddProjectRequest{
											Name:                 integration.ProjectName(),
											ProjectRoleAssertion: true,
										},
									},
									{
										ProjectId: projectIDs[3],
										Project: &management.AddProjectRequest{
											Name:                 integration.ProjectName(),
											ProjectRoleAssertion: false,
										},
									},
								},
								ProjectRoles: []*management.AddProjectRoleRequest{
									{
										ProjectId:   projectIDs[2],
										RoleKey:     "role1",
										DisplayName: "role1",
									},
									{
										ProjectId:   projectIDs[2],
										RoleKey:     "role2",
										DisplayName: "role2",
									},
									{
										ProjectId:   projectIDs[3],
										RoleKey:     "role3",
										DisplayName: "role3",
									},
									{
										ProjectId:   projectIDs[3],
										RoleKey:     "role4",
										DisplayName: "role4",
									},
								},
								ProjectGrants: []*v1.DataProjectGrant{
									{
										GrantId: grantIDs[4],
										ProjectGrant: &management.AddProjectGrantRequest{
											ProjectId:    projectIDs[2],
											GrantedOrgId: orgIDs[4],
											RoleKeys:     []string{"role1", "role2"},
										},
									},
									{
										GrantId: grantIDs[4],
										ProjectGrant: &management.AddProjectGrantRequest{
											ProjectId:    projectIDs[2],
											GrantedOrgId: orgIDs[4],
											RoleKeys:     []string{"role1", "role2"},
										},
									},
								},
							},
						},
					},
				},
				Timeout: time.Minute.String(),
			},
			want: &admin.ImportDataResponse{
				Errors: []*admin.ImportDataError{
					{
						Type:    "project_grant",
						Id:      orgIDs[3] + "_" + projectIDs[2] + "_" + orgIDs[4],
						Message: "ID=V3-DKcYh Message=Errors.Project.Grant.AlreadyExists Parent=(ERROR: duplicate key value violates unique constraint \"unique_constraints_pkey\" (SQLSTATE 23505))",
					},
				},
				Success: &admin.ImportDataSuccess{
					Orgs: []*admin.ImportDataSuccessOrg{
						{
							OrgId: orgIDs[4],
						},
						{
							OrgId:      orgIDs[3],
							ProjectIds: projectIDs[2:4],
							ProjectRoles: []string{
								projectIDs[2] + "_role1",
								projectIDs[2] + "_role2",
								projectIDs[3] + "_role3",
								projectIDs[3] + "_role4",
							},
							ProjectGrants: []*admin.ImportDataSuccessProjectGrant{
								{
									GrantId:   grantIDs[4],
									ProjectId: projectIDs[2],
									OrgId:     orgIDs[4],
								},
							},
						},
					},
				},
			},
		},
		{
			name: "duplicate project grant member error",
			req: &admin.ImportDataRequest{
				Data: &admin.ImportDataRequest_DataOrgs{
					DataOrgs: &admin.ImportDataOrg{
						Orgs: []*admin.DataOrg{
							{
								OrgId: orgIDs[6],
								Org: &management.AddOrgRequest{
									Name: integration.OrganizationName(),
								},
							},
							{
								OrgId: orgIDs[5],
								Org: &management.AddOrgRequest{
									Name: integration.OrganizationName(),
								},
								Projects: []*v1.DataProject{
									{
										ProjectId: projectIDs[4],
										Project: &management.AddProjectRequest{
											Name:                 integration.ProjectName(),
											ProjectRoleAssertion: true,
										},
									},
								},
								ProjectRoles: []*management.AddProjectRoleRequest{
									{
										ProjectId:   projectIDs[4],
										RoleKey:     "role1",
										DisplayName: "role1",
									},
									{
										ProjectId:   projectIDs[4],
										RoleKey:     "role2",
										DisplayName: "role2",
									},
								},
								HumanUsers: []*v1.DataHumanUser{
									{
										UserId: userIDs[2],
										User: &management.ImportHumanUserRequest{
											UserName: integration.Username(),
											Profile: &management.ImportHumanUserRequest_Profile{
												FirstName:         integration.FirstName(),
												LastName:          integration.LastName(),
												DisplayName:       integration.Username(),
												PreferredLanguage: integration.Language(),
											},
											Email: &management.ImportHumanUserRequest_Email{
												Email:           integration.Email(),
												IsEmailVerified: true,
											},
										},
									},
								},
								ProjectGrants: []*v1.DataProjectGrant{
									{
										GrantId: grantIDs[5],
										ProjectGrant: &management.AddProjectGrantRequest{
											ProjectId:    projectIDs[4],
											GrantedOrgId: orgIDs[6],
											RoleKeys:     []string{"role1", "role2"},
										},
									},
								},
								ProjectGrantMembers: []*management.AddProjectGrantMemberRequest{
									{
										ProjectId: projectIDs[4],
										GrantId:   grantIDs[5],
										UserId:    userIDs[2],
										Roles:     []string{"PROJECT_GRANT_OWNER"},
									},
									{
										ProjectId: projectIDs[4],
										GrantId:   grantIDs[5],
										UserId:    userIDs[2],
										Roles:     []string{"PROJECT_GRANT_OWNER"},
									},
								},
							},
						},
					},
				},
				Timeout: time.Minute.String(),
			},
			want: &admin.ImportDataResponse{
				Errors: []*admin.ImportDataError{
					{
						Type:    "project_grant_member",
						Id:      orgIDs[5] + "_" + projectIDs[4] + "_" + grantIDs[5] + "_" + userIDs[2],
						Message: "ID=PROJECT-37fug Message=Errors.AlreadyExists",
					},
				},
				Success: &admin.ImportDataSuccess{
					Orgs: []*admin.ImportDataSuccessOrg{
						{
							OrgId: orgIDs[6],
						},
						{
							OrgId:      orgIDs[5],
							ProjectIds: projectIDs[4:5],
							ProjectRoles: []string{
								projectIDs[4] + "_role1",
								projectIDs[4] + "_role2",
							},
							HumanUserIds: userIDs[2:3],
							ProjectGrants: []*admin.ImportDataSuccessProjectGrant{
								{
									GrantId:   grantIDs[5],
									ProjectId: projectIDs[4],
									OrgId:     orgIDs[6],
								},
							},
							ProjectGrantMembers: []*admin.ImportDataSuccessProjectGrantMember{
								{
									ProjectId: projectIDs[4],
									GrantId:   grantIDs[5],
									UserId:    userIDs[2],
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.ImportData(AdminCTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.EqualProto(t, tt.want, got)
		})
	}
}

func generateIDs(n int) []string {
	ids := make([]string, n)
	for i := range ids {
		ids[i] = uuid.NewString()
	}
	return ids
}
