package handler

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	org_model "github.com/caos/zitadel/internal/org/repository/view/model"

	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing"
)

const (
	orgTable = "management.orgs"
)

type Org struct {
	handler
}

func NewOrg(h handler) *Org {
	o := &Org{h}
	go o.awaitEvents(context.Background())
	return o
}

func (o *Org) awaitEvents(ctx context.Context) {
	feed := make(chan *es_models.Event)
	subscriber := eventstore.Subscribe(feed, model.OrgAggregate)
	for {
		select {
		case <-ctx.Done():
			subscriber.Unsubscribe()
			return
		case event := <-feed:
			err := o.Reduce(event)
			if err != nil {
				err = o.OnError(event, err)
				logging.Log("HANDL-FOdlj").OnError(err).Info("unable to save failed event")
			}
		}
	}
}

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

func (o *Org) Reduce(event *es_models.Event) error {
	org := new(org_model.OrgView)

	switch event.Type {
	case model.OrgAdded:
		org.AppendEvent(event)
	case model.OrgChanged:
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
