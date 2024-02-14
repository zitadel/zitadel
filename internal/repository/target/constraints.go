package target

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	UniqueTarget    = "target"
	DuplicateTarget = "Errors.Target.AlreadyExists"
)

func NewAddUniqueConstraints(name string) []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		eventstore.NewAddEventUniqueConstraint(
			UniqueTarget,
			name,
			DuplicateTarget,
		),
	}
}

func NewRemoveUniqueConstraints(name string) []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		eventstore.NewRemoveUniqueConstraint(
			UniqueTarget,
			name,
		),
	}
}

func NewUpdateUniqueConstraints(oldName, name string) []*eventstore.UniqueConstraint {
	return append(
		append(
			[]*eventstore.UniqueConstraint{},
			NewRemoveUniqueConstraints(oldName)...,
		),
		NewAddUniqueConstraints(name)...,
	)
}
