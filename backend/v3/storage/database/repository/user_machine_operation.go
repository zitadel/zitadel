package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
)

type machineOperation struct {
	userOperation
}

// SetDescription implements domain.MachineOperation.
func (m *machineOperation) SetDescription(ctx context.Context, description string) error {
	return m.QueryExecutor.Exec(ctx, `UPDATE machines SET description = $1 WHERE id = $2`, description, m.clauses)
}

var _ domain.MachineOperation = (*machineOperation)(nil)
