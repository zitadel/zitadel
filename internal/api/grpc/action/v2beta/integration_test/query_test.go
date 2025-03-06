//go:build integration

package action_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	action "github.com/zitadel/zitadel/pkg/grpc/resources/action/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
)

func TestServer_GetTarget(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
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
				ctx: instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
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
					name := gofakeit.Name()
					resp := instance.CreateTarget(ctx, t, name, "https://example.com", domain.TargetTypeWebhook, false)
					request.Id = resp.GetDetails().GetId()
					response.Target.Config.Name = name
					response.Target.Details = resp.GetDetails()
					response.Target.SigningKey = resp.GetSigningKey()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.GetTarget{
					Details: &resource_object.Details{
						Created: timestamppb.Now(),
						Changed: timestamppb.Now(),
					},
					Config: &action.Target{
						Endpoint: "https://example.com",
						TargetType: &action.Target_RestWebhook{
							RestWebhook: &action.SetRESTWebhook{},
						},
						Timeout: durationpb.New(10 * time.Second),
					},
				},
			},
		},
		{
			name: "get, async, ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) error {
					name := gofakeit.Name()
					resp := instance.CreateTarget(ctx, t, name, "https://example.com", domain.TargetTypeAsync, false)
					request.Id = resp.GetDetails().GetId()
					response.Target.Config.Name = name
					response.Target.Details = resp.GetDetails()
					response.Target.SigningKey = resp.GetSigningKey()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.GetTarget{
					Details: &resource_object.Details{
						Created: timestamppb.Now(),
						Changed: timestamppb.Now(),
					},
					Config: &action.Target{
						Endpoint: "https://example.com",
						TargetType: &action.Target_RestAsync{
							RestAsync: &action.SetRESTAsync{},
						},
						Timeout: durationpb.New(10 * time.Second),
					},
				},
			},
		},
		{
			name: "get, webhook interruptOnError, ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) error {
					name := gofakeit.Name()
					resp := instance.CreateTarget(ctx, t, name, "https://example.com", domain.TargetTypeWebhook, true)
					request.Id = resp.GetDetails().GetId()
					response.Target.Config.Name = name
					response.Target.Details = resp.GetDetails()
					response.Target.SigningKey = resp.GetSigningKey()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.GetTarget{
					Details: &resource_object.Details{
						Created: timestamppb.Now(),
						Changed: timestamppb.Now(),
					},
					Config: &action.Target{
						Endpoint: "https://example.com",
						TargetType: &action.Target_RestWebhook{
							RestWebhook: &action.SetRESTWebhook{
								InterruptOnError: true,
							},
						},
						Timeout: durationpb.New(10 * time.Second),
					},
				},
			},
		},
		{
			name: "get, call, ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) error {
					name := gofakeit.Name()
					resp := instance.CreateTarget(ctx, t, name, "https://example.com", domain.TargetTypeCall, false)
					request.Id = resp.GetDetails().GetId()
					response.Target.Config.Name = name
					response.Target.Details = resp.GetDetails()
					response.Target.SigningKey = resp.GetSigningKey()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.GetTarget{
					Details: &resource_object.Details{
						Created: timestamppb.Now(),
						Changed: timestamppb.Now(),
					},
					Config: &action.Target{
						Endpoint: "https://example.com",
						TargetType: &action.Target_RestCall{
							RestCall: &action.SetRESTCall{
								InterruptOnError: false,
							},
						},
						Timeout: durationpb.New(10 * time.Second),
					},
				},
			},
		},
		{
			name: "get, call interruptOnError, ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) error {
					name := gofakeit.Name()
					resp := instance.CreateTarget(ctx, t, name, "https://example.com", domain.TargetTypeCall, true)
					request.Id = resp.GetDetails().GetId()
					response.Target.Config.Name = name
					response.Target.Details = resp.GetDetails()
					response.Target.SigningKey = resp.GetSigningKey()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.GetTarget{
					Details: &resource_object.Details{
						Created: timestamppb.Now(),
						Changed: timestamppb.Now(),
					},
					Config: &action.Target{
						Endpoint: "https://example.com",
						TargetType: &action.Target_RestCall{
							RestCall: &action.SetRESTCall{
								InterruptOnError: true,
							},
						},
						Timeout: durationpb.New(10 * time.Second),
					},
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
				got, err := instance.Client.ActionV3Alpha.GetTarget(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(ttt, err, "Error: "+err.Error())
					return
				}
				if !assert.NoError(ttt, err) {
					return
				}

				wantTarget := tt.want.GetTarget()
				gotTarget := got.GetTarget()
				integration.AssertResourceDetails(ttt, wantTarget.GetDetails(), gotTarget.GetDetails())
				assert.EqualExportedValues(ttt, wantTarget.GetConfig(), gotTarget.GetConfig())
				assert.Equal(ttt, wantTarget.GetSigningKey(), gotTarget.GetSigningKey())
			}, retryDuration, tick, "timeout waiting for expected target result")
		})
	}
}

