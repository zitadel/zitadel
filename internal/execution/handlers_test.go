package execution

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/action"
	execution_rp "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func Test_EventExecution(t *testing.T) {
	type args struct {
		event   eventstore.Event
		targets []*query.ExecutionTarget
	}
	type res struct {
		targets     []Target
		contextInfo *ContextInfoEvent
		columns     []handler.Column
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
				event: eventForEventExecution(
					"aggID",
					session.AggregateType,
					session.AggregateVersion,
					session.AddedType,
					time.Date(2024, 1, 1, 1, 1, 1, 1, time.Local),
					"userID",
					1,
					[]byte("{\"attr\":\"value\"}"),
				),
				targets: []*query.ExecutionTarget{{
					InstanceID:       "instanceID",
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
						InstanceID:       "instanceID",
						ExecutionID:      "executionID",
						TargetID:         "targetID",
						TargetType:       domain.TargetTypeWebhook,
						Endpoint:         "endpoint",
						Timeout:          time.Minute,
						InterruptOnError: true,
						SigningKey:       "key",
					},
				},
				contextInfo: &ContextInfoEvent{
					AggregateID:   "aggID",
					AggregateType: "session",
					ResourceOwner: "resourceOwner",
					InstanceID:    "instanceID",
					Version:       "v1",
					Sequence:      1,
					EventType:     "session.added",
					CreatedAt:     time.Date(2024, 1, 1, 1, 1, 1, 1, time.Local).Format(time.RFC3339),
					UserID:        "userID",
					EventPayload:  []byte("{\"attr\":\"value\"}"),
				},
				columns: []handler.Column{
					handler.NewCol(ExecutionInstanceID, "instanceID"),
					handler.NewCol(ExecutionResourceOwner, "resourceOwner"),
					handler.NewCol(ExecutionAggregateType, eventstore.AggregateType("session")),
					handler.NewCol(ExecutionAggregateVersion, eventstore.Version("v1")),
					handler.NewCol(ExecutionAggregateID, "aggID"),
					handler.NewCol(ExecutionSequence, uint64(1)),
					handler.NewCol(ExecutionEventType, eventstore.EventType("session.added")),
					handler.NewCol(ExecutionCreatedAt, time.Date(2024, 1, 1, 1, 1, 1, 1, time.Local)),
					handler.NewCol(ExecutionEventUserIDCol, "userID"),
					handler.NewCol(ExecutionEventDataCol, []byte("{\"attr\":\"value\"}")),
					handler.NewCol(ExecutionTargetsDataCol, []byte("[{\"InstanceID\":\"instanceID\",\"ExecutionID\":\"executionID\",\"TargetID\":\"targetID\",\"TargetType\":0,\"Endpoint\":\"endpoint\",\"Timeout\":60000000000,\"InterruptOnError\":true,\"SigningKey\":\"key\"}]")),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventExecution, err := NewEventExecution(tt.args.event, tt.args.targets)
			if tt.res.wantErr {
				assert.Error(t, err)
				assert.Nil(t, eventExecution)
				return
			}
			assert.NoError(t, err)
			targets, err := eventExecution.Targets()
			assert.NoError(t, err)
			assert.Equal(t, tt.res.targets, targets)
			assert.Equal(t, tt.res.contextInfo, eventExecution.ContextInfo())
			assert.Equal(t, tt.res.columns, eventExecution.Columns())
		})
	}
}

