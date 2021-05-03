package iam

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	UniqueConstraintsMigratedEventType eventstore.EventType = "iam.unique.constraints.migrated"
)

type MigrateUniqueConstraintEvent struct {
	eventstore.BaseEvent `json:"-"`

	uniqueConstraintMigrations []*domain.UniqueConstraintMigration `json:"-"`
}

func NewAddMigrateUniqueConstraint(uniqueMigration *domain.UniqueConstraintMigration) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueMigration.UniqueType,
		uniqueMigration.UniqueField,
		uniqueMigration.ErrorMessage)
}

func (e *MigrateUniqueConstraintEvent) Data() interface{} {
	return nil
}

func (e *MigrateUniqueConstraintEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	constraints := make([]*eventstore.EventUniqueConstraint, len(e.uniqueConstraintMigrations))
	for i, uniqueMigration := range e.uniqueConstraintMigrations {
		constraints[i] = NewAddMigrateUniqueConstraint(uniqueMigration)
	}
	return constraints
}

func NewMigrateUniqueConstraintEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	uniqueConstraintMigrations []*domain.UniqueConstraintMigration) *MigrateUniqueConstraintEvent {
	return &MigrateUniqueConstraintEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UniqueConstraintsMigratedEventType,
		),
		uniqueConstraintMigrations: uniqueConstraintMigrations,
	}
}

func MigrateUniqueConstraintEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &MigrateUniqueConstraintEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
