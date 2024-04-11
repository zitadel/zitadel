//go:build integration

package action_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/integration"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v3alpha"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
)

func TestServer_GetTargetByID(t *testing.T) {
	ensureFeatureEnabled(t)
	type args struct {
		ctx context.Context
		dep func(context.Context, *action.GetTargetByIDRequest, *action.GetTargetByIDResponse) error
		req *action.GetTargetByIDRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *action.GetTargetByIDResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			args: args{
				ctx: Tester.WithAuthorization(context.Background(), integration.OrgOwner),
				req: &action.GetTargetByIDRequest{},
			},
			wantErr: true,
		},
		{
			name: "not found",
			args: args{
				ctx: CTX,
				req: &action.GetTargetByIDRequest{TargetId: "notexisting"},
			},
			wantErr: true,
		},
		{
			name: "get, ok",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.GetTargetByIDRequest, response *action.GetTargetByIDResponse) error {
					name := fmt.Sprint(time.Now().UnixNano() + 1)
					resp := Tester.CreateTargetWithNameAndType(ctx, t, name, false, false)
					request.TargetId = resp.GetId()

					response.Target.TargetId = resp.GetId()
					response.Target.Name = name
					response.Target.Details.ResourceOwner = resp.GetDetails().GetResourceOwner()
					response.Target.Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Target.Details.Sequence = resp.GetDetails().GetSequence()
					return nil
				},
				req: &action.GetTargetByIDRequest{},
			},
			want: &action.GetTargetByIDResponse{
				Target: &action.Target{
					Details: &object.Details{
						ResourceOwner: Tester.Instance.InstanceID(),
					},
					TargetType: &action.Target_RestWebhook{
						RestWebhook: &action.SetRESTWebhook{
							Url: "https://example.com",
						},
					},
					Timeout: durationpb.New(10 * time.Second),
				},
			},
		},
		{
			name: "get, async, ok",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.GetTargetByIDRequest, response *action.GetTargetByIDResponse) error {
					name := fmt.Sprint(time.Now().UnixNano() + 1)
					resp := Tester.CreateTargetWithNameAndType(ctx, t, name, true, false)
					request.TargetId = resp.GetId()

					response.Target.TargetId = resp.GetId()
					response.Target.Name = name
					response.Target.Details.ResourceOwner = resp.GetDetails().GetResourceOwner()
					response.Target.Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Target.Details.Sequence = resp.GetDetails().GetSequence()
					return nil
				},
				req: &action.GetTargetByIDRequest{},
			},
			want: &action.GetTargetByIDResponse{
				Target: &action.Target{
					Details: &object.Details{
						ResourceOwner: Tester.Instance.InstanceID(),
					},
					TargetType: &action.Target_RestWebhook{
						RestWebhook: &action.SetRESTWebhook{
							Url: "https://example.com",
						},
					},
					Timeout:       durationpb.New(10 * time.Second),
					ExecutionType: &action.Target_IsAsync{IsAsync: true},
				},
			},
		},
		{
			name: "get, interruptOnError, ok",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.GetTargetByIDRequest, response *action.GetTargetByIDResponse) error {
					name := fmt.Sprint(time.Now().UnixNano() + 1)
					resp := Tester.CreateTargetWithNameAndType(ctx, t, name, false, true)
					request.TargetId = resp.GetId()

					response.Target.TargetId = resp.GetId()
					response.Target.Name = name
					response.Target.Details.ResourceOwner = resp.GetDetails().GetResourceOwner()
					response.Target.Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Target.Details.Sequence = resp.GetDetails().GetSequence()
					return nil
				},
				req: &action.GetTargetByIDRequest{},
			},
			want: &action.GetTargetByIDResponse{
				Target: &action.Target{
					Details: &object.Details{
						ResourceOwner: Tester.Instance.InstanceID(),
					},
					TargetType: &action.Target_RestWebhook{
						RestWebhook: &action.SetRESTWebhook{
							Url: "https://example.com",
						},
					},
					Timeout:       durationpb.New(10 * time.Second),
					ExecutionType: &action.Target_InterruptOnError{InterruptOnError: true},
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
				got, getErr := Client.GetTargetByID(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(ttt, getErr, "Error: "+getErr.Error())
				} else {
					assert.NoError(ttt, getErr)
				}
				if getErr != nil {
					fmt.Println("Error: " + getErr.Error())
					return
				}

				integration.AssertDetails(t, tt.want.GetTarget(), got.GetTarget())

				assert.Equal(t, tt.want.Target, got.Target)

			}, retryDuration, time.Millisecond*100, "timeout waiting for expected execution result")
		})
	}
}

