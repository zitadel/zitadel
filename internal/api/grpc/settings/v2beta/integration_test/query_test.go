//go:build integration

package settings_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	settings "github.com/zitadel/zitadel/pkg/grpc/settings/v2beta"
)

func TestServer_ListOrganizationSettings(t *testing.T) {
	instance := integration.NewInstance(CTX)
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx context.Context
		dep func(*settings.ListOrganizationSettingsRequest, *settings.ListOrganizationSettingsResponse)
		req *settings.ListOrganizationSettingsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *settings.ListOrganizationSettingsResponse
		wantErr bool
	}{
		{
			name: "list by id, unauthenticated",
			args: args{
				ctx: CTX,
				dep: func(request *settings.ListOrganizationSettingsRequest, response *settings.ListOrganizationSettingsResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp.GetOrganizationId(), true)

					request.Filters[0].Filter = &settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter{
						InOrganizationIdsFilter: &filter.InIDsFilter{
							Ids: []string{orgResp.GetOrganizationId()},
						},
					}
				},
				req: &settings.ListOrganizationSettingsRequest{
					Filters: []*settings.OrganizationSettingsSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, no permission",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(request *settings.ListOrganizationSettingsRequest, response *settings.ListOrganizationSettingsResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp.GetOrganizationId(), true)

					request.Filters[0].Filter = &settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter{
						InOrganizationIdsFilter: &filter.InIDsFilter{
							Ids: []string{orgResp.GetOrganizationId()},
						},
					}
				},
				req: &settings.ListOrganizationSettingsRequest{
					Filters: []*settings.OrganizationSettingsSearchFilter{{}},
				},
			},
			want: &settings.ListOrganizationSettingsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				OrganizationSettings: []*settings.OrganizationSettings{},
			},
		},
		{
			name: "list by id, missing permission",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *settings.ListOrganizationSettingsRequest, response *settings.ListOrganizationSettingsResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp.GetOrganizationId(), true)

					request.Filters[0].Filter = &settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter{
						InOrganizationIdsFilter: &filter.InIDsFilter{
							Ids: []string{orgResp.GetOrganizationId()},
						},
					}
				},
				req: &settings.ListOrganizationSettingsRequest{
					Filters: []*settings.OrganizationSettingsSearchFilter{{}},
				},
			},
			want: &settings.ListOrganizationSettingsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				OrganizationSettings: []*settings.OrganizationSettings{},
			},
		},
		{
			name: "list, not found",
			args: args{
				ctx: iamOwnerCtx,

				req: &settings.ListOrganizationSettingsRequest{
					Filters: []*settings.OrganizationSettingsSearchFilter{{
						Filter: &settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter{
							InOrganizationIdsFilter: &filter.InIDsFilter{
								Ids: []string{"notexisting"},
							},
						},
					}},
				},
			},
			want: &settings.ListOrganizationSettingsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
				OrganizationSettings: []*settings.OrganizationSettings{},
			},
		},
		{
			name: "list single id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *settings.ListOrganizationSettingsRequest, response *settings.ListOrganizationSettingsResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					settingsResp := instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp.GetOrganizationId(), true)

					request.Filters[0].Filter = &settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter{
						InOrganizationIdsFilter: &filter.InIDsFilter{
							Ids: []string{orgResp.GetOrganizationId()},
						},
					}
					response.OrganizationSettings[0] = &settings.OrganizationSettings{
						OrganizationId:              orgResp.GetOrganizationId(),
						CreationDate:                settingsResp.GetSetDate(),
						ChangeDate:                  settingsResp.GetSetDate(),
						OrganizationScopedUsernames: true,
					}
				},
				req: &settings.ListOrganizationSettingsRequest{
					Filters: []*settings.OrganizationSettingsSearchFilter{{}},
				},
			},
			want: &settings.ListOrganizationSettingsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				OrganizationSettings: []*settings.OrganizationSettings{{}},
			},
		},
		{
			name: "list multiple id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *settings.ListOrganizationSettingsRequest, response *settings.ListOrganizationSettingsResponse) {
					orgResp1 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					settingsResp1 := instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp1.GetOrganizationId(), true)
					orgResp2 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					settingsResp2 := instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp2.GetOrganizationId(), true)
					orgResp3 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					settingsResp3 := instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp3.GetOrganizationId(), true)

					request.Filters[0].Filter = &settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter{
						InOrganizationIdsFilter: &filter.InIDsFilter{
							Ids: []string{orgResp1.GetOrganizationId(), orgResp2.GetOrganizationId(), orgResp3.GetOrganizationId()},
						},
					}
					response.OrganizationSettings[2] = &settings.OrganizationSettings{
						OrganizationId:              orgResp1.GetOrganizationId(),
						CreationDate:                settingsResp1.GetSetDate(),
						ChangeDate:                  settingsResp1.GetSetDate(),
						OrganizationScopedUsernames: true,
					}
					response.OrganizationSettings[1] = &settings.OrganizationSettings{
						OrganizationId:              orgResp2.GetOrganizationId(),
						CreationDate:                settingsResp2.GetSetDate(),
						ChangeDate:                  settingsResp2.GetSetDate(),
						OrganizationScopedUsernames: true,
					}
					response.OrganizationSettings[0] = &settings.OrganizationSettings{
						OrganizationId:              orgResp3.GetOrganizationId(),
						CreationDate:                settingsResp3.GetSetDate(),
						ChangeDate:                  settingsResp3.GetSetDate(),
						OrganizationScopedUsernames: true,
					}
				},
				req: &settings.ListOrganizationSettingsRequest{
					Filters: []*settings.OrganizationSettingsSearchFilter{{}},
				},
			},
			want: &settings.ListOrganizationSettingsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				OrganizationSettings: []*settings.OrganizationSettings{{}, {}, {}},
			},
		},
		{
			name: "list multiple id, only org scoped usernames",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *settings.ListOrganizationSettingsRequest, response *settings.ListOrganizationSettingsResponse) {
					orgResp1 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp1.GetOrganizationId(), false)
					orgResp2 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					settingsResp2 := instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp2.GetOrganizationId(), true)
					orgResp3 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					instance.SetOrganizationSettings(iamOwnerCtx, t, orgResp3.GetOrganizationId(), false)

					request.Filters[0].Filter = &settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter{
						InOrganizationIdsFilter: &filter.InIDsFilter{
							Ids: []string{orgResp1.GetOrganizationId(), orgResp2.GetOrganizationId(), orgResp3.GetOrganizationId()},
						},
					}
					request.Filters[1].Filter = &settings.OrganizationSettingsSearchFilter_OrganizationScopedUsernamesFilter{
						OrganizationScopedUsernamesFilter: &settings.OrganizationScopedUsernamesFilter{
							OrganizationScopedUsernames: true,
						},
					}
					response.OrganizationSettings[0] = &settings.OrganizationSettings{
						OrganizationId:              orgResp2.GetOrganizationId(),
						CreationDate:                settingsResp2.GetSetDate(),
						ChangeDate:                  settingsResp2.GetSetDate(),
						OrganizationScopedUsernames: true,
					}
				},
				req: &settings.ListOrganizationSettingsRequest{
					Filters: []*settings.OrganizationSettingsSearchFilter{{}, {}},
				},
			},
			want: &settings.ListOrganizationSettingsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				OrganizationSettings: []*settings.OrganizationSettings{{}},
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
				got, listErr := instance.Client.SettingsV2beta.ListOrganizationSettings(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr)
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.OrganizationSettings, len(tt.want.OrganizationSettings)) {
					for i := range tt.want.OrganizationSettings {
						assert.EqualExportedValues(ttt, tt.want.OrganizationSettings[i], got.OrganizationSettings[i])
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
