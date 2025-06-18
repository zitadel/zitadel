//go:build integration

package authorization_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	authorization "github.com/zitadel/zitadel/pkg/grpc/authorization/v2beta"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_ListAuthorizations(t *testing.T) {
	iamOwnerCtx := Instance.WithAuthorizationToken(EmptyCTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx context.Context
		dep func(*authorization.ListAuthorizationsRequest, *authorization.ListAuthorizationsResponse)
		req *authorization.ListAuthorizationsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *authorization.ListAuthorizationsResponse
		wantErr bool
	}{
		{
			name: "list by user id, unauthenticated",
			args: args{
				ctx: EmptyCTX,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx)

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), false)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, no permission",
			args: args{
				ctx: Instance.WithAuthorizationToken(EmptyCTX, integration.UserTypeNoPermission),
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx)

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), false)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, missing permission",
			args: args{
				ctx: Instance.WithAuthorizationToken(EmptyCTX, integration.UserTypeOrgOwner),
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx)

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), false)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{},
			},
		},
		{
			name: "list, not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{Filter: &authorization.AuthorizationsSearchFilter_UserId{
							UserId: &filter.IDFilter{
								Id: "notexisting",
							},
						},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list single id, project",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx)

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					response.Authorizations[0] = createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), false)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single id, project grant",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx)

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					response.Authorizations[0] = createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), true)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.dep != nil {
				tt.args.dep(tt.args.req, tt.want)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(iamOwnerCtx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := Instance.Client.AuthorizationV2Beta.ListAuthorizations(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr)
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Authorizations, len(tt.want.Authorizations)) {
					for i := range tt.want.Authorizations {
						assert.EqualExportedValues(ttt, tt.want.Authorizations[i], got.Authorizations[i])
					}
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func assertPaginationResponse(t *assert.CollectT, expected *filter.PaginationResponse, actual *filter.PaginationResponse) {
	assert.Equal(t, expected.AppliedLimit, actual.AppliedLimit)
	assert.Equal(t, expected.TotalResult, actual.TotalResult)
}

func createAuthorization(ctx context.Context, instance *integration.Instance, t *testing.T, orgID, userID string, grant bool) *authorization.Authorization {
	projectName := gofakeit.AppName()
	projectResp := instance.CreateProject(ctx, t, orgID, projectName, false, false)
	userResp, err := instance.Client.UserV2.GetUserByID(ctx, &user.GetUserByIDRequest{UserId: userID})

	if grant {
		grantedOrgName := gofakeit.Company()
		grantedOrg := instance.CreateOrganization(ctx, grantedOrgName, gofakeit.Email())
		instance.CreateProjectGrant(ctx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())

		userGrantResp := instance.CreateProjectGrantUserGrant(ctx, orgID, projectResp.GetId(), grantedOrg.GetOrganizationId(), userID)
		require.NoError(t, err)
		return &authorization.Authorization{
			Id: userGrantResp.GetUserGrantId(),
			Grant: &authorization.Authorization_ProjectGrant{
				ProjectGrant: &authorization.ProjectGrant{
					Id:             grantedOrg.GetOrganizationId(),
					ProjectId:      projectResp.GetId(),
					ProjectName:    projectName,
					OrganizationId: orgID,
				},
			},
			OrganizationId: grantedOrg.GetOrganizationId(),
			CreationDate:   userGrantResp.Details.GetCreationDate(),
			ChangeDate:     userGrantResp.Details.GetCreationDate(),
			State:          1,
			User: &authorization.User{
				Id:                 userID,
				PreferredLoginName: userResp.User.GetPreferredLoginName(),
				DisplayName:        userResp.User.GetHuman().GetProfile().GetDisplayName(),
				AvatarUrl:          userResp.User.GetHuman().GetProfile().GetAvatarUrl(),
				OrganizationId:     userResp.GetUser().GetDetails().GetResourceOwner(),
			},
			Roles: []string{},
		}
	}
	userGrantResp := instance.CreateProjectUserGrant(t, ctx, projectResp.GetId(), userID)
	return &authorization.Authorization{
		Id: userGrantResp.GetUserGrantId(),
		Grant: &authorization.Authorization_Project{
			Project: &authorization.Project{
				Id:             projectResp.GetId(),
				Name:           projectName,
				OrganizationId: orgID,
			},
		},
		OrganizationId: orgID,
		CreationDate:   userGrantResp.Details.GetCreationDate(),
		ChangeDate:     userGrantResp.Details.GetCreationDate(),
		State:          1,
		User: &authorization.User{
			Id:                 userID,
			PreferredLoginName: userResp.User.GetPreferredLoginName(),
			DisplayName:        userResp.User.GetHuman().GetProfile().GetDisplayName(),
			AvatarUrl:          userResp.User.GetHuman().GetProfile().GetAvatarUrl(),
			OrganizationId:     userResp.GetUser().GetDetails().GetResourceOwner(),
		},
	}
}
