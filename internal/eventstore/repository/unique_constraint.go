package repository

//UniqueCheck represents all information about a unique attribute
type UniqueConstraint struct {
	//UniqueField is the field which should be unique
	UniqueField string

	//UniqueType is the type of the unique field
	UniqueType string

	//Action defines if unique constraint should be added or removed
	Action UniqueConstraintAction

	//ErrorMessage is the message key which should be returned if constraint is violated
	ErrorMessage string
}

type UniqueConstraintAction int32

const (
	UniqueConstraintAdd UniqueConstraintAction = iota
	UniqueConstraintRemoved

	uniqueConstraintActionCount
)

func (f UniqueConstraintAction) Valid() bool {
	return f >= 0 && f < uniqueConstraintActionCount
}