func eventForEventExecution(aggID, aggType, version, eventType string, creation time.Time, userID string, sequence uint64, data []byte) eventstore.Event {
	return &eventstore.BaseEvent{
		Agg: &eventstore.Aggregate{
			ID:            aggID,
			Type:          eventstore.AggregateType(aggType),
			ResourceOwner: "resourceOwner",
			InstanceID:    "instanceID",
			Version:       eventstore.Version(version),
		},
		EventType: eventstore.EventType(eventType),
		Seq:       sequence,
		Creation:  creation,
		User:      userID,
		Data:      data,
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
	type args struct {
		event func(t *testing.T) eventstore.Event
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.Event) (*handler.Statement, error)
		want   wantReduce
	}{
		{
			name: "reduce, action, error",
			args: args{
				event: getEvent(
					testEvent(
						action.AddedEventType,
						action.AggregateType,
						[]byte(`{"name": "name", "script":"name(){}","timeout": 3000000000, "allowedToFail": true}`),
					),
					action.AddedEventMapper,
				),
			},
			reduce: (newMockEventExecutionsHandler(newMockEventExecutionHandlerQuery(0, fmt.Errorf("failed query")))).reduce,
			want: wantReduce{
				err: func(err error) bool {
					return err.Error() == "failed query"
				},
			},
		},
		{
			name: "reduce, action, none",
			args: args{
				event: getEvent(
					testEvent(
						action.AddedEventType,
						action.AggregateType,
						[]byte(`{"name": "name", "script":"name(){}","timeout": 3000000000, "allowedToFail": true}`),
					),
					action.AddedEventMapper,
				),
			},
			reduce: (newMockEventExecutionsHandler(newMockEventExecutionHandlerQuery(0, nil))).reduce,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("action"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{},
				},
			},
		},
		{
			name: "reduce, action, single",
			args: args{
				event: getEvent(
					testEvent(
						action.AddedEventType,
						action.AggregateType,
						[]byte(`{"name": "name", "script":"name(){}","timeout": 3000000000, "allowedToFail": true}`),
					),
					action.AddedEventMapper,
				),
			},
			reduce: (newMockEventExecutionsHandler(newMockEventExecutionHandlerQuery(1, nil))).reduce,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("action"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.execution_handler (instance_id, resource_owner, aggregate_type, aggregate_version, aggregate_id, sequence, event_type, created_at, user_id, event_data, targets_data) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
							expectedArgs: []interface{}{
								"instance-id",
								"ro-id",
								eventstore.AggregateType(action.AggregateType),
								eventstore.Version(action.AggregateVersion),
								"agg-id",
								uint64(15),
								action.AddedEventType,
								anyArg{},
								"editor-user",
								[]byte(`{"name": "name", "script":"name(){}","timeout": 3000000000, "allowedToFail": true}`),
								anyArg{},
							},
						},
					},
				},
			},
		},
		{
			name: "reduce, action, multiple",
			args: args{
				event: getEvent(
					testEvent(
						action.AddedEventType,
						action.AggregateType,
						[]byte(`{"name": "name", "script":"name(){}","timeout": 3000000000, "allowedToFail": true}`),
					),
					action.AddedEventMapper,
				),
			},
			reduce: (newMockEventExecutionsHandler(newMockEventExecutionHandlerQuery(3, nil))).reduce,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("action"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.execution_handler (instance_id, resource_owner, aggregate_type, aggregate_version, aggregate_id, sequence, event_type, created_at, user_id, event_data, targets_data) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
							expectedArgs: []interface{}{
								"instance-id",
								"ro-id",
								eventstore.AggregateType(action.AggregateType),
								eventstore.Version(action.AggregateVersion),
								"agg-id",
								uint64(15),
								action.AddedEventType,
								anyArg{},
								"editor-user",
								[]byte(`{"name": "name", "script":"name(){}","timeout": 3000000000, "allowedToFail": true}`),
								mockTargetsToBytes(3),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := tt.args.event(t)
			got, err := tt.reduce(event)
			assertReduce(t, got, err, ExecutionHandlerTable, tt.want)
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

func mockTargetsToBytes(count int) []byte {
	targets := mockTargets(count)
	data, _ := json.Marshal(targets)
	return data
}

func newMockEventExecutionHandlerQuery(count int, err error) *mockEventExecutionHandlerQuery {
	return &mockEventExecutionHandlerQuery{
		targets: mockTargets(count),
		err:     err,
	}
}

type mockEventExecutionHandlerQuery struct {
	targets []*query.ExecutionTarget
	err     error
}

func (q *mockEventExecutionHandlerQuery) TargetsByExecutionID(ctx context.Context, ids []string) (execution []*query.ExecutionTarget, err error) {
	return q.targets, q.err
}

func newMockEventExecutionsHandler(query eventExecutionsHandlerQueries) *eventExecutionsHandler {
	return &eventExecutionsHandler{
		query: query,
	}
}
