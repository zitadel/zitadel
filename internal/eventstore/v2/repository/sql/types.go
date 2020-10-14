package sql

import (
	"database/sql/driver"
)

// Data represents a byte array that may be null.
// Data implements the sql.Scanner interface
type Data []byte

// Scan implements the Scanner interface.
func (data *Data) Scan(value interface{}) error {
	if value == nil {
		*data = nil
		return nil
	}
	*data = Data(value.([]byte))
	return nil
}

// Value implements the driver Valuer interface.
func (data Data) Value() (driver.Value, error) {
	if len(data) == 0 {
		return nil, nil
	}
	return []byte(data), nil
}

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
