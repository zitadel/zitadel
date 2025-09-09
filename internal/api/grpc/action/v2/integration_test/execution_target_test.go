//go:build integration

package action_test

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/grpc/server/middleware"
	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	saml_api "github.com/zitadel/zitadel/internal/api/saml"
	"github.com/zitadel/zitadel/internal/domain"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/action/v2"
	"github.com/zitadel/zitadel/pkg/grpc/app"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/metadata"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
	saml_pb "github.com/zitadel/zitadel/pkg/grpc/saml/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

const (
	redirectURIImplicit = "http://localhost:9999/callback"
)

var (
	loginV2 = &app.LoginVersion{Version: &app.LoginVersion_LoginV2{LoginV2: &app.LoginV2{BaseUri: nil}}}
)

func TestServer_ExecutionTarget(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	fullMethod := action.ActionService_GetTarget_FullMethodName

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(context.Context, *action.GetTargetRequest, *action.GetTargetResponse) (closeF func(), calledF func() bool)
		clean   func(context.Context)
		req     *action.GetTargetRequest
		want    *action.GetTargetResponse
		wantErr bool
	}{
		{
			name: "GetTarget, request and response, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) (func(), func() bool) {

				orgID := instance.DefaultOrg.Id
				projectID := ""
				userID := instance.Users.Get(integration.UserTypeIAMOwner).ID

				// create target for target changes
				targetCreatedName := integration.TargetName()
				targetCreatedURL := "https://nonexistent"

				targetCreated := instance.CreateTarget(ctx, t, targetCreatedName, targetCreatedURL, target_domain.TargetTypeCall, false)

				// request received by target
				wantRequest := &middleware.ContextInfoRequest{
					FullMethod: fullMethod,
					InstanceID: instance.ID(),
					OrgID:      orgID,
					ProjectID:  projectID,
					UserID:     userID,
					Request:    middleware.Message{Message: request},
					Headers:    map[string][]string{"Content-Type": {"application/grpc"}, "Host": {instance.Host()}},
				}
				changedRequest := &action.GetTargetRequest{Id: targetCreated.GetId()}
				// replace original request with different targetID
				urlRequest, closeRequest, calledRequest, _ := integration.TestServerCallProto(wantRequest, 0, http.StatusOK, changedRequest)

				targetRequest := waitForTarget(ctx, t, instance, urlRequest, target_domain.TargetTypeCall, false)

				waitForExecutionOnCondition(ctx, t, instance, conditionRequestFullMethod(fullMethod), []string{targetRequest.GetId()})

				// expected response from the GetTarget
				expectedResponse := &action.GetTargetResponse{
					Target: &action.Target{
						Id:           targetCreated.GetId(),
						CreationDate: targetCreated.GetCreationDate(),
						ChangeDate:   targetCreated.GetCreationDate(),
						Name:         targetCreatedName,
						TargetType: &action.Target_RestCall{
							RestCall: &action.RESTCall{
								InterruptOnError: false,
							},
						},
						Timeout:    durationpb.New(5 * time.Second),
						Endpoint:   targetCreatedURL,
						SigningKey: targetCreated.GetSigningKey(),
					},
				}

				changedResponse := &action.GetTargetResponse{
					Target: &action.Target{
						Id:           "changed",
						CreationDate: targetCreated.GetCreationDate(),
						ChangeDate:   targetCreated.GetCreationDate(),
						Name:         targetCreatedName,
						TargetType: &action.Target_RestCall{
							RestCall: &action.RESTCall{
								InterruptOnError: false,
							},
						},
						Timeout:    durationpb.New(5 * time.Second),
						Endpoint:   targetCreatedURL,
						SigningKey: targetCreated.GetSigningKey(),
					},
				}
				// content for update
				response.Target = &action.Target{
					Id:           "changed",
					CreationDate: targetCreated.GetCreationDate(),
					ChangeDate:   targetCreated.GetCreationDate(),
					Name:         targetCreatedName,
					TargetType: &action.Target_RestCall{
						RestCall: &action.RESTCall{
							InterruptOnError: false,
						},
					},
					Timeout:    durationpb.New(5 * time.Second),
					Endpoint:   targetCreatedURL,
					SigningKey: targetCreated.GetSigningKey(),
				}

				// response received by target
				wantResponse := &middleware.ContextInfoResponse{
					FullMethod: fullMethod,
					InstanceID: instance.ID(),
					OrgID:      orgID,
					ProjectID:  projectID,
					UserID:     userID,
					Request:    middleware.Message{Message: changedRequest},
					Response:   middleware.Message{Message: expectedResponse},
					Headers:    map[string][]string{"Content-Type": {"application/grpc"}, "Host": {instance.Host()}},
				}
				// after request with different targetID, return changed response
				targetResponseURL, closeResponse, calledResponse, _ := integration.TestServerCallProto(wantResponse, 0, http.StatusOK, changedResponse)

				targetResponse := waitForTarget(ctx, t, instance, targetResponseURL, target_domain.TargetTypeCall, false)
				waitForExecutionOnCondition(ctx, t, instance, conditionResponseFullMethod(fullMethod), []string{targetResponse.GetId()})
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
					}
			},
			clean: func(ctx context.Context) {
				instance.DeleteExecution(ctx, t, conditionRequestFullMethod(fullMethod))
				instance.DeleteExecution(ctx, t, conditionResponseFullMethod(fullMethod))
			},
			req: &action.GetTargetRequest{
				Id: "something",
			},
			want: &action.GetTargetResponse{
				// defined in the dependency function
			},
		},
		{
			name: "GetTarget, request, interrupt",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) (func(), func() bool) {
				orgID := instance.DefaultOrg.Id
				projectID := ""
				userID := instance.Users.Get(integration.UserTypeIAMOwner).ID

				// request received by target
				wantRequest := &middleware.ContextInfoRequest{FullMethod: fullMethod, InstanceID: instance.ID(), OrgID: orgID, ProjectID: projectID, UserID: userID, Request: middleware.Message{Message: request}}
				urlRequest, closeRequest, calledRequest, _ := integration.TestServerCallProto(wantRequest, 0, http.StatusInternalServerError, nil)

				targetRequest := waitForTarget(ctx, t, instance, urlRequest, target_domain.TargetTypeCall, true)
				waitForExecutionOnCondition(ctx, t, instance, conditionRequestFullMethod(fullMethod), []string{targetRequest.GetId()})
				// GetTarget with used target
				request.Id = targetRequest.GetId()
				return func() {
						closeRequest()
					}, func() bool {
						return calledRequest() == 1
					}
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
			dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) (func(), func() bool) {
				orgID := instance.DefaultOrg.Id
				projectID := ""
				userID := instance.Users.Get(integration.UserTypeIAMOwner).ID

				// create target for target changes
				targetCreatedName := integration.TargetName()
				targetCreatedURL := "https://nonexistent"

				targetCreated := instance.CreateTarget(ctx, t, targetCreatedName, targetCreatedURL, target_domain.TargetTypeCall, false)

				// GetTarget with used target
				request.Id = targetCreated.GetId()

				// expected response from the GetTarget
				expectedResponse := &action.GetTargetResponse{
					Target: &action.Target{
						Id:           targetCreated.GetId(),
						CreationDate: targetCreated.GetCreationDate(),
						ChangeDate:   targetCreated.GetCreationDate(),
						Name:         targetCreatedName,
						TargetType: &action.Target_RestCall{
							RestCall: &action.RESTCall{
								InterruptOnError: false,
							},
						},
						Timeout:    durationpb.New(5 * time.Second),
						Endpoint:   targetCreatedURL,
						SigningKey: targetCreated.GetSigningKey(),
					},
				}

				// response received by target
				wantResponse := &middleware.ContextInfoResponse{
					FullMethod: fullMethod,
					InstanceID: instance.ID(),
					OrgID:      orgID,
					ProjectID:  projectID,
					UserID:     userID,
					Request:    middleware.Message{Message: request},
					Response:   middleware.Message{Message: expectedResponse},
				}
				// after request with different targetID, return changed response
				targetResponseURL, closeResponse, calledResponse, _ := integration.TestServerCallProto(wantResponse, 0, http.StatusInternalServerError, nil)

				targetResponse := waitForTarget(ctx, t, instance, targetResponseURL, target_domain.TargetTypeCall, true)
				waitForExecutionOnCondition(ctx, t, instance, conditionResponseFullMethod(fullMethod), []string{targetResponse.GetId()})
				return func() {
						closeResponse()
					}, func() bool {
						return calledResponse() == 1
					}
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
			closeF, calledF := tt.dep(tt.ctx, tt.req, tt.want)
			defer closeF()

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.ctx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := instance.Client.ActionV2.GetTarget(tt.ctx, tt.req)
				if tt.wantErr {
					require.Error(ttt, err)
					return
				}
				require.NoError(ttt, err)
				assert.EqualExportedValues(ttt, tt.want.GetTarget(), got.GetTarget())

			}, retryDuration, tick, "timeout waiting for expected execution result")

			if tt.clean != nil {
				tt.clean(tt.ctx)
			}
			require.True(t, calledF())
		})
	}
}

