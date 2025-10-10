package database

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	for _, test := range []struct {
		name string
		err  error
		want string
	}{
		{
			name: "missing condition without column",
			err:  NewMissingConditionError(nil),
			want: "missing condition for column",
		},
		{
			name: "missing condition with column",
			err:  NewMissingConditionError(NewColumn("table", "column")),
			want: "missing condition for column on table.column",
		},
		{
			name: "no row found without original error",
			err:  NewNoRowFoundError(nil),
			want: "no row found",
		},
		{
			name: "no row found with original error",
			err:  NewNoRowFoundError(errors.New("original error")),
			want: "no row found: original error",
		},
		{
			name: "multiple rows found without original error",
			err:  NewMultipleRowsFoundError(nil),
			want: "multiple rows found",
		},
		{
			name: "multiple rows found with original error",
			err:  NewMultipleRowsFoundError(errors.New("original error")),
			want: "multiple rows found: original error",
		},
		{
			name: "check violation without original error",
			err:  NewCheckError("table", "constraint", nil),
			want: `integrity violation of type "check" on "table" (constraint: "constraint")`,
		},
		{
			name: "check violation with original error",
			err:  NewCheckError("table", "constraint", errors.New("original error")),
			want: `integrity violation of type "check" on "table" (constraint: "constraint"): original error`,
		},
		{
			name: "unique violation without original error",
			err:  NewUniqueError("table", "constraint", nil),
			want: `integrity violation of type "unique" on "table" (constraint: "constraint")`,
		},
		{
			name: "unique violation with original error",
			err:  NewUniqueError("table", "constraint", errors.New("original error")),
			want: `integrity violation of type "unique" on "table" (constraint: "constraint"): original error`,
		},
		{
			name: "foreign key violation without original error",
			err:  NewForeignKeyError("table", "constraint", nil),
			want: `integrity violation of type "foreign" on "table" (constraint: "constraint")`,
		},
		{
			name: "foreign key violation with original error",
			err:  NewForeignKeyError("table", "constraint", errors.New("original error")),
			want: `integrity violation of type "foreign" on "table" (constraint: "constraint"): original error`,
		},
		{
			name: "not null violation without original error",
			err:  NewNotNullError("table", "constraint", nil),
			want: `integrity violation of type "not null" on "table" (constraint: "constraint")`,
		},
		{
			name: "not null violation with original error",
			err:  NewNotNullError("table", "constraint", errors.New("original error")),
			want: `integrity violation of type "not null" on "table" (constraint: "constraint"): original error`,
		},
		{
			name: "unknown error",
			err:  NewUnknownError(errors.New("original error")),
			want: `unknown database error: original error`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.err.Error())
		})
	}
}

