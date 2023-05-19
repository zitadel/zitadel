package eventstore

import (
	"database/sql"
)

type Eventstore struct {
	client *sql.DB
}
