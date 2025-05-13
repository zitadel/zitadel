package v3

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Org struct {
	instanceID string
	id         string

	name string

	createdAt time.Time
	updatedAt time.Time
	deletedAt time.Time
}

// Columns implements [object].
func (Org) Columns(table Table) []Column {
	return []Column{
		&column{
			table: table,
			name:  columnNameInstanceID,
		},
		&column{
			table: table,
			name:  columnNameID,
		},
		&column{
			table: table,
			name:  columnNameName,
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
func (o Org) Scan(row database.Scanner) error {
	return row.Scan(
		&o.instanceID,
		&o.id,
		&o.name,
		&o.createdAt,
		&o.updatedAt,
		&o.deletedAt,
	)
}

type orgTable struct {
	*table
}

func OrgTable() *orgTable {
	table := &orgTable{
		table: newTable[Org]("zitadel", "orgs"),
	}

	table.possibleJoins = func(table Table) map[string]Column {
		switch on := table.(type) {
		case *instanceTable:
			return map[string]Column{
				columnNameInstanceID: on.IDColumn(),
			}
		case *orgTable:
			return map[string]Column{
				columnNameInstanceID: on.InstanceIDColumn(),
				columnNameID:         on.IDColumn(),
			}
		case *userTable:
			return map[string]Column{
				columnNameInstanceID: on.InstanceIDColumn(),
				columnNameID:         on.IDColumn(),
			}
		default:
			return nil
		}
	}

	return table
}

func (o *orgTable) InstanceIDColumn() Column {
	return o.columns[columnNameInstanceID]
}

func (o *orgTable) IDColumn() Column {
	return o.columns[columnNameID]
}

func (o *orgTable) NameColumn() Column {
	return o.columns[columnNameName]
}

func (o *orgTable) CreatedAtColumn() Column {
	return o.columns[columnNameCreatedAt]
}

func (o *orgTable) UpdatedAtColumn() Column {
	return o.columns[columnNameUpdatedAt]
}

func (o *orgTable) DeletedAtColumn() Column {
	return o.columns[columnNameDeletedAt]
}
