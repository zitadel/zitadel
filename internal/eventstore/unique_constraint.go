package eventstore

type EventUniqueConstraint struct {
	// UniqueType is the table name for the unique constraint
	UniqueType string
	//UniqueField is the unique key
	UniqueField string
	//Action defines if unique constraint should be added or removed
	Action UniqueConstraintAction
	//ErrorMessage defines the translation file key for the error message
	ErrorMessage string
}

type UniqueConstraintAction int32

const (
	UniqueConstraintAdd UniqueConstraintAction = iota
	UniqueConstraintRemove
)

func NewAddEventUniqueConstraint(
	uniqueType,
	uniqueField,
	errMessage string) *EventUniqueConstraint {
	return &EventUniqueConstraint{
		UniqueType:   uniqueType,
		UniqueField:  uniqueField,
		ErrorMessage: errMessage,
		Action:       UniqueConstraintAdd,
	}
}

func NewRemoveEventUniqueConstraint(
	uniqueType,
	uniqueField string) *EventUniqueConstraint {
	return &EventUniqueConstraint{
		UniqueType:  uniqueType,
		UniqueField: uniqueField,
		Action:      UniqueConstraintRemove,
	}
}
