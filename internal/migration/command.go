package migration

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/service"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// SetupStep is the command pushed on the eventstore
type SetupStep struct {
	eventstore.BaseEvent `json:"-"`
	migration            Migration
	Name                 string `json:"name"`
	Error                any    `json:"error,omitempty"`
	LastRun              any    `json:"lastRun,omitempty"`
}

func setupStartedCmd(ctx context.Context, migration Migration) eventstore.Command {
	ctx = authz.SetCtxData(service.WithService(ctx, "system"), authz.CtxData{UserID: "system", OrgID: "SYSTEM", ResourceOwner: "SYSTEM"})
	return &SetupStep{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			eventstore.NewAggregate(ctx, aggregateID, aggregateType, "v1"),
			StartedType),
		migration: migration,
		Name:      migration.String(),
	}
}

func setupDoneCmd(ctx context.Context, migration Migration, err error) eventstore.Command {
	ctx = authz.SetCtxData(service.WithService(ctx, "system"), authz.CtxData{UserID: "system", OrgID: "SYSTEM", ResourceOwner: "SYSTEM"})
	typ := doneType
	var lastRun interface{}
	if repeatable, ok := migration.(RepeatableMigration); ok {
		typ = repeatableDoneType
		lastRun = repeatable
	}

	s := &SetupStep{
		migration: migration,
		Name:      migration.String(),
		LastRun:   lastRun,
	}
	if err != nil {
		typ = failedType
		s.Error = err.Error()
	}

	s.BaseEvent = *eventstore.NewBaseEventForPush(
		ctx,
		eventstore.NewAggregate(ctx, aggregateID, aggregateType, "v1"),
		typ)

	return s
}

func (s *SetupStep) Payload() interface{} {
	return s
}

func (s *SetupStep) UniqueConstraints() []*eventstore.UniqueConstraint {
	switch s.Type() {
	case StartedType:
		return []*eventstore.UniqueConstraint{
			eventstore.NewAddGlobalUniqueConstraint("migration_started", s.migration.String(), "Errors.Step.Started.AlreadyExists"),
		}
	case failedType,
		repeatableDoneType:
		return []*eventstore.UniqueConstraint{
			eventstore.NewRemoveGlobalUniqueConstraint("migration_started", s.migration.String()),
		}
	default:
		return []*eventstore.UniqueConstraint{
			eventstore.NewAddGlobalUniqueConstraint("migration_done", s.migration.String(), "Errors.Step.Done.AlreadyExists"),
		}
	}
}

func RegisterMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(aggregateType, StartedType, SetupMapper)
	es.RegisterFilterEventMapper(aggregateType, doneType, SetupMapper)
	es.RegisterFilterEventMapper(aggregateType, failedType, SetupMapper)
	es.RegisterFilterEventMapper(aggregateType, repeatableDoneType, SetupMapper)
}

func SetupMapper(event eventstore.Event) (eventstore.Event, error) {
	step := &SetupStep{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(step)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-hYp7M", "unable to unmarshal step")
	}

	return step, nil
}
