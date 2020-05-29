package handler

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_events "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_events "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	grant_es_model "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/usergrant/repository/view/model"
	"strings"
	"time"
)

type UserGrant struct {
	handler
	eventstore    eventstore.Eventstore
	projectEvents *proj_event.ProjectEventstore
	userEvents    *usr_events.UserEventstore
	orgEvents     *org_events.OrgEventstore
	iamProjectID  string
}

const (
	userGrantTable = "auth.user_grants"
)

func (u *UserGrant) MinimumCycleDuration() time.Duration { return u.cycleDuration }

func (u *UserGrant) ViewModel() string {
	return userGrantTable
}

func (u *UserGrant) EventQuery() (*models.SearchQuery, error) {
	if u.iamProjectID == "" {
		u.setIamProjectID()
	}
	sequence, err := u.view.GetLatestUserGrantSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(grant_es_model.UserGrantAggregate, usr_es_model.UserAggregate, proj_es_model.ProjectAggregate).
		LatestSequenceFilter(sequence), nil
}

func (u *UserGrant) Process(event *models.Event) (err error) {
	switch event.AggregateType {
	case grant_es_model.UserGrantAggregate:
		err = u.processUserGrant(event)
	case usr_es_model.UserAggregate:
		err = u.processUser(event)
	case proj_es_model.ProjectAggregate:
		err = u.processProject(event)
	}
	return err
}

