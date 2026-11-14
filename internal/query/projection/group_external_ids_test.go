package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func Test_GroupExternalIDsReduces(t *testing.T) {
	projection := &groupExternalIDsProjection{}
	tests := []struct {
		name      string
		event     func(*testing.T) eventstore.Event
		reduce    func(eventstore.Event) (*handler.Statement, error)
		statement string
		args      []interface{}
	}{
		{name: "persist source binding", event: getEvent(testEvent(group.GroupExternalIDSetEventType, group.AggregateType, []byte(`{"scope":"entra","externalId":"object-id"}`)), eventstore.GenericEventMapper[group.GroupExternalIDSetEvent]), reduce: projection.reduceExternalIDSet, statement: "INSERT INTO projections.group_external_ids1 (instance_id, group_id, resource_owner, scope, external_id) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (instance_id, group_id, scope) DO UPDATE SET (resource_owner, external_id) = (EXCLUDED.resource_owner, EXCLUDED.external_id)", args: []interface{}{"instance-id", "agg-id", "ro-id", "entra", "object-id"}},
		{name: "remove only one source", event: getEvent(testEvent(group.GroupExternalIDSetEventType, group.AggregateType, []byte(`{"scope":"entra","externalId":""}`)), eventstore.GenericEventMapper[group.GroupExternalIDSetEvent]), reduce: projection.reduceExternalIDSet, statement: "DELETE FROM projections.group_external_ids1 WHERE (instance_id = $1) AND (group_id = $2) AND (scope = $3)", args: []interface{}{"instance-id", "agg-id", "entra"}},
		{name: "delete group bindings", event: getEvent(testEvent(group.GroupRemovedEventType, group.AggregateType, nil), eventstore.GenericEventMapper[group.GroupRemovedEvent]), reduce: projection.reduceGroupRemoved, statement: "DELETE FROM projections.group_external_ids1 WHERE (group_id = $1) AND (resource_owner = $2) AND (instance_id = $3)", args: []interface{}{"agg-id", "ro-id", "instance-id"}},
		{name: "delete organization bindings", event: getEvent(testEvent(org.OrgRemovedEventType, org.AggregateType, nil), org.OrgRemovedEventMapper), reduce: projection.reduceOwnerRemoved, statement: "DELETE FROM projections.group_external_ids1 WHERE (resource_owner = $1) AND (instance_id = $2)", args: []interface{}{"agg-id", "instance-id"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := tt.event(t)
			statement, err := tt.reduce(event)
			assertReduce(t, statement, err, GroupExternalIDsProjectionTable, wantReduce{aggregateType: event.Aggregate().Type, sequence: 15, executer: &testExecuter{executions: []execution{{expectedStmt: tt.statement, expectedArgs: tt.args}}}})
		})
	}
}
