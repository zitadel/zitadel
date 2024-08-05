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
	settings_object "github.com/zitadel/zitadel/pkg/grpc/settings/object/v3alpha"
)

func executionTargetsSingleTarget(id string) []*action.ExecutionTargetType {
	return []*action.ExecutionTargetType{{Type: &action.ExecutionTargetType_Target{Target: id}}}
}

func executionTargetsSingleInclude(include *action.Condition) []*action.ExecutionTargetType {
	return []*action.ExecutionTargetType{{Type: &action.ExecutionTargetType_Include{Include: include}}}
}

func TestServer_SetExecution_Request(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

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
			ctx:  CTX,
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
			ctx:  CTX,
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
			ctx:  CTX,
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
				Details: &settings_object.Details{
					ChangeDate: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   Tester.Instance.InstanceID(),
					},
				},
			},
		},
		{
			name: "service, not existing",
			ctx:  CTX,
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
			ctx:  CTX,
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
				Details: &settings_object.Details{
					ChangeDate: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   Tester.Instance.InstanceID(),
					},
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
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
				Details: &settings_object.Details{
					ChangeDate: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   Tester.Instance.InstanceID(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We want to have the same response no matter how often we call the function
			Client.SetExecution(tt.ctx, tt.req)
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			integration.AssertSettingsDetails(t, tt.want.Details, got.Details)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_SetExecution_Request_Include(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)
	executionCond := &action.Condition{
		ConditionType: &action.Condition_Request{
			Request: &action.RequestExecution{
				Condition: &action.RequestExecution_All{
					All: true,
				},
			},
		},
	}
	Tester.SetExecution(CTX, t,
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
	Tester.SetExecution(CTX, t,
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
	Tester.SetExecution(CTX, t,
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
			ctx:  CTX,
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
			ctx:  CTX,
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
				Details: &settings_object.Details{
					ChangeDate: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   Tester.Instance.InstanceID(),
					},
				},
			},
		},
		{
			name: "service, ok",
			ctx:  CTX,
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
				Details: &settings_object.Details{
					ChangeDate: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   Tester.Instance.InstanceID(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We want to have the same response no matter how often we call the function
			Client.SetExecution(tt.ctx, tt.req)
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertSettingsDetails(t, tt.want.Details, got.Details)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_SetExecution_Response(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

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
			ctx:  CTX,
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
			ctx:  CTX,
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
			ctx:  CTX,
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
				Details: &settings_object.Details{
					ChangeDate: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   Tester.Instance.InstanceID(),
					},
				},
			},
		},
		{
			name: "service, not existing",
			ctx:  CTX,
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
			ctx:  CTX,
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
				Details: &settings_object.Details{
					ChangeDate: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   Tester.Instance.InstanceID(),
					},
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
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
				Details: &settings_object.Details{
					ChangeDate: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   Tester.Instance.InstanceID(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We want to have the same response no matter how often we call the function
			Client.SetExecution(tt.ctx, tt.req)
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertSettingsDetails(t, tt.want.Details, got.Details)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_SetExecution_Event(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

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
			ctx:  CTX,
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
				ctx:  CTX,
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
			ctx:  CTX,
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
				Details: &settings_object.Details{
					ChangeDate: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   Tester.Instance.InstanceID(),
					},
				},
			},
		},
		/*
			// TODO:

			{
				name: "group, not existing",
				ctx:  CTX,
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
			ctx:  CTX,
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
				Details: &settings_object.Details{
					ChangeDate: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   Tester.Instance.InstanceID(),
					},
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
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
				Details: &settings_object.Details{
					ChangeDate: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   Tester.Instance.InstanceID(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We want to have the same response no matter how often we call the function
			Client.SetExecution(tt.ctx, tt.req)
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertSettingsDetails(t, tt.want.Details, got.Details)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_SetExecution_Function(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

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
			ctx:  CTX,
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
			ctx:  CTX,
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
			ctx:  CTX,
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
				Details: &settings_object.Details{
					ChangeDate: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   Tester.Instance.InstanceID(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We want to have the same response no matter how often we call the function
			Client.SetExecution(tt.ctx, tt.req)
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertSettingsDetails(t, tt.want.Details, got.Details)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}
