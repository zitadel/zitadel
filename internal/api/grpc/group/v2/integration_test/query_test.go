//go:build integration

package group_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	authorization "github.com/zitadel/zitadel/pkg/grpc/authorization/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_GetGroup(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	type args struct {
		ctx context.Context
		req *group_v2.GetGroupRequest
		dep func(*group_v2.GetGroupRequest, *group_v2.GetGroupResponse)
	}
	tests := []struct {
		name        string
		args        args
		want        *group_v2.GetGroupResponse
		wantErrCode codes.Code
		wantErrMsg  string
	}{
		{
			name: "unauthenticated",
			args: args{
				ctx: context.Background(),
				req: &group_v2.GetGroupRequest{},
			},
			wantErrCode: codes.Unauthenticated,
			wantErrMsg:  "auth header missing",
		},
		{
			name: "missing id",
			args: args{
				ctx: iamOwnerCtx,
				req: &group_v2.GetGroupRequest{},
			},
			wantErrCode: codes.InvalidArgument,
			wantErrMsg:  "invalid GetGroupRequest.Id: value length must be between 1 and 200 runes, inclusive",
		},
		{
			name: "get group, instance owner, not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &group_v2.GetGroupRequest{
					Id: "group1",
				},
			},
			wantErrCode: codes.NotFound,
			wantErrMsg:  "Errors.Group.NotFound (QUERY-SG4WbR)",
		},
		{
			name: "get group, missing permission, error",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(req *group_v2.GetGroupRequest, resp *group_v2.GetGroupResponse) {
					groupName := integration.GroupName()
					group := instance.CreateGroup(iamOwnerCtx, t, instance.DefaultOrg.GetId(), groupName)

					req.Id = group.GetId()
				},
				req: &group_v2.GetGroupRequest{},
			},
			wantErrCode: codes.NotFound,
			wantErrMsg:  "membership not found (AUTHZ-cdgFk)",
		},
		{
			name: "get group, organization owner, missing permission, error",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(req *group_v2.GetGroupRequest, resp *group_v2.GetGroupResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupName := integration.GroupName()
					group := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

					req.Id = group.GetId()
				},
				req: &group_v2.GetGroupRequest{},
			},
			wantErrCode: codes.NotFound,
			wantErrMsg:  "membership not found (AUTHZ-cdgFk)",
		},
		{
			name: "get group, organization owner, with permission, found",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(req *group_v2.GetGroupRequest, resp *group_v2.GetGroupResponse) {
					groupName := integration.GroupName()
					group := instance.CreateGroup(iamOwnerCtx, t, instance.DefaultOrg.GetId(), groupName)

					req.Id = group.GetId()
					resp.Group = &group_v2.Group{
						Id:             group.GetId(),
						Name:           groupName,
						Description:    "",
						OrganizationId: instance.DefaultOrg.GetId(),
						ChangeDate:     group.GetCreationDate(),
						CreationDate:   group.GetCreationDate(),
					}
				},
				req: &group_v2.GetGroupRequest{},
			},
			want: &group_v2.GetGroupResponse{
				Group: &group_v2.Group{},
			},
		},
		{
			name: "get group, instance owner, found",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(req *group_v2.GetGroupRequest, resp *group_v2.GetGroupResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupName := integration.GroupName()
					group := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

					req.Id = group.GetId()
					resp.Group = &group_v2.Group{
						Id:             group.GetId(),
						Name:           groupName,
						Description:    "",
						OrganizationId: orgResp.GetOrganizationId(),
						ChangeDate:     group.GetCreationDate(),
						CreationDate:   group.GetCreationDate(),
					}
				},
				req: &group_v2.GetGroupRequest{},
			},
			want: &group_v2.GetGroupResponse{
				Group: &group_v2.Group{},
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
				got, err := instance.Client.GroupV2.GetGroup(tt.args.ctx, tt.args.req)
				if tt.wantErrCode != codes.OK {
					require.Error(ttt, err)
					assert.Equal(ttt, tt.wantErrCode, status.Code(err))
					assert.Equal(ttt, tt.wantErrMsg, status.Convert(err).Message())
					return
				}
				require.NoError(ttt, err)
				assert.EqualExportedValues(t, tt.want.Group, got.Group, "want: %v, got: %v", tt.want.Group, got.Group)
			}, retryDuration, tick, "timeout waiting for expected result")
		})
	}
}

