package view

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
)

type View struct {
	Db     *gorm.DB
	client *database.DB
}

func StartView(sqlClient *database.DB) (*View, error) {
	gorm, err := gorm.Open("postgres", sqlClient.DB)
	if err != nil {
		return nil, err
	}
	return &View{
		Db:     gorm,
		client: sqlClient,
	}, nil
}

func (v *View) Health() (err error) {
	return v.Db.DB().Ping()
}

func (v *View) TimeTravel(ctx context.Context, tableName string) string {
	return tableName + v.client.Timetravel(call.Took(ctx))
}
