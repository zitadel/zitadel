package handler

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view"
	usr_view_model "github.com/caos/zitadel/internal/user/repository/view/model"

	"github.com/caos/logging"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	projectGrantMemberTable = "management.project_grant_members"
)

type ProjectGrantMember struct {
	handler
	subscription *v1.Subscription
}

func newProjectGrantMember(
	handler handler,
) *ProjectGrantMember {
	h := &ProjectGrantMember{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *ProjectGrantMember) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (p *ProjectGrantMember) ViewModel() string {
	return projectGrantMemberTable
}

func (_ *ProjectGrantMember) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{proj_es_model.ProjectAggregate, usr_es_model.UserAggregate}
}

func (p *ProjectGrantMember) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestProjectGrantMemberSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *ProjectGrantMember) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectGrantMemberSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *ProjectGrantMember) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case proj_es_model.ProjectAggregate:
		err = p.processProjectGrantMember(event)
	case usr_es_model.UserAggregate:
		err = p.processUser(event)
	}
	return err
}

func (p *ProjectGrantMember) processProjectGrantMember(event *es_models.Event) (err error) {
	member := new(view_model.ProjectGrantMemberView)
	switch event.Type {
	case proj_es_model.ProjectGrantMemberAdded:
		err = member.AppendEvent(event)
		if err != nil {
			return err
		}
		err = p.fillData(member)
	case proj_es_model.ProjectGrantMemberChanged:
		err = member.SetData(event)
		if err != nil {
			return err
		}
		member, err = p.view.ProjectGrantMemberByIDs(member.GrantID, member.UserID)
		if err != nil {
			return err
		}
		err = member.AppendEvent(event)
	case proj_es_model.ProjectGrantMemberRemoved:
		err = member.SetData(event)
		if err != nil {
			return err
		}
		return p.view.DeleteProjectGrantMember(member.GrantID, member.UserID, event)
	case proj_es_model.ProjectRemoved:
		err = p.view.DeleteProjectGrantMembersByProjectID(event.AggregateID)
		if err != nil {
			return err
		}
		return p.view.ProcessedProjectGrantMemberSequence(event)
	default:
		return p.view.ProcessedProjectGrantMemberSequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutProjectGrantMember(member, event)
}

func (p *ProjectGrantMember) processUser(event *es_models.Event) (err error) {
	switch event.Type {
	case usr_es_model.UserProfileChanged,
		usr_es_model.UserEmailChanged,
		usr_es_model.HumanProfileChanged,
		usr_es_model.HumanEmailChanged,
		usr_es_model.MachineChanged:
		members, err := p.view.ProjectGrantMembersByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(members) == 0 {
			return p.view.ProcessedProjectGrantMemberSequence(event)
		}
		user, err := p.getUserByID(event.AggregateID)
		if err != nil {
			return err
		}
		for _, member := range members {
			p.fillUserData(member, user)
		}
		return p.view.PutProjectGrantMembers(members, event)
	default:
		return p.view.ProcessedProjectGrantMemberSequence(event)
	}
}

func (p *ProjectGrantMember) fillData(member *view_model.ProjectGrantMemberView) (err error) {
	user, err := p.getUserByID(member.UserID)
	if err != nil {
		return err
	}
	p.fillUserData(member, user)
	return nil
}

func (p *ProjectGrantMember) fillUserData(member *view_model.ProjectGrantMemberView, user *usr_view_model.UserView) error {
	org, err := p.getOrgByID(context.Background(), user.ResourceOwner)
	policy := org.OrgIamPolicy
	if policy == nil {
		policy, err = p.getDefaultOrgIAMPolicy(context.TODO())
		if err != nil {
			return err
		}
	}
	member.UserName = user.UserName
	member.PreferredLoginName = user.GenerateLoginName(org.GetPrimaryDomain().Domain, policy.UserLoginMustBeDomain)
	if user.HumanView != nil {
		member.FirstName = user.FirstName
		member.LastName = user.LastName
		member.DisplayName = user.DisplayName
		member.Email = user.Email
	}
	if user.MachineView != nil {
		member.DisplayName = user.MachineView.Name
	}
	return nil
}

func (p *ProjectGrantMember) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-kls93", "id", event.AggregateID).WithError(err).Warn("something went wrong in projectmember handler")
	return spooler.HandleError(event, err, p.view.GetLatestProjectGrantMemberFailedEvent, p.view.ProcessedProjectGrantMemberFailedEvent, p.view.ProcessedProjectGrantMemberSequence, p.errorCountUntilSkip)
}

func (p *ProjectGrantMember) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateProjectGrantMemberSpoolerRunTimestamp)
}

func (u *ProjectGrantMember) getUserByID(userID string) (*usr_view_model.UserView, error) {
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

func (u *ProjectGrantMember) getUserEvents(userID string, sequence uint64) ([]*es_models.Event, error) {
	query, err := view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}

	return u.es.FilterEvents(context.Background(), query)
}

func (u *ProjectGrantMember) getOrgByID(ctx context.Context, orgID string) (*org_model.Org, error) {
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
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-3nd7s", "Errors.Org.NotFound")
	}

	return org_es_model.OrgToModel(esOrg), nil
}

func (u *ProjectGrantMember) getDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicy, error) {
	existingIAM, err := u.getIAMByID(ctx)
	if err != nil {
		return nil, err
	}
	if existingIAM.DefaultOrgIAMPolicy == nil {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-3Bf7s", "Errors.IAM.OrgIAMPolicy.NotExisting")
	}
	return existingIAM.DefaultOrgIAMPolicy, nil
}

func (u *ProjectGrantMember) getIAMByID(ctx context.Context) (*iam_model.IAM, error) {
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