func TestServer_ExecutionTarget_Event(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	event := "session.added"
	urlRequest, closeF, calledF, resetF := integration.TestServerCall(nil, 0, http.StatusOK, nil)
	defer closeF()

	targetResponse := waitForTarget(isolatedIAMOwnerCTX, t, instance, urlRequest, target_domain.TargetTypeWebhook, true)
	waitForExecutionOnCondition(isolatedIAMOwnerCTX, t, instance, conditionEvent(event), []string{targetResponse.GetId()})

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
			name:          "event, 50 session.added, ok",
			ctx:           isolatedIAMOwnerCTX,
			eventCount:    50,
			expectedCalls: 50,
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

func TestServer_ExecutionTarget_Event_LongerThanTargetTimeout(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	event := "session.added"
	// call takes longer than timeout of target
	urlRequest, closeF, calledF, resetF := integration.TestServerCall(nil, 5*time.Second, http.StatusOK, nil)
	defer closeF()

	targetResponse := waitForTarget(isolatedIAMOwnerCTX, t, instance, urlRequest, target_domain.TargetTypeWebhook, true)
	waitForExecutionOnCondition(isolatedIAMOwnerCTX, t, instance, conditionEvent(event), []string{targetResponse.GetId()})

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

func TestServer_ExecutionTarget_Event_LongerThanTransactionTimeout(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	event := "session.added"
	urlRequest, closeF, calledF, resetF := integration.TestServerCall(nil, 1*time.Second, http.StatusOK, nil)
	defer closeF()

	targetResponse := waitForTarget(isolatedIAMOwnerCTX, t, instance, urlRequest, target_domain.TargetTypeWebhook, true)
	waitForExecutionOnCondition(isolatedIAMOwnerCTX, t, instance, conditionEvent(event), []string{targetResponse.GetId()})

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
			name:          "event, 5 session.added, ok",
			ctx:           isolatedIAMOwnerCTX,
			eventCount:    5,
			expectedCalls: 5,
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

func waitForExecutionOnCondition(ctx context.Context, t *testing.T, instance *integration.Instance, condition *action.Condition, targets []string) {
	instance.SetExecution(ctx, t, condition, targets)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		got, err := instance.Client.ActionV2.ListExecutions(ctx, &action.ListExecutionsRequest{
			Filters: []*action.ExecutionSearchFilter{
				{Filter: &action.ExecutionSearchFilter_InConditionsFilter{
					InConditionsFilter: &action.InConditionsFilter{Conditions: []*action.Condition{condition}},
				}},
			},
		})
		if !assert.NoError(ttt, err) {
			return
		}
		if !assert.Len(ttt, got.GetExecutions(), 1) {
			return
		}
		gotTargets := got.GetExecutions()[0].GetTargets()
		// always first check length, otherwise its failed anyway
		if assert.Len(ttt, gotTargets, len(targets)) {
			for i := range targets {
				assert.EqualExportedValues(ttt, targets[i], gotTargets[i])
			}
		}
	}, retryDuration, tick, "timeout waiting for expected execution result")
}

func waitForTarget(ctx context.Context, t *testing.T, instance *integration.Instance, endpoint string, ty target_domain.TargetType, interrupt bool) *action.CreateTargetResponse {
	resp := instance.CreateTarget(ctx, t, "", endpoint, ty, interrupt)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		got, err := instance.Client.ActionV2.ListTargets(ctx, &action.ListTargetsRequest{
			Filters: []*action.TargetSearchFilter{
				{Filter: &action.TargetSearchFilter_InTargetIdsFilter{
					InTargetIdsFilter: &action.InTargetIDsFilter{TargetIds: []string{resp.GetId()}},
				}},
			},
		})
		if !assert.NoError(ttt, err) {
			return
		}
		if !assert.Len(ttt, got.GetTargets(), 1) {
			return
		}
		config := got.GetTargets()[0]
		assert.Equal(ttt, config.GetEndpoint(), endpoint)
		switch ty {
		case target_domain.TargetTypeWebhook:
			if !assert.NotNil(ttt, config.GetRestWebhook()) {
				return
			}
			assert.Equal(ttt, interrupt, config.GetRestWebhook().GetInterruptOnError())
		case target_domain.TargetTypeAsync:
			assert.NotNil(ttt, config.GetRestAsync())
		case target_domain.TargetTypeCall:
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

func conditionFunction(function string) *action.Condition {
	return &action.Condition{
		ConditionType: &action.Condition_Function{
			Function: &action.FunctionExecution{
				Name: function,
			},
		},
	}
}

func TestServer_ExecutionTargetPreUserinfo(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	ctxLoginClient := instance.WithAuthorizationToken(CTX, integration.UserTypeLogin)

	client, err := instance.CreateOIDCImplicitFlowClient(isolatedIAMCtx, t, redirectURIImplicit, loginV2)
	require.NoError(t, err)

	type want struct {
		addedClaims     map[string]any
		addedLogClaims  map[string][]string
		setUserMetadata []*metadata.Metadata
	}
	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, t *testing.T, req *oidc_pb.CreateCallbackRequest) (string, func())
		req     *oidc_pb.CreateCallbackRequest
		want    want
		wantErr bool
	}{
		{
			name: "append claim",
			ctx:  ctxLoginClient,
			dep: func(ctx context.Context, t *testing.T, req *oidc_pb.CreateCallbackRequest) (string, func()) {
				response := &oidc_api.ContextInfoResponse{
					AppendClaims: []*oidc_api.AppendClaim{
						{Key: "added", Value: "value"},
					},
				}
				return expectPreUserinfoExecution(ctx, t, instance, client.GetClientId(), req, response)
			},
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := instance.CreateOIDCAuthRequestImplicitWithoutLoginClientHeader(isolatedIAMCtx, client.GetClientId(), redirectURIImplicit)
					require.NoError(t, err)
					return authRequestID
				}(),
			},
			want: want{
				addedClaims: map[string]any{
					"added": "value",
				},
			},
			wantErr: false,
		},
		{
			name: "append log claim",
			ctx:  ctxLoginClient,
			dep: func(ctx context.Context, t *testing.T, req *oidc_pb.CreateCallbackRequest) (string, func()) {
				response := &oidc_api.ContextInfoResponse{
					AppendLogClaims: []string{
						"addedLog",
					},
				}
				return expectPreUserinfoExecution(ctx, t, instance, client.GetClientId(), req, response)
			},
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := instance.CreateOIDCAuthRequestImplicitWithoutLoginClientHeader(isolatedIAMCtx, client.GetClientId(), redirectURIImplicit)
					require.NoError(t, err)
					return authRequestID
				}(),
			},
			want: want{
				addedLogClaims: map[string][]string{
					"urn:zitadel:iam:action:function/preuserinfo:log": {"addedLog"},
				},
			},
			wantErr: false,
		},
		{
			name: "set user metadata",
			ctx:  ctxLoginClient,
			dep: func(ctx context.Context, t *testing.T, req *oidc_pb.CreateCallbackRequest) (string, func()) {
				response := &oidc_api.ContextInfoResponse{
					SetUserMetadata: []*domain.Metadata{
						{Key: "key", Value: []byte("value")},
					},
				}
				return expectPreUserinfoExecution(ctx, t, instance, client.GetClientId(), req, response)
			},
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := instance.CreateOIDCAuthRequestImplicitWithoutLoginClientHeader(isolatedIAMCtx, client.GetClientId(), redirectURIImplicit)
					require.NoError(t, err)
					return authRequestID
				}(),
			},
			want: want{
				setUserMetadata: []*metadata.Metadata{
					{Key: "key", Value: []byte("value")},
				},
			},
			wantErr: false,
		},
		{
			name: "full usage",
			ctx:  ctxLoginClient,
			dep: func(ctx context.Context, t *testing.T, req *oidc_pb.CreateCallbackRequest) (string, func()) {
				response := &oidc_api.ContextInfoResponse{
					SetUserMetadata: []*domain.Metadata{
						{Key: "key1", Value: []byte("value1")},
						{Key: "key2", Value: []byte("value2")},
						{Key: "key3", Value: []byte("value3")},
					},
					AppendLogClaims: []string{
						"addedLog1",
						"addedLog2",
						"addedLog3",
					},
					AppendClaims: []*oidc_api.AppendClaim{
						{Key: "added1", Value: "value1"},
						{Key: "added2", Value: "value2"},
						{Key: "added3", Value: "value3"},
					},
				}
				return expectPreUserinfoExecution(ctx, t, instance, client.GetClientId(), req, response)
			},
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := instance.CreateOIDCAuthRequestImplicitWithoutLoginClientHeader(isolatedIAMCtx, client.GetClientId(), redirectURIImplicit)
					require.NoError(t, err)
					return authRequestID
				}(),
			},
			want: want{
				addedClaims: map[string]any{
					"added1": "value1",
					"added2": "value2",
					"added3": "value3",
				},
				setUserMetadata: []*metadata.Metadata{
					{Key: "key1", Value: []byte("value1")},
					{Key: "key2", Value: []byte("value2")},
					{Key: "key3", Value: []byte("value3")},
				},
				addedLogClaims: map[string][]string{
					"urn:zitadel:iam:action:function/preuserinfo:log": {"addedLog1", "addedLog2", "addedLog3"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, closeF := tt.dep(isolatedIAMCtx, t, tt.req)
			defer closeF()

			got, err := instance.Client.OIDCv2.CreateCallback(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			callbackUrl, err := url.Parse(strings.Replace(got.GetCallbackUrl(), "#", "?", 1))
			require.NoError(t, err)
			claims := getIDTokenClaimsFromCallbackURL(tt.ctx, t, instance, client.GetClientId(), callbackUrl)

			for k, v := range tt.want.addedClaims {
				value, ok := claims[k]
				if !assert.True(t, ok) {
					return
				}
				assert.Equal(t, v, value)
			}
			for k, v := range tt.want.addedLogClaims {
				value, ok := claims[k]
				if !assert.True(t, ok) {
					return
				}
				assert.ElementsMatch(t, v, value)
			}
			if len(tt.want.setUserMetadata) > 0 {
				checkForSetMetadata(isolatedIAMCtx, t, instance, userID, tt.want.setUserMetadata)
			}
		})
	}
}

