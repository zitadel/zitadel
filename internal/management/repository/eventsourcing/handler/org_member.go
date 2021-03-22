package handler

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view"

	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	org_view_model "github.com/caos/zitadel/internal/org/repository/view/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	usr_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	orgMemberTable = "management.org_members"
)

type OrgMember struct {
	handler
	subscription *v1.Subscription
}

func newOrgMember(
	handler handler,
) *OrgMember {
	h := &OrgMember{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *OrgMember) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (m *OrgMember) ViewModel() string {
	return orgMemberTable
}

func (_ *OrgMember) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate, usr_es_model.UserAggregate}
}

func (p *OrgMember) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestOrgMemberSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *OrgMember) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := m.view.GetLatestOrgMemberSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *OrgMember) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate:
		err = m.processOrgMember(event)
	case usr_es_model.UserAggregate:
		err = m.processUser(event)
	}
	return err
}

func (m *OrgMember) processOrgMember(event *es_models.Event) (err error) {
	member := new(org_view_model.OrgMemberView)
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
		return m.view.DeleteOrgMember(event.AggregateID, member.UserID, event)
	default:
		return m.view.ProcessedOrgMemberSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutOrgMember(member, event)
}

func (m *OrgMember) processUser(event *es_models.Event) (err error) {
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
			return m.view.ProcessedOrgMemberSequence(event)
		}
		user, err := m.getUserByID(event.AggregateID)
		if err != nil {
			return err
		}
		for _, member := range members {
			m.fillUserData(member, user)
		}
		return m.view.PutOrgMembers(members, event)
	case usr_es_model.UserRemoved:
		return m.view.DeleteOrgMembersByUserID(event.AggregateID, event)
	default:
		return m.view.ProcessedOrgMemberSequence(event)
	}
}

func (m *OrgMember) fillData(member *org_view_model.OrgMemberView) (err error) {
	user, err := m.getUserByID(member.UserID)
	if err != nil {
		return err
	}
	return m.fillUserData(member, user)
}

func (m *OrgMember) fillUserData(member *org_view_model.OrgMemberView, user *usr_view_model.UserView) error {
	org, err := m.getOrgByID(context.Background(), user.ResourceOwner)
	policy := org.OrgIamPolicy
	if policy == nil {
		policy, err = m.getDefaultOrgIAMPolicy(context.TODO())
		if err != nil {
			return err
		}
	}
	member.UserName = user.UserName
	member.PreferredLoginName = user.GenerateLoginName(org.GetPrimaryDomain().Domain, policy.UserLoginMustBeDomain)
	if user.HumanView != nil {
		member.FirstName = user.FirstName
		member.LastName = user.LastName
		member.DisplayName = user.FirstName + " " + user.LastName
		member.Email = user.Email
	}
	if user.MachineView != nil {
		member.DisplayName = user.MachineView.Name
	}
	return nil
}

func (m *OrgMember) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-u73es", "id", event.AggregateID).WithError(err).Warn("something went wrong in orgmember handler")
	return spooler.HandleError(event, err, m.view.GetLatestOrgMemberFailedEvent, m.view.ProcessedOrgMemberFailedEvent, m.view.ProcessedOrgMemberSequence, m.errorCountUntilSkip)
}

func (o *OrgMember) OnSuccess() error {
	return spooler.HandleSuccess(o.view.UpdateOrgMemberSpoolerRunTimestamp)
}

func (u *OrgMember) getUserByID(userID string) (*usr_view_model.UserView, error) {
	user, usrErr := u.view.UserByID(userID)
	if usrErr != nil && !caos_errs.IsNotFound(usrErr) {
		return nil, usrErr
	}
	if user == nil {
		user = &usr_view_model.UserView{}
	}
	events, err := u.getUserEvents(userID, user.Sequence)
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
		return nil, caos_errs.ThrowNotFound(nil, "HANDLER-m9dos", "Errors.User.NotFound")
	}
	return &userCopy, nil
}

func (u *OrgMember) getUserEvents(userID string, sequence uint64) ([]*es_models.Event, error) {
	query, err := view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}

	return u.es.FilterEvents(context.Background(), query)
}

func (u *OrgMember) getOrgByID(ctx context.Context, orgID string) (*org_model.Org, error) {
	query, err := org_view.OrgByIDQuery(orgID, 0)
	if err != nil {
		return nil, err
	}

	esOrg := &model.Org{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: orgID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, esOrg.AppendEvents, query)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-kVLb2", "Errors.Org.NotFound")
	}

	return model.OrgToModel(esOrg), nil
}

func (u *OrgMember) getDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicy, error) {
	existingIAM, err := u.getIAMByID(ctx)
	if err != nil {
		return nil, err
	}
	if existingIAM.DefaultOrgIAMPolicy == nil {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-2Fj8s", "Errors.IAM.OrgIAMPolicy.NotExisting")
	}
	return existingIAM.DefaultOrgIAMPolicy, nil
}

func (u *OrgMember) getIAMByID(ctx context.Context) (*iam_model.IAM, error) {
	query, err := iam_view.IAMByIDQuery(domain.IAMID, 0)
	if err != nil {
		return nil, err
	}
	iam := &iam_es_model.IAM{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: domain.IAMID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, iam.AppendEvents, query)
	if err != nil && caos_errs.IsNotFound(err) && iam.Sequence == 0 {
		return nil, err
	}
	return iam_es_model.IAMToModel(iam), nil
}
