package eventstore

type EventUniqueConstraint interface {
	// TableName is the table name for the unique constraint
	TableName() string
	//UniqueField is the unique key
	UniqueField() string
	//Action defines if unique constraint should be added or removed
	Action() UniqueConstraintAction
}

type UniqueConstraintAction int32

const (
	UniqueConstraintAdd UniqueConstraintAction = iota
	UniqueConstraintRemoved

	uniqueConstraintActionCount
)
