package sql

import (
	"database/sql/driver"
)

// Sequence represents a number that may be null.
// Sequence implements the sql.Scanner interface
type Sequence uint64

// Scan implements the Scanner interface.
func (seq *Sequence) Scan(value interface{}) error {
	if value == nil {
		*seq = 0
		return nil
	}
	*seq = Sequence(value.(int64))
	return nil
}

// Value implements the driver Valuer interface.
func (seq Sequence) Value() (driver.Value, error) {
	if seq == 0 {
		return nil, nil
	}
	return int64(seq), nil
}
