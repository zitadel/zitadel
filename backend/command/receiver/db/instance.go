package db

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/receiver"
	"github.com/zitadel/zitadel/backend/storage/database"
)

// NewInstance returns a new instance receiver.
func NewInstance(client database.QueryExecutor) receiver.InstanceManipulator {
	return &instance{client: client}
}

// instance is the sql interface for instances.
type instance struct {
	client database.QueryExecutor
}

// ByID implements receiver.InstanceReader.
func (i *instance) ByID(ctx context.Context, id string) (*receiver.Instance, error) {
	var instance receiver.Instance
	err := i.client.QueryRow(ctx, "SELECT id, name, state FROM instances WHERE id = $1", id).
		Scan(
			&instance.ID,
			&instance.Name,
			&instance.State,
		)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

// AddDomain implements [receiver.InstanceManipulator].
func (i *instance) AddDomain(ctx context.Context, instance *receiver.Instance, domain *receiver.Domain) error {
	return i.client.Exec(ctx, "INSERT INTO instance_domains (instance_id, domain, is_primary) VALUES ($1, $2, $3)", instance.ID, domain.Name, domain.IsPrimary)
}

// Create implements [receiver.InstanceManipulator].
func (i *instance) Create(ctx context.Context, instance *receiver.Instance) error {
	return i.client.Exec(ctx, "INSERT INTO instances (id, name, state) VALUES ($1, $2, $3)", instance.ID, instance.Name, instance.State)
}

// Delete implements [receiver.InstanceManipulator].
func (i *instance) Delete(ctx context.Context, instance *receiver.Instance) error {
	return i.client.Exec(ctx, "DELETE FROM instances WHERE id = $1", instance.ID)
}

// SetPrimaryDomain implements [receiver.InstanceManipulator].
func (i *instance) SetPrimaryDomain(ctx context.Context, instance *receiver.Instance, domain *receiver.Domain) error {
	return i.client.Exec(ctx, "UPDATE instance_domains SET is_primary = domain = $1 WHERE instance_id = $2", domain.Name, instance.ID)
}

var (
	_ receiver.InstanceManipulator = (*instance)(nil)
	_ receiver.InstanceReader      = (*instance)(nil)
)
