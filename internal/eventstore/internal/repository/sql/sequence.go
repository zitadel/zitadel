package sql

import (
	"database/sql/driver"
)

// Sequence represents a number that may be null.
// Sequence implements the sql.Scanner interface so
type Sequence uint64

// Scan implements the Scanner interface.
func (n *Sequence) Scan(value interface{}) error {
	if value == nil {
		*n = 0
		return nil
	}
	*n = Sequence(value.(int64))
	return nil
}

// Value implements the driver Valuer interface.
func (seq Sequence) Value() (driver.Value, error) {
	if seq == 0 {
		return nil, nil
	}
	return int64(seq), nil
}
