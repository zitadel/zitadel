package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(OrgAdded, OrgAddedEventMapper).
		RegisterFilterEventMapper(OrgChanged, OrgChangedEventMapper).
		//RegisterFilterEventMapper(IDPOIDCConfigAddedEventType, IDPOIDCConfigAddedEventMapper).
		RegisterFilterEventMapper(OrgDomainAdded, DomainAddedEventMapper).
		RegisterFilterEventMapper(OrgDomainVerificationAdded, DomainVerificationAddedEventMapper).
		RegisterFilterEventMapper(OrgDomainVerificationFailed, DomainVerificationFailedEventMapper).
		RegisterFilterEventMapper(OrgDomainVerified, DomainVerifiedEventMapper).
		RegisterFilterEventMapper(OrgDomainPrimarySet, DomainPrimarySetEventMapper).
		RegisterFilterEventMapper(OrgDomainRemoved, DomainRemovedEventMapper)
}
