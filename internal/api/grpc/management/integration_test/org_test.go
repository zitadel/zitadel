//go:build integration

package management_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/member"
	"github.com/zitadel/zitadel/pkg/grpc/object"
)

var iamRoles = []string{
	"SELF_MANAGEMENT_GLOBAL",
	"ORG_OWNER",
	"ORG_USER_MANAGER",
	"ORG_OWNER_VIEWER",
	"ORG_SETTINGS_MANAGER",
	"ORG_USER_PERMISSION_EDITOR",
	"ORG_PROJECT_PERMISSION_EDITOR",
	"ORG_PROJECT_CREATOR",
	"ORG_USER_SELF_MANAGER",
	"ORG_ADMIN_IMPERSONATOR",
	"ORG_END_USER_IMPERSONATOR",
}

func TestServer_ListOrgMemberRoles(t *testing.T) {
	got, err := Client.ListOrgMemberRoles(OrgCTX, &mgmt_pb.ListOrgMemberRolesRequest{})
	require.NoError(t, err)
	assert.ElementsMatch(t, iamRoles, got.GetResult())
}

func TestServer_ListOrgMembers(t *testing.T) {
	user := Instance.CreateHumanUserVerified(OrgCTX, Instance.DefaultOrg.Id, integration.Email(), integration.Phone())
	_, err := Client.AddOrgMember(OrgCTX, &mgmt_pb.AddOrgMemberRequest{
		UserId: user.GetUserId(),
		Roles:  iamRoles[1:],
	})
	require.NoError(t, err)
	type args struct {
		ctx context.Context
		req *mgmt_pb.ListOrgMembersRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *mgmt_pb.ListOrgMembersResponse
		wantErr bool
	}{
		{
			name: "permission error",
			args: args{
				ctx: CTX,
				req: &mgmt_pb.ListOrgMembersRequest{
					Query: &object.ListQuery{},
					Queries: []*member.SearchQuery{{
						Query: &member.SearchQuery_UserIdQuery{
							UserIdQuery: &member.UserIDQuery{
								UserId: user.GetUserId(),
							},
						},
					}},
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.ListOrgMembersRequest{
					Query: &object.ListQuery{},
					Queries: []*member.SearchQuery{{
						Query: &member.SearchQuery_UserIdQuery{
							UserIdQuery: &member.UserIDQuery{
								UserId: user.GetUserId(),
							},
						},
					}},
				},
			},
			want: &mgmt_pb.ListOrgMembersResponse{
				Result: []*member.Member{{
					UserId: user.GetUserId(),
					Roles:  iamRoles[1:],
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, time.Minute)
			assert.EventuallyWithT(t, func(ct *assert.CollectT) {
				got, err := Client.ListOrgMembers(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ct, err)
					return
				}
				require.NoError(ct, err)
				wantResult := tt.want.GetResult()
				gotResult := got.GetResult()

				if assert.Len(ct, gotResult, len(wantResult)) {
					for i, want := range wantResult {
						assert.Equal(ct, want.GetUserId(), gotResult[i].GetUserId())
						assert.ElementsMatch(ct, want.GetRoles(), gotResult[i].GetRoles())
					}
				}
			}, retryDuration, tick)
		})
	}
}

func TestServer_AddOrgMember(t *testing.T) {
	user := Instance.CreateHumanUserVerified(OrgCTX, Instance.DefaultOrg.Id, integration.Email(), integration.Phone())
	type args struct {
		ctx context.Context
		req *mgmt_pb.AddOrgMemberRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *mgmt_pb.AddOrgMemberResponse
		wantErr bool
	}{
		{
			name: "permission error",
			args: args{
				ctx: CTX,
				req: &mgmt_pb.AddOrgMemberRequest{
					UserId: user.GetUserId(),
					Roles:  iamRoles,
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.AddOrgMemberRequest{
					UserId: user.GetUserId(),
					Roles:  iamRoles[1:],
				},
			},
			want: &mgmt_pb.AddOrgMemberResponse{
				Details: &object.ObjectDetails{
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "unknown roles error",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.AddOrgMemberRequest{
					UserId: user.GetUserId(),
					Roles:  []string{"FOO", "BAR"},
				},
			},
			wantErr: true,
		},
		{
			name: "iam role error",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.AddOrgMemberRequest{
					UserId: user.GetUserId(),
					Roles:  []string{"IAM_OWNER"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.AddOrgMember(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_UpdateOrgMember(t *testing.T) {
	user := Instance.CreateHumanUserVerified(OrgCTX, Instance.DefaultOrg.Id, integration.Email(), integration.Phone())
	_, err := Client.AddOrgMember(OrgCTX, &mgmt_pb.AddOrgMemberRequest{
		UserId: user.GetUserId(),
		Roles:  []string{"ORG_OWNER"},
	})
	require.NoError(t, err)

	type args struct {
		ctx context.Context
		req *mgmt_pb.UpdateOrgMemberRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *mgmt_pb.UpdateOrgMemberResponse
		wantErr bool
	}{
		{
			name: "permission error",
			args: args{
				ctx: CTX,
				req: &mgmt_pb.UpdateOrgMemberRequest{
					UserId: user.GetUserId(),
					Roles:  iamRoles,
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.UpdateOrgMemberRequest{
					UserId: user.GetUserId(),
					Roles:  iamRoles[1:],
				},
			},
			want: &mgmt_pb.UpdateOrgMemberResponse{
				Details: &object.ObjectDetails{
					ResourceOwner: Instance.DefaultOrg.Id,
					ChangeDate:    timestamppb.Now(),
				},
			},
		},
		{
			name: "unknown roles error",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.UpdateOrgMemberRequest{
					UserId: user.GetUserId(),
					Roles:  []string{"FOO", "BAR"},
				},
			},
			wantErr: true,
		},
		{
			name: "iam role error",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.UpdateOrgMemberRequest{
					UserId: user.GetUserId(),
					Roles:  []string{"IAM_OWNER"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.UpdateOrgMember(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_RemoveIAMMember(t *testing.T) {
	user := Instance.CreateHumanUserVerified(OrgCTX, Instance.DefaultOrg.Id, integration.Email(), integration.Phone())
	_, err := Client.AddOrgMember(OrgCTX, &mgmt_pb.AddOrgMemberRequest{
		UserId: user.GetUserId(),
		Roles:  []string{"ORG_OWNER"},
	})
	require.NoError(t, err)

	type args struct {
		ctx context.Context
		req *mgmt_pb.RemoveOrgMemberRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *mgmt_pb.RemoveOrgMemberResponse
		wantErr bool
	}{
		{
			name: "permission error",
			args: args{
				ctx: CTX,
				req: &mgmt_pb.RemoveOrgMemberRequest{
					UserId: user.GetUserId(),
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: OrgCTX,
				req: &mgmt_pb.RemoveOrgMemberRequest{
					UserId: user.GetUserId(),
				},
			},
			want: &mgmt_pb.RemoveOrgMemberResponse{
				Details: &object.ObjectDetails{
					ResourceOwner: Instance.DefaultOrg.Id,
					ChangeDate:    timestamppb.Now(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.RemoveOrgMember(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}
