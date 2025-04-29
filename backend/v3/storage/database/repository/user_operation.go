package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userOperation struct {
	database.QueryExecutor
	clauses []domain.UserClause
}

// Delete implements [domain.UserOperation].
func (u *userOperation) Delete(ctx context.Context) error {
	return u.QueryExecutor.Exec(ctx, `DELETE FROM users WHERE id = $1`, u.clauses)
}

// SetUsername implements [domain.UserOperation].
func (u *userOperation) SetUsername(ctx context.Context, username string) error {
	var stmt statement

	stmt.builder.WriteString(`UPDATE users SET username = $1 WHERE `)
	stmt.appendArg(username)
	clausesToSQL(&stmt, u.clauses)
	return u.QueryExecutor.Exec(ctx, stmt.builder.String(), stmt.args...)
}

var _ domain.UserOperation = (*userOperation)(nil)

func UserIDQuery(id string) domain.UserClause {
	return textClause[string]{
		clause: clause[domain.TextOperation]{
			field: userFields[domain.UserFieldID],
			op:    domain.TextOperationEqual,
		},
		value: id,
	}
}

func HumanEmailQuery(op domain.TextOperation, email string) domain.UserClause {
	return textClause[string]{
		clause: clause[domain.TextOperation]{
			field: userFields[domain.UserHumanFieldEmail],
			op:    op,
		},
		value: email,
	}
}

func HumanEmailVerifiedQuery(op domain.BoolOperation) domain.UserClause {
	return boolClause[domain.BoolOperation]{
		clause: clause[domain.BoolOperation]{
			field: userFields[domain.UserHumanFieldEmailVerified],
			op:    op,
		},
	}
}

func clausesToSQL(stmt *statement, clauses []domain.UserClause) {
	for _, clause := range clauses {

		stmt.builder.WriteString(userFields[clause.Field()].String())
		stmt.builder.WriteString(clause.Operation().String())
		stmt.appendArg(clause.Args()...)
	}
}
