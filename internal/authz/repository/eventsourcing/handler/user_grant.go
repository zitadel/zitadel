package handler

import (
	"context"
	"strings"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_events "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/usergrant/repository/view/model"
)

type UserGrant struct {
	handler
	eventstore   eventstore.Eventstore
	iamEvents    *iam_events.IamEventstore
	iamID        string
	iamProjectID string
}

const (
	userGrantTable = "authz.user_grants"
)

func (u *UserGrant) MinimumCycleDuration() time.Duration { return u.cycleDuration }

func (u *UserGrant) ViewModel() string {
	return userGrantTable
}

func (u *UserGrant) EventQuery() (*models.SearchQuery, error) {
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
		AggregateTypeFilter(iam_es_model.IamAggregate, org_es_model.OrgAggregate, proj_es_model.ProjectAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (u *UserGrant) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case proj_es_model.ProjectAggregate:
		err = u.processProject(event)
	case iam_es_model.IamAggregate:
		err = u.processIamMember(event, "IAM", false)
	case org_es_model.OrgAggregate:
		return u.processOrg(event)
	}
	return err
}

func (u *UserGrant) processProject(event *models.Event) (err error) {
	switch event.Type {
	case proj_es_model.ProjectMemberAdded, proj_es_model.ProjectMemberChanged, proj_es_model.ProjectMemberRemoved:
		member := new(proj_es_model.ProjectMember)
		member.SetData(event)
		return u.processMember(event, "PROJECT", event.AggregateID, member.UserID, member.Roles)
	case proj_es_model.ProjectGrantMemberAdded, proj_es_model.ProjectGrantMemberChanged, proj_es_model.ProjectGrantMemberRemoved:
		member := new(proj_es_model.ProjectGrantMember)
		member.SetData(event)
		return u.processMember(event, "PROJECT_GRANT", member.GrantID, member.UserID, member.Roles)
	default:
		return u.view.ProcessedUserGrantSequence(event.Sequence)
	}
	return nil
}

func (u *UserGrant) processOrg(event *models.Event) (err error) {
	switch event.Type {
	case org_es_model.OrgMemberAdded, org_es_model.OrgMemberChanged, org_es_model.OrgMemberRemoved:
		member := new(org_es_model.OrgMember)
		member.SetData(event)
		return u.processMember(event, "ORG", "", member.UserID, member.Roles)
	default:
		return u.view.ProcessedUserGrantSequence(event.Sequence)
	}
	return nil
}

func (u *UserGrant) processIamMember(event *models.Event, rolePrefix string, suffix bool) error {
	member := new(iam_es_model.IamMember)

	switch event.Type {
	case iam_es_model.IamMemberAdded, iam_es_model.IamMemberChanged:
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
		return u.view.PutUserGrant(grant, grant.Sequence)
	case iam_es_model.IamMemberRemoved:
		member.SetData(event)
		grant, err := u.view.UserGrantByIDs(u.iamID, u.iamProjectID, member.UserID)
		if err != nil {
			return err
		}
		return u.view.DeleteUserGrant(grant.ID, event.Sequence)
	default:
		return u.view.ProcessedUserGrantSequence(event.Sequence)
	}
}

func (u *UserGrant) processMember(event *models.Event, rolePrefix, roleSuffix string, userID string, roleKeys []string) error {
	switch event.Type {
	case org_es_model.OrgMemberAdded, proj_es_model.ProjectMemberAdded, proj_es_model.ProjectGrantMemberAdded,
		org_es_model.OrgMemberChanged, proj_es_model.ProjectMemberChanged, proj_es_model.ProjectGrantMemberChanged:

		grant, err := u.view.UserGrantByIDs(event.ResourceOwner, u.iamProjectID, userID)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
		if roleSuffix != "" {
			roleKeys = suffixRoles(event.AggregateID, roleKeys)
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
		return u.view.PutUserGrant(grant, event.Sequence)
	case org_es_model.OrgMemberRemoved,
		proj_es_model.ProjectMemberRemoved,
		proj_es_model.ProjectGrantMemberRemoved:

		grant, err := u.view.UserGrantByIDs(event.ResourceOwner, u.iamProjectID, userID)
		if err != nil {
			return err
		}
		return u.view.DeleteUserGrant(grant.ID, event.Sequence)
	default:
		return u.view.ProcessedUserGrantSequence(event.Sequence)
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
	for _, existing := range existingRoles {
		if !strings.HasPrefix(existing, rolePrefix) {
			mergedRoles = append(mergedRoles, existing)
			continue
		}
		if suffix != "" && !strings.HasSuffix(existing, suffix) {
			mergedRoles = append(mergedRoles, existing)
		}
	}
	return append(mergedRoles, newRoles...)
}

func (u *UserGrant) setIamProjectID() error {
	if u.iamProjectID != "" {
		return nil
	}
	iam, err := u.iamEvents.IamByID(context.Background(), u.iamID)
	if err != nil {
		return err
	}
	if !iam.SetUpDone {
		return caos_errs.ThrowPreconditionFailed(nil, "HANDL-s5DTs", "Setup not done")
	}
	u.iamProjectID = iam.IamProjectID
	return nil
}

func (u *UserGrant) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-VcVoJ", "id", event.AggregateID).WithError(err).Warn("something went wrong in user grant handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserGrantFailedEvent, u.view.ProcessedUserGrantFailedEvent, u.view.ProcessedUserGrantSequence, u.errorCountUntilSkip)
}
