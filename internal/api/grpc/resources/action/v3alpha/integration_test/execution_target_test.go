//go:build integration

package action_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/v2/internal/api/grpc/server/middleware"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/integration"
	object "github.com/zitadel/zitadel/v2/pkg/grpc/object/v3alpha"
	action "github.com/zitadel/zitadel/v2/pkg/grpc/resources/action/v3alpha"
	resource_object "github.com/zitadel/zitadel/v2/pkg/grpc/resources/object/v3alpha"
)

func TestServer_ExecutionTarget(t *testing.T) {
	t.Parallel()
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
				targetRequest := instance.CreateTarget(ctx, t, "", urlRequest, domain.TargetTypeCall, false)
				instance.SetExecution(ctx, t, conditionRequestFullMethod(fullMethod), executionTargetsSingleTarget(targetRequest.GetDetails().GetId()))

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
				targetResponse := instance.CreateTarget(ctx, t, "", targetResponseURL, domain.TargetTypeCall, false)
				instance.SetExecution(ctx, t, conditionResponseFullMethod(fullMethod), executionTargetsSingleTarget(targetResponse.GetDetails().GetId()))

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

				targetRequest := instance.CreateTarget(ctx, t, "", urlRequest, domain.TargetTypeCall, true)
				instance.SetExecution(ctx, t, conditionRequestFullMethod(fullMethod), executionTargetsSingleTarget(targetRequest.GetDetails().GetId()))
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
				targetResponse := instance.CreateTarget(ctx, t, "", targetResponseURL, domain.TargetTypeCall, true)
				instance.SetExecution(ctx, t, conditionResponseFullMethod(fullMethod), executionTargetsSingleTarget(targetResponse.GetDetails().GetId()))

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
			retryDuration := 5 * time.Second
			if ctxDeadline, ok := isolatedIAMOwnerCTX.Deadline(); ok {
				retryDuration = time.Until(ctxDeadline)
			}

			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := instance.Client.ActionV3Alpha.GetTarget(tt.ctx, tt.req)
				if tt.wantErr {
					assert.Error(ttt, err, "Error: "+err.Error())
				} else {
					assert.NoError(ttt, err)
				}
				if err != nil {
					return
				}

				integration.AssertResourceDetails(t, tt.want.GetTarget().GetDetails(), got.GetTarget().GetDetails())
				assert.Equal(t, tt.want.GetTarget().GetConfig(), got.GetTarget().GetConfig())
			}, retryDuration, time.Millisecond*100, "timeout waiting for expected execution result")

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
