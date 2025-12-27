package execution_test

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
	"github.com/zitadel/zitadel/internal/execution"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/repository/action"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type fieldsWorker struct {
	now execution.NowFunc
}
type argsWorker struct {
	job *river.Job[*exec_repo.Request]
}
type wantWorker struct {
	targets        []target
	sendStatusCode int
	err            assert.ErrorAssertionFunc
}

type target target_domain.Target

func (t *target) validate(expectedBody []byte) func(*testing.T, []byte) bool {
	switch t.PayloadType {
	case target_domain.PayloadTypeUnspecified,
		target_domain.PayloadTypeJSON:
		return validateJSONPayload(expectedBody)
	case target_domain.PayloadTypeJWT:
		return validateJWTPayload(expectedBody)
	case target_domain.PayloadTypeJWE:
		return validateJWEPayload(expectedBody)
	default:
		return validateJSONPayload(expectedBody)
	}
}

func newExecutionWorker(f fieldsWorker) *execution.Worker {
	return execution.NewWorker(
		execution.WorkerConfig{
			Workers:             1,
			TransactionDuration: 5 * time.Second,
			MaxTtl:              5 * time.Minute,
		},
		nil,
		mockGetActiveSigningWebKey,
		f.now,
	)
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
						targets:        mockTargets(target_domain.PayloadTypeJSON),
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
						targets:        mockTargets(target_domain.PayloadTypeJSON),
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
						targets:        mockTargets(target_domain.PayloadTypeJSON),
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
						targets:        mockTargets(target_domain.PayloadTypeJSON, target_domain.PayloadTypeJWT, target_domain.PayloadTypeJWE),
						sendStatusCode: http.StatusOK,
						err:            nil,
					}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, a, w := tt.test()

			body, err := json.Marshal(exec_repo.ContextInfoFromRequest(a.job.Args))
			require.NoError(t, err)

			closeFuncs := make([]func(), len(w.targets))
			calledFuncs := make([]func() bool, len(w.targets))
			for i := range w.targets {
				url, closeF, calledF := listen(t, &callTestServer{
					method:      http.MethodPost,
					expectBody:  w.targets[i].validate(body),
					timeout:     time.Second,
					statusCode:  w.sendStatusCode,
					respondBody: nil,
				})
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

func mockTarget(payloadType target_domain.PayloadType) target {
	return target{
		ExecutionID:      "executionID",
		TargetID:         "targetID",
		TargetType:       target_domain.TargetTypeWebhook,
		Endpoint:         "endpoint",
		Timeout:          time.Minute,
		InterruptOnError: true,
		PayloadType:      payloadType,
		EncryptionKey:    encryptionKey,
		EncryptionKeyID:  encryptionKeyID,
	}
}

func mockTargets(payloadTypes ...target_domain.PayloadType) []target {
	targets := make([]target, len(payloadTypes))
	for i, payloadType := range payloadTypes {
		targets[i] = mockTarget(payloadType)
	}
	return targets
}
