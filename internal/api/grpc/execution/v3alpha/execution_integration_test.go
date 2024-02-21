//go:build integration

package execution_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
)

func TestServer_SetExecution(t *testing.T) {
	targetResp, err := Tester.Client.ExecutionV3.CreateTarget(CTX, &execution.CreateTargetRequest{
		Name: fmt.Sprint(time.Now().UnixNano() + 1),
		TargetType: &execution.CreateTargetRequest_RestWebhook{
			RestWebhook: &execution.SetRESTWebhook{
				Url: "https://example.com",
			},
		},
		Timeout: durationpb.New(10 * time.Second),
		ExecutionType: &execution.CreateTargetRequest_InterruptOnError{
			InterruptOnError: true,
		},
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *execution.SetExecutionRequest
		want    *execution.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &execution.SetExecutionRequest{
				Condition: &execution.SetConditions{
					ConditionType: &execution.SetConditions_Request{
						Request: &execution.SetRequestExecution{
							Condition: &execution.SetRequestExecution_All{All: true},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "request no condition, error",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.SetConditions{
					ConditionType: &execution.SetConditions_Request{
						Request: &execution.SetRequestExecution{},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			wantErr: true,
		},
		{
			name: "request method, not existing",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.SetConditions{
					ConditionType: &execution.SetConditions_Request{
						Request: &execution.SetRequestExecution{
							Condition: &execution.SetRequestExecution_Method{
								Method: "/zitadel.session.v2beta.NotExistingService/List",
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			wantErr: true,
		},
		{
			name: "request method, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.SetConditions{
					ConditionType: &execution.SetConditions_Request{
						Request: &execution.SetRequestExecution{
							Condition: &execution.SetRequestExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/ListSessions",
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "request service, not existing",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.SetConditions{
					ConditionType: &execution.SetConditions_Request{
						Request: &execution.SetRequestExecution{
							Condition: &execution.SetRequestExecution_Service{
								Service: "NotExistingService",
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			wantErr: true,
		},
		{
			name: "request service, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.SetConditions{
					ConditionType: &execution.SetConditions_Request{
						Request: &execution.SetRequestExecution{
							Condition: &execution.SetRequestExecution_Service{
								Service: "zitadel.session.v2beta.SessionService",
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "request all, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.SetConditions{
					ConditionType: &execution.SetConditions_Request{
						Request: &execution.SetRequestExecution{
							Condition: &execution.SetRequestExecution_All{
								All: true,
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}
