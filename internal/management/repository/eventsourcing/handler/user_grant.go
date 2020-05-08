package handler

import (
	es_model "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/usergrant/repository/eventsourcing"
	view_model "github.com/caos/zitadel/internal/usergrant/repository/view/model"
)

type UserGrant struct {
	handler
	eventstore eventstore.Eventstore
}

const (
	userGrantTable = "management.user_grants"
)

func (p *UserGrant) MinimumCycleDuration() time.Duration { return p.cycleDuration }

func (p *UserGrant) ViewModel() string {
	return userGrantTable
}

func (p *UserGrant) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestUserGrantSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.UserGrantQuery(sequence), nil
}

func (p *UserGrant) Process(event *models.Event) (err error) {
	grant := new(view_model.UserGrantView)
	switch event.Type {
	case es_model.UserGrantAdded:
		grant.AppendEvent(event)
	case es_model.UserGrantChanged,
		es_model.UserGrantDeactivated,
		es_model.UserGrantReactivated:
		grant, err = p.view.UserGrantByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = grant.AppendEvent(event)
	case es_model.UserGrantRemoved:
		err = p.view.DeleteUserGrant(event.AggregateID, event.Sequence)
	default:
		return p.view.ProcessedUserGrantSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return p.view.PutUserGrant(grant)
}

func (p *UserGrant) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-8is4s", "id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, p.view.GetLatestUserGrantFailedEvent, p.view.ProcessedUserGrantFailedEvent, p.view.ProcessedUserGrantSequence, p.errorCountUntilSkip)
}
