package handler

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
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

func (_ *UserMembership) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{iam_es_model.IAMAggregate, org_es_model.OrgAggregate, proj_es_model.ProjectAggregate, model.UserAggregate}
}

func (u *UserMembership) CurrentSequence(event *models.Event) (uint64, error) {
	sequence, err := u.view.GetLatestUserMembershipSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *UserMembership) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := m.view.GetLatestUserMembershipSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *UserMembership) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case iam_es_model.IAMAggregate:
		err = m.processIam(event)
	case org_es_model.OrgAggregate:
		err = m.processOrg(event)
	case proj_es_model.ProjectAggregate:
		err = m.processProject(event)
	case model.UserAggregate:
		err = m.processUser(event)
	}
	return err
}

func (m *UserMembership) processIam(event *es_models.Event) (err error) {
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
		return m.view.DeleteUserMembership(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeIam, event.Sequence, event.CreationDate)
	default:
		return m.view.ProcessedUserMembershipSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return m.view.PutUserMembership(member, event.Sequence, event.CreationDate)
}

func (m *UserMembership) fillIamDisplayName(member *usr_es_model.UserMembershipView) {
	member.DisplayName = member.AggregateID
}

func (m *UserMembership) processOrg(event *es_models.Event) (err error) {
	member := new(usr_es_model.UserMembershipView)
	err = member.AppendEvent(event)
	if err != nil {
		return err
	}
	switch event.Type {
	case org_es_model.OrgMemberAdded:
		err = m.fillOrgDisplayName(member)
	case org_es_model.OrgMemberChanged:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeOrganisation)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case org_es_model.OrgMemberRemoved:
		return m.view.DeleteUserMembership(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeOrganisation, event.Sequence, event.CreationDate)
	case org_es_model.OrgChanged:
		return m.updateOrgDisplayName(event)
	default:
		return m.view.ProcessedUserMembershipSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return m.view.PutUserMembership(member, event.Sequence, event.CreationDate)
}

func (m *UserMembership) fillOrgDisplayName(member *usr_es_model.UserMembershipView) (err error) {
	org, err := m.orgEvents.OrgByID(context.Background(), org_model.NewOrg(member.AggregateID))
	if err != nil {
		return err
	}
	member.DisplayName = org.Name
	return nil
}

func (m *UserMembership) updateOrgDisplayName(event *es_models.Event) error {
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
	return m.view.BulkPutUserMemberships(memberships, event.Sequence, event.CreationDate)
}

func (m *UserMembership) processProject(event *es_models.Event) (err error) {
	member := new(usr_es_model.UserMembershipView)
	err = member.AppendEvent(event)
	if err != nil {
		return err
	}
	switch event.Type {
	case proj_es_model.ProjectMemberAdded, proj_es_model.ProjectGrantMemberAdded:
		err = m.fillProjectDisplayName(member)
	case proj_es_model.ProjectMemberChanged:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeProject)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case proj_es_model.ProjectMemberRemoved:
		return m.view.DeleteUserMembership(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeProject, event.Sequence, event.CreationDate)
	case proj_es_model.ProjectGrantMemberChanged:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, member.ObjectID, usr_model.MemberTypeProjectGrant)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case proj_es_model.ProjectGrantMemberRemoved:
		return m.view.DeleteUserMembership(member.UserID, event.AggregateID, member.ObjectID, usr_model.MemberTypeProjectGrant, event.Sequence, event.CreationDate)
	case proj_es_model.ProjectChanged:
		return m.updateProjectDisplayName(event)
	case proj_es_model.ProjectRemoved:
		return m.view.DeleteUserMembershipsByAggregateID(event.AggregateID, event.Sequence, event.CreationDate)
	case proj_es_model.ProjectGrantRemoved:
		return m.view.DeleteUserMembershipsByAggregateIDAndObjectID(event.AggregateID, member.ObjectID, event.Sequence, event.CreationDate)
	default:
		return m.view.ProcessedUserMembershipSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return m.view.PutUserMembership(member, event.Sequence, event.CreationDate)
}

func (m *UserMembership) fillProjectDisplayName(member *usr_es_model.UserMembershipView) (err error) {
	project, err := m.projectEvents.ProjectByID(context.Background(), member.AggregateID)
	if err != nil {
		return err
	}
	member.DisplayName = project.Name
	return nil
}

func (m *UserMembership) updateProjectDisplayName(event *es_models.Event) error {
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
	return m.view.BulkPutUserMemberships(memberships, event.Sequence, event.CreationDate)
}

func (m *UserMembership) processUser(event *es_models.Event) (err error) {
	switch event.Type {
	case model.UserRemoved:
		return m.view.DeleteUserMembershipsByUserID(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return m.view.ProcessedUserMembershipSequence(event.Sequence, event.CreationDate)
	}
}

func (m *UserMembership) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-Fwer2", "id", event.AggregateID).WithError(err).Warn("something went wrong in user membership handler")
	return spooler.HandleError(event, err, m.view.GetLatestUserMembershipFailedEvent, m.view.ProcessedUserMembershipFailedEvent, m.view.ProcessedUserMembershipSequence, m.errorCountUntilSkip)
}

func (m *UserMembership) OnSuccess() error {
	return spooler.HandleSuccess(m.view.UpdateUserMembershipSpoolerRunTimestamp)
}
