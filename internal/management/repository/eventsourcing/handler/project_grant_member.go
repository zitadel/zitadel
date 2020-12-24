package handler

import (
	"context"

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
)

const (
	projectGrantMemberTable = "management.project_grant_members"
)

type ProjectGrantMember struct {
	handler
	userEvents   *usr_event.UserEventstore
	subscription *eventstore.Subscription
}

func newProjectGrantMember(
	handler handler,
	userEvents *usr_event.UserEventstore,
) *ProjectGrantMember {
	h := &ProjectGrantMember{
		handler:    handler,
		userEvents: userEvents,
	}

	h.subscribe()

	return h
}

func (m *ProjectGrantMember) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (p *ProjectGrantMember) ViewModel() string {
	return projectGrantMemberTable
}

func (_ *ProjectGrantMember) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{proj_es_model.ProjectAggregate, usr_es_model.UserAggregate}
}

func (p *ProjectGrantMember) CurrentSequence(event *models.Event) (uint64, error) {
	sequence, err := p.view.GetLatestProjectGrantMemberSequence(string(event.AggregateType))
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *ProjectGrantMember) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectGrantMemberSequence("")
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *ProjectGrantMember) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case proj_es_model.ProjectAggregate:
		err = p.processProjectGrantMember(event)
	case usr_es_model.UserAggregate:
		err = p.processUser(event)
	}
	return err
}

func (p *ProjectGrantMember) processProjectGrantMember(event *models.Event) (err error) {
	member := new(view_model.ProjectGrantMemberView)
	switch event.Type {
	case proj_es_model.ProjectGrantMemberAdded:
		err = member.AppendEvent(event)
		if err != nil {
			return err
		}
		err = p.fillData(member)
	case proj_es_model.ProjectGrantMemberChanged:
		err = member.SetData(event)
		if err != nil {
			return err
		}
		member, err = p.view.ProjectGrantMemberByIDs(member.GrantID, member.UserID)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case proj_es_model.ProjectGrantMemberRemoved:
		err = member.SetData(event)
		if err != nil {
			return err
		}
		return p.view.DeleteProjectGrantMember(member.GrantID, member.UserID, event)
	case proj_es_model.ProjectRemoved:
		err = p.view.DeleteProjectGrantMembersByProjectID(event.AggregateID)
		if err != nil {
			return err
		}
		return p.view.ProcessedProjectGrantMemberSequence(event)
	default:
		return p.view.ProcessedProjectGrantMemberSequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutProjectGrantMember(member, event)
}

func (p *ProjectGrantMember) processUser(event *models.Event) (err error) {
	switch event.Type {
	case usr_es_model.UserProfileChanged,
		usr_es_model.UserEmailChanged,
		usr_es_model.HumanProfileChanged,
		usr_es_model.HumanEmailChanged,
		usr_es_model.MachineChanged:
		members, err := p.view.ProjectGrantMembersByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(members) == 0 {
			return p.view.ProcessedProjectGrantMemberSequence(event)
		}
		user, err := p.userEvents.UserByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		for _, member := range members {
			p.fillUserData(member, user)
		}
		return p.view.PutProjectGrantMembers(members, event)
	default:
		return p.view.ProcessedProjectGrantMemberSequence(event)
	}
}

func (p *ProjectGrantMember) fillData(member *view_model.ProjectGrantMemberView) (err error) {
	user, err := p.userEvents.UserByID(context.Background(), member.UserID)
	if err != nil {
		return err
	}
	p.fillUserData(member, user)
	return nil
}

func (p *ProjectGrantMember) fillUserData(member *view_model.ProjectGrantMemberView, user *usr_model.User) {
	member.UserName = user.UserName
	if user.Human != nil {
		member.FirstName = user.FirstName
		member.LastName = user.LastName
		member.DisplayName = user.FirstName + " " + user.LastName
		member.Email = user.EmailAddress
	}
	if user.Machine != nil {
		member.DisplayName = user.Machine.Name
	}
}

func (p *ProjectGrantMember) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-kls93", "id", event.AggregateID).WithError(err).Warn("something went wrong in projectmember handler")
	return spooler.HandleError(event, err, p.view.GetLatestProjectGrantMemberFailedEvent, p.view.ProcessedProjectGrantMemberFailedEvent, p.view.ProcessedProjectGrantMemberSequence, p.errorCountUntilSkip)
}

func (p *ProjectGrantMember) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateProjectGrantMemberSpoolerRunTimestamp)
}
