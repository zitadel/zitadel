//go:build integration

package action_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v2beta"
)

func executionTargetsSingleTarget(id string) []*action.ExecutionTargetType {
	return []*action.ExecutionTargetType{{Type: &action.ExecutionTargetType_Target{Target: id}}}
}

func executionTargetsSingleInclude(include *action.Condition) []*action.ExecutionTargetType {
	return []*action.ExecutionTargetType{{Type: &action.ExecutionTargetType_Include{Include: include}}}
}

func TestServer_SetExecution_Request(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	targetResp := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name        string
		ctx         context.Context
		req         *action.SetExecutionRequest
		wantSetDate bool
		wantErr     bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantSetDate: true,
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantSetDate: true,
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantSetDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We want to have the same response no matter how often we call the function
			creationDate := time.Now().UTC()
			got, err := instance.Client.ActionV2beta.SetExecution(tt.ctx, tt.req)
			setDate := time.Now().UTC()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assertSetExecutionResponse(t, creationDate, setDate, tt.wantSetDate, got)

			// cleanup to not impact other requests
			instance.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func assertSetExecutionResponse(t *testing.T, creationDate, setDate time.Time, expectedSetDate bool, actualResp *action.SetExecutionResponse) {
	if expectedSetDate {
		if !setDate.IsZero() {
			assert.WithinRange(t, actualResp.GetSetDate().AsTime(), creationDate, setDate)
		} else {
			assert.WithinRange(t, actualResp.GetSetDate().AsTime(), creationDate, time.Now().UTC())
		}
	} else {
		assert.Nil(t, actualResp.SetDate)
	}
}

func TestServer_SetExecution_Request_Include(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	targetResp := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)
	executionCond := &action.Condition{
		ConditionType: &action.Condition_Request{
			Request: &action.RequestExecution{
				Condition: &action.RequestExecution_All{
					All: true,
				},
			},
		},
	}
	instance.SetExecution(isolatedIAMOwnerCTX, t,
		executionCond,
		executionTargetsSingleTarget(targetResp.GetId()),
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
	instance.SetExecution(isolatedIAMOwnerCTX, t,
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
	instance.SetExecution(isolatedIAMOwnerCTX, t,
		circularExecutionMethod,
		executionTargetsSingleInclude(circularExecutionService),
	)

	tests := []struct {
		name        string
		ctx         context.Context
		req         *action.SetExecutionRequest
		wantSetDate bool
		wantErr     bool
	}{
		{
			name: "method, circular error",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: circularExecutionService,
				Targets:   executionTargetsSingleInclude(circularExecutionMethod),
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
				Targets: executionTargetsSingleInclude(executionCond),
			},
			wantSetDate: true,
		},
		{
			name: "service, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Service{
								Service: "zitadel.user.v2beta.UserService",
							},
						},
					},
				},
				Targets: executionTargetsSingleInclude(executionCond),
			},
			wantSetDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			got, err := instance.Client.ActionV2beta.SetExecution(tt.ctx, tt.req)
			setDate := time.Now().UTC()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assertSetExecutionResponse(t, creationDate, setDate, tt.wantSetDate, got)

			// cleanup to not impact other requests
			instance.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_SetExecution_Response(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	targetResp := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name        string
		ctx         context.Context
		req         *action.SetExecutionRequest
		wantSetDate bool
		wantErr     bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantSetDate: true,
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantSetDate: true,
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantSetDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			got, err := instance.Client.ActionV2beta.SetExecution(tt.ctx, tt.req)
			setDate := time.Now().UTC()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assertSetExecutionResponse(t, creationDate, setDate, tt.wantSetDate, got)

			// cleanup to not impact other requests
			instance.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_SetExecution_Event(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	targetResp := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name        string
		ctx         context.Context
		req         *action.SetExecutionRequest
		wantSetDate bool
		wantErr     bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantSetDate: true,
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantSetDate: true,
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantSetDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			got, err := instance.Client.ActionV2beta.SetExecution(tt.ctx, tt.req)
			setDate := time.Now().UTC()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assertSetExecutionResponse(t, creationDate, setDate, tt.wantSetDate, got)

			// cleanup to not impact other requests
			instance.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_SetExecution_Function(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	targetResp := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name        string
		ctx         context.Context
		req         *action.SetExecutionRequest
		wantSetDate bool
		wantErr     bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
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
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantErr: true,
		},
		{
			name: "function, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Function{
						Function: &action.FunctionExecution{Name: "presamlresponse"},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantSetDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			got, err := instance.Client.ActionV2beta.SetExecution(tt.ctx, tt.req)
			setDate := time.Now().UTC()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assertSetExecutionResponse(t, creationDate, setDate, tt.wantSetDate, got)

			// cleanup to not impact other requests
			instance.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}
