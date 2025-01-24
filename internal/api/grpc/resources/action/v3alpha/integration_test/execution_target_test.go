//go:build integration

package action_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/grpc/server/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	action "github.com/zitadel/zitadel/pkg/grpc/resources/action/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestServer_ExecutionTarget(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	fullMethod := "/zitadel.resources.action.v3alpha.ZITADELActions/GetTarget"

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(context.Context, *action.GetTargetRequest, *action.GetTargetResponse) (func(), func() bool, error)
		clean   func(context.Context)
		req     *action.GetTargetRequest
		want    *action.GetTargetResponse
		wantErr bool
	}{
		{
			name: "GetTarget, request and response, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) (func(), func() bool, error) {

				orgID := instance.DefaultOrg.Id
				projectID := ""
				userID := instance.Users.Get(integration.UserTypeIAMOwner).ID

				// create target for target changes
				targetCreatedName := gofakeit.Name()
				targetCreatedURL := "https://nonexistent"

				targetCreated := instance.CreateTarget(ctx, t, targetCreatedName, targetCreatedURL, domain.TargetTypeCall, false)

				// request received by target
				wantRequest := &middleware.ContextInfoRequest{FullMethod: fullMethod, InstanceID: instance.ID(), OrgID: orgID, ProjectID: projectID, UserID: userID, Request: request}
				changedRequest := &action.GetTargetRequest{Id: targetCreated.GetDetails().GetId()}
				// replace original request with different targetID
				urlRequest, closeRequest, calledRequest, _ := integration.TestServerCall(wantRequest, 0, http.StatusOK, changedRequest)

				targetRequest := waitForTarget(ctx, t, instance, urlRequest, domain.TargetTypeCall, false)

				waitForExecutionOnCondition(ctx, t, instance, conditionRequestFullMethod(fullMethod), executionTargetsSingleTarget(targetRequest.GetDetails().GetId()))

				// expected response from the GetTarget
				expectedResponse := &action.GetTargetResponse{
					Target: &action.GetTarget{
						Config: &action.Target{
							Name:     targetCreatedName,
							Endpoint: targetCreatedURL,
							TargetType: &action.Target_RestCall{
								RestCall: &action.SetRESTCall{
									InterruptOnError: false,
								},
							},
							Timeout: durationpb.New(5 * time.Second),
						},
						Details: targetCreated.GetDetails(),
					},
				}
				// has to be set separately because of the pointers
				response.Target = &action.GetTarget{
					Details: targetCreated.GetDetails(),
					Config: &action.Target{
						Name: targetCreatedName,
						TargetType: &action.Target_RestCall{
							RestCall: &action.SetRESTCall{
								InterruptOnError: false,
							},
						},
						Timeout:  durationpb.New(5 * time.Second),
						Endpoint: targetCreatedURL,
					},
				}

				// content for partial update
				changedResponse := &action.GetTargetResponse{
					Target: &action.GetTarget{
						Details: &resource_object.Details{
							Id: targetCreated.GetDetails().GetId(),
						},
					},
				}

				// response received by target
				wantResponse := &middleware.ContextInfoResponse{
					FullMethod: fullMethod,
					InstanceID: instance.ID(),
					OrgID:      orgID,
					ProjectID:  projectID,
					UserID:     userID,
					Request:    changedRequest,
					Response:   expectedResponse,
				}
				// after request with different targetID, return changed response
				targetResponseURL, closeResponse, calledResponse, _ := integration.TestServerCall(wantResponse, 0, http.StatusOK, changedResponse)

				targetResponse := waitForTarget(ctx, t, instance, targetResponseURL, domain.TargetTypeCall, false)
				waitForExecutionOnCondition(ctx, t, instance, conditionResponseFullMethod(fullMethod), executionTargetsSingleTarget(targetResponse.GetDetails().GetId()))
				return func() {
						closeRequest()
						closeResponse()
					}, func() bool {
						if calledRequest() != 1 {
							return false
						}
						if calledResponse() != 1 {
							return false
						}
						return true
					}, nil
			},
			clean: func(ctx context.Context) {
				instance.DeleteExecution(ctx, t, conditionRequestFullMethod(fullMethod))
				instance.DeleteExecution(ctx, t, conditionResponseFullMethod(fullMethod))
			},
			req: &action.GetTargetRequest{
				Id: "something",
			},
			want: &action.GetTargetResponse{
				Target: &action.GetTarget{
					Details: &resource_object.Details{
						Id: "changed",
						Owner: &object.Owner{
							Type: object.OwnerType_OWNER_TYPE_INSTANCE,
							Id:   instance.ID(),
						},
					},
				},
			},
		},
		{
			name: "GetTarget, request, interrupt",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) (func(), func() bool, error) {

				fullMethod := "/zitadel.resources.action.v3alpha.ZITADELActions/GetTarget"
				orgID := instance.DefaultOrg.Id
				projectID := ""
				userID := instance.Users.Get(integration.UserTypeIAMOwner).ID

				// request received by target
				wantRequest := &middleware.ContextInfoRequest{FullMethod: fullMethod, InstanceID: instance.ID(), OrgID: orgID, ProjectID: projectID, UserID: userID, Request: request}
				urlRequest, closeRequest, calledRequest, _ := integration.TestServerCall(wantRequest, 0, http.StatusInternalServerError, &action.GetTargetRequest{Id: "notchanged"})

				targetRequest := waitForTarget(ctx, t, instance, urlRequest, domain.TargetTypeCall, true)
				waitForExecutionOnCondition(ctx, t, instance, conditionRequestFullMethod(fullMethod), executionTargetsSingleTarget(targetRequest.GetDetails().GetId()))
				// GetTarget with used target
				request.Id = targetRequest.GetDetails().GetId()
				return func() {
						closeRequest()
					}, func() bool {
						return calledRequest() == 1
					}, nil
			},
			clean: func(ctx context.Context) {
				instance.DeleteExecution(ctx, t, conditionRequestFullMethod(fullMethod))
			},
			req:     &action.GetTargetRequest{},
			wantErr: true,
		},
		{
			name: "GetTarget, response, interrupt",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) (func(), func() bool, error) {

				fullMethod := "/zitadel.resources.action.v3alpha.ZITADELActions/GetTarget"
				orgID := instance.DefaultOrg.Id
				projectID := ""
				userID := instance.Users.Get(integration.UserTypeIAMOwner).ID

				// create target for target changes
				targetCreatedName := gofakeit.Name()
				targetCreatedURL := "https://nonexistent"

				targetCreated := instance.CreateTarget(ctx, t, targetCreatedName, targetCreatedURL, domain.TargetTypeCall, false)

				// GetTarget with used target
				request.Id = targetCreated.GetDetails().GetId()

				// expected response from the GetTarget
				expectedResponse := &action.GetTargetResponse{
					Target: &action.GetTarget{
						Details: targetCreated.GetDetails(),
						Config: &action.Target{
							Name:     targetCreatedName,
							Endpoint: targetCreatedURL,
							TargetType: &action.Target_RestCall{
								RestCall: &action.SetRESTCall{
									InterruptOnError: false,
								},
							},
							Timeout: durationpb.New(5 * time.Second),
						},
					},
				}
				// content for partial update
				changedResponse := &action.GetTargetResponse{
					Target: &action.GetTarget{
						Details: &resource_object.Details{
							Id: "changed",
						},
					},
				}

				// response received by target
				wantResponse := &middleware.ContextInfoResponse{
					FullMethod: fullMethod,
					InstanceID: instance.ID(),
					OrgID:      orgID,
					ProjectID:  projectID,
					UserID:     userID,
					Request:    request,
					Response:   expectedResponse,
				}
				// after request with different targetID, return changed response
				targetResponseURL, closeResponse, calledResponse, _ := integration.TestServerCall(wantResponse, 0, http.StatusInternalServerError, changedResponse)

				targetResponse := waitForTarget(ctx, t, instance, targetResponseURL, domain.TargetTypeCall, true)
				waitForExecutionOnCondition(ctx, t, instance, conditionResponseFullMethod(fullMethod), executionTargetsSingleTarget(targetResponse.GetDetails().GetId()))
				return func() {
						closeResponse()
					}, func() bool {
						return calledResponse() == 1
					}, nil
			},
			clean: func(ctx context.Context) {
				instance.DeleteExecution(ctx, t, conditionResponseFullMethod(fullMethod))
			},
			req:     &action.GetTargetRequest{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			close, called, err := tt.dep(tt.ctx, tt.req, tt.want)
			require.NoError(t, err)
			defer close()

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.ctx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := instance.Client.ActionV3Alpha.GetTarget(tt.ctx, tt.req)
				if tt.wantErr {
					require.Error(ttt, err)
					return
				}
				require.NoError(ttt, err)

				integration.AssertResourceDetails(ttt, tt.want.GetTarget().GetDetails(), got.GetTarget().GetDetails())
				tt.want.Target.Details = got.GetTarget().GetDetails()
				assert.EqualExportedValues(ttt, tt.want.GetTarget().GetConfig(), got.GetTarget().GetConfig())

			}, retryDuration, tick, "timeout waiting for expected execution result")

			if tt.clean != nil {
				tt.clean(tt.ctx)
			}
			require.True(t, called())
		})
	}
}

