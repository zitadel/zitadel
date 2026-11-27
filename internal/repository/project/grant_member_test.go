package project

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func TestGrantMemberEvents_Fields(t *testing.T) {
	ctx := context.Background()
	agg := &NewAggregate("project-id", "org-id").Aggregate
	agg.InstanceID = "instance-id"

	member := eventstore.Object{Type: ProjectGrantMemberRoleSearchType, ID: "user-id", Revision: 1}
	setRole := func(role string) *eventstore.FieldOperation {
		return eventstore.SetField(agg, member, "grant-id",
			&eventstore.Value{Value: role, ShouldIndex: true},
			eventstore.FieldTypeInstanceID,
			eventstore.FieldTypeResourceOwner,
			eventstore.FieldTypeAggregateType,
			eventstore.FieldTypeAggregateID,
			eventstore.FieldTypeObjectType,
			eventstore.FieldTypeObjectID,
			eventstore.FieldTypeFieldName,
			eventstore.FieldTypeValue,
		)
	}
	// removes the roles of the user in this grant only, not in other grants of the project
	removeMembership := &eventstore.FieldOperation{Remove: map[eventstore.FieldType]any{
		eventstore.FieldTypeInstanceID:     "instance-id",
		eventstore.FieldTypeResourceOwner:  "org-id",
		eventstore.FieldTypeAggregateType:  eventstore.AggregateType(AggregateType),
		eventstore.FieldTypeAggregateID:    "project-id",
		eventstore.FieldTypeObjectType:     ProjectGrantMemberRoleSearchType,
		eventstore.FieldTypeObjectID:       "user-id",
		eventstore.FieldTypeObjectRevision: uint8(1),
		eventstore.FieldTypeFieldName:      "grant-id",
	}}

	tests := []struct {
		name  string
		event interface {
			Fields() []*eventstore.FieldOperation
		}
		want []*eventstore.FieldOperation
	}{
		{
			name:  "added sets every role",
			event: NewProjectGrantMemberAddedEvent(ctx, agg, "user-id", "grant-id", "PROJECT_GRANT_OWNER", "PROJECT_GRANT_VIEWER"),
			want:  []*eventstore.FieldOperation{setRole("PROJECT_GRANT_OWNER"), setRole("PROJECT_GRANT_VIEWER")},
		},
		{
			name:  "changed replaces the roles of the membership",
			event: NewProjectGrantMemberChangedEvent(ctx, agg, "user-id", "grant-id", "PROJECT_GRANT_VIEWER"),
			want:  []*eventstore.FieldOperation{removeMembership, setRole("PROJECT_GRANT_VIEWER")},
		},
		{
			name:  "changed without roles removes the membership roles",
			event: NewProjectGrantMemberChangedEvent(ctx, agg, "user-id", "grant-id"),
			want:  []*eventstore.FieldOperation{removeMembership},
		},
		{
			name:  "removed",
			event: NewProjectGrantMemberRemovedEvent(ctx, agg, "user-id", "grant-id"),
			want:  []*eventstore.FieldOperation{removeMembership},
		},
		{
			name:  "cascade removed",
			event: NewProjectGrantMemberCascadeRemovedEvent(ctx, agg, "user-id", "grant-id"),
			want:  []*eventstore.FieldOperation{removeMembership},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.event.Fields())
		})
	}
}

func TestGrantRemovedEvent_Fields_removesMemberRoles(t *testing.T) {
	agg := &NewAggregate("project-id", "org-id").Aggregate
	agg.InstanceID = "instance-id"

	fields := NewGrantRemovedEvent(context.Background(), agg, "grant-id", "granted-org-id").Fields()

	// the roles of all members of the grant, whatever user they belong to
	assert.Contains(t, fields, &eventstore.FieldOperation{Remove: map[eventstore.FieldType]any{
		eventstore.FieldTypeInstanceID:     "instance-id",
		eventstore.FieldTypeResourceOwner:  "org-id",
		eventstore.FieldTypeAggregateType:  eventstore.AggregateType(AggregateType),
		eventstore.FieldTypeAggregateID:    "project-id",
		eventstore.FieldTypeObjectType:     ProjectGrantMemberRoleSearchType,
		eventstore.FieldTypeObjectRevision: uint8(1),
		eventstore.FieldTypeFieldName:      "grant-id",
	}})
}
