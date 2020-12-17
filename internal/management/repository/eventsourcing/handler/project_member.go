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

type ProjectMember struct {
	handler
	userEvents *usr_event.UserEventstore
}

const (
	projectMemberTable = "management.project_members"
)

func (p *ProjectMember) ViewModel() string {
	return projectMemberTable
}

func (_ *ProjectMember) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{proj_es_model.ProjectAggregate, usr_es_model.UserAggregate}
}

func (p *ProjectMember) CurrentSequence(event *models.Event) (uint64, error) {
	sequence, err := p.view.GetLatestProjectMemberSequence(string(event.AggregateType))
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *ProjectMember) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectMemberSequence("")
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
		user, err := p.userEvents.UserByID(context.Background(), event.AggregateID)
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
	user, err := p.userEvents.UserByID(context.Background(), member.UserID)
	if err != nil {
		return err
	}
	p.fillUserData(member, user)
	return nil
}

func (p *ProjectMember) fillUserData(member *view_model.ProjectMemberView, user *usr_model.User) {
	member.UserName = user.UserName
	if user.Human != nil {
		member.FirstName = user.FirstName
		member.LastName = user.LastName
		member.Email = user.EmailAddress
		member.DisplayName = user.FirstName + " " + user.LastName
	}
	if user.Machine != nil {
		member.DisplayName = user.Machine.Name
	}
}
func (p *ProjectMember) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-u73es", "id", event.AggregateID).WithError(err).Warn("something went wrong in projectmember handler")
	return spooler.HandleError(event, err, p.view.GetLatestProjectMemberFailedEvent, p.view.ProcessedProjectMemberFailedEvent, p.view.ProcessedProjectMemberSequence, p.errorCountUntilSkip)
}

func (p *ProjectMember) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateProjectMemberSpoolerRunTimestamp)
}
