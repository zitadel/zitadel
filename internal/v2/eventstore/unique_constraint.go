package eventstore

type UniqueConstraint struct {
	// UniqueType is the table name for the unique constraint
	UniqueType string
	// UniqueField is the unique key
	UniqueField string
	// Action defines if unique constraint should be added or removed
	Action UniqueConstraintAction
	// ErrorMessage defines the translation file key for the error message
	ErrorMessage string
	// IsGlobal defines if the unique constraint is globally unique or just within a single instance
	IsGlobal bool
}

type UniqueConstraintAction int8

const (
	UniqueConstraintAdd UniqueConstraintAction = iota
	UniqueConstraintRemove
	UniqueConstraintInstanceRemove

	uniqueConstraintActionCount
)

func (f UniqueConstraintAction) Valid() bool {
	return f >= 0 && f < uniqueConstraintActionCount
}

func NewAddEventUniqueConstraint(
	uniqueType,
	uniqueField,
	errMessage string) *UniqueConstraint {
	return &UniqueConstraint{
		UniqueType:   uniqueType,
		UniqueField:  uniqueField,
		ErrorMessage: errMessage,
		Action:       UniqueConstraintAdd,
	}
}

func NewRemoveUniqueConstraint(
	uniqueType,
	uniqueField string) *UniqueConstraint {
	return &UniqueConstraint{
		UniqueType:  uniqueType,
		UniqueField: uniqueField,
		Action:      UniqueConstraintRemove,
	}
}

func NewRemoveInstanceUniqueConstraints() *UniqueConstraint {
	return &UniqueConstraint{
		Action: UniqueConstraintInstanceRemove,
	}
}

func NewAddGlobalUniqueConstraint(
	uniqueType,
	uniqueField,
	errMessage string) *UniqueConstraint {
	return &UniqueConstraint{
		UniqueType:   uniqueType,
		UniqueField:  uniqueField,
		ErrorMessage: errMessage,
		IsGlobal:     true,
		Action:       UniqueConstraintAdd,
	}
}

func NewRemoveGlobalUniqueConstraint(
	uniqueType,
	uniqueField string) *UniqueConstraint {
	return &UniqueConstraint{
		UniqueType:  uniqueType,
		UniqueField: uniqueField,
		IsGlobal:    true,
		Action:      UniqueConstraintRemove,
	}
}
