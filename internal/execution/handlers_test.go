package execution

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/execution/mock"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/action"
	execution_rp "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_EventExecution(t *testing.T) {
	type args struct {
		event   eventstore.Event
		targets []*query.ExecutionTarget
	}
	type res struct {
		targets     []Target
		contextInfo *execution_rp.ContextInfoEvent
		wantErr     bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"session added, ok",
			args{
				event: &eventstore.BaseEvent{
					Agg: &eventstore.Aggregate{
						ID:            "aggID",
						Type:          session.AggregateType,
						ResourceOwner: "resourceOwner",
						InstanceID:    "instanceID",
						Version:       session.AggregateVersion,
					},
					EventType: session.AddedType,
					Seq:       1,
					Creation:  time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC),
					User:      userID,
					Data:      []byte(`{"ID":"","Seq":1,"Pos":0,"Creation":"2024-01-01T01:01:01.000000001Z"}`),
				},
				targets: []*query.ExecutionTarget{{
					InstanceID:       instanceID,
					ExecutionID:      "executionID",
					TargetID:         "targetID",
					TargetType:       domain.TargetTypeWebhook,
					Endpoint:         "endpoint",
					Timeout:          time.Minute,
					InterruptOnError: true,
					SigningKey:       "key",
				}},
			},
			res{
				targets: []Target{
					&query.ExecutionTarget{
						InstanceID:       instanceID,
						ExecutionID:      "executionID",
						TargetID:         "targetID",
						TargetType:       domain.TargetTypeWebhook,
						Endpoint:         "endpoint",
						Timeout:          time.Minute,
						InterruptOnError: true,
						SigningKey:       "key",
					},
				},
				contextInfo: &execution_rp.ContextInfoEvent{
					AggregateID:   "aggID",
					AggregateType: "session",
					ResourceOwner: "resourceOwner",
					InstanceID:    "instanceID",
					Version:       "v1",
					Sequence:      1,
					EventType:     "session.added",
					CreatedAt:     time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC).Format(time.RFC3339Nano),
					UserID:        userID,
					EventPayload:  []byte(`{"ID":"","Seq":1,"Pos":0,"Creation":"2024-01-01T01:01:01.000000001Z"}`),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := NewRequest(tt.args.event, tt.args.targets)
			if tt.res.wantErr {
				assert.Error(t, err)
				assert.Nil(t, request)
				return
			}
			assert.NoError(t, err)
			targets, err := TargetsFromRequest(request)
			assert.NoError(t, err)
			assert.Equal(t, tt.res.targets, targets)
			assert.Equal(t, tt.res.contextInfo, execution_rp.ContextInfoFromRequest(request))
		})
	}
}

func Test_groupsFromEventType(t *testing.T) {
	type args struct {
		eventType eventstore.EventType
	}
	type res struct {
		groups []string
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"user human mfa init skipped, ok",
			args{
				eventType: user.HumanMFAInitSkippedType,
			},
			res{
				groups: []string{
					"user.human.mfa.init.skipped",
					"user.human.mfa.init.*",
					"user.human.mfa.*",
					"user.human.*",
					"user.*",
				},
			},
		},
		{
			"session added, ok",
			args{
				eventType: session.AddedType,
			},
			res{
				groups: []string{
					"session.added",
					"session.*",
				},
			},
		},
		{
			"user added, ok",
			args{
				eventType: user.HumanAddedType,
			},
			res{
				groups: []string{
					"user.human.added",
					"user.human.*",
					"user.*",
				},
			},
		},
		{
			"execution set, ok",
			args{
				eventType: execution_rp.SetEventV2Type,
			},
			res{
				groups: []string{
					"execution.v2.set",
					"execution.v2.*",
					"execution.*",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.res.groups, groupsFromEventType(string(tt.args.eventType)))
		})
	}
}

func Test_idsForEventType(t *testing.T) {
	type args struct {
		eventType eventstore.EventType
	}
	type res struct {
		groups []string
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"session added, ok",
			args{
				eventType: session.AddedType,
			},
			res{
				groups: []string{
					"event/session.added",
					"event/session.*",
					"event",
				},
			},
		},
		{
			"user added, ok",
			args{
				eventType: user.HumanAddedType,
			},
			res{
				groups: []string{
					"event/user.human.added",
					"event/user.human.*",
					"event/user.*",
					"event",
				},
			},
		},
		{
			"execution set, ok",
			args{
				eventType: execution_rp.SetEventV2Type,
			},
			res{
				groups: []string{
					"event/execution.v2.set",
					"event/execution.v2.*",
					"event/execution.*",
					"event",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.res.groups, idsForEventType(string(tt.args.eventType)))
		})
	}
}

