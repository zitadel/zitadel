package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/org"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

type UniqueConstraintReadModel struct {
	eventstore.WriteModel

	UniqueConstraints []*domain.UniqueConstraintMigration
}

func NewUniqueConstraintReadModel() *UniqueConstraintReadModel {
	return &UniqueConstraintReadModel{}
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
		case *org.DomainVerificationAddedEvent:
			rm.addUniqueConstraint(e.AggregateID(), e.AggregateID(), org.NewAddOrgNameUniqueConstraint(e.Domain))
		case *org.DomainRemovedEvent:
			rm.removeUniqueConstraint(e.AggregateID(), e.AggregateID())
		case *iam.IDPConfigAddedEvent:

		case *iam.IDPConfigChangedEvent:

		case *iam.IDPConfigRemovedEvent:

		case *org.IDPConfigAddedEvent:

		case *org.IDPConfigChangedEvent:

		case *org.IDPConfigRemovedEvent:

		case *iam.MailTextAddedEvent:

		case *org.MailTextAddedEvent:

		case *org.MailTextRemovedEvent:

		case *project.ProjectAddedEvent:

		case *project.ProjectChangeEvent:

		case *project.ProjectRemovedEvent:

		case *project.ApplicationAddedEvent:

		case *project.ApplicationChangedEvent:

		case *project.ApplicationRemovedEvent:

		case *project.GrantAddedEvent:
		case *project.GrantRemovedEvent:
		case *project.GrantMemberAddedEvent:
		case *project.GrantMemberRemovedEvent:
		case *project.RoleAddedEvent:
		case *project.RoleRemovedEvent:
		}
	}
	return nil
}

func (rm *UniqueConstraintReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AggregateIDs(rm.AggregateID).
		ResourceOwner(rm.ResourceOwner).
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
		)
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
		if uniqueConstraint.AggregateID == aggregateID && uniqueConstraint.ObjectID == objectID {
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

func (rm *UniqueConstraintReadModel) removeUniqueConstraint(aggregateID, objectID string) {
	for i, uniqueConstraint := range rm.UniqueConstraints {
		if uniqueConstraint.AggregateID == aggregateID && uniqueConstraint.ObjectID == objectID {
			copy(rm.UniqueConstraints[i:], rm.UniqueConstraints[i+1:])
			rm.UniqueConstraints[len(rm.UniqueConstraints)-1] = nil
			rm.UniqueConstraints = rm.UniqueConstraints[:len(rm.UniqueConstraints)-1]
			return
		}
	}
}
