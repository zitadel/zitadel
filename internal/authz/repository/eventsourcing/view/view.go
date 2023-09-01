package view

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/query"
)

type View struct {
	Db          *gorm.DB
	Query       *query.Queries
	idGenerator id.Generator
	client      *database.DB
}

func StartView(sqlClient *database.DB, idGenerator id.Generator, queries *query.Queries) (*View, error) {
	gorm, err := gorm.Open("postgres", sqlClient.DB)
	if err != nil {
		return nil, err
	}
	return &View{
		Db:          gorm,
		idGenerator: idGenerator,
		Query:       queries,
		client:      sqlClient,
	}, nil
}

func (v *View) Health() (err error) {
	return v.Db.DB().Ping()
}

func (v *View) TimeTravel(ctx context.Context, tableName string) string {
	return tableName + v.client.Timetravel(call.Took(ctx))
}
