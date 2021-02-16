package iam

import (
	"context"
	"github.com/caos/zitadel/internal/v2/domain"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
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

func NewMigrateUniqueConstraintEvent(ctx context.Context, uniqueConstraintMigrations []*domain.UniqueConstraintMigration) *MigrateUniqueConstraintEvent {
	return &MigrateUniqueConstraintEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
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
