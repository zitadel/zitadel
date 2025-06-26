package repository

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
func (m userMachine) Update(ctx context.Context, condition database.Condition, changes ...database.Change) error {
	builder := database.StatementBuilder{}
	builder.WriteString("UPDATE user_machines SET ")
	database.Changes(changes).Write(&builder)
	m.writeCondition(&builder, condition)
	m.writeReturning()

	_, err := m.client.Exec(ctx, builder.String(), builder.Args()...)
	return err
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
	builder := database.StatementBuilder{}
	builder.WriteString(" RETURNING ")
	m.columns().Write(&builder)
}
