package sql

import (
	"context"
	"log"

	"github.com/zitadel/zitadel/backend/repository"
)

const instanceByIDQuery = `SELECT id, name FROM instances WHERE id = $1`

func (q *querier[C]) InstanceByID(ctx context.Context, id string) (*repository.Instance, error) {
	log.Println("sql.instance.byID")
	row := q.client.QueryRow(ctx, instanceByIDQuery, id)
	var instance repository.Instance
	if err := row.Scan(&instance.ID, &instance.Name); err != nil {
		return nil, err
	}
	return &instance, nil
}

const instanceByDomainQuery = `SELECT i.id, i.name FROM instances i JOIN instance_domains id ON i.id = id.instance_id WHERE id.domain = $1`

func (q *querier[C]) InstanceByDomain(ctx context.Context, domain string) (*repository.Instance, error) {
	log.Println("sql.instance.byDomain")
	row := q.client.QueryRow(ctx, instanceByDomainQuery, domain)
	var instance repository.Instance
	if err := row.Scan(&instance.ID, &instance.Name); err != nil {
		return nil, err
	}
	return &instance, nil
}

func (q *querier[C]) ListInstances(ctx context.Context, request *repository.ListRequest) (res []*repository.Instance, err error) {
	log.Println("sql.instance.list")
	rows, err := q.client.Query(ctx, "SELECT id, name FROM instances")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var instance repository.Instance
		if err = rows.Scan(&instance.ID, &instance.Name); err != nil {
			return nil, err
		}
		res = append(res, &instance)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (e *executor[C]) CreateInstance(ctx context.Context, instance *repository.Instance) (*repository.Instance, error) {
	log.Println("sql.instance.create")
	err := e.client.Exec(ctx, "INSERT INTO instances (id, name) VALUES ($1, $2)", instance.ID, instance.Name)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
