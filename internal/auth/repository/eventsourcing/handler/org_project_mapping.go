package handler

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
	"github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
	view_model "github.com/zitadel/zitadel/internal/project/repository/view/model"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
)

const (
	orgProjectMappingTable = "auth.org_project_mapping2"
)

type OrgProjectMapping struct {
	handler
	subscription *v1.Subscription
}

func newOrgProjectMapping(
	ctx context.Context,
	handler handler,
) *OrgProjectMapping {
	h := &OrgProjectMapping{
		handler: handler,
	}

	h.subscribe(ctx)

	return h
}

func (p *OrgProjectMapping) subscribe(ctx context.Context) {
	p.subscription = p.es.Subscribe(p.AggregateTypes()...)
	go func() {
		for event := range p.subscription.Events {
			query.ReduceEvent(ctx, p, event)
		}
	}()
}

func (p *OrgProjectMapping) ViewModel() string {
	return orgProjectMappingTable
}

func (p *OrgProjectMapping) Subscription() *v1.Subscription {
	return p.subscription
}

func (_ *OrgProjectMapping) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{project.AggregateType, instance.AggregateType}
}

func (p *OrgProjectMapping) CurrentSequence(instanceID string) (uint64, error) {
	sequence, err := p.view.GetLatestOrgProjectMappingSequence(instanceID)
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *OrgProjectMapping) EventQuery(instanceIDs []string) (*es_models.SearchQuery, error) {
	sequences, err := p.view.GetLatestOrgProjectMappingSequences(instanceIDs)
	if err != nil {
		return nil, err
	}
	return newSearchQuery(sequences, p.AggregateTypes(), instanceIDs), nil
}

func (p *OrgProjectMapping) Reduce(event *es_models.Event) (err error) {
	mapping := new(view_model.OrgProjectMapping)
	switch eventstore.EventType(event.Type) {
	case project.ProjectAddedType:
		mapping.OrgID = event.ResourceOwner
		mapping.ProjectID = event.AggregateID
		mapping.InstanceID = event.InstanceID
	case project.ProjectRemovedType:
		err := p.view.DeleteOrgProjectMappingsByProjectID(event.AggregateID, event.InstanceID)
		if err == nil {
			return p.view.ProcessedOrgProjectMappingSequence(event)
		}
	case project.GrantAddedType:
		projectGrant := new(view_model.ProjectGrant)
		err := projectGrant.SetData(event)
		if err != nil {
			return err
		}
		mapping.OrgID = projectGrant.GrantedOrgID
		mapping.ProjectID = event.AggregateID
		mapping.ProjectGrantID = projectGrant.GrantID
		mapping.InstanceID = event.InstanceID
	case project.GrantRemovedType:
		projectGrant := new(view_model.ProjectGrant)
		err := projectGrant.SetData(event)
		if err != nil {
			return err
		}
		err = p.view.DeleteOrgProjectMappingsByProjectGrantID(event.AggregateID, event.InstanceID)
		if err == nil {
			return p.view.ProcessedOrgProjectMappingSequence(event)
		}
	case instance.InstanceRemovedEventType:
		return p.view.DeleteInstanceOrgProjectMappings(event)
	case org.OrgRemovedEventType:
		return p.view.UpdateOwnerRemovedOrgProjectMappings(event)
	default:
		return p.view.ProcessedOrgProjectMappingSequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutOrgProjectMapping(mapping, event)
}

func (p *OrgProjectMapping) OnError(event *es_models.Event, err error) error {
	logging.WithFields("id", event.AggregateID).WithError(err).Warn("something went wrong in org project mapping handler")
	return spooler.HandleError(event, err, p.view.GetLatestOrgProjectMappingFailedEvent, p.view.ProcessedOrgProjectMappingFailedEvent, p.view.ProcessedOrgProjectMappingSequence, p.errorCountUntilSkip)
}

func (p *OrgProjectMapping) OnSuccess(instanceIDs []string) error {
	return spooler.HandleSuccess(p.view.UpdateOrgProjectMappingSpoolerRunTimestamp, instanceIDs)
}
