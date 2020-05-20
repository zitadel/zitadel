package handler

import (
	"time"

	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing"
	"github.com/caos/zitadel/internal/org/repository/view"
)

type Org struct {
	handler
}

const (
	orgTable = "management.orgs"
)

func (o *Org) MinimumCycleDuration() time.Duration { return o.cycleDuration }

func (o *Org) ViewModel() string {
	return orgTable
}

func (o *Org) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := o.view.GetLatestOrgSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.OrgQuery(sequence), nil
}

func (o *Org) Process(event *es_models.Event) error {
	org := new(view.OrgView)

	switch event.Type {
	case org_model.OrgAdded:
		org.AppendEvent(event)
	case org_model.OrgChanged:
		err := org.SetData(event)
		if err != nil {
			return err
		}
		org, err = o.view.OrgByID(org.ID)
		if err != nil {
			return err
		}
		err = org.AppendEvent(event)
		if err != nil {
			return err
		}
	default:
		return o.view.ProcessedOrgSequence(event.Sequence)
	}

	return o.view.PutOrg(org)
}

func (o *Org) OnError(event *es_models.Event, spoolerErr error) error {
	logging.LogWithFields("SPOOL-ls9ew", "id", event.AggregateID).WithError(spoolerErr).Warn("something went wrong in project app handler")
	return spooler.HandleError(event, spoolerErr, o.view.GetLatestOrgFailedEvent, o.view.ProcessedOrgFailedEvent, o.view.ProcessedOrgSequence, o.errorCountUntilSkip)
}
