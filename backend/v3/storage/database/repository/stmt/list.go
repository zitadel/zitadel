package stmt

type ListEntry interface {
	Number | Text | any
}

type ListCondition[E ListEntry, T any] struct {
	condition[[]E, T, ListOperation]
}

func (lc *ListCondition[E, T]) Apply(stmt *statement[T]) {
	placeholder := stmt.appendArg(lc.value)

	switch lc.op {
	case ListOperationEqual, ListOperationNotEqual:
		lc.field.Apply(stmt)
		operation[T, ListOperation]{lc.op}.Apply(stmt)
		stmt.builder.WriteString(placeholder)
	case ListOperationContainsAny, ListOperationContainsAll:
		lc.field.Apply(stmt)
		operation[T, ListOperation]{lc.op}.Apply(stmt)
		stmt.builder.WriteString(placeholder)
	case ListOperationNotContainsAny, ListOperationNotContainsAll:
		stmt.builder.WriteString("NOT (")
		lc.field.Apply(stmt)
		operation[T, ListOperation]{lc.op}.Apply(stmt)
		stmt.builder.WriteString(placeholder)
		stmt.builder.WriteString(")")
	default:
		panic("unknown list operation")
	}
}

type ListOperation uint8

const (
	// ListOperationEqual checks if the arrays are equal including the order of the elements
	ListOperationEqual ListOperation = iota + 1
	// ListOperationNotEqual checks if the arrays are not equal including the order of the elements
	ListOperationNotEqual

	// ListOperationContains checks if the array column contains all the values of the specified array
	ListOperationContainsAll
	// ListOperationContainsAny checks if the arrays have at least one value in common
	ListOperationContainsAny
	// ListOperationContainsAll checks if the array column contains all the values of the specified array

	// ListOperationNotContainsAll checks if the specified array is not contained by the column
	ListOperationNotContainsAll
	// ListOperationNotContainsAny checks if the arrays column contains none of the values of the specified array
	ListOperationNotContainsAny
)

var listOperations = map[ListOperation]string{
	// ListOperationEqual checks if the lists are equal
	ListOperationEqual: " = ",
	// ListOperationNotEqual checks if the lists are not equal
	ListOperationNotEqual: " <> ",
	// ListOperationContainsAny checks if the arrays have at least one value in common
	ListOperationContainsAny: " && ",
	// ListOperationContainsAll checks if the array column contains all the values of the specified array
	ListOperationContainsAll: " @> ",
	// ListOperationNotContainsAny checks if the arrays column contains none of the values of the specified array
	ListOperationNotContainsAny: " && ", // Base operator for NOT (A && B)
	// ListOperationNotContainsAll checks if the array column is not contained by the specified array
	ListOperationNotContainsAll: " <@ ", // Base operator for NOT (A <@ B)
}

func (lo ListOperation) String() string {
	return listOperations[lo]
}