func TestServer_ListGroups(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	type args struct {
		ctx context.Context
		req *group_v2.ListGroupsRequest
		dep func(*group_v2.ListGroupsRequest, *group_v2.ListGroupsResponse)
	}
	tests := []struct {
		name        string
		args        args
		want        *group_v2.ListGroupsResponse
		wantErrCode codes.Code
		wantErrMsg  string
	}{
		{
			name: "list groups, unauthenticated",
			args: args{
				ctx: CTX,
				req: &group_v2.ListGroupsRequest{},
			},
			wantErrCode: codes.Unauthenticated,
			wantErrMsg:  "auth header missing",
		},
		{
			name: "no permission, empty list, TotalResult count returned",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(req *group_v2.ListGroupsRequest, resp *group_v2.ListGroupsResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupName := integration.GroupName()
					group1 := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{group1.GetId()},
						},
					}
				},
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Groups: []*group_v2.Group{},
			},
		},
		{
			name: "org owner, missing permission, empty list, TotalResult count returned",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(req *group_v2.ListGroupsRequest, resp *group_v2.ListGroupsResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupName := integration.GroupName()
					group1 := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{group1.GetId()},
						},
					}
				},
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Groups: []*group_v2.Group{},
			},
		},
		{
			name: "org owner, with permission, ok",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(req *group_v2.ListGroupsRequest, resp *group_v2.ListGroupsResponse) {
					groupName := integration.GroupName()
					group1 := instance.CreateGroup(iamOwnerCtx, t, instance.DefaultOrg.GetId(), groupName)

					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{group1.GetId()},
						},
					}
					resp.Groups[0] = &group_v2.Group{
						Id:             group1.GetId(),
						Name:           groupName,
						Description:    "",
						OrganizationId: instance.DefaultOrg.GetId(),
						CreationDate:   group1.GetCreationDate(),
						ChangeDate:     group1.GetCreationDate(),
					}
				},
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Groups: []*group_v2.Group{
					{},
				},
			},
		},
		{
			name: "instance owner, group ID not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{
						{
							Filter: &group_v2.GroupsSearchFilter_GroupIds{
								GroupIds: &filter.InIDsFilter{
									Ids: []string{"random-group"},
								},
							},
						},
					},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list single group by ID, ok",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(req *group_v2.ListGroupsRequest, resp *group_v2.ListGroupsResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupName := integration.GroupName()
					group1 := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

					resp.Groups[0] = &group_v2.Group{
						Id:             group1.GetId(),
						Name:           groupName,
						Description:    "",
						OrganizationId: orgResp.GetOrganizationId(),
						CreationDate:   group1.GetCreationDate(),
						ChangeDate:     group1.GetCreationDate(),
					}
					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{group1.GetId()},
						},
					}
				},
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Groups: []*group_v2.Group{
					{},
				},
			},
		},
		{
			name: "list multiple groups by IDs, ok",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(req *group_v2.ListGroupsRequest, resp *group_v2.ListGroupsResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupName1 := integration.GroupName()
					group1 := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName1)

					resp.Groups[1] = &group_v2.Group{
						Id:             group1.GetId(),
						Name:           groupName1,
						Description:    "",
						OrganizationId: orgResp.GetOrganizationId(),
						CreationDate:   group1.GetCreationDate(),
						ChangeDate:     group1.GetCreationDate(),
					}
					groupName2 := integration.GroupName()
					group2 := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName2)

					resp.Groups[0] = &group_v2.Group{
						Id:             group2.GetId(),
						Name:           groupName2,
						Description:    "",
						OrganizationId: orgResp.GetOrganizationId(),
						CreationDate:   group2.GetCreationDate(),
						ChangeDate:     group2.GetCreationDate(),
					}
					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{group1.GetId(), group2.GetId()},
						},
					}
				},
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Groups: []*group_v2.Group{
					{}, {},
				},
			},
		},
		{
			name: "list group by name, ok",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(req *group_v2.ListGroupsRequest, resp *group_v2.ListGroupsResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupName := integration.GroupName()
					group1 := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

					resp.Groups[0] = &group_v2.Group{
						Id:             group1.GetId(),
						Name:           groupName,
						Description:    "",
						OrganizationId: orgResp.GetOrganizationId(),
						CreationDate:   group1.GetCreationDate(),
						ChangeDate:     group1.GetCreationDate(),
					}
					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_NameFilter{
						NameFilter: &group_v2.GroupNameFilter{
							Name: groupName,
						},
					}
				},
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Groups: []*group_v2.Group{
					{},
				},
			},
		},
		{
			name: "list by organization ID, ok",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(req *group_v2.ListGroupsRequest, resp *group_v2.ListGroupsResponse) {
					org1 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupName2 := integration.GroupName()
					group2 := instance.CreateGroup(iamOwnerCtx, t, org1.GetOrganizationId(), groupName2)

					resp.Groups[2] = &group_v2.Group{
						Id:             group2.GetId(),
						Name:           groupName2,
						Description:    "",
						OrganizationId: org1.GetOrganizationId(),
						CreationDate:   group2.GetCreationDate(),
						ChangeDate:     group2.GetCreationDate(),
					}
					groupName1 := integration.GroupName()
					group1 := instance.CreateGroup(iamOwnerCtx, t, org1.GetOrganizationId(), groupName1)

					resp.Groups[1] = &group_v2.Group{
						Id:             group1.GetId(),
						Name:           groupName1,
						Description:    "",
						OrganizationId: org1.GetOrganizationId(),
						CreationDate:   group1.GetCreationDate(),
						ChangeDate:     group1.GetCreationDate(),
					}
					groupName0 := integration.GroupName()
					group0 := instance.CreateGroup(iamOwnerCtx, t, org1.GetOrganizationId(), groupName0)

					resp.Groups[0] = &group_v2.Group{
						Id:             group0.GetId(),
						Name:           groupName0,
						Description:    "",
						OrganizationId: org1.GetOrganizationId(),
						CreationDate:   group0.GetCreationDate(),
						ChangeDate:     group0.GetCreationDate(),
					}
					org2 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					org2GroupName0 := integration.GroupName()
					_ = instance.CreateGroup(iamOwnerCtx, t, org2.GetOrganizationId(), org2GroupName0)

					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_OrganizationId{
						OrganizationId: &filter.IDFilter{
							Id: org1.GetOrganizationId(),
						},
					}
				},
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Groups: []*group_v2.Group{
					{}, {}, {},
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
				got, err := instance.Client.GroupV2.ListGroups(tt.args.ctx, tt.args.req)
				if tt.wantErrCode != codes.OK {
					require.Error(ttt, err)
					assert.Equal(ttt, tt.wantErrCode, status.Code(err))
					assert.Equal(ttt, tt.wantErrMsg, status.Convert(err).Message())
					return
				}
				require.NoError(ttt, err)
				if assert.Len(ttt, got.Groups, len(tt.want.Groups)) {
					for i := range got.Groups {
						assert.EqualExportedValues(ttt, tt.want.Groups[i], got.Groups[i], "want: %v, got: %v", tt.want.Groups[i], got.Groups[i])
					}
				}
				assert.Equal(ttt, tt.want.Pagination.AppliedLimit, got.Pagination.AppliedLimit)
				assert.Equal(ttt, tt.want.Pagination.TotalResult, got.Pagination.TotalResult)
			}, retryDuration, tick, "timeout waiting for expected result")
		})
	}
}