func TestServer_ListTargets(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	type args struct {
		ctx context.Context
		dep func(context.Context, *action.SearchTargetsRequest, *action.SearchTargetsResponse) error
		req *action.SearchTargetsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *action.SearchTargetsResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			args: args{
				ctx: instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
				req: &action.SearchTargetsRequest{},
			},
			wantErr: true,
		},
		{
			name: "list, not found",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.SearchTargetsRequest{
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
			want: &action.SearchTargetsResponse{
				Details: &resource_object.ListDetails{
					TotalResult:  0,
					AppliedLimit: 100,
				},
				Result: []*action.GetTarget{},
			},
		},
		{
			name: "list single id",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.SearchTargetsRequest, response *action.SearchTargetsResponse) error {
					name := gofakeit.Name()
					resp := instance.CreateTarget(ctx, t, name, "https://example.com", domain.TargetTypeWebhook, false)
					request.Filters[0].Filter = &action.TargetSearchFilter_InTargetIdsFilter{
						InTargetIdsFilter: &action.InTargetIDsFilter{
							TargetIds: []string{resp.GetDetails().GetId()},
						},
					}
					response.Details.Timestamp = resp.GetDetails().GetChanged()

					response.Result[0].Details = resp.GetDetails()
					response.Result[0].Config.Name = name
					return nil
				},
				req: &action.SearchTargetsRequest{
					Filters: []*action.TargetSearchFilter{{}},
				},
			},
			want: &action.SearchTargetsResponse{
				Details: &resource_object.ListDetails{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Result: []*action.GetTarget{
					{
						Details: &resource_object.Details{
							Created: timestamppb.Now(),
							Changed: timestamppb.Now(),
						},
						Config: &action.Target{
							Endpoint: "https://example.com",
							TargetType: &action.Target_RestWebhook{
								RestWebhook: &action.SetRESTWebhook{
									InterruptOnError: false,
								},
							},
							Timeout: durationpb.New(10 * time.Second),
						},
					},
				},
			},
		}, {
			name: "list single name",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.SearchTargetsRequest, response *action.SearchTargetsResponse) error {
					name := gofakeit.Name()
					resp := instance.CreateTarget(ctx, t, name, "https://example.com", domain.TargetTypeWebhook, false)
					request.Filters[0].Filter = &action.TargetSearchFilter_TargetNameFilter{
						TargetNameFilter: &action.TargetNameFilter{
							TargetName: name,
						},
					}
					response.Details.Timestamp = resp.GetDetails().GetChanged()

					response.Result[0].Details = resp.GetDetails()
					response.Result[0].Config.Name = name
					return nil
				},
				req: &action.SearchTargetsRequest{
					Filters: []*action.TargetSearchFilter{{}},
				},
			},
			want: &action.SearchTargetsResponse{
				Details: &resource_object.ListDetails{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Result: []*action.GetTarget{
					{
						Details: &resource_object.Details{
							Created: timestamppb.Now(),
							Changed: timestamppb.Now(),
							Owner: &object.Owner{
								Type: object.OwnerType_OWNER_TYPE_INSTANCE,
								Id:   instance.ID(),
							},
						},
						Config: &action.Target{
							Endpoint: "https://example.com",
							TargetType: &action.Target_RestWebhook{
								RestWebhook: &action.SetRESTWebhook{
									InterruptOnError: false,
								},
							},
							Timeout: durationpb.New(10 * time.Second),
						},
					},
				},
			},
		},
		{
			name: "list multiple id",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.SearchTargetsRequest, response *action.SearchTargetsResponse) error {
					name1 := gofakeit.Name()
					name2 := gofakeit.Name()
					name3 := gofakeit.Name()
					resp1 := instance.CreateTarget(ctx, t, name1, "https://example.com", domain.TargetTypeWebhook, false)
					resp2 := instance.CreateTarget(ctx, t, name2, "https://example.com", domain.TargetTypeCall, true)
					resp3 := instance.CreateTarget(ctx, t, name3, "https://example.com", domain.TargetTypeAsync, false)
					request.Filters[0].Filter = &action.TargetSearchFilter_InTargetIdsFilter{
						InTargetIdsFilter: &action.InTargetIDsFilter{
							TargetIds: []string{resp1.GetDetails().GetId(), resp2.GetDetails().GetId(), resp3.GetDetails().GetId()},
						},
					}
					response.Details.Timestamp = resp3.GetDetails().GetChanged()

					response.Result[0].Details = resp1.GetDetails()
					response.Result[0].Config.Name = name1
					response.Result[1].Details = resp2.GetDetails()
					response.Result[1].Config.Name = name2
					response.Result[2].Details = resp3.GetDetails()
					response.Result[2].Config.Name = name3
					return nil
				},
				req: &action.SearchTargetsRequest{
					Filters: []*action.TargetSearchFilter{{}},
				},
			},
			want: &action.SearchTargetsResponse{
				Details: &resource_object.ListDetails{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Result: []*action.GetTarget{
					{
						Details: &resource_object.Details{
							Created: timestamppb.Now(),
							Changed: timestamppb.Now(),
							Owner: &object.Owner{
								Type: object.OwnerType_OWNER_TYPE_INSTANCE,
								Id:   instance.ID(),
							},
						},
						Config: &action.Target{
							Endpoint: "https://example.com",
							TargetType: &action.Target_RestWebhook{
								RestWebhook: &action.SetRESTWebhook{
									InterruptOnError: false,
								},
							},
							Timeout: durationpb.New(10 * time.Second),
						},
					},
					{
						Details: &resource_object.Details{
							Created: timestamppb.Now(),
							Changed: timestamppb.Now(),
							Owner: &object.Owner{
								Type: object.OwnerType_OWNER_TYPE_INSTANCE,
								Id:   instance.ID(),
							},
						},
						Config: &action.Target{
							Endpoint: "https://example.com",
							TargetType: &action.Target_RestCall{
								RestCall: &action.SetRESTCall{
									InterruptOnError: true,
								},
							},
							Timeout: durationpb.New(10 * time.Second),
						},
					},
					{
						Details: &resource_object.Details{
							Created: timestamppb.Now(),
							Changed: timestamppb.Now(),
							Owner: &object.Owner{
								Type: object.OwnerType_OWNER_TYPE_INSTANCE,
								Id:   instance.ID(),
							},
						},
						Config: &action.Target{
							Endpoint: "https://example.com",
							TargetType: &action.Target_RestAsync{
								RestAsync: &action.SetRESTAsync{},
							},
							Timeout: durationpb.New(10 * time.Second),
						},
					},
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

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(isolatedIAMOwnerCTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := instance.Client.ActionV3Alpha.SearchTargets(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr, "Error: "+listErr.Error())
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Result, len(tt.want.Result)) {
					for i := range tt.want.Result {
						integration.AssertResourceDetails(ttt, tt.want.Result[i].GetDetails(), got.Result[i].GetDetails())
						assert.EqualExportedValues(ttt, tt.want.Result[i].GetConfig(), got.Result[i].GetConfig())
						assert.NotEmpty(ttt, got.Result[i].GetSigningKey())
					}
				}
				integration.AssertResourceListDetails(ttt, tt.want, got)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func TestServer_SearchExecutions(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	targetResp := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false)

	type args struct {
		ctx context.Context
		dep func(context.Context, *action.SearchExecutionsRequest, *action.SearchExecutionsResponse) error
		req *action.SearchExecutionsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *action.SearchExecutionsResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			args: args{
				ctx: instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
				req: &action.SearchExecutionsRequest{},
			},
			wantErr: true,
		},
		{
			name: "list request single condition",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.SearchExecutionsRequest, response *action.SearchExecutionsResponse) error {
					cond := request.Filters[0].GetInConditionsFilter().GetConditions()[0]
					resp := instance.SetExecution(ctx, t, cond, executionTargetsSingleTarget(targetResp.GetDetails().GetId()))

					response.Details.Timestamp = resp.GetDetails().GetChanged()
					// Set expected response with used values for SetExecution
					response.Result[0].Details = resp.GetDetails()
					response.Result[0].Condition = cond
					return nil
				},
				req: &action.SearchExecutionsRequest{
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
			want: &action.SearchExecutionsResponse{
				Details: &resource_object.ListDetails{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Result: []*action.GetExecution{
					{
						Details: &resource_object.Details{
							Created: timestamppb.Now(),
							Changed: timestamppb.Now(),
						},
						Condition: &action.Condition{
							ConditionType: &action.Condition_Request{
								Request: &action.RequestExecution{
									Condition: &action.RequestExecution_Method{
										Method: "/zitadel.session.v2.SessionService/GetSession",
									},
								},
							},
						},
						Execution: &action.Execution{
							Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
						},
					},
				},
			},
		},
		{
			name: "list request single target",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.SearchExecutionsRequest, response *action.SearchExecutionsResponse) error {
					target := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false)
					// add target as Filter to the request
					request.Filters[0] = &action.ExecutionSearchFilter{
						Filter: &action.ExecutionSearchFilter_TargetFilter{
							TargetFilter: &action.TargetFilter{
								TargetId: target.GetDetails().GetId(),
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
					targets := executionTargetsSingleTarget(target.GetDetails().GetId())
					resp := instance.SetExecution(ctx, t, cond, targets)

					response.Details.Timestamp = resp.GetDetails().GetChanged()

					response.Result[0].Details = resp.GetDetails()
					response.Result[0].Condition = cond
					response.Result[0].Execution.Targets = targets
					return nil
				},
				req: &action.SearchExecutionsRequest{
					Filters: []*action.ExecutionSearchFilter{{}},
				},
			},
			want: &action.SearchExecutionsResponse{
				Details: &resource_object.ListDetails{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Result: []*action.GetExecution{
					{
						Details: &resource_object.Details{
							Created: timestamppb.Now(),
							Changed: timestamppb.Now(),
						},
						Condition: &action.Condition{},
						Execution: &action.Execution{
							Targets: executionTargetsSingleTarget(""),
						},
					},
				},
			},
		}, {
			name: "list request single include",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.SearchExecutionsRequest, response *action.SearchExecutionsResponse) error {
					cond := &action.Condition{
						ConditionType: &action.Condition_Request{
							Request: &action.RequestExecution{
								Condition: &action.RequestExecution_Method{
									Method: "/zitadel.management.v1.ManagementService/GetAction",
								},
							},
						},
					}
					instance.SetExecution(ctx, t, cond, executionTargetsSingleTarget(targetResp.GetDetails().GetId()))
					request.Filters[0].GetIncludeFilter().Include = cond

					includeCond := &action.Condition{
						ConditionType: &action.Condition_Request{
							Request: &action.RequestExecution{
								Condition: &action.RequestExecution_Method{
									Method: "/zitadel.management.v1.ManagementService/ListActions",
								},
							},
						},
					}
					includeTargets := executionTargetsSingleInclude(cond)
					resp2 := instance.SetExecution(ctx, t, includeCond, includeTargets)

					response.Details.Timestamp = resp2.GetDetails().GetChanged()

					response.Result[0].Details = resp2.GetDetails()
					response.Result[0].Condition = includeCond
					response.Result[0].Execution = &action.Execution{
						Targets: includeTargets,
					}
					return nil
				},
				req: &action.SearchExecutionsRequest{
					Filters: []*action.ExecutionSearchFilter{{
						Filter: &action.ExecutionSearchFilter_IncludeFilter{
							IncludeFilter: &action.IncludeFilter{},
						},
					}},
				},
			},
			want: &action.SearchExecutionsResponse{
				Details: &resource_object.ListDetails{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Result: []*action.GetExecution{
					{
						Details: &resource_object.Details{
							Created: timestamppb.Now(),
							Changed: timestamppb.Now(),
						},
					},
				},
			},
		},
		{
			name: "list multiple conditions",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.SearchExecutionsRequest, response *action.SearchExecutionsResponse) error {

					cond1 := request.Filters[0].GetInConditionsFilter().GetConditions()[0]
					targets1 := executionTargetsSingleTarget(targetResp.GetDetails().GetId())
					resp1 := instance.SetExecution(ctx, t, cond1, targets1)
					response.Result[0].Details = resp1.GetDetails()
					response.Result[0].Condition = cond1
					response.Result[0].Execution = &action.Execution{
						Targets: targets1,
					}

					cond2 := request.Filters[0].GetInConditionsFilter().GetConditions()[1]
					targets2 := executionTargetsSingleTarget(targetResp.GetDetails().GetId())
					resp2 := instance.SetExecution(ctx, t, cond2, targets2)
					response.Result[1].Details = resp2.GetDetails()
					response.Result[1].Condition = cond2
					response.Result[1].Execution = &action.Execution{
						Targets: targets2,
					}

					cond3 := request.Filters[0].GetInConditionsFilter().GetConditions()[2]
					targets3 := executionTargetsSingleTarget(targetResp.GetDetails().GetId())
					resp3 := instance.SetExecution(ctx, t, cond3, targets3)
					response.Result[2].Details = resp3.GetDetails()
					response.Result[2].Condition = cond3
					response.Result[2].Execution = &action.Execution{
						Targets: targets3,
					}
					response.Details.Timestamp = resp3.GetDetails().GetChanged()
					return nil
				},
				req: &action.SearchExecutionsRequest{
					Filters: []*action.ExecutionSearchFilter{{
						Filter: &action.ExecutionSearchFilter_InConditionsFilter{
							InConditionsFilter: &action.InConditionsFilter{
								Conditions: []*action.Condition{
									{
										ConditionType: &action.Condition_Request{
											Request: &action.RequestExecution{
												Condition: &action.RequestExecution_Method{
													Method: "/zitadel.session.v2.SessionService/GetSession",
												},
											},
										},
									},
									{
										ConditionType: &action.Condition_Request{
											Request: &action.RequestExecution{
												Condition: &action.RequestExecution_Method{
													Method: "/zitadel.session.v2.SessionService/CreateSession",
												},
											},
										},
									},
									{
										ConditionType: &action.Condition_Request{
											Request: &action.RequestExecution{
												Condition: &action.RequestExecution_Method{
													Method: "/zitadel.session.v2.SessionService/SetSession",
												},
											},
										},
									},
								},
							},
						},
					}},
				},
			},
			want: &action.SearchExecutionsResponse{
				Details: &resource_object.ListDetails{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Result: []*action.GetExecution{
					{
						Details: &resource_object.Details{
							Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instance.ID()},
						},
					}, {
						Details: &resource_object.Details{
							Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instance.ID()},
						},
					}, {
						Details: &resource_object.Details{
							Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instance.ID()},
						},
					},
				},
			},
		},
		{
			name: "list multiple conditions all types",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				dep: func(ctx context.Context, request *action.SearchExecutionsRequest, response *action.SearchExecutionsResponse) error {
					targets := executionTargetsSingleTarget(targetResp.GetDetails().GetId())
					for i, cond := range request.Filters[0].GetInConditionsFilter().GetConditions() {
						resp := instance.SetExecution(ctx, t, cond, targets)
						response.Result[i].Details = resp.GetDetails()
						response.Result[i].Condition = cond
						response.Result[i].Execution = &action.Execution{
							Targets: targets,
						}
						// filled with info of last sequence
						response.Details.Timestamp = resp.GetDetails().GetChanged()
					}

					return nil
				},
				req: &action.SearchExecutionsRequest{
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
			want: &action.SearchExecutionsResponse{
				Details: &resource_object.ListDetails{
					TotalResult:  10,
					AppliedLimit: 100,
				},
				Result: []*action.GetExecution{
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instance.ID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instance.ID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instance.ID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instance.ID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instance.ID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instance.ID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instance.ID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instance.ID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instance.ID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instance.ID()}}},
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

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(isolatedIAMOwnerCTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := instance.Client.ActionV3Alpha.SearchExecutions(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr, "Error: "+listErr.Error())
					return
				}
				require.NoError(ttt, listErr)
				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Result, len(tt.want.Result)) {
					for i := range tt.want.Result {
						// as not sorted, all elements have to be checked
						// workaround as oneof elements can only be checked with assert.EqualExportedValues()
						if j, found := containExecution(got.Result, tt.want.Result[i]); found {
							integration.AssertResourceDetails(ttt, tt.want.Result[i].GetDetails(), got.Result[j].GetDetails())
							got.Result[j].Details = tt.want.Result[i].GetDetails()
							assert.EqualExportedValues(ttt, tt.want.Result[i], got.Result[j])
						}
					}
				}
				integration.AssertResourceListDetails(ttt, tt.want, got)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func containExecution(executionList []*action.GetExecution, execution *action.GetExecution) (int, bool) {
	for i, exec := range executionList {
		if reflect.DeepEqual(exec.Details, execution.Details) {
			return i, true
		}
	}
	return 0, false
}
