package cockroach

import (
	"database/sql/driver"

	"github.com/jackc/pgtype"

	"github.com/zitadel/zitadel/internal/database/dialect"
)

// var _ database.Array[string] = (*TextArray[string])(nil)

type TextArray[T dialect.Text] []T

// Scan implements the database/sql Scanner interface.
func (dst *TextArray[T]) Scan(src any) error {
	d := new(pgtype.TextArray)

	if err := d.Scan(src); err != nil {
		return err
	}

	if d.Status == pgtype.Null {
		return nil
	}

	*dst = make([]T, len(d.Elements))
	for i, element := range d.Elements {
		(*dst)[i] = T(element.String)
	}

	return nil
}

// Value implements the database/sql/driver Valuer interface.
func (src TextArray[T]) Value() (driver.Value, error) {
	s := pgtype.TextArray{Elements: make([]pgtype.Text, len(src))}
	for i, elem := range src {
		s.Elements[i] = pgtype.Text{String: string(elem)}
	}
	buf, err := s.EncodeText(nil, nil)
	if err != nil || buf == nil {
		return nil, err
	}

	return string(buf), nil
}