func TestServer_ListGroups_WithPermissionV2(t *testing.T) {
	ensureFeaturePermissionV2Enabled(t, instancePermissionV2)
	iamOwnerCtx := instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	type args struct {
		ctx context.Context
		req *group_v2.ListGroupsRequest
		dep func(*group_v2.ListGroupsRequest, *group_v2.ListGroupsResponse)
	}
	tests := []struct {
		name        string
		args        args
		want        *group_v2.ListGroupsResponse
		wantErrCode codes.Code
		wantErrMsg  string
	}{
		{
			name: "list groups, unauthenticated",
			args: args{
				ctx: CTX,
				req: &group_v2.ListGroupsRequest{},
			},
			wantErrCode: codes.Unauthenticated,
			wantErrMsg:  "auth header missing",
		},
		{
			name: "no permission, empty list, TotalResult count set to 0",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(req *group_v2.ListGroupsRequest, resp *group_v2.ListGroupsResponse) {
					orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupName := integration.GroupName()
					group1 := instancePermissionV2.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{group1.GetId()},
						},
					}
				},
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
				Groups: []*group_v2.Group{},
			},
		},
		{
			name: "org owner, missing permission, empty list, TotalResult count set to 0",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(req *group_v2.ListGroupsRequest, resp *group_v2.ListGroupsResponse) {
					orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupName := integration.GroupName()
					group1 := instancePermissionV2.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{group1.GetId()},
						},
					}
				},
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
				Groups: []*group_v2.Group{},
			},
		},
		{
			name: "org owner, with permission, ok",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(req *group_v2.ListGroupsRequest, resp *group_v2.ListGroupsResponse) {
					groupName := integration.GroupName()
					group1 := instancePermissionV2.CreateGroup(iamOwnerCtx, t, instancePermissionV2.DefaultOrg.GetId(), groupName)

					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{group1.GetId()},
						},
					}
					resp.Groups[0] = &group_v2.Group{
						Id:             group1.GetId(),
						Name:           groupName,
						Description:    "",
						OrganizationId: instancePermissionV2.DefaultOrg.GetId(),
						CreationDate:   group1.GetCreationDate(),
						ChangeDate:     group1.GetCreationDate(),
					}
				},
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Groups: []*group_v2.Group{
					{},
				},
			},
		},
		{
			name: "instance owner, group ID not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{
						{
							Filter: &group_v2.GroupsSearchFilter_GroupIds{
								GroupIds: &filter.InIDsFilter{
									Ids: []string{"random-group"},
								},
							},
						},
					},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list single group by ID, ok",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(req *group_v2.ListGroupsRequest, resp *group_v2.ListGroupsResponse) {
					orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupName := integration.GroupName()
					group1 := instancePermissionV2.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

					resp.Groups[0] = &group_v2.Group{
						Id:             group1.GetId(),
						Name:           groupName,
						Description:    "",
						OrganizationId: orgResp.GetOrganizationId(),
						CreationDate:   group1.GetCreationDate(),
						ChangeDate:     group1.GetCreationDate(),
					}
					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{group1.GetId()},
						},
					}
				},
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Groups: []*group_v2.Group{
					{},
				},
			},
		},
		{
			name: "list multiple groups by IDs, ok",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(req *group_v2.ListGroupsRequest, resp *group_v2.ListGroupsResponse) {
					orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupName1 := integration.GroupName()
					group1 := instancePermissionV2.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName1)

					resp.Groups[1] = &group_v2.Group{
						Id:             group1.GetId(),
						Name:           groupName1,
						Description:    "",
						OrganizationId: orgResp.GetOrganizationId(),
						CreationDate:   group1.GetCreationDate(),
						ChangeDate:     group1.GetCreationDate(),
					}
					groupName2 := integration.GroupName()
					group2 := instancePermissionV2.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName2)

					resp.Groups[0] = &group_v2.Group{
						Id:             group2.GetId(),
						Name:           groupName2,
						Description:    "",
						OrganizationId: orgResp.GetOrganizationId(),
						CreationDate:   group2.GetCreationDate(),
						ChangeDate:     group2.GetCreationDate(),
					}
					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{group1.GetId(), group2.GetId()},
						},
					}
				},
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Groups: []*group_v2.Group{
					{}, {},
				},
			},
		},
		{
			name: "list group by name, ok",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(req *group_v2.ListGroupsRequest, resp *group_v2.ListGroupsResponse) {
					orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupName := integration.GroupName()
					group1 := instancePermissionV2.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

					resp.Groups[0] = &group_v2.Group{
						Id:             group1.GetId(),
						Name:           groupName,
						Description:    "",
						OrganizationId: orgResp.GetOrganizationId(),
						CreationDate:   group1.GetCreationDate(),
						ChangeDate:     group1.GetCreationDate(),
					}
					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_NameFilter{
						NameFilter: &group_v2.GroupNameFilter{
							Name: groupName,
						},
					}
				},
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Groups: []*group_v2.Group{
					{},
				},
			},
		},
		{
			name: "list by organization ID, ok",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(req *group_v2.ListGroupsRequest, resp *group_v2.ListGroupsResponse) {
					org1 := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupName2 := integration.GroupName()
					group2 := instancePermissionV2.CreateGroup(iamOwnerCtx, t, org1.GetOrganizationId(), groupName2)

					resp.Groups[2] = &group_v2.Group{
						Id:             group2.GetId(),
						Name:           groupName2,
						Description:    "",
						OrganizationId: org1.GetOrganizationId(),
						CreationDate:   group2.GetCreationDate(),
						ChangeDate:     group2.GetCreationDate(),
					}
					groupName1 := integration.GroupName()
					group1 := instancePermissionV2.CreateGroup(iamOwnerCtx, t, org1.GetOrganizationId(), groupName1)

					resp.Groups[1] = &group_v2.Group{
						Id:             group1.GetId(),
						Name:           groupName1,
						Description:    "",
						OrganizationId: org1.GetOrganizationId(),
						CreationDate:   group1.GetCreationDate(),
						ChangeDate:     group1.GetCreationDate(),
					}
					groupName0 := integration.GroupName()
					group0 := instancePermissionV2.CreateGroup(iamOwnerCtx, t, org1.GetOrganizationId(), groupName0)

					resp.Groups[0] = &group_v2.Group{
						Id:             group0.GetId(),
						Name:           groupName0,
						Description:    "",
						OrganizationId: org1.GetOrganizationId(),
						CreationDate:   group0.GetCreationDate(),
						ChangeDate:     group0.GetCreationDate(),
					}
					org2 := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					org2GroupName0 := integration.GroupName()
					_ = instancePermissionV2.CreateGroup(iamOwnerCtx, t, org2.GetOrganizationId(), org2GroupName0)

					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_OrganizationId{
						OrganizationId: &filter.IDFilter{
							Id: org1.GetOrganizationId(),
						},
					}
				},
				req: &group_v2.ListGroupsRequest{
					Filters: []*group_v2.GroupsSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Groups: []*group_v2.Group{
					{}, {}, {},
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
				got, err := instancePermissionV2.Client.GroupV2.ListGroups(tt.args.ctx, tt.args.req)
				if tt.wantErrCode != codes.OK {
					require.Error(ttt, err)
					assert.Equal(ttt, tt.wantErrCode, status.Code(err))
					assert.Equal(ttt, tt.wantErrMsg, status.Convert(err).Message())
					return
				}
				require.NoError(ttt, err)
				if assert.Len(ttt, got.Groups, len(tt.want.Groups)) {
					for i := range got.Groups {
						assert.EqualExportedValues(ttt, tt.want.Groups[i], got.Groups[i], "want: %v, got: %v", tt.want.Groups[i], got.Groups[i])
					}
				}
				assert.Equal(ttt, tt.want.Pagination.AppliedLimit, got.Pagination.AppliedLimit)
				assert.Equal(ttt, tt.want.Pagination.TotalResult, got.Pagination.TotalResult)
			}, retryDuration, tick, "timeout waiting for expected result")
		})
	}
}

