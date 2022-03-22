package migration

import "github.com/caos/zitadel/internal/eventstore"

//SetupStep is the command pushed on the eventstore
type SetupStep struct {
	typ       eventstore.EventType
	migration Migration
	Name      string `json:"name"`
	done      bool
}

func setupStartedCmd(migration Migration) eventstore.Command {
	return &SetupStep{
		migration: migration,
		typ:       startedType,
		Name:      migration.String(),
	}
}

func setupDoneCmd(migration Migration, err error) eventstore.Command {
	s := &SetupStep{
		typ:       doneType,
		migration: migration,
		Name:      migration.String(),
	}

	if err != nil {
		s.typ = failedType
	}

	return s
}

func (s *SetupStep) Aggregate() eventstore.Aggregate {
	return eventstore.Aggregate{
		ID:            aggregateID,
		Type:          aggregateType,
		ResourceOwner: "SYSTEM",
		Version:       "v1",
	}
}

func (s *SetupStep) EditorService() string {
	return "system"
}

func (s *SetupStep) EditorUser() string {
	return "system"
}

func (s *SetupStep) Type() eventstore.EventType {
	return s.typ
}

func (s *SetupStep) Data() interface{} {
	return s
}

func (s *SetupStep) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	switch s.typ {
	case startedType:
		return []*eventstore.EventUniqueConstraint{
			eventstore.NewAddEventUniqueConstraint("migration_started", s.migration.String(), "Errors.Step.Started.AlreadyExists"),
		}
	case failedType:
		return []*eventstore.EventUniqueConstraint{
			eventstore.NewRemoveEventUniqueConstraint("migration_started", s.migration.String()),
		}
	default:
		return []*eventstore.EventUniqueConstraint{
			eventstore.NewAddEventUniqueConstraint("migration_done", s.migration.String(), "Errors.Step.Done.AlreadyExists"),
		}
	}
}
