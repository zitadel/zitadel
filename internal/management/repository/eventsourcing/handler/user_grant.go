package handler

import (
	"context"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	proj_view "github.com/caos/zitadel/internal/project/repository/view"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/user/repository/view"
	usr_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
	grant_es_model "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/usergrant/repository/view/model"
)

const (
	userGrantTable = "management.user_grants"
)

type UserGrant struct {
	handler
	subscription *v1.Subscription
}

func newUserGrant(
	handler handler,
) *UserGrant {
	h := &UserGrant{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *UserGrant) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (u *UserGrant) ViewModel() string {
	return userGrantTable
}

func (u *UserGrant) Subscription() *v1.Subscription {
	return u.subscription
}

func (_ *UserGrant) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{grant_es_model.UserGrantAggregate, usr_es_model.UserAggregate, proj_es_model.ProjectAggregate, org_es_model.OrgAggregate}
}

func (u *UserGrant) CurrentSequence() (uint64, error) {
	sequence, err := u.view.GetLatestUserGrantSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (u *UserGrant) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := u.view.GetLatestUserGrantSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(u.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (u *UserGrant) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case grant_es_model.UserGrantAggregate:
		err = u.processUserGrant(event)
	case usr_es_model.UserAggregate:
		err = u.processUser(event)
	case proj_es_model.ProjectAggregate:
		err = u.processProject(event)
	case org_es_model.OrgAggregate:
		err = u.processOrg(event)
	}
	return err
}

func (u *UserGrant) processUserGrant(event *es_models.Event) (err error) {
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
		return u.view.DeleteUserGrant(event.AggregateID, event)
	default:
		return u.view.ProcessedUserGrantSequence(event)
	}
	if err != nil {
		return err
	}
	return u.view.PutUserGrant(grant, event)
}

func (u *UserGrant) processUser(event *es_models.Event) (err error) {
	switch event.Type {
	case usr_es_model.UserProfileChanged,
		usr_es_model.UserEmailChanged,
		usr_es_model.HumanProfileChanged,
		usr_es_model.HumanEmailChanged,
		usr_es_model.MachineChanged,
		usr_es_model.HumanAvatarAdded,
		usr_es_model.HumanAvatarRemoved:
		grants, err := u.view.UserGrantsByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(grants) == 0 {
			return u.view.ProcessedUserGrantSequence(event)
		}
		user, err := u.getUserByID(event.AggregateID)
		if err != nil {
			return err
		}
		for _, grant := range grants {
			u.fillUserData(grant, user)
		}
		return u.view.PutUserGrants(grants, event)
	default:
		return u.view.ProcessedUserGrantSequence(event)
	}
}

func (u *UserGrant) processProject(event *es_models.Event) (err error) {
	switch event.Type {
	case proj_es_model.ProjectChanged:
		grants, err := u.view.UserGrantsByProjectID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(grants) == 0 {
			return u.view.ProcessedUserGrantSequence(event)
		}
		project, err := u.getProjectByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		for _, grant := range grants {
			u.fillProjectData(grant, project)
		}
		return u.view.PutUserGrants(grants, event)
	default:
		return u.view.ProcessedUserGrantSequence(event)
	}
}

func (u *UserGrant) processOrg(event *es_models.Event) (err error) {
	switch event.Type {
	case org_es_model.OrgChanged:
		grants, err := u.view.UserGrantsByOrgID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(grants) == 0 {
			return u.view.ProcessedUserGrantSequence(event)
		}
		org, err := u.getOrgByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		for _, grant := range grants {
			u.fillOrgData(grant, org)
		}
		return u.view.PutUserGrants(grants, event)
	default:
		return u.view.ProcessedUserGrantSequence(event)
	}
}

func (u *UserGrant) fillData(grant *view_model.UserGrantView, resourceOwner string) (err error) {
	user, err := u.getUserByID(grant.UserID)
	if err != nil {
		return err
	}
	u.fillUserData(grant, user)
	project, err := u.getProjectByID(context.Background(), grant.ProjectID)
	if err != nil {
		return err
	}
	u.fillProjectData(grant, project)

	org, err := u.getOrgByID(context.TODO(), resourceOwner)
	if err != nil {
		return err
	}
	u.fillOrgData(grant, org)
	return nil
}

func (u *UserGrant) fillUserData(grant *view_model.UserGrantView, user *usr_view_model.UserView) {
	grant.UserName = user.UserName
	grant.UserResourceOwner = user.ResourceOwner
	if user.HumanView != nil {
		grant.FirstName = user.FirstName
		grant.LastName = user.LastName
		grant.DisplayName = user.FirstName + " " + user.LastName
		grant.Email = user.Email
		grant.AvatarKey = user.AvatarKey
	}
	if user.MachineView != nil {
		grant.DisplayName = user.MachineView.Name
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

func (u *UserGrant) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-8is4s", "id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserGrantFailedEvent, u.view.ProcessedUserGrantFailedEvent, u.view.ProcessedUserGrantSequence, u.errorCountUntilSkip)
}

func (u *UserGrant) OnSuccess() error {
	return spooler.HandleSuccess(u.view.UpdateUserGrantSpoolerRunTimestamp)
}

func (u *UserGrant) getUserByID(userID string) (*usr_view_model.UserView, error) {
	user, usrErr := u.view.UserByID(userID)
	if usrErr != nil && !caos_errs.IsNotFound(usrErr) {
		return nil, usrErr
	}
	if user == nil {
		user = &usr_view_model.UserView{}
	}
	events, err := u.getUserEvents(userID, user.Sequence)
	if err != nil {
		return user, usrErr
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return user, nil
		}
	}
	if userCopy.State == int32(usr_model.UserStateDeleted) {
		return nil, caos_errs.ThrowNotFound(nil, "HANDLER-m9dos", "Errors.User.NotFound")
	}
	return &userCopy, nil
}

func (u *UserGrant) getUserEvents(userID string, sequence uint64) ([]*es_models.Event, error) {
	query, err := view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}

	return u.es.FilterEvents(context.Background(), query)
}

func (u *UserGrant) getOrgByID(ctx context.Context, orgID string) (*org_model.Org, error) {
	query, err := org_view.OrgByIDQuery(orgID, 0)
	if err != nil {
		return nil, err
	}

	esOrg := &org_es_model.Org{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: orgID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, esOrg.AppendEvents, query)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-kVLb2", "Errors.Org.NotFound")
	}

	return org_es_model.OrgToModel(esOrg), nil
}

func (u *UserGrant) getProjectByID(ctx context.Context, projID string) (*proj_model.Project, error) {
	query, err := proj_view.ProjectByIDQuery(projID, 0)
	if err != nil {
		return nil, err
	}
	esProject := &proj_es_model.Project{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: projID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, esProject.AppendEvents, query)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	if esProject.Sequence == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-Dfb42", "Errors.Project.NotFound")
	}

	return proj_es_model.ProjectToModel(esProject), nil
}
