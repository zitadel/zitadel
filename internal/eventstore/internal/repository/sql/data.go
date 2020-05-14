package sql

import "database/sql/driver"

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
