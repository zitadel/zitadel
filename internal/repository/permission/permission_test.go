package permission

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func TestRemovedEvent_Fields(t *testing.T) {
	aggregate := NewAggregate("instance-id")
	event := NewRemovedEvent(context.Background(), aggregate, "IAM_OWNER", "org.read")

	assert.Equal(t, []*eventstore.FieldOperation{
		{
			Remove: map[eventstore.FieldType]any{
				eventstore.FieldTypeInstanceID:     aggregate.InstanceID,
				eventstore.FieldTypeResourceOwner:  aggregate.ResourceOwner,
				eventstore.FieldTypeAggregateType:  aggregate.Type,
				eventstore.FieldTypeAggregateID:    aggregate.ID,
				eventstore.FieldTypeObjectType:     RolePermissionType,
				eventstore.FieldTypeObjectID:       "IAM_OWNER",
				eventstore.FieldTypeObjectRevision: RolePermissionRevision,
				eventstore.FieldTypeFieldName:      PermissionSearchField,
				eventstore.FieldTypeValue:          "org.read",
			},
		},
	}, event.Fields())
}