func expectPreUserinfoExecution(ctx context.Context, t *testing.T, instance *integration.Instance, clientID string, req *oidc_pb.CreateCallbackRequest, response *oidc_api.ContextInfoResponse) (string, func()) {
	userEmail := integration.Email()
	userPhone := integration.Phone()
	userResp := instance.CreateHumanUserVerified(ctx, instance.DefaultOrg.Id, userEmail, userPhone)

	sessionResp := createSession(ctx, t, instance, userResp.GetUserId())
	req.CallbackKind = &oidc_pb.CreateCallbackRequest_Session{
		Session: &oidc_pb.Session{
			SessionId:    sessionResp.GetSessionId(),
			SessionToken: sessionResp.GetSessionToken(),
		},
	}
	expectedContextInfo := contextInfoForUserOIDC(instance, "function/preuserinfo", clientID, userResp, userEmail, userPhone)

	targetURL, closeF, _, _ := integration.TestServerCall(expectedContextInfo, 0, http.StatusOK, response)

	targetResp := waitForTarget(ctx, t, instance, targetURL, target_domain.TargetTypeCall, true)
	waitForExecutionOnCondition(ctx, t, instance, conditionFunction("preuserinfo"), []string{targetResp.GetId()})
	return userResp.GetUserId(), closeF
}