func TestActionProjection_reduces(t *testing.T) {
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockQueue) (fields, args, want)
	}{
		{
			name: "reduce, action, error",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, q *mock.MockQueue) (f fields, a args, w want) {
				queries.EXPECT().TargetsByExecutionID(gomock.Any(), gomock.Any()).Return(nil, zerrors.ThrowInternal(nil, "QUERY-37ardr0pki", "Errors.Query.CloseRows"))
				return fields{
						queries: queries,
						queue:   q,
					}, args{
						event: &action.AddedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   eventID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           action.AddedEventType,
								Data:          []byte(eventData),
								EditorUser:    userID,
								Seq:           1,
								AggregateType: action.AggregateType,
								Version:       action.AggregateVersion,
							}),
							Name:          "name",
							Script:        "name(){}",
							Timeout:       3 * time.Second,
							AllowedToFail: true,
						},
						mapper: action.AddedEventMapper,
					}, want{
						err: func(tt assert.TestingT, err error, i ...interface{}) bool {
							return errors.Is(err, zerrors.ThrowInternal(nil, "QUERY-37ardr0pki", "Errors.Query.CloseRows"))
						},
					}
			},
		},

		{
			name: "reduce, action, none",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, q *mock.MockQueue) (f fields, a args, w want) {
				queries.EXPECT().TargetsByExecutionID(gomock.Any(), gomock.Any()).Return([]*query.ExecutionTarget{}, nil)
				return fields{
						queries: queries,
						queue:   q,
					}, args{
						event: &action.AddedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   eventID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           action.AddedEventType,
								Data:          []byte(eventData),
								EditorUser:    userID,
								Seq:           1,
								AggregateType: action.AggregateType,
								Version:       action.AggregateVersion,
							}),
							Name:          "name",
							Script:        "name(){}",
							Timeout:       3 * time.Second,
							AllowedToFail: true,
						},
						mapper: action.AddedEventMapper,
					}, want{
						noOperation: true,
					}
			},
		},
		{
			name: "reduce, action, single",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, q *mock.MockQueue) (f fields, a args, w want) {
				targets := mockTargets(1)
				queries.EXPECT().TargetsByExecutionID(gomock.Any(), gomock.Any()).Return(targets, nil)
				createdAt := time.Now().UTC()
				q.EXPECT().Insert(
					gomock.Any(),
					&execution_rp.Request{
						Aggregate: &eventstore.Aggregate{
							InstanceID:    instanceID,
							Type:          action.AggregateType,
							Version:       action.AggregateVersion,
							ID:            eventID,
							ResourceOwner: orgID,
						},
						Sequence:    1,
						CreatedAt:   createdAt,
						EventType:   action.AddedEventType,
						UserID:      userID,
						EventData:   []byte(eventData),
						TargetsData: mockTargetsToBytes(targets),
					},
					gomock.Any(),
				).Return(nil)
				return fields{
						queries: queries,
						queue:   q,
					}, args{
						event: &action.AddedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   eventID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  createdAt,
								Typ:           action.AddedEventType,
								Data:          []byte(eventData),
								EditorUser:    userID,
								Seq:           1,
								AggregateType: action.AggregateType,
								Version:       action.AggregateVersion,
							}),
							Name:          "name",
							Script:        "name(){}",
							Timeout:       3 * time.Second,
							AllowedToFail: true,
						},
						mapper: action.AddedEventMapper,
					}, w
			},
		},
		{
			name: "reduce, action, multiple",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, q *mock.MockQueue) (f fields, a args, w want) {
				targets := mockTargets(3)
				queries.EXPECT().TargetsByExecutionID(gomock.Any(), gomock.Any()).Return(targets, nil)
				createdAt := time.Now().UTC()
				q.EXPECT().Insert(
					gomock.Any(),
					&execution_rp.Request{
						Aggregate: &eventstore.Aggregate{
							InstanceID:    instanceID,
							Type:          action.AggregateType,
							Version:       action.AggregateVersion,
							ID:            eventID,
							ResourceOwner: orgID,
						},
						Sequence:    1,
						CreatedAt:   createdAt,
						EventType:   action.AddedEventType,
						UserID:      userID,
						EventData:   []byte(eventData),
						TargetsData: mockTargetsToBytes(targets),
					},
					gomock.Any(),
				).Return(nil)
				return fields{
						queries: queries,
						queue:   q,
					}, args{
						event: &action.AddedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   eventID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  createdAt,
								Typ:           action.AddedEventType,
								Data:          []byte(eventData),
								EditorUser:    userID,
								Seq:           1,
								AggregateType: action.AggregateType,
								Version:       action.AggregateVersion,
							}),
							Name:          "name",
							Script:        "name(){}",
							Timeout:       3 * time.Second,
							AllowedToFail: true,
						},
						mapper: action.AddedEventMapper,
					}, w
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			queries := mock.NewMockQueries(ctrl)
			queue := mock.NewMockQueue(ctrl)
			f, a, w := tt.test(ctrl, queries, queue)

			event, err := a.mapper(a.event)
			assert.NoError(t, err)

			stmt, err := newEventExecutionsHandler(queries, f).reduce(event)
			if w.err != nil {
				w.err(t, err)
				return
			}
			assert.NoError(t, err)

			if w.noOperation {
				assert.Nil(t, stmt.Execute)
				return
			}
			err = stmt.Execute(t.Context(), nil, "")
			if w.stmtErr != nil {
				w.stmtErr(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func mockTarget() *query.ExecutionTarget {
	return &query.ExecutionTarget{
		InstanceID:       "instanceID",
		ExecutionID:      "executionID",
		TargetID:         "targetID",
		TargetType:       domain.TargetTypeWebhook,
		Endpoint:         "endpoint",
		Timeout:          time.Minute,
		InterruptOnError: true,
		SigningKey:       "key",
	}
}

func mockTargets(count int) []*query.ExecutionTarget {
	var targets []*query.ExecutionTarget
	if count > 0 {
		targets = make([]*query.ExecutionTarget, count)
		for i := range targets {
			targets[i] = mockTarget()
		}
	}
	return targets
}

func mockTargetsToBytes(targets []*query.ExecutionTarget) []byte {
	data, _ := json.Marshal(targets)
	return data
}

func newEventExecutionsHandler(queries *mock.MockQueries, f fields) *eventHandler {
	return &eventHandler{
		queue: f.queue,
		query: queries,
	}
}
