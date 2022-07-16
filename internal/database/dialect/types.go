package dialect

import (
	"database/sql"
	"database/sql/driver"
)

type User struct {
	Username string
	Password string
	SSL      SSL
}

type SSL struct {
	// type of connection security
	Mode string
	// RootCert Path to the CA certificate
	RootCert string
	// Cert Path to the client certificate
	Cert string
	// Key Path to the client private key
	Key string
}

type Text interface {
	~string | ~[]byte
}

type Field[T any] interface {
	// Set([]T)
	// Get() []T
	// T

	sql.Scanner
	driver.Valuer
}

type Array[T any] interface {
	// Set([]T)
	// Get() []T
	[]T

	sql.Scanner
	driver.Valuer
}

type TextArray[T Text] interface {
	Array[T]
}

// import (
// 	"database/sql"
// 	"database/sql/driver"

// 	"github.com/jackc/pgtype"
// )

// type text interface {
// 	~string | ~[]byte
// }

// type TextArray[T any] struct {
// 	asdf []T

// 	sql.Scanner
// 	driver.Valuer
// }

// func RegisterBlablaScanner()

// type PostgresTextArray struct{
// 	pgtype.TextArray
// }

// var _ TextArray[string] = (*PostgresTextArray[string])(nil)

/*




type PostgresTextArray[T text] pgtype.TextArray

func (a PostgresTextArray[T]) Scan(src any) error {

	return nil
}

func (a PostgresTextArray[T]) Value() (driver.Value, error) {
	return nil, nil
}

func (a PostgresTextArray[T]) Get() []T {
	return nil
}

func (a PostgresTextArray[T]) Set([]T) {
	// return nil
}

*/

/*
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

	switch s := src.(type) {
	case string:
		return (*pgtype.NumericArray)(dst).DecodeText(nil, []byte(s))
	case []byte:
		return (*pgtype.NumericArray)(dst).DecodeText(nil, s)
	}

	return fmt.Errorf("unable to parse int array")
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


*/
