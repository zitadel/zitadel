package view

import (
	"github.com/jinzhu/gorm"
	"github.com/zitadel/zitadel/internal/database"
)

type View struct {
	Db *gorm.DB
}

func StartView(sqlClient *database.DB) (*View, error) {
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
