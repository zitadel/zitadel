//go:build integration

package action_test

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/integration"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v2beta"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
)

func TestServer_GetTarget(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	type args struct {
		ctx context.Context
		dep func(context.Context, *action.GetTargetRequest, *action.GetTargetResponse) error
		req *action.GetTargetRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *action.GetTargetResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			args: args{
				ctx: instance.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner),
				req: &action.GetTargetRequest{},
			},
			wantErr: true,
		},
		{
			name: "not found",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.GetTargetRequest{Id: "notexisting"},
			},
			wantErr: true,
		},
		{
			name: "get, ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) error {
					name := integration.TargetName()
					resp := instance.CreateTarget(ctx, t, name, "https://example.com", target_domain.TargetTypeWebhook, false)
					request.Id = resp.GetId()
					response.Target.Id = resp.GetId()
					response.Target.Name = name
					response.Target.CreationDate = resp.GetCreationDate()
					response.Target.ChangeDate = resp.GetCreationDate()
					response.Target.SigningKey = resp.GetSigningKey()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.Target{
					Endpoint: "https://example.com",
					TargetType: &action.Target_RestWebhook{
						RestWebhook: &action.RESTWebhook{},
					},
					Timeout: durationpb.New(5 * time.Second),
				},
			},
		},
		{
			name: "get, async, ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) error {
					name := integration.TargetName()
					resp := instance.CreateTarget(ctx, t, name, "https://example.com", target_domain.TargetTypeAsync, false)
					request.Id = resp.GetId()
					response.Target.Id = resp.GetId()
					response.Target.Name = name
					response.Target.CreationDate = resp.GetCreationDate()
					response.Target.ChangeDate = resp.GetCreationDate()
					response.Target.SigningKey = resp.GetSigningKey()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.Target{
					Endpoint: "https://example.com",
					TargetType: &action.Target_RestAsync{
						RestAsync: &action.RESTAsync{},
					},
					Timeout: durationpb.New(5 * time.Second),
				},
			},
		},
		{
			name: "get, webhook interruptOnError, ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) error {
					name := integration.TargetName()
					resp := instance.CreateTarget(ctx, t, name, "https://example.com", target_domain.TargetTypeWebhook, true)
					request.Id = resp.GetId()
					response.Target.Id = resp.GetId()
					response.Target.Name = name
					response.Target.CreationDate = resp.GetCreationDate()
					response.Target.ChangeDate = resp.GetCreationDate()
					response.Target.SigningKey = resp.GetSigningKey()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.Target{
					Endpoint: "https://example.com",
					TargetType: &action.Target_RestWebhook{
						RestWebhook: &action.RESTWebhook{
							InterruptOnError: true,
						},
					},
					Timeout: durationpb.New(5 * time.Second),
				},
			},
		},
		{
			name: "get, call, ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) error {
					name := integration.TargetName()
					resp := instance.CreateTarget(ctx, t, name, "https://example.com", target_domain.TargetTypeCall, false)
					request.Id = resp.GetId()
					response.Target.Id = resp.GetId()
					response.Target.Name = name
					response.Target.CreationDate = resp.GetCreationDate()
					response.Target.ChangeDate = resp.GetCreationDate()
					response.Target.SigningKey = resp.GetSigningKey()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.Target{
					Endpoint: "https://example.com",
					TargetType: &action.Target_RestCall{
						RestCall: &action.RESTCall{
							InterruptOnError: false,
						},
					},
					Timeout: durationpb.New(5 * time.Second),
				},
			},
		},
		{
			name: "get, call interruptOnError, ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) error {
					name := integration.TargetName()
					resp := instance.CreateTarget(ctx, t, name, "https://example.com", target_domain.TargetTypeCall, true)
					request.Id = resp.GetId()
					response.Target.Id = resp.GetId()
					response.Target.Name = name
					response.Target.CreationDate = resp.GetCreationDate()
					response.Target.ChangeDate = resp.GetCreationDate()
					response.Target.SigningKey = resp.GetSigningKey()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.Target{
					Endpoint: "https://example.com",
					TargetType: &action.Target_RestCall{
						RestCall: &action.RESTCall{
							InterruptOnError: true,
						},
					},
					Timeout: durationpb.New(5 * time.Second),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.dep != nil {
				err := tt.args.dep(tt.args.ctx, tt.args.req, tt.want)
				require.NoError(t, err)
			}
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(isolatedIAMOwnerCTX, 2*time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := instance.Client.ActionV2beta.GetTarget(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(ttt, err, "Error: "+err.Error())
					return
				}
				assert.NoError(ttt, err)
				assert.EqualExportedValues(ttt, tt.want, got)
			}, retryDuration, tick, "timeout waiting for expected target Executions")
		})
	}
}

