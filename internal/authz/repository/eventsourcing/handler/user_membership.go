package handler

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	userMembershipTable = "authz.user_memberships"
)

type UserMembership struct {
	handler
	orgEvents     *org_event.OrgEventstore
	projectEvents *proj_event.ProjectEventstore
	subscription  *eventstore.Subscription
}

func newUserMembership(
	handler handler,
	orgEvents *org_event.OrgEventstore,
	projectEvents *proj_event.ProjectEventstore,
) *UserMembership {
	h := &UserMembership{
		handler:       handler,
		orgEvents:     orgEvents,
		projectEvents: projectEvents,
	}

	h.subscribe()

	return h
}

func (m *UserMembership) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (m *UserMembership) ViewModel() string {
	return userMembershipTable
}

func (_ *UserMembership) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{iam_es_model.IAMAggregate, org_es_model.OrgAggregate, proj_es_model.ProjectAggregate, model.UserAggregate}
}

func (m *UserMembership) CurrentSequence() (uint64, error) {
	sequence, err := m.view.GetLatestUserMembershipSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *UserMembership) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestUserMembershipSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *UserMembership) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case iam_es_model.IAMAggregate:
		err = m.processIAM(event)
	case org_es_model.OrgAggregate:
		err = m.processOrg(event)
	case proj_es_model.ProjectAggregate:
		err = m.processProject(event)
	case model.UserAggregate:
		err = m.processUser(event)
	}
	return err
}

func (m *UserMembership) processIAM(event *models.Event) (err error) {
	member := new(usr_es_model.UserMembershipView)
	err = member.AppendEvent(event)
	if err != nil {
		return err
	}
	switch event.Type {
	case iam_es_model.IAMMemberAdded:
		m.fillIamDisplayName(member)
	case iam_es_model.IAMMemberChanged:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeIam)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case iam_es_model.IAMMemberRemoved:
		return m.view.DeleteUserMembership(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeIam, event)
	default:
		return m.view.ProcessedUserMembershipSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutUserMembership(member, event)
}

func (m *UserMembership) fillIamDisplayName(member *usr_es_model.UserMembershipView) {
	member.DisplayName = member.AggregateID
	member.ResourceOwnerName = member.ResourceOwner
}

func (m *UserMembership) processOrg(event *models.Event) (err error) {
	member := new(usr_es_model.UserMembershipView)
	err = member.AppendEvent(event)
	if err != nil {
		return err
	}
	switch event.Type {
	case org_es_model.OrgMemberAdded:
		err = m.fillOrgName(member)
	case org_es_model.OrgMemberChanged:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeOrganisation)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case org_es_model.OrgMemberRemoved:
		return m.view.DeleteUserMembership(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeOrganisation, event)
	case org_es_model.OrgChanged:
		return m.updateOrgName(event)
	default:
		return m.view.ProcessedUserMembershipSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutUserMembership(member, event)
}

func (m *UserMembership) fillOrgName(member *usr_es_model.UserMembershipView) (err error) {
	org, err := m.orgEvents.OrgByID(context.Background(), org_model.NewOrg(member.ResourceOwner))
	if err != nil {
		return err
	}
	member.ResourceOwnerName = org.Name
	if member.AggregateID == org.AggregateID {
		member.DisplayName = org.Name
	}
	return nil
}

func (m *UserMembership) updateOrgName(event *models.Event) error {
	org, err := m.orgEvents.OrgByID(context.Background(), org_model.NewOrg(event.AggregateID))
	if err != nil {
		return err
	}

	memberships, err := m.view.UserMembershipsByResourceOwner(event.ResourceOwner)
	if err != nil {
		return err
	}
	for _, membership := range memberships {
		membership.ResourceOwnerName = org.Name
		if membership.AggregateID == event.AggregateID {
			membership.DisplayName = org.Name
		}
	}
	return m.view.BulkPutUserMemberships(memberships, event)
}

func (m *UserMembership) processProject(event *models.Event) (err error) {
	member := new(usr_es_model.UserMembershipView)
	err = member.AppendEvent(event)
	if err != nil {
		return err
	}
	switch event.Type {
	case proj_es_model.ProjectMemberAdded, proj_es_model.ProjectGrantMemberAdded:
		err = m.fillProjectDisplayName(member)
		if err != nil {
			return err
		}
		err = m.fillOrgName(member)
	case proj_es_model.ProjectMemberChanged:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeProject)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case proj_es_model.ProjectMemberRemoved:
		return m.view.DeleteUserMembership(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeProject, event)
	case proj_es_model.ProjectGrantMemberChanged:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, member.ObjectID, usr_model.MemberTypeProjectGrant)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case proj_es_model.ProjectGrantMemberRemoved:
		return m.view.DeleteUserMembership(member.UserID, event.AggregateID, member.ObjectID, usr_model.MemberTypeProjectGrant, event)
	case proj_es_model.ProjectChanged:
		return m.updateProjectDisplayName(event)
	case proj_es_model.ProjectRemoved:
		return m.view.DeleteUserMembershipsByAggregateID(event.AggregateID, event)
	case proj_es_model.ProjectGrantRemoved:
		return m.view.DeleteUserMembershipsByAggregateIDAndObjectID(event.AggregateID, member.ObjectID, event)
	default:
		return m.view.ProcessedUserMembershipSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutUserMembership(member, event)
}

func (m *UserMembership) fillProjectDisplayName(member *usr_es_model.UserMembershipView) (err error) {
	project, err := m.projectEvents.ProjectByID(context.Background(), member.AggregateID)
	if err != nil {
		return err
	}
	member.DisplayName = project.Name
	return nil
}

func (m *UserMembership) updateProjectDisplayName(event *models.Event) error {
	project, err := m.projectEvents.ProjectByID(context.Background(), event.AggregateID)
	if err != nil {
		return err
	}

	memberships, err := m.view.UserMembershipsByAggregateID(event.AggregateID)
	if err != nil {
		return err
	}
	for _, membership := range memberships {
		membership.DisplayName = project.Name
	}
	return m.view.BulkPutUserMemberships(memberships, event)
}

func (m *UserMembership) processUser(event *models.Event) (err error) {
	switch event.Type {
	case model.UserRemoved:
		return m.view.DeleteUserMembershipsByUserID(event.AggregateID, event)
	default:
		return m.view.ProcessedUserMembershipSequence(event)
	}
}

func (m *UserMembership) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Ms3fj", "id", event.AggregateID).WithError(err).Warn("something went wrong in user membership handler")
	return spooler.HandleError(event, err, m.view.GetLatestUserMembershipFailedEvent, m.view.ProcessedUserMembershipFailedEvent, m.view.ProcessedUserMembershipSequence, m.errorCountUntilSkip)
}

func (m *UserMembership) OnSuccess() error {
	return spooler.HandleSuccess(m.view.UpdateUserMembershipSpoolerRunTimestamp)
}
