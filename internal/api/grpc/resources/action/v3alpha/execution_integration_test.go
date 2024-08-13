//go:build integration

package action_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	action "github.com/zitadel/zitadel/pkg/grpc/resources/action/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
)

func executionTargetsSingleTarget(id string) []*action.ExecutionTargetType {
	return []*action.ExecutionTargetType{{Type: &action.ExecutionTargetType_Target{Target: id}}}
}

func executionTargetsSingleInclude(include *action.Condition) []*action.ExecutionTargetType {
	return []*action.ExecutionTargetType{{Type: &action.ExecutionTargetType_Include{Include: include}}}
}

func TestServer_SetExecution_Request(t *testing.T) {
	_, instanceID, _, isolatedIAMOwnerCTX := Tester.UseIsolatedInstance(t, IAMOwnerCTX, SystemCTX)
	ensureFeatureEnabled(t, isolatedIAMOwnerCTX)
	targetResp := Tester.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.SetExecutionRequest
		want    *action.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_All{All: true},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			wantErr: true,
		},
		{
			name: "method, not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Method{
								Method: "/zitadel.session.v2beta.NotExistingService/List",
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			wantErr: true,
		},
		{
			name: "method, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/ListSessions",
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			want: &action.SetExecutionResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instanceID,
					},
				},
			},
		},
		{
			name: "service, not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Service{
								Service: "NotExistingService",
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			wantErr: true,
		},
		{
			name: "service, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Service{
								Service: "zitadel.session.v2beta.SessionService",
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			want: &action.SetExecutionResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instanceID,
					},
				},
			},
		},
		{
			name: "all, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_All{
								All: true,
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			want: &action.SetExecutionResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instanceID,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We want to have the same response no matter how often we call the function
			Tester.Client.ActionV3.SetExecution(tt.ctx, tt.req)
			got, err := Tester.Client.ActionV3.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			integration.AssertResourceDetails(t, tt.want.Details, got.Details)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_SetExecution_Request_Include(t *testing.T) {
	_, instanceID, _, isolatedIAMOwnerCTX := Tester.UseIsolatedInstance(t, IAMOwnerCTX, SystemCTX)
	ensureFeatureEnabled(t, isolatedIAMOwnerCTX)
	targetResp := Tester.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)
	executionCond := &action.Condition{
		ConditionType: &action.Condition_Request{
			Request: &action.RequestExecution{
				Condition: &action.RequestExecution_All{
					All: true,
				},
			},
		},
	}
	Tester.SetExecution(isolatedIAMOwnerCTX, t,
		executionCond,
		executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
	)

	circularExecutionService := &action.Condition{
		ConditionType: &action.Condition_Request{
			Request: &action.RequestExecution{
				Condition: &action.RequestExecution_Service{
					Service: "zitadel.session.v2beta.SessionService",
				},
			},
		},
	}
	Tester.SetExecution(isolatedIAMOwnerCTX, t,
		circularExecutionService,
		executionTargetsSingleInclude(executionCond),
	)
	circularExecutionMethod := &action.Condition{
		ConditionType: &action.Condition_Request{
			Request: &action.RequestExecution{
				Condition: &action.RequestExecution_Method{
					Method: "/zitadel.session.v2beta.SessionService/ListSessions",
				},
			},
		},
	}
	Tester.SetExecution(isolatedIAMOwnerCTX, t,
		circularExecutionMethod,
		executionTargetsSingleInclude(circularExecutionService),
	)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.SetExecutionRequest
		want    *action.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "method, circular error",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: circularExecutionService,
				Execution: &action.Execution{
					Targets: executionTargetsSingleInclude(circularExecutionMethod),
				},
			},
			wantErr: true,
		},
		{
			name: "method, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/ListSessions",
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleInclude(executionCond),
				},
			},
			want: &action.SetExecutionResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instanceID,
					},
				},
			},
		},
		{
			name: "service, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Service{
								Service: "zitadel.session.v2beta.SessionService",
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleInclude(executionCond),
				},
			},
			want: &action.SetExecutionResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instanceID,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We want to have the same response no matter how often we call the function
			Tester.Client.ActionV3.SetExecution(tt.ctx, tt.req)
			got, err := Tester.Client.ActionV3.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertResourceDetails(t, tt.want.Details, got.Details)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_SetExecution_Response(t *testing.T) {
	_, instanceID, _, isolatedIAMOwnerCTX := Tester.UseIsolatedInstance(t, IAMOwnerCTX, SystemCTX)
	ensureFeatureEnabled(t, isolatedIAMOwnerCTX)
	targetResp := Tester.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.SetExecutionRequest
		want    *action.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_All{All: true},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			wantErr: true,
		},
		{
			name: "method, not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_Method{
								Method: "/zitadel.session.v2beta.NotExistingService/List",
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			wantErr: true,
		},
		{
			name: "method, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/ListSessions",
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			want: &action.SetExecutionResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instanceID,
					},
				},
			},
		},
		{
			name: "service, not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_Service{
								Service: "NotExistingService",
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			wantErr: true,
		},
		{
			name: "service, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_Service{
								Service: "zitadel.session.v2beta.SessionService",
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			want: &action.SetExecutionResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instanceID,
					},
				},
			},
		},
		{
			name: "all, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_All{
								All: true,
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			want: &action.SetExecutionResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instanceID,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We want to have the same response no matter how often we call the function
			Tester.Client.ActionV3.SetExecution(tt.ctx, tt.req)
			got, err := Tester.Client.ActionV3.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertResourceDetails(t, tt.want.Details, got.Details)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_SetExecution_Event(t *testing.T) {
	_, instanceID, _, isolatedIAMOwnerCTX := Tester.UseIsolatedInstance(t, IAMOwnerCTX, SystemCTX)
	ensureFeatureEnabled(t, isolatedIAMOwnerCTX)
	targetResp := Tester.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.SetExecutionRequest
		want    *action.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{
							Condition: &action.EventExecution_All{
								All: true,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			wantErr: true,
		},
		/*
			//TODO event existing check

			{
				name: "event, not existing",
				ctx:  isolatedIAMOwnerCTX,
				req: &action.SetExecutionRequest{
					Condition: &action.Condition{
						ConditionType: &action.Condition_Event{
							Event: &action.EventExecution{
								Condition: &action.EventExecution_Event{
									Event: "xxx",
								},
							},
						},
					},
					Targets: []string{targetResp.GetId()},
				},
				wantErr: true,
			},
		*/
		{
			name: "event, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{
							Condition: &action.EventExecution_Event{
								Event: "xxx",
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			want: &action.SetExecutionResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instanceID,
					},
				},
			},
		},
		/*
			// TODO:

			{
				name: "group, not existing",
				ctx:  isolatedIAMOwnerCTX,
				req: &action.SetExecutionRequest{
					Condition: &action.Condition{
						ConditionType: &action.Condition_Event{
							Event: &action.EventExecution{
								Condition: &action.EventExecution_Group{
									Group: "xxx",
								},
							},
						},
					},
					Targets: []string{targetResp.GetId()},
				},
				wantErr: true,
			},
		*/
		{
			name: "group, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{
							Condition: &action.EventExecution_Group{
								Group: "xxx",
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			want: &action.SetExecutionResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instanceID,
					},
				},
			},
		},
		{
			name: "all, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{
							Condition: &action.EventExecution_All{
								All: true,
							},
						},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			want: &action.SetExecutionResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instanceID,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We want to have the same response no matter how often we call the function
			Tester.Client.ActionV3.SetExecution(tt.ctx, tt.req)
			got, err := Tester.Client.ActionV3.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertResourceDetails(t, tt.want.Details, got.Details)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_SetExecution_Function(t *testing.T) {
	_, instanceID, _, isolatedIAMOwnerCTX := Tester.UseIsolatedInstance(t, IAMOwnerCTX, SystemCTX)
	ensureFeatureEnabled(t, isolatedIAMOwnerCTX)
	targetResp := Tester.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.SetExecutionRequest
		want    *action.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_All{All: true},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			wantErr: true,
		},
		{
			name: "function, not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Function{
						Function: &action.FunctionExecution{Name: "xxx"},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			wantErr: true,
		},
		{
			name: "function, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Function{
						Function: &action.FunctionExecution{Name: "Action.Flow.Type.ExternalAuthentication.Action.TriggerType.PostAuthentication"},
					},
				},
				Execution: &action.Execution{
					Targets: executionTargetsSingleTarget(targetResp.GetDetails().GetId()),
				},
			},
			want: &action.SetExecutionResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instanceID,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We want to have the same response no matter how often we call the function
			Tester.Client.ActionV3.SetExecution(tt.ctx, tt.req)
			got, err := Tester.Client.ActionV3.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertResourceDetails(t, tt.want.Details, got.Details)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}
