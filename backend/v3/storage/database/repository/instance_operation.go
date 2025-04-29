package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type instanceOperation struct {
	database.QueryExecutor
	id string
}

const addInstanceAdminStmt = `INSERT INTO instance_admins (instance_id, user_id, roles) VALUES ($1, $2, $3)`

// AddAdmin implements [domain.InstanceOperation].
func (i *instanceOperation) AddAdmin(ctx context.Context, userID string, roles []string) error {
	return i.QueryExecutor.Exec(ctx, addInstanceAdminStmt, i.id, userID, roles)
}

// Delete implements [domain.InstanceOperation].
func (i *instanceOperation) Delete(ctx context.Context) error {
	return i.QueryExecutor.Exec(ctx, `DELETE FROM instances WHERE id = $1`, i.id)
}

const removeInstanceAdminStmt = `DELETE FROM instance_admins WHERE instance_id = $1 AND user_id = $2`

// RemoveAdmin implements [domain.InstanceOperation].
func (i *instanceOperation) RemoveAdmin(ctx context.Context, userID string) error {
	return i.QueryExecutor.Exec(ctx, removeInstanceAdminStmt, i.id, userID)
}

const setInstanceAdminRolesStmt = `UPDATE instance_admins SET roles = $1 WHERE instance_id = $2 AND user_id = $3`

// SetAdminRoles implements [domain.InstanceOperation].
func (i *instanceOperation) SetAdminRoles(ctx context.Context, userID string, roles []string) error {
	return i.QueryExecutor.Exec(ctx, setInstanceAdminRolesStmt, roles, i.id, userID)
}

const updateInstanceStmt = `UPDATE instances SET name = $1, updated_at = $2 WHERE id = $3 RETURNING updated_at`

// Update implements [domain.InstanceOperation].
func (i *instanceOperation) Update(ctx context.Context, instance *domain.Instance) error {
	return i.QueryExecutor.QueryRow(ctx, updateInstanceStmt,
		instance.Name,
		instance.UpdatedAt,
		i.id,
	).Scan(&instance.UpdatedAt)
}

var _ domain.InstanceOperation = (*instanceOperation)(nil)
