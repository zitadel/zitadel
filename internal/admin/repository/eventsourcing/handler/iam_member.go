package handler

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

type IamMember struct {
	handler
	userEvents *usr_event.UserEventstore
}

const (
	iamMemberTable = "adminapi.iam_members"
)

func (m *IamMember) ViewModel() string {
	return iamMemberTable
}

func (m *IamMember) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestIAMMemberSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.IAMAggregate, usr_es_model.UserAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *IamMember) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate:
		err = m.processIamMember(event)
	case usr_es_model.UserAggregate:
		err = m.processUser(event)
	}
	return err
}

func (m *IamMember) processIamMember(event *models.Event) (err error) {
	member := new(iam_model.IAMMemberView)
	switch event.Type {
	case model.IAMMemberAdded:
		err = member.AppendEvent(event)
		if err != nil {
			return err
		}
		err = m.fillData(member)
	case model.IAMMemberChanged:
		err := member.SetData(event)
		if err != nil {
			return err
		}
		member, err = m.view.IAMMemberByIDs(event.AggregateID, member.UserID)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case model.IAMMemberRemoved:
		err := member.SetData(event)
		if err != nil {
			return err
		}
		return m.view.DeleteIAMMember(event.AggregateID, member.UserID, event.Sequence, event.CreationDate)
	default:
		return m.view.ProcessedIAMMemberSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return m.view.PutIAMMember(member, member.Sequence, event.CreationDate)
}

func (m *IamMember) processUser(event *models.Event) (err error) {
	switch event.Type {
	case usr_es_model.UserProfileChanged,
		usr_es_model.UserEmailChanged,
		usr_es_model.HumanProfileChanged,
		usr_es_model.HumanEmailChanged,
		usr_es_model.MachineChanged:
		members, err := m.view.IAMMembersByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(members) == 0 {
			return m.view.ProcessedIAMMemberSequence(event.Sequence, event.CreationDate)
		}
		user, err := m.userEvents.UserByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		for _, member := range members {
			m.fillUserData(member, user)
		}
		return m.view.PutIAMMembers(members, event.Sequence, event.CreationDate)
	case usr_es_model.UserRemoved:
		return m.view.DeleteIAMMembersByUserID(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return m.view.ProcessedIAMMemberSequence(event.Sequence, event.CreationDate)
	}
}

func (m *IamMember) fillData(member *iam_model.IAMMemberView) (err error) {
	user, err := m.userEvents.UserByID(context.Background(), member.UserID)
	if err != nil {
		return err
	}
	m.fillUserData(member, user)
	return nil
}

func (m *IamMember) fillUserData(member *iam_model.IAMMemberView, user *usr_model.User) {
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
func (m *IamMember) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Ld9ow", "id", event.AggregateID).WithError(err).Warn("something went wrong in iammember handler")
	return spooler.HandleError(event, err, m.view.GetLatestIAMMemberFailedEvent, m.view.ProcessedIAMMemberFailedEvent, m.view.ProcessedIAMMemberSequence, m.errorCountUntilSkip)
}

func (m *IamMember) OnSuccess() error {
	return spooler.HandleSuccess(m.view.UpdateIAMMemberSpoolerRunTimestamp)
}
