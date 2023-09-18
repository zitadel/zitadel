package view

import (
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func OrgByIDQuery(id, instanceID string, latestSequence uint64) (*eventstore.SearchQueryBuilder, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dke74", "id should be filled")
	}
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(instanceID).
		AwaitOpenTransactions().
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(id).
		SequenceGreater(latestSequence).
		EventTypes(
			org.OrgAddedEventType,
			org.OrgChangedEventType,
			org.OrgDeactivatedEventType,
			org.OrgReactivatedEventType,
			org.OrgDomainAddedEventType,
			org.OrgDomainVerificationAddedEventType,
			org.OrgDomainVerifiedEventType,
			org.OrgDomainPrimarySetEventType,
			org.OrgDomainRemovedEventType,
			org.DomainPolicyAddedEventType,
			org.DomainPolicyChangedEventType,
			org.DomainPolicyRemovedEventType,
		).
		Builder(), nil
}
