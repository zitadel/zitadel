package database

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type TextArray[T ~string] pgtype.FlatArray[T]

// Scan implements the [database/sql.Scanner] interface.
func (s *TextArray[T]) Scan(src any) error {
	var typedArray []string
	err := pgtype.NewMap().SQLScanner(&typedArray).Scan(src)
	if err != nil {
		return err
	}

	(*s) = make(TextArray[T], len(typedArray))
	for i, value := range typedArray {
		(*s)[i] = T(value)
	}

	return nil
}

// Value implements the [database/sql/driver.Valuer] interface.
func (s TextArray[T]) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}

	typed := make([]string, len(s))

	for i, value := range s {
		typed[i] = string(value)
	}

	return []byte("{" + strings.Join(typed, ",") + "}"), nil
}

type ByteArray[T ~byte] pgtype.FlatArray[T]

// Scan implements the [database/sql.Scanner] interface.
func (s *ByteArray[T]) Scan(src any) error {
	var typedArray []byte
	typedArray, ok := src.([]byte)
	if !ok {
		// tests use a different src type
		err := pgtype.NewMap().SQLScanner(&typedArray).Scan(src)
		if err != nil {
			return err
		}
	}

	(*s) = make(ByteArray[T], len(typedArray))
	for i, value := range typedArray {
		(*s)[i] = T(value)
	}

	return nil
}

// Value implements the [database/sql/driver.Valuer] interface.
func (s ByteArray[T]) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	typed := make([]byte, len(s))

	for i, value := range s {
		typed[i] = byte(value)
	}

	return typed, nil
}

type numberField interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64 | ~int | ~uint
}

type numberTypeField interface {
	int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | int | uint
}

var _ sql.Scanner = (*NumberArray[int8])(nil)

type NumberArray[F numberField] pgtype.FlatArray[F]

// Scan implements the [database/sql.Scanner] interface.
func (a *NumberArray[F]) Scan(src any) (err error) {
	var (
		mapper  func()
		scanner sql.Scanner
	)

	//nolint: exhaustive
	// only defined types
	switch reflect.TypeOf(*a).Elem().Kind() {
	case reflect.Int8:
		mapper, scanner = castedScan[int8](a)
	case reflect.Uint8:
		// we provide int16 is a workaround because pgx thinks we want to scan a byte array if we provide uint8
		mapper, scanner = castedScan[int16](a)
	case reflect.Int16:
		mapper, scanner = castedScan[int16](a)
	case reflect.Uint16:
		mapper, scanner = castedScan[uint16](a)
	case reflect.Int32:
		mapper, scanner = castedScan[int32](a)
	case reflect.Uint32:
		mapper, scanner = castedScan[uint32](a)
	case reflect.Int64:
		mapper, scanner = castedScan[int64](a)
	case reflect.Uint64:
		mapper, scanner = castedScan[uint64](a)
	case reflect.Int:
		mapper, scanner = castedScan[int](a)
	case reflect.Uint:
		mapper, scanner = castedScan[uint](a)
	}

	if err = scanner.Scan(src); err != nil {
		return err
	}
	mapper()

	return nil
}

func castedScan[T numberTypeField, F numberField](a *NumberArray[F]) (mapper func(), scanner sql.Scanner) {
	var typedArray []T

	mapper = func() {
		(*a) = make(NumberArray[F], len(typedArray))
		for i, value := range typedArray {
			(*a)[i] = F(value)
		}
	}
	scanner = pgtype.NewMap().SQLScanner(&typedArray)

	return mapper, scanner
}

type Map[V any] map[string]V

// Scan implements the [database/sql.Scanner] interface.
func (m *Map[V]) Scan(src any) error {
	if src == nil {
		return nil
	}

	bytes := src.([]byte)
	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, &m)
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
	switch duration := src.(type) {
	case *time.Duration:
		*d = Duration(*duration)
		return nil
	case time.Duration:
		*d = Duration(duration)
		return nil
	case *pgtype.Interval:
		*d = intervalToDuration(duration)
		return nil
	case pgtype.Interval:
		*d = intervalToDuration(&duration)
		return nil
	case int64:
		*d = Duration(duration)
		return nil
	}
	interval := new(pgtype.Interval)
	if err := interval.Scan(src); err != nil {
		return err
	}
	*d = intervalToDuration(interval)
	return nil
}

func intervalToDuration(interval *pgtype.Interval) Duration {
	return Duration(time.Duration(interval.Microseconds*1000) + time.Duration(interval.Days)*24*time.Hour + time.Duration(interval.Months)*30*24*time.Hour)
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
