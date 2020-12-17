package handler

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	iamMemberTable = "adminapi.iam_members"
)

type IAMMember struct {
	handler
	userEvents   *usr_event.UserEventstore
	subscription *eventstore.Subscription
}

func newIAMMember(handler handler, userEvents *usr_event.UserEventstore) *IAMMember {
	iamMember := &IAMMember{
		handler:    handler,
		userEvents: userEvents,
	}

	iamMember.subscribe()

	return iamMember
}

func (m *IAMMember) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (m *IAMMember) CurrentSequence(event *es_models.Event) (uint64, error) {
	sequence, err := m.view.GetLatestIAMMemberSequence(string(event.AggregateType))
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *IAMMember) ViewModel() string {
	return iamMemberTable
}

func (m *IAMMember) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.IAMAggregate, usr_es_model.UserAggregate}
}

func (m *IAMMember) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := m.view.GetLatestIAMMemberSequence("")
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *IAMMember) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate:
		err = m.processIamMember(event)
	case usr_es_model.UserAggregate:
		err = m.processUser(event)
	}
	return err
}

func (m *IAMMember) processIamMember(event *es_models.Event) (err error) {
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
		return m.view.DeleteIAMMember(event.AggregateID, member.UserID, event)
	default:
		return m.view.ProcessedIAMMemberSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutIAMMember(member, event)
}

func (m *IAMMember) processUser(event *es_models.Event) (err error) {
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
			return m.view.ProcessedIAMMemberSequence(event)
		}
		user, err := m.userEvents.UserByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		for _, member := range members {
			m.fillUserData(member, user)
		}
		return m.view.PutIAMMembers(members, event)
	case usr_es_model.UserRemoved:
		return m.view.DeleteIAMMembersByUserID(event.AggregateID, event)
	default:
		return m.view.ProcessedIAMMemberSequence(event)
	}
}

func (m *IAMMember) fillData(member *iam_model.IAMMemberView) (err error) {
	user, err := m.userEvents.UserByID(context.Background(), member.UserID)
	if err != nil {
		return err
	}
	m.fillUserData(member, user)
	return nil
}

func (m *IAMMember) fillUserData(member *iam_model.IAMMemberView, user *usr_model.User) {
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
func (m *IAMMember) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-Ld9ow", "id", event.AggregateID).WithError(err).Warn("something went wrong in iammember handler")
	return spooler.HandleError(event, err, m.view.GetLatestIAMMemberFailedEvent, m.view.ProcessedIAMMemberFailedEvent, m.view.ProcessedIAMMemberSequence, m.errorCountUntilSkip)
}

func (m *IAMMember) OnSuccess() error {
	return spooler.HandleSuccess(m.view.UpdateIAMMemberSpoolerRunTimestamp)
}
