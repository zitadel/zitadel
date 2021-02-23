package handler

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v1"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	"strings"

	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	proj_view "github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	grant_es_model "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/usergrant/repository/view/model"
	"github.com/caos/zitadel/internal/v2/domain"
)

const (
	userGrantTable = "auth.user_grants"
)

type UserGrant struct {
	handler
	iamID        string
	iamProjectID string
	subscription *v1.Subscription
}

func newUserGrant(
	handler handler,
	iamID string,
) *UserGrant {
	h := &UserGrant{
		handler: handler,
		iamID:   iamID,
	}

	h.subscribe()

	return h
}

func (k *UserGrant) subscribe() {
	k.subscription = k.es.Subscribe(k.AggregateTypes()...)
	go func() {
		for event := range k.subscription.Events {
			query.ReduceEvent(k, event)
		}
	}()
}

func (u *UserGrant) ViewModel() string {
	return userGrantTable
}

func (_ *UserGrant) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{grant_es_model.UserGrantAggregate, iam_es_model.IAMAggregate, org_es_model.OrgAggregate, usr_es_model.UserAggregate, proj_es_model.ProjectAggregate}
}

func (u *UserGrant) CurrentSequence() (uint64, error) {
	sequence, err := u.view.GetLatestUserGrantSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (u *UserGrant) EventQuery() (*es_models.SearchQuery, error) {
	if u.iamProjectID == "" {
		err := u.setIamProjectID()
		if err != nil {
			return nil, err
		}
	}
	sequence, err := u.view.GetLatestUserGrantSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(u.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (u *UserGrant) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case grant_es_model.UserGrantAggregate:
		err = u.processUserGrant(event)
	case usr_es_model.UserAggregate:
		err = u.processUser(event)
	case proj_es_model.ProjectAggregate:
		err = u.processProject(event)
	case iam_es_model.IAMAggregate:
		err = u.processIAMMember(event, "IAM", false)
	case org_es_model.OrgAggregate:
		return u.processOrg(event)
	}
	return err
}

func (u *UserGrant) processUserGrant(event *es_models.Event) (err error) {
	grant := new(view_model.UserGrantView)
	switch event.Type {
	case grant_es_model.UserGrantAdded:
		err = grant.AppendEvent(event)
		if err != nil {
			return err
		}
		err = u.fillData(grant, event.ResourceOwner)
	case grant_es_model.UserGrantChanged,
		grant_es_model.UserGrantCascadeChanged,
		grant_es_model.UserGrantDeactivated,
		grant_es_model.UserGrantReactivated:
		grant, err = u.view.UserGrantByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = grant.AppendEvent(event)
	case grant_es_model.UserGrantRemoved, grant_es_model.UserGrantCascadeRemoved:
		return u.view.DeleteUserGrant(event.AggregateID, event)
	default:
		return u.view.ProcessedUserGrantSequence(event)
	}
	if err != nil {
		return err
	}
	return u.view.PutUserGrant(grant, event)
}

func (u *UserGrant) processUser(event *es_models.Event) (err error) {
	switch event.Type {
	case usr_es_model.UserProfileChanged,
		usr_es_model.UserEmailChanged,
		usr_es_model.HumanProfileChanged,
		usr_es_model.HumanEmailChanged,
		usr_es_model.MachineChanged:
		grants, err := u.view.UserGrantsByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(grants) == 0 {
			return u.view.ProcessedUserGrantSequence(event)
		}
		user, err := u.getUserByID(event.AggregateID)
		if err != nil {
			return err
		}
		for _, grant := range grants {
			u.fillUserData(grant, user)
		}
		return u.view.PutUserGrants(grants, event)
	default:
		return u.view.ProcessedUserGrantSequence(event)
	}
}

func (u *UserGrant) processProject(event *es_models.Event) (err error) {
	switch event.Type {
	case proj_es_model.ProjectChanged:
		grants, err := u.view.UserGrantsByProjectID(event.AggregateID)
		if err != nil {
			return err
		}
		project, err := u.getProjectByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		for _, grant := range grants {
			u.fillProjectData(grant, project)
		}
		return u.view.PutUserGrants(grants, event)
	case proj_es_model.ProjectMemberAdded, proj_es_model.ProjectMemberChanged, proj_es_model.ProjectMemberRemoved:
		member := new(proj_es_model.ProjectMember)
		member.SetData(event)
		return u.processMember(event, "PROJECT", event.AggregateID, member.UserID, member.Roles)
	case proj_es_model.ProjectGrantMemberAdded, proj_es_model.ProjectGrantMemberChanged, proj_es_model.ProjectGrantMemberRemoved:
		member := new(proj_es_model.ProjectGrantMember)
		member.SetData(event)
		return u.processMember(event, "PROJECT_GRANT", member.GrantID, member.UserID, member.Roles)
	default:
		return u.view.ProcessedUserGrantSequence(event)
	}
}

func (u *UserGrant) processOrg(event *es_models.Event) (err error) {
	switch event.Type {
	case org_es_model.OrgMemberAdded, org_es_model.OrgMemberChanged, org_es_model.OrgMemberRemoved:
		member := new(org_es_model.OrgMember)
		member.SetData(event)
		return u.processMember(event, "ORG", "", member.UserID, member.Roles)
	default:
		return u.view.ProcessedUserGrantSequence(event)
	}
}

func (u *UserGrant) processIAMMember(event *es_models.Event, rolePrefix string, suffix bool) error {
	member := new(iam_es_model.IAMMember)

	switch event.Type {
	case iam_es_model.IAMMemberAdded, iam_es_model.IAMMemberChanged:
		member.SetData(event)

		grant, err := u.view.UserGrantByIDs(u.iamID, u.iamProjectID, member.UserID)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
		if errors.IsNotFound(err) {
			grant = &view_model.UserGrantView{
				ID:            u.iamProjectID + member.UserID,
				ResourceOwner: u.iamID,
				OrgName:       u.iamID,
				ProjectID:     u.iamProjectID,
				UserID:        member.UserID,
				RoleKeys:      member.Roles,
				CreationDate:  event.CreationDate,
			}
			if suffix {
				grant.RoleKeys = suffixRoles(event.AggregateID, grant.RoleKeys)
			}
		} else {
			newRoles := member.Roles
			if grant.RoleKeys != nil {
				grant.RoleKeys = mergeExistingRoles(rolePrefix, "", grant.RoleKeys, newRoles)
			} else {
				grant.RoleKeys = newRoles
			}
		}
		grant.Sequence = event.Sequence
		grant.ChangeDate = event.CreationDate
		return u.view.PutUserGrant(grant, event)
	case iam_es_model.IAMMemberRemoved:
		member.SetData(event)
		grant, err := u.view.UserGrantByIDs(u.iamID, u.iamProjectID, member.UserID)
		if err != nil {
			return err
		}
		return u.view.DeleteUserGrant(grant.ID, event)
	default:
		return u.view.ProcessedUserGrantSequence(event)
	}
}

func (u *UserGrant) processMember(event *es_models.Event, rolePrefix, roleSuffix string, userID string, roleKeys []string) error {
	switch event.Type {
	case org_es_model.OrgMemberAdded, proj_es_model.ProjectMemberAdded, proj_es_model.ProjectGrantMemberAdded,
		org_es_model.OrgMemberChanged, proj_es_model.ProjectMemberChanged, proj_es_model.ProjectGrantMemberChanged:

		grant, err := u.view.UserGrantByIDs(event.ResourceOwner, u.iamProjectID, userID)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
		if roleSuffix != "" {
			roleKeys = suffixRoles(roleSuffix, roleKeys)
		}
		if errors.IsNotFound(err) {
			grant = &view_model.UserGrantView{
				ID:            u.iamProjectID + event.ResourceOwner + userID,
				ResourceOwner: event.ResourceOwner,
				ProjectID:     u.iamProjectID,
				UserID:        userID,
				RoleKeys:      roleKeys,
				CreationDate:  event.CreationDate,
			}
			u.fillData(grant, event.ResourceOwner)
		} else {
			newRoles := roleKeys
			if grant.RoleKeys != nil {
				grant.RoleKeys = mergeExistingRoles(rolePrefix, roleSuffix, grant.RoleKeys, newRoles)
			} else {
				grant.RoleKeys = newRoles
			}
		}
		grant.Sequence = event.Sequence
		grant.ChangeDate = event.CreationDate
		return u.view.PutUserGrant(grant, event)
	case org_es_model.OrgMemberRemoved,
		proj_es_model.ProjectMemberRemoved,
		proj_es_model.ProjectGrantMemberRemoved:

		grant, err := u.view.UserGrantByIDs(event.ResourceOwner, u.iamProjectID, userID)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
		if errors.IsNotFound(err) {
			return u.view.ProcessedUserGrantSequence(event)
		}
		if roleSuffix != "" {
			roleKeys = suffixRoles(roleSuffix, roleKeys)
		}
		if grant.RoleKeys == nil {
			return u.view.ProcessedUserGrantSequence(event)
		}
		grant.RoleKeys = mergeExistingRoles(rolePrefix, roleSuffix, grant.RoleKeys, nil)
		return u.view.PutUserGrant(grant, event)
	default:
		return u.view.ProcessedUserGrantSequence(event)
	}
}

func suffixRoles(suffix string, roles []string) []string {
	suffixedRoles := make([]string, len(roles))
	for i := 0; i < len(roles); i++ {
		suffixedRoles[i] = roles[i] + ":" + suffix
	}
	return suffixedRoles
}

func mergeExistingRoles(rolePrefix, suffix string, existingRoles, newRoles []string) []string {
	mergedRoles := make([]string, 0)
	for _, existingRole := range existingRoles {
		if !strings.HasPrefix(existingRole, rolePrefix) {
			mergedRoles = append(mergedRoles, existingRole)
			continue
		}
		if suffix != "" && !strings.HasSuffix(existingRole, suffix) {
			mergedRoles = append(mergedRoles, existingRole)
		}
	}
	return append(mergedRoles, newRoles...)
}

func (u *UserGrant) setIamProjectID() error {
	if u.iamProjectID != "" {
		return nil
	}
	iam, err := u.getIAMByID(context.Background())
	if err != nil {
		return err
	}

	if iam.SetUpDone < domain.StepCount-1 {
		return caos_errs.ThrowPreconditionFailed(nil, "HANDL-s5DTs", "Setup not done")
	}
	u.iamProjectID = iam.IAMProjectID
	return nil
}

func (u *UserGrant) fillData(grant *view_model.UserGrantView, resourceOwner string) (err error) {
	user, err := u.getUserByID(grant.UserID)
	if err != nil {
		return err
	}
	u.fillUserData(grant, user)
	project, err := u.getProjectByID(context.Background(), grant.ProjectID)
	if err != nil {
		return err
	}
	u.fillProjectData(grant, project)

	org, err := u.getOrgByID(context.TODO(), resourceOwner)
	if err != nil {
		return err
	}
	u.fillOrgData(grant, org)
	return nil
}

func (u *UserGrant) fillUserData(grant *view_model.UserGrantView, user *model.UserView) {
	grant.UserName = user.UserName
	if user.HumanView != nil {
		grant.FirstName = user.FirstName
		grant.LastName = user.LastName
		grant.DisplayName = user.FirstName + " " + user.LastName
		grant.Email = user.Email
	}
	if user.MachineView != nil {
		grant.DisplayName = user.MachineView.Name
	}
}

func (u *UserGrant) fillProjectData(grant *view_model.UserGrantView, project *proj_model.Project) {
	grant.ProjectName = project.Name
	grant.ProjectOwner = project.ResourceOwner
}

func (u *UserGrant) fillOrgData(grant *view_model.UserGrantView, org *org_model.Org) {
	grant.OrgName = org.Name
	for _, domain := range org.Domains {
		if domain.Primary {
			grant.OrgPrimaryDomain = domain.Domain
			break
		}
	}
}

func (u *UserGrant) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-UZmc7", "id", event.AggregateID).WithError(err).Warn("something went wrong in user grant handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserGrantFailedEvent, u.view.ProcessedUserGrantFailedEvent, u.view.ProcessedUserGrantSequence, u.errorCountUntilSkip)
}

