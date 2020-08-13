package handler

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

type ProjectGrantMember struct {
	handler
	userEvents *usr_event.UserEventstore
}

const (
	projectGrantMemberTable = "management.project_grant_members"
)

func (p *ProjectGrantMember) ViewModel() string {
	return projectGrantMemberTable
}

func (p *ProjectGrantMember) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectGrantMemberSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(proj_es_model.ProjectAggregate, usr_es_model.UserAggregate).
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
		p.fillData(member)
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
		return p.view.DeleteProjectGrantMember(member.GrantID, member.UserID, event.Sequence)
	case proj_es_model.ProjectRemoved:
		return p.view.DeleteProjectGrantMembersByProjectID(event.AggregateID)

	default:
		return p.view.ProcessedProjectGrantMemberSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return p.view.PutProjectGrantMember(member, member.Sequence)
}

func (p *ProjectGrantMember) processUser(event *models.Event) (err error) {
	switch event.Type {
	case usr_es_model.UserProfileChanged,
		usr_es_model.UserEmailChanged:
		members, err := p.view.ProjectGrantMembersByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(members) == 0 {
			return p.view.ProcessedProjectGrantMemberSequence(event.Sequence)
		}
		user, err := p.userEvents.UserByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		for _, member := range members {
			p.fillUserData(member, user)
			err = p.view.PutProjectGrantMember(member, event.Sequence)
			if err != nil {
				return err
			}
		}
	default:
		return p.view.ProcessedProjectGrantMemberSequence(event.Sequence)
	}
	return nil
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
	member.FirstName = user.FirstName
	member.LastName = user.LastName
	member.Email = user.EmailAddress
	member.DisplayName = user.DisplayName
}

func (p *ProjectGrantMember) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-kls93", "id", event.AggregateID).WithError(err).Warn("something went wrong in projectmember handler")
	return spooler.HandleError(event, err, p.view.GetLatestProjectGrantMemberFailedEvent, p.view.ProcessedProjectGrantMemberFailedEvent, p.view.ProcessedProjectGrantMemberSequence, p.errorCountUntilSkip)
}
