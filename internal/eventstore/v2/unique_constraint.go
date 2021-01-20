package eventstore

type EventUniqueConstraint struct {
	// TableName is the table name for the unique constraint
	TableName string
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

	uniqueConstraintActionCount
)

func NewAddEventUniqueConstraint(
	tableName,
	uniqueField,
	errMessage string) *EventUniqueConstraint {
	return &EventUniqueConstraint{
		TableName:    tableName,
		UniqueField:  uniqueField,
		ErrorMessage: errMessage,
		Action:       UniqueConstraintAdd,
	}
}

func NewRemoveEventUniqueConstraint(
	tableName,
	uniqueField string) *EventUniqueConstraint {
	return &EventUniqueConstraint{
		TableName:   tableName,
		UniqueField: uniqueField,
		Action:      UniqueConstraintRemove,
	}
}
