package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"
)

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
)

var (
	//go:embed unique_constraints_delete.sql
	deleteConstraintStmt string
	//go:embed unique_constraints_add.sql
	addConstraintStmt string
)

func handleUniqueConstraints(ctx context.Context, tx *sql.Tx, commands []Command) error {
	deletePlaceholders := make([]string, 0)
	deleteArgs := make([]any, 0)

	addPlaceholders := make([]string, 0)
	addArgs := make([]any, 0)

	for _, command := range commands {
		for _, constraint := range command.UniqueConstraints() {
			switch constraint.Action {
			case UniqueConstraintAdd:
				addPlaceholders = append(addPlaceholders, fmt.Sprintf("($%d, $%d, $%d)", len(addArgs)+1, len(addArgs)+2, len(addArgs)+3))
				addArgs = append(addArgs, command.Aggregate().InstanceID, constraint.UniqueType, constraint.UniqueField)
			case UniqueConstraintRemove:
				deletePlaceholders = append(deletePlaceholders, fmt.Sprintf("(instance_id = $%d AND unique_type = $%d AND unique_field = $%d)", len(deleteArgs)+1, len(deleteArgs)+2, len(deleteArgs)+3))
				deleteArgs = append(deleteArgs, command.Aggregate().InstanceID, constraint.UniqueType, constraint.UniqueField)
			case UniqueConstraintInstanceRemove:
				deletePlaceholders = append(deletePlaceholders, fmt.Sprintf("(instance_id = $%d)", len(deleteArgs)+1))
				deleteArgs = append(deleteArgs, command.Aggregate().InstanceID)
			}
		}
	}

	if len(deletePlaceholders) > 0 {
		_, err := tx.ExecContext(ctx, fmt.Sprintf(deleteConstraintStmt, strings.Join(deletePlaceholders, " OR ")), deleteArgs...)
		if err != nil {
			return err
		}
	}
	if len(addPlaceholders) > 0 {
		_, err := tx.ExecContext(ctx, fmt.Sprintf(addConstraintStmt, strings.Join(addPlaceholders, ", ")), addArgs...)
		return err
	}
	return nil
}
