package database

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/jackc/pgtype"
)

type TextArray[t ~string] []t

// Scan implements the [database/sql.Scanner] interface.
func (s *TextArray[t]) Scan(src any) error {
	array := new(pgtype.TextArray)
	if err := array.Scan(src); err != nil {
		return err
	}
	return array.AssignTo(s)
}

// Value implements the [database/sql/driver.Valuer] interface.
func (s TextArray[t]) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}

	array := pgtype.TextArray{}
	if err := array.Set(s); err != nil {
		return nil, err
	}

	return array.Value()
}

type arrayField interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32
}

type Array[F arrayField] []F

// Scan implements the [database/sql.Scanner] interface.
func (a *Array[F]) Scan(src any) error {
	array := new(pgtype.Int8Array)
	if err := array.Scan(src); err != nil {
		return err
	}
	elements := make([]int64, len(array.Elements))
	if err := array.AssignTo(&elements); err != nil {
		return err
	}
	*a = make([]F, len(elements))
	for i, element := range elements {
		(*a)[i] = F(element)
	}
	return nil
}

// Value implements the [database/sql/driver.Valuer] interface.
func (a Array[F]) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}

	array := pgtype.Int8Array{}
	if err := array.Set(a); err != nil {
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

type Duration time.Duration

// Scan implements the [database/sql.Scanner] interface.
func (d *Duration) Scan(src any) error {
	interval := new(pgtype.Interval)
	if err := interval.Scan(src); err != nil {
		return err
	}
	*d = Duration(time.Duration(interval.Microseconds*1000) + time.Duration(interval.Days)*24*time.Hour + time.Duration(interval.Months)*30*24*time.Hour)
	return nil
}

// NullDuration can be used for NULL intervals.
// If Valid is false, the scanned value was NULL
// This behavior is similar to [database/sql.NullString]
type NullDuration struct {
	Valid    bool
	Duration time.Duration
}

// Scan implements the [database/sql.Scanner] interface.
func (d *NullDuration) Scan(src any) error {
	if src == nil {
		d.Duration, d.Valid = 0, false
		return nil
	}
	duration := new(Duration)
	if err := duration.Scan(src); err != nil {
		return err
	}
	d.Duration, d.Valid = time.Duration(*duration), true
	return nil
}