func (u *UserGrant) processUserGrant(event *models.Event) (err error) {
	grant := new(view_model.UserGrantView)
	switch event.Type {
	case grant_es_model.UserGrantAdded:
		err = grant.AppendEvent(event)
		if err != nil {
			return err
		}
		err = u.fillData(grant, event.ResourceOwner)
	case grant_es_model.UserGrantChanged,
		grant_es_model.UserGrantDeactivated,
		grant_es_model.UserGrantReactivated:
		grant, err = u.view.UserGrantByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = grant.AppendEvent(event)
	case grant_es_model.UserGrantRemoved:
		err = u.view.DeleteUserGrant(event.AggregateID, event.Sequence)
	default:
		return u.view.ProcessedUserGrantSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return u.view.PutUserGrant(grant, grant.Sequence)
}

func (u *UserGrant) processUser(event *models.Event) (err error) {
	switch event.Type {
	case usr_es_model.UserProfileChanged,
		usr_es_model.UserEmailChanged:
		grants, err := u.view.UserGrantsByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		user, err := u.userEvents.UserByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		for _, grant := range grants {
			u.fillUserData(grant, user)
			err = u.view.PutUserGrant(grant, event.Sequence)
			if err != nil {
				return err
			}
		}
	default:
		return u.view.ProcessedUserGrantSequence(event.Sequence)
	}
	return nil
}

func (u *UserGrant) processProject(event *models.Event) (err error) {
	switch event.Type {
	case proj_es_model.ProjectChanged:
		grants, err := u.view.UserGrantsByProjectID(event.AggregateID)
		if err != nil {
			return err
		}
		project, err := u.projectEvents.ProjectByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		for _, grant := range grants {
			u.fillProjectData(grant, project)
			return u.view.PutUserGrant(grant, event.Sequence)
		}
	default:
		return u.view.ProcessedUserGrantSequence(event.Sequence)
	}
	return nil
}

//
//func (u *UserGrant) processMember(rolePrefix string, event *models.Event, suffix bool) error {
//	member := &model.Member{}
//
//	switch event.Type {
//	case org_es_model.OrgMemberAdded, proj_es_model.ProjectMemberAdded, proj_es_model.ProjectGrantMemberAdded,
//		org_pkg.ChangeMember, proj_es_model.ProjectMemberChanged, proj_es_model.ProjectGrantMemberChanged:
//
//		if err := proto.FromPBStruct(member, event.Data); err != nil {
//			logging.Log("VIEW-Nt7QE").WithError(err).Debug("unable to map data into user")
//			return err
//		}
//
//		grant, err := u.cache.LatestGrant(event.ResourceOwner, u.iamProjectID, member.UserID)
//		if err != nil && !errors.IsNotFound(err) {
//			return err
//		}
//		if suffix {
//			member.Roles = suffixRoles(event.GetAggregateId(), member.Roles)
//		}
//		if errors.IsNotFound(err) {
//			grant = &model.Grant{
//				OrgID:     event.ResourceOwner,
//				ProjectID: g.iamProjectID,
//				UserID:    member.UserID,
//				Roles:     member.Roles,
//			}
//			g.fillOrg(grant)
//		} else {
//			newRoles := member.Roles
//			if grant.Roles != nil {
//				grant.Roles = mergeExistingRoles(rolePrefix, grant.Roles, newRoles)
//			} else {
//				grant.Roles = newRoles
//			}
//		}
//		grant.CurrentSequence = event.Sequence
//		return g.cache.PutGrant(grant)
//	case org_pkg.RemoveMember,
//		pro_pkg.RemovedMember,
//		pro_pkg.RemovedGrantMember:
//
//		if err := proto.FromPBStruct(member, event.Data); err != nil {
//			logging.Log("VIEW-Nt7QE").WithError(err).Debug("unable to map data into user")
//			return err
//		}
//		grant, err := g.cache.LatestGrant(event.ResourceOwner, g.iamProjectID, member.UserID)
//		if err != nil {
//			return err
//		}
//		return g.cache.RemoveGrant(grant, event.Sequence)
//	default:
//		return g.cache.ProcessedGrantSequence(event.Sequence)
//	}
//}
//
//func (u *UserGrant) processIamMember(rolePrefix string, event *models.Event, suffix bool) error {
//	member := &model.Member{}
//
//	switch event.EventType() {
//	case iam_pkg.AddedIamMember, iam_pkg.ChangedIamMember:
//		if err := proto.FromPBStruct(member, event.Data); err != nil {
//			logging.Log("VIEW-UeJ58").WithError(err).Debug("unable to map data into member")
//			return err
//		}
//		grant, err := g.cache.LatestGrant(iamResourceOwner, g.iamProjectID, member.UserID)
//		if err != nil && !errors.IsNotFound(err) {
//			return err
//		}
//		if errors.IsNotFound(err) {
//			grant = &model.Grant{
//				OrgID:     iamResourceOwner,
//				OrgName:   iamResourceOwner,
//				ProjectID: g.iamProjectID,
//				UserID:    member.UserID,
//				Roles:     member.Roles,
//			}
//			if suffix {
//				grant.Roles = suffixRoles(event.GetAggregateId(), grant.Roles)
//			}
//		} else {
//			newRoles := member.Roles
//			if grant.Roles != nil {
//				grant.Roles = mergeExistingRoles(rolePrefix, grant.Roles, newRoles)
//			} else {
//				grant.Roles = newRoles
//			}
//
//		}
//		grant.CurrentSequence = event.Sequence
//		return g.cache.PutGrant(grant)
//	case iam_pkg.RemovedIamMember:
//		if err := proto.FromPBStruct(member, event.Data); err != nil {
//			logging.Log("VIEW-Mi7Er").WithError(err).Debug("unable to map data into user")
//			return err
//		}
//		grant, err := g.cache.LatestGrant(iamResourceOwner, g.iamProjectID, member.UserID)
//		if err != nil {
//			return err
//		}
//		return g.cache.RemoveGrant(grant, event.Sequence)
//	default:
//		return g.cache.ProcessedGrantSequence(event.Sequence)
//	}
//}

func suffixRoles(suffix string, roles []string) []string {
	suffixedRoles := make([]string, len(roles))
	for i := 0; i < len(roles); i++ {
		suffixedRoles[i] = roles[i] + ":" + suffix
	}
	return suffixedRoles
}

func mergeExistingRoles(rolePrefix string, existingRoles, newRoles []string) []string {
	mergedRoles := make([]string, 0)
	for _, existing := range existingRoles {
		if !strings.HasPrefix(existing, rolePrefix) {
			mergedRoles = append(mergedRoles, existing)
		}
	}
	return append(mergedRoles, newRoles...)
}

func (u *UserGrant) setIamProjectID() {
	filter := es_models.NewSearchQuery().
		AggregateTypeFilter(iam_es_model.IamAggregate).
		LatestSequenceFilter(0)

	events, err := u.eventstore.FilterEvents(context.Background(), filter)
	if err != nil {
		return
	}
	if len(events) == 0 {
		return
	}
	for _, e := range events {
		if e.Type == iam_es_model.IamProjectSet {
			iam := &iam_es_model.Iam{}
			iam.SetData(e)
			u.iamProjectID = iam.IamProjectID
		}
	}
}

func (u *UserGrant) fillData(grant *view_model.UserGrantView, resourceOwner string) (err error) {
	user, err := u.userEvents.UserByID(context.Background(), grant.UserID)
	if err != nil {
		return err
	}
	u.fillUserData(grant, user)
	project, err := u.projectEvents.ProjectByID(context.Background(), grant.ProjectID)
	if err != nil {
		return err
	}
	u.fillProjectData(grant, project)

	org, err := u.orgEvents.OrgByID(context.TODO(), org_model.NewOrg(resourceOwner))
	if err != nil {
		return err
	}
	u.fillOrgData(grant, org)
	return nil
}

func (u *UserGrant) fillUserData(grant *view_model.UserGrantView, user *usr_model.User) {
	grant.UserName = user.UserName
	grant.FirstName = user.FirstName
	grant.LastName = user.LastName
	grant.Email = user.EmailAddress
}

func (u *UserGrant) fillProjectData(grant *view_model.UserGrantView, project *proj_model.Project) {
	grant.ProjectName = project.Name
}

func (u *UserGrant) fillOrgData(grant *view_model.UserGrantView, org *org_model.Org) {
	grant.OrgDomain = org.Domain
	grant.OrgName = org.Name
}

func (u *UserGrant) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-8is4s", "id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserGrantFailedEvent, u.view.ProcessedUserGrantFailedEvent, u.view.ProcessedUserGrantSequence, u.errorCountUntilSkip)
}