func (u *UserGrant) OnSuccess() error {
	return spooler.HandleSuccess(u.view.UpdateUserGrantSpoolerRunTimestamp)
}

func (u *UserGrant) getUserByID(userID string) (*model.UserView, error) {
	user, usrErr := u.view.UserByID(userID)
	if usrErr != nil && !caos_errs.IsNotFound(usrErr) {
		return nil, usrErr
	}
	if user == nil {
		user = &model.UserView{}
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

func (u *UserGrant) getUserEvents(userID string, sequence uint64) ([]*es_models.Event, error) {
	query, err := view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}

	return u.es.FilterEvents(context.Background(), query)
}

func (u *UserGrant) getOrgByID(ctx context.Context, orgID string) (*org_model.Org, error) {
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

func (u *UserGrant) getProjectByID(ctx context.Context, projID string) (*proj_model.Project, error) {
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
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if esProject.Sequence == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-DAfng", "Errors.Project.NotFound")
	}

	return proj_es_model.ProjectToModel(esProject), nil
}

func (u *UserGrant) getIAMByID(ctx context.Context) (*iam_model.IAM, error) {
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
	if err != nil && errors.IsNotFound(err) && iam.Sequence == 0 {
		return nil, err
	}
	return iam_es_model.IAMToModel(iam), nil
}
