//go:build integration

package group_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

func TestServer_ListGroups(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	type args struct {
		ctx context.Context
		req *group_v2.ListGroupsRequest
		dep func(*group_v2.ListGroupsRequest, *group_v2.ListGroupsResponse)
	}
	tests := []struct {
		name    string
		args    args
		want    *group_v2.ListGroupsResponse
		wantErr bool
	}{
		{
			name: "list groups, unauthenticated",
			args: args{
				ctx: CTX,
				req: &group_v2.ListGroupsRequest{},
			},
			wantErr: true,
		},
		{
			name: "group ID not found",
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
						Id:           group1.GetId(),
						Name:         groupName,
						CreationDate: group1.GetCreationDate(),
						ChangeDate:   group1.GetCreationDate(),
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
						Id:           group1.GetId(),
						Name:         groupName1,
						CreationDate: group1.GetCreationDate(),
						ChangeDate:   group1.GetCreationDate(),
					}
					groupName2 := integration.GroupName()
					group2 := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName2)

					resp.Groups[0] = &group_v2.Group{
						Id:           group2.GetId(),
						Name:         groupName2,
						CreationDate: group2.GetCreationDate(),
						ChangeDate:   group2.GetCreationDate(),
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
						Id:           group1.GetId(),
						Name:         groupName,
						CreationDate: group1.GetCreationDate(),
						ChangeDate:   group1.GetCreationDate(),
					}
					req.Filters[0].Filter = &group_v2.GroupsSearchFilter_NameQuery{
						NameQuery: &group_v2.GroupNameQuery{
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
						Id:           group2.GetId(),
						Name:         groupName2,
						CreationDate: group2.GetCreationDate(),
						ChangeDate:   group2.GetCreationDate(),
					}
					groupName1 := integration.GroupName()
					group1 := instance.CreateGroup(iamOwnerCtx, t, org1.GetOrganizationId(), groupName1)

					resp.Groups[1] = &group_v2.Group{
						Id:           group1.GetId(),
						Name:         groupName1,
						CreationDate: group1.GetCreationDate(),
						ChangeDate:   group1.GetCreationDate(),
					}
					groupName0 := integration.GroupName()
					group0 := instance.CreateGroup(iamOwnerCtx, t, org1.GetOrganizationId(), groupName0)

					resp.Groups[0] = &group_v2.Group{
						Id:           group0.GetId(),
						Name:         groupName0,
						CreationDate: group0.GetCreationDate(),
						ChangeDate:   group0.GetCreationDate(),
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
				if tt.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				if assert.Len(t, got.Groups, len(tt.want.Groups)) {
					for i := range got.Groups {
						assert.EqualExportedValues(t, tt.want.Groups[i], got.Groups[i], "want: %v, got: %v", tt.want.Groups[i], got.Groups[i])
					}
				}
				assert.Equal(t, tt.want.Pagination.AppliedLimit, got.Pagination.AppliedLimit)
				assert.Equal(t, tt.want.Pagination.TotalResult, got.Pagination.TotalResult)
			}, retryDuration, tick, "timeout waiting for expected result")
		})
	}
}
