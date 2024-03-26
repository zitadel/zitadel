//go:build integration

package execution_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
)

func TestServer_ExecutionTarget_Request(t *testing.T) {
	targetResp := Tester.CreateTarget(CTX, t)

	tests := []struct {
		name    string
		dep     func(context.Context, *execution.SetExecutionRequest) (int, error)
		req     *execution.SetExecutionRequest
		want    interface{}
		wantErr bool
	}{
		{
			name: "all, ok",
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_All{
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