func createSession(ctx context.Context, t *testing.T, instance *integration.Instance, userID string) *session.CreateSessionResponse {
	sessionResp, err := instance.Client.SessionV2.CreateSession(ctx, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: userID,
				},
			},
		},
	})
	require.NoError(t, err)
	return sessionResp
}

func checkForSetMetadata(ctx context.Context, t *testing.T, instance *integration.Instance, userID string, metadataExpected []*metadata.Metadata) {
	integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	assert.EventuallyWithT(t, func(ct *assert.CollectT) {
		metadataResp, err := instance.Client.Mgmt.ListUserMetadata(ctx, &management.ListUserMetadataRequest{Id: userID})
		if !assert.NoError(ct, err) {
			return
		}
		for _, dataExpected := range metadataExpected {
			found := false
			for _, dataCheck := range metadataResp.GetResult() {
				if dataExpected.Key == dataCheck.Key {
					found = true
					if !assert.Equal(ct, dataExpected.Value, dataCheck.Value) {
						return
					}
				}
			}
			if !assert.True(ct, found) {
				return
			}
		}
	}, retryDuration, tick)
}

func getIDTokenClaimsFromCallbackURL(ctx context.Context, t *testing.T, instance *integration.Instance, clientID string, callbackURL *url.URL) map[string]any {
	accessToken := callbackURL.Query().Get("access_token")
	idToken := callbackURL.Query().Get("id_token")

	provider, err := instance.CreateRelyingParty(ctx, clientID, redirectURIImplicit, oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopePhone)
	require.NoError(t, err)
	claims, err := rp.VerifyTokens[*oidc.IDTokenClaims](context.Background(), accessToken, idToken, provider.IDTokenVerifier())
	require.NoError(t, err)
	return claims.Claims
}

