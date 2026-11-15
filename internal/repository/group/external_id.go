package group

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const GroupExternalIDSetEventType = groupEventTypePrefix + "external_id.set"

// GroupExternalIDSetEvent binds a provisioning source's identity to a native group.
// An empty external ID removes the binding for that source.
type GroupExternalIDSetEvent struct {
	eventstore.BaseEvent `json:"-"`
	Scope                string `json:"scope"`
	ExternalID           string `json:"externalId"`
	OldExternalID        string `json:"oldExternalId,omitempty"`
}

func NewGroupExternalIDSetEvent(ctx context.Context, aggregate *eventstore.Aggregate, scope, oldID, externalID string) *GroupExternalIDSetEvent {
	return &GroupExternalIDSetEvent{BaseEvent: *eventstore.NewBaseEventForPush(ctx, aggregate, GroupExternalIDSetEventType), Scope: scope, ExternalID: externalID, OldExternalID: oldID}
}
func (e *GroupExternalIDSetEvent) SetBaseEvent(base *eventstore.BaseEvent) { e.BaseEvent = *base }
func (e *GroupExternalIDSetEvent) Payload() any                            { return e }
func (e *GroupExternalIDSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	var constraints []*eventstore.UniqueConstraint
	if e.OldExternalID != "" {
		constraints = append(constraints, NewRemoveGroupExternalIDUniqueConstraint(e.Aggregate().ID, e.Aggregate().ResourceOwner, e.Scope, e.OldExternalID))
	}
	if e.ExternalID != "" {
		constraints = append(constraints, eventstore.NewAddEventUniqueConstraint("group_external_id", externalIDConstraintKey(e.Aggregate().ResourceOwner, e.Scope, e.ExternalID), "Errors.Group.AlreadyExists").
			WithOwners(eventstore.OwnerTag(eventstore.UniqueConstraintOwnerOrg, e.Aggregate().ResourceOwner)).
			WithClaimant(AggregateType, e.Aggregate().ID))
	}
	return constraints
}
func NewRemoveGroupExternalIDUniqueConstraint(groupID, orgID, scope, externalID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint("group_external_id", externalIDConstraintKey(orgID, scope, externalID)).WithClaimant(AggregateType, groupID)
}
func externalIDConstraintKey(orgID, scope, externalID string) string {
	// The constraint store lowercases keys; hash an unambiguous tuple to preserve case-exact IDs.
	data, _ := json.Marshal([3]string{orgID, scope, externalID})
	return fmt.Sprintf("%x", sha256.Sum256(data))
}