func TestServer_ListTargets(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	type args struct {
		ctx context.Context
		dep func(context.Context, *action.ListTargetsRequest, *action.ListTargetsResponse)
		req *action.ListTargetsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *action.ListTargetsResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			args: args{
				ctx: instance.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner),
				req: &action.ListTargetsRequest{},
			},
			wantErr: true,
		},
		{
			name: "list, not found",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.ListTargetsRequest{
					Filters: []*action.TargetSearchFilter{
						{Filter: &action.TargetSearchFilter_InTargetIdsFilter{
							InTargetIdsFilter: &action.InTargetIDsFilter{
								TargetIds: []string{"notfound"},
							},
						},
						},
					},
				},
			},
			want: &action.ListTargetsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
				Targets: []*action.Target{},
			},
		},
		{
			name: "list single id",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.ListTargetsRequest, response *action.ListTargetsResponse) {
					name := integration.TargetName()
					resp := instance.CreateTarget(ctx, t, name, "https://example.com", target_domain.TargetTypeWebhook, false)
					request.Filters[0].Filter = &action.TargetSearchFilter_InTargetIdsFilter{
						InTargetIdsFilter: &action.InTargetIDsFilter{
							TargetIds: []string{resp.GetId()},
						},
					}

					response.Targets[0].Id = resp.GetId()
					response.Targets[0].Name = name
					response.Targets[0].CreationDate = resp.GetCreationDate()
					response.Targets[0].ChangeDate = resp.GetCreationDate()
					response.Targets[0].SigningKey = resp.GetSigningKey()
				},
				req: &action.ListTargetsRequest{
					Filters: []*action.TargetSearchFilter{{}},
				},
			},
			want: &action.ListTargetsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Targets: []*action.Target{
					{
						Endpoint: "https://example.com",
						TargetType: &action.Target_RestWebhook{
							RestWebhook: &action.RESTWebhook{
								InterruptOnError: false,
							},
						},
						Timeout: durationpb.New(5 * time.Second),
					},
				},
			},
		}, {
			name: "list single name",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.ListTargetsRequest, response *action.ListTargetsResponse) {
					name := integration.TargetName()
					resp := instance.CreateTarget(ctx, t, name, "https://example.com", target_domain.TargetTypeWebhook, false)
					request.Filters[0].Filter = &action.TargetSearchFilter_TargetNameFilter{
						TargetNameFilter: &action.TargetNameFilter{
							TargetName: name,
						},
					}

					response.Targets[0].Id = resp.GetId()
					response.Targets[0].Name = name
					response.Targets[0].CreationDate = resp.GetCreationDate()
					response.Targets[0].ChangeDate = resp.GetCreationDate()
					response.Targets[0].SigningKey = resp.GetSigningKey()
				},
				req: &action.ListTargetsRequest{
					Filters: []*action.TargetSearchFilter{{}},
				},
			},
			want: &action.ListTargetsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Targets: []*action.Target{
					{
						Endpoint: "https://example.com",
						TargetType: &action.Target_RestWebhook{
							RestWebhook: &action.RESTWebhook{
								InterruptOnError: false,
							},
						},
						Timeout: durationpb.New(5 * time.Second),
					},
				},
			},
		},
		{
			name: "list multiple id",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.ListTargetsRequest, response *action.ListTargetsResponse) {
					name1 := integration.TargetName()
					name2 := integration.TargetName()
					name3 := integration.TargetName()
					resp1 := instance.CreateTarget(ctx, t, name1, "https://example.com", target_domain.TargetTypeWebhook, false)
					resp2 := instance.CreateTarget(ctx, t, name2, "https://example.com", target_domain.TargetTypeCall, true)
					resp3 := instance.CreateTarget(ctx, t, name3, "https://example.com", target_domain.TargetTypeAsync, false)
					request.Filters[0].Filter = &action.TargetSearchFilter_InTargetIdsFilter{
						InTargetIdsFilter: &action.InTargetIDsFilter{
							TargetIds: []string{resp1.GetId(), resp2.GetId(), resp3.GetId()},
						},
					}

					response.Targets[2].Id = resp1.GetId()
					response.Targets[2].Name = name1
					response.Targets[2].CreationDate = resp1.GetCreationDate()
					response.Targets[2].ChangeDate = resp1.GetCreationDate()
					response.Targets[2].SigningKey = resp1.GetSigningKey()

					response.Targets[1].Id = resp2.GetId()
					response.Targets[1].Name = name2
					response.Targets[1].CreationDate = resp2.GetCreationDate()
					response.Targets[1].ChangeDate = resp2.GetCreationDate()
					response.Targets[1].SigningKey = resp2.GetSigningKey()

					response.Targets[0].Id = resp3.GetId()
					response.Targets[0].Name = name3
					response.Targets[0].CreationDate = resp3.GetCreationDate()
					response.Targets[0].ChangeDate = resp3.GetCreationDate()
					response.Targets[0].SigningKey = resp3.GetSigningKey()
				},
				req: &action.ListTargetsRequest{
					Filters: []*action.TargetSearchFilter{{}},
				},
			},
			want: &action.ListTargetsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Targets: []*action.Target{
					{
						Endpoint: "https://example.com",
						TargetType: &action.Target_RestAsync{
							RestAsync: &action.RESTAsync{},
						},
						Timeout: durationpb.New(5 * time.Second),
					},
					{
						Endpoint: "https://example.com",
						TargetType: &action.Target_RestCall{
							RestCall: &action.RESTCall{
								InterruptOnError: true,
							},
						},
						Timeout: durationpb.New(5 * time.Second),
					},
					{
						Endpoint: "https://example.com",
						TargetType: &action.Target_RestWebhook{
							RestWebhook: &action.RESTWebhook{
								InterruptOnError: false,
							},
						},
						Timeout: durationpb.New(5 * time.Second),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.dep != nil {
				tt.args.dep(tt.args.ctx, tt.args.req, tt.want)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(isolatedIAMOwnerCTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := instance.Client.ActionV2beta.ListTargets(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr, "Error: "+listErr.Error())
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Targets, len(tt.want.Targets)) {
					for i := range tt.want.Targets {
						assert.EqualExportedValues(ttt, tt.want.Targets[i], got.Targets[i])
					}
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution Executions")
		})
	}
}

