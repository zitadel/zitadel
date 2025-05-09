package v4

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userMachine struct {
	*user
}

var _ domain.MachineRepository = (*userMachine)(nil)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

// Update implements [domain.MachineRepository].
func (m userMachine) Update(ctx context.Context, condition database.Condition, changes ...database.Change) (err error) {
	m.builder.WriteString("UPDATE user_machines SET ")
	database.Changes(changes).Write(&m.builder)
	m.writeCondition(condition)
	m.writeReturning()

	return m.client.Exec(ctx, m.builder.String(), m.builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetDescription implements [domain.machineChanges].
func (m userMachine) SetDescription(description string) database.Change {
	return database.NewChange(m.DescriptionColumn(), description)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// DescriptionCondition implements [domain.machineConditions].
func (m userMachine) DescriptionCondition(op database.TextOperation, description string) database.Condition {
	return database.NewTextCondition(m.DescriptionColumn(), op, description)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// DescriptionColumn implements [domain.machineColumns].
func (m userMachine) DescriptionColumn() database.Column {
	return database.NewColumn("description")
}

func (m userMachine) columns() database.Columns {
	return append(m.user.columns(), m.DescriptionColumn())
}

func (m *userMachine) writeReturning() {
	m.builder.WriteString(" RETURNING ")
	m.columns().Write(&m.builder)
}
