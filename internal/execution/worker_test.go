package execution

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/repository/action"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type fieldsWorker struct {
	now nowFunc
}
type argsWorker struct {
	job *river.Job[*exec_repo.Request]
}
type wantWorker struct {
	targets        []target_domain.Target
	sendStatusCode int
	err            assert.ErrorAssertionFunc
}

func newExecutionWorker(f fieldsWorker) *Worker {
	return &Worker{
		config: WorkerConfig{
			Workers:             1,
			TransactionDuration: 5 * time.Second,
			MaxTtl:              5 * time.Minute,
		},
		now: f.now,
	}
}

const (
	userID     = "user1"
	orgID      = "orgID"
	instanceID = "instanceID"
	eventID    = "eventID"
	eventData  = `{"name":"name","script":"name(){}","timeout":3000000000,"allowedToFail":true}`
)

func Test_handleEventExecution(t *testing.T) {
	testNow := time.Now
	tests := []struct {
		name string
		test func() (fieldsWorker, argsWorker, wantWorker)
	}{
		{
			"max TTL",
			func() (fieldsWorker, argsWorker, wantWorker) {
				return fieldsWorker{
						now: testNow,
					},
					argsWorker{
						job: &river.Job[*exec_repo.Request]{
							JobRow: &rivertype.JobRow{
								CreatedAt: time.Now().Add(-1 * time.Hour),
							},
							Args: &exec_repo.Request{
								Aggregate: &eventstore.Aggregate{
									InstanceID:    instanceID,
									ID:            eventID,
									ResourceOwner: instanceID,
								},
								Sequence:  1,
								CreatedAt: time.Now().Add(-1 * time.Hour),
								EventType: user.HumanInviteCodeAddedType,
								UserID:    userID,
								EventData: []byte(eventData),
							},
						},
					},
					wantWorker{
						targets:        mockTargets(1),
						sendStatusCode: http.StatusOK,
						err: func(tt assert.TestingT, err error, i ...interface{}) bool {
							return errors.Is(err, new(river.JobCancelError))
						},
					}
			},
		},
		{
			"none",
			func() (fieldsWorker, argsWorker, wantWorker) {
				return fieldsWorker{
						now: testNow,
					},
					argsWorker{
						job: &river.Job[*exec_repo.Request]{
							JobRow: &rivertype.JobRow{
								CreatedAt: time.Now(),
							},
							Args: &exec_repo.Request{
								Aggregate: &eventstore.Aggregate{
									InstanceID:    instanceID,
									ID:            eventID,
									ResourceOwner: instanceID,
								},
								Sequence:  1,
								CreatedAt: time.Now(),
								EventType: user.HumanInviteCodeAddedType,
								UserID:    userID,
								EventData: []byte(eventData),
							},
						},
					},
					wantWorker{
						targets:        mockTargets(0),
						sendStatusCode: http.StatusOK,
						err:            nil,
					}
			},
		},
		{
			"single",
			func() (fieldsWorker, argsWorker, wantWorker) {
				return fieldsWorker{
						now: testNow,
					},
					argsWorker{
						job: &river.Job[*exec_repo.Request]{
							JobRow: &rivertype.JobRow{
								CreatedAt: time.Now(),
							},
							Args: &exec_repo.Request{
								Aggregate: &eventstore.Aggregate{
									InstanceID:    instanceID,
									Type:          action.AggregateType,
									Version:       action.AggregateVersion,
									ID:            eventID,
									ResourceOwner: orgID,
								},
								Sequence:  1,
								CreatedAt: time.Now().UTC(),
								EventType: action.AddedEventType,
								UserID:    userID,
								EventData: []byte(eventData),
							},
						},
					},
					wantWorker{
						targets:        mockTargets(1),
						sendStatusCode: http.StatusOK,
						err:            nil,
					}
			},
		},
		{
			"single, failed 400",
			func() (fieldsWorker, argsWorker, wantWorker) {
				return fieldsWorker{
						now: testNow,
					},
					argsWorker{
						job: &river.Job[*exec_repo.Request]{
							JobRow: &rivertype.JobRow{
								CreatedAt: time.Now(),
							},
							Args: &exec_repo.Request{
								Aggregate: &eventstore.Aggregate{
									InstanceID:    instanceID,
									Type:          action.AggregateType,
									Version:       action.AggregateVersion,
									ID:            eventID,
									ResourceOwner: orgID,
								},
								Sequence:  1,
								CreatedAt: time.Now().UTC(),
								EventType: action.AddedEventType,
								UserID:    userID,
								EventData: []byte(eventData),
							},
						},
					},
					wantWorker{
						targets:        mockTargets(1),
						sendStatusCode: http.StatusBadRequest,
						err: func(tt assert.TestingT, err error, i ...interface{}) bool {
							return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "EXEC-dra6yamk98", "Errors.Execution.Failed"))
						},
					}
			},
		},
		{
			"multiple",
			func() (fieldsWorker, argsWorker, wantWorker) {
				return fieldsWorker{
						now: testNow,
					},
					argsWorker{
						job: &river.Job[*exec_repo.Request]{
							JobRow: &rivertype.JobRow{
								CreatedAt: time.Now(),
							},
							Args: &exec_repo.Request{
								Aggregate: &eventstore.Aggregate{
									InstanceID:    instanceID,
									Type:          action.AggregateType,
									Version:       action.AggregateVersion,
									ID:            eventID,
									ResourceOwner: orgID,
								},
								Sequence:  1,
								CreatedAt: time.Now().UTC(),
								EventType: action.AddedEventType,
								UserID:    userID,
								EventData: []byte(eventData),
							},
						},
					},
					wantWorker{
						targets:        mockTargets(3),
						sendStatusCode: http.StatusOK,
						err:            nil,
					}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, a, w := tt.test()

			closeFuncs := make([]func(), len(w.targets))
			calledFuncs := make([]func() bool, len(w.targets))
			for i := range w.targets {
				url, closeF, calledF := testServerCall(
					exec_repo.ContextInfoFromRequest(a.job.Args),
					time.Second,
					w.sendStatusCode,
					nil,
				)
				w.targets[i].Endpoint = url
				closeFuncs[i] = closeF
				calledFuncs[i] = calledF
			}

			data, err := json.Marshal(w.targets)
			require.NoError(t, err)
			a.job.Args.TargetsData = data

			err = newExecutionWorker(f).Work(
				authz.WithInstanceID(context.Background(), instanceID),
				a.job,
			)

			if w.err != nil {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			for _, closeF := range closeFuncs {
				closeF()
			}
			for _, calledF := range calledFuncs {
				assert.True(t, calledF())
			}
		})
	}
}

func mockTarget() target_domain.Target {
	return target_domain.Target{
		InstanceID:       "instanceID",
		ExecutionID:      "executionID",
		TargetID:         "targetID",
		TargetType:       target_domain.TargetTypeWebhook,
		Endpoint:         "endpoint",
		Timeout:          time.Minute,
		InterruptOnError: true,
		SigningKeyDec:    "key",
	}
}

func mockTargets(count int) []target_domain.Target {
	var targets []target_domain.Target
	if count > 0 {
		targets = make([]target_domain.Target, count)
		for i := range targets {
			targets[i] = mockTarget()
		}
	}
	return targets
}

func mockTargetsToBytes(targets []target_domain.Target) []byte {
	data, _ := json.Marshal(targets)
	return data
}
