package view

import (
	"database/sql"
	"github.com/caos/zitadel/internal/id"

	"github.com/jinzhu/gorm"
)

type View struct {
	Db          *gorm.DB
	idGenerator id.Generator
}

func StartView(sqlClient *sql.DB, idGenerator id.Generator) (*View, error) {
	gorm, err := gorm.Open("postgres", sqlClient)
	if err != nil {
		return nil, err
	}
	return &View{
		Db:          gorm,
		idGenerator: idGenerator,
	}, nil
}

func (v *View) Health() (err error) {
	return v.Db.DB().Ping()
}
