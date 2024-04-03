//go:build integration

package execution_test

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
	"github.com/zitadel/zitadel/internal/integration"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
)

func TestServer_ExecutionTarget(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(context.Context, *execution.GetTargetByIDRequest, *execution.GetTargetByIDResponse) (func(), error)
		req     *execution.GetTargetByIDRequest
		want    *execution.GetTargetByIDResponse
		wantErr bool
	}{
		{
			name: "GetTargetByID, request and response, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *execution.GetTargetByIDRequest, response *execution.GetTargetByIDResponse) (func(), error) {

				fullMethod := "/zitadel.execution.v3alpha.ExecutionService/GetTargetByID"
				instanceID := Tester.Instance.InstanceID()
				orgID := Tester.Organisation.ID
				projectID := ""
				userID := Tester.Users.Get(integration.FirstInstanceUsersKey, integration.IAMOwner).ID

				// create target for target changes
				targetCreatedName := fmt.Sprint("GetTargetByID", time.Now().UnixNano()+1)
				targetCreatedURL := "https://nonexistent"
				targetCreated := Tester.CreateTarget(ctx, t, targetCreatedName, targetCreatedURL, true, false, true)

				// request received by target
				wantRequest := &middleware.ContextInfoRequest{FullMethod: fullMethod, InstanceID: instanceID, OrgID: orgID, ProjectID: projectID, UserID: userID, Request: request}
				changedRequest := &execution.GetTargetByIDRequest{TargetId: targetCreated.GetId()}
				// replace original request with different targetID
				urlRequest, closeRequest := testServerCall(wantRequest, 0, http.StatusOK, changedRequest)
				targetRequest := Tester.CreateTarget(ctx, t, "", urlRequest, true, false, true)
				Tester.SetExecution(ctx, t, conditionRequestFullMethod(fullMethod), []string{targetRequest.GetId()}, []string{})
				// GetTargetByID with used target
				request.TargetId = targetRequest.GetId()

				// expected response from the GetTargetByID
				expectedResponse := &execution.GetTargetByIDResponse{
					Target: &execution.Target{
						TargetId: targetCreated.GetId(),
						Details:  targetCreated.GetDetails(),
						Name:     targetCreatedName,
						TargetType: &execution.Target_RestRequestResponse{
							RestRequestResponse: &execution.SetRESTRequestResponse{
								Url: targetCreatedURL,
							},
						},
						Timeout:       durationpb.New(10 * time.Second),
						ExecutionType: &execution.Target_InterruptOnError{InterruptOnError: true},
					},
				}
				// has to be set separately because of the pointers
				response.Target = &execution.Target{
					TargetId: targetCreated.GetId(),
					Details:  targetCreated.GetDetails(),
					Name:     targetCreatedName,
					TargetType: &execution.Target_RestRequestResponse{
						RestRequestResponse: &execution.SetRESTRequestResponse{
							Url: targetCreatedURL,
						},
					},
					Timeout:       durationpb.New(10 * time.Second),
					ExecutionType: &execution.Target_InterruptOnError{InterruptOnError: true},
				}

				// content for partial update
				changedResponse := &execution.GetTargetByIDResponse{
					Target: &execution.Target{
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
				targetResponse := Tester.CreateTarget(ctx, t, "", targetResponseURL, true, false, true)
				Tester.SetExecution(ctx, t, conditionResponseFullMethod(fullMethod), []string{targetResponse.GetId()}, []string{})

				return func() {
					closeRequest()
					closeResponse()
				}, nil
			},
			req:  &execution.GetTargetByIDRequest{},
			want: &execution.GetTargetByIDResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				close, err := tt.dep(tt.ctx, tt.req, tt.want)
				require.NoError(t, err)
				defer close()
			}

			got, err := Client.GetTargetByID(tt.ctx, tt.req)
			if err != nil {
				fmt.Println("error")
				fmt.Println(err.Error())
			}
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want.GetTarget(), got.GetTarget())

			assert.Equal(t, tt.want.Target.TargetId, got.Target.TargetId)
		})
	}
}

func conditionRequestFullMethod(fullMethod string) *execution.Condition {
	return &execution.Condition{
		ConditionType: &execution.Condition_Request{
			Request: &execution.RequestExecution{
				Condition: &execution.RequestExecution_Method{
					Method: fullMethod,
				},
			},
		},
	}
}

func conditionResponseFullMethod(fullMethod string) *execution.Condition {
	return &execution.Condition{
		ConditionType: &execution.Condition_Response{
			Response: &execution.ResponseExecution{
				Condition: &execution.ResponseExecution_Method{
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
