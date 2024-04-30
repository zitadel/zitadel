//go:build integration

package action_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/grpc/server/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v3alpha"
)

func TestServer_ExecutionTarget(t *testing.T) {
	ensureFeatureEnabled(t)

	fullMethod := "/zitadel.action.v3alpha.ActionService/GetTargetByID"

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(context.Context, *action.GetTargetByIDRequest, *action.GetTargetByIDResponse) (func(), error)
		clean   func(context.Context)
		req     *action.GetTargetByIDRequest
		want    *action.GetTargetByIDResponse
		wantErr bool
	}{
		{
			name: "GetTargetByID, request and response, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *action.GetTargetByIDRequest, response *action.GetTargetByIDResponse) (func(), error) {

				instanceID := Tester.Instance.InstanceID()
				orgID := Tester.Organisation.ID
				projectID := ""
				userID := Tester.Users.Get(integration.FirstInstanceUsersKey, integration.IAMOwner).ID

				// create target for target changes
				targetCreatedName := fmt.Sprint("GetTargetByID", time.Now().UnixNano()+1)
				targetCreatedURL := "https://nonexistent"

				targetCreated := Tester.CreateTarget(ctx, t, targetCreatedName, targetCreatedURL, domain.TargetTypeCall, false)

				// request received by target
				wantRequest := &middleware.ContextInfoRequest{FullMethod: fullMethod, InstanceID: instanceID, OrgID: orgID, ProjectID: projectID, UserID: userID, Request: request}
				changedRequest := &action.GetTargetByIDRequest{TargetId: targetCreated.GetId()}
				// replace original request with different targetID
				urlRequest, closeRequest := testServerCall(wantRequest, 0, http.StatusOK, changedRequest)
				targetRequest := Tester.CreateTarget(ctx, t, "", urlRequest, domain.TargetTypeCall, false)
				Tester.SetExecution(ctx, t, conditionRequestFullMethod(fullMethod), executionTargetsSingleTarget(targetRequest.GetId()))
				// GetTargetByID with used target
				request.TargetId = targetRequest.GetId()

				// expected response from the GetTargetByID
				expectedResponse := &action.GetTargetByIDResponse{
					Target: &action.Target{
						TargetId: targetCreated.GetId(),
						Details:  targetCreated.GetDetails(),
						Name:     targetCreatedName,
						Endpoint: targetCreatedURL,
						TargetType: &action.Target_RestCall{
							RestCall: &action.SetRESTCall{
								InterruptOnError: false,
							},
						},
						Timeout: durationpb.New(10 * time.Second),
					},
				}
				// has to be set separately because of the pointers
				response.Target = &action.Target{
					TargetId: targetCreated.GetId(),
					Details:  targetCreated.GetDetails(),
					Name:     targetCreatedName,
					TargetType: &action.Target_RestCall{
						RestCall: &action.SetRESTCall{
							InterruptOnError: false,
						},
					},
					Timeout: durationpb.New(10 * time.Second),
				}

				// content for partial update
				changedResponse := &action.GetTargetByIDResponse{
					Target: &action.Target{
						TargetId: "changed",
					},
				}
				// change partial updated content on returned response
				response.Target.TargetId = changedResponse.Target.TargetId

				// response received by target
				wantResponse := &middleware.ContextInfoResponse{
					FullMethod: fullMethod,
					InstanceID: instanceID,
					OrgID:      orgID,
					ProjectID:  projectID,
					UserID:     userID,
					Request:    changedRequest,
					Response:   expectedResponse,
				}
				// after request with different targetID, return changed response
				targetResponseURL, closeResponse := testServerCall(wantResponse, 0, http.StatusOK, changedResponse)
				targetResponse := Tester.CreateTarget(ctx, t, "", targetResponseURL, domain.TargetTypeCall, false)
				Tester.SetExecution(ctx, t, conditionResponseFullMethod(fullMethod), executionTargetsSingleTarget(targetResponse.GetId()))

				return func() {
					closeRequest()
					closeResponse()
				}, nil
			},
			clean: func(ctx context.Context) {
				Tester.DeleteExecution(ctx, t, conditionRequestFullMethod(fullMethod))
				Tester.DeleteExecution(ctx, t, conditionResponseFullMethod(fullMethod))
			},
			req:  &action.GetTargetByIDRequest{},
			want: &action.GetTargetByIDResponse{},
		},
		/*{
			name: "GetTargetByID, request, interrupt",
			ctx:  CTX,
			dep: func(ctx context.Context, request *action.GetTargetByIDRequest, response *action.GetTargetByIDResponse) (func(), error) {

				fullMethod := "/zitadel.action.v3alpha.ActionService/GetTargetByID"
				instanceID := Tester.Instance.InstanceID()
				orgID := Tester.Organisation.ID
				projectID := ""
				userID := Tester.Users.Get(integration.FirstInstanceUsersKey, integration.IAMOwner).ID

				// request received by target
				wantRequest := &middleware.ContextInfoRequest{FullMethod: fullMethod, InstanceID: instanceID, OrgID: orgID, ProjectID: projectID, UserID: userID, Request: request}
				urlRequest, closeRequest := testServerCall(wantRequest, 0, http.StatusInternalServerError, &action.GetTargetByIDRequest{TargetId: "notchanged"})

				targetRequest := Tester.CreateTarget(ctx, t, "", urlRequest, domain.TargetTypeCall, true)
				Tester.SetExecution(ctx, t, conditionRequestFullMethod(fullMethod), executionTargetsSingleTarget(targetRequest.GetId()))
				// GetTargetByID with used target
				request.TargetId = targetRequest.GetId()

				return func() {
					closeRequest()
				}, nil
			},
			clean: func(ctx context.Context) {
				Tester.DeleteExecution(ctx, t, conditionRequestFullMethod(fullMethod))
			},
			req:     &action.GetTargetByIDRequest{},
			wantErr: true,
		},
		{
			name: "GetTargetByID, response, interrupt",
			ctx:  CTX,
			dep: func(ctx context.Context, request *action.GetTargetByIDRequest, response *action.GetTargetByIDResponse) (func(), error) {

				fullMethod := "/zitadel.action.v3alpha.ActionService/GetTargetByID"
				instanceID := Tester.Instance.InstanceID()
				orgID := Tester.Organisation.ID
				projectID := ""
				userID := Tester.Users.Get(integration.FirstInstanceUsersKey, integration.IAMOwner).ID

				// create target for target changes
				targetCreatedName := fmt.Sprint("GetTargetByID", time.Now().UnixNano()+1)
				targetCreatedURL := "https://nonexistent"

				targetCreated := Tester.CreateTarget(ctx, t, targetCreatedName, targetCreatedURL, domain.TargetTypeCall, false)

				// GetTargetByID with used target
				request.TargetId = targetCreated.GetId()

				// expected response from the GetTargetByID
				expectedResponse := &action.GetTargetByIDResponse{
					Target: &action.Target{
						TargetId: targetCreated.GetId(),
						Details:  targetCreated.GetDetails(),
						Name:     targetCreatedName,
						Endpoint: targetCreatedURL,
						TargetType: &action.Target_RestCall{
							RestCall: &action.SetRESTCall{
								InterruptOnError: false,
							},
						},
						Timeout: durationpb.New(10 * time.Second),
					},
				}

				// content for partial update
				changedResponse := &action.GetTargetByIDResponse{
					Target: &action.Target{
						TargetId: "changed",
					},
				}

				// response received by target
				wantResponse := &middleware.ContextInfoResponse{
					FullMethod: fullMethod,
					InstanceID: instanceID,
					OrgID:      orgID,
					ProjectID:  projectID,
					UserID:     userID,
					Request:    request,
					Response:   expectedResponse,
				}
				// after request with different targetID, return changed response
				targetResponseURL, closeResponse := testServerCall(wantResponse, 0, http.StatusInternalServerError, changedResponse)
				targetResponse := Tester.CreateTarget(ctx, t, "", targetResponseURL, domain.TargetTypeCall, true)
				Tester.SetExecution(ctx, t, conditionResponseFullMethod(fullMethod), executionTargetsSingleTarget(targetResponse.GetId()))

				return func() {
					closeResponse()
				}, nil
			},
			clean: func(ctx context.Context) {
				Tester.DeleteExecution(ctx, t, conditionResponseFullMethod(fullMethod))
			},
			req:     &action.GetTargetByIDRequest{},
			wantErr: true,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				close, err := tt.dep(tt.ctx, tt.req, tt.want)
				require.NoError(t, err)
				defer close()
			}

			got, err := Client.GetTargetByID(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want.GetTarget(), got.GetTarget())

			assert.Equal(t, tt.want.Target.TargetId, got.Target.TargetId)

			if tt.clean != nil {
				tt.clean(tt.ctx)
			}
		})
	}
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
