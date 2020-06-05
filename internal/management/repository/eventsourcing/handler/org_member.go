package handler

import (
	"context"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	view_model "github.com/caos/zitadel/internal/org/repository/view"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

type OrgMember struct {
	handler
	userEvents *usr_event.UserEventstore
}

const (
	orgMemberTable = "management.org_members"
)

func (m *OrgMember) MinimumCycleDuration() time.Duration { return m.cycleDuration }

func (m *OrgMember) ViewModel() string {
	return orgMemberTable
}

func (m *OrgMember) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestOrgMemberSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgAggregate, usr_es_model.UserAggregate).
		LatestSequenceFilter(sequence), nil
}

func (m *OrgMember) Process(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate:
		err = m.processOrgMember(event)
	case usr_es_model.UserAggregate:
		err = m.processUser(event)
	}
	return err
}

func (m *OrgMember) processOrgMember(event *models.Event) (err error) {
	member := new(view_model.OrgMemberView)
	switch event.Type {
	case model.OrgMemberAdded:
		member.AppendEvent(event)
		m.fillData(member)
	case model.OrgMemberChanged:
		err := member.SetData(event)
		if err != nil {
			return err
		}
		member, err = m.view.OrgMemberByIDs(event.AggregateID, member.UserID)
		if err != nil {
			return err
		}
		member.AppendEvent(event)
	case model.OrgMemberRemoved:
		err := member.SetData(event)
		if err != nil {
			return err
		}
		return m.view.DeleteOrgMember(event.AggregateID, member.UserID, event.Sequence)
	default:
		return m.view.ProcessedOrgMemberSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutOrgMember(member, member.Sequence)
}

func (m *OrgMember) processUser(event *models.Event) (err error) {
	switch event.Type {
	case usr_es_model.UserProfileChanged,
		usr_es_model.UserEmailChanged:
		members, err := m.view.OrgMembersByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		user, err := m.userEvents.UserByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		for _, member := range members {
			m.fillUserData(member, user)
			err = m.view.PutOrgMember(member, event.Sequence)
			if err != nil {
				return err
			}
		}
	default:
		return m.view.ProcessedOrgMemberSequence(event.Sequence)
	}
	return nil
}

func (m *OrgMember) fillData(member *view_model.OrgMemberView) (err error) {
	user, err := m.userEvents.UserByID(context.Background(), member.UserID)
	if err != nil {
		return err
	}
	m.fillUserData(member, user)
	return nil
}

func (m *OrgMember) fillUserData(member *view_model.OrgMemberView, user *usr_model.User) {
	member.UserName = user.UserName
	member.FirstName = user.FirstName
	member.LastName = user.LastName
	member.Email = user.EmailAddress
}
func (m *OrgMember) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-u73es", "id", event.AggregateID).WithError(err).Warn("something went wrong in orgmember handler")
	return spooler.HandleError(event, err, m.view.GetLatestOrgMemberFailedEvent, m.view.ProcessedOrgMemberFailedEvent, m.view.ProcessedOrgMemberSequence, m.errorCountUntilSkip)
}
