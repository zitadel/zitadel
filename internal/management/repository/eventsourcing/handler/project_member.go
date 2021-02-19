package handler

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/user/repository/view"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	usr_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	projectMemberTable = "management.project_members"
)

type ProjectMember struct {
	handler
	userEvents   *usr_event.UserEventstore
	subscription *eventstore.Subscription
}

func newProjectMember(
	handler handler,
	userEvents *usr_event.UserEventstore,
) *ProjectMember {
	h := &ProjectMember{
		handler:    handler,
		userEvents: userEvents,
	}

	h.subscribe()

	return h
}

func (m *ProjectMember) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (p *ProjectMember) ViewModel() string {
	return projectMemberTable
}

func (_ *ProjectMember) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{proj_es_model.ProjectAggregate, usr_es_model.UserAggregate}
}

func (p *ProjectMember) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestProjectMemberSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *ProjectMember) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectMemberSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *ProjectMember) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case proj_es_model.ProjectAggregate:
		err = p.processProjectMember(event)
	case usr_es_model.UserAggregate:
		err = p.processUser(event)
	}
	return err
}

func (p *ProjectMember) processProjectMember(event *models.Event) (err error) {
	member := new(view_model.ProjectMemberView)
	switch event.Type {
	case proj_es_model.ProjectMemberAdded:
		err = member.AppendEvent(event)
		if err != nil {
			return err
		}
		p.fillData(member)
	case proj_es_model.ProjectMemberChanged:
		err = member.SetData(event)
		if err != nil {
			return err
		}
		member, err = p.view.ProjectMemberByIDs(event.AggregateID, member.UserID)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case proj_es_model.ProjectMemberRemoved:
		err = member.SetData(event)
		if err != nil {
			return err
		}
		return p.view.DeleteProjectMember(event.AggregateID, member.UserID, event)
	case proj_es_model.ProjectRemoved:
		return p.view.DeleteProjectMembersByProjectID(event.AggregateID)
	default:
		return p.view.ProcessedProjectMemberSequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutProjectMember(member, event)
}

func (p *ProjectMember) processUser(event *models.Event) (err error) {
	switch event.Type {
	case usr_es_model.UserProfileChanged,
		usr_es_model.UserEmailChanged,
		usr_es_model.HumanProfileChanged,
		usr_es_model.HumanEmailChanged,
		usr_es_model.MachineChanged:
		members, err := p.view.ProjectMembersByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(members) == 0 {
			return p.view.ProcessedProjectMemberSequence(event)
		}
		user, err := p.getUserByID(event.AggregateID)
		if err != nil {
			return err
		}
		for _, member := range members {
			p.fillUserData(member, user)
		}
		return p.view.PutProjectMembers(members, event)
	default:
		return p.view.ProcessedProjectMemberSequence(event)
	}
	return nil
}

func (p *ProjectMember) fillData(member *view_model.ProjectMemberView) (err error) {
	user, err := p.getUserByID(member.UserID)
	if err != nil {
		return err
	}
	p.fillUserData(member, user)
	return nil
}

func (p *ProjectMember) fillUserData(member *view_model.ProjectMemberView, user *usr_view_model.UserView) {
	member.UserName = user.UserName
	if user.HumanView != nil {
		member.FirstName = user.FirstName
		member.LastName = user.LastName
		member.Email = user.Email
		member.DisplayName = user.FirstName + " " + user.LastName
	}
	if user.MachineView != nil {
		member.DisplayName = user.MachineView.Name
	}
}
func (p *ProjectMember) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-u73es", "id", event.AggregateID).WithError(err).Warn("something went wrong in projectmember handler")
	return spooler.HandleError(event, err, p.view.GetLatestProjectMemberFailedEvent, p.view.ProcessedProjectMemberFailedEvent, p.view.ProcessedProjectMemberSequence, p.errorCountUntilSkip)
}

func (p *ProjectMember) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateProjectMemberSpoolerRunTimestamp)
}

func (u *ProjectMember) getUserByID(userID string) (*usr_view_model.UserView, error) {
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

func (u *ProjectMember) getUserEvents(userID string, sequence uint64) ([]*models.Event, error) {
	query, err := view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}

	return u.es.FilterEvents(context.Background(), query)
}
