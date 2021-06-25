package command

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/idpconfig"
	"github.com/caos/zitadel/internal/repository/member"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/repository/usergrant"
)

type UniqueConstraintReadModel struct {
	eventstore.WriteModel

	UniqueConstraints []*domain.UniqueConstraintMigration
	commandProvider   commandProvider
	ctx               context.Context
}

type commandProvider interface {
	getOrgIAMPolicy(ctx context.Context, orgID string) (*domain.OrgIAMPolicy, error)
}

func NewUniqueConstraintReadModel(ctx context.Context, provider commandProvider) *UniqueConstraintReadModel {
	return &UniqueConstraintReadModel{
		ctx:             ctx,
		commandProvider: provider,
	}
}

func (rm *UniqueConstraintReadModel) AppendEvents(events ...eventstore.EventReader) {
	rm.WriteModel.AppendEvents(events...)
}

func (rm *UniqueConstraintReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *org.OrgAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, org.NewAddOrgNameUniqueConstraint(e.Name))
		case *org.OrgChangedEvent:
			rm.changeUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, org.NewAddOrgNameUniqueConstraint(e.Name))
		case *org.DomainVerifiedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, org.NewAddOrgDomainUniqueConstraint(e.Domain))
		case *org.DomainRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, org.UniqueOrgDomain)
		case *iam.IDPConfigAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.ConfigID, idpconfig.NewAddIDPConfigNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner))
		case *iam.IDPConfigChangedEvent:
			if e.Name == nil {
				continue
			}
			rm.changeUniqueConstraint(e.Aggregate().ID, e.ConfigID, idpconfig.NewAddIDPConfigNameUniqueConstraint(*e.Name, e.Aggregate().ResourceOwner))
		case *iam.IDPConfigRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.ConfigID, idpconfig.UniqueIDPConfigNameType)
		case *org.IDPConfigAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.ConfigID, idpconfig.NewAddIDPConfigNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner))
		case *org.IDPConfigChangedEvent:
			if e.Name == nil {
				continue
			}
			rm.changeUniqueConstraint(e.Aggregate().ID, e.ConfigID, idpconfig.NewAddIDPConfigNameUniqueConstraint(*e.Name, e.Aggregate().ResourceOwner))
		case *org.IDPConfigRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.ConfigID, idpconfig.UniqueIDPConfigNameType)
		case *iam.MailTextAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.MailTextType+e.Language, policy.NewAddMailTextUniqueConstraint(e.Aggregate().ID, e.MailTextType, e.Language))
		case *org.MailTextAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.MailTextType+e.Language, policy.NewAddMailTextUniqueConstraint(e.Aggregate().ID, e.MailTextType, e.Language))
		case *org.MailTextRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.MailTextType+e.Language, policy.UniqueMailText)
		case *project.ProjectAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, project.NewAddProjectNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner))
		case *project.ProjectChangeEvent:
			if e.Name == nil {
				continue
			}
			rm.changeUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, project.NewAddProjectNameUniqueConstraint(*e.Name, e.Aggregate().ResourceOwner))
		case *project.ProjectRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, project.UniqueProjectnameType)
			rm.listRemoveUniqueConstraint(e.Aggregate().ID, project.UniqueAppNameType)
			rm.listRemoveUniqueConstraint(e.Aggregate().ID, member.UniqueMember)
			rm.listRemoveUniqueConstraint(e.Aggregate().ID, project.UniqueRoleType)
			rm.listRemoveUniqueConstraint(e.Aggregate().ID, project.UniqueGrantType)
			rm.listRemoveUniqueConstraint(e.Aggregate().ID, project.UniqueProjectGrantMemberType)
		case *project.ApplicationAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.AppID, project.NewAddApplicationUniqueConstraint(e.Name, e.Aggregate().ID))
		case *project.ApplicationChangedEvent:
			rm.changeUniqueConstraint(e.Aggregate().ID, e.AppID, project.NewAddApplicationUniqueConstraint(e.Name, e.Aggregate().ID))
		case *project.ApplicationRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.AppID, project.UniqueAppNameType)
		case *project.GrantAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.GrantID, project.NewAddProjectGrantUniqueConstraint(e.GrantedOrgID, e.Aggregate().ID))
		case *project.GrantRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.GrantID, project.UniqueGrantType)
		case *project.GrantMemberAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.GrantID+e.UserID, project.NewAddProjectGrantMemberUniqueConstraint(e.Aggregate().ID, e.UserID, e.GrantID))
		case *project.GrantMemberRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.GrantID+e.UserID, project.UniqueProjectGrantMemberType)
		case *project.GrantMemberCascadeRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.GrantID+e.UserID, project.UniqueProjectGrantMemberType)
		case *project.RoleAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.Key, project.NewAddProjectRoleUniqueConstraint(e.Key, e.Aggregate().ID))
		case *project.RoleRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.Key, project.UniqueRoleType)
		case *user.HumanAddedEvent:
			policy, err := rm.commandProvider.getOrgIAMPolicy(rm.ctx, e.Aggregate().ResourceOwner)
			if err != nil {
				logging.Log("COMMAND-0k9Gs").WithError(err).Error("could not read policy for human added event unique constraint")
				continue
			}
			rm.addUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, user.NewAddUsernameUniqueConstraint(e.UserName, e.Aggregate().ResourceOwner, policy.UserLoginMustBeDomain))
		case *user.HumanRegisteredEvent:
			policy, err := rm.commandProvider.getOrgIAMPolicy(rm.ctx, e.Aggregate().ResourceOwner)
			if err != nil {
				logging.Log("COMMAND-m9fod").WithError(err).Error("could not read policy for human registered event unique constraint")
				continue
			}
			rm.addUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, user.NewAddUsernameUniqueConstraint(e.UserName, e.Aggregate().ResourceOwner, policy.UserLoginMustBeDomain))
		case *user.MachineAddedEvent:
			policy, err := rm.commandProvider.getOrgIAMPolicy(rm.ctx, e.Aggregate().ResourceOwner)
			if err != nil {
				logging.Log("COMMAND-2n8vs").WithError(err).Error("could not read policy for machine added event unique constraint")
				continue
			}
			rm.addUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, user.NewAddUsernameUniqueConstraint(e.UserName, e.Aggregate().ResourceOwner, policy.UserLoginMustBeDomain))
		case *user.UserRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, user.UniqueUsername)
			rm.listRemoveUniqueConstraint(e.Aggregate().ID, user.UniqueExternalIDPType)
		case *user.UsernameChangedEvent:
			policy, err := rm.commandProvider.getOrgIAMPolicy(rm.ctx, e.Aggregate().ResourceOwner)
			if err != nil {
				logging.Log("COMMAND-5n8gk").WithError(err).Error("could not read policy for username changed event unique constraint")
				continue
			}
			rm.changeUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, user.NewAddUsernameUniqueConstraint(e.UserName, e.Aggregate().ResourceOwner, policy.UserLoginMustBeDomain))
		case *user.DomainClaimedEvent:
			policy, err := rm.commandProvider.getOrgIAMPolicy(rm.ctx, e.Aggregate().ResourceOwner)
			if err != nil {
				logging.Log("COMMAND-xb8uf").WithError(err).Error("could not read policy for domain claimed event unique constraint")
				continue
			}
			rm.changeUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, user.NewAddUsernameUniqueConstraint(e.UserName, e.Aggregate().ResourceOwner, policy.UserLoginMustBeDomain))
		case *user.HumanExternalIDPAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.IDPConfigID+e.ExternalUserID, user.NewAddExternalIDPUniqueConstraint(e.IDPConfigID, e.ExternalUserID))
		case *user.HumanExternalIDPRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.IDPConfigID+e.ExternalUserID, user.UniqueExternalIDPType)
		case *user.HumanExternalIDPCascadeRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.IDPConfigID+e.ExternalUserID, user.UniqueExternalIDPType)
		case *usergrant.UserGrantAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, usergrant.NewAddUserGrantUniqueConstraint(e.Aggregate().ResourceOwner, e.UserID, e.ProjectID, e.ProjectGrantID))
		case *usergrant.UserGrantRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, usergrant.UniqueUserGrant)
		case *usergrant.UserGrantCascadeRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.Aggregate().ID, usergrant.UniqueUserGrant)
		case *iam.MemberAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.UserID, member.NewAddMemberUniqueConstraint(e.Aggregate().ID, e.UserID))
		case *iam.MemberRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.UserID, member.UniqueMember)
		case *iam.MemberCascadeRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.UserID, member.UniqueMember)
		case *org.MemberAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.UserID, member.NewAddMemberUniqueConstraint(e.Aggregate().ID, e.UserID))
		case *org.MemberRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.UserID, member.UniqueMember)
		case *org.MemberCascadeRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.UserID, member.UniqueMember)
		case *project.MemberAddedEvent:
			rm.addUniqueConstraint(e.Aggregate().ID, e.UserID, member.NewAddMemberUniqueConstraint(e.Aggregate().ID, e.UserID))
		case *project.MemberRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.UserID, member.UniqueMember)
		case *project.MemberCascadeRemovedEvent:
			rm.removeUniqueConstraint(e.Aggregate().ID, e.UserID, member.UniqueMember)
		}
	}
	return nil
}

