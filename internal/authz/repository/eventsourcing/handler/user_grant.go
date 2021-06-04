package handler

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	"strings"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/usergrant/repository/view/model"
)

const (
	userGrantTable = "authz.user_grants"
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
	return []es_models.AggregateType{iam_es_model.IAMAggregate, org_es_model.OrgAggregate, proj_es_model.ProjectAggregate}
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
		AggregateTypeFilter(iam_es_model.IAMAggregate, org_es_model.OrgAggregate, proj_es_model.ProjectAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (u *UserGrant) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case proj_es_model.ProjectAggregate:
		err = u.processProject(event)
	case iam_es_model.IAMAggregate:
		err = u.processIAMMember(event, "IAM", false)
	case org_es_model.OrgAggregate:
		return u.processOrg(event)
	}
	return err
}

func (u *UserGrant) processProject(event *es_models.Event) (err error) {
	switch event.Type {
	case proj_es_model.ProjectMemberAdded, proj_es_model.ProjectMemberChanged,
		proj_es_model.ProjectMemberRemoved, proj_es_model.ProjectMemberCascadeRemoved:
		member := new(proj_es_model.ProjectMember)
		member.SetData(event)
		return u.processMember(event, "PROJECT", event.AggregateID, member.UserID, member.Roles)
	case proj_es_model.ProjectGrantMemberAdded, proj_es_model.ProjectGrantMemberChanged,
		proj_es_model.ProjectGrantMemberRemoved,
		proj_es_model.ProjectGrantMemberCascadeRemoved:
		member := new(proj_es_model.ProjectGrantMember)
		member.SetData(event)
		return u.processMember(event, "PROJECT_GRANT", member.GrantID, member.UserID, member.Roles)
	default:
		return u.view.ProcessedUserGrantSequence(event)
	}
}

func (u *UserGrant) processOrg(event *es_models.Event) (err error) {
	switch event.Type {
	case org_es_model.OrgMemberAdded, org_es_model.OrgMemberChanged,
		org_es_model.OrgMemberRemoved, org_es_model.OrgMemberCascadeRemoved:
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
	case iam_es_model.IAMMemberRemoved,
		iam_es_model.IAMMemberCascadeRemoved:
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
		org_es_model.OrgMemberCascadeRemoved,
		proj_es_model.ProjectMemberRemoved,
		proj_es_model.ProjectMemberCascadeRemoved,
		proj_es_model.ProjectGrantMemberRemoved,
		proj_es_model.ProjectGrantMemberCascadeRemoved:

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

func (u *UserGrant) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-VcVoJ", "id", event.AggregateID).WithError(err).Warn("something went wrong in user grant handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserGrantFailedEvent, u.view.ProcessedUserGrantFailedEvent, u.view.ProcessedUserGrantSequence, u.errorCountUntilSkip)
}

func (u *UserGrant) OnSuccess() error {
	return spooler.HandleSuccess(u.view.UpdateUserGrantSpoolerRunTimestamp)
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
