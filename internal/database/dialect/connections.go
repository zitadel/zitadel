package dialect

import (
	"context"
	"errors"
	"reflect"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrIllegalMaxOpenConns = errors.New("MaxOpenConns of the database must be higher than 3 or 0 for unlimited")
	ErrIllegalMaxIdleConns = errors.New("MaxIdleConns of the database must be higher than 3 or 0 for unlimited")
)

// ConnectionConfig defines the Max Open and Idle connections for a DB connection pool.
type ConnectionConfig struct {
	MaxOpenConns,
	MaxIdleConns uint32
	AfterConnect  []func(ctx context.Context, c *pgx.Conn) error
	BeforeAcquire []func(ctx context.Context, c *pgx.Conn) error
	AfterRelease  []func(c *pgx.Conn) error
}

var afterConnectFuncs = []func(ctx context.Context, c *pgx.Conn) error{
	func(ctx context.Context, c *pgx.Conn) error {
		pgxdecimal.Register(c.TypeMap())
		return nil
	},
}

func RegisterAfterConnect(f func(ctx context.Context, c *pgx.Conn) error) {
	afterConnectFuncs = append(afterConnectFuncs, f)
}

var beforeAcquireFuncs []func(ctx context.Context, c *pgx.Conn) error

func RegisterBeforeAcquire(f func(ctx context.Context, c *pgx.Conn) error) {
	beforeAcquireFuncs = append(beforeAcquireFuncs, f)
}

var afterReleaseFuncs []func(c *pgx.Conn) error

func RegisterAfterRelease(f func(c *pgx.Conn) error) {
	afterReleaseFuncs = append(afterReleaseFuncs, f)
}

func RegisterDefaultPgTypeVariants[T any](m *pgtype.Map, name, arrayName string) {
	// T
	var value T
	m.RegisterDefaultPgType(value, name)

	// *T
	valueType := reflect.TypeOf(value)
	m.RegisterDefaultPgType(reflect.New(valueType).Interface(), name)

	// []T
	sliceType := reflect.SliceOf(valueType)
	m.RegisterDefaultPgType(reflect.MakeSlice(sliceType, 0, 0).Interface(), arrayName)

	// *[]T
	m.RegisterDefaultPgType(reflect.New(sliceType).Interface(), arrayName)

	// []*T
	sliceOfPointerType := reflect.SliceOf(reflect.TypeOf(reflect.New(valueType).Interface()))
	m.RegisterDefaultPgType(reflect.MakeSlice(sliceOfPointerType, 0, 0).Interface(), arrayName)

	// *[]*T
	m.RegisterDefaultPgType(reflect.New(sliceOfPointerType).Interface(), arrayName)
}

// NewConnectionConfig calculates [ConnectionConfig] values from the passed ratios
// and returns the config applicable for the requested purpose.
//
// openConns and idleConns must be at least 3 or 0, which means no limit.
// The pusherRatio and spoolerRatio must be between 0 and 1.
func NewConnectionConfig(openConns, idleConns uint32) *ConnectionConfig {
	return &ConnectionConfig{
		MaxOpenConns:  openConns,
		MaxIdleConns:  idleConns,
		AfterConnect:  afterConnectFuncs,
		BeforeAcquire: beforeAcquireFuncs,
		AfterRelease:  afterReleaseFuncs,
	}
}