func TestServer_ListTargets(t *testing.T) {
	ensureFeatureEnabled(t)
	type args struct {
		ctx context.Context
		dep func(context.Context, *action.ListTargetsRequest, *action.ListTargetsResponse) error
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
				ctx: Tester.WithAuthorization(context.Background(), integration.OrgOwner),
				req: &action.ListTargetsRequest{},
			},
			wantErr: true,
		},
		{
			name: "list, not found",
			args: args{
				ctx: CTX,
				req: &action.ListTargetsRequest{
					Queries: []*action.TargetSearchQuery{
						{Query: &action.TargetSearchQuery_InTargetIdsQuery{
							InTargetIdsQuery: &action.InTargetIDsQuery{
								TargetIds: []string{"notfound"},
							},
						},
						},
					},
				},
			},
			want: &action.ListTargetsResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
				},
				Result: []*action.Target{},
			},
		},
		{
			name: "list single id",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.ListTargetsRequest, response *action.ListTargetsResponse) error {
					name := fmt.Sprint(time.Now().UnixNano() + 1)
					resp := Tester.CreateTargetWithNameAndType(ctx, t, name, false, false)
					request.Queries[0].Query = &action.TargetSearchQuery_InTargetIdsQuery{
						InTargetIdsQuery: &action.InTargetIDsQuery{
							TargetIds: []string{resp.GetId()},
						},
					}
					response.Details.Timestamp = resp.GetDetails().GetChangeDate()
					response.Details.ProcessedSequence = resp.GetDetails().GetSequence()

					response.Result[0].Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Result[0].Details.Sequence = resp.GetDetails().GetSequence()
					response.Result[0].TargetId = resp.GetId()
					response.Result[0].Name = name
					return nil
				},
				req: &action.ListTargetsRequest{
					Queries: []*action.TargetSearchQuery{{}},
				},
			},
			want: &action.ListTargetsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				Result: []*action.Target{
					{
						Details: &object.Details{
							ResourceOwner: Tester.Instance.InstanceID(),
						},
						TargetType: &action.Target_RestWebhook{
							RestWebhook: &action.SetRESTWebhook{
								Url: "https://example.com",
							},
						},
						Timeout: durationpb.New(10 * time.Second),
					},
				},
			},
		}, {
			name: "list single name",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.ListTargetsRequest, response *action.ListTargetsResponse) error {
					name := fmt.Sprint(time.Now().UnixNano() + 1)
					resp := Tester.CreateTargetWithNameAndType(ctx, t, name, false, false)
					request.Queries[0].Query = &action.TargetSearchQuery_TargetNameQuery{
						TargetNameQuery: &action.TargetNameQuery{
							TargetName: name,
						},
					}
					response.Details.Timestamp = resp.GetDetails().GetChangeDate()
					response.Details.ProcessedSequence = resp.GetDetails().GetSequence()

					response.Result[0].Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Result[0].Details.Sequence = resp.GetDetails().GetSequence()
					response.Result[0].TargetId = resp.GetId()
					response.Result[0].Name = name
					return nil
				},
				req: &action.ListTargetsRequest{
					Queries: []*action.TargetSearchQuery{{}},
				},
			},
			want: &action.ListTargetsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				Result: []*action.Target{
					{
						Details: &object.Details{
							ResourceOwner: Tester.Instance.InstanceID(),
						},
						TargetType: &action.Target_RestWebhook{
							RestWebhook: &action.SetRESTWebhook{
								Url: "https://example.com",
							},
						},
						Timeout: durationpb.New(10 * time.Second),
					},
				},
			},
		},
		{
			name: "list multiple id",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.ListTargetsRequest, response *action.ListTargetsResponse) error {
					name1 := fmt.Sprint(time.Now().UnixNano() + 1)
					name2 := fmt.Sprint(time.Now().UnixNano() + 3)
					name3 := fmt.Sprint(time.Now().UnixNano() + 5)
					resp1 := Tester.CreateTargetWithNameAndType(ctx, t, name1, false, false)
					resp2 := Tester.CreateTargetWithNameAndType(ctx, t, name2, true, false)
					resp3 := Tester.CreateTargetWithNameAndType(ctx, t, name3, false, true)
					request.Queries[0].Query = &action.TargetSearchQuery_InTargetIdsQuery{
						InTargetIdsQuery: &action.InTargetIDsQuery{
							TargetIds: []string{resp1.GetId(), resp2.GetId(), resp3.GetId()},
						},
					}
					response.Details.Timestamp = resp3.GetDetails().GetChangeDate()
					response.Details.ProcessedSequence = resp3.GetDetails().GetSequence()

					response.Result[0].Details.ChangeDate = resp1.GetDetails().GetChangeDate()
					response.Result[0].Details.Sequence = resp1.GetDetails().GetSequence()
					response.Result[0].TargetId = resp1.GetId()
					response.Result[0].Name = name1
					response.Result[1].Details.ChangeDate = resp2.GetDetails().GetChangeDate()
					response.Result[1].Details.Sequence = resp2.GetDetails().GetSequence()
					response.Result[1].TargetId = resp2.GetId()
					response.Result[1].Name = name2
					response.Result[2].Details.ChangeDate = resp3.GetDetails().GetChangeDate()
					response.Result[2].Details.Sequence = resp3.GetDetails().GetSequence()
					response.Result[2].TargetId = resp3.GetId()
					response.Result[2].Name = name3
					return nil
				},
				req: &action.ListTargetsRequest{
					Queries: []*action.TargetSearchQuery{{}},
				},
			},
			want: &action.ListTargetsResponse{
				Details: &object.ListDetails{
					TotalResult: 3,
				},
				Result: []*action.Target{
					{
						Details: &object.Details{
							ResourceOwner: Tester.Instance.InstanceID(),
						},
						TargetType: &action.Target_RestWebhook{
							RestWebhook: &action.SetRESTWebhook{
								Url: "https://example.com",
							},
						},
						Timeout: durationpb.New(10 * time.Second),
					},
					{
						Details: &object.Details{
							ResourceOwner: Tester.Instance.InstanceID(),
						},
						TargetType: &action.Target_RestWebhook{
							RestWebhook: &action.SetRESTWebhook{
								Url: "https://example.com",
							},
						},
						Timeout:       durationpb.New(10 * time.Second),
						ExecutionType: &action.Target_IsAsync{IsAsync: true},
					},
					{
						Details: &object.Details{
							ResourceOwner: Tester.Instance.InstanceID(),
						},
						TargetType: &action.Target_RestWebhook{
							RestWebhook: &action.SetRESTWebhook{
								Url: "https://example.com",
							},
						},
						Timeout:       durationpb.New(10 * time.Second),
						ExecutionType: &action.Target_InterruptOnError{InterruptOnError: true},
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
				got, listErr := Client.ListTargets(tt.args.ctx, tt.args.req)
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
				integration.AssertListDetails(t, tt.want, got)
			}, retryDuration, time.Millisecond*100, "timeout waiting for expected execution result")
		})
	}
}

