package database

import (
	"database/sql"
	"database/sql/driver"

	"github.com/jackc/pgtype"
	"golang.org/x/exp/constraints"

	"github.com/zitadel/zitadel/internal/errors"
)

type text interface {
	~string | ~[]byte
}

type TextArray[T text] pgtype.TextArray

type IntegerArray[T constraints.Integer] pgtype.NumericArray

var _ sql.Scanner = (*TextArray[string])(nil)

// Scan implements the database/sql Scanner interface.
func (dst *TextArray[T]) Scan(src any) error {
	if src == nil {
		return (*pgtype.TextArray)(dst).DecodeText(nil, nil)
	}

	t, ok := src.(T)
	if !ok {
		return errors.ThrowInvalidArgumentf(nil, "DATAB-TODaM", "unexpected type %T | expected %T", t, *(new(T)))
	}

	return (*pgtype.TextArray)(dst).DecodeText(nil, []byte(t))
}

// Value implements the database/sql/driver Valuer interface.
func (src TextArray[T]) Value() (driver.Value, error) {
	buf, err := (pgtype.TextArray)(src).EncodeText(nil, nil)
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}

	return string(buf), nil
}

func (data *TextArray[T]) Data() []T {
	res := make([]T, len(data.Elements))
	for i, e := range data.Elements {
		res[i] = T(e.String)
	}
	return res
}

func (dst *IntegerArray[T]) Scan(src any) error {
	if src == nil {
		return (*pgtype.NumericArray)(dst).DecodeText(nil, nil)
	}

	return (*pgtype.NumericArray)(dst).DecodeText(nil, src.([]byte))
}

// Value implements the database/sql/driver Valuer interface.
func (src IntegerArray[T]) Value() (driver.Value, error) {
	buf, err := (pgtype.NumericArray)(src).EncodeText(nil, nil)
	if err != nil || buf == nil {
		return nil, err
	}

	return string(buf), nil
}

func (data *IntegerArray[T]) Data() []T {
	res := make([]T, len(data.Elements))
	for i, e := range data.Elements {
		if e.Int.IsInt64() {
			res[i] = T(e.Int.Int64())
		} else {
			res[i] = T(e.Int.Uint64())
		}
	}
	return res
}
