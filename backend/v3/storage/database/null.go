package database

import "database/sql"

// Null is a type alias for sql.Null[T] to avoid importing the sql package into repositories.
type Null[T any] = sql.Null[T]