type CustomAccessTokenClaims struct {
	oidc.TokenClaims
	Added1 string   `json:"added1,omitempty"`
	Added2 string   `json:"added2,omitempty"`
	Added3 string   `json:"added3,omitempty"`
	Log    []string `json:"urn:zitadel:iam:action:function/preaccesstoken:log,omitempty"`
}

func getAccessTokenClaims(ctx context.Context, t *testing.T, instance *integration.Instance, callbackURL *url.URL) *CustomAccessTokenClaims {
	accessToken := callbackURL.Query().Get("access_token")

	verifier := op.NewAccessTokenVerifier(instance.OIDCIssuer(), rp.NewRemoteKeySet(http.DefaultClient, instance.OIDCIssuer()+"/oauth/v2/keys"))

	claims, err := op.VerifyAccessToken[*CustomAccessTokenClaims](ctx, accessToken, verifier)
	require.NoError(t, err)
	return claims
}

func contextInfoForUserOIDC(instance *integration.Instance, function string, clientID string, userResp *user.AddHumanUserResponse, email, phone string) *oidc_api.ContextInfo {
	return &oidc_api.ContextInfo{
		Function: function,
		UserInfo: &oidc.UserInfo{
			Subject: userResp.GetUserId(),
		},
		User: &query.User{
			ID:                 userResp.GetUserId(),
			CreationDate:       userResp.Details.ChangeDate.AsTime(),
			ChangeDate:         userResp.Details.ChangeDate.AsTime(),
			ResourceOwner:      instance.DefaultOrg.GetId(),
			Sequence:           userResp.Details.Sequence,
			State:              1,
			Username:           email,
			PreferredLoginName: email,
			Human: &query.Human{
				FirstName:              "Mickey",
				LastName:               "Mouse",
				NickName:               "Mickey",
				DisplayName:            "Mickey Mouse",
				AvatarKey:              "",
				PreferredLanguage:      language.Dutch,
				Gender:                 2,
				Email:                  domain.EmailAddress(email),
				IsEmailVerified:        true,
				Phone:                  domain.PhoneNumber(phone),
				IsPhoneVerified:        true,
				PasswordChangeRequired: false,
				PasswordChanged:        time.Time{},
				MFAInitSkipped:         time.Time{},
			},
		},
		UserMetadata: nil,
		Application: &oidc_api.ContextInfoApplication{
			ClientID: clientID,
		},
		Org: &query.UserInfoOrg{
			ID:            instance.DefaultOrg.GetId(),
			Name:          instance.DefaultOrg.GetName(),
			PrimaryDomain: instance.DefaultOrg.GetPrimaryDomain(),
		},
		UserGrants: nil,
		Response:   nil,
	}
}

