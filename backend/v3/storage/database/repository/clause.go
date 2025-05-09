package repository

import (
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/domain"
)

type field interface {
	fmt.Stringer
}

type fieldDescriptor struct {
	schema string
	table  string
	name   string
}

func (f fieldDescriptor) String() string {
	return f.schema + "." + f.table + "." + f.name
}

type ignoreCaseFieldDescriptor struct {
	fieldDescriptor
	fieldNameSuffix string
}

func (f ignoreCaseFieldDescriptor) String() string {
	return f.fieldDescriptor.String() + f.fieldNameSuffix
}

type textFieldDescriptor struct {
	field
	isIgnoreCase bool
}

type clause[Op domain.Operation] struct {
	field field
	op    Op
}

const (
	schema    = "zitadel"
	userTable = "users"
)

var userFields = map[domain.UserField]field{
	domain.UserFieldInstanceID: fieldDescriptor{
		schema: schema,
		table:  userTable,
		name:   "instance_id",
	},
	domain.UserFieldOrgID: fieldDescriptor{
		schema: schema,
		table:  userTable,
		name:   "org_id",
	},
	domain.UserFieldID: fieldDescriptor{
		schema: schema,
		table:  userTable,
		name:   "id",
	},
	domain.UserFieldUsername: textFieldDescriptor{
		field: ignoreCaseFieldDescriptor{
			fieldDescriptor: fieldDescriptor{
				schema: schema,
				table:  userTable,
				name:   "username",
			},
			fieldNameSuffix: "_lower",
		},
	},
	domain.UserHumanFieldEmail: textFieldDescriptor{
		field: ignoreCaseFieldDescriptor{
			fieldDescriptor: fieldDescriptor{
				schema: schema,
				table:  userTable,
				name:   "email",
			},
			fieldNameSuffix: "_lower",
		},
	},
	domain.UserHumanFieldEmailVerified: fieldDescriptor{
		schema: schema,
		table:  userTable,
		name:   "email_is_verified",
	},
}

type textClause[V domain.Text] struct {
	clause[domain.TextOperation]
	value V
}

var textOp map[domain.TextOperation]string = map[domain.TextOperation]string{
	domain.TextOperationEqual:                " = ",
	domain.TextOperationNotEqual:             " <> ",
	domain.TextOperationStartsWith:           " LIKE ",
	domain.TextOperationStartsWithIgnoreCase: " LIKE ",
}

func (tc textClause[V]) Write(stmt *statement) {
	placeholder := stmt.appendArg(tc.value)
	var (
		left, right string
	)
	switch tc.clause.op {
	case domain.TextOperationEqual:
		left = tc.clause.field.String()
		right = placeholder
	case domain.TextOperationNotEqual:
		left = tc.clause.field.String()
		right = placeholder
	case domain.TextOperationStartsWith:
		left = tc.clause.field.String()
		right = placeholder + "%"
	case domain.TextOperationStartsWithIgnoreCase:
		left = tc.clause.field.String()
		if _, ok := tc.clause.field.(ignoreCaseFieldDescriptor); !ok {
			left = "LOWER(" + left + ")"
		}
		right = "LOWER(" + placeholder + "%)"
	}

	stmt.builder.WriteString(left)
	stmt.builder.WriteString(textOp[tc.clause.op])
	stmt.builder.WriteString(right)
}

type boolClause[V domain.Bool] struct {
	clause[domain.BoolOperation]
	value V
}

func (bc boolClause[V]) Write(stmt *statement) {
	if !bc.value {
		stmt.builder.WriteString("NOT ")
	}
	stmt.builder.WriteString(bc.clause.field.String())
}

type numberClause[V domain.Number] struct {
	clause[domain.NumberOperation]
	value V
}

var numberOp map[domain.NumberOperation]string = map[domain.NumberOperation]string{
	domain.NumberOperationEqual:              " = ",
	domain.NumberOperationNotEqual:           " <> ",
	domain.NumberOperationLessThan:           " < ",
	domain.NumberOperationLessThanOrEqual:    " <= ",
	domain.NumberOperationGreaterThan:        " > ",
	domain.NumberOperationGreaterThanOrEqual: " >= ",
}

func (nc numberClause[V]) Write(stmt *statement) {
	stmt.builder.WriteString(nc.clause.field.String())
	stmt.builder.WriteString(numberOp[nc.clause.op])
	stmt.builder.WriteString(stmt.appendArg(nc.value))
}
