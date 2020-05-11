package handler

import (
	"context"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	usr_events "github.com/caos/zitadel/internal/user/repository/eventsourcing"
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
	eventstore    eventstore.Eventstore
	projectEvents *proj_event.ProjectEventstore
	userEvents    *usr_events.UserEventstore
}

const (
	userGrantTable = "management.user_grants"
)

func (u *UserGrant) MinimumCycleDuration() time.Duration { return u.cycleDuration }

func (u *UserGrant) ViewModel() string {
	return userGrantTable
}

func (u *UserGrant) EventQuery() (*models.SearchQuery, error) {
	sequence, err := u.view.GetLatestUserGrantSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.UserGrantQuery(sequence), nil
}

func (u *UserGrant) Process(event *models.Event) (err error) {
	grant := new(view_model.UserGrantView)
	switch event.Type {
	case es_model.UserGrantAdded:
		grant.AppendEvent(event)
		if err != nil {
			return err
		}
		err = u.fillData(grant)
	case es_model.UserGrantChanged,
		es_model.UserGrantDeactivated,
		es_model.UserGrantReactivated:
		grant, err = u.view.UserGrantByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = grant.AppendEvent(event)
	case es_model.UserGrantRemoved:
		err = u.view.DeleteUserGrant(event.AggregateID, event.Sequence)
	default:
		return u.view.ProcessedUserGrantSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return u.view.PutUserGrant(grant)
}

func (u *UserGrant) fillData(grant *view_model.UserGrantView) (err error) {
	err = u.fillUser(grant)
	if err != nil {
		return err
	}
	err = u.fillProject(grant)
	if err != nil {
		return err
	}
	return u.fillOrg(grant)
}

func (u *UserGrant) fillUser(grant *view_model.UserGrantView) error {
	user, err := u.userEvents.UserByID(context.Background(), grant.UserID)
	if err != nil {
		return err
	}
	grant.UserName = user.UserName
	grant.FirstName = user.FirstName
	grant.LastName = user.LastName
	grant.Email = user.EmailAddress
	return nil
}

func (u *UserGrant) fillProject(grant *view_model.UserGrantView) error {
	project, err := u.projectEvents.ProjectByID(context.Background(), grant.ProjectID)
	if err != nil {
		return err
	}
	grant.ProjectName = project.Name
	return nil
}

func (u *UserGrant) fillOrg(grant *view_model.UserGrantView) error {
	//TODO: get ORG
	return nil
}

func (u *UserGrant) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-8is4s", "id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserGrantFailedEvent, u.view.ProcessedUserGrantFailedEvent, u.view.ProcessedUserGrantSequence, u.errorCountUntilSkip)
}