func TestUnwrap(t *testing.T) {
	originalErr := errors.New("original error")
	for _, test := range []struct {
		name string
		err  error
		want error
	}{
		{
			name: "missing condition without column",
			err:  NewMissingConditionError(nil),
			want: nil,
		},
		{
			name: "missing condition with column",
			err:  NewMissingConditionError(NewColumn("table", "column")),
			want: nil,
		},
		{
			name: "no row found without original error",
			err:  NewNoRowFoundError(nil),
			want: nil,
		},
		{
			name: "no row found with original error",
			err:  NewNoRowFoundError(errors.New("original error")),
			want: originalErr,
		},
		{
			name: "multiple rows found without original error",
			err:  NewMultipleRowsFoundError(nil),
			want: nil,
		},
		{
			name: "multiple rows found with original error",
			err:  NewMultipleRowsFoundError(originalErr),
			want: originalErr,
		},
		{
			name: "check violation without original error",
			err:  NewCheckError("table", "constraint", nil),
			want: &IntegrityViolationError{
				integrityType: IntegrityTypeCheck,
				table:         "table",
				constraint:    "constraint",
				original:      nil,
			},
		},
		{
			name: "check violation with original error",
			err:  NewCheckError("table", "constraint", originalErr),
			want: &IntegrityViolationError{
				integrityType: IntegrityTypeCheck,
				table:         "table",
				constraint:    "constraint",
				original:      originalErr,
			},
		},
		{
			name: "unique violation without original error",
			err:  NewUniqueError("table", "constraint", nil),
			want: &IntegrityViolationError{
				integrityType: IntegrityTypeUnique,
				table:         "table",
				constraint:    "constraint",
				original:      nil,
			},
		},
		{
			name: "unique violation with original error",
			err:  NewUniqueError("table", "constraint", originalErr),
			want: &IntegrityViolationError{
				integrityType: IntegrityTypeUnique,
				table:         "table",
				constraint:    "constraint",
				original:      originalErr,
			},
		},
		{
			name: "foreign key violation without original error",
			err:  NewForeignKeyError("table", "constraint", nil),
			want: &IntegrityViolationError{
				integrityType: IntegrityTypeForeign,
				table:         "table",
				constraint:    "constraint",
				original:      nil,
			},
		},
		{
			name: "foreign key violation with original error",
			err:  NewForeignKeyError("table", "constraint", originalErr),
			want: &IntegrityViolationError{
				integrityType: IntegrityTypeForeign,
				table:         "table",
				constraint:    "constraint",
				original:      originalErr,
			},
		},
		{
			name: "not null violation without original error",
			err:  NewNotNullError("table", "constraint", nil),
			want: &IntegrityViolationError{
				integrityType: IntegrityTypeNotNull,
				table:         "table",
				constraint:    "constraint",
				original:      nil,
			},
		},
		{
			name: "not null violation with original error",
			err:  NewNotNullError("table", "constraint", originalErr),
			want: &IntegrityViolationError{
				integrityType: IntegrityTypeNotNull,
				table:         "table",
				constraint:    "constraint",
				original:      originalErr,
			},
		},
		{
			name: "unwrap integrity violation error",
			err:  errors.Unwrap(NewNotNullError("table", "constraint", originalErr)),
			want: originalErr,
		},
		{
			name: "unknown error",
			err:  NewUnknownError(originalErr),
			want: originalErr,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, errors.Unwrap(test.err))
		})
	}
}

func TestIs(t *testing.T) {
	originalErr := errors.New("original error")
	for _, test := range []struct {
		name string
		err  error
		want error
	}{
		{
			name: "missing condition",
			err:  NewMissingConditionError(NewColumn("table", "column")),
			want: new(MissingConditionError),
		},
		{
			name: "no row found",
			err:  NewNoRowFoundError(errors.New("original error")),
			want: new(NoRowFoundError),
		},
		{
			name: "multiple rows found",
			err:  NewMultipleRowsFoundError(originalErr),
			want: new(MultipleRowsFoundError),
		},
		{
			name: "check violation is for integrity",
			err:  NewCheckError("table", "constraint", nil),
			want: new(IntegrityViolationError),
		},
		{
			name: "check violation is check violation",
			err:  NewCheckError("table", "constraint", nil),
			want: new(CheckError),
		},
		{
			name: "unique violation is for integrity",
			err:  NewUniqueError("table", "constraint", nil),
			want: new(IntegrityViolationError),
		},
		{
			name: "unique violation is unique violation",
			err:  NewUniqueError("table", "constraint", nil),
			want: new(UniqueError),
		},
		{
			name: "foreign key violation is for integrity",
			err:  NewForeignKeyError("table", "constraint", nil),
			want: new(IntegrityViolationError),
		},
		{
			name: "foreign key violation is foreign key violation",
			err:  NewForeignKeyError("table", "constraint", nil),
			want: new(ForeignKeyError),
		},
		{
			name: "not null violation is for integrity",
			err:  NewNotNullError("table", "constraint", nil),
			want: new(IntegrityViolationError),
		},
		{
			name: "not null violation is not null violation",
			err:  NewNotNullError("table", "constraint", nil),
			want: new(NotNullError),
		},
		{
			name: "unknown error",
			err:  NewUnknownError(originalErr),
			want: new(UnknownError),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.ErrorIs(t, test.err, test.want)
		})
	}
}
