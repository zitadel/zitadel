package storage

// import (
// 	"context"
// 	"time"

// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/pgx/v5/pgconn"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type EventStore interface {
// 	Write(ctx context.Context, commands ...Command) error
// }

// type Command interface {
// 	Aggregate() Aggregate

// 	Type() string
// 	Payload() any
// 	Revision() uint8
// 	Creator() string

// 	Operations() []Operation
// }

// type Event interface {
// 	ID() string

// 	Aggregate() Aggregate

// 	CreatedAt() time.Time
// 	Type() string
// 	UnmarshalPayload(ptr any) error

// 	Revision() uint8
// 	Creator() string

// 	Position() float64
// 	InPositionOrder() uint16
// }

// type Aggregate struct {
// 	Type     string
// 	ID       string
// 	Instance string
// }

// type Operation interface {
// 	Exec(ctx context.Context, tx TX) error
// }

// type TX interface {
// 	Query() // Placeholder for future query methods
// }

// type createInstance struct {
// 	id      string
// 	creator string
// 	Name    string
// }

// var (
// 	_ Command = (*createInstance)(nil)
// )

// func (c *createInstance) Aggregate() Aggregate {
// 	return Aggregate{
// 		Type: "instance",
// 		ID:   c.id,
// 	}
// }

// func (c *createInstance) Type() string {
// 	return "instance.created"
// }

// func (c *createInstance) Payload() any {
// 	return c
// }

// func (c *createInstance) Revision() uint8 {
// 	return 1
// }

// func (c *createInstance) Creator() string {
// 	return c.creator
// }

// func (c *createInstance) Operations() []Operation {
// 	return []Operation{}
// }

// type executor[R any] interface {
// 	Exec(ctx context.Context, query string, args ...any) (R, error)
// }

// type querier[R any, Q Query[R]] interface {
// 	Query(ctx context.Context, query Q) (R, error)
// 	// Query(ctx context.Context, query string, args ...any) (R, error)
// }

// type Query[R any] interface {
// 	Query(ctx context.Context) (R, error)
// }

// var _ executor[pgconn.CommandTag] = (*pgxExecutor)(nil)

// type pgxExecutor struct {
// 	conn *pgx.Conn
// }

// func (e *pgxExecutor) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
// 	return e.conn.Exec(ctx, query, args...)
// }

// type pgxQuerier struct {
// 	conn *pgx.Conn
// }

// type pgxCreateInstanceQuery struct {
// 	Name string
// }

// type instance struct{
// 	id string
// 	name string
// }

// type instanceByIDQuery[R Rows] struct{
// 	id string
// }

// type Rows interface {
// 	Next() bool
// 	Scan(dest ...any) error
// }

// func (q *instanceByIDQuery[R]) Query(ctx context.Context) (*instance, error) {
// 	return nil, nil
// }

// type pgxInstanceRepository struct {
// 	pool pgxpool.Pool
// }

// func (r *pgxInstanceRepository) InstanceByID(ctx context.Context, id string) (*instance, error) {
// 	query := &instanceByIDQuery[pgx.Rows]{}
// 	query.
// 	return nil, nil
// }

// func (q *pgxQuerier) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
// 	return q.conn.Query(ctx, query, args...)
// }
