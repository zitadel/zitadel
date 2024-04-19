package database

import (
	"context"
	"database/sql"

	"github.com/zitadel/logging"
)

type Tx interface {
	Commit() error
	Rollback() error
}

func CloseTx(tx Tx, err error) error {
	if err != nil {
		rollbackErr := tx.Rollback()
		logging.OnError(rollbackErr).Debug("unable to rollback")
		return err
	}

	return tx.Commit()
}

type DestMapper[R any] func(index int, scan func(dest ...any) error) (*R, error)

type Rows interface {
	Close() error
	Err() error
	Next() bool
	Scan(dest ...any) error
}

func MapRows[R any](rows Rows, mapper DestMapper[R]) (result []*R, err error) {
	defer func() {
		closeErr := rows.Close()
		logging.OnError(closeErr).Debug("unable to close rows")

		if err == nil && rows.Err() != nil {
			result = nil
			err = rows.Err()
		}
	}()
	for i := 0; rows.Next(); i++ {
		res, err := mapper(i, rows.Scan)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}

	return result, nil
}

func MapRowsToObject(rows Rows, mapper func(scan func(dest ...any) error) error) (err error) {
	defer func() {
		closeErr := rows.Close()
		logging.OnError(closeErr).Debug("unable to close rows")

		if err == nil && rows.Err() != nil {
			err = rows.Err()
		}
	}()
	for rows.Next() {
		err = mapper(rows.Scan)
		if err != nil {
			return err
		}
	}
	return nil
}

type Querier interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}
