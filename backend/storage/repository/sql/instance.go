package sql

import (
	"context"

	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/storage/repository"
)

func NewInstance(client database.QueryExecutor) repository.InstanceRepository {
	return &Instance{client: client}
}

type Instance struct {
	client database.QueryExecutor
}

const instanceByDomainQuery = `SELECT i.id, i.name FROM instances i JOIN instance_domains id ON i.id = id.instance_id WHERE id.domain = $1`

// ByDomain implements [InstanceRepository].
func (r *Instance) ByDomain(ctx context.Context, domain string) (*repository.Instance, error) {
	row := r.client.QueryRow(ctx, instanceByDomainQuery, domain)
	var instance repository.Instance
	if err := row.Scan(&instance.ID, &instance.Name); err != nil {
		return nil, err
	}
	return &instance, nil
}

const instanceByIDQuery = `SELECT id, name FROM instances WHERE id = $1`

// ByID implements [InstanceRepository].
func (r *Instance) ByID(ctx context.Context, id string) (*repository.Instance, error) {
	row := r.client.QueryRow(ctx, instanceByIDQuery, id)
	var instance repository.Instance
	if err := row.Scan(&instance.ID, &instance.Name); err != nil {
		return nil, err
	}
	return &instance, nil
}

// SetUp implements [InstanceRepository].
func (r *Instance) SetUp(ctx context.Context, instance *repository.Instance) error {
	return r.client.Exec(ctx, "INSERT INTO instances (id, name) VALUES ($1, $2)", instance.ID, instance.Name)
}
