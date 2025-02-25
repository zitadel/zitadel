//go:build integration

package action_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
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
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/app"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/metadata"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
	action "github.com/zitadel/zitadel/pkg/grpc/resources/action/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
	saml_pb "github.com/zitadel/zitadel/pkg/grpc/saml/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

const (
	redirectURI         = "https://callback"
	logoutRedirectURI   = "https://logged-out"
	redirectURIImplicit = "http://localhost:9999/callback"
)

var (
	loginV2 = &app.LoginVersion{Version: &app.LoginVersion_LoginV2{LoginV2: &app.LoginV2{BaseUri: nil}}}
)

func TestServer_ExecutionTarget(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	fullMethod := "/zitadel.resources.action.v3alpha.ZITADELActions/GetTarget"

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(context.Context, *action.GetTargetRequest, *action.GetTargetResponse) (func(), error)
		clean   func(context.Context)
		req     *action.GetTargetRequest
		want    *action.GetTargetResponse
		wantErr bool
	}{
		{
			name: "GetTarget, request and response, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) (func(), error) {

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
				urlRequest, closeRequest := testServerCall(wantRequest, 0, http.StatusOK, changedRequest)

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
							Timeout: durationpb.New(10 * time.Second),
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
						Timeout:  durationpb.New(10 * time.Second),
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
				targetResponseURL, closeResponse := testServerCall(wantResponse, 0, http.StatusOK, changedResponse)

				targetResponse := waitForTarget(ctx, t, instance, targetResponseURL, domain.TargetTypeCall, false)
				waitForExecutionOnCondition(ctx, t, instance, conditionResponseFullMethod(fullMethod), executionTargetsSingleTarget(targetResponse.GetDetails().GetId()))
				return func() {
					closeRequest()
					closeResponse()
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
			dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) (func(), error) {

				fullMethod := "/zitadel.resources.action.v3alpha.ZITADELActions/GetTarget"
				orgID := instance.DefaultOrg.Id
				projectID := ""
				userID := instance.Users.Get(integration.UserTypeIAMOwner).ID

				// request received by target
				wantRequest := &middleware.ContextInfoRequest{FullMethod: fullMethod, InstanceID: instance.ID(), OrgID: orgID, ProjectID: projectID, UserID: userID, Request: request}
				urlRequest, closeRequest := testServerCall(wantRequest, 0, http.StatusInternalServerError, &action.GetTargetRequest{Id: "notchanged"})

				targetRequest := waitForTarget(ctx, t, instance, urlRequest, domain.TargetTypeCall, true)
				waitForExecutionOnCondition(ctx, t, instance, conditionRequestFullMethod(fullMethod), executionTargetsSingleTarget(targetRequest.GetDetails().GetId()))
				// GetTarget with used target
				request.Id = targetRequest.GetDetails().GetId()
				return func() {
					closeRequest()
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
			dep: func(ctx context.Context, request *action.GetTargetRequest, response *action.GetTargetResponse) (func(), error) {

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
							Timeout: durationpb.New(10 * time.Second),
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
				targetResponseURL, closeResponse := testServerCall(wantResponse, 0, http.StatusInternalServerError, changedResponse)

				targetResponse := waitForTarget(ctx, t, instance, targetResponseURL, domain.TargetTypeCall, true)
				waitForExecutionOnCondition(ctx, t, instance, conditionResponseFullMethod(fullMethod), executionTargetsSingleTarget(targetResponse.GetDetails().GetId()))
				return func() {
					closeResponse()
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
			if tt.dep != nil {
				close, err := tt.dep(tt.ctx, tt.req, tt.want)
				require.NoError(t, err)
				defer close()
			}
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(isolatedIAMOwnerCTX, time.Minute)
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

func testServerCall(
	reqBody interface{},
	sleep time.Duration,
	statusCode int,
	respBody interface{},
) (string, func()) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(reqBody)
		if err != nil {
			http.Error(w, "error, marshall: "+err.Error(), http.StatusInternalServerError)
			return
		}

		sentBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error, read body: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !reflect.DeepEqual(data, sentBody) {
			http.Error(w, "error, equal:\n"+string(data)+"\nsent:\n"+string(sentBody), http.StatusInternalServerError)
			return
		}
		if statusCode != http.StatusOK {
			http.Error(w, "error, statusCode", statusCode)
			return
		}

		time.Sleep(sleep)

		w.Header().Set("Content-Type", "application/json")
		resp, err := json.Marshal(respBody)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
		if _, err := io.WriteString(w, string(resp)); err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
	}

	server := httptest.NewServer(http.HandlerFunc(handler))

	return server.URL, server.Close
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
	ensureFeatureEnabled(t, instance)
	isolatedIAMCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	ctxLoginClient := instance.WithAuthorization(CTX, integration.UserTypeLogin)

	client, err := instance.CreateOIDCImplicitFlowClient(isolatedIAMCtx, redirectURIImplicit, loginV2)
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
				return expectPreUserinfoExecution(ctx, t, instance, req, response)
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
				return expectPreUserinfoExecution(ctx, t, instance, req, response)
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
				return expectPreUserinfoExecution(ctx, t, instance, req, response)
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
				return expectPreUserinfoExecution(ctx, t, instance, req, response)
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

func expectPreUserinfoExecution(ctx context.Context, t *testing.T, instance *integration.Instance, req *oidc_pb.CreateCallbackRequest, response *oidc_api.ContextInfoResponse) (string, func()) {
	userEmail := gofakeit.Email()
	userPhone := "+41" + gofakeit.Phone()
	userResp := instance.CreateHumanUserVerified(ctx, instance.DefaultOrg.Id, userEmail, userPhone)

	sessionResp := createSession(ctx, t, instance, userResp.GetUserId())
	req.CallbackKind = &oidc_pb.CreateCallbackRequest_Session{
		Session: &oidc_pb.Session{
			SessionId:    sessionResp.GetSessionId(),
			SessionToken: sessionResp.GetSessionToken(),
		},
	}
	expectedContextInfo := contextInfoForUserOIDC(instance, "function/preuserinfo", userResp, userEmail, userPhone)

	targetURL, closeF := testServerCall(expectedContextInfo, 0, http.StatusOK, response)

	targetResp := waitForTarget(ctx, t, instance, targetURL, domain.TargetTypeCall, true)
	waitForExecutionOnCondition(ctx, t, instance, conditionFunction("preuserinfo"), executionTargetsSingleTarget(targetResp.GetDetails().GetId()))
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

func contextInfoForUserOIDC(instance *integration.Instance, function string, userResp *user.AddHumanUserResponse, email, phone string) *oidc_api.ContextInfo {
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
	ensureFeatureEnabled(t, instance)
	isolatedIAMCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	ctxLoginClient := instance.WithAuthorization(CTX, integration.UserTypeLogin)

	client, err := instance.CreateOIDCImplicitFlowClient(isolatedIAMCtx, redirectURIImplicit, loginV2)
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
				return expectPreAccessTokenExecution(ctx, t, instance, req, response)
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
				return expectPreAccessTokenExecution(ctx, t, instance, req, response)
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
				return expectPreAccessTokenExecution(ctx, t, instance, req, response)
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
				return expectPreAccessTokenExecution(ctx, t, instance, req, response)
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

func expectPreAccessTokenExecution(ctx context.Context, t *testing.T, instance *integration.Instance, req *oidc_pb.CreateCallbackRequest, response *oidc_api.ContextInfoResponse) (string, func()) {
	userEmail := gofakeit.Email()
	userPhone := "+41" + gofakeit.Phone()
	userResp := instance.CreateHumanUserVerified(ctx, instance.DefaultOrg.Id, userEmail, userPhone)

	sessionResp := createSession(ctx, t, instance, userResp.GetUserId())
	req.CallbackKind = &oidc_pb.CreateCallbackRequest_Session{
		Session: &oidc_pb.Session{
			SessionId:    sessionResp.GetSessionId(),
			SessionToken: sessionResp.GetSessionToken(),
		},
	}
	expectedContextInfo := contextInfoForUserOIDC(instance, "function/preaccesstoken", userResp, userEmail, userPhone)

	targetURL, closeF := testServerCall(expectedContextInfo, 0, http.StatusOK, response)

	targetResp := waitForTarget(ctx, t, instance, targetURL, domain.TargetTypeCall, true)
	waitForExecutionOnCondition(ctx, t, instance, conditionFunction("preaccesstoken"), executionTargetsSingleTarget(targetResp.GetDetails().GetId()))
	return userResp.GetUserId(), closeF
}

func TestServer_ExecutionTargetPreSAMLResponse(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	ctxLoginClient := instance.WithAuthorization(CTX, integration.UserTypeLogin)

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
					_, samlRequestID, err := instance.CreateSAMLAuthRequest(spMiddlewarePost, instance.Users[integration.UserTypeOrgOwner].ID, acsPost, gofakeit.BitcoinAddress(), saml.HTTPPostBinding)
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
					_, samlRequestID, err := instance.CreateSAMLAuthRequest(spMiddlewarePost, instance.Users[integration.UserTypeOrgOwner].ID, acsPost, gofakeit.BitcoinAddress(), saml.HTTPPostBinding)
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
					_, samlRequestID, err := instance.CreateSAMLAuthRequest(spMiddlewarePost, instance.Users[integration.UserTypeOrgOwner].ID, acsPost, gofakeit.BitcoinAddress(), saml.HTTPPostBinding)
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
	userEmail := gofakeit.Email()
	userPhone := "+41" + gofakeit.Phone()
	userResp := instance.CreateHumanUserVerified(ctx, instance.DefaultOrg.Id, userEmail, userPhone)

	sessionResp := createSession(ctx, t, instance, userResp.GetUserId())
	req.ResponseKind = &saml_pb.CreateResponseRequest_Session{
		Session: &saml_pb.Session{
			SessionId:    sessionResp.GetSessionId(),
			SessionToken: sessionResp.GetSessionToken(),
		},
	}
	expectedContextInfo := contextInfoForUserSAML(instance, "function/presamlresponse", userResp, userEmail, userPhone)

	targetURL, closeF := testServerCall(expectedContextInfo, 0, http.StatusOK, response)

	targetResp := waitForTarget(ctx, t, instance, targetURL, domain.TargetTypeCall, true)
	waitForExecutionOnCondition(ctx, t, instance, conditionFunction("presamlresponse"), executionTargetsSingleTarget(targetResp.GetDetails().GetId()))

	return userResp.GetUserId(), closeF
}

func createSAMLSP(t *testing.T, idpMetadata *saml.EntityDescriptor, binding string) (string, *samlsp.Middleware) {
	rootURL := "example." + gofakeit.DomainName()
	spMiddleware, err := integration.CreateSAMLSP("https://"+rootURL, idpMetadata, binding)
	require.NoError(t, err)
	return rootURL, spMiddleware
}

func createSAMLApplication(ctx context.Context, t *testing.T, instance *integration.Instance, idpMetadata *saml.EntityDescriptor, binding string, projectRoleCheck, hasProjectCheck bool) (string, string, *samlsp.Middleware) {
	project, err := instance.CreateProjectWithPermissionCheck(ctx, projectRoleCheck, hasProjectCheck)
	require.NoError(t, err)
	rootURL, sp := createSAMLSP(t, idpMetadata, binding)
	_, err = instance.CreateSAMLClient(ctx, project.GetId(), sp)
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
