package handler

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	org_model "github.com/caos/zitadel/internal/org/repository/view/model"
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
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *OrgMember) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate:
		err = m.processOrgMember(event)
	case usr_es_model.UserAggregate:
		err = m.processUser(event)
	}
	return err
}

func (m *OrgMember) processOrgMember(event *models.Event) (err error) {
	member := new(org_model.OrgMemberView)
	switch event.Type {
	case model.OrgMemberAdded:
		err = member.AppendEvent(event)
		if err != nil {
			return err
		}
		err = m.fillData(member)
	case model.OrgMemberChanged:
		err = member.SetData(event)
		if err != nil {
			return err
		}
		member, err = m.view.OrgMemberByIDs(event.AggregateID, member.UserID)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case model.OrgMemberRemoved:
		err = member.SetData(event)
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
		usr_es_model.UserEmailChanged,
		usr_es_model.HumanProfileChanged,
		usr_es_model.HumanEmailChanged,
		usr_es_model.MachineChanged:
		members, err := m.view.OrgMembersByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(members) == 0 {
			return m.view.ProcessedOrgMemberSequence(event.Sequence)
		}
		user, err := m.userEvents.UserByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		for _, member := range members {
			m.fillUserData(member, user)
		}
		return m.view.PutOrgMembers(members, event.Sequence)
	case usr_es_model.UserRemoved:
		return m.view.DeleteOrgMembersByUserID(event.AggregateID, event.Sequence)
	default:
		return m.view.ProcessedOrgMemberSequence(event.Sequence)
	}
	return nil
}

func (m *OrgMember) fillData(member *org_model.OrgMemberView) (err error) {
	user, err := m.userEvents.UserByID(context.Background(), member.UserID)
	if err != nil {
		return err
	}
	m.fillUserData(member, user)
	return nil
}

func (m *OrgMember) fillUserData(member *org_model.OrgMemberView, user *usr_model.User) {
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
func (m *OrgMember) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-u73es", "id", event.AggregateID).WithError(err).Warn("something went wrong in orgmember handler")
	return spooler.HandleError(event, err, m.view.GetLatestOrgMemberFailedEvent, m.view.ProcessedOrgMemberFailedEvent, m.view.ProcessedOrgMemberSequence, m.errorCountUntilSkip)
}

func (o *OrgMember) OnSuccess() error {
	return spooler.HandleSuccess(o.view.UpdateOrgMemberSpoolerRunTimestamp)
}