func (rm *UniqueConstraintReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().AggregateTypes(
		iam.AggregateType,
		org.AggregateType,
		project.AggregateType,
		user.AggregateType,
		usergrant.AggregateType).
		EventTypes(
			org.OrgAddedEventType,
			org.OrgChangedEventType,
			org.OrgDomainVerifiedEventType,
			org.OrgDomainRemovedEventType,
			iam.IDPConfigAddedEventType,
			iam.IDPConfigChangedEventType,
			iam.IDPConfigRemovedEventType,
			org.IDPConfigAddedEventType,
			org.IDPConfigChangedEventType,
			org.IDPConfigRemovedEventType,
			iam.MailTextAddedEventType,
			org.MailTextAddedEventType,
			org.MailTextRemovedEventType,
			project.ProjectAddedType,
			project.ProjectChangedType,
			project.ProjectRemovedType,
			project.ApplicationAddedType,
			project.ApplicationChangedType,
			project.ApplicationRemovedType,
			project.GrantAddedType,
			project.GrantRemovedType,
			project.GrantMemberAddedType,
			project.GrantMemberRemovedType,
			project.GrantMemberCascadeRemovedType,
			project.RoleAddedType,
			project.RoleRemovedType,
			user.UserV1AddedType,
			user.UserV1RegisteredType,
			user.HumanAddedType,
			user.HumanRegisteredType,
			user.MachineAddedEventType,
			user.UserUserNameChangedType,
			user.UserDomainClaimedType,
			user.UserRemovedType,
			user.HumanExternalIDPAddedType,
			user.HumanExternalIDPRemovedType,
			user.HumanExternalIDPCascadeRemovedType,
			usergrant.UserGrantAddedType,
			usergrant.UserGrantRemovedType,
			usergrant.UserGrantCascadeRemovedType,
			iam.MemberAddedEventType,
			iam.MemberRemovedEventType,
			iam.MemberCascadeRemovedEventType,
			org.MemberAddedEventType,
			org.MemberRemovedEventType,
			org.MemberCascadeRemovedEventType,
			project.MemberAddedType,
			project.MemberRemovedType,
			project.MemberCascadeRemovedType).
		Builder()
}

