package repository

import (
	"context"
	"log"
)

const InstanceByIDStmt = `SELECT id, name FROM instances WHERE id = $1`

func (q *querier) InstanceByID(ctx context.Context, id string) (*Instance, error) {
	log.Println("sql.instance.byID")
	row := q.client.QueryRow(ctx, InstanceByIDStmt, id)
	var instance Instance
	if err := row.Scan(&instance.ID, &instance.Name); err != nil {
		return nil, err
	}
	return &instance, nil
}

const instanceByDomainQuery = `SELECT i.id, i.name FROM instances i JOIN instance_domains id ON i.id = id.instance_id WHERE id.domain = $1`

func (q *querier) InstanceByDomain(ctx context.Context, domain string) (*Instance, error) {
	log.Println("sql.instance.byDomain")
	row := q.client.QueryRow(ctx, instanceByDomainQuery, domain)
	var instance Instance
	if err := row.Scan(&instance.ID, &instance.Name); err != nil {
		return nil, err
	}
	return &instance, nil
}

func (q *querier) ListInstances(ctx context.Context, request *ListRequest) (res []*Instance, err error) {
	log.Println("sql.instance.list")
	rows, err := q.client.Query(ctx, "SELECT id, name FROM instances")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var instance Instance
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

const InstanceCreateStmt = `INSERT INTO instances (id, name) VALUES ($1, $2)`

func (e *executor) CreateInstance(ctx context.Context, instance *Instance) (*Instance, error) {
	log.Println("sql.instance.create")
	err := e.client.Exec(ctx, InstanceCreateStmt, instance.ID, instance.Name)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
