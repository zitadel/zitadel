package execution_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/execution"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
)

func TestExecutionWorker_RetryDelivery(t *testing.T) {
	for _, targetType := range []target_domain.TargetType{target_domain.TargetTypeWebhook, target_domain.TargetTypeCall, target_domain.TargetTypeAsync} {
		for _, interrupt := range []bool{false, true} {
			for _, status := range []int{http.StatusTooManyRequests, http.StatusServiceUnavailable} {
				t.Run(fmt.Sprintf("type=%d/status=%d/interrupt=%t", targetType, status, interrupt), func(t *testing.T) {
					var failedCalls, successfulCalls atomic.Int32
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.Path == "/failing" && failedCalls.Add(1) == 1 {
							w.WriteHeader(status)
							return
						}
						if r.URL.Path == "/successful" {
							successfulCalls.Add(1)
						}
						w.WriteHeader(http.StatusNoContent)
					}))
					defer server.Close()
					now := time.Now()
					worker := execution.NewWorker(execution.WorkerConfig{MaxTtl: 5 * time.Minute}, nil, nil, func() time.Time { return now }, server.Client())
					targets, err := json.Marshal([]target_domain.Target{
						{TargetID: "failing", TargetType: targetType, Endpoint: server.URL + "/failing", Timeout: time.Second, InterruptOnError: interrupt},
						{TargetID: "successful", Endpoint: server.URL + "/successful", Timeout: time.Second},
					})
					require.NoError(t, err)
					job := &river.Job[*exec_repo.Request]{
						JobRow: &rivertype.JobRow{CreatedAt: now},
						Args:   &exec_repo.Request{Aggregate: &eventstore.Aggregate{InstanceID: "instance"}, TargetsData: targets},
					}
					err = worker.Work(t.Context(), job)
					require.ErrorContains(t, err, "target failing:")
					require.False(t, errors.Is(err, new(river.JobCancelError)), "delivery failures must reach River's retry policy")
					require.EqualValues(t, 1, failedCalls.Load(), "async delivery must finish before the job returns")
					wantSuccessful := int32(1)
					if interrupt {
						wantSuccessful = 0
					}
					require.Equal(t, wantSuccessful, successfulCalls.Load())
					require.NoError(t, worker.Work(t.Context(), job))
					require.EqualValues(t, 2, failedCalls.Load())
					require.Equal(t, wantSuccessful+1, successfulCalls.Load(), "retries can repeat successful deliveries")
					now = now.Add(5*time.Minute + time.Nanosecond)
					require.ErrorIs(t, worker.Work(t.Context(), job), new(river.JobCancelError))
					require.EqualValues(t, 2, failedCalls.Load(), "expired retries must not deliver")
				})
			}
		}
	}
}

func TestExecutionWorker_RetryTransportFailure(t *testing.T) {
	now := time.Now()
	worker := execution.NewWorker(execution.WorkerConfig{MaxTtl: time.Minute}, nil, nil, func() time.Time { return now }, http.DefaultClient)
	targets, err := json.Marshal([]target_domain.Target{{TargetID: "cancelled", Endpoint: "http://localhost", Timeout: time.Second}})
	require.NoError(t, err)
	job := &river.Job[*exec_repo.Request]{JobRow: &rivertype.JobRow{CreatedAt: now}, Args: &exec_repo.Request{
		Aggregate: &eventstore.Aggregate{InstanceID: "instance"}, TargetsData: targets,
	}}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	err = worker.Work(ctx, job)
	require.ErrorIs(t, err, context.Canceled)
	require.False(t, errors.Is(err, new(river.JobCancelError)))
	job.Args.TargetsData = []byte("invalid")
	require.ErrorIs(t, worker.Work(t.Context(), job), new(river.JobCancelError))
}
