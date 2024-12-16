package dialect

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrNegativeRatio       = errors.New("ratio cannot be negative")
	ErrHighSumRatio        = errors.New("sum of pusher and projection ratios must be < 1")
	ErrIllegalMaxOpenConns = errors.New("MaxOpenConns of the database must be higher than 3 or 0 for unlimited")
	ErrIllegalMaxIdleConns = errors.New("MaxIdleConns of the database must be higher than 3 or 0 for unlimited")
	ErrInvalidPurpose      = errors.New("DBPurpose out of range")
)

// ConnectionConfig defines the Max Open and Idle connections for a DB connection pool.
type ConnectionConfig struct {
	MaxOpenConns,
	MaxIdleConns uint32
	AfterConnect []func(ctx context.Context, c *pgx.Conn) error
}

// takeRatio of MaxOpenConns and MaxIdleConns from config and returns
// a new ConnectionConfig with the resulting values.
func (c *ConnectionConfig) takeRatio(ratio float64) (*ConnectionConfig, error) {
	if ratio < 0 {
		return nil, ErrNegativeRatio
	}

	out := &ConnectionConfig{
		MaxOpenConns: uint32(ratio * float64(c.MaxOpenConns)),
		MaxIdleConns: uint32(ratio * float64(c.MaxIdleConns)),
		AfterConnect: c.AfterConnect,
	}
	if c.MaxOpenConns != 0 && out.MaxOpenConns < 1 && ratio > 0 {
		out.MaxOpenConns = 1
	}
	if c.MaxIdleConns != 0 && out.MaxIdleConns < 1 && ratio > 0 {
		out.MaxIdleConns = 1
	}

	return out, nil
}

var afterConnectFuncs []func(ctx context.Context, c *pgx.Conn) error

func RegisterAfterConnect(f func(ctx context.Context, c *pgx.Conn) error) {
	afterConnectFuncs = append(afterConnectFuncs, f)
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
func NewConnectionConfig(openConns, idleConns uint32, pusherRatio, projectionRatio float64, purpose DBPurpose) (*ConnectionConfig, error) {
	if openConns != 0 && openConns < 3 {
		return nil, ErrIllegalMaxOpenConns
	}
	if idleConns != 0 && idleConns < 3 {
		return nil, ErrIllegalMaxIdleConns
	}
	if pusherRatio+projectionRatio >= 1 {
		return nil, ErrHighSumRatio
	}

	queryConfig := &ConnectionConfig{
		MaxOpenConns: openConns,
		MaxIdleConns: idleConns,
		AfterConnect: afterConnectFuncs,
	}
	pusherConfig, err := queryConfig.takeRatio(pusherRatio)
	if err != nil {
		return nil, fmt.Errorf("event pusher: %w", err)
	}

	spoolerConfig, err := queryConfig.takeRatio(projectionRatio)
	if err != nil {
		return nil, fmt.Errorf("projection spooler: %w", err)
	}

	// subtract the claimed amount
	if queryConfig.MaxOpenConns > 0 {
		queryConfig.MaxOpenConns -= pusherConfig.MaxOpenConns + spoolerConfig.MaxOpenConns
	}
	if queryConfig.MaxIdleConns > 0 {
		queryConfig.MaxIdleConns -= pusherConfig.MaxIdleConns + spoolerConfig.MaxIdleConns
	}

	switch purpose {
	case DBPurposeQuery:
		return queryConfig, nil
	case DBPurposeEventPusher:
		return pusherConfig, nil
	case DBPurposeProjectionSpooler:
		return spoolerConfig, nil
	default:
		return nil, fmt.Errorf("%w: %v", ErrInvalidPurpose, purpose)
	}
}
