package eventstore

import "database/sql/driver"

// Payload represents a byte array that may be null.
// Payload implements the sql.Scanner interface
type Payload []byte

// Scan implements the Scanner interface.
func (data *Payload) Scan(value interface{}) error {
	if value == nil {
		*data = nil
		return nil
	}
	*data = Payload(value.([]byte))
	return nil
}

// Value implements the driver Valuer interface.
func (data Payload) Value() (driver.Value, error) {
	if len(data) == 0 {
		return nil, nil
	}
	return []byte(data), nil
}
