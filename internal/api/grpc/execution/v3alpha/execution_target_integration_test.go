//go:build integration

package execution_test

/*
import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
)

func TestServer_ExecutionTarget_Request(t *testing.T) {

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(context.Context, *execution.GetTargetByIDRequest, *execution.GetTargetByIDResponse) (func(), error)
		req     *execution.GetTargetByIDRequest
		want    *execution.GetTargetByIDResponse
		wantErr bool
	}{
		{
			name: "GetTargetByID, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *execution.GetTargetByIDRequest, response *execution.GetTargetByIDResponse) (func(), error) {

				// create target which can be used for query
				fullMethod := "/zitadel.execution.v3alpha.ExecutionService/GetTargetByID"
				notExistingTarget := Tester.CreateTargetURL(ctx, t, "https://nonexistent")

				// condition on GetTargetByID
				cond := &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_Method{
								Method: fullMethod,
							},
						},
					},
				}

				// start target with response
				url, close := testServerCall(0, http.StatusOK, &execution.GetTargetByIDRequest{TargetId: notExistingTarget.GetId()})
				targetResp := Tester.CreateTargetURL(ctx, t, url)

				// GetTargetByID with used target
				request.TargetId = notExistingTarget.GetId()
				Tester.SetExecution(ctx, t, cond, []string{targetResp.GetId()}, []string{})
				response.Target = &execution.Target{
					TargetId: notExistingTarget.GetId(),
					Details:  notExistingTarget.GetDetails(),
				}

				return close, nil
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

func testServerCall(
	sleep time.Duration,
	statusCode int,
	respBody interface{},
) (string, func()) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if statusCode != http.StatusOK {
			http.Error(w, "error", statusCode)
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
*/
