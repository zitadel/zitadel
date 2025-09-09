package view

import (
	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/query"
)

type View struct {
	Db     *gorm.DB
	client *database.DB
	Query  *query.Queries
}

func StartView(sqlClient *database.DB, queries *query.Queries) (*View, error) {
	gorm, err := gorm.Open("postgres", sqlClient.DB.(*sql.Pool).DB)
	if err != nil {
		return nil, err
	}
	return &View{
		Db:     gorm,
		Query:  queries,
		client: sqlClient,
	}, nil
}

func (v *View) Health() (err error) {
	return v.Db.DB().Ping()
}
