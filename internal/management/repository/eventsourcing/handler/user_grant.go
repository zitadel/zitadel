package handler

import (
	"context"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_events "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_events "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	grant_es_model "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	view_model "github.com/caos/zitadel/internal/usergrant/repository/view/model"
)

type UserGrant struct {
	handler
	eventstore    eventstore.Eventstore
	projectEvents *proj_event.ProjectEventstore
	userEvents    *usr_events.UserEventstore
	orgEvents     *org_events.OrgEventstore
}

const (
	userGrantTable = "management.user_grants"
)

func (u *UserGrant) ViewModel() string {
	return userGrantTable
}

func (u *UserGrant) EventQuery() (*models.SearchQuery, error) {
	sequence, err := u.view.GetLatestUserGrantSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(grant_es_model.UserGrantAggregate, usr_es_model.UserAggregate, proj_es_model.ProjectAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (u *UserGrant) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case grant_es_model.UserGrantAggregate:
		err = u.processUserGrant(event)
	case usr_es_model.UserAggregate:
		err = u.processUser(event)
	case proj_es_model.ProjectAggregate:
		err = u.processProject(event)
	}
	return err
}

func (u *UserGrant) processUserGrant(event *models.Event) (err error) {
	grant := new(view_model.UserGrantView)
	switch event.Type {
	case grant_es_model.UserGrantAdded:
		err = grant.AppendEvent(event)
		if err != nil {
			return err
		}
		err = u.fillData(grant, event.ResourceOwner)
	case grant_es_model.UserGrantChanged,
		grant_es_model.UserGrantCascadeChanged,
		grant_es_model.UserGrantDeactivated,
		grant_es_model.UserGrantReactivated:
		grant, err = u.view.UserGrantByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = grant.AppendEvent(event)
	case grant_es_model.UserGrantRemoved, grant_es_model.UserGrantCascadeRemoved:
		return u.view.DeleteUserGrant(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return u.view.ProcessedUserGrantSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return u.view.PutUserGrant(grant, grant.Sequence, event.CreationDate)
}

func (u *UserGrant) processUser(event *models.Event) (err error) {
	switch event.Type {
	case usr_es_model.UserProfileChanged,
		usr_es_model.UserEmailChanged,
		usr_es_model.HumanProfileChanged,
		usr_es_model.HumanEmailChanged,
		usr_es_model.MachineChanged:
		grants, err := u.view.UserGrantsByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(grants) == 0 {
			return u.view.ProcessedUserGrantSequence(event.Sequence, event.CreationDate)
		}
		user, err := u.userEvents.UserByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		for _, grant := range grants {
			u.fillUserData(grant, user)
		}
		return u.view.PutUserGrants(grants, event.Sequence, event.CreationDate)
	default:
		return u.view.ProcessedUserGrantSequence(event.Sequence, event.CreationDate)
	}
	return nil
}

func (u *UserGrant) processProject(event *models.Event) (err error) {
	switch event.Type {
	case proj_es_model.ProjectChanged:
		grants, err := u.view.UserGrantsByProjectID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(grants) == 0 {
			return u.view.ProcessedUserGrantSequence(event.Sequence, event.CreationDate)
		}
		project, err := u.projectEvents.ProjectByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		for _, grant := range grants {
			u.fillProjectData(grant, project)
		}
		return u.view.PutUserGrants(grants, event.Sequence, event.CreationDate)
	default:
		return u.view.ProcessedUserGrantSequence(event.Sequence, event.CreationDate)
	}
	return nil
}

func (u *UserGrant) fillData(grant *view_model.UserGrantView, resourceOwner string) (err error) {
	user, err := u.userEvents.UserByID(context.Background(), grant.UserID)
	if err != nil {
		return err
	}
	u.fillUserData(grant, user)
	project, err := u.projectEvents.ProjectByID(context.Background(), grant.ProjectID)
	if err != nil {
		return err
	}
	u.fillProjectData(grant, project)

	org, err := u.orgEvents.OrgByID(context.TODO(), org_model.NewOrg(resourceOwner))
	if err != nil {
		return err
	}
	u.fillOrgData(grant, org)
	return nil
}

func (u *UserGrant) fillUserData(grant *view_model.UserGrantView, user *usr_model.User) {
	grant.UserName = user.UserName
	if user.Human != nil {
		grant.FirstName = user.FirstName
		grant.LastName = user.LastName
		grant.DisplayName = user.FirstName + " " + user.LastName
		grant.Email = user.EmailAddress
	}
	if user.Machine != nil {
		grant.DisplayName = user.Machine.Name
	}
}

func (u *UserGrant) fillProjectData(grant *view_model.UserGrantView, project *proj_model.Project) {
	grant.ProjectName = project.Name
	grant.ProjectOwner = project.ResourceOwner
}

func (u *UserGrant) fillOrgData(grant *view_model.UserGrantView, org *org_model.Org) {
	grant.OrgName = org.Name
	for _, domain := range org.Domains {
		if domain.Primary {
			grant.OrgPrimaryDomain = domain.Domain
			break
		}
	}
}

func (u *UserGrant) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-8is4s", "id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserGrantFailedEvent, u.view.ProcessedUserGrantFailedEvent, u.view.ProcessedUserGrantSequence, u.errorCountUntilSkip)
}

func (u *UserGrant) OnSuccess() error {
	return spooler.HandleSuccess(u.view.UpdateUserGrantSpoolerRunTimestamp)
}
