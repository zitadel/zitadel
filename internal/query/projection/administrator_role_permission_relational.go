package projection

import (
	"context"
	"database/sql"

	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/permission"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (p *relationalTablesProjection) reduceAdministratorRolePermissionAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*permission.AddedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-v3ARP01", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AdministratorRoleRepository()
		_, err := repo.AddPermissions(ctx, v3_sql.SQLTx(tx), e.Aggregate().InstanceID, e.Role, e.Permission)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceAdministratorRolePermissionRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*permission.RemovedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-v3ARP02", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AdministratorRoleRepository()
		_, err := repo.RemovePermissions(ctx, v3_sql.SQLTx(tx), e.Aggregate().InstanceID, e.Role, e.Permission)
		return err
	}), nil
}
