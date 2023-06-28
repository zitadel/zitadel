package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgtype"
	pgtype_v5 "github.com/jackc/pgx/v5/pgtype"
)

type Array[T any] []T

func (a *Array[T]) Scan(src any) error {
	pgTypeMap := pgtype_v5.NewMap()
	var dst []T
	if err := pgTypeMap.SQLScanner(&dst).Scan(src); err != nil {
		return err
	}
	*a = Array[T](dst)
	return nil
}

func (a Array[T]) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	src := pgtype_v5.FlatArray[T](a)

	pgTypeMap := pgtype_v5.NewMap()
	arrayType, ok1 := pgTypeMap.TypeForValue(src)
	elementType, ok2 := pgTypeMap.TypeForValue(src[0])
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("%T not registered as sql driver.Value", src)
	}
	codec := pgtype_v5.ArrayCodec{ElementType: elementType}

	buf, err := codec.PlanEncode(pgTypeMap, arrayType.OID, pgtype_v5.TextFormatCode, src).Encode(src, nil)
	if err != nil {
		return nil, err
	}
	return string(buf), err
}

type EnumArray[F ~int32] []F

// Scan implements the [database/sql.Scanner] interface.
func (s *EnumArray[F]) Scan(src any) error {
	array := new(Array[int32])
	if err := array.Scan(src); err != nil {
		return err
	}
	s.setArray(*array)
	return nil
}

func (s *EnumArray[F]) setArray(array Array[int32]) {
	out := make(EnumArray[F], len(array))
	for k, v := range array {
		out[k] = F(v)
	}
	*s = out
}

// Value implements the [database/sql/driver.Valuer] interface.
func (s EnumArray[F]) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	array := s.toArray()
	return array.Value()
}

func (enums EnumArray[F]) toArray() Array[int32] {
	out := make(Array[int32], len(enums))
	for k, v := range enums {
		out[k] = int32(v)
	}
	return out
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
