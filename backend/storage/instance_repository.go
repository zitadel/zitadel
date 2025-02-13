package storage

// import (
// 	"context"

// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type Instance struct {
// 	ID   string
// 	Name string
// }

// // type InstanceRepository interface {
// // 	InstanceByID(id string) (*Instance, error)
// // }

// type Row interface {
// 	Scan(dest ...interface{}) error
// }

// type RowQuerier[R Row] interface {
// 	QueryRow(ctx context.Context, query string, args ...any) R
// }

// type Rows interface {
// 	Row
// 	Next() bool
// 	Close()
// 	Err() error
// }

// type RowsQuerier[R Rows] interface {
// 	Query(ctx context.Context, query string, args ...any) (R, error)
// }

// type Executor interface {
// 	Exec(ctx context.Context, query string, args ...any) error
// }

// type instanceByIDQuery[R Row, Q RowQuerier[R]] struct {
// 	id      string
// 	querier Q
// }

// func (q instanceByIDQuery[R, Q]) Query(ctx context.Context) (*Instance, error) {
// 	row := q.querier.QueryRow(ctx, "SELECT * FROM instances WHERE id = $1", q.id)
// 	var instance Instance
// 	if err := row.Scan(&instance.ID, &instance.Name); err != nil {
// 		return nil, err
// 	}
// 	return &instance, nil
// }

// type InstanceRepository interface {
// 	ByIDQuery(ctx context.Context, id string) (*Instance, error)
// 	ByDomainQuery(ctx context.Context, domain string) (*Instance, error)
// }

// type InstanceRepositorySQL struct {
// 	pool *pgxpool.Pool
// }

// type InstanceRepositoryMap struct {
// 	instances map[string]*Instance
// 	domains   *InstanceDomainRepositoryMao
// }

// type InstanceDomainRepositoryMao struct {
// 	domains map[string]string
// }

// func GetInstanceByID[R Row, Q RowQuerier[R]](ctx context.Context, querier Q, id string) (*Instance, error) {
// 	row := querier.QueryRow(ctx, "SELECT * FROM instances WHERE id = $1", id)
// 	var instance Instance
// 	if err := row.Scan(&instance.ID, &instance.Name); err != nil {
// 		return nil, err
// 	}
// 	return &instance, nil
// }

// const instanceByDomainQuery = `SELECT
// 	i.*
// FROM
// 	instances i
// JOIN
// 	instance_domains id
//   ON
// 	id.instance_id = i.id
// WHERE
// 	id.domain = $1`

// func GetInstanceByDomain[R Row, Q RowQuerier[R]](ctx context.Context, querier Q, domain string) (*Instance, error) {
// 	row := querier.QueryRow(ctx, instanceByDomainQuery, domain)
// 	var instance Instance
// 	if err := row.Scan(); err != nil {
// 		return nil, err
// 	}
// 	return &instance, nil
// }

// func CreateInstance[E Executor](ctx context.Context, executor E, instance *Instance) error {
// 	return executor.Exec(ctx, "INSERT INTO instances (id, name) VALUES ($1, $2)", instance.ID, instance.Name)
// }

// func bla(ctx context.Context) {
// 	var c *pgxpool.Tx
// 	instance, err := instanceByIDQuery[pgx.Row, *pgxpool.Tx]{querier: c}.Query(ctx)
// 	_, _ = instance, err
// }