func TestServer_ListGroupUsers(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	type args struct {
		ctx context.Context
		req *group_v2.ListGroupUsersRequest
		dep func(*group_v2.ListGroupUsersRequest, *group_v2.ListGroupUsersResponse)
	}
	tests := []struct {
		name        string
		args        args
		want        *group_v2.ListGroupUsersResponse
		wantErrCode codes.Code
		wantErrMsg  string
	}{
		{
			name: "list groups, unauthenticated",
			args: args{
				ctx: CTX,
				req: &group_v2.ListGroupUsersRequest{},
			},
			wantErrCode: codes.Unauthenticated,
			wantErrMsg:  "auth header missing",
		},
		{
			name: "no permission, empty list, TotalResult count returned",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(req *group_v2.ListGroupUsersRequest, resp *group_v2.ListGroupUsersResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupID, _, _ := addUsersToGroup(iamOwnerCtx, t, instance, orgResp.GetOrganizationId(), 2)

					req.Filters[0].Filter = &group_v2.GroupUsersSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{groupID},
						},
					}
				},
				req: &group_v2.ListGroupUsersRequest{
					Filters: []*group_v2.GroupUsersSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupUsersResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				GroupUsers: []*group_v2.GroupUser{},
			},
		},
		{
			name: "org owner, missing permission, empty list, TotalResult count returned",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(req *group_v2.ListGroupUsersRequest, resp *group_v2.ListGroupUsersResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupID, _, _ := addUsersToGroup(iamOwnerCtx, t, instance, orgResp.GetOrganizationId(), 2)

					req.Filters[0].Filter = &group_v2.GroupUsersSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{groupID},
						},
					}
				},
				req: &group_v2.ListGroupUsersRequest{
					Filters: []*group_v2.GroupUsersSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupUsersResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				GroupUsers: []*group_v2.GroupUser{},
			},
		},
		{
			name: "org owner, with permission, ok",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(req *group_v2.ListGroupUsersRequest, resp *group_v2.ListGroupUsersResponse) {
					groupID, users, addUsersResp := addUsersToGroup(iamOwnerCtx, t, instance, instance.DefaultOrg.GetId(), 2)

					user0, err := instance.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: users[0]})
					require.NoError(t, err)

					user1, err := instance.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: users[1]})
					require.NoError(t, err)

					req.Filters[0].Filter = &group_v2.GroupUsersSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{groupID},
						},
					}
					resp.GroupUsers[0] = &group_v2.GroupUser{
						GroupId:        groupID,
						OrganizationId: instance.DefaultOrg.GetId(),
						User: &authorization.User{
							Id:                 user0.User.GetUserId(),
							OrganizationId:     user0.User.Details.ResourceOwner,
							PreferredLoginName: user0.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp.GetChangeDate(),
					}
					resp.GroupUsers[1] = &group_v2.GroupUser{
						GroupId:        groupID,
						OrganizationId: instance.DefaultOrg.GetId(),
						User: &authorization.User{
							Id:                 user1.User.GetUserId(),
							OrganizationId:     instance.DefaultOrg.GetId(),
							PreferredLoginName: user1.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp.GetChangeDate(),
					}
				},
				req: &group_v2.ListGroupUsersRequest{
					Filters: []*group_v2.GroupUsersSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupUsersResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				GroupUsers: []*group_v2.GroupUser{
					{}, {},
				},
			},
		},
		{
			name: "instance owner, no matching results, ok",
			args: args{
				ctx: iamOwnerCtx,
				req: &group_v2.ListGroupUsersRequest{
					Filters: []*group_v2.GroupUsersSearchFilter{
						{
							Filter: &group_v2.GroupUsersSearchFilter_GroupIds{
								GroupIds: &filter.InIDsFilter{
									Ids: []string{"random-group"},
								},
							},
						},
					},
				},
			},
			want: &group_v2.ListGroupUsersResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list by multiple group IDs, ok",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(req *group_v2.ListGroupUsersRequest, resp *group_v2.ListGroupUsersResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					group0, group0Users, addUsersResp0 := addUsersToGroup(iamOwnerCtx, t, instance, orgResp.GetOrganizationId(), 2)
					group0User0, err := instance.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: group0Users[0]})
					require.NoError(t, err)
					group0User1, err := instance.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: group0Users[1]})
					require.NoError(t, err)

					group1, group1Users, addUsersResp1 := addUsersToGroup(iamOwnerCtx, t, instance, orgResp.GetOrganizationId(), 2)
					group1User0, err := instance.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: group1Users[0]})
					require.NoError(t, err)
					group1User1, err := instance.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: group1Users[1]})
					require.NoError(t, err)

					resp.GroupUsers[0] = &group_v2.GroupUser{
						GroupId:        group1,
						OrganizationId: orgResp.GetOrganizationId(),
						User: &authorization.User{
							Id:                 group1User0.User.GetUserId(),
							OrganizationId:     group1User0.User.Details.ResourceOwner,
							PreferredLoginName: group1User0.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp1.GetChangeDate(),
					}
					resp.GroupUsers[1] = &group_v2.GroupUser{
						GroupId:        group1,
						OrganizationId: orgResp.GetOrganizationId(),
						User: &authorization.User{
							Id:                 group1User1.User.GetUserId(),
							OrganizationId:     group1User1.User.Details.ResourceOwner,
							PreferredLoginName: group1User1.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp1.GetChangeDate(),
					}
					resp.GroupUsers[2] = &group_v2.GroupUser{
						GroupId:        group0,
						OrganizationId: orgResp.GetOrganizationId(),
						User: &authorization.User{
							Id:                 group0User0.User.GetUserId(),
							OrganizationId:     group0User0.User.Details.ResourceOwner,
							PreferredLoginName: group0User0.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp0.GetChangeDate(),
					}
					resp.GroupUsers[3] = &group_v2.GroupUser{
						GroupId:        group0,
						OrganizationId: orgResp.GetOrganizationId(),
						User: &authorization.User{
							Id:                 group0User1.User.GetUserId(),
							OrganizationId:     group0User1.User.Details.ResourceOwner,
							PreferredLoginName: group0User1.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp0.GetChangeDate(),
					}

					req.Filters[0].Filter = &group_v2.GroupUsersSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{group0, group1},
						},
					}
				},
				req: &group_v2.ListGroupUsersRequest{
					Filters: []*group_v2.GroupUsersSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupUsersResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				GroupUsers: []*group_v2.GroupUser{
					{}, {}, {}, {},
				},
			},
		},
		{
			name: "list by multiple user IDs in user groups from different organizations, ok",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(req *group_v2.ListGroupUsersRequest, resp *group_v2.ListGroupUsersResponse) {
					org0 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					group0, group0Users, addUsersResp0 := addUsersToGroup(iamOwnerCtx, t, instance, org0.GetOrganizationId(), 2)
					group0User0, err := instance.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: group0Users[0]})
					require.NoError(t, err)

					org1 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					group1, group1Users, addUsersResp1 := addUsersToGroup(iamOwnerCtx, t, instance, org1.GetOrganizationId(), 2)
					group1User1, err := instance.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: group1Users[1]})
					require.NoError(t, err)

					resp.GroupUsers[0] = &group_v2.GroupUser{
						GroupId:        group1,
						OrganizationId: org1.GetOrganizationId(),
						User: &authorization.User{
							Id:                 group1User1.User.GetUserId(),
							OrganizationId:     group1User1.User.Details.ResourceOwner,
							PreferredLoginName: group1User1.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp1.GetChangeDate(),
					}
					resp.GroupUsers[1] = &group_v2.GroupUser{
						GroupId:        group0,
						OrganizationId: org0.GetOrganizationId(),
						User: &authorization.User{
							Id:                 group0User0.User.GetUserId(),
							OrganizationId:     group0User0.User.Details.ResourceOwner,
							PreferredLoginName: group0User0.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp0.GetChangeDate(),
					}
					req.Filters[0].Filter = &group_v2.GroupUsersSearchFilter_UserIds{
						UserIds: &filter.InIDsFilter{
							Ids: []string{group0User0.User.GetUserId(), group1User1.User.GetUserId()},
						},
					}
				},
				req: &group_v2.ListGroupUsersRequest{
					Filters: []*group_v2.GroupUsersSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupUsersResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				GroupUsers: []*group_v2.GroupUser{
					{}, {},
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
				got, err := instance.Client.GroupV2.ListGroupUsers(tt.args.ctx, tt.args.req)
				if tt.wantErrCode != codes.OK {
					require.Error(ttt, err)
					assert.Equal(ttt, tt.wantErrCode, status.Code(err))
					assert.Equal(ttt, tt.wantErrMsg, status.Convert(err).Message())
					return
				}
				require.NoError(ttt, err)
				if assert.Len(ttt, got.GroupUsers, len(tt.want.GroupUsers)) {
					for i := range got.GroupUsers {
						assert.EqualExportedValues(ttt, tt.want.GroupUsers[i], got.GroupUsers[i], "want: %v, got: %v", tt.want.GroupUsers[i], got.GroupUsers[i])
					}
				}
				assert.Equal(ttt, tt.want.Pagination.AppliedLimit, got.Pagination.AppliedLimit)
				assert.Equal(ttt, tt.want.Pagination.TotalResult, got.Pagination.TotalResult)
			}, retryDuration, tick, "timeout waiting for expected result")
		})
	}
}

