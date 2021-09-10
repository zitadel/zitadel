package view

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

type View struct {
	Db *gorm.DB
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
