package v3

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type User struct {
	instanceID string
	orgID      string
	id         string
	username   string

	createdAt time.Time
	updatedAt time.Time
	deletedAt time.Time
}

// Columns implements [object].
func (u User) Columns(table Table) []Column {
	return []Column{
		&column{
			table: table,
			name:  columnNameInstanceID,
		},
		&column{
			table: table,
			name:  columnNameOrgID,
		},
		&column{
			table: table,
			name:  columnNameID,
		},
		&columnIgnoreCase{
			column: column{
				table: table,
				name:  userTableUsernameColumn,
			},
			suffix: "_lower",
		},
		&column{
			table: table,
			name:  columnNameCreatedAt,
		},
		&column{
			table: table,
			name:  columnNameUpdatedAt,
		},
		&column{
			table: table,
			name:  columnNameDeletedAt,
		},
	}
}

// Scan implements [object].
func (u User) Scan(row database.Scanner) error {
	return row.Scan(
		&u.instanceID,
		&u.orgID,
		&u.id,
		&u.username,
		&u.createdAt,
		&u.updatedAt,
		&u.deletedAt,
	)
}

type userTable struct {
	*table
}

const (
	userTableUsernameColumn = "username"
)

func UserTable() *userTable {
	table := &userTable{
		table: newTable[User]("zitadel", "users"),
	}

	table.possibleJoins = func(table Table) map[string]Column {
		switch on := table.(type) {
		case *userTable:
			return map[string]Column{
				columnNameInstanceID: on.InstanceIDColumn(),
				columnNameOrgID:      on.OrgIDColumn(),
				columnNameID:         on.IDColumn(),
			}
		case *orgTable:
			return map[string]Column{
				columnNameInstanceID: on.InstanceIDColumn(),
				columnNameOrgID:      on.IDColumn(),
			}
		case *instanceTable:
			return map[string]Column{
				columnNameInstanceID: on.IDColumn(),
			}
		default:
			return nil
		}
	}

	return table
}

func (t *userTable) InstanceIDColumn() Column {
	return t.columns[columnNameInstanceID]
}

func (t *userTable) OrgIDColumn() Column {
	return t.columns[columnNameOrgID]
}

func (t *userTable) IDColumn() Column {
	return t.columns[columnNameID]
}

func (t *userTable) UsernameColumn() Column {
	return t.columns[userTableUsernameColumn]
}

func (t *userTable) CreatedAtColumn() Column {
	return t.columns[columnNameCreatedAt]
}

func (t *userTable) UpdatedAtColumn() Column {
	return t.columns[columnNameUpdatedAt]
}

func (t *userTable) DeletedAtColumn() Column {
	return t.columns[columnNameDeletedAt]
}

func NewUserQuery() Query[User] {
	q := NewQuery[User](UserTable())
	return q
}

type userByIDCondition[T Text] struct {
	id T
}

func UserByID[T Text](id T) Condition {
	return &userByIDCondition[T]{id: id}
}

// writeOn implements Condition.
func (u *userByIDCondition[T]) writeOn(builder statementBuilder) {
	NewTextCondition(builder.table().(*userTable).IDColumn(), TextOperatorEqual, u.id).writeOn(builder)
}

var _ Condition = (*userByIDCondition[string])(nil)

type userByUsernameCondition[T Text] struct {
	username T
	operator TextOperator
}

func UserByUsername[T Text](username T, operator TextOperator) Condition {
	return &userByUsernameCondition[T]{username: username, operator: operator}
}

// writeOn implements Condition.
func (u *userByUsernameCondition[T]) writeOn(builder statementBuilder) {
	NewTextCondition(builder.table().(*userTable).UsernameColumn(), u.operator, u.username).writeOn(builder)
}

var _ Condition = (*userByUsernameCondition[string])(nil)
