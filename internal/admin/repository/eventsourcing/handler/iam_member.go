package handler

import (
	"context"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
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

func (m *IamMember) MinimumCycleDuration() time.Duration { return m.cycleDuration }

func (m *IamMember) ViewModel() string {
	return iamMemberTable
}

func (m *IamMember) EventQuery() (*models.SearchQuery, error) {
	sequence, _, err := m.view.GetLatestIamMemberSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.IamAggregate, usr_es_model.UserAggregate).
		LatestSequenceFilter(sequence), nil
}

func (m *IamMember) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.IamAggregate:
		err = m.processIamMember(event)
	case usr_es_model.UserAggregate:
		err = m.processUser(event)
	}
	return err
}

func (m *IamMember) processIamMember(event *models.Event) (err error) {
	member := new(iam_model.IamMemberView)
	switch event.Type {
	case model.IamMemberAdded:
		member.AppendEvent(event)
		m.fillData(member)
	case model.IamMemberChanged:
		err := member.SetData(event)
		if err != nil {
			return err
		}
		member, err = m.view.IamMemberByIDs(event.AggregateID, member.UserID)
		if err != nil {
			return err
		}
		member.AppendEvent(event)
	case model.IamMemberRemoved:
		err := member.SetData(event)
		if err != nil {
			return err
		}
		return m.view.DeleteIamMember(event.AggregateID, member.UserID, event.Sequence)
	default:
		return m.view.ProcessedIamMemberSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutIamMember(member, member.Sequence)
}

func (m *IamMember) processUser(event *models.Event) (err error) {
	switch event.Type {
	case usr_es_model.UserProfileChanged,
		usr_es_model.UserEmailChanged:
		members, err := m.view.IamMembersByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		user, err := m.userEvents.UserByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		for _, member := range members {
			m.fillUserData(member, user)
			err = m.view.PutIamMember(member, event.Sequence)
			if err != nil {
				return err
			}
		}
	default:
		return m.view.ProcessedIamMemberSequence(event.Sequence)
	}
	return nil
}

func (m *IamMember) fillData(member *iam_model.IamMemberView) (err error) {
	user, err := m.userEvents.UserByID(context.Background(), member.UserID)
	if err != nil {
		return err
	}
	m.fillUserData(member, user)
	return nil
}

func (m *IamMember) fillUserData(member *iam_model.IamMemberView, user *usr_model.User) {
	member.UserName = user.UserName
	member.FirstName = user.FirstName
	member.LastName = user.LastName
	member.Email = user.EmailAddress
}
func (m *IamMember) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Ld9ow", "id", event.AggregateID).WithError(err).Warn("something went wrong in iammember handler")
	return spooler.HandleError(event, err, m.view.GetLatestIamMemberFailedEvent, m.view.ProcessedIamMemberFailedEvent, m.view.ProcessedIamMemberSequence, m.errorCountUntilSkip)
}
