package database

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgtype"
)

type StringArray []string

// Scan implements the [database/sql.Scanner] interface.
func (s *StringArray) Scan(src any) error {
	array := new(pgtype.TextArray)
	if err := array.Scan(src); err != nil {
		return err
	}
	if err := array.AssignTo(s); err != nil {
		return err
	}
	return nil
}

// Value implements the [database/sql/driver.Valuer] interface.
func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}

	array := pgtype.TextArray{}
	if err := array.Set(s); err != nil {
		return nil, err
	}

	return array.Value()
}

type enumField interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32
}

type EnumArray[F enumField] []F

// Scan implements the [database/sql.Scanner] interface.
func (s *EnumArray[F]) Scan(src any) error {
	array := new(pgtype.Int2Array)
	if err := array.Scan(src); err != nil {
		return err
	}
	ints := make([]int32, 0, len(array.Elements))
	if err := array.AssignTo(&ints); err != nil {
		return err
	}
	*s = make([]F, len(ints))
	for i, a := range ints {
		(*s)[i] = F(a)
	}
	return nil
}

// Value implements the [database/sql/driver.Valuer] interface.
func (s EnumArray[F]) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}

	array := pgtype.Int2Array{}
	if err := array.Set(s); err != nil {
		return nil, err
	}

	return array.Value()
}

type Map[V any] map[string]V

// Scan implements the [database/sql.Scanner] interface.
func (m *Map[V]) Scan(src any) error {
	bytea := new(pgtype.Bytea)
	if err := bytea.Scan(src); err != nil {
		return err
	}
	if len(bytea.Bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytea.Bytes, &m)
}

// Value implements the [database/sql/driver.Valuer] interface.
func (m Map[V]) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return json.Marshal(m)
}
