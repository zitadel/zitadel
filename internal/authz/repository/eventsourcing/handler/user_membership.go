package handler

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	proj_view "github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/repository/user"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	userMembershipTable = "authz.user_memberships"
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

func (m *UserMembership) Subscription() *v1.Subscription {
	return m.subscription
}

func (_ *UserMembership) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{instance.AggregateType, org.AggregateType, project.AggregateType, user.AggregateType}
}

func (m *UserMembership) CurrentSequence(instanceID string) (uint64, error) {
	sequence, err := m.view.GetLatestUserMembershipSequence(instanceID)
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *UserMembership) EventQuery() (*es_models.SearchQuery, error) {
	sequences, err := m.view.GetLatestUserMembershipSequences()
	if err != nil {
		return nil, err
	}
	query := es_models.NewSearchQuery()
	instances := make([]string, 0)
	for _, sequence := range sequences {
		for _, instance := range instances {
			if sequence.InstanceID == instance {
				break
			}
		}
		instances = append(instances, sequence.InstanceID)
		query.AddQuery().
			AggregateTypeFilter(m.AggregateTypes()...).
			LatestSequenceFilter(sequence.CurrentSequence).
			InstanceIDFilter(sequence.InstanceID)
	}
	return query.AddQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(0).
		IgnoredInstanceIDsFilter(instances...).
		SearchQuery(), nil
}

func (m *UserMembership) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case instance.AggregateType:
		err = m.processIAM(event)
	case org.AggregateType:
		err = m.processOrg(event)
	case project.AggregateType:
		err = m.processProject(event)
	case user.AggregateType:
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
	switch eventstore.EventType(event.Type) {
	case instance.MemberAddedEventType:
		m.fillIamDisplayName(member)
	case instance.MemberChangedEventType:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, event.AggregateID, event.InstanceID, usr_model.MemberTypeIam)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case instance.MemberRemovedEventType,
		instance.MemberCascadeRemovedEventType:
		return m.view.DeleteUserMembership(member.UserID, event.AggregateID, event.AggregateID, event.InstanceID, usr_model.MemberTypeIam, event)
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

func (m *UserMembership) processOrg(event *es_models.Event) (err error) {
	member := new(usr_es_model.UserMembershipView)
	err = member.AppendEvent(event)
	if err != nil {
		return err
	}
	switch eventstore.EventType(event.Type) {
	case org.MemberAddedEventType:
		err = m.fillOrgName(member)
	case org.MemberChangedEventType:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, event.AggregateID, event.InstanceID, usr_model.MemberTypeOrganisation)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case org.MemberRemovedEventType,
		org.MemberCascadeRemovedEventType:
		return m.view.DeleteUserMembership(member.UserID, event.AggregateID, event.AggregateID, event.InstanceID, usr_model.MemberTypeOrganisation, event)
	case org.OrgChangedEventType:
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
	org, err := m.getOrgByID(context.Background(), member.ResourceOwner)
	if err != nil {
		return err
	}
	member.ResourceOwnerName = org.Name
	if member.AggregateID == org.AggregateID {
		member.DisplayName = org.Name
	}
	return nil
}

func (m *UserMembership) updateOrgName(event *es_models.Event) error {
	org, err := m.getOrgByID(context.Background(), event.AggregateID)
	if err != nil {
		return err
	}

	memberships, err := m.view.UserMembershipsByResourceOwner(event.ResourceOwner, event.InstanceID)
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

func (m *UserMembership) processProject(event *es_models.Event) (err error) {
	member := new(usr_es_model.UserMembershipView)
	err = member.AppendEvent(event)
	if err != nil {
		return err
	}
	switch eventstore.EventType(event.Type) {
	case project.MemberAddedType, project.GrantMemberAddedType:
		err = m.fillProjectDisplayName(member)
		if err != nil {
			return err
		}
		err = m.fillOrgName(member)
	case project.MemberChangedType:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, event.AggregateID, event.InstanceID, usr_model.MemberTypeProject)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case project.MemberRemovedType, project.MemberCascadeRemovedType:
		return m.view.DeleteUserMembership(member.UserID, event.AggregateID, event.AggregateID, event.InstanceID, usr_model.MemberTypeProject, event)
	case project.GrantMemberChangedType:
		member, err = m.view.UserMembershipByIDs(member.UserID, event.AggregateID, member.ObjectID, event.InstanceID, usr_model.MemberTypeProjectGrant)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case project.GrantMemberRemovedType,
		project.GrantMemberCascadeRemovedType:
		return m.view.DeleteUserMembership(member.UserID, event.AggregateID, member.ObjectID, member.InstanceID, usr_model.MemberTypeProjectGrant, event)
	case project.ProjectChangedType:
		return m.updateProjectDisplayName(event)
	case project.ProjectRemovedType:
		return m.view.DeleteUserMembershipsByAggregateID(event.AggregateID, event.InstanceID, event)
	case project.GrantRemovedType:
		return m.view.DeleteUserMembershipsByAggregateIDAndObjectID(event.AggregateID, member.ObjectID, member.InstanceID, event)
	default:
		return m.view.ProcessedUserMembershipSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutUserMembership(member, event)
}

func (m *UserMembership) fillProjectDisplayName(member *usr_es_model.UserMembershipView) (err error) {
	project, err := m.getProjectByID(context.Background(), member.AggregateID, member.InstanceID)
	if err != nil {
		return err
	}
	member.DisplayName = project.Name
	return nil
}

func (m *UserMembership) updateProjectDisplayName(event *es_models.Event) error {
	proj := new(proj_es_model.Project)
	err := proj.SetData(event)
	if err != nil {
		return err
	}
	if proj.Name == "" {
		return m.view.ProcessedUserMembershipSequence(event)
	}

	memberships, err := m.view.UserMembershipsByAggregateID(event.AggregateID, event.InstanceID)
	if err != nil {
		return err
	}
	for _, membership := range memberships {
		membership.DisplayName = proj.Name
	}
	return m.view.BulkPutUserMemberships(memberships, event)
}

func (m *UserMembership) processUser(event *es_models.Event) (err error) {
	switch eventstore.EventType(event.Type) {
	case user.UserRemovedType:
		return m.view.DeleteUserMembershipsByUserID(event.AggregateID, event.InstanceID, event)
	default:
		return m.view.ProcessedUserMembershipSequence(event)
	}
}

func (m *UserMembership) OnError(event *es_models.Event, err error) error {
	logging.WithFields("id", event.AggregateID).WithError(err).Warn("something went wrong in user membership handler")
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
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-3m9vs", "Errors.Org.NotFound")
	}

	return org_es_model.OrgToModel(esOrg), nil
}

func (u *UserMembership) getProjectByID(ctx context.Context, projID, instanceID string) (*proj_model.Project, error) {
	query, err := proj_view.ProjectByIDQuery(projID, instanceID, 0)
	if err != nil {
		return nil, err
	}
	esProject := &proj_es_model.Project{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: projID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, esProject.AppendEvents, query)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if esProject.Sequence == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-Dfrt2", "Errors.Project.NotFound")
	}

	return proj_es_model.ProjectToModel(esProject), nil
}