func TestServer_ExecutionTargetPreAccessToken(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	ctxLoginClient := instance.WithAuthorizationToken(CTX, integration.UserTypeLogin)

	client, err := instance.CreateOIDCImplicitFlowClient(isolatedIAMCtx, t, redirectURIImplicit, loginV2)
	require.NoError(t, err)

	type want struct {
		addedClaims     *CustomAccessTokenClaims
		addedLogClaims  map[string][]string
		setUserMetadata []*metadata.Metadata
	}
	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, t *testing.T, req *oidc_pb.CreateCallbackRequest) (string, func())
		req     *oidc_pb.CreateCallbackRequest
		want    want
		wantErr bool
	}{
		{
			name: "append claim",
			ctx:  ctxLoginClient,
			dep: func(ctx context.Context, t *testing.T, req *oidc_pb.CreateCallbackRequest) (string, func()) {
				response := &oidc_api.ContextInfoResponse{
					AppendClaims: []*oidc_api.AppendClaim{
						{Key: "added1", Value: "value"},
					},
				}
				return expectPreAccessTokenExecution(ctx, t, instance, client.GetClientId(), req, response)
			},
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := instance.CreateOIDCAuthRequestImplicitWithoutLoginClientHeader(isolatedIAMCtx, client.GetClientId(), redirectURIImplicit)
					require.NoError(t, err)
					return authRequestID
				}(),
			},
			want: want{
				addedClaims: &CustomAccessTokenClaims{
					Added1: "value",
				},
			},
			wantErr: false,
		},
		{
			name: "append log claim",
			ctx:  ctxLoginClient,
			dep: func(ctx context.Context, t *testing.T, req *oidc_pb.CreateCallbackRequest) (string, func()) {
				response := &oidc_api.ContextInfoResponse{
					AppendLogClaims: []string{
						"addedLog",
					},
				}
				return expectPreAccessTokenExecution(ctx, t, instance, client.GetClientId(), req, response)
			},
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := instance.CreateOIDCAuthRequestImplicitWithoutLoginClientHeader(isolatedIAMCtx, client.GetClientId(), redirectURIImplicit)
					require.NoError(t, err)
					return authRequestID
				}(),
			},
			want: want{
				addedClaims: &CustomAccessTokenClaims{
					Log: []string{"addedLog"},
				},
			},
			wantErr: false,
		},
		{
			name: "set user metadata",
			ctx:  ctxLoginClient,
			dep: func(ctx context.Context, t *testing.T, req *oidc_pb.CreateCallbackRequest) (string, func()) {
				response := &oidc_api.ContextInfoResponse{
					SetUserMetadata: []*domain.Metadata{
						{Key: "key", Value: []byte("value")},
					},
				}
				return expectPreAccessTokenExecution(ctx, t, instance, client.GetClientId(), req, response)
			},
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := instance.CreateOIDCAuthRequestImplicitWithoutLoginClientHeader(isolatedIAMCtx, client.GetClientId(), redirectURIImplicit)
					require.NoError(t, err)
					return authRequestID
				}(),
			},
			want: want{
				setUserMetadata: []*metadata.Metadata{
					{Key: "key", Value: []byte("value")},
				},
			},
			wantErr: false,
		},
		{
			name: "full usage",
			ctx:  ctxLoginClient,
			dep: func(ctx context.Context, t *testing.T, req *oidc_pb.CreateCallbackRequest) (string, func()) {
				response := &oidc_api.ContextInfoResponse{
					SetUserMetadata: []*domain.Metadata{
						{Key: "key1", Value: []byte("value1")},
						{Key: "key2", Value: []byte("value2")},
						{Key: "key3", Value: []byte("value3")},
					},
					AppendLogClaims: []string{
						"addedLog1",
						"addedLog2",
						"addedLog3",
					},
					AppendClaims: []*oidc_api.AppendClaim{
						{Key: "added1", Value: "value1"},
						{Key: "added2", Value: "value2"},
						{Key: "added3", Value: "value3"},
					},
				}
				return expectPreAccessTokenExecution(ctx, t, instance, client.GetClientId(), req, response)
			},
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := instance.CreateOIDCAuthRequestImplicitWithoutLoginClientHeader(isolatedIAMCtx, client.GetClientId(), redirectURIImplicit)
					require.NoError(t, err)
					return authRequestID
				}(),
			},
			want: want{
				addedClaims: &CustomAccessTokenClaims{
					Added1: "value1",
					Added2: "value2",
					Added3: "value3",
					Log:    []string{"addedLog1", "addedLog2", "addedLog3"},
				},
				setUserMetadata: []*metadata.Metadata{
					{Key: "key1", Value: []byte("value1")},
					{Key: "key2", Value: []byte("value2")},
					{Key: "key3", Value: []byte("value3")},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, closeF := tt.dep(isolatedIAMCtx, t, tt.req)
			defer closeF()

			got, err := instance.Client.OIDCv2.CreateCallback(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			callbackUrl, err := url.Parse(strings.Replace(got.GetCallbackUrl(), "#", "?", 1))
			require.NoError(t, err)
			claims := getAccessTokenClaims(tt.ctx, t, instance, callbackUrl)

			if tt.want.addedClaims != nil {
				assert.Equal(t, tt.want.addedClaims.Added1, claims.Added1)
				assert.Equal(t, tt.want.addedClaims.Added2, claims.Added2)
				assert.Equal(t, tt.want.addedClaims.Added3, claims.Added3)
				assert.Equal(t, tt.want.addedClaims.Log, claims.Log)
			}
			if len(tt.want.setUserMetadata) > 0 {
				checkForSetMetadata(isolatedIAMCtx, t, instance, userID, tt.want.setUserMetadata)
			}

		})
	}
}