func TestServer_ListGroupUsers_WithPermissionV2(t *testing.T) {
	ensureFeaturePermissionV2Enabled(t, instancePermissionV2)
	iamOwnerCtx := instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	type args struct {
		ctx context.Context
		req *group_v2.ListGroupUsersRequest
		dep func(*group_v2.ListGroupUsersRequest, *group_v2.ListGroupUsersResponse)
	}
	tests := []struct {
		name        string
		args        args
		want        *group_v2.ListGroupUsersResponse
		wantErrCode codes.Code
		wantErrMsg  string
	}{
		{
			name: "list groups, unauthenticated",
			args: args{
				ctx: CTX,
				req: &group_v2.ListGroupUsersRequest{},
			},
			wantErrCode: codes.Unauthenticated,
			wantErrMsg:  "auth header missing",
		},
		{
			name: "no permission, empty list, TotalResult count set to 0",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(req *group_v2.ListGroupUsersRequest, resp *group_v2.ListGroupUsersResponse) {
					orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupID, _, _ := addUsersToGroup(iamOwnerCtx, t, instancePermissionV2, orgResp.GetOrganizationId(), 2)

					req.Filters[0].Filter = &group_v2.GroupUsersSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{groupID},
						},
					}
				},
				req: &group_v2.ListGroupUsersRequest{
					Filters: []*group_v2.GroupUsersSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupUsersResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
				GroupUsers: []*group_v2.GroupUser{},
			},
		},
		{
			name: "org owner, missing permission, empty list, TotalResult count set to 0",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(req *group_v2.ListGroupUsersRequest, resp *group_v2.ListGroupUsersResponse) {
					orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					groupID, _, _ := addUsersToGroup(iamOwnerCtx, t, instancePermissionV2, orgResp.GetOrganizationId(), 2)

					req.Filters[0].Filter = &group_v2.GroupUsersSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{groupID},
						},
					}
				},
				req: &group_v2.ListGroupUsersRequest{
					Filters: []*group_v2.GroupUsersSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupUsersResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
				GroupUsers: []*group_v2.GroupUser{},
			},
		},
		{
			name: "org owner, with permission, ok",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(req *group_v2.ListGroupUsersRequest, resp *group_v2.ListGroupUsersResponse) {
					groupID, users, addUsersResp := addUsersToGroup(iamOwnerCtx, t, instancePermissionV2, instancePermissionV2.DefaultOrg.GetId(), 2)

					user0, err := instancePermissionV2.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: users[0]})
					require.NoError(t, err)

					user1, err := instancePermissionV2.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: users[1]})
					require.NoError(t, err)

					req.Filters[0].Filter = &group_v2.GroupUsersSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{groupID},
						},
					}
					resp.GroupUsers[0] = &group_v2.GroupUser{
						GroupId:        groupID,
						OrganizationId: instancePermissionV2.DefaultOrg.GetId(),
						User: &authorization.User{
							Id:                 user0.User.GetUserId(),
							OrganizationId:     user0.User.Details.ResourceOwner,
							PreferredLoginName: user0.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp.GetChangeDate(),
					}
					resp.GroupUsers[1] = &group_v2.GroupUser{
						GroupId:        groupID,
						OrganizationId: instancePermissionV2.DefaultOrg.GetId(),
						User: &authorization.User{
							Id:                 user1.User.GetUserId(),
							OrganizationId:     instancePermissionV2.DefaultOrg.GetId(),
							PreferredLoginName: user1.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp.GetChangeDate(),
					}
				},
				req: &group_v2.ListGroupUsersRequest{
					Filters: []*group_v2.GroupUsersSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupUsersResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				GroupUsers: []*group_v2.GroupUser{
					{}, {},
				},
			},
		},
		{
			name: "instance owner, no matching results, ok",
			args: args{
				ctx: iamOwnerCtx,
				req: &group_v2.ListGroupUsersRequest{
					Filters: []*group_v2.GroupUsersSearchFilter{
						{
							Filter: &group_v2.GroupUsersSearchFilter_GroupIds{
								GroupIds: &filter.InIDsFilter{
									Ids: []string{"random-group"},
								},
							},
						},
					},
				},
			},
			want: &group_v2.ListGroupUsersResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list by multiple group IDs, ok",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(req *group_v2.ListGroupUsersRequest, resp *group_v2.ListGroupUsersResponse) {
					orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					group0, group0Users, addUsersResp0 := addUsersToGroup(iamOwnerCtx, t, instancePermissionV2, orgResp.GetOrganizationId(), 2)
					group0User0, err := instancePermissionV2.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: group0Users[0]})
					require.NoError(t, err)
					group0User1, err := instancePermissionV2.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: group0Users[1]})
					require.NoError(t, err)

					group1, group1Users, addUsersResp1 := addUsersToGroup(iamOwnerCtx, t, instancePermissionV2, orgResp.GetOrganizationId(), 2)
					group1User0, err := instancePermissionV2.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: group1Users[0]})
					require.NoError(t, err)
					group1User1, err := instancePermissionV2.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: group1Users[1]})
					require.NoError(t, err)

					resp.GroupUsers[0] = &group_v2.GroupUser{
						GroupId:        group1,
						OrganizationId: orgResp.GetOrganizationId(),
						User: &authorization.User{
							Id:                 group1User0.User.GetUserId(),
							OrganizationId:     group1User0.User.Details.ResourceOwner,
							PreferredLoginName: group1User0.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp1.GetChangeDate(),
					}
					resp.GroupUsers[1] = &group_v2.GroupUser{
						GroupId:        group1,
						OrganizationId: orgResp.GetOrganizationId(),
						User: &authorization.User{
							Id:                 group1User1.User.GetUserId(),
							OrganizationId:     group1User1.User.Details.ResourceOwner,
							PreferredLoginName: group1User1.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp1.GetChangeDate(),
					}
					resp.GroupUsers[2] = &group_v2.GroupUser{
						GroupId:        group0,
						OrganizationId: orgResp.GetOrganizationId(),
						User: &authorization.User{
							Id:                 group0User0.User.GetUserId(),
							OrganizationId:     group0User0.User.Details.ResourceOwner,
							PreferredLoginName: group0User0.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp0.GetChangeDate(),
					}
					resp.GroupUsers[3] = &group_v2.GroupUser{
						GroupId:        group0,
						OrganizationId: orgResp.GetOrganizationId(),
						User: &authorization.User{
							Id:                 group0User1.User.GetUserId(),
							OrganizationId:     group0User1.User.Details.ResourceOwner,
							PreferredLoginName: group0User1.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp0.GetChangeDate(),
					}

					req.Filters[0].Filter = &group_v2.GroupUsersSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{group0, group1},
						},
					}
				},
				req: &group_v2.ListGroupUsersRequest{
					Filters: []*group_v2.GroupUsersSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupUsersResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				GroupUsers: []*group_v2.GroupUser{
					{}, {}, {}, {},
				},
			},
		},
		{
			name: "list by multiple user IDs in user groups from different organizations with permission in one org, ok",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(req *group_v2.ListGroupUsersRequest, resp *group_v2.ListGroupUsersResponse) {
					group0, group0Users, addUsersResp0 := addUsersToGroup(iamOwnerCtx, t, instancePermissionV2, instancePermissionV2.DefaultOrg.GetId(), 2)
					group0User0, err := instancePermissionV2.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: group0Users[0]})
					require.NoError(t, err)

					org1 := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					_, group1Users, _ := addUsersToGroup(iamOwnerCtx, t, instancePermissionV2, org1.GetOrganizationId(), 2)
					group1User1, err := instancePermissionV2.Client.UserV2.GetUserByID(iamOwnerCtx, &user.GetUserByIDRequest{UserId: group1Users[1]})
					require.NoError(t, err)

					resp.GroupUsers[0] = &group_v2.GroupUser{
						GroupId:        group0,
						OrganizationId: instancePermissionV2.DefaultOrg.GetId(),
						User: &authorization.User{
							Id:                 group0User0.User.GetUserId(),
							OrganizationId:     group0User0.User.Details.ResourceOwner,
							PreferredLoginName: group0User0.User.GetPreferredLoginName(),
							DisplayName:        "Mickey Mouse",
						},
						CreationDate: addUsersResp0.GetChangeDate(),
					}
					req.Filters[0].Filter = &group_v2.GroupUsersSearchFilter_UserIds{
						UserIds: &filter.InIDsFilter{
							Ids: []string{group0User0.User.GetUserId(), group1User1.User.GetUserId()},
						},
					}
				},
				req: &group_v2.ListGroupUsersRequest{
					Filters: []*group_v2.GroupUsersSearchFilter{{}},
				},
			},
			want: &group_v2.ListGroupUsersResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				GroupUsers: []*group_v2.GroupUser{
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
				got, err := instancePermissionV2.Client.GroupV2.ListGroupUsers(tt.args.ctx, tt.args.req)
				if tt.wantErrCode != codes.OK {
					require.Error(ttt, err)
					assert.Equal(ttt, tt.wantErrCode, status.Code(err))
					assert.Equal(ttt, tt.wantErrMsg, status.Convert(err).Message())
					return
				}
				require.NoError(ttt, err)
				if assert.Len(ttt, got.GroupUsers, len(tt.want.GroupUsers)) {
					for i := range got.GroupUsers {
						assert.EqualExportedValues(ttt, tt.want.GroupUsers[i], got.GroupUsers[i], "want: %v, got: %v", tt.want.GroupUsers[i], got.GroupUsers[i])
					}
				}
				assert.Equal(ttt, tt.want.Pagination.AppliedLimit, got.Pagination.AppliedLimit)
				assert.Equal(ttt, tt.want.Pagination.TotalResult, got.Pagination.TotalResult)
			}, retryDuration, tick, "timeout waiting for expected result")
		})
	}
}

func addUsersToGroup(ctx context.Context, t *testing.T, instance *integration.Instance, orgID string, numUsers int) (string, []string, *group_v2.AddUsersToGroupResponse) {
	groupName := integration.GroupName()
	group1 := instance.CreateGroup(ctx, t, orgID, groupName)
	users := make([]string, 0, numUsers)
	for i := 0; i < numUsers; i++ {
		u := instance.CreateHumanUserVerified(ctx, orgID, integration.Email(), integration.Phone())
		users = append(users, u.GetUserId())
	}
	addUsersResp := instance.AddUsersToGroup(ctx, t, group1.GetId(), users)
	return group1.GetId(), users, addUsersResp
}
