package repository

import (
	"database/sql"

	es_stor "github.com/caos/eventstore-lib/pkg/storage"
	"github.com/jinzhu/gorm"

	//sql import
	_ "github.com/lib/pq"
	// postgres is for gorm dialect defintion
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type SQL struct {
	address string
	dialect string
	client  *gorm.DB

	sqlClient *sql.DB
}

func (db *SQL) Start(options ...es_stor.Option) (err error) {
	db.sqlClient, err = sql.Open("postgres", db.address)
	if err != nil {
		return err
	}

	db.client, err = gorm.Open("postgres", db.sqlClient)
	return err
}

func (db *SQL) Health() error {
	return db.client.DB().Ping()
}