func expectPreAccessTokenExecution(ctx context.Context, t *testing.T, instance *integration.Instance, clientID string, req *oidc_pb.CreateCallbackRequest, response *oidc_api.ContextInfoResponse) (string, func()) {
	userEmail := integration.Email()
	userPhone := integration.Phone()
	userResp := instance.CreateHumanUserVerified(ctx, instance.DefaultOrg.Id, userEmail, userPhone)

	sessionResp := createSession(ctx, t, instance, userResp.GetUserId())
	req.CallbackKind = &oidc_pb.CreateCallbackRequest_Session{
		Session: &oidc_pb.Session{
			SessionId:    sessionResp.GetSessionId(),
			SessionToken: sessionResp.GetSessionToken(),
		},
	}
	expectedContextInfo := contextInfoForUserOIDC(instance, "function/preaccesstoken", clientID, userResp, userEmail, userPhone)

	targetURL, closeF, _, _ := integration.TestServerCall(expectedContextInfo, 0, http.StatusOK, response)

	targetResp := waitForTarget(ctx, t, instance, targetURL, target_domain.TargetTypeCall, true)
	waitForExecutionOnCondition(ctx, t, instance, conditionFunction("preaccesstoken"), []string{targetResp.GetId()})
	return userResp.GetUserId(), closeF
}

func TestServer_ExecutionTargetPreSAMLResponse(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	ctxLoginClient := instance.WithAuthorizationToken(CTX, integration.UserTypeLogin)

	idpMetadata, err := instance.GetSAMLIDPMetadata()
	require.NoError(t, err)

	acsPost := idpMetadata.IDPSSODescriptors[0].SingleSignOnServices[1]
	_, _, spMiddlewarePost := createSAMLApplication(isolatedIAMCtx, t, instance, idpMetadata, saml.HTTPPostBinding, false, false)

	type want struct {
		addedAttributes map[string][]saml.AttributeValue
		setUserMetadata []*metadata.Metadata
	}
	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, t *testing.T, req *saml_pb.CreateResponseRequest) (string, func())
		req     *saml_pb.CreateResponseRequest
		want    want
		wantErr bool
	}{
		{
			name: "append attribute",
			ctx:  ctxLoginClient,
			dep: func(ctx context.Context, t *testing.T, req *saml_pb.CreateResponseRequest) (string, func()) {
				response := &saml_api.ContextInfoResponse{
					AppendAttribute: []*saml_api.AppendAttribute{
						{Name: "added", NameFormat: "format", Value: []string{"value"}},
					},
				}
				return expectPreSAMLResponseExecution(ctx, t, instance, req, response)
			},
			req: &saml_pb.CreateResponseRequest{
				SamlRequestId: func() string {
					_, samlRequestID, err := instance.CreateSAMLAuthRequest(spMiddlewarePost, instance.Users[integration.UserTypeOrgOwner].ID, acsPost, integration.RelayState(), saml.HTTPPostBinding)
					require.NoError(t, err)
					return samlRequestID
				}(),
			},
			want: want{
				addedAttributes: map[string][]saml.AttributeValue{
					"added": {saml.AttributeValue{Value: "value"}},
				},
			},
			wantErr: false,
		},
		{
			name: "set user metadata",
			ctx:  ctxLoginClient,
			dep: func(ctx context.Context, t *testing.T, req *saml_pb.CreateResponseRequest) (string, func()) {
				response := &saml_api.ContextInfoResponse{
					SetUserMetadata: []*domain.Metadata{
						{Key: "key", Value: []byte("value")},
					},
				}
				return expectPreSAMLResponseExecution(ctx, t, instance, req, response)
			},
			req: &saml_pb.CreateResponseRequest{
				SamlRequestId: func() string {
					_, samlRequestID, err := instance.CreateSAMLAuthRequest(spMiddlewarePost, instance.Users[integration.UserTypeOrgOwner].ID, acsPost, integration.RelayState(), saml.HTTPPostBinding)
					require.NoError(t, err)
					return samlRequestID
				}(),
			},
			want: want{
				setUserMetadata: []*metadata.Metadata{
					{Key: "key", Value: []byte("value")},
				},
			},
			wantErr: false,
		},
		{
			name: "set user metadata",
			ctx:  ctxLoginClient,
			dep: func(ctx context.Context, t *testing.T, req *saml_pb.CreateResponseRequest) (string, func()) {
				response := &saml_api.ContextInfoResponse{
					AppendAttribute: []*saml_api.AppendAttribute{
						{Name: "added1", NameFormat: "format", Value: []string{"value1"}},
						{Name: "added2", NameFormat: "format", Value: []string{"value2"}},
						{Name: "added3", NameFormat: "format", Value: []string{"value3"}},
					},
					SetUserMetadata: []*domain.Metadata{
						{Key: "key1", Value: []byte("value1")},
						{Key: "key2", Value: []byte("value2")},
						{Key: "key3", Value: []byte("value3")},
					},
				}
				return expectPreSAMLResponseExecution(ctx, t, instance, req, response)
			},
			req: &saml_pb.CreateResponseRequest{
				SamlRequestId: func() string {
					_, samlRequestID, err := instance.CreateSAMLAuthRequest(spMiddlewarePost, instance.Users[integration.UserTypeOrgOwner].ID, acsPost, integration.RelayState(), saml.HTTPPostBinding)
					require.NoError(t, err)
					return samlRequestID
				}(),
			},
			want: want{
				addedAttributes: map[string][]saml.AttributeValue{
					"added1": {saml.AttributeValue{Value: "value1"}},
					"added2": {saml.AttributeValue{Value: "value2"}},
					"added3": {saml.AttributeValue{Value: "value3"}},
				},
				setUserMetadata: []*metadata.Metadata{
					{Key: "key1", Value: []byte("value1")},
					{Key: "key2", Value: []byte("value2")},
					{Key: "key3", Value: []byte("value3")},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, closeF := tt.dep(isolatedIAMCtx, t, tt.req)
			defer closeF()

			got, err := instance.Client.SAMLv2.CreateResponse(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			attributes := getSAMLResponseAttributes(t, got.GetPost().GetSamlResponse(), spMiddlewarePost)
			for k, v := range tt.want.addedAttributes {
				found := false
				for _, attribute := range attributes {
					if attribute.Name == k {
						found = true
						assert.Equal(t, v, attribute.Values)
					}
				}
				if !assert.True(t, found) {
					return
				}
			}
			if len(tt.want.setUserMetadata) > 0 {
				checkForSetMetadata(isolatedIAMCtx, t, instance, userID, tt.want.setUserMetadata)
			}
		})
	}
}

func expectPreSAMLResponseExecution(ctx context.Context, t *testing.T, instance *integration.Instance, req *saml_pb.CreateResponseRequest, response *saml_api.ContextInfoResponse) (string, func()) {
	userEmail := integration.Email()
	userPhone := integration.Phone()
	userResp := instance.CreateHumanUserVerified(ctx, instance.DefaultOrg.Id, userEmail, userPhone)

	sessionResp := createSession(ctx, t, instance, userResp.GetUserId())
	req.ResponseKind = &saml_pb.CreateResponseRequest_Session{
		Session: &saml_pb.Session{
			SessionId:    sessionResp.GetSessionId(),
			SessionToken: sessionResp.GetSessionToken(),
		},
	}
	expectedContextInfo := contextInfoForUserSAML(instance, "function/presamlresponse", userResp, userEmail, userPhone)

	targetURL, closeF, _, _ := integration.TestServerCall(expectedContextInfo, 0, http.StatusOK, response)

	targetResp := waitForTarget(ctx, t, instance, targetURL, target_domain.TargetTypeCall, true)
	waitForExecutionOnCondition(ctx, t, instance, conditionFunction("presamlresponse"), []string{targetResp.GetId()})

	return userResp.GetUserId(), closeF
}

func createSAMLSP(t *testing.T, idpMetadata *saml.EntityDescriptor, binding string) (string, *samlsp.Middleware) {
	rootURL := "example." + integration.DomainName()
	spMiddleware, err := integration.CreateSAMLSP("https://"+rootURL, idpMetadata, binding)
	require.NoError(t, err)
	return rootURL, spMiddleware
}

func createSAMLApplication(ctx context.Context, t *testing.T, instance *integration.Instance, idpMetadata *saml.EntityDescriptor, binding string, projectRoleCheck, hasProjectCheck bool) (string, string, *samlsp.Middleware) {
	project := instance.CreateProject(ctx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), projectRoleCheck, hasProjectCheck)
	rootURL, sp := createSAMLSP(t, idpMetadata, binding)
	_, err := instance.CreateSAMLClient(ctx, project.GetId(), sp)
	require.NoError(t, err)
	return project.GetId(), rootURL, sp
}