func (rm *UniqueConstraintReadModel) getUniqueConstraint(aggregateID, objectID, constraintType string) *domain.UniqueConstraintMigration {
	for _, uniqueConstraint := range rm.UniqueConstraints {
		if uniqueConstraint.AggregateID == aggregateID && uniqueConstraint.ObjectID == objectID && uniqueConstraint.UniqueType == constraintType {
			return uniqueConstraint
		}
	}
	return nil
}

func (rm *UniqueConstraintReadModel) addUniqueConstraint(aggregateID, objectID string, constraint *eventstore.EventUniqueConstraint) {
	migrateUniqueConstraint := &domain.UniqueConstraintMigration{
		AggregateID:  aggregateID,
		ObjectID:     objectID,
		UniqueType:   constraint.UniqueType,
		UniqueField:  constraint.UniqueField,
		ErrorMessage: constraint.ErrorMessage,
	}
	rm.UniqueConstraints = append(rm.UniqueConstraints, migrateUniqueConstraint)
}

func (rm *UniqueConstraintReadModel) changeUniqueConstraint(aggregateID, objectID string, constraint *eventstore.EventUniqueConstraint) {
	for i, uniqueConstraint := range rm.UniqueConstraints {
		if uniqueConstraint.AggregateID == aggregateID && uniqueConstraint.ObjectID == objectID && uniqueConstraint.UniqueType == constraint.UniqueType {
			rm.UniqueConstraints[i] = &domain.UniqueConstraintMigration{
				AggregateID:  aggregateID,
				ObjectID:     objectID,
				UniqueType:   constraint.UniqueType,
				UniqueField:  constraint.UniqueField,
				ErrorMessage: constraint.ErrorMessage,
			}
			return
		}
	}
}

func (rm *UniqueConstraintReadModel) removeUniqueConstraint(aggregateID, objectID, constraintType string) {
	for i, uniqueConstraint := range rm.UniqueConstraints {
		if uniqueConstraint.AggregateID == aggregateID && uniqueConstraint.ObjectID == objectID && uniqueConstraint.UniqueType == constraintType {
			copy(rm.UniqueConstraints[i:], rm.UniqueConstraints[i+1:])
			rm.UniqueConstraints[len(rm.UniqueConstraints)-1] = nil
			rm.UniqueConstraints = rm.UniqueConstraints[:len(rm.UniqueConstraints)-1]
			return
		}
	}
}

func (rm *UniqueConstraintReadModel) listRemoveUniqueConstraint(aggregateID, constraintType string) {
	for i := len(rm.UniqueConstraints) - 1; i >= 0; i-- {
		if rm.UniqueConstraints[i].AggregateID == aggregateID && rm.UniqueConstraints[i].UniqueType == constraintType {
			copy(rm.UniqueConstraints[i:], rm.UniqueConstraints[i+1:])
			rm.UniqueConstraints[len(rm.UniqueConstraints)-1] = nil
			rm.UniqueConstraints = rm.UniqueConstraints[:len(rm.UniqueConstraints)-1]
		}
	}
}
