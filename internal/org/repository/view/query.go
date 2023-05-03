package view

import (
	"github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func OrgByIDQuery(id, instanceID string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dke74", "id should be filled")
	}
	return es_models.NewSearchQuery().
		AddQuery().
		AggregateTypeFilter(org.AggregateType).
		LatestSequenceFilter(latestSequence).
		InstanceIDFilter(instanceID).
		AggregateIDFilter(id).
		EventTypesFilter(
			es_models.EventType(org.OrgAddedEventType),
			es_models.EventType(org.OrgChangedEventType),
			es_models.EventType(org.OrgDeactivatedEventType),
			es_models.EventType(org.OrgReactivatedEventType),
			es_models.EventType(org.OrgDomainAddedEventType),
			es_models.EventType(org.OrgDomainVerificationAddedEventType),
			es_models.EventType(org.OrgDomainVerifiedEventType),
			es_models.EventType(org.OrgDomainPrimarySetEventType),
			es_models.EventType(org.OrgDomainRemovedEventType),
			es_models.EventType(org.DomainPolicyAddedEventType),
			es_models.EventType(org.DomainPolicyChangedEventType),
			es_models.EventType(org.DomainPolicyRemovedEventType),
		).
		SearchQuery(), nil
}
