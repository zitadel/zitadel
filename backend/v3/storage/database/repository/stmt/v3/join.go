package v3

type join struct {
	table      Table
	conditions []joinCondition
}

type joinCondition struct {
	left  Column
	right Column
}
