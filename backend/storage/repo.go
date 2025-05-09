package storage

// import (
// 	"context"

// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type row interface {
// 	Scan(dest ...any) error
// }

// type rows interface {
// 	row
// 	Next() bool
// 	Close()
// 	Err() error
// }

// type querier interface {
// 	NewRowQuery()
// 	NewRowsQuery()
// }

// type Client interface {
// 	// querier
// 	InstanceRepository() instanceRepository
// }

// type instanceRepository interface {
// 	ByIDQuery(ctx context.Context, id string) (*Instance, error)
// 	ByDomainQuery(ctx context.Context, domain string) (*Instance, error)
// }

// type pgxClient pgxpool.Pool

// func (c *pgxClient) Begin(ctx context.Context) (*pgxTx, error) {
// 	tx, err := (*pgxpool.Pool)(c).Begin(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return (*pgxTx)(tx.(*pgxpool.Tx)), nil
// }

// func (c *pgxClient) InstanceRepository() instanceRepository {
// 	return &pgxInstanceRepository[pgxpool.Pool]{client: (*pgxpool.Pool)(c)}
// }

// type pgxTx pgxpool.Tx

// func (c *pgxTx) InstanceRepository() instanceRepository {
// 	return &pgxInstanceRepository[pgxpool.Tx]{client: (*pgxpool.Tx)(c)}
// }

// type pgxInstanceRepository[C pgxpool.Pool | pgxpool.Tx] struct {
// 	client *C
// }

// func (r *pgxInstanceRepository[C]) ByIDQuery(ctx context.Context, id string) (*Instance, error) {
// 	// return r.client
// 	pgx.Tx
// 	return nil, nil
// }

// func (r *pgxInstanceRepository[C]) ByDomainQuery(ctx context.Context, domain string) (*Instance, error) {

// 	return nil, nil
// }

// func blabla() {
// 	var client Client = &pgxClient{&pgxpool.Pool{}}

// 	client.be
// }
