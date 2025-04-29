package v4

import "context"

type Machine struct {
	Description string
}

func (Machine) userTrait() {}

func (m Machine) Type() UserType {
	return UserTypeMachine
}

const UserTypeMachine UserType = "machine"

var _ userTrait = (*Machine)(nil)

type userMachine struct {
	*user
}

func (u *user) Machine() *userMachine {
	return &userMachine{user: u}
}

func (m userMachine) Update(ctx context.Context, cols ...Change) (*Machine, error) {
	return nil, nil
}

func (userMachine) DescriptionColumn() Column {
	return column{"m.description"}
}

func (m userMachine) SetDescription(description string) Change {
	return newChange(m.DescriptionColumn(), description)
}

func (m userMachine) DescriptionCondition(op TextOperator, description string) Condition {
	return newTextCondition(m.DescriptionColumn(), op, description)
}
