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

func (m userMachine) Update(ctx context.Context, condition Condition, changes ...Change) ([]*Machine, error) {
	m.builder.WriteString("UPDATE user_machines SET ")
	Changes(changes).writeTo(&m.builder)
	m.writeCondition(condition)
	m.writeReturning()

	var machines []*Machine
	rows, err := m.client.Query(ctx, m.builder.String(), m.builder.args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		machine := new(Machine)
		if err := rows.Scan(&machine.Description); err != nil {
			return nil, err
		}
		machines = append(machines, machine)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return machines, nil
}

func (userMachine) DescriptionColumn() Column {
	return column{"description"}
}

func (m userMachine) SetDescription(description string) Change {
	return newChange(m.DescriptionColumn(), description)
}

func (m userMachine) DescriptionCondition(op TextOperator, description string) Condition {
	return newTextCondition(m.DescriptionColumn(), op, description)
}

func (m userMachine) columns() Columns {
	return append(m.user.columns(), m.DescriptionColumn())
}

func (m *userMachine) writeReturning() {
	m.builder.WriteString(" RETURNING ")
	m.columns().writeTo(&m.builder)
}
