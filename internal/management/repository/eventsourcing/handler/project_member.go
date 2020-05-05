package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
	"time"
)

type ProjectMember struct {
	handler
	//TODO: Add UserEvents
}

const (
	projectMemberTable = "management.project_members"
)

func (p *ProjectMember) MinimumCycleDuration() time.Duration { return p.cycleDuration }

func (p *ProjectMember) ViewModel() string {
	return projectMemberTable
}

func (p *ProjectMember) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectMemberSequence()
	if err != nil {
		return nil, err
	}
	//TODO: AddUserAggregate
	return es_models.NewSearchQuery().
		AggregateTypeFilter(es_model.ProjectAggregate).
		LatestSequenceFilter(sequence), nil
}

func (p *ProjectMember) Process(event *models.Event) (err error) {
	member := new(view_model.ProjectMemberView)
	switch event.Type {
	case es_model.ProjectMemberAdded:
		member.AppendEvent(event)
		//TODO: getUserData
	case es_model.ProjectMemberChanged:
		err := member.SetData(event)
		if err != nil {
			return err
		}
		member, err = p.view.ProjectMemberByIDs(event.AggregateID, member.UserID)
		if err != nil {
			return err
		}
		member.AppendEvent(event)
	case es_model.ProjectMemberRemoved:
		err := member.SetData(event)
		if err != nil {
			return err
		}
		return p.view.DeleteProjectMember(event.AggregateID, member.UserID, event.Sequence)
	default:
		return p.view.ProcessedProjectMemberSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return p.view.PutProjectMember(member)
}

func (p *ProjectMember) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-u73es", "id", event.AggregateID).WithError(err).Warn("something went wrong in projectmember handler")
	failedEvent, err := p.view.GetLatestProjectMemberFailedEvent(event.Sequence)
	if err != nil {
		return err
	}
	failedEvent.FailureCount++
	failedEvent.ErrMsg = err.Error()
	err = p.view.ProcessedProjectMemberFailedEvent(failedEvent)
	if err != nil {
		return err
	}
	if p.errorCountUntilSkip == failedEvent.FailureCount {
		return p.view.ProcessedProjectMemberSequence(event.Sequence)
	}
	return nil
}
