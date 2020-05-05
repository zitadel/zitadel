package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
	"time"
)

type ProjectGrantMember struct {
	handler
	//TODO: Add UserEvents
}

const (
	projectGrantMemberTable = "management.project_grant_members"
)

func (p *ProjectGrantMember) MinimumCycleDuration() time.Duration { return p.cycleDuration }

func (p *ProjectGrantMember) ViewModel() string {
	return projectGrantMemberTable
}

func (p *ProjectGrantMember) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectMemberSequence()
	if err != nil {
		return nil, err
	}
	//TODO: AddUserAggregate
	return es_models.NewSearchQuery().
		AggregateTypeFilter(es_model.ProjectAggregate).
		LatestSequenceFilter(sequence), nil
}

func (p *ProjectGrantMember) Process(event *models.Event) (err error) {
	member := new(view_model.ProjectGrantMemberView)
	switch event.Type {
	case es_model.ProjectGrantMemberAdded:
		member.AppendEvent(event)
		//TODO: getUserData
	case es_model.ProjectGrantMemberChanged:
		err := member.SetData(event)
		if err != nil {
			return err
		}
		member, err = p.view.ProjectGrantMemberByIDs(member.GrantID, member.UserID)
		if err != nil {
			return err
		}
		member.AppendEvent(event)
	case es_model.ProjectGrantMemberRemoved:
		err := member.SetData(event)
		if err != nil {
			return err
		}
		return p.view.DeleteProjectGrantMember(event.AggregateID, member.UserID, event.Sequence)
	default:
		return p.view.ProcessedProjectGrantMemberSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return p.view.PutProjectGrantMember(member)
}

func (p *ProjectGrantMember) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-kls93", "id", event.AggregateID).WithError(err).Warn("something went wrong in projectmember handler")
	failedEvent, err := p.view.GetLatestProjectGrantMemberFailedEvent(event.Sequence)
	if err != nil {
		return err
	}
	failedEvent.FailureCount++
	failedEvent.ErrMsg = err.Error()
	err = p.view.ProcessedProjectGrantMemberFailedEvent(failedEvent)
	if err != nil {
		return err
	}
	if p.errorCountUntilSkip == failedEvent.FailureCount {
		return p.view.ProcessedProjectGrantMemberSequence(event.Sequence)
	}
	return nil
}
