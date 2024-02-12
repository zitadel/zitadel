package execution

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	UniqueExecution    = "execution"
	DuplicateExecution = "Errors.Execution.AlreadyExists"
)

func NewAddUniqueConstraints(name string) []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		eventstore.NewAddEventUniqueConstraint(
			UniqueExecution,
			name,
			DuplicateExecution,
		),
	}
}

func NewRemoveUniqueConstraints(name string) []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		eventstore.NewRemoveUniqueConstraint(
			UniqueExecution,
			name,
		),
	}
}
