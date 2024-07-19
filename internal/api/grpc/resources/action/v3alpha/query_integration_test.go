//go:build integration

package action_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	action "github.com/zitadel/zitadel/pkg/grpc/resources/action/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
)

func TestServer_GetTargetByID(t *testing.T) {
	ensureFeatureEnabled(t)
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
				ctx: Tester.WithAuthorization(context.Background(), integration.OrgOwner),
				req: &action.GetTargetRequest{},
			},
			wantErr: true,
		},
		{
			name: "not found",
			args: args{
				ctx: CTX,
				req: &action.GetTargetRequest{Id: "notexisting"},
			},
			wantErr: true,
		},
		{
			name: "get, ok",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) error {
					name := fmt.Sprint(time.Now().UnixNano() + 1)
					resp := Tester.CreateTarget(ctx, t, name, "https://example.com", domain.TargetTypeWebhook, false)
					request.Id = resp.GetDetails().GetId()
					response.Target.Target.Name = name
					response.Target.Details.Id = resp.GetDetails().GetId()
					response.Target.Details.Owner = resp.GetDetails().GetOwner()
					response.Target.Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Target.Details.Sequence = resp.GetDetails().GetSequence()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.GetTarget{
					Details: &resource_object.Details{
						Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
					},
					Target: &action.Target{
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
				ctx: CTX,
				dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) error {
					name := fmt.Sprint(time.Now().UnixNano() + 1)
					resp := Tester.CreateTarget(ctx, t, name, "https://example.com", domain.TargetTypeAsync, false)
					request.Id = resp.GetDetails().GetId()
					response.Target.Details.Id = resp.GetDetails().GetId()
					response.Target.Target.Name = name
					response.Target.Details.Owner = resp.GetDetails().GetOwner()
					response.Target.Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Target.Details.Sequence = resp.GetDetails().GetSequence()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.GetTarget{
					Details: &resource_object.Details{
						Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
					},
					Target: &action.Target{
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
				ctx: CTX,
				dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) error {
					name := fmt.Sprint(time.Now().UnixNano() + 1)
					resp := Tester.CreateTarget(ctx, t, name, "https://example.com", domain.TargetTypeWebhook, true)
					request.Id = resp.GetDetails().GetId()
					response.Target.Details.Id = resp.GetDetails().GetId()
					response.Target.Target.Name = name
					response.Target.Details.Owner = resp.GetDetails().GetOwner()
					response.Target.Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Target.Details.Sequence = resp.GetDetails().GetSequence()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.GetTarget{
					Details: &resource_object.Details{
						Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
					},
					Target: &action.Target{
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
				ctx: CTX,
				dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) error {
					name := fmt.Sprint(time.Now().UnixNano() + 1)
					resp := Tester.CreateTarget(ctx, t, name, "https://example.com", domain.TargetTypeCall, false)
					request.Id = resp.GetDetails().GetId()
					response.Target.Details.Id = resp.GetDetails().GetId()
					response.Target.Target.Name = name
					response.Target.Details.Owner = resp.GetDetails().GetOwner()
					response.Target.Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Target.Details.Sequence = resp.GetDetails().GetSequence()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.GetTarget{
					Details: &resource_object.Details{
						Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
					},
					Target: &action.Target{
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
				ctx: CTX,
				dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) error {
					name := fmt.Sprint(time.Now().UnixNano() + 1)
					resp := Tester.CreateTarget(ctx, t, name, "https://example.com", domain.TargetTypeCall, true)
					request.Id = resp.GetDetails().GetId()
					response.Target.Details.Id = resp.GetDetails().GetId()
					response.Target.Target.Name = name
					response.Target.Details.Owner = resp.GetDetails().GetOwner()
					response.Target.Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Target.Details.Sequence = resp.GetDetails().GetSequence()
					return nil
				},
				req: &action.GetTargetRequest{},
			},
			want: &action.GetTargetResponse{
				Target: &action.GetTarget{
					Details: &resource_object.Details{
						Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
					},
					Target: &action.Target{
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

			retryDuration := 5 * time.Second
			if ctxDeadline, ok := CTX.Deadline(); ok {
				retryDuration = time.Until(ctxDeadline)
			}

			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, getErr := Client.GetTarget(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(ttt, getErr, "Error: "+getErr.Error())
				} else {
					assert.NoError(ttt, getErr)

					integration.AssertResourceDetails(t, tt.want.GetTarget().GetDetails(), got.GetTarget().GetDetails())

					assert.Equal(t, tt.want.Target, got.Target)
				}

			}, retryDuration, time.Millisecond*100, "timeout waiting for expected execution result")
		})
	}
}

func TestServer_ListTargets(t *testing.T) {
	ensureFeatureEnabled(t)
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
				ctx: Tester.WithAuthorization(context.Background(), integration.OrgOwner),
				req: &action.SearchTargetsRequest{},
			},
			wantErr: true,
		},
		{
			name: "list, not found",
			args: args{
				ctx: CTX,
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
					TotalResult: 0,
				},
				Result: []*action.GetTarget{},
			},
		},
		{
			name: "list single id",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.SearchTargetsRequest, response *action.SearchTargetsResponse) error {
					name := fmt.Sprint(time.Now().UnixNano() + 1)
					resp := Tester.CreateTarget(ctx, t, name, "https://example.com", domain.TargetTypeWebhook, false)
					request.Filters[0].Filter = &action.TargetSearchFilter_InTargetIdsFilter{
						InTargetIdsFilter: &action.InTargetIDsFilter{
							TargetIds: []string{resp.GetDetails().GetId()},
						},
					}
					response.Details.Timestamp = resp.GetDetails().GetChangeDate()
					//response.Details.ProcessedSequence = resp.GetDetails().GetSequence()

					response.Result[0].Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Result[0].Details.Sequence = resp.GetDetails().GetSequence()
					response.Result[0].Details.Id = resp.GetDetails().GetId()
					response.Result[0].Target.Name = name
					return nil
				},
				req: &action.SearchTargetsRequest{
					Filters: []*action.TargetSearchFilter{{}},
				},
			},
			want: &action.SearchTargetsResponse{
				Details: &resource_object.ListDetails{
					TotalResult: 1,
				},
				Result: []*action.GetTarget{
					{
						Details: &resource_object.Details{
							Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
						},
						Target: &action.Target{
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
				ctx: CTX,
				dep: func(ctx context.Context, request *action.SearchTargetsRequest, response *action.SearchTargetsResponse) error {
					name := fmt.Sprint(time.Now().UnixNano() + 1)
					resp := Tester.CreateTarget(ctx, t, name, "https://example.com", domain.TargetTypeWebhook, false)
					request.Filters[0].Filter = &action.TargetSearchFilter_TargetNameFilter{
						TargetNameFilter: &action.TargetNameFilter{
							TargetName: name,
						},
					}
					response.Details.Timestamp = resp.GetDetails().GetChangeDate()
					response.Details.ProcessedSequence = resp.GetDetails().GetSequence()

					response.Result[0].Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Result[0].Details.Sequence = resp.GetDetails().GetSequence()
					response.Result[0].Details.Id = resp.GetDetails().GetId()
					response.Result[0].Target.Name = name
					return nil
				},
				req: &action.SearchTargetsRequest{
					Filters: []*action.TargetSearchFilter{{}},
				},
			},
			want: &action.SearchTargetsResponse{
				Details: &resource_object.ListDetails{
					TotalResult: 1,
				},
				Result: []*action.GetTarget{
					{
						Details: &resource_object.Details{
							Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
						},
						Target: &action.Target{
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
				ctx: CTX,
				dep: func(ctx context.Context, request *action.SearchTargetsRequest, response *action.SearchTargetsResponse) error {
					name1 := fmt.Sprint(time.Now().UnixNano() + 1)
					name2 := fmt.Sprint(time.Now().UnixNano() + 3)
					name3 := fmt.Sprint(time.Now().UnixNano() + 5)
					resp1 := Tester.CreateTarget(ctx, t, name1, "https://example.com", domain.TargetTypeWebhook, false)
					resp2 := Tester.CreateTarget(ctx, t, name2, "https://example.com", domain.TargetTypeCall, true)
					resp3 := Tester.CreateTarget(ctx, t, name3, "https://example.com", domain.TargetTypeAsync, false)
					request.Filters[0].Filter = &action.TargetSearchFilter_InTargetIdsFilter{
						InTargetIdsFilter: &action.InTargetIDsFilter{
							TargetIds: []string{resp1.GetDetails().GetId(), resp2.GetDetails().GetId(), resp3.GetDetails().GetId()},
						},
					}
					response.Details.Timestamp = resp3.GetDetails().GetChangeDate()
					response.Details.ProcessedSequence = resp3.GetDetails().GetSequence()

					response.Result[0].Details.ChangeDate = resp1.GetDetails().GetChangeDate()
					response.Result[0].Details.Sequence = resp1.GetDetails().GetSequence()
					response.Result[0].Details.Id = resp1.GetDetails().GetId()
					response.Result[0].Target.Name = name1
					response.Result[1].Details.ChangeDate = resp2.GetDetails().GetChangeDate()
					response.Result[1].Details.Sequence = resp2.GetDetails().GetSequence()
					response.Result[1].Details.Id = resp2.GetDetails().GetId()
					response.Result[1].Target.Name = name2
					response.Result[2].Details.ChangeDate = resp3.GetDetails().GetChangeDate()
					response.Result[2].Details.Sequence = resp3.GetDetails().GetSequence()
					response.Result[2].Details.Id = resp3.GetDetails().GetId()
					response.Result[2].Target.Name = name3
					return nil
				},
				req: &action.SearchTargetsRequest{
					Filters: []*action.TargetSearchFilter{{}},
				},
			},
			want: &action.SearchTargetsResponse{
				Details: &resource_object.ListDetails{
					TotalResult: 3,
				},
				Result: []*action.GetTarget{
					{
						Details: &resource_object.Details{
							Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
						},
						Target: &action.Target{
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
							Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
						},
						Target: &action.Target{
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
							Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
						},
						Target: &action.Target{
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

			retryDuration := 5 * time.Second
			if ctxDeadline, ok := CTX.Deadline(); ok {
				retryDuration = time.Until(ctxDeadline)
			}

			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := Client.SearchTargets(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(ttt, listErr, "Error: "+listErr.Error())
				} else {
					assert.NoError(ttt, listErr)
				}
				if listErr != nil {
					return
				}
				// always first check length, otherwise its failed anyway
				assert.Len(ttt, got.Result, len(tt.want.Result))
				for i := range tt.want.Result {
					assert.Contains(ttt, got.Result, tt.want.Result[i])
				}
				integration.AssertResourceListDetails(t, tt.want, got)
			}, retryDuration, time.Millisecond*100, "timeout waiting for expected execution result")
		})
	}
}

func TestServer_SearchExecutions(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t, "", "https://example.com", domain.TargetTypeWebhook, false)

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
				ctx: Tester.WithAuthorization(context.Background(), integration.OrgOwner),
				req: &action.SearchExecutionsRequest{},
			},
			wantErr: true,
		},
		{
			name: "list request single condition",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.SearchExecutionsRequest, response *action.SearchExecutionsResponse) error {
					cond := request.Filters[0].GetInConditionsFilter().GetConditions()[0]
					resp := Tester.SetExecution(ctx, t, cond, executionTargetsSingleTarget(targetResp.GetDetails().GetId()))

					response.Details.Timestamp = resp.GetDetails().GetChangeDate()
					// response.Details.ProcessedSequence = resp.GetDetails().GetSequence()

					// Set expected response with used values for SetExecution
					response.Result[0].Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Result[0].Details.Sequence = resp.GetDetails().GetSequence()
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
												Method: "/zitadel.session.v2beta.SessionService/GetSession",
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
					TotalResult: 1,
				},
				Result: []*action.GetExecution{
					{
						Details: &resource_object.Details{
							Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
						},
						Condition: &action.Condition{
							ConditionType: &action.Condition_Request{
								Request: &action.RequestExecution{
									Condition: &action.RequestExecution_Method{
										Method: "/zitadel.session.v2beta.SessionService/GetSession",
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
				ctx: CTX,
				dep: func(ctx context.Context, request *action.SearchExecutionsRequest, response *action.SearchExecutionsResponse) error {
					target := Tester.CreateTarget(CTX, t, "", "https://example.com", domain.TargetTypeWebhook, false)
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
					resp := Tester.SetExecution(ctx, t, cond, targets)

					response.Details.Timestamp = resp.GetDetails().GetChangeDate()
					response.Details.ProcessedSequence = resp.GetDetails().GetSequence()

					response.Result[0].Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Result[0].Details.Sequence = resp.GetDetails().GetSequence()
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
					TotalResult: 1,
				},
				Result: []*action.GetExecution{
					{
						Details: &resource_object.Details{
							Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
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
				ctx: CTX,
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
					Tester.SetExecution(ctx, t, cond, executionTargetsSingleTarget(targetResp.GetDetails().GetId()))
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
					resp2 := Tester.SetExecution(ctx, t, includeCond, includeTargets)

					response.Details.Timestamp = resp2.GetDetails().GetChangeDate()
					response.Details.ProcessedSequence = resp2.GetDetails().GetSequence()

					response.Result[0].Details.ChangeDate = resp2.GetDetails().GetChangeDate()
					response.Result[0].Details.Sequence = resp2.GetDetails().GetSequence()
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
					TotalResult: 1,
				},
				Result: []*action.GetExecution{
					{
						Details: &resource_object.Details{
							Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
						},
					},
				},
			},
		},
		{
			name: "list multiple conditions",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.SearchExecutionsRequest, response *action.SearchExecutionsResponse) error {

					cond1 := request.Filters[0].GetInConditionsFilter().GetConditions()[0]
					targets1 := executionTargetsSingleTarget(targetResp.GetDetails().GetId())
					resp1 := Tester.SetExecution(ctx, t, cond1, targets1)
					response.Result[0].Details.ChangeDate = resp1.GetDetails().GetChangeDate()
					response.Result[0].Details.Sequence = resp1.GetDetails().GetSequence()
					response.Result[0].Condition = cond1
					response.Result[0].Execution = &action.Execution{
						Targets: targets1,
					}

					cond2 := request.Filters[0].GetInConditionsFilter().GetConditions()[1]
					targets2 := executionTargetsSingleTarget(targetResp.GetDetails().GetId())
					resp2 := Tester.SetExecution(ctx, t, cond2, targets2)
					response.Result[1].Details.ChangeDate = resp2.GetDetails().GetChangeDate()
					response.Result[1].Details.Sequence = resp2.GetDetails().GetSequence()
					response.Result[1].Condition = cond2
					response.Result[1].Execution = &action.Execution{
						Targets: targets2,
					}

					cond3 := request.Filters[0].GetInConditionsFilter().GetConditions()[2]
					targets3 := executionTargetsSingleTarget(targetResp.GetDetails().GetId())
					resp3 := Tester.SetExecution(ctx, t, cond3, targets3)
					response.Result[2].Details.ChangeDate = resp3.GetDetails().GetChangeDate()
					response.Result[2].Details.Sequence = resp3.GetDetails().GetSequence()
					response.Result[2].Condition = cond3
					response.Result[2].Execution = &action.Execution{
						Targets: targets3,
					}
					response.Details.Timestamp = resp3.GetDetails().GetChangeDate()
					response.Details.ProcessedSequence = resp3.GetDetails().GetSequence()
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
													Method: "/zitadel.session.v2beta.SessionService/GetSession",
												},
											},
										},
									},
									{
										ConditionType: &action.Condition_Request{
											Request: &action.RequestExecution{
												Condition: &action.RequestExecution_Method{
													Method: "/zitadel.session.v2beta.SessionService/CreateSession",
												},
											},
										},
									},
									{
										ConditionType: &action.Condition_Request{
											Request: &action.RequestExecution{
												Condition: &action.RequestExecution_Method{
													Method: "/zitadel.session.v2beta.SessionService/SetSession",
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
					TotalResult: 3,
				},
				Result: []*action.GetExecution{
					{
						Details: &resource_object.Details{
							Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
						},
					}, {
						Details: &resource_object.Details{
							Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
						},
					}, {
						Details: &resource_object.Details{
							Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()},
						},
					},
				},
			},
		},
		{
			name: "list multiple conditions all types",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.SearchExecutionsRequest, response *action.SearchExecutionsResponse) error {
					targets := executionTargetsSingleTarget(targetResp.GetDetails().GetId())
					for i, cond := range request.Filters[0].GetInConditionsFilter().GetConditions() {
						resp := Tester.SetExecution(ctx, t, cond, targets)
						response.Result[i].Details.ChangeDate = resp.GetDetails().GetChangeDate()
						response.Result[i].Details.Sequence = resp.GetDetails().GetSequence()
						response.Result[i].Condition = cond
						response.Result[i].Execution = &action.Execution{
							Targets: targets,
						}
						// filled with info of last sequence
						response.Details.Timestamp = resp.GetDetails().GetChangeDate()
						response.Details.ProcessedSequence = resp.GetDetails().GetSequence()
					}

					return nil
				},
				req: &action.SearchExecutionsRequest{
					Filters: []*action.ExecutionSearchFilter{{
						Filter: &action.ExecutionSearchFilter_InConditionsFilter{
							InConditionsFilter: &action.InConditionsFilter{
								Conditions: []*action.Condition{
									{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_Method{Method: "/zitadel.session.v2beta.SessionService/GetSession"}}}},
									{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_Service{Service: "zitadel.session.v2beta.SessionService"}}}},
									{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_All{All: true}}}},
									{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_Method{Method: "/zitadel.session.v2beta.SessionService/GetSession"}}}},
									{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_Service{Service: "zitadel.session.v2beta.SessionService"}}}},
									{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_All{All: true}}}},
									{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_Event{Event: "user.added"}}}},
									{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_Group{Group: "user"}}}},
									{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_All{All: true}}}},
									{ConditionType: &action.Condition_Function{Function: &action.FunctionExecution{Name: "Action.Flow.Type.ExternalAuthentication.Action.TriggerType.PostAuthentication"}}},
								},
							},
						},
					}},
				},
			},
			want: &action.SearchExecutionsResponse{
				Details: &resource_object.ListDetails{
					TotalResult: 10,
				},
				Result: []*action.GetExecution{
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()}}},
					{Details: &resource_object.Details{Owner: &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: Tester.Instance.InstanceID()}}},
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

			retryDuration := 5 * time.Second
			if ctxDeadline, ok := CTX.Deadline(); ok {
				retryDuration = time.Until(ctxDeadline)
			}

			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := Client.SearchExecutions(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(t, listErr, "Error: "+listErr.Error())
				} else {
					assert.NoError(t, listErr)
				}
				if listErr != nil {
					return
				}
				// always first check length, otherwise its failed anyway
				assert.Len(t, got.Result, len(tt.want.Result))
				for i := range tt.want.Result {
					// as not sorted, all elements have to be checked
					// workaround as oneof elements can only be checked with assert.EqualExportedValues()
					if j, found := containExecution(got.Result, tt.want.Result[i]); found {
						assert.EqualExportedValues(t, tt.want.Result[i], got.Result[j])
					}
				}
				integration.AssertResourceListDetails(t, tt.want, got)
			}, retryDuration, time.Millisecond*100, "timeout waiting for expected execution result")
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
