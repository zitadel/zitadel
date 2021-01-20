package repository

//UniqueCheck represents all information about a unique attribute
type UniqueConstraint struct {
	//UniqueField is the field which should be unique
	UniqueField string

	//TableName is the table name for the unique field
	TableName string

	//Type describes the cause of the event (e.g. user.added)
	// it should always be in past-form
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
