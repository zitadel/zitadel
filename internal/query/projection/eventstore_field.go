package projection

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
)

const (
	fieldsProjectGrant      = "project_grant_fields"
	fieldsOrgDomainVerified = "org_domain_verified_fields"
	fieldsInstanceDomain    = "instance_domain_fields"
)

func newFillProjectGrantFields(config handler.Config) *handler.FieldHandler {
	return handler.NewFieldHandler(
		&config,
		fieldsProjectGrant,
		map[eventstore.AggregateType][]eventstore.EventType{
			org.AggregateType:     nil,
			project.AggregateType: nil,
		},
	)
}

func newFillOrgDomainVerifiedFields(config handler.Config) *handler.FieldHandler {
	return handler.NewFieldHandler(
		&config,
		fieldsOrgDomainVerified,
		map[eventstore.AggregateType][]eventstore.EventType{
			org.AggregateType: {
				org.OrgDomainAddedEventType,
				org.OrgDomainVerifiedEventType,
				org.OrgDomainRemovedEventType,
			},
		},
	)
}

func newFillInstanceDomainFields(config handler.Config) *handler.FieldHandler {
	return handler.NewFieldHandler(
		&config,
		fieldsInstanceDomain,
		map[eventstore.AggregateType][]eventstore.EventType{
			instance.AggregateType: {
				instance.InstanceDomainAddedEventType,
				instance.InstanceDomainRemovedEventType,
				instance.InstanceRemovedEventType,
			},
		},
	)
}