func assertPaginationResponse(t *assert.CollectT, expected *filter.PaginationResponse, actual *filter.PaginationResponse) {
	assert.Equal(t, expected.AppliedLimit, actual.AppliedLimit)
	assert.Equal(t, expected.TotalResult, actual.TotalResult)
}

func TestServer_ListExecutions(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	targetResp := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false)

	type args struct {
		ctx context.Context
		dep func(context.Context, *action.ListExecutionsRequest, *action.ListExecutionsResponse)
		req *action.ListExecutionsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *action.ListExecutionsResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			args: args{
				ctx: instance.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner),
				req: &action.ListExecutionsRequest{},
			},
			wantErr: true,
		},
		{
			name: "list request single condition",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.ListExecutionsRequest, response *action.ListExecutionsResponse) {
					cond := request.Filters[0].GetInConditionsFilter().GetConditions()[0]
					resp := setExecution(ctx, t, instance, cond, []string{targetResp.GetId()})

					// Set expected response with used values for SetExecution
					response.Executions[0].CreationDate = resp.GetSetDate()
					response.Executions[0].ChangeDate = resp.GetSetDate()
					response.Executions[0].Condition = cond
				},
				req: &action.ListExecutionsRequest{
					Filters: []*action.ExecutionSearchFilter{{
						Filter: &action.ExecutionSearchFilter_InConditionsFilter{
							InConditionsFilter: &action.InConditionsFilter{
								Conditions: []*action.Condition{{
									ConditionType: &action.Condition_Request{
										Request: &action.RequestExecution{
											Condition: &action.RequestExecution_Method{
												Method: "/zitadel.session.v2.SessionService/GetSession",
											},
										},
									},
								}},
							},
						},
					}},
				},
			},
			want: &action.ListExecutionsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Executions: []*action.Execution{
					{
						Condition: &action.Condition{
							ConditionType: &action.Condition_Request{
								Request: &action.RequestExecution{
									Condition: &action.RequestExecution_Method{
										Method: "/zitadel.session.v2.SessionService/GetSession",
									},
								},
							},
						},
						Targets: []string{targetResp.GetId()},
					},
				},
			},
		},
		{
			name: "list request single target",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.ListExecutionsRequest, response *action.ListExecutionsResponse) {
					target := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false)
					// add target as Filter to the request
					request.Filters[0] = &action.ExecutionSearchFilter{
						Filter: &action.ExecutionSearchFilter_TargetFilter{
							TargetFilter: &action.TargetFilter{
								TargetId: target.GetId(),
							},
						},
					}
					cond := &action.Condition{
						ConditionType: &action.Condition_Request{
							Request: &action.RequestExecution{
								Condition: &action.RequestExecution_Method{
									Method: "/zitadel.management.v1.ManagementService/UpdateAction",
								},
							},
						},
					}
					resp := setExecution(ctx, t, instance, cond, []string{target.GetId()})

					response.Executions[0].CreationDate = resp.GetSetDate()
					response.Executions[0].ChangeDate = resp.GetSetDate()
					response.Executions[0].Condition = cond
					response.Executions[0].Targets = []string{target.GetId()}
				},
				req: &action.ListExecutionsRequest{
					Filters: []*action.ExecutionSearchFilter{{}},
				},
			},
			want: &action.ListExecutionsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Executions: []*action.Execution{
					{
						Condition: &action.Condition{},
						Targets:   []string{""},
					},
				},
			},
		},
		{
			name: "list multiple conditions",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.ListExecutionsRequest, response *action.ListExecutionsResponse) {

					request.Filters[0] = &action.ExecutionSearchFilter{
						Filter: &action.ExecutionSearchFilter_InConditionsFilter{
							InConditionsFilter: &action.InConditionsFilter{
								Conditions: []*action.Condition{
									{ConditionType: &action.Condition_Request{
										Request: &action.RequestExecution{
											Condition: &action.RequestExecution_Method{
												Method: "/zitadel.session.v2.SessionService/GetSession",
											},
										},
									}},
									{ConditionType: &action.Condition_Request{
										Request: &action.RequestExecution{
											Condition: &action.RequestExecution_Method{
												Method: "/zitadel.session.v2.SessionService/CreateSession",
											},
										},
									}},
									{ConditionType: &action.Condition_Request{
										Request: &action.RequestExecution{
											Condition: &action.RequestExecution_Method{
												Method: "/zitadel.session.v2.SessionService/SetSession",
											},
										},
									}},
								},
							},
						},
					}

					cond1 := request.Filters[0].GetInConditionsFilter().GetConditions()[0]
					resp1 := setExecution(ctx, t, instance, cond1, []string{targetResp.GetId()})
					response.Executions[2] = &action.Execution{
						CreationDate: resp1.GetSetDate(),
						ChangeDate:   resp1.GetSetDate(),
						Condition:    cond1,
						Targets:      []string{targetResp.GetId()},
					}

					cond2 := request.Filters[0].GetInConditionsFilter().GetConditions()[1]
					resp2 := setExecution(ctx, t, instance, cond2, []string{targetResp.GetId()})
					response.Executions[1] = &action.Execution{
						CreationDate: resp2.GetSetDate(),
						ChangeDate:   resp2.GetSetDate(),
						Condition:    cond2,
						Targets:      []string{targetResp.GetId()},
					}

					cond3 := request.Filters[0].GetInConditionsFilter().GetConditions()[2]
					resp3 := setExecution(ctx, t, instance, cond3, []string{targetResp.GetId()})
					response.Executions[0] = &action.Execution{
						CreationDate: resp3.GetSetDate(),
						ChangeDate:   resp3.GetSetDate(),
						Condition:    cond3,
						Targets:      []string{targetResp.GetId()},
					}
				},
				req: &action.ListExecutionsRequest{
					Filters: []*action.ExecutionSearchFilter{
						{},
					},
				},
			},
			want: &action.ListExecutionsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Executions: []*action.Execution{
					{}, {}, {},
				},
			},
		},
		{
			name: "list multiple conditions all types",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.ListExecutionsRequest, response *action.ListExecutionsResponse) {
					conditions := request.Filters[0].GetInConditionsFilter().GetConditions()
					for i, cond := range conditions {
						resp := setExecution(ctx, t, instance, cond, []string{targetResp.GetId()})
						response.Executions[(len(conditions)-1)-i] = &action.Execution{
							CreationDate: resp.GetSetDate(),
							ChangeDate:   resp.GetSetDate(),
							Condition:    cond,
							Targets:      []string{targetResp.GetId()},
						}
					}
				},
				req: &action.ListExecutionsRequest{
					Filters: []*action.ExecutionSearchFilter{{
						Filter: &action.ExecutionSearchFilter_InConditionsFilter{
							InConditionsFilter: &action.InConditionsFilter{
								Conditions: []*action.Condition{
									{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_Method{Method: "/zitadel.session.v2.SessionService/GetSession"}}}},
									{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_Service{Service: "zitadel.session.v2.SessionService"}}}},
									{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_All{All: true}}}},
									{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_Method{Method: "/zitadel.session.v2.SessionService/GetSession"}}}},
									{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_Service{Service: "zitadel.session.v2.SessionService"}}}},
									{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_All{All: true}}}},
									{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_Event{Event: "user.added"}}}},
									{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_Group{Group: "user"}}}},
									{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_All{All: true}}}},
									{ConditionType: &action.Condition_Function{Function: &action.FunctionExecution{Name: "presamlresponse"}}},
								},
							},
						},
					}},
				},
			},
			want: &action.ListExecutionsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  10,
					AppliedLimit: 100,
				},
				Executions: []*action.Execution{
					{},
					{},
					{},
					{},
					{},
					{},
					{},
					{},
					{},
					{},
				},
			},
		},
		{
			name: "list multiple conditions all types, sort id",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.ListExecutionsRequest, response *action.ListExecutionsResponse) {
					conditions := request.Filters[0].GetInConditionsFilter().GetConditions()
					for i, cond := range conditions {
						resp := setExecution(ctx, t, instance, cond, []string{targetResp.GetId()})
						response.Executions[i] = &action.Execution{
							CreationDate: resp.GetSetDate(),
							ChangeDate:   resp.GetSetDate(),
							Condition:    cond,
							Targets:      []string{targetResp.GetId()},
						}
					}
				},
				req: &action.ListExecutionsRequest{
					SortingColumn: gu.Ptr(action.ExecutionFieldName_EXECUTION_FIELD_NAME_ID),
					Filters: []*action.ExecutionSearchFilter{{
						Filter: &action.ExecutionSearchFilter_InConditionsFilter{
							InConditionsFilter: &action.InConditionsFilter{
								Conditions: []*action.Condition{
									{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_Method{Method: "/zitadel.session.v2.SessionService/GetSession"}}}},
									{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_Service{Service: "zitadel.session.v2.SessionService"}}}},
									{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_All{All: true}}}},
									{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_Method{Method: "/zitadel.session.v2.SessionService/GetSession"}}}},
									{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_Service{Service: "zitadel.session.v2.SessionService"}}}},
									{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_All{All: true}}}},
									{ConditionType: &action.Condition_Function{Function: &action.FunctionExecution{Name: "presamlresponse"}}},
									{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_Event{Event: "user.added"}}}},
									{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_Group{Group: "user"}}}},
									{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_All{All: true}}}},
								},
							},
						},
					}},
				},
			},
			want: &action.ListExecutionsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  10,
					AppliedLimit: 100,
				},
				Executions: []*action.Execution{
					{},
					{},
					{},
					{},
					{},
					{},
					{},
					{},
					{},
					{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.dep != nil {
				tt.args.dep(tt.args.ctx, tt.args.req, tt.want)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(isolatedIAMOwnerCTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := instance.Client.ActionV2beta.ListExecutions(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr, "Error: "+listErr.Error())
					return
				}
				require.NoError(ttt, listErr)
				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Executions, len(tt.want.Executions)) {
					assert.EqualExportedValues(ttt, got.Executions, tt.want.Executions)
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution Executions")
		})
	}
}
