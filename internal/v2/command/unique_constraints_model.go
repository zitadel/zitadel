package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/idpconfig"
	"github.com/caos/zitadel/internal/v2/repository/member"
	"github.com/caos/zitadel/internal/v2/repository/org"
	"github.com/caos/zitadel/internal/v2/repository/policy"
	"github.com/caos/zitadel/internal/v2/repository/project"
	"github.com/caos/zitadel/internal/v2/repository/user"
	"github.com/caos/zitadel/internal/v2/repository/usergrant"
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
			rm.addUniqueConstraint(e.AggregateID(), e.AggregateID(), org.NewAddOrgNameUniqueConstraint(e.Name))
		case *org.OrgChangedEvent:
			rm.changeUniqueConstraint(e.AggregateID(), e.AggregateID(), org.NewAddOrgNameUniqueConstraint(e.Name))
		case *org.DomainVerifiedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.AggregateID(), org.NewAddOrgNameUniqueConstraint(e.Domain))
		case *org.DomainRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.AggregateID(), org.UniqueOrgDomain)
		case *iam.IDPConfigAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.ConfigID, idpconfig.NewAddIDPConfigNameUniqueConstraint(e.Name, e.ResourceOwner()))
		case *iam.IDPConfigChangedEvent:
			rm.changeUniqueConstraint(e.AggregateID(), e.ConfigID, idpconfig.NewAddIDPConfigNameUniqueConstraint(*e.Name, e.ResourceOwner()))
		case *iam.IDPConfigRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.ConfigID, idpconfig.UniqueIDPConfigNameType)
		case *org.IDPConfigAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.ConfigID, idpconfig.NewAddIDPConfigNameUniqueConstraint(e.Name, e.ResourceOwner()))
		case *org.IDPConfigChangedEvent:
			rm.changeUniqueConstraint(e.AggregateID(), e.ConfigID, idpconfig.NewAddIDPConfigNameUniqueConstraint(*e.Name, e.ResourceOwner()))
		case *org.IDPConfigRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.ConfigID, idpconfig.UniqueIDPConfigNameType)
		case *iam.MailTextAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.MailTextType+e.Language, policy.NewAddMailTextUniqueConstraint(e.AggregateID(), e.MailTextType, e.Language))
		case *org.MailTextAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.MailTextType+e.Language, policy.NewAddMailTextUniqueConstraint(e.AggregateID(), e.MailTextType, e.Language))
		case *org.MailTextRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.MailTextType+e.Language, policy.UniqueMailText)
		case *project.ProjectAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.AggregateID(), project.NewAddProjectNameUniqueConstraint(e.Name, e.ResourceOwner()))
		case *project.ProjectChangeEvent:
			rm.changeUniqueConstraint(e.AggregateID(), e.AggregateID(), project.NewAddProjectNameUniqueConstraint(*e.Name, e.ResourceOwner()))
		case *project.ProjectRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.AggregateID(), project.UniqueProjectnameType)
		case *project.ApplicationAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.AppID, project.NewAddApplicationUniqueConstraint(e.Name, e.AggregateID()))
		case *project.ApplicationChangedEvent:
			rm.changeUniqueConstraint(e.AggregateID(), e.AppID, project.NewAddApplicationUniqueConstraint(e.Name, e.AggregateID()))
		case *project.ApplicationRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.AppID, project.UniqueAppNameType)
		case *project.GrantAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.GrantID, project.NewAddProjectGrantUniqueConstraint(e.GrantedOrgID, e.AggregateID()))
		case *project.GrantRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.GrantID, project.UniqueGrantType)
		case *project.GrantMemberAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.GrantID+e.UserID, project.NewAddProjectGrantMemberUniqueConstraint(e.AggregateID(), e.UserID, e.GrantID))
		case *project.GrantMemberRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.GrantID+e.UserID, project.UniqueProjectGrantMemberType)
		case *project.RoleAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.Key, project.NewAddProjectRoleUniqueConstraint(e.Key, e.AggregateID()))
		case *project.RoleRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.Key, project.UniqueRoleType)
		case *user.HumanAddedEvent:
			policy, err := rm.commandProvider.getOrgIAMPolicy(rm.ctx, e.ResourceOwner())
			if err != nil {
				logging.Log("COMMAND-0k9Gs").WithError(err).Error("could not read policy for human added event unique constraint")
				continue
			}
			rm.addUniqueConstraint(e.AggregateID(), e.AggregateID(), user.NewAddUsernameUniqueConstraint(e.UserName, e.ResourceOwner(), policy.UserLoginMustBeDomain))
		case *user.HumanRegisteredEvent:
			policy, err := rm.commandProvider.getOrgIAMPolicy(rm.ctx, e.ResourceOwner())
			if err != nil {
				logging.Log("COMMAND-m9fod").WithError(err).Error("could not read policy for human registered event unique constraint")
				continue
			}
			rm.addUniqueConstraint(e.AggregateID(), e.AggregateID(), user.NewAddUsernameUniqueConstraint(e.UserName, e.ResourceOwner(), policy.UserLoginMustBeDomain))
		case *user.MachineAddedEvent:
			policy, err := rm.commandProvider.getOrgIAMPolicy(rm.ctx, e.ResourceOwner())
			if err != nil {
				logging.Log("COMMAND-2n8vs").WithError(err).Error("could not read policy for machine added event unique constraint")
				continue
			}
			rm.addUniqueConstraint(e.AggregateID(), e.AggregateID(), user.NewAddUsernameUniqueConstraint(e.UserName, e.ResourceOwner(), policy.UserLoginMustBeDomain))
		case *user.UserRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.AggregateID(), user.UniqueUsername)
		case *user.UsernameChangedEvent:
			policy, err := rm.commandProvider.getOrgIAMPolicy(rm.ctx, e.ResourceOwner())
			if err != nil {
				logging.Log("COMMAND-5n8gk").WithError(err).Error("could not read policy for username changed event unique constraint")
				continue
			}
			rm.changeUniqueConstraint(e.AggregateID(), e.AggregateID(), user.NewAddUsernameUniqueConstraint(e.UserName, e.ResourceOwner(), policy.UserLoginMustBeDomain))
		case *user.DomainClaimedEvent:
			policy, err := rm.commandProvider.getOrgIAMPolicy(rm.ctx, e.ResourceOwner())
			if err != nil {
				logging.Log("COMMAND-xb8uf").WithError(err).Error("could not read policy for domain claimed event unique constraint")
				continue
			}
			rm.changeUniqueConstraint(e.AggregateID(), e.AggregateID(), user.NewAddUsernameUniqueConstraint(e.UserName, e.ResourceOwner(), policy.UserLoginMustBeDomain))
		case *user.HumanExternalIDPAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.IDPConfigID+e.ExternalUserID, user.NewAddExternalIDPUniqueConstraint(e.IDPConfigID, e.ExternalUserID))
		case *user.HumanExternalIDPRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.IDPConfigID+e.ExternalUserID, user.UniqueExternalIDPType)
		case *user.HumanExternalIDPCascadeRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.IDPConfigID+e.ExternalUserID, user.UniqueExternalIDPType)
		case *usergrant.UserGrantAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.AggregateID(), usergrant.NewAddUserGrantUniqueConstraint(e.ResourceOwner(), e.UserID, e.ProjectID, e.ProjectGrantID))
		case *usergrant.UserGrantRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.AggregateID(), usergrant.UniqueUserGrant)
		case *usergrant.UserGrantCascadeRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.AggregateID(), usergrant.UniqueUserGrant)
		case *iam.MemberAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.UserID, member.NewAddMemberUniqueConstraint(e.AggregateID(), e.UserID))
		case *iam.MemberRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.UserID, member.UniqueMember)
		case *org.MemberAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.UserID, member.NewAddMemberUniqueConstraint(e.AggregateID(), e.UserID))
		case *org.MemberRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.UserID, member.UniqueMember)
		case *project.MemberAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.UserID, member.NewAddMemberUniqueConstraint(e.AggregateID(), e.UserID))
		case *project.MemberRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.UserID, member.UniqueMember)
		}
	}
	return nil
}

func (rm *UniqueConstraintReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType, org.AggregateType, project.AggregateType, user.AggregateType, usergrant.AggregateType).
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
			project.RoleAddedType,
			project.RoleRemovedType,
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
			org.MemberAddedEventType,
			org.MemberRemovedEventType,
			project.MemberAddedType,
			project.MemberRemovedType,
		)
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
