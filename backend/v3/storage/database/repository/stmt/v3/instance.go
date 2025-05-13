package v3

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Instance struct {
	id   string
	name string

	createdAt time.Time
	updatedAt time.Time
	deletedAt time.Time
}

// Columns implements [object].
func (Instance) Columns(table Table) []Column {
	return []Column{
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
func (i Instance) Scan(row database.Scanner) error {
	return row.Scan(
		&i.id,
		&i.name,
		&i.createdAt,
		&i.updatedAt,
		&i.deletedAt,
	)
}

type instanceTable struct {
	*table
}

func InstanceTable() *instanceTable {
	table := &instanceTable{
		table: newTable[Instance]("zitadel", "instances"),
	}

	table.possibleJoins = func(t Table) map[string]Column {
		switch on := t.(type) {
		case *instanceTable:
			return map[string]Column{
				columnNameID: on.IDColumn(),
			}
		case *orgTable:
			return map[string]Column{
				columnNameID: on.InstanceIDColumn(),
			}
		case *userTable:
			return map[string]Column{
				columnNameID: on.InstanceIDColumn(),
			}
		default:
			return nil
		}
	}

	return table
}

func (i *instanceTable) IDColumn() Column {
	return i.columns[columnNameID]
}

func (i *instanceTable) NameColumn() Column {
	return i.columns[columnNameName]
}

func (i *instanceTable) CreatedAtColumn() Column {
	return i.columns[columnNameCreatedAt]
}

func (i *instanceTable) UpdatedAtColumn() Column {
	return i.columns[columnNameUpdatedAt]
}

func (i *instanceTable) DeletedAtColumn() Column {
	return i.columns[columnNameDeletedAt]
}
