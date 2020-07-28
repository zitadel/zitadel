package handler

import (
	"context"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type UserMembership struct {
	handler
	orgEvents     *org_event.OrgEventstore
	projectEvents *proj_event.ProjectEventstore
}

const (
	userMembershipTable = "management.user_memberships"
)

func (m *UserMembership) ViewModel() string {
	return userMembershipTable
}

func (m *UserMembership) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestUserMembershipSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(org_es_model.OrgAggregate, proj_es_model.ProjectAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *UserMembership) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case org_es_model.OrgAggregate:
		err = m.processOrg(event)
	case proj_es_model.ProjectAggregate:
		err = m.processProject(event)
	}
	return err
}

func (m *UserMembership) processOrg(event *models.Event) (err error) {
	member := new(usr_es_model.UserMembershipView)
	member.AppendEvent(event)
	switch event.Type {
	case org_es_model.OrgMemberAdded:
		m.fillOrgDisplayName(member)
	case org_es_model.OrgMemberChanged:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeOrganisation)
		if err != nil {
			return err
		}
		member.AppendEvent(event)
	case org_es_model.OrgMemberRemoved:
		err = m.view.DeleteUserMembership(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeOrganisation, event.Sequence)
	case org_es_model.OrgChanged:
		m.updateOrgDisplayName(event)
	default:
		return m.view.ProcessedUserMembershipSequence(event.Sequence)
	}
	return m.view.PutUserMembership(member, event.Sequence)
}

func (m *UserMembership) fillOrgDisplayName(member *usr_es_model.UserMembershipView) (err error) {
	org, err := m.orgEvents.OrgByID(context.Background(), org_model.NewOrg(member.AggregateID))
	if err != nil {
		return err
	}
	member.DisplayName = org.Name
	return nil
}

func (m *UserMembership) updateOrgDisplayName(event *models.Event) error {
	org, err := m.orgEvents.OrgByID(context.Background(), org_model.NewOrg(event.AggregateID))
	if err != nil {
		return err
	}

	memberships, err := m.view.UserMembershipsByAggregateID(event.AggregateID)
	if err != nil {
		return err
	}
	for _, membership := range memberships {
		membership.DisplayName = org.Name
	}
	return m.view.BulkPutUserMemberships(memberships, event.Sequence)
}

func (m *UserMembership) processProject(event *models.Event) (err error) {
	member := new(usr_es_model.UserMembershipView)
	member.AppendEvent(event)
	switch event.Type {
	case proj_es_model.ProjectMemberAdded, proj_es_model.ProjectGrantMemberAdded:
		m.fillProjectDisplayName(member)
	case proj_es_model.ProjectMemberChanged:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeProject)
		if err != nil {
			return err
		}
		member.AppendEvent(event)
	case proj_es_model.ProjectMemberRemoved:
		err = m.view.DeleteUserMembership(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeProject, event.Sequence)
	case proj_es_model.ProjectGrantMemberChanged:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, member.ObjectID, usr_model.MemberTypeProjectGrant)
		if err != nil {
			return err
		}
		member.AppendEvent(event)
	case proj_es_model.ProjectGrantMemberRemoved:
		err = m.view.DeleteUserMembership(member.UserID, event.AggregateID, member.ObjectID, usr_model.MemberTypeProjectGrant, event.Sequence)
	case proj_es_model.ProjectChanged:
		m.updateProjectDisplayName(event)
	default:
		return m.view.ProcessedUserMembershipSequence(event.Sequence)
	}
	return m.view.PutUserMembership(member, event.Sequence)
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
	return m.view.BulkPutUserMemberships(memberships, event.Sequence)
}

func (m *UserMembership) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Ms3fj", "id", event.AggregateID).WithError(err).Warn("something went wrong in orgmember handler")
	return spooler.HandleError(event, err, m.view.GetLatestUserMembershipFailedEvent, m.view.ProcessedUserMembershipFailedEvent, m.view.ProcessedUserMembershipSequence, m.errorCountUntilSkip)
}
