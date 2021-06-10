package handler

import (
	"context"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	proj_view "github.com/caos/zitadel/internal/project/repository/view"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	userMembershipTable = "management.user_memberships"
)

type UserMembership struct {
	handler
	subscription *v1.Subscription
}

func newUserMembership(
	handler handler,
) *UserMembership {
	h := &UserMembership{
		handler: handler,
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

func (u *UserMembership) CurrentSequence() (uint64, error) {
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

func (m *UserMembership) processIAM(event *es_models.Event) (err error) {
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
	case iam_es_model.IAMMemberRemoved,
		iam_es_model.IAMMemberCascadeRemoved:
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
	case org_es_model.OrgMemberRemoved, org_es_model.OrgMemberCascadeRemoved:
		return m.view.DeleteUserMembership(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeOrganisation, event)
	case org_es_model.OrgChanged:
		return m.updateOrgDisplayName(event)
	default:
		return m.view.ProcessedUserMembershipSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutUserMembership(member, event)
}

func (m *UserMembership) fillOrgDisplayName(member *usr_es_model.UserMembershipView) (err error) {
	org, err := m.getOrgByID(context.Background(), member.AggregateID)
	if err != nil {
		return err
	}
	member.DisplayName = org.Name
	return nil
}

func (m *UserMembership) updateOrgDisplayName(event *es_models.Event) error {
	org, err := m.getOrgByID(context.Background(), event.AggregateID)
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
	return m.view.BulkPutUserMemberships(memberships, event)
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
	case proj_es_model.ProjectMemberRemoved, proj_es_model.ProjectMemberCascadeRemoved:
		return m.view.DeleteUserMembership(member.UserID, event.AggregateID, event.AggregateID, usr_model.MemberTypeProject, event)
	case proj_es_model.ProjectGrantMemberChanged:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, member.ObjectID, usr_model.MemberTypeProjectGrant)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case proj_es_model.ProjectGrantMemberRemoved,
		proj_es_model.ProjectGrantMemberCascadeRemoved:
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
	project, err := m.getProjectByID(context.Background(), member.AggregateID)
	if err != nil {
		return err
	}
	member.DisplayName = project.Name
	return nil
}

func (m *UserMembership) updateProjectDisplayName(event *es_models.Event) error {
	project, err := m.getProjectByID(context.Background(), event.AggregateID)
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

func (m *UserMembership) processUser(event *es_models.Event) (err error) {
	switch event.Type {
	case model.UserRemoved:
		return m.view.DeleteUserMembershipsByUserID(event.AggregateID, event)
	default:
		return m.view.ProcessedUserMembershipSequence(event)
	}
}

func (m *UserMembership) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-Fwer2", "id", event.AggregateID).WithError(err).Warn("something went wrong in user membership handler")
	return spooler.HandleError(event, err, m.view.GetLatestUserMembershipFailedEvent, m.view.ProcessedUserMembershipFailedEvent, m.view.ProcessedUserMembershipSequence, m.errorCountUntilSkip)
}

func (m *UserMembership) OnSuccess() error {
	return spooler.HandleSuccess(m.view.UpdateUserMembershipSpoolerRunTimestamp)
}

func (u *UserMembership) getOrgByID(ctx context.Context, orgID string) (*org_model.Org, error) {
	query, err := org_view.OrgByIDQuery(orgID, 0)
	if err != nil {
		return nil, err
	}

	esOrg := &org_es_model.Org{
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

	return org_es_model.OrgToModel(esOrg), nil
}

func (u *UserMembership) getProjectByID(ctx context.Context, projID string) (*proj_model.Project, error) {
	query, err := proj_view.ProjectByIDQuery(projID, 0)
	if err != nil {
		return nil, err
	}
	esProject := &proj_es_model.Project{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: projID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, esProject.AppendEvents, query)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	if esProject.Sequence == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-Bg32b", "Errors.Project.NotFound")
	}

	return proj_es_model.ProjectToModel(esProject), nil
}