func TestServer_ListExecutions_Request(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t)

	type args struct {
		ctx context.Context
		dep func(context.Context, *action.ListExecutionsRequest, *action.ListExecutionsResponse) error
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
				ctx: Tester.WithAuthorization(context.Background(), integration.OrgOwner),
				req: &action.ListExecutionsRequest{},
			},
			wantErr: true,
		},
		{
			name: "list single condition",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.ListExecutionsRequest, response *action.ListExecutionsResponse) error {
					resp := Tester.SetExecution(ctx, t, request.Queries[0].GetInConditionsQuery().GetConditions()[0], []string{targetResp.GetId()}, []string{})

					response.Details.Timestamp = resp.GetDetails().GetChangeDate()
					response.Details.ProcessedSequence = resp.GetDetails().GetSequence()

					response.Result[0].Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Result[0].Details.Sequence = resp.GetDetails().GetSequence()
					return nil
				},
				req: &action.ListExecutionsRequest{
					Queries: []*action.SearchQuery{{
						Query: &action.SearchQuery_InConditionsQuery{
							InConditionsQuery: &action.InConditionsQuery{
								Conditions: []*action.Condition{{
									ConditionType: &action.Condition_Request{
										Request: &action.RequestExecution{
											Condition: &action.RequestExecution_Method{
												Method: "/zitadel.session.v2beta.SessionService/GetSession",
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
			want: &action.ListExecutionsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				Result: []*action.Execution{
					{
						Details: &object.Details{
							ResourceOwner: Tester.Instance.InstanceID(),
						},
						ExecutionId: "request./zitadel.session.v2beta.SessionService/GetSession",
						Targets:     []string{targetResp.GetId()},
					},
				},
			},
		},
		{
			name: "list single target",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.ListExecutionsRequest, response *action.ListExecutionsResponse) error {
					target := Tester.CreateTarget(ctx, t)
					// add target as query to the request
					request.Queries[0] = &action.SearchQuery{
						Query: &action.SearchQuery_TargetQuery{
							TargetQuery: &action.TargetQuery{
								TargetId: target.GetId(),
							},
						},
					}
					resp := Tester.SetExecution(ctx, t, &action.Condition{
						ConditionType: &action.Condition_Request{
							Request: &action.RequestExecution{
								Condition: &action.RequestExecution_Method{
									Method: "/zitadel.management.v1.ManagementService/UpdateAction",
								},
							},
						},
					}, []string{target.GetId()}, []string{})

					response.Details.Timestamp = resp.GetDetails().GetChangeDate()
					response.Details.ProcessedSequence = resp.GetDetails().GetSequence()

					response.Result[0].Details.ChangeDate = resp.GetDetails().GetChangeDate()
					response.Result[0].Details.Sequence = resp.GetDetails().GetSequence()
					response.Result[0].Targets[0] = target.GetId()
					return nil
				},
				req: &action.ListExecutionsRequest{
					Queries: []*action.SearchQuery{{}},
				},
			},
			want: &action.ListExecutionsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				Result: []*action.Execution{
					{
						Details: &object.Details{
							ResourceOwner: Tester.Instance.InstanceID(),
						},
						ExecutionId: "request./zitadel.management.v1.ManagementService/UpdateAction",
						Targets:     []string{""},
					},
				},
			},
		}, {
			name: "list single include",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.ListExecutionsRequest, response *action.ListExecutionsResponse) error {
					Tester.SetExecution(ctx, t, &action.Condition{
						ConditionType: &action.Condition_Request{
							Request: &action.RequestExecution{
								Condition: &action.RequestExecution_Method{
									Method: "/zitadel.management.v1.ManagementService/GetAction",
								},
							},
						},
					}, []string{targetResp.GetId()}, []string{})
					resp2 := Tester.SetExecution(ctx, t, &action.Condition{
						ConditionType: &action.Condition_Request{
							Request: &action.RequestExecution{
								Condition: &action.RequestExecution_Method{
									Method: "/zitadel.management.v1.ManagementService/ListActions",
								},
							},
						},
					}, []string{}, []string{"request./zitadel.management.v1.ManagementService/GetAction"})

					response.Details.Timestamp = resp2.GetDetails().GetChangeDate()
					response.Details.ProcessedSequence = resp2.GetDetails().GetSequence()

					response.Result[0].Details.ChangeDate = resp2.GetDetails().GetChangeDate()
					response.Result[0].Details.Sequence = resp2.GetDetails().GetSequence()
					return nil
				},
				req: &action.ListExecutionsRequest{
					Queries: []*action.SearchQuery{{
						Query: &action.SearchQuery_IncludeQuery{
							IncludeQuery: &action.IncludeQuery{Include: "request./zitadel.management.v1.ManagementService/GetAction"},
						},
					}},
				},
			},
			want: &action.ListExecutionsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				Result: []*action.Execution{
					{
						Details: &object.Details{
							ResourceOwner: Tester.Instance.InstanceID(),
						},
						ExecutionId: "request./zitadel.management.v1.ManagementService/ListActions",
						Includes:    []string{"request./zitadel.management.v1.ManagementService/GetAction"},
					},
				},
			},
		},
		{
			name: "list multiple conditions",
			args: args{
				ctx: CTX,
				dep: func(ctx context.Context, request *action.ListExecutionsRequest, response *action.ListExecutionsResponse) error {

					resp1 := Tester.SetExecution(ctx, t, request.Queries[0].GetInConditionsQuery().GetConditions()[0], []string{targetResp.GetId()}, []string{})
					response.Result[0].Details.ChangeDate = resp1.GetDetails().GetChangeDate()
					response.Result[0].Details.Sequence = resp1.GetDetails().GetSequence()

					resp2 := Tester.SetExecution(ctx, t, request.Queries[0].GetInConditionsQuery().GetConditions()[1], []string{targetResp.GetId()}, []string{})
					response.Result[1].Details.ChangeDate = resp2.GetDetails().GetChangeDate()
					response.Result[1].Details.Sequence = resp2.GetDetails().GetSequence()

					resp3 := Tester.SetExecution(ctx, t, request.Queries[0].GetInConditionsQuery().GetConditions()[2], []string{targetResp.GetId()}, []string{})
					response.Details.Timestamp = resp3.GetDetails().GetChangeDate()
					response.Details.ProcessedSequence = resp3.GetDetails().GetSequence()
					response.Result[2].Details.ChangeDate = resp3.GetDetails().GetChangeDate()
					response.Result[2].Details.Sequence = resp3.GetDetails().GetSequence()
					return nil
				},
				req: &action.ListExecutionsRequest{
					Queries: []*action.SearchQuery{{
						Query: &action.SearchQuery_InConditionsQuery{
							InConditionsQuery: &action.InConditionsQuery{
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
			want: &action.ListExecutionsResponse{
				Details: &object.ListDetails{
					TotalResult: 3,
				},
				Result: []*action.Execution{
					{
						Details: &object.Details{
							ResourceOwner: Tester.Instance.InstanceID(),
						},
						ExecutionId: "request./zitadel.session.v2beta.SessionService/GetSession",
						Targets:     []string{targetResp.GetId()},
					}, {
						Details: &object.Details{
							ResourceOwner: Tester.Instance.InstanceID(),
						},
						ExecutionId: "request./zitadel.session.v2beta.SessionService/CreateSession",
						Targets:     []string{targetResp.GetId()},
					}, {
						Details: &object.Details{
							ResourceOwner: Tester.Instance.InstanceID(),
						},
						ExecutionId: "request./zitadel.session.v2beta.SessionService/SetSession",
						Targets:     []string{targetResp.GetId()},
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
				got, listErr := Client.ListExecutions(tt.args.ctx, tt.args.req)
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
				integration.AssertListDetails(t, tt.want, got)
			}, retryDuration, time.Millisecond*100, "timeout waiting for expected execution result")
		})
	}
}
