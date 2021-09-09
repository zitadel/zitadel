package view

import (
	"database/sql"

	"github.com/caos/zitadel/internal/query"
	"github.com/jinzhu/gorm"
)

type View struct {
	Db    *gorm.DB
	query *query.Queries
}

func StartView(sqlClient *sql.DB) (*View, error) {
	gorm, err := gorm.Open("postgres", sqlClient)
	if err != nil {
		return nil, err
	}
	return &View{
		Db: gorm,
	}, nil
}

func (v *View) Health() (err error) {
	return v.Db.DB().Ping()
}
