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

	"github.com/go-jose/go-jose/v4"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/execution"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/zerrors"
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
					if interrupt && targetType != target_domain.TargetTypeAsync {
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

func TestExecutionWorker_AsyncTargetBoundByJobContext(t *testing.T) {
	const jobTimeout = 50 * time.Millisecond
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.URL.Path == "/slow" {
			select {
			case <-r.Context().Done():
			case <-time.After(time.Second):
			}
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	now := time.Now()
	worker := execution.NewWorker(execution.WorkerConfig{TransactionDuration: jobTimeout, MaxTtl: time.Minute}, nil, nil, func() time.Time { return now }, server.Client())
	newJob := func(endpoint string) *river.Job[*exec_repo.Request] {
		targets, err := json.Marshal([]target_domain.Target{
			{TargetID: "async", TargetType: target_domain.TargetTypeAsync, Endpoint: server.URL + endpoint, Timeout: 10 * time.Second},
		})
		require.NoError(t, err)
		return &river.Job[*exec_repo.Request]{JobRow: &rivertype.JobRow{CreatedAt: now}, Args: &exec_repo.Request{
			Aggregate: &eventstore.Aggregate{InstanceID: "instance"}, TargetsData: targets,
		}}
	}

	t.Run("job timeout stops a slow async target", func(t *testing.T) {
		job := newJob("/slow")
		ctx, cancel := context.WithTimeout(t.Context(), worker.Timeout(job))
		defer cancel()
		started := time.Now()
		err := worker.Work(ctx, job)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.False(t, errors.Is(err, new(river.JobCancelError)), "a timed out delivery must stay retryable")
		require.Less(t, time.Since(started), 500*time.Millisecond)
	})

	t.Run("cancelled worker does not deliver", func(t *testing.T) {
		calls.Store(0)
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		err := worker.Work(ctx, newJob("/fast"))
		require.ErrorIs(t, err, context.Canceled)
		require.Zero(t, calls.Load())
	})
}

func TestExecutionWorker_PermanentErrorCancelsJob(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()
	now := time.Now()
	worker := execution.NewWorker(execution.WorkerConfig{MaxTtl: time.Minute}, nil, mockGetActiveSigningWebKey, func() time.Time { return now }, server.Client())
	permanent := []target_domain.Target{
		{TargetID: "unknownType", TargetType: target_domain.TargetType(99), Endpoint: server.URL, Timeout: time.Second},
		{TargetID: "missingKey", TargetType: target_domain.TargetTypeWebhook, PayloadType: target_domain.PayloadTypeJWE, Endpoint: server.URL, Timeout: time.Second},
	}
	newJob := func(targets []target_domain.Target) *river.Job[*exec_repo.Request] {
		data, err := json.Marshal(targets)
		require.NoError(t, err)
		return &river.Job[*exec_repo.Request]{JobRow: &rivertype.JobRow{CreatedAt: now}, Args: &exec_repo.Request{
			Aggregate: &eventstore.Aggregate{InstanceID: "instance"}, TargetsData: data,
		}}
	}

	err := worker.Work(t.Context(), newJob(permanent))
	require.ErrorIs(t, err, new(river.JobCancelError))
	require.ErrorContains(t, err, "target unknownType:")
	require.ErrorContains(t, err, "target missingKey:")
	require.Zero(t, calls.Load())

	err = worker.Work(t.Context(), newJob(append(permanent, target_domain.Target{TargetID: "unavailable", TargetType: target_domain.TargetTypeWebhook, Endpoint: server.URL, Timeout: time.Second})))
	require.ErrorContains(t, err, "target unavailable:")
	require.False(t, errors.Is(err, new(river.JobCancelError)), "a retryable delivery error must keep the job retryable")
	require.EqualValues(t, 1, calls.Load())

	keyQueryFailure := func(context.Context) (*jose.JSONWebKey, error) {
		return nil, zerrors.ThrowInternal(errors.New("connection reset"), "QUERY-Shoo0", "Errors.Internal")
	}
	worker = execution.NewWorker(execution.WorkerConfig{MaxTtl: time.Minute}, nil, keyQueryFailure, func() time.Time { return now }, server.Client())
	err = worker.Work(t.Context(), newJob([]target_domain.Target{{TargetID: "jwt", TargetType: target_domain.TargetTypeWebhook, PayloadType: target_domain.PayloadTypeJWT, Endpoint: server.URL, Timeout: time.Second}}))
	require.ErrorContains(t, err, "target jwt:")
	require.False(t, errors.Is(err, new(river.JobCancelError)), "a signing key query failure can be transient and must stay retryable")
	require.EqualValues(t, 1, calls.Load())
}

func TestExecutionWorker_AsyncTargetsRunConcurrently(t *testing.T) {
	const delay = 200 * time.Millisecond
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		select {
		case <-r.Context().Done():
		case <-time.After(delay):
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	now := time.Now()
	worker := execution.NewWorker(execution.WorkerConfig{TransactionDuration: time.Minute, MaxTtl: time.Minute}, nil, nil, func() time.Time { return now }, server.Client())
	targets, err := json.Marshal([]target_domain.Target{
		{TargetID: "async1", TargetType: target_domain.TargetTypeAsync, Endpoint: server.URL, Timeout: 10 * time.Second},
		{TargetID: "async2", TargetType: target_domain.TargetTypeAsync, Endpoint: server.URL, Timeout: 10 * time.Second},
	})
	require.NoError(t, err)
	job := &river.Job[*exec_repo.Request]{JobRow: &rivertype.JobRow{CreatedAt: now}, Args: &exec_repo.Request{
		Aggregate: &eventstore.Aggregate{InstanceID: "instance"}, TargetsData: targets,
	}}

	started := time.Now()
	require.NoError(t, worker.Work(t.Context(), job))
	require.EqualValues(t, 2, calls.Load())
	require.Less(t, time.Since(started), 2*delay, "async targets of one job must not run one after another")

	ctx, cancel := context.WithTimeout(t.Context(), delay/4)
	defer cancel()
	started = time.Now()
	err = worker.Work(ctx, job)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.ErrorContains(t, err, "target async1:")
	require.ErrorContains(t, err, "target async2:")
	require.Less(t, time.Since(started), delay, "cancellation must stop concurrent async targets")
}
