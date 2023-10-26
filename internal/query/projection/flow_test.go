package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func TestFlowProjection_reduces(t *testing.T) {
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
			name: "reduceTriggerActionsSetEventType",
			args: args{
				event: getEvent(
					testEvent(
						org.TriggerActionsSetEventType,
						org.AggregateType,
						[]byte(`{"flowType": 1, "triggerType": 1, "actionIDs": ["id1", "id2"]}`),
					), org.TriggerActionsSetEventMapper),
			},
			reduce: (&flowProjection{}).reduceTriggerActionsSetEventType,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.flow_triggers2 WHERE (flow_type = $1) AND (trigger_type = $2) AND (resource_owner = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								domain.FlowTypeExternalAuthentication,
								domain.TriggerTypePostAuthentication,
								"ro-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.flow_triggers2 (resource_owner, instance_id, flow_type, change_date, sequence, trigger_type, action_id, trigger_sequence) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"ro-id",
								"instance-id",
								domain.FlowTypeExternalAuthentication,
								anyArg{},
								uint64(15),
								domain.TriggerTypePostAuthentication,
								"id1",
								0,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.flow_triggers2 (resource_owner, instance_id, flow_type, change_date, sequence, trigger_type, action_id, trigger_sequence) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"ro-id",
								"instance-id",
								domain.FlowTypeExternalAuthentication,
								anyArg{},
								uint64(15),
								domain.TriggerTypePostAuthentication,
								"id2",
								1,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceFlowClearedEventType",
			args: args{
				event: getEvent(
					testEvent(
						org.FlowClearedEventType,
						org.AggregateType,
						[]byte(`{"flowType": 1}`),
					), org.FlowClearedEventMapper),
			},
			reduce: (&flowProjection{}).reduceFlowClearedEventType,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.flow_triggers2 WHERE (flow_type = $1) AND (resource_owner = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								domain.FlowTypeExternalAuthentication,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceOwnerRemoved",
			reduce: (&flowProjection{}).reduceOwnerRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.flow_triggers2 WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(FlowInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.flow_triggers2 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := baseEvent(t)
			got, err := tt.reduce(event)
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, FlowTriggerTable, tt.want)
		})
	}
}
