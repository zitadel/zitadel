package sql

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/repository/orchestrate/handler"
	"github.com/zitadel/zitadel/backend/storage/database"
)

const instanceByDomainQuery = `SELECT i.id, i.name FROM instances i JOIN instance_domains id ON i.id = id.instance_id WHERE id.domain = $1`

func InstanceByDomain(client database.Querier) handler.Handle[string, *repository.Instance] {
	return func(ctx context.Context, domain string) (*repository.Instance, error) {
		row := client.QueryRow(ctx, instanceByDomainQuery, domain)
		var instance repository.Instance
		if err := row.Scan(&instance.ID, &instance.Name); err != nil {
			return nil, err
		}
		return &instance, nil
	}
}

const instanceByIDQuery = `SELECT id, name FROM instances WHERE id = $1`

func InstanceByID(client database.Querier) handler.Handle[string, *repository.Instance] {
	return func(ctx context.Context, id string) (*repository.Instance, error) {
		row := client.QueryRow(ctx, instanceByIDQuery, id)
		var instance repository.Instance
		if err := row.Scan(&instance.ID, &instance.Name); err != nil {
			return nil, err
		}
		return &instance, nil
	}
}

func SetUpInstance(tx database.Transaction) handler.Handle[*repository.Instance, *repository.Instance] {
	return func(ctx context.Context, instance *repository.Instance) (*repository.Instance, error) {
		err := tx.Exec(ctx, "INSERT INTO instances (id, name) VALUES ($1, $2)", instance.ID, instance.Name)
		if err != nil {
			return nil, err
		}
		return instance, nil
	}
}
