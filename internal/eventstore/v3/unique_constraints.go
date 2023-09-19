package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed unique_constraints_delete.sql
	deleteConstraintStmt string
	//go:embed unique_constraints_add.sql
	addConstraintStmt string
)

func handleUniqueConstraints(ctx context.Context, tx *sql.Tx, commands []eventstore.Command) error {
	deletePlaceholders := make([]string, 0)
	deleteArgs := make([]any, 0)

	addPlaceholders := make([]string, 0)
	addArgs := make([]any, 0)

	for _, command := range commands {
		for _, constraint := range command.UniqueConstraints() {
			switch constraint.Action {
			case eventstore.UniqueConstraintAdd:
				addPlaceholders = append(addPlaceholders, fmt.Sprintf("($%d, $%d, $%d)", len(addArgs)+1, len(addArgs)+2, len(addArgs)+3))
				addArgs = append(addArgs, command.Aggregate().InstanceID, constraint.UniqueType, constraint.UniqueField)
			case eventstore.UniqueConstraintRemove:
				deletePlaceholders = append(deletePlaceholders, fmt.Sprintf("(instance_id = $%d AND unique_type = $%d AND unique_field = $%d)", len(deleteArgs)+1, len(deleteArgs)+2, len(deleteArgs)+3))
				deleteArgs = append(deleteArgs, command.Aggregate().InstanceID, constraint.UniqueType, constraint.UniqueField)
			case eventstore.UniqueConstraintInstanceRemove:
				deletePlaceholders = append(deletePlaceholders, fmt.Sprintf("(instance_id = $%d)", len(deleteArgs)+1))
				deleteArgs = append(deleteArgs, command.Aggregate().InstanceID)
			}
		}
	}

	if len(deletePlaceholders) > 0 {
		_, err := tx.ExecContext(ctx, fmt.Sprintf(deleteConstraintStmt, strings.Join(deletePlaceholders, " OR ")), deleteArgs...)
		if err != nil {
			logging.WithError(err).Warn("delete unique constraint failed")
			return errors.ThrowInternal(err, "V3-C8l3V", "Errors.Internal")
		}
	}
	if len(addPlaceholders) > 0 {
		_, err := tx.ExecContext(ctx, fmt.Sprintf(addConstraintStmt, strings.Join(addPlaceholders, ", ")), addArgs...)
		if err != nil {
			logging.WithError(err).Warn("add unique constraint failed")
			return errors.ThrowInternal(err, "V3-DKcYh", "Errors.Internal")
		}
	}
	return nil
}
