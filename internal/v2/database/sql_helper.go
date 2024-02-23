package database

import "database/sql"

func CloseTx(tx *sql.Tx, err error) error {
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

type DestMapper[R any] func(index int, scan func(dest ...any) error) (*R, error)

func MapRows[R any](rows *sql.Rows, mapper DestMapper[R]) (result []*R, err error) {
	defer func() {
		rows.Close()
		if rows.Err() != nil {
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

func MapRowsToObject(rows *sql.Rows, mapper func(scan func(dest ...any) error) error) (err error) {
	defer func() {
		rows.Close()
		if rows.Err() != nil {
			err = rows.Err()
		}
	}()
	for i := 0; rows.Next(); i++ {
		err = mapper(rows.Scan)
		if err != nil {
			return err
		}
	}
	return nil
}
