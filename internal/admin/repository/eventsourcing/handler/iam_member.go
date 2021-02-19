package handler

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/user/repository/view"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	iamMemberTable = "adminapi.iam_members"
)

type IAMMember struct {
	handler
	subscription *eventstore.Subscription
}

func newIAMMember(handler handler) *IAMMember {
	iamMember := &IAMMember{
		handler: handler,
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

func (m *IAMMember) CurrentSequence() (uint64, error) {
	sequence, err := m.view.GetLatestIAMMemberSequence()
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
	sequence, err := m.view.GetLatestIAMMemberSequence()
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
		err = m.processIAMMember(event)
	case usr_es_model.UserAggregate:
		err = m.processUser(event)
	}
	return err
}

func (m *IAMMember) processIAMMember(event *es_models.Event) (err error) {
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
		user, err := m.getUserByID(event.AggregateID)
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
	user, err := m.getUserByID(member.UserID)
	if err != nil {
		return err
	}
	m.fillUserData(member, user)
	return nil
}

func (m *IAMMember) fillUserData(member *iam_model.IAMMemberView, user *view_model.UserView) {
	member.UserName = user.UserName
	if user.HumanView != nil {
		member.FirstName = user.FirstName
		member.LastName = user.LastName
		member.DisplayName = user.FirstName + " " + user.LastName
		member.Email = user.Email
	}
	if user.MachineView != nil {
		member.DisplayName = user.MachineView.Name
	}
}
func (m *IAMMember) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-Ld9ow", "id", event.AggregateID).WithError(err).Warn("something went wrong in iammember handler")
	return spooler.HandleError(event, err, m.view.GetLatestIAMMemberFailedEvent, m.view.ProcessedIAMMemberFailedEvent, m.view.ProcessedIAMMemberSequence, m.errorCountUntilSkip)
}

func (m *IAMMember) OnSuccess() error {
	return spooler.HandleSuccess(m.view.UpdateIAMMemberSpoolerRunTimestamp)
}

func (m *IAMMember) getUserByID(userID string) (*view_model.UserView, error) {
	user, usrErr := m.view.UserByID(userID)
	if usrErr != nil && !caos_errs.IsNotFound(usrErr) {
		return nil, usrErr
	}
	if user == nil {
		user = &view_model.UserView{}
	}
	events, err := m.getUserEvents(userID, user.Sequence)
	if err != nil {
		return user, usrErr
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return user, nil
		}
	}
	if userCopy.State == int32(usr_model.UserStateDeleted) {
		return nil, caos_errs.ThrowNotFound(nil, "HANDLER-4n9fs", "Errors.User.NotFound")
	}
	return &userCopy, nil
}

func (m *IAMMember) getUserEvents(userID string, sequence uint64) ([]*models.Event, error) {
	query, err := view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}

	return m.es.FilterEvents(context.Background(), query)
}