func getSAMLResponseAttributes(t *testing.T, samlResponse string, sp *samlsp.Middleware) []saml.Attribute {
	data, err := base64.StdEncoding.DecodeString(samlResponse)
	require.NoError(t, err)
	sp.ServiceProvider.AllowIDPInitiated = true
	assertion, err := sp.ServiceProvider.ParseXMLResponse(data, []string{})
	require.NoError(t, err)
	return assertion.AttributeStatements[0].Attributes
}

func contextInfoForUserSAML(instance *integration.Instance, function string, userResp *user.AddHumanUserResponse, email, phone string) *saml_api.ContextInfo {
	return &saml_api.ContextInfo{
		Function: function,
		User: &query.User{
			ID:                 userResp.GetUserId(),
			CreationDate:       userResp.Details.ChangeDate.AsTime(),
			ChangeDate:         userResp.Details.ChangeDate.AsTime(),
			ResourceOwner:      instance.DefaultOrg.GetId(),
			Sequence:           userResp.Details.Sequence,
			State:              1,
			Type:               domain.UserTypeHuman,
			Username:           email,
			PreferredLoginName: email,
			LoginNames:         []string{email},
			Human: &query.Human{
				FirstName:              "Mickey",
				LastName:               "Mouse",
				NickName:               "Mickey",
				DisplayName:            "Mickey Mouse",
				AvatarKey:              "",
				PreferredLanguage:      language.Dutch,
				Gender:                 2,
				Email:                  domain.EmailAddress(email),
				IsEmailVerified:        true,
				Phone:                  domain.PhoneNumber(phone),
				IsPhoneVerified:        true,
				PasswordChangeRequired: false,
				PasswordChanged:        time.Time{},
				MFAInitSkipped:         time.Time{},
			},
		},
		UserGrants: nil,
		Response:   nil,
	}
}