func TestServer_ExecutionTarget_Event(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	event := "session.added"
	urlRequest, closeF, calledF, resetF := integration.TestServerCall(nil, 0, http.StatusOK, nil)
	defer closeF()

	targetRequest := instance.CreateTarget(isolatedIAMOwnerCTX, t, "session-added-target-"+gofakeit.AppName(), urlRequest, domain.TargetTypeWebhook, true)
	instance.SetExecution(isolatedIAMOwnerCTX, t, conditionEvent(event), executionTargetsSingleTarget(targetRequest.GetDetails().GetId()))
	defer func() {
		instance.DeleteExecution(isolatedIAMOwnerCTX, t, conditionEvent(event))
	}()

	tests := []struct {
		name          string
		ctx           context.Context
		eventCount    int
		expectedCalls int
		clean         func(context.Context)
		wantErr       bool
	}{
		{
			name:          "event, 1 session.added, ok",
			ctx:           isolatedIAMOwnerCTX,
			eventCount:    1,
			expectedCalls: 1,
		},
		{
			name:          "event, 10 session.added, ok",
			ctx:           isolatedIAMOwnerCTX,
			eventCount:    10,
			expectedCalls: 10,
		},
		{
			name:          "event, 100 session.added, ok",
			ctx:           isolatedIAMOwnerCTX,
			eventCount:    100,
			expectedCalls: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset the count of the target
			resetF()

			for i := 0; i < tt.eventCount; i++ {
				_, err := instance.Client.SessionV2.CreateSession(tt.ctx, &session.CreateSessionRequest{})
				require.NoError(t, err)
			}

			// wait for called target
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.ctx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				assert.True(ttt, calledF() == tt.expectedCalls)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func TestServer_ExecutionTarget_Event_10Sec(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	event := "session.added"
	// call takes longer than timeout of target
	urlRequest, closeF, calledF, resetF := integration.TestServerCall(nil, 10*time.Second, http.StatusOK, nil)
	defer closeF()

	targetRequest := instance.CreateTarget(isolatedIAMOwnerCTX, t, "session-added-target-"+gofakeit.AppName(), urlRequest, domain.TargetTypeWebhook, true)
	instance.SetExecution(isolatedIAMOwnerCTX, t, conditionEvent(event), executionTargetsSingleTarget(targetRequest.GetDetails().GetId()))
	defer func() {
		instance.DeleteExecution(isolatedIAMOwnerCTX, t, conditionEvent(event))
	}()

	tests := []struct {
		name          string
		ctx           context.Context
		eventCount    int
		expectedCalls int
		clean         func(context.Context)
		wantErr       bool
	}{
		{
			name:          "event, 1 session.added, error logs",
			ctx:           isolatedIAMOwnerCTX,
			eventCount:    1,
			expectedCalls: 1,
		},
		{
			name:          "event, 5 session.added, error logs",
			ctx:           isolatedIAMOwnerCTX,
			eventCount:    5,
			expectedCalls: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset the count of the target
			resetF()

			for i := 0; i < tt.eventCount; i++ {
				_, err := instance.Client.SessionV2.CreateSession(tt.ctx, &session.CreateSessionRequest{})
				require.NoError(t, err)
			}

			// wait for called target
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.ctx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				assert.True(ttt, calledF() == tt.expectedCalls)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func TestServer_ExecutionTarget_Event_1Sec(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	event := "session.added"
	urlRequest, closeF, calledF, resetF := integration.TestServerCall(nil, 1*time.Second, http.StatusOK, nil)
	defer closeF()

	targetRequest := instance.CreateTarget(isolatedIAMOwnerCTX, t, "session-added-target-"+gofakeit.AppName(), urlRequest, domain.TargetTypeWebhook, true)
	instance.SetExecution(isolatedIAMOwnerCTX, t, conditionEvent(event), executionTargetsSingleTarget(targetRequest.GetDetails().GetId()))
	defer func() {
		instance.DeleteExecution(isolatedIAMOwnerCTX, t, conditionEvent(event))
	}()

	tests := []struct {
		name          string
		ctx           context.Context
		eventCount    int
		expectedCalls int
		clean         func(context.Context)
		wantErr       bool
	}{
		{
			name:          "event, 1 session.added, ok",
			ctx:           isolatedIAMOwnerCTX,
			eventCount:    1,
			expectedCalls: 1,
		},
		{
			name:          "event, 5 session.added, ok",
			ctx:           isolatedIAMOwnerCTX,
			eventCount:    5,
			expectedCalls: 5,
		},
		{
			name:          "event, 10 session.added, ok",
			ctx:           isolatedIAMOwnerCTX,
			eventCount:    10,
			expectedCalls: 10,
		},
		{
			name:          "event, 20 session.added, ok",
			ctx:           isolatedIAMOwnerCTX,
			eventCount:    20,
			expectedCalls: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset the count of the target
			resetF()

			for i := 0; i < tt.eventCount; i++ {
				_, err := instance.Client.SessionV2.CreateSession(tt.ctx, &session.CreateSessionRequest{})
				require.NoError(t, err)
			}

			// wait for called target
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.ctx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				assert.True(ttt, calledF() == tt.expectedCalls)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func waitForExecutionOnCondition(ctx context.Context, t *testing.T, instance *integration.Instance, condition *action.Condition, targets []*action.ExecutionTargetType) {
	instance.SetExecution(ctx, t, condition, targets)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		got, err := instance.Client.ActionV3Alpha.SearchExecutions(ctx, &action.SearchExecutionsRequest{
			Filters: []*action.ExecutionSearchFilter{
				{Filter: &action.ExecutionSearchFilter_InConditionsFilter{
					InConditionsFilter: &action.InConditionsFilter{Conditions: []*action.Condition{condition}},
				}},
			},
		})
		if !assert.NoError(ttt, err) {
			return
		}
		if !assert.Len(ttt, got.GetResult(), 1) {
			return
		}
		gotTargets := got.GetResult()[0].GetExecution().GetTargets()
		// always first check length, otherwise its failed anyway
		if assert.Len(ttt, gotTargets, len(targets)) {
			for i := range targets {
				assert.EqualExportedValues(ttt, targets[i].GetType(), gotTargets[i].GetType())
			}
		}
	}, retryDuration, tick, "timeout waiting for expected execution result")
	return
}

func waitForTarget(ctx context.Context, t *testing.T, instance *integration.Instance, endpoint string, ty domain.TargetType, interrupt bool) *action.CreateTargetResponse {
	resp := instance.CreateTarget(ctx, t, "", endpoint, ty, interrupt)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		got, err := instance.Client.ActionV3Alpha.SearchTargets(ctx, &action.SearchTargetsRequest{
			Filters: []*action.TargetSearchFilter{
				{Filter: &action.TargetSearchFilter_InTargetIdsFilter{
					InTargetIdsFilter: &action.InTargetIDsFilter{TargetIds: []string{resp.GetDetails().GetId()}},
				}},
			},
		})
		if !assert.NoError(ttt, err) {
			return
		}
		if !assert.Len(ttt, got.GetResult(), 1) {
			return
		}
		config := got.GetResult()[0].GetConfig()
		assert.Equal(ttt, config.GetEndpoint(), endpoint)
		switch ty {
		case domain.TargetTypeWebhook:
			if !assert.NotNil(ttt, config.GetRestWebhook()) {
				return
			}
			assert.Equal(ttt, interrupt, config.GetRestWebhook().GetInterruptOnError())
		case domain.TargetTypeAsync:
			assert.NotNil(ttt, config.GetRestAsync())
		case domain.TargetTypeCall:
			if !assert.NotNil(ttt, config.GetRestCall()) {
				return
			}
			assert.Equal(ttt, interrupt, config.GetRestCall().GetInterruptOnError())
		}
	}, retryDuration, tick, "timeout waiting for expected execution result")
	return resp
}

func conditionRequestFullMethod(fullMethod string) *action.Condition {
	return &action.Condition{
		ConditionType: &action.Condition_Request{
			Request: &action.RequestExecution{
				Condition: &action.RequestExecution_Method{
					Method: fullMethod,
				},
			},
		},
	}
}

func conditionResponseFullMethod(fullMethod string) *action.Condition {
	return &action.Condition{
		ConditionType: &action.Condition_Response{
			Response: &action.ResponseExecution{
				Condition: &action.ResponseExecution_Method{
					Method: fullMethod,
				},
			},
		},
	}
}

func conditionEvent(event string) *action.Condition {
	return &action.Condition{
		ConditionType: &action.Condition_Event{
			Event: &action.EventExecution{
				Condition: &action.EventExecution_Event{
					Event: event,
				},
			},
		},
	}
}
