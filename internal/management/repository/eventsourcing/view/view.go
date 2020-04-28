package view

import (
	"database/sql"
	"github.com/caos/zitadel/internal/view"
	"github.com/jinzhu/gorm"
)

type View struct {
	Db *gorm.DB
}

func StartView(conf view.ViewConfig) (*View, *sql.DB, error) {
	viewDB, err := view.Start(conf)
	if err != nil {
		return nil, nil, err
	}
	return &View{
		Db: viewDB.GORM,
	}, viewDB.SQL, nil
}

func (v *View) Health() (err error) {
	return v.Db.DB().Ping()
}
